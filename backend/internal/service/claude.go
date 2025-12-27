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

	"github.com/rs/zerolog"
)

// ClaudeMessenger is the interface for sending messages to Claude (real or mock).
type ClaudeMessenger interface {
	SendMessage(ctx context.Context, systemPrompt string, messages []ClaudeMessage) (*ClaudeStream, error)
}

// ClaudeVision is the interface for image analysis with Claude Vision.
type ClaudeVision interface {
	AnalyzeImage(ctx context.Context, imageData []byte, mimeType, prompt string) (string, error)
}

const (
	defaultBaseURL      = "https://api.anthropic.com/v1/messages"
	anthropicVersion    = "2023-06-01"
	defaultSystemPrompt = `You are Go Chat. You create files for users.

FILE FORMAT REQUIREMENT - READ CAREFULLY:
Every code block MUST include a filename after the language, separated by a colon.
Each file MUST start with YAML front matter containing metadata for the App Map.

CORRECT FORMAT (files will be saved with metadata):
` + "```" + `html:index.html
---
short_description: Main homepage structure with navigation and hero section
long_description: This HTML file defines the complete structure of the landing page including a responsive navigation bar, hero section with call-to-action, and footer.
functional_group: Homepage
---
<!DOCTYPE html>
<html>content</html>
` + "```" + `

` + "```" + `css:styles.css
---
short_description: Visual styling for the homepage
long_description: Contains all CSS rules for the homepage including responsive layout, typography, colors, and component styles.
functional_group: Homepage
---
body { }
` + "```" + `

` + "```" + `javascript:app.js
---
short_description: Homepage interactive functionality
long_description: Handles user interactions on the homepage including navigation menu toggle, form validation, and dynamic content loading.
functional_group: Homepage
---
code here
` + "```" + `

METADATA FIELDS (all required in YAML front matter):
- short_description: One sentence (max 100 chars) describing what this file does
- long_description: Detailed explanation of the file's purpose and contents (2-3 sentences)
- functional_group: The feature area this file belongs to (e.g., Homepage, Contact Form, Navigation, Backend Services, Authentication, Configuration, Utilities)

WRONG FORMAT (files will NOT be saved):
` + "```" + `html
content
` + "```" + `

The pattern is: ` + "```" + `LANGUAGE:FILENAME followed by YAML front matter

Always include both the language AND filename with a colon between them.
Always include the YAML front matter metadata block.
Without the filename, the code will not be saved to the project.

Be concise. Generate working code with complete metadata.`
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
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content_block,omitempty"`
	Delta *struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta,omitempty"`
	Message *struct {
		ID   string `json:"id"`
		Role string `json:"role"`
	} `json:"message,omitempty"`
}

// ClaudeStream represents a streaming response from Claude.
type ClaudeStream struct {
	chunks chan string
	err    error
	done   chan struct{}
	resp   *http.Response
}

// Chunks returns a channel of text chunks from the stream.
func (s *ClaudeStream) Chunks() <-chan string {
	return s.chunks
}

// Err returns any error that occurred during streaming.
func (s *ClaudeStream) Err() error {
	return s.err
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

	// Build request
	reqBody := claudeRequest{
		Model:     s.config.Model,
		MaxTokens: s.config.MaxTokens,
		System:    systemPrompt,
		Messages:  messages,
		Stream:    true,
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

// processStream reads SSE events from the response and sends text chunks.
func (s *ClaudeService) processStream(resp *http.Response, stream *ClaudeStream) {
	defer close(stream.chunks)
	defer close(stream.done)

	scanner := bufio.NewScanner(resp.Body)

	// Increase buffer size for large responses
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var currentEvent string

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

			// Extract text from content_block_delta events
			if event.Type == "content_block_delta" && event.Delta != nil && event.Delta.Type == "text_delta" {
				stream.chunks <- event.Delta.Text
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
