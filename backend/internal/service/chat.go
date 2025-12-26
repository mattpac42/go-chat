package service

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/pkg/markdown"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

// filenamePattern matches common filename patterns like index.html, script.js, etc.
var filenamePattern = regexp.MustCompile(`\b([a-zA-Z0-9_-]+\.(?:html|css|js|ts|tsx|jsx|go|py|rb|java|json|xml|yaml|yml|md|sh|sql))\b`)

// ChatConfig holds configuration for the chat service.
type ChatConfig struct {
	ContextMessageLimit int
}

// ChatService orchestrates chat interactions between WebSocket, Claude, and database.
type ChatService struct {
	config           ChatConfig
	claudeService    ClaudeMessenger
	discoveryService *DiscoveryService
	repo             repository.ProjectRepository
	fileRepo         repository.FileRepository
	fileMetadataRepo repository.FileMetadataRepository
	logger           zerolog.Logger
}

// NewChatService creates a new chat service.
func NewChatService(
	config ChatConfig,
	claudeService ClaudeMessenger,
	discoveryService *DiscoveryService,
	repo repository.ProjectRepository,
	fileRepo repository.FileRepository,
	fileMetadataRepo repository.FileMetadataRepository,
	logger zerolog.Logger,
) *ChatService {
	if config.ContextMessageLimit <= 0 {
		config.ContextMessageLimit = 20
	}
	return &ChatService{
		config:           config,
		claudeService:    claudeService,
		discoveryService: discoveryService,
		repo:             repo,
		fileRepo:         fileRepo,
		fileMetadataRepo: fileMetadataRepo,
		logger:           logger,
	}
}

// ChatResult contains the result of processing a chat message.
type ChatResult struct {
	Message    *model.Message
	Role       model.Role
	Content    string
	CodeBlocks []model.CodeBlock
}

// ProcessMessage handles a user message and streams the AI response.
// It saves both the user message and assistant response to the database.
// The onChunk callback is called for each streaming chunk received.
func (s *ChatService) ProcessMessage(
	ctx context.Context,
	projectID uuid.UUID,
	content string,
	onChunk func(chunk string),
) (*ChatResult, error) {
	// Verify project exists
	_, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Get or create discovery state for the project (if discovery service is configured)
	var discovery *model.ProjectDiscovery
	if s.discoveryService != nil {
		discovery, err = s.discoveryService.GetOrCreateDiscovery(ctx, projectID)
		if err != nil {
			s.logger.Warn().
				Err(err).
				Str("projectId", projectID.String()).
				Msg("failed to get discovery state, falling back to default mode")
			// Continue without discovery mode - don't fail the message
		}
	}

	// Save user message first
	_, err = s.repo.CreateMessage(ctx, projectID, model.RoleUser, content)
	if err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Get conversation history for context
	messages, err := s.repo.GetMessages(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Convert to Claude messages format with limit
	claudeMessages := s.buildClaudeMessages(messages)

	// Get appropriate system prompt (discovery-aware or default)
	systemPrompt := s.getSystemPrompt(ctx, projectID, discovery)

	s.logger.Debug().
		Str("projectId", projectID.String()).
		Int("contextMessages", len(claudeMessages)).
		Bool("discoveryMode", discovery != nil && !discovery.Stage.IsComplete()).
		Msg("sending message to Claude")

	// Send to Claude
	stream, err := s.claudeService.SendMessage(ctx, systemPrompt, claudeMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to send message to Claude: %w", err)
	}
	defer stream.Close()

	// Collect full response while streaming chunks
	var fullResponse strings.Builder
	for chunk := range stream.Chunks() {
		fullResponse.WriteString(chunk)
		if onChunk != nil {
			onChunk(chunk)
		}
	}

	if err := stream.Err(); err != nil {
		return nil, fmt.Errorf("stream error: %w", err)
	}

	responseContent := fullResponse.String()

	// If in discovery mode, extract and save discovery data from response
	if discovery != nil && !discovery.Stage.IsComplete() {
		if err := s.discoveryService.ExtractAndSaveData(ctx, discovery.ID, responseContent); err != nil {
			s.logger.Warn().
				Err(err).
				Str("projectId", projectID.String()).
				Str("discoveryId", discovery.ID.String()).
				Msg("failed to extract discovery data from response")
			// Continue - don't fail the message for discovery extraction errors
		}

		// Strip discovery metadata from response for display
		responseContent = StripMetadata(responseContent)
	}

	// Extract code blocks with metadata from response
	markdownBlocks := markdown.ExtractCodeBlocksWithMetadata(responseContent)

	// Try to infer filenames from user message if Claude didn't provide them
	markdownBlocks = inferFilenamesFromUserMessageWithMetadata(content, markdownBlocks)

	codeBlocks := make([]model.CodeBlock, len(markdownBlocks))
	for i, block := range markdownBlocks {
		codeBlocks[i] = model.CodeBlock{
			Language:   block.Language,
			Code:       block.Code,
			StartIndex: block.StartIndex,
			EndIndex:   block.EndIndex,
		}
		// Log each extracted code block for debugging
		hasMetadata := block.Metadata != nil
		s.logger.Info().
			Str("projectId", projectID.String()).
			Str("language", block.Language).
			Str("filename", block.Filename).
			Bool("hasMetadata", hasMetadata).
			Int("codeLength", len(block.Code)).
			Msg("extracted code block")
	}

	// Save files extracted from code blocks (only those with filenames)
	if s.fileRepo != nil {
		for _, block := range markdownBlocks {
			if block.Filename != "" {
				file, err := s.fileRepo.SaveFile(ctx, projectID, block.Filename, block.Language, block.Code)
				if err != nil {
					s.logger.Warn().
						Err(err).
						Str("projectId", projectID.String()).
						Str("filename", block.Filename).
						Msg("failed to save extracted file")
					continue
				}

				s.logger.Info().
					Str("projectId", projectID.String()).
					Str("filename", block.Filename).
					Msg("saved extracted file")

				// Save metadata if available and repository is configured
				if block.Metadata != nil && s.fileMetadataRepo != nil {
					_, err := s.fileMetadataRepo.Upsert(
						ctx,
						file.ID,
						block.Metadata.ShortDescription,
						block.Metadata.LongDescription,
						block.Metadata.FunctionalGroup,
					)
					if err != nil {
						s.logger.Warn().
							Err(err).
							Str("projectId", projectID.String()).
							Str("filename", block.Filename).
							Msg("failed to save file metadata")
					} else {
						s.logger.Info().
							Str("projectId", projectID.String()).
							Str("filename", block.Filename).
							Str("functionalGroup", block.Metadata.FunctionalGroup).
							Msg("saved file metadata")
					}
				}
			} else {
				s.logger.Info().
					Str("projectId", projectID.String()).
					Str("language", block.Language).
					Msg("code block has no filename - not saving")
			}
		}
	}

	// Save assistant response
	assistantMsg, err := s.repo.CreateMessage(ctx, projectID, model.RoleAssistant, responseContent)
	if err != nil {
		return nil, fmt.Errorf("failed to save assistant message: %w", err)
	}

	s.logger.Debug().
		Str("projectId", projectID.String()).
		Int("responseLength", len(responseContent)).
		Int("codeBlocks", len(codeBlocks)).
		Msg("completed message processing")

	return &ChatResult{
		Message:    assistantMsg,
		Role:       model.RoleAssistant,
		Content:    responseContent,
		CodeBlocks: codeBlocks,
	}, nil
}

// buildClaudeMessages converts database messages to Claude API format.
// It applies the context message limit, keeping the most recent messages.
func (s *ChatService) buildClaudeMessages(messages []model.Message) []ClaudeMessage {
	// Apply context limit - keep the most recent messages
	startIdx := 0
	if len(messages) > s.config.ContextMessageLimit {
		startIdx = len(messages) - s.config.ContextMessageLimit
	}

	claudeMessages := make([]ClaudeMessage, 0, len(messages)-startIdx)
	for _, msg := range messages[startIdx:] {
		claudeMessages = append(claudeMessages, ClaudeMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	return claudeMessages
}

// getSystemPrompt returns the appropriate system prompt based on discovery state.
// If the project is in discovery mode, returns the stage-specific discovery prompt.
// Otherwise, returns the default code-generation prompt.
func (s *ChatService) getSystemPrompt(ctx context.Context, projectID uuid.UUID, discovery *model.ProjectDiscovery) string {
	// If discovery service is configured and project is in discovery mode
	if s.discoveryService != nil && discovery != nil && !discovery.Stage.IsComplete() {
		prompt, err := s.discoveryService.GetSystemPrompt(ctx, projectID)
		if err != nil {
			s.logger.Warn().
				Err(err).
				Str("projectId", projectID.String()).
				Msg("failed to get discovery prompt, using default")
		} else if prompt != "" {
			s.logger.Debug().
				Str("projectId", projectID.String()).
				Str("stage", string(discovery.Stage)).
				Msg("using discovery system prompt")
			return prompt
		}
	}

	// Fall back to default prompt
	return DefaultSystemPrompt()
}

// inferFilenamesFromUserMessageWithMetadata extracts filenames mentioned in the user's message
// and assigns them to code blocks based on matching file extensions.
// This version works with CodeBlockWithMetadata which includes App Map metadata.
func inferFilenamesFromUserMessageWithMetadata(userMessage string, blocks []markdown.CodeBlockWithMetadata) []markdown.CodeBlockWithMetadata {
	// Find all filenames mentioned in the user message
	matches := filenamePattern.FindAllString(userMessage, -1)
	if len(matches) == 0 {
		return blocks
	}

	// Build a map of extension to filename
	extToFilename := make(map[string]string)
	for _, filename := range matches {
		ext := strings.ToLower(filepath.Ext(filename))
		extToFilename[ext] = filename
	}

	// Map extensions to languages
	extToLang := map[string][]string{
		".html": {"html", "htm"},
		".css":  {"css", "scss", "sass"},
		".js":   {"javascript", "js"},
		".ts":   {"typescript", "ts"},
		".tsx":  {"tsx"},
		".jsx":  {"jsx"},
		".go":   {"go", "golang"},
		".py":   {"python", "py"},
		".json": {"json"},
		".yaml": {"yaml", "yml"},
		".yml":  {"yaml", "yml"},
		".md":   {"markdown", "md"},
		".sh":   {"bash", "sh", "shell"},
		".sql":  {"sql"},
	}

	// Assign filenames to blocks that don't have one
	result := make([]markdown.CodeBlockWithMetadata, len(blocks))
	usedFilenames := make(map[string]bool)

	for i, block := range blocks {
		result[i] = block

		// Skip if block already has a filename
		if block.Filename != "" {
			usedFilenames[block.Filename] = true
			continue
		}

		// Try to find a matching filename based on language
		blockLang := strings.ToLower(block.Language)
		for ext, filename := range extToFilename {
			if usedFilenames[filename] {
				continue
			}
			if langs, ok := extToLang[ext]; ok {
				for _, lang := range langs {
					if blockLang == lang {
						result[i].Filename = filename
						usedFilenames[filename] = true
						break
					}
				}
			}
			if result[i].Filename != "" {
				break
			}
		}
	}

	return result
}
