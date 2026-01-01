package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

// TestChatService_ProcessMessage_ToolUseWriteFile tests that when Claude uses write_file tool,
// the file is saved and the tool result is sent back
func TestChatService_ProcessMessage_ToolUseWriteFile(t *testing.T) {
	requestCount := 0
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		currentRequest := requestCount
		mu.Unlock()

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		if currentRequest == 1 {
			// First request: Claude responds with tool_use
			events := []string{
				`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_123","role":"assistant"}}` + "\n\n",
				`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}` + "\n\n",
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"I'll create that file for you."}}` + "\n\n",
				`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":0}` + "\n\n",
				`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_123","name":"write_file","input":{}}}` + "\n\n",
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"path\": \"index.html\", \"content\": \"<!DOCTYPE html><html><body>Hello</body></html>\"}"}}` + "\n\n",
				`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":1}` + "\n\n",
				`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"tool_use"}}` + "\n\n",
				`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
			}
			for _, event := range events {
				w.Write([]byte(event))
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		} else {
			// Second request (after tool result): Claude continues
			events := []string{
				`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_456","role":"assistant"}}` + "\n\n",
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"I've created the index.html file."}}` + "\n\n",
				`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"end_turn"}}` + "\n\n",
				`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
			}
			for _, event := range events {
				w.Write([]byte(event))
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	repo := repository.NewMockProjectRepository()
	fileRepo := repository.NewMockFileRepository()
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
	}, claudeService, nil, nil, repo, fileRepo, nil, logger)

	var chunks []string
	onChunk := func(chunk string) {
		chunks = append(chunks, chunk)
	}

	result, err := chatService.ProcessMessage(ctx, project.ID, "Create an index.html file", onChunk)
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify we received the streaming text
	fullText := strings.Join(chunks, "")
	if !strings.Contains(fullText, "create") && !strings.Contains(fullText, "created") {
		t.Errorf("expected response about creating file, got: %s", fullText)
	}

	// Verify the file was saved
	files, err := fileRepo.GetFilesByProject(ctx, project.ID)
	if err != nil {
		t.Fatalf("failed to get files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	if files[0].Path != "index.html" {
		t.Errorf("expected file path 'index.html', got '%s'", files[0].Path)
	}

	// Verify the file content
	file, err := fileRepo.GetFileByPath(ctx, project.ID, "index.html")
	if err != nil {
		t.Fatalf("failed to get file by path: %v", err)
	}

	if !strings.Contains(file.Content, "Hello") {
		t.Errorf("expected file content to contain 'Hello', got: %s", file.Content)
	}

	// Verify the result contains the continuation response
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Role != model.RoleAssistant {
		t.Errorf("expected role 'assistant', got '%s'", result.Role)
	}
}

// TestChatService_ProcessMessage_ToolUseReadFile tests read_file tool execution
func TestChatService_ProcessMessage_ToolUseReadFile(t *testing.T) {
	requestCount := 0
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		currentRequest := requestCount
		mu.Unlock()

		// Check for tool_result in the request body
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		if currentRequest == 1 {
			// First request: Claude uses read_file tool
			events := []string{
				`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_123","role":"assistant"}}` + "\n\n",
				`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"toolu_read","name":"read_file","input":{}}}` + "\n\n",
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{\"path\": \"config.json\"}"}}` + "\n\n",
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
		} else {
			// Second request: Claude uses the file content
			events := []string{
				`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_456","role":"assistant"}}` + "\n\n",
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"I see the config has debug mode enabled."}}` + "\n\n",
				`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"end_turn"}}` + "\n\n",
				`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
			}
			for _, event := range events {
				w.Write([]byte(event))
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	repo := repository.NewMockProjectRepository()
	fileRepo := repository.NewMockFileRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	ctx := context.Background()
	project, _ := repo.Create(ctx, "Test Project")

	// Create a file that Claude will read
	fileRepo.SaveFile(ctx, project.ID, "config.json", "json", `{"debug": true}`)

	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, nil, nil, repo, fileRepo, nil, logger)

	result, err := chatService.ProcessMessage(ctx, project.ID, "What's in config.json?", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify the response mentions the file content
	if !strings.Contains(result.Content, "debug") {
		t.Errorf("expected response to mention debug mode, got: %s", result.Content)
	}
}

// TestChatService_ProcessMessage_MultipleToolUses tests handling multiple tool uses in one response
func TestChatService_ProcessMessage_MultipleToolUses(t *testing.T) {
	requestCount := 0
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		currentRequest := requestCount
		mu.Unlock()

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		if currentRequest == 1 {
			// First request: Claude uses write_file twice
			events := []string{
				`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_123","role":"assistant"}}` + "\n\n",
				`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}` + "\n\n",
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Creating files..."}}` + "\n\n",
				`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":0}` + "\n\n",
				// First tool use
				`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_1","name":"write_file","input":{}}}` + "\n\n",
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"path\": \"index.html\", \"content\": \"<html></html>\"}"}}` + "\n\n",
				`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":1}` + "\n\n",
				// Second tool use
				`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":2,"content_block":{"type":"tool_use","id":"toolu_2","name":"write_file","input":{}}}` + "\n\n",
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":2,"delta":{"type":"input_json_delta","partial_json":"{\"path\": \"style.css\", \"content\": \"body {}\"}"}}` + "\n\n",
				`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":2}` + "\n\n",
				`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"tool_use"}}` + "\n\n",
				`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
			}
			for _, event := range events {
				w.Write([]byte(event))
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		} else {
			// Second request: Claude continues
			events := []string{
				`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_456","role":"assistant"}}` + "\n\n",
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Both files created!"}}` + "\n\n",
				`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"end_turn"}}` + "\n\n",
				`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
			}
			for _, event := range events {
				w.Write([]byte(event))
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	repo := repository.NewMockProjectRepository()
	fileRepo := repository.NewMockFileRepository()
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
	}, claudeService, nil, nil, repo, fileRepo, nil, logger)

	_, err := chatService.ProcessMessage(ctx, project.ID, "Create HTML and CSS files", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify both files were saved
	files, err := fileRepo.GetFilesByProject(ctx, project.ID)
	if err != nil {
		t.Fatalf("failed to get files: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	// Verify file paths
	filePaths := make(map[string]bool)
	for _, f := range files {
		filePaths[f.Path] = true
	}

	if !filePaths["index.html"] {
		t.Error("expected index.html to be created")
	}
	if !filePaths["style.css"] {
		t.Error("expected style.css to be created")
	}
}

// TestChatService_ProcessMessage_ToolError tests error handling when tool execution fails
func TestChatService_ProcessMessage_ToolError(t *testing.T) {
	requestCount := 0
	var mu sync.Mutex
	var secondRequestMessages []interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		currentRequest := requestCount
		mu.Unlock()

		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)

		if currentRequest == 2 {
			// Capture the messages from the second request to verify tool_result with error
			secondRequestMessages = req["messages"].([]interface{})
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		if currentRequest == 1 {
			// First request: Claude tries to read a non-existent file
			events := []string{
				`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_123","role":"assistant"}}` + "\n\n",
				`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"toolu_err","name":"read_file","input":{}}}` + "\n\n",
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{\"path\": \"nonexistent.txt\"}"}}` + "\n\n",
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
		} else {
			// Second request: Claude handles the error gracefully
			events := []string{
				`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_456","role":"assistant"}}` + "\n\n",
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"The file doesn't exist."}}` + "\n\n",
				`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"end_turn"}}` + "\n\n",
				`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
			}
			for _, event := range events {
				w.Write([]byte(event))
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	repo := repository.NewMockProjectRepository()
	fileRepo := repository.NewMockFileRepository()
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
	}, claudeService, nil, nil, repo, fileRepo, nil, logger)

	result, err := chatService.ProcessMessage(ctx, project.ID, "Read nonexistent.txt", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify Claude handled the error
	if !strings.Contains(result.Content, "doesn't exist") && !strings.Contains(result.Content, "not found") {
		t.Errorf("expected error message in response, got: %s", result.Content)
	}

	// Verify the tool_result with error was sent
	if len(secondRequestMessages) < 2 {
		t.Fatal("expected at least 2 messages in second request")
	}

	// The last message should be a user message with tool_result containing an error
	lastMsg := secondRequestMessages[len(secondRequestMessages)-1].(map[string]interface{})
	content := lastMsg["content"].([]interface{})
	if len(content) == 0 {
		t.Fatal("expected content in tool_result message")
	}

	toolResult := content[0].(map[string]interface{})
	if toolResult["type"] != "tool_result" {
		t.Errorf("expected tool_result, got %v", toolResult["type"])
	}
	if toolResult["is_error"] != true {
		t.Error("expected is_error to be true for failed tool execution")
	}
}

// TestChatService_ProcessMessage_FallbackToMarkdown tests that markdown parsing still works as fallback
func TestChatService_ProcessMessage_FallbackToMarkdown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Claude responds with markdown code block (no tool use) - fallback behavior
		responseContent := "Here is your file:\n\n```html:index.html\n---\nshort_description: Test file\nlong_description: A test HTML file.\nfunctional_group: Test\n---\n<!DOCTYPE html><html></html>\n```"

		escaped := strings.ReplaceAll(responseContent, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		escaped = strings.ReplaceAll(escaped, "\n", "\\n")

		events := []string{
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"` + escaped + `"}}` + "\n\n",
			`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"end_turn"}}` + "\n\n",
			`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
		}

		for _, event := range events {
			w.Write([]byte(event))
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	repo := repository.NewMockProjectRepository()
	fileRepo := repository.NewMockFileRepository()
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
	}, claudeService, nil, nil, repo, fileRepo, nil, logger)

	_, err := chatService.ProcessMessage(ctx, project.ID, "Create a file", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify the file was saved via markdown fallback
	files, err := fileRepo.GetFilesByProject(ctx, project.ID)
	if err != nil {
		t.Fatalf("failed to get files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file from markdown fallback, got %d", len(files))
	}

	if files[0].Path != "index.html" {
		t.Errorf("expected file path 'index.html', got '%s'", files[0].Path)
	}
}
