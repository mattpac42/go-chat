package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

func TestChatService_ProcessMessage_Success(t *testing.T) {
	// Create a mock Claude server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		events := []string{
			`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_123","role":"assistant"}}` + "\n\n",
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello! "}}` + "\n\n",
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"How can I help?"}}` + "\n\n",
			`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
		}

		for _, event := range events {
			w.Write([]byte(event))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}))
	defer server.Close()

	// Setup dependencies
	logger := zerolog.Nop()
	repo := repository.NewMockProjectRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	// Create a project
	ctx := context.Background()
	project, err := repo.Create(ctx, "Test Project")
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	// Create chat service
	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, repo, logger)

	// Collect chunks
	var chunks []string
	var mu sync.Mutex

	onChunk := func(chunk string) {
		mu.Lock()
		chunks = append(chunks, chunk)
		mu.Unlock()
	}

	// Process message
	result, err := chatService.ProcessMessage(ctx, project.ID, "Hi there!", onChunk)
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify result
	if result.Role != model.RoleAssistant {
		t.Errorf("expected role 'assistant', got '%s'", result.Role)
	}

	if result.Content != "Hello! How can I help?" {
		t.Errorf("expected content 'Hello! How can I help?', got '%s'", result.Content)
	}

	// Verify chunks were received
	if len(chunks) != 2 {
		t.Errorf("expected 2 chunks, got %d", len(chunks))
	}

	// Verify messages were saved to repository
	messages, err := repo.GetMessages(ctx, project.ID)
	if err != nil {
		t.Fatalf("failed to get messages: %v", err)
	}

	// Should have 2 messages: user message and assistant response
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}

	if messages[0].Role != model.RoleUser || messages[0].Content != "Hi there!" {
		t.Errorf("unexpected user message: %+v", messages[0])
	}

	if messages[1].Role != model.RoleAssistant || messages[1].Content != "Hello! How can I help?" {
		t.Errorf("unexpected assistant message: %+v", messages[1])
	}
}

func TestChatService_ProcessMessage_WithContext(t *testing.T) {
	// Track messages sent to Claude
	var sentMessages []ClaudeMessage

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req claudeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}
		sentMessages = req.Messages

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		events := []string{
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Response"}}` + "\n\n",
			`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
		}

		for _, event := range events {
			w.Write([]byte(event))
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	repo := repository.NewMockProjectRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	ctx := context.Background()
	project, _ := repo.Create(ctx, "Test Project")

	// Add some existing messages
	repo.CreateMessage(ctx, project.ID, model.RoleUser, "First message")
	repo.CreateMessage(ctx, project.ID, model.RoleAssistant, "First response")

	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, repo, logger)

	// Process a new message
	_, err := chatService.ProcessMessage(ctx, project.ID, "Second message", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify context was included
	if len(sentMessages) != 3 {
		t.Fatalf("expected 3 messages sent to Claude, got %d", len(sentMessages))
	}

	if sentMessages[0].Content != "First message" {
		t.Errorf("expected first message 'First message', got '%s'", sentMessages[0].Content)
	}
	if sentMessages[1].Content != "First response" {
		t.Errorf("expected second message 'First response', got '%s'", sentMessages[1].Content)
	}
	if sentMessages[2].Content != "Second message" {
		t.Errorf("expected third message 'Second message', got '%s'", sentMessages[2].Content)
	}
}

func TestChatService_ProcessMessage_ContextLimit(t *testing.T) {
	var sentMessages []ClaudeMessage

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req claudeRequest
		json.NewDecoder(r.Body).Decode(&req)
		sentMessages = req.Messages

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"OK"}}` + "\n\n"))
		w.Write([]byte(`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n"))
	}))
	defer server.Close()

	logger := zerolog.Nop()
	repo := repository.NewMockProjectRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	ctx := context.Background()
	project, _ := repo.Create(ctx, "Test Project")

	// Add more messages than the limit
	for i := 0; i < 25; i++ {
		repo.CreateMessage(ctx, project.ID, model.RoleUser, "Message "+string(rune('A'+i)))
		repo.CreateMessage(ctx, project.ID, model.RoleAssistant, "Response "+string(rune('A'+i)))
	}

	// Set limit to 10
	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 10,
	}, claudeService, repo, logger)

	_, err := chatService.ProcessMessage(ctx, project.ID, "New message", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Should only include last 10 messages (the new message is already saved and included in context)
	if len(sentMessages) != 10 {
		t.Errorf("expected 10 messages (limited context including new message), got %d", len(sentMessages))
	}

	// The last message should be the new one
	if sentMessages[len(sentMessages)-1].Content != "New message" {
		t.Errorf("expected last message to be 'New message', got '%s'", sentMessages[len(sentMessages)-1].Content)
	}
}

func TestChatService_ProcessMessage_ProjectNotFound(t *testing.T) {
	logger := zerolog.Nop()
	repo := repository.NewMockProjectRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
	}, logger)

	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, repo, logger)

	ctx := context.Background()
	nonExistentID := uuid.New()

	_, err := chatService.ProcessMessage(ctx, nonExistentID, "Hello", func(string) {})
	if err == nil {
		t.Fatal("expected error for non-existent project")
	}
}

func TestChatService_ProcessMessage_CodeBlockExtraction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Response with a code block - using escaped newlines for valid JSON
		// The actual response will be: Here's some code:\n```python\nprint('hello')\n```\nThat's it!
		events := []string{
			"event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"Here's some code:\\n```python\\nprint('hello')\\n```\\nThat's it!\"}}\n\n",
			`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
		}

		for _, event := range events {
			w.Write([]byte(event))
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	repo := repository.NewMockProjectRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	ctx := context.Background()
	project, _ := repo.Create(ctx, "Test Project")

	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, repo, logger)

	result, err := chatService.ProcessMessage(ctx, project.ID, "Show me code", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify code blocks were extracted
	if len(result.CodeBlocks) != 1 {
		t.Fatalf("expected 1 code block, got %d", len(result.CodeBlocks))
	}

	if result.CodeBlocks[0].Language != "python" {
		t.Errorf("expected language 'python', got '%s'", result.CodeBlocks[0].Language)
	}

	if result.CodeBlocks[0].Code != "print('hello')" {
		t.Errorf("expected code \"print('hello')\", got '%s'", result.CodeBlocks[0].Code)
	}
}

func TestChatService_ProcessMessage_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := zerolog.Nop()
	repo := repository.NewMockProjectRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	project, _ := repo.Create(context.Background(), "Test Project")

	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, repo, logger)

	_, err := chatService.ProcessMessage(ctx, project.ID, "Hello", func(string) {})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
