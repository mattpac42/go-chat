package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

// ClaudeMessenger is the interface for sending messages to Claude (real or mock).
type ClaudeMessenger interface {
	SendMessage(ctx context.Context, systemPrompt string, messages []ClaudeMessage) (*ClaudeStream, error)
	SendMessageWithToolResults(ctx context.Context, systemPrompt string, messages []ClaudeMessage, assistantContent []ContentBlock, toolResults []ToolResult) (*ClaudeStream, error)
}

// ContentBlock represents a content block in Claude's response (for re-sending in continuation).
type ContentBlock struct {
	Type  string                 `json:"type"`
	Text  string                 `json:"text,omitempty"`
	ID    string                 `json:"id,omitempty"`
	Name  string                 `json:"name,omitempty"`
	Input map[string]interface{} `json:"input,omitempty"`
}

// ToolResult represents a tool result to send back to Claude.
type ToolResult struct {
	Type      string `json:"type"`
	ToolUseID string `json:"tool_use_id"`
	Content   string `json:"content"`
	IsError   bool   `json:"is_error,omitempty"`
}

// ClaudeVision is the interface for image analysis with Claude Vision.
type ClaudeVision interface {
	AnalyzeImage(ctx context.Context, imageData []byte, mimeType, prompt string) (string, error)
}

const (
	defaultBaseURL      = "https://api.anthropic.com/v1/messages"
	anthropicVersion    = "2023-06-01"
	defaultSystemPrompt = `You are Go Chat. You create files for users.

IMPORTANT: When creating files, ALWAYS use the write_file tool. Do not output code blocks with filenames - use the tool instead.
When reading existing files, use the read_file tool.

For each file you create with write_file, provide a brief explanation of what the file does.

FALLBACK FORMAT (only if tools are unavailable):
If for any reason you cannot use tools, use this code block format:
` + "```" + `html:index.html
---
short_description: Main homepage structure with navigation and hero section
long_description: This HTML file defines the complete structure of the landing page including a responsive navigation bar, hero section with call-to-action, and footer.
functional_group: Homepage
---
<!DOCTYPE html>
<html>content</html>
` + "```" + `

Be concise. Generate working code.`
)

// ClaudeConfig holds configuration for the Claude service.
type ClaudeConfig struct {
	APIKey    string
	Model     string
	MaxTokens int
	BaseURL   string // Optional, for testing
}

// ClaudeMessage represents a message in the Claude conversation.
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Validate checks if the message is valid.
func (m ClaudeMessage) Validate() error {
	if m.Content == "" {
		return errors.New("message content cannot be empty")
	}
	if m.Role != "user" && m.Role != "assistant" {
		return errors.New("message role must be 'user' or 'assistant'")
	}
	return nil
}

// claudeRequest is the request body for the Claude API.
type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	System    string          `json:"system"`
	Messages  []ClaudeMessage `json:"messages"`
	Stream    bool            `json:"stream"`
	Tools     []ClaudeTool    `json:"tools,omitempty"`
}

// claudeRequestWithContent is the request body when messages include content arrays (for tool results).
type claudeRequestWithContent struct {
	Model     string                   `json:"model"`
	MaxTokens int                      `json:"max_tokens"`
	System    string                   `json:"system"`
	Messages  []map[string]interface{} `json:"messages"`
	Stream    bool                     `json:"stream"`
	Tools     []ClaudeTool             `json:"tools,omitempty"`
}

// ClaudeTool represents a tool definition for the Claude API.
type ClaudeTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

// ToolUseBlock represents a tool_use content block from Claude's response.
type ToolUseBlock struct {
	Type  string                 `json:"type"`
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

// getFileTools returns the tool definitions for file operations.
func getFileTools() []ClaudeTool {
	return []ClaudeTool{
		{
			Name:        "write_file",
			Description: "Create or overwrite a file in the project. Use this for ALL file creation.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "File path relative to project root (e.g., 'src/index.html', 'styles/main.css')",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Complete file content to write",
					},
				},
				"required": []string{"path", "content"},
			},
		},
		{
			Name:        "read_file",
			Description: "Read the contents of a file in the project",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "File path relative to project root",
					},
				},
				"required": []string{"path"},
			},
		},
	}
}

// claudeVisionRequest is the request body for the Claude Vision API.
type claudeVisionRequest struct {
	Model     string                `json:"model"`
	MaxTokens int                   `json:"max_tokens"`
	Messages  []claudeVisionMessage `json:"messages"`
}

// claudeVisionMessage represents a message with multimodal content.
type claudeVisionMessage struct {
	Role    string                       `json:"role"`
	Content []claudeVisionMessageContent `json:"content"`
}

// claudeVisionMessageContent represents content in a Vision message.
type claudeVisionMessageContent struct {
	Type   string                    `json:"type"`
	Text   string                    `json:"text,omitempty"`
	Source *claudeVisionImageSource  `json:"source,omitempty"`
}

// claudeVisionImageSource represents an image source for Vision API.
type claudeVisionImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

// claudeVisionResponse represents a non-streaming response from Claude.
type claudeVisionResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
}

// claudeErrorResponse represents an error from the Claude API.
type claudeErrorResponse struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// claudeStreamEvent represents a streaming event from the Claude API.
type claudeStreamEvent struct {
	Type         string `json:"type"`
	Index        int    `json:"index,omitempty"`
	ContentBlock *struct {
		Type  string `json:"type"`
		Text  string `json:"text,omitempty"`
		ID    string `json:"id,omitempty"`    // For tool_use blocks
		Name  string `json:"name,omitempty"`  // For tool_use blocks
		Input json.RawMessage `json:"input,omitempty"` // For tool_use blocks
	} `json:"content_block,omitempty"`
	Delta *struct {
		Type        string `json:"type"`
		Text        string `json:"text,omitempty"`
		PartialJSON string `json:"partial_json,omitempty"` // For tool input streaming
		StopReason  string `json:"stop_reason,omitempty"`  // For message_delta
	} `json:"delta,omitempty"`
	Message *struct {
		ID   string `json:"id"`
		Role string `json:"role"`
	} `json:"message,omitempty"`
}

// ClaudeStream represents a streaming response from Claude.
type ClaudeStream struct {
	chunks     chan string
	err        error
	done       chan struct{}
	resp       *http.Response
	toolUses   []ToolUseBlock
	stopReason string
	mu         sync.Mutex
}

// Chunks returns a channel of text chunks from the stream.
func (s *ClaudeStream) Chunks() <-chan string {
	return s.chunks
}

// Err returns any error that occurred during streaming.
func (s *ClaudeStream) Err() error {
	return s.err
}

// ToolUses returns any tool_use blocks received during streaming.
func (s *ClaudeStream) ToolUses() []ToolUseBlock {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.toolUses
}

// StopReason returns the stop reason from the message (e.g., "end_turn", "tool_use").
func (s *ClaudeStream) StopReason() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopReason
}

// Close closes the stream and releases resources.
func (s *ClaudeStream) Close() error {
	if s.resp != nil && s.resp.Body != nil {
		return s.resp.Body.Close()
	}
	return nil
}

// ClaudeService handles communication with the Claude API.
type ClaudeService struct {
	config ClaudeConfig
	client *http.Client
	logger zerolog.Logger
}

// NewClaudeService creates a new Claude service.
func NewClaudeService(config ClaudeConfig, logger zerolog.Logger) *ClaudeService {
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}
	return &ClaudeService{
		config: config,
		client: &http.Client{},
		logger: logger,
	}
}

// DefaultSystemPrompt returns the default system prompt for Go Chat.
func DefaultSystemPrompt() string {
	return defaultSystemPrompt
}

// SendMessage sends messages to Claude and returns a streaming response.
func (s *ClaudeService) SendMessage(ctx context.Context, systemPrompt string, messages []ClaudeMessage) (*ClaudeStream, error) {
	// Validate messages
	for _, msg := range messages {
		if err := msg.Validate(); err != nil {
			return nil, fmt.Errorf("invalid message: %w", err)
		}
	}

	// Build request with file tools
	reqBody := claudeRequest{
		Model:     s.config.Model,
		MaxTokens: s.config.MaxTokens,
		System:    systemPrompt,
		Messages:  messages,
		Stream:    true,
		Tools:     getFileTools(),
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.BaseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.config.APIKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	s.logger.Debug().
		Str("model", s.config.Model).
		Int("messageCount", len(messages)).
		Msg("sending request to Claude API")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var errResp claudeErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
			return nil, fmt.Errorf("Claude API error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("Claude API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Start streaming
	stream := &ClaudeStream{
		chunks: make(chan string, 100),
		done:   make(chan struct{}),
		resp:   resp,
	}

	go s.processStream(resp, stream)

	return stream, nil
}

// SendMessageWithToolResults sends messages to Claude including tool results from previous tool uses.
// This is used to continue a conversation after executing tools.
func (s *ClaudeService) SendMessageWithToolResults(
	ctx context.Context,
	systemPrompt string,
	messages []ClaudeMessage,
	assistantContent []ContentBlock,
	toolResults []ToolResult,
) (*ClaudeStream, error) {
	// Build messages array with content blocks for tool results
	var msgArray []map[string]interface{}

	// Add previous conversation messages (string content)
	for _, msg := range messages {
		msgArray = append(msgArray, map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	// Add assistant message with the tool_use content blocks
	if len(assistantContent) > 0 {
		contentArray := make([]map[string]interface{}, len(assistantContent))
		for i, block := range assistantContent {
			contentArray[i] = map[string]interface{}{
				"type": block.Type,
			}
			if block.Type == "text" {
				contentArray[i]["text"] = block.Text
			} else if block.Type == "tool_use" {
				contentArray[i]["id"] = block.ID
				contentArray[i]["name"] = block.Name
				contentArray[i]["input"] = block.Input
			}
		}
		msgArray = append(msgArray, map[string]interface{}{
			"role":    "assistant",
			"content": contentArray,
		})
	}

	// Add user message with tool results
	if len(toolResults) > 0 {
		resultArray := make([]map[string]interface{}, len(toolResults))
		for i, result := range toolResults {
			resultArray[i] = map[string]interface{}{
				"type":        "tool_result",
				"tool_use_id": result.ToolUseID,
				"content":     result.Content,
			}
			if result.IsError {
				resultArray[i]["is_error"] = true
			}
		}
		msgArray = append(msgArray, map[string]interface{}{
			"role":    "user",
			"content": resultArray,
		})
	}

	// Build request with file tools
	reqBody := claudeRequestWithContent{
		Model:     s.config.Model,
		MaxTokens: s.config.MaxTokens,
		System:    systemPrompt,
		Messages:  msgArray,
		Stream:    true,
		Tools:     getFileTools(),
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.BaseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.config.APIKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	s.logger.Debug().
		Str("model", s.config.Model).
		Int("messageCount", len(msgArray)).
		Int("toolResults", len(toolResults)).
		Msg("sending request to Claude API with tool results")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var errResp claudeErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
			return nil, fmt.Errorf("Claude API error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("Claude API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Start streaming
	stream := &ClaudeStream{
		chunks: make(chan string, 100),
		done:   make(chan struct{}),
		resp:   resp,
	}

	go s.processStream(resp, stream)

	return stream, nil
}

// processStream reads SSE events from the response and sends text chunks.
func (s *ClaudeService) processStream(resp *http.Response, stream *ClaudeStream) {
	defer close(stream.chunks)
	defer close(stream.done)

	scanner := bufio.NewScanner(resp.Body)

	// Increase buffer size for large responses
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var currentEvent string

	// Track current content blocks by index (for tool_use blocks that stream JSON)
	type contentBlockState struct {
		blockType   string
		toolID      string
		toolName    string
		partialJSON strings.Builder
	}
	contentBlocks := make(map[int]*contentBlockState)

	for scanner.Scan() {
		line := scanner.Text()

		// Parse SSE format
		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			if data == "[DONE]" {
				break
			}

			var event claudeStreamEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				s.logger.Warn().Err(err).Str("data", data).Msg("failed to parse SSE event")
				continue
			}

			// Handle content_block_start - initialize tracking for tool_use blocks
			if event.Type == "content_block_start" && event.ContentBlock != nil {
				state := &contentBlockState{
					blockType: event.ContentBlock.Type,
				}
				if event.ContentBlock.Type == "tool_use" {
					state.toolID = event.ContentBlock.ID
					state.toolName = event.ContentBlock.Name
					s.logger.Debug().
						Str("toolID", state.toolID).
						Str("toolName", state.toolName).
						Int("index", event.Index).
						Msg("started tool_use content block")
				}
				contentBlocks[event.Index] = state
			}

			// Extract text from content_block_delta events
			if event.Type == "content_block_delta" && event.Delta != nil {
				if event.Delta.Type == "text_delta" {
					stream.chunks <- event.Delta.Text
				} else if event.Delta.Type == "input_json_delta" {
					// Accumulate partial JSON for tool input
					if state, ok := contentBlocks[event.Index]; ok {
						state.partialJSON.WriteString(event.Delta.PartialJSON)
					}
				}
			}

			// Handle content_block_stop - finalize tool_use blocks
			if event.Type == "content_block_stop" {
				if state, ok := contentBlocks[event.Index]; ok && state.blockType == "tool_use" {
					// Parse the accumulated JSON input
					var input map[string]interface{}
					jsonStr := state.partialJSON.String()
					if jsonStr != "" {
						if err := json.Unmarshal([]byte(jsonStr), &input); err != nil {
							s.logger.Warn().
								Err(err).
								Str("json", jsonStr).
								Msg("failed to parse tool input JSON")
						}
					}

					// Add the tool use to the stream
					toolUse := ToolUseBlock{
						Type:  "tool_use",
						ID:    state.toolID,
						Name:  state.toolName,
						Input: input,
					}
					stream.mu.Lock()
					stream.toolUses = append(stream.toolUses, toolUse)
					stream.mu.Unlock()

					s.logger.Debug().
						Str("toolID", state.toolID).
						Str("toolName", state.toolName).
						Interface("input", input).
						Msg("completed tool_use content block")

					delete(contentBlocks, event.Index)
				}
			}

			// Handle message_delta for stop_reason
			if event.Type == "message_delta" && event.Delta != nil && event.Delta.StopReason != "" {
				stream.mu.Lock()
				stream.stopReason = event.Delta.StopReason
				stream.mu.Unlock()
			}

			// Log other event types for debugging
			if event.Type == "error" {
				s.logger.Error().Str("event", currentEvent).Msg("received error event from Claude")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		stream.err = fmt.Errorf("stream read error: %w", err)
	}
}

// AnalyzeImage sends an image to Claude Vision API and returns the analysis as text.
func (s *ClaudeService) AnalyzeImage(ctx context.Context, imageData []byte, mimeType, prompt string) (string, error) {
	// Encode image data to base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)

	// Build the multimodal request
	reqBody := claudeVisionRequest{
		Model:     s.config.Model,
		MaxTokens: s.config.MaxTokens,
		Messages: []claudeVisionMessage{
			{
				Role: "user",
				Content: []claudeVisionMessageContent{
					{
						Type: "image",
						Source: &claudeVisionImageSource{
							Type:      "base64",
							MediaType: mimeType,
							Data:      base64Data,
						},
					},
					{
						Type: "text",
						Text: prompt,
					},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal vision request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.BaseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create vision request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.config.APIKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	s.logger.Debug().
		Str("model", s.config.Model).
		Str("mimeType", mimeType).
		Int("imageSize", len(imageData)).
		Msg("sending vision request to Claude API")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send vision request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read vision response: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		var errResp claudeErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
			return "", fmt.Errorf("Claude Vision API error: %s", errResp.Error.Message)
		}
		return "", fmt.Errorf("Claude Vision API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var visionResp claudeVisionResponse
	if err := json.Unmarshal(body, &visionResp); err != nil {
		return "", fmt.Errorf("failed to parse vision response: %w", err)
	}

	// Extract text from the response
	var result strings.Builder
	for _, content := range visionResp.Content {
		if content.Type == "text" {
			result.WriteString(content.Text)
		}
	}

	return result.String(), nil
}
