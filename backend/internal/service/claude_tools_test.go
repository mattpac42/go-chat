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

// TestClaudeService_ToolDefinitions verifies that tools are included in the request
func TestClaudeService_ToolDefinitions(t *testing.T) {
	var receivedRequest map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		events := []string{
			`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_123","role":"assistant"}}` + "\n\n",
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}` + "\n\n",
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
		{Role: "user", Content: "Create a file"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := service.SendMessage(ctx, "You are a helpful assistant", messages)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	defer stream.Close()

	// Drain the stream
	for range stream.Chunks() {
	}

	// Verify tools were included in the request
	tools, ok := receivedRequest["tools"].([]interface{})
	if !ok {
		t.Fatal("expected 'tools' array in request")
	}

	// Verify we have the expected tools
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolMap := tool.(map[string]interface{})
		name := toolMap["name"].(string)
		toolNames[name] = true
	}

	if !toolNames["write_file"] {
		t.Error("expected 'write_file' tool to be defined")
	}
	if !toolNames["read_file"] {
		t.Error("expected 'read_file' tool to be defined")
	}
}

// TestClaudeService_WriteFileToolSchema verifies the write_file tool schema
func TestClaudeService_WriteFileToolSchema(t *testing.T) {
	var receivedRequest map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedRequest)

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n"))
	}))
	defer server.Close()

	logger := zerolog.Nop()
	service := NewClaudeService(ClaudeConfig{
		APIKey:    "test-api-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	ctx := context.Background()
	stream, _ := service.SendMessage(ctx, "test", []ClaudeMessage{{Role: "user", Content: "test"}})
	defer stream.Close()
	for range stream.Chunks() {
	}

	tools := receivedRequest["tools"].([]interface{})
	var writeFileTool map[string]interface{}
	for _, tool := range tools {
		toolMap := tool.(map[string]interface{})
		if toolMap["name"] == "write_file" {
			writeFileTool = toolMap
			break
		}
	}

	if writeFileTool == nil {
		t.Fatal("write_file tool not found")
	}

	inputSchema := writeFileTool["input_schema"].(map[string]interface{})
	properties := inputSchema["properties"].(map[string]interface{})

	if _, ok := properties["path"]; !ok {
		t.Error("write_file tool missing 'path' property")
	}
	if _, ok := properties["content"]; !ok {
		t.Error("write_file tool missing 'content' property")
	}

	required := inputSchema["required"].([]interface{})
	requiredFields := make(map[string]bool)
	for _, r := range required {
		requiredFields[r.(string)] = true
	}

	if !requiredFields["path"] || !requiredFields["content"] {
		t.Error("write_file tool should require 'path' and 'content'")
	}
}

// TestClaudeStream_ToolUseBlock verifies parsing of tool_use content blocks
func TestClaudeStream_ToolUseBlock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Simulate Claude returning a tool_use block
		events := []string{
			`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_123","role":"assistant"}}` + "\n\n",
			`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"toolu_123","name":"write_file","input":{}}}` + "\n\n",
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{\"path\": \"index.html\", \"content\": \"<html>Hello</html>\"}"}}` + "\n\n",
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

	ctx := context.Background()
	stream, err := service.SendMessage(ctx, "test", []ClaudeMessage{{Role: "user", Content: "test"}})
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	defer stream.Close()

	// Drain text chunks
	for range stream.Chunks() {
	}

	// Verify tool use was captured
	toolUses := stream.ToolUses()
	if len(toolUses) != 1 {
		t.Fatalf("expected 1 tool use, got %d", len(toolUses))
	}

	if toolUses[0].Name != "write_file" {
		t.Errorf("expected tool name 'write_file', got '%s'", toolUses[0].Name)
	}
	if toolUses[0].ID != "toolu_123" {
		t.Errorf("expected tool ID 'toolu_123', got '%s'", toolUses[0].ID)
	}
	if toolUses[0].Input["path"] != "index.html" {
		t.Errorf("expected path 'index.html', got '%v'", toolUses[0].Input["path"])
	}
}

// TestClaudeStream_MixedTextAndToolUse verifies parsing when Claude returns both text and tool_use
func TestClaudeStream_MixedTextAndToolUse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Claude often returns text explaining what it's doing, then a tool use
		events := []string{
			`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_123","role":"assistant"}}` + "\n\n",
			// First content block: text
			`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}` + "\n\n",
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"I'll create the file for you."}}` + "\n\n",
			`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":0}` + "\n\n",
			// Second content block: tool_use
			`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_456","name":"write_file","input":{}}}` + "\n\n",
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"path\": \"style.css\", \"content\": \"body { color: red; }\"}"}}` + "\n\n",
			`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":1}` + "\n\n",
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

	ctx := context.Background()
	stream, err := service.SendMessage(ctx, "test", []ClaudeMessage{{Role: "user", Content: "test"}})
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	defer stream.Close()

	// Collect text chunks
	var textContent strings.Builder
	for chunk := range stream.Chunks() {
		textContent.WriteString(chunk)
	}

	// Verify both text and tool use were captured
	if textContent.String() != "I'll create the file for you." {
		t.Errorf("expected text 'I'll create the file for you.', got '%s'", textContent.String())
	}

	toolUses := stream.ToolUses()
	if len(toolUses) != 1 {
		t.Fatalf("expected 1 tool use, got %d", len(toolUses))
	}
	if toolUses[0].Name != "write_file" {
		t.Errorf("expected tool name 'write_file', got '%s'", toolUses[0].Name)
	}
}

// TestClaudeStream_StopReasonToolUse verifies we can detect when Claude stopped to use a tool
func TestClaudeStream_StopReasonToolUse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		events := []string{
			`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_123","role":"assistant"}}` + "\n\n",
			`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"toolu_789","name":"write_file","input":{}}}` + "\n\n",
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{\"path\": \"test.js\", \"content\": \"console.log('hi');\"}"}}` + "\n\n",
			`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":0}` + "\n\n",
			`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"tool_use"}}` + "\n\n",
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

	ctx := context.Background()
	stream, err := service.SendMessage(ctx, "test", []ClaudeMessage{{Role: "user", Content: "test"}})
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	defer stream.Close()

	// Drain the stream
	for range stream.Chunks() {
	}

	if stream.StopReason() != "tool_use" {
		t.Errorf("expected stop_reason 'tool_use', got '%s'", stream.StopReason())
	}
}

// TestToolUseBlock_Types verifies the ToolUseBlock struct
func TestToolUseBlock_Types(t *testing.T) {
	block := ToolUseBlock{
		Type:  "tool_use",
		ID:    "toolu_123",
		Name:  "write_file",
		Input: map[string]interface{}{
			"path":    "index.html",
			"content": "<html>Hello</html>",
		},
	}

	if block.Type != "tool_use" {
		t.Errorf("expected type 'tool_use', got '%s'", block.Type)
	}
	if block.ID != "toolu_123" {
		t.Errorf("expected ID 'toolu_123', got '%s'", block.ID)
	}
	if block.Name != "write_file" {
		t.Errorf("expected name 'write_file', got '%s'", block.Name)
	}
	if block.Input["path"] != "index.html" {
		t.Errorf("expected path 'index.html', got '%v'", block.Input["path"])
	}
}
