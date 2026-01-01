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
	config              ChatConfig
	claudeService       ClaudeMessenger
	discoveryService    *DiscoveryService
	agentContextService *AgentContextService
	repo                repository.ProjectRepository
	fileRepo            repository.FileRepository
	fileMetadataRepo    repository.FileMetadataRepository
	logger              zerolog.Logger
}

// NewChatService creates a new chat service.
func NewChatService(
	config ChatConfig,
	claudeService ClaudeMessenger,
	discoveryService *DiscoveryService,
	agentContextService *AgentContextService,
	repo repository.ProjectRepository,
	fileRepo repository.FileRepository,
	fileMetadataRepo repository.FileMetadataRepository,
	logger zerolog.Logger,
) *ChatService {
	if config.ContextMessageLimit <= 0 {
		config.ContextMessageLimit = 20
	}
	return &ChatService{
		config:              config,
		claudeService:       claudeService,
		discoveryService:    discoveryService,
		agentContextService: agentContextService,
		repo:                repo,
		fileRepo:            fileRepo,
		fileMetadataRepo:    fileMetadataRepo,
		logger:              logger,
	}
}

// ChatResult contains the result of processing a chat message.
type ChatResult struct {
	Message    *model.Message
	Role       model.Role
	Content    string
	CodeBlocks []model.CodeBlock
	AgentType  *string // "product_manager", "designer", "developer", or nil
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

	// Determine agent context and system prompt
	var agentType *string
	var agentContext *model.AgentContext

	// Get agent context (only when not in discovery mode)
	if s.agentContextService != nil && (discovery == nil || discovery.Stage.IsComplete()) {
		var err error
		agentContext, err = s.agentContextService.GetContextForMessage(ctx, projectID, content)
		if err != nil {
			s.logger.Warn().
				Err(err).
				Str("projectId", projectID.String()).
				Msg("failed to get agent context, continuing without agent type")
		} else if agentContext != nil {
			agentTypeStr := string(agentContext.Agent)
			agentType = &agentTypeStr
			s.logger.Debug().
				Str("projectId", projectID.String()).
				Str("agentType", agentTypeStr).
				Msg("determined agent type for message")
		}
	}

	// Get appropriate system prompt (discovery-aware, agent-specific, or default)
	systemPrompt := s.getSystemPrompt(ctx, projectID, discovery, agentContext)

	s.logger.Debug().
		Str("projectId", projectID.String()).
		Int("contextMessages", len(claudeMessages)).
		Bool("discoveryMode", discovery != nil && !discovery.Stage.IsComplete()).
		Msg("sending message to Claude")

	// Send to Claude and handle tool use loop
	responseContent, err := s.processStreamWithTools(ctx, projectID, systemPrompt, claudeMessages, onChunk)
	if err != nil {
		return nil, err
	}

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

	// Save assistant response with agent type
	assistantMsg, err := s.repo.CreateMessageWithAgent(ctx, projectID, model.RoleAssistant, responseContent, agentType)
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
		AgentType:  agentType,
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

// getSystemPrompt returns the appropriate system prompt based on state.
// Priority: 1) Discovery prompt (if in discovery mode)
//           2) Agent-specific prompt (if agent context available)
//           3) Default code-generation prompt
func (s *ChatService) getSystemPrompt(ctx context.Context, projectID uuid.UUID, discovery *model.ProjectDiscovery, agentContext *model.AgentContext) string {
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

	// If agent context is available, use agent-specific prompt
	if agentContext != nil && s.agentContextService != nil {
		prompt, err := s.agentContextService.GetSystemPrompt(ctx, agentContext)
		if err != nil {
			s.logger.Warn().
				Err(err).
				Str("projectId", projectID.String()).
				Str("agent", string(agentContext.Agent)).
				Msg("failed to get agent prompt, using default")
		} else if prompt != "" {
			s.logger.Debug().
				Str("projectId", projectID.String()).
				Str("agent", string(agentContext.Agent)).
				Msg("using agent-specific system prompt")
			return prompt
		}
	}

	// Fall back to default prompt
	return DefaultSystemPrompt()
}

// processStreamWithTools handles streaming from Claude, executing tools, and continuing
// the conversation until Claude returns a final response (not a tool_use).
func (s *ChatService) processStreamWithTools(
	ctx context.Context,
	projectID uuid.UUID,
	systemPrompt string,
	claudeMessages []ClaudeMessage,
	onChunk func(chunk string),
) (string, error) {
	// Initial request
	stream, err := s.claudeService.SendMessage(ctx, systemPrompt, claudeMessages)
	if err != nil {
		return "", fmt.Errorf("failed to send message to Claude: %w", err)
	}

	var fullResponse strings.Builder
	maxIterations := 10 // Prevent infinite loops

	for iteration := 0; iteration < maxIterations; iteration++ {
		// Collect response while streaming
		var iterationResponse strings.Builder
		for chunk := range stream.Chunks() {
			iterationResponse.WriteString(chunk)
			fullResponse.WriteString(chunk)
			if onChunk != nil {
				onChunk(chunk)
			}
		}

		if err := stream.Err(); err != nil {
			stream.Close()
			return "", fmt.Errorf("stream error: %w", err)
		}

		// Check for tool uses
		toolUses := stream.ToolUses()
		stopReason := stream.StopReason()
		stream.Close()

		// If no tool uses or stop reason is not tool_use, we're done
		if len(toolUses) == 0 || stopReason != "tool_use" {
			break
		}

		s.logger.Debug().
			Int("toolCount", len(toolUses)).
			Str("projectId", projectID.String()).
			Msg("executing tools")

		// Execute tools and collect results
		var toolResults []ToolResult
		var assistantContent []ContentBlock

		// Add text content block if there was text
		if iterationResponse.Len() > 0 {
			assistantContent = append(assistantContent, ContentBlock{
				Type: "text",
				Text: iterationResponse.String(),
			})
		}

		// Add tool_use blocks and execute them
		for _, toolUse := range toolUses {
			assistantContent = append(assistantContent, ContentBlock{
				Type:  "tool_use",
				ID:    toolUse.ID,
				Name:  toolUse.Name,
				Input: toolUse.Input,
			})

			result := s.executeTool(ctx, projectID, toolUse)
			toolResults = append(toolResults, result)

			s.logger.Debug().
				Str("toolName", toolUse.Name).
				Str("toolID", toolUse.ID).
				Bool("isError", result.IsError).
				Msg("tool executed")
		}

		// Continue conversation with tool results
		stream, err = s.claudeService.SendMessageWithToolResults(
			ctx,
			systemPrompt,
			claudeMessages,
			assistantContent,
			toolResults,
		)
		if err != nil {
			return "", fmt.Errorf("failed to continue with tool results: %w", err)
		}
	}

	return fullResponse.String(), nil
}

// executeTool executes a single tool and returns the result.
func (s *ChatService) executeTool(ctx context.Context, projectID uuid.UUID, toolUse ToolUseBlock) ToolResult {
	result := ToolResult{
		Type:      "tool_result",
		ToolUseID: toolUse.ID,
	}

	switch toolUse.Name {
	case "write_file":
		path, _ := toolUse.Input["path"].(string)
		content, _ := toolUse.Input["content"].(string)

		if path == "" {
			result.Content = "Error: path is required"
			result.IsError = true
			return result
		}

		if s.fileRepo == nil {
			result.Content = "Error: file operations not available"
			result.IsError = true
			return result
		}

		// Infer language from file extension
		language := inferLanguageFromPath(path)

		_, err := s.fileRepo.SaveFile(ctx, projectID, path, language, content)
		if err != nil {
			result.Content = fmt.Sprintf("Error writing file: %v", err)
			result.IsError = true
			s.logger.Warn().
				Err(err).
				Str("path", path).
				Str("projectId", projectID.String()).
				Msg("failed to write file via tool")
		} else {
			result.Content = fmt.Sprintf("File written successfully: %s", path)
			s.logger.Info().
				Str("path", path).
				Str("projectId", projectID.String()).
				Msg("wrote file via tool")
		}

	case "read_file":
		path, _ := toolUse.Input["path"].(string)

		if path == "" {
			result.Content = "Error: path is required"
			result.IsError = true
			return result
		}

		if s.fileRepo == nil {
			result.Content = "Error: file operations not available"
			result.IsError = true
			return result
		}

		file, err := s.fileRepo.GetFileByPath(ctx, projectID, path)
		if err != nil {
			result.Content = fmt.Sprintf("Error reading file: %v", err)
			result.IsError = true
			s.logger.Warn().
				Err(err).
				Str("path", path).
				Str("projectId", projectID.String()).
				Msg("failed to read file via tool")
		} else {
			result.Content = file.Content
			s.logger.Debug().
				Str("path", path).
				Str("projectId", projectID.String()).
				Int("contentLength", len(file.Content)).
				Msg("read file via tool")
		}

	default:
		result.Content = fmt.Sprintf("Unknown tool: %s", toolUse.Name)
		result.IsError = true
	}

	return result
}

// inferLanguageFromPath determines the programming language from a file path.
func inferLanguageFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".html", ".htm":
		return "html"
	case ".css":
		return "css"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".tsx":
		return "tsx"
	case ".jsx":
		return "jsx"
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".md":
		return "markdown"
	case ".sh":
		return "bash"
	case ".sql":
		return "sql"
	case ".xml":
		return "xml"
	case ".java":
		return "java"
	case ".rb":
		return "ruby"
	case ".rs":
		return "rust"
	case ".c":
		return "c"
	case ".cpp", ".cc":
		return "cpp"
	case ".h", ".hpp":
		return "cpp"
	case ".php":
		return "php"
	default:
		return "text"
	}
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
