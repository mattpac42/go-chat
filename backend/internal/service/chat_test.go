package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

	// Create chat service (nil for discoveryService, fileRepo and fileMetadataRepo in basic tests)
	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, nil, nil, repo, nil, nil, logger)

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
	}, claudeService, nil, nil, repo, nil, nil, logger)

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
	}, claudeService, nil, nil, repo, nil, nil, logger)

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
	}, claudeService, nil, nil, repo, nil, nil, logger)

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
	}, claudeService, nil, nil, repo, nil, nil, logger)

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
	}, claudeService, nil, nil, repo, nil, nil, logger)

	_, err := chatService.ProcessMessage(ctx, project.ID, "Hello", func(string) {})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestChatService_ProcessMessage_FileMetadataExtraction(t *testing.T) {
	// Create a mock Claude server that returns files with YAML front matter metadata
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Response with code blocks containing YAML front matter metadata
		responseContent := `Here is your homepage:

` + "```html:index.html" + `
---
short_description: Main landing page structure
long_description: The homepage HTML with navigation and hero section.
functional_group: Homepage
---
<!DOCTYPE html>
<html>
<head><title>Home</title></head>
<body><h1>Welcome</h1></body>
</html>
` + "```" + `

And here is the styling:

` + "```css:styles.css" + `
---
short_description: Homepage styles
long_description: CSS rules for the homepage layout and typography.
functional_group: Homepage
---
body { margin: 0; padding: 0; }
h1 { color: blue; }
` + "```"

		// Properly escape for JSON
		escaped := strings.ReplaceAll(responseContent, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		escaped = strings.ReplaceAll(escaped, "\n", "\\n")

		events := []string{
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"` + escaped + `"}}` + "\n\n",
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
	fileMetadataRepo := repository.NewMockFileMetadataRepository()
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
	}, claudeService, nil, nil, repo, fileRepo, fileMetadataRepo, logger)

	_, err := chatService.ProcessMessage(ctx, project.ID, "Create a homepage", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify files were saved
	files, err := fileRepo.GetFilesByProject(ctx, project.ID)
	if err != nil {
		t.Fatalf("failed to get files: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	// Verify metadata was saved for the files
	// Get the actual file IDs to query metadata
	htmlFile, err := fileRepo.GetFileByPath(ctx, project.ID, "index.html")
	if err != nil {
		t.Fatalf("failed to get index.html: %v", err)
	}

	htmlMetadata, err := fileMetadataRepo.GetByFileID(ctx, htmlFile.ID)
	if err != nil {
		t.Fatalf("failed to get metadata for index.html: %v", err)
	}

	if htmlMetadata.ShortDescription != "Main landing page structure" {
		t.Errorf("expected short_description 'Main landing page structure', got '%s'", htmlMetadata.ShortDescription)
	}
	if htmlMetadata.FunctionalGroup != "Homepage" {
		t.Errorf("expected functional_group 'Homepage', got '%s'", htmlMetadata.FunctionalGroup)
	}

	cssFile, err := fileRepo.GetFileByPath(ctx, project.ID, "styles.css")
	if err != nil {
		t.Fatalf("failed to get styles.css: %v", err)
	}

	cssMetadata, err := fileMetadataRepo.GetByFileID(ctx, cssFile.ID)
	if err != nil {
		t.Fatalf("failed to get metadata for styles.css: %v", err)
	}

	if cssMetadata.ShortDescription != "Homepage styles" {
		t.Errorf("expected short_description 'Homepage styles', got '%s'", cssMetadata.ShortDescription)
	}
	if cssMetadata.FunctionalGroup != "Homepage" {
		t.Errorf("expected functional_group 'Homepage', got '%s'", cssMetadata.FunctionalGroup)
	}
}

func TestChatService_ProcessMessage_FileWithoutMetadata(t *testing.T) {
	// Create a mock Claude server that returns files WITHOUT metadata
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Response with code blocks without YAML front matter
		responseContent := `Here is your script:

` + "```javascript:app.js" + `
console.log('Hello World');
` + "```"

		escaped := strings.ReplaceAll(responseContent, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		escaped = strings.ReplaceAll(escaped, "\n", "\\n")

		events := []string{
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"` + escaped + `"}}` + "\n\n",
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
	fileMetadataRepo := repository.NewMockFileMetadataRepository()
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
	}, claudeService, nil, nil, repo, fileRepo, fileMetadataRepo, logger)

	_, err := chatService.ProcessMessage(ctx, project.ID, "Create a script", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify file was saved
	jsFile, err := fileRepo.GetFileByPath(ctx, project.ID, "app.js")
	if err != nil {
		t.Fatalf("failed to get app.js: %v", err)
	}

	// Verify the file content doesn't include the metadata (should just be the code)
	expectedContent := "console.log('Hello World');"
	if jsFile.Content != expectedContent {
		t.Errorf("expected content '%s', got '%s'", expectedContent, jsFile.Content)
	}

	// Verify NO metadata was saved for this file (since it had none)
	_, err = fileMetadataRepo.GetByFileID(ctx, jsFile.ID)
	if err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound for file without metadata, got: %v", err)
	}
}

func TestChatService_ProcessMessage_DiscoveryMode(t *testing.T) {
	// Track the system prompt sent to Claude
	var sentSystemPrompt string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req claudeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}
		sentSystemPrompt = req.System

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Response with discovery metadata
		responseContent := `Welcome! I'm here to help you turn your idea into a working application. First, tell me a bit about yourself - what do you do?<!--DISCOVERY_DATA:{"stage_complete":false,"extracted":{}}-->`

		escaped := strings.ReplaceAll(responseContent, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		escaped = strings.ReplaceAll(escaped, "\n", "\\n")

		events := []string{
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"` + escaped + `"}}` + "\n\n",
			`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
		}

		for _, event := range events {
			w.Write([]byte(event))
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	projectRepo := repository.NewMockProjectRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	discoveryService := NewDiscoveryService(discoveryRepo, nil, logger)

	ctx := context.Background()
	project, _ := projectRepo.Create(ctx, "Test Project")

	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, discoveryService, nil, projectRepo, nil, nil, logger)

	result, err := chatService.ProcessMessage(ctx, project.ID, "Hello", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify discovery system prompt was used (should contain "Product Guide")
	if !strings.Contains(sentSystemPrompt, "Product Guide") {
		t.Error("expected discovery system prompt to be used")
	}

	// Verify metadata was stripped from response
	if strings.Contains(result.Content, "DISCOVERY_DATA") {
		t.Error("expected discovery metadata to be stripped from response")
	}

	// Verify the clean response content
	expectedContent := "Welcome! I'm here to help you turn your idea into a working application. First, tell me a bit about yourself - what do you do?"
	if result.Content != expectedContent {
		t.Errorf("expected content '%s', got '%s'", expectedContent, result.Content)
	}
}

func TestChatService_ProcessMessage_DiscoveryModeAdvancesStage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Response with stage_complete:true and extracted data
		responseContent := `Great! Running a bakery sounds wonderful. What's your biggest challenge?<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"business_context":"Runs a local bakery"}}-->`

		escaped := strings.ReplaceAll(responseContent, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		escaped = strings.ReplaceAll(escaped, "\n", "\\n")

		events := []string{
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"` + escaped + `"}}` + "\n\n",
			`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
		}

		for _, event := range events {
			w.Write([]byte(event))
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	projectRepo := repository.NewMockProjectRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	discoveryService := NewDiscoveryService(discoveryRepo, nil, logger)

	ctx := context.Background()
	project, _ := projectRepo.Create(ctx, "Test Project")

	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, discoveryService, nil, projectRepo, nil, nil, logger)

	// Process message - should create discovery in welcome stage and advance to problem stage
	_, err := chatService.ProcessMessage(ctx, project.ID, "I run a bakery", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify discovery was created and stage advanced
	discovery, err := discoveryRepo.GetByProjectID(ctx, project.ID)
	if err != nil {
		t.Fatalf("failed to get discovery: %v", err)
	}

	// Stage should have advanced from welcome to problem
	if discovery.Stage != model.StageProblem {
		t.Errorf("expected stage '%s', got '%s'", model.StageProblem, discovery.Stage)
	}

	// Verify business context was extracted
	if discovery.BusinessContext == nil || *discovery.BusinessContext != "Runs a local bakery" {
		t.Errorf("expected business context 'Runs a local bakery', got '%v'", discovery.BusinessContext)
	}
}

func TestChatService_ProcessMessage_WithoutDiscoveryService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req claudeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		// Verify default system prompt is used (should contain "Go Chat" and "FILE FORMAT REQUIREMENT")
		if !strings.Contains(req.System, "FILE FORMAT REQUIREMENT") {
			t.Error("expected default system prompt with file format requirements")
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		events := []string{
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello!"}}` + "\n\n",
			`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
		}

		for _, event := range events {
			w.Write([]byte(event))
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	projectRepo := repository.NewMockProjectRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	ctx := context.Background()
	project, _ := projectRepo.Create(ctx, "Test Project")

	// Create chat service WITHOUT discovery service (nil)
	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, nil, nil, projectRepo, nil, nil, logger)

	result, err := chatService.ProcessMessage(ctx, project.ID, "Hello", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	if result.Content != "Hello!" {
		t.Errorf("expected content 'Hello!', got '%s'", result.Content)
	}
}

func TestChatService_ProcessMessage_DiscoveryCompleteUsesDefaultPrompt(t *testing.T) {
	var sentSystemPrompt string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req claudeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}
		sentSystemPrompt = req.System

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		events := []string{
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Here is your code..."}}` + "\n\n",
			`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
		}

		for _, event := range events {
			w.Write([]byte(event))
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	projectRepo := repository.NewMockProjectRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	discoveryService := NewDiscoveryService(discoveryRepo, nil, logger)

	ctx := context.Background()
	project, _ := projectRepo.Create(ctx, "Test Project")

	// Create a discovery and mark it complete
	discovery, _ := discoveryRepo.Create(ctx, project.ID)
	discoveryRepo.MarkComplete(ctx, discovery.ID)

	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, discoveryService, nil, projectRepo, nil, nil, logger)

	_, err := chatService.ProcessMessage(ctx, project.ID, "Create a homepage", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify default system prompt was used (should contain file format requirements, not Product Guide)
	if strings.Contains(sentSystemPrompt, "Product Guide") {
		t.Error("expected default system prompt after discovery complete, got discovery prompt")
	}
	if !strings.Contains(sentSystemPrompt, "FILE FORMAT REQUIREMENT") {
		t.Error("expected default system prompt with FILE FORMAT REQUIREMENT")
	}
}

func TestChatService_ProcessMessage_AgentTypeTracking(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		events := []string{
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Here is your code..."}}` + "\n\n",
			`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
		}

		for _, event := range events {
			w.Write([]byte(event))
		}
	}))
	defer server.Close()

	logger := zerolog.Nop()
	projectRepo := repository.NewMockProjectRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	prdRepo := repository.NewMockPRDRepository()
	claudeService := NewClaudeService(ClaudeConfig{
		APIKey:    "test-key",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		BaseURL:   server.URL,
	}, logger)

	discoveryService := NewDiscoveryService(discoveryRepo, nil, logger)
	agentContextService := NewAgentContextService(prdRepo, projectRepo, discoveryRepo, logger)

	ctx := context.Background()
	project, _ := projectRepo.Create(ctx, "Test Project")

	// Create a discovery and mark it complete (so we use agent context, not discovery mode)
	discovery, _ := discoveryRepo.Create(ctx, project.ID)
	discoveryRepo.MarkComplete(ctx, discovery.ID)

	chatService := NewChatService(ChatConfig{
		ContextMessageLimit: 20,
	}, claudeService, discoveryService, agentContextService, projectRepo, nil, nil, logger)

	result, err := chatService.ProcessMessage(ctx, project.ID, "Build me a homepage", func(string) {})
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Verify agent type was set (should be developer for "build" type message)
	if result.AgentType == nil {
		t.Error("expected agent type to be set")
	} else if *result.AgentType != "developer" {
		t.Errorf("expected agent type 'developer', got '%s'", *result.AgentType)
	}

	// Verify agent type is persisted in the message
	if result.Message.AgentType == nil {
		t.Error("expected message agent type to be set")
	} else if *result.Message.AgentType != "developer" {
		t.Errorf("expected message agent type 'developer', got '%s'", *result.Message.AgentType)
	}

	// Verify message is stored with agent type
	messages, err := projectRepo.GetMessages(ctx, project.ID)
	if err != nil {
		t.Fatalf("failed to get messages: %v", err)
	}

	// Find the assistant message
	var assistantMsg *model.Message
	for i := range messages {
		if messages[i].Role == model.RoleAssistant {
			assistantMsg = &messages[i]
			break
		}
	}

	if assistantMsg == nil {
		t.Fatal("expected assistant message to be saved")
	}

	if assistantMsg.AgentType == nil {
		t.Error("expected stored message to have agent type")
	} else if *assistantMsg.AgentType != "developer" {
		t.Errorf("expected stored message agent type 'developer', got '%s'", *assistantMsg.AgentType)
	}
}
