package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestClaudeService_SendMessage_Success(t *testing.T) {
	// Create a mock server that returns a streaming response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if r.Header.Get("x-api-key") != "test-api-key" {
			t.Errorf("expected api key header")
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Errorf("expected anthropic-version header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		if reqBody["model"] != "claude-sonnet-4-20250514" {
			t.Errorf("expected model claude-sonnet-4-20250514, got %v", reqBody["model"])
		}
		if reqBody["stream"] != true {
			t.Errorf("expected stream to be true")
		}

		// Return streaming SSE response
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		events := []string{
			`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_123","role":"assistant"}}` + "\n\n",
			`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}` + "\n\n",
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}` + "\n\n",
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" World"}}` + "\n\n",
			`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":0}` + "\n\n",
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

	logger := zerolog.Nop()
	service := NewClaudeService(ClaudeConfig{
		APIKey:    "test-api-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	messages := []ClaudeMessage{
		{Role: "user", Content: "Hello"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := service.SendMessage(ctx, "You are a helpful assistant", messages)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	defer stream.Close()

	var chunks []string
	for chunk := range stream.Chunks() {
		chunks = append(chunks, chunk)
	}

	if err := stream.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	fullContent := strings.Join(chunks, "")
	if fullContent != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", fullContent)
	}
}

func TestClaudeService_SendMessage_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"type":"error","error":{"type":"invalid_request_error","message":"Invalid API key"}}`))
	}))
	defer server.Close()

	logger := zerolog.Nop()
	service := NewClaudeService(ClaudeConfig{
		APIKey:    "invalid-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	messages := []ClaudeMessage{
		{Role: "user", Content: "Hello"},
	}

	ctx := context.Background()
	_, err := service.SendMessage(ctx, "You are a helpful assistant", messages)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "Invalid API key") {
		t.Errorf("expected error to contain 'Invalid API key', got: %v", err)
	}
}

func TestClaudeService_SendMessage_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := zerolog.Nop()
	service := NewClaudeService(ClaudeConfig{
		APIKey:    "test-api-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	messages := []ClaudeMessage{
		{Role: "user", Content: "Hello"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := service.SendMessage(ctx, "You are a helpful assistant", messages)

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestClaudeMessage_Validation(t *testing.T) {
	tests := []struct {
		name    string
		message ClaudeMessage
		valid   bool
	}{
		{
			name:    "valid user message",
			message: ClaudeMessage{Role: "user", Content: "Hello"},
			valid:   true,
		},
		{
			name:    "valid assistant message",
			message: ClaudeMessage{Role: "assistant", Content: "Hi there"},
			valid:   true,
		},
		{
			name:    "empty content",
			message: ClaudeMessage{Role: "user", Content: ""},
			valid:   false,
		},
		{
			name:    "invalid role",
			message: ClaudeMessage{Role: "system", Content: "Hello"},
			valid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()
			if tt.valid && err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}
