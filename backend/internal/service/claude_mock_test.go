package service

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewMockClaudeService(t *testing.T) {
	fixturesPath := filepath.Join("..", "..", "testdata", "discovery")

	mock, err := NewMockClaudeService(fixturesPath)
	if err != nil {
		t.Fatalf("failed to create mock service: %v", err)
	}

	// Verify fixtures were loaded
	if fixture, ok := mock.GetFixture("welcome_response"); !ok {
		t.Error("expected welcome_response fixture to be loaded")
	} else if fixture.Stage != StageWelcome {
		t.Errorf("expected stage welcome, got %s", fixture.Stage)
	}
}

func TestMockClaudeServiceSimple(t *testing.T) {
	mock := NewMockClaudeServiceSimple()

	// Add a fixture programmatically
	mock.AddFixture("test_fixture", &DiscoveryFixture{
		Stage:    StageWelcome,
		Response: "Test response",
		Metadata: DiscoveryFixtureMetadata{
			StageComplete: false,
			NextStage:     "problem",
		},
	})

	fixture, ok := mock.GetFixture("test_fixture")
	if !ok {
		t.Error("expected test_fixture to be found")
	}
	if fixture.Response != "Test response" {
		t.Errorf("expected 'Test response', got '%s'", fixture.Response)
	}
}

func TestMockClaudeServiceSendMessage(t *testing.T) {
	mock := NewMockClaudeServiceSimple()

	// Add test fixtures
	mock.AddFixture("welcome_response", &DiscoveryFixture{
		Stage:    StageWelcome,
		Response: "Welcome! Tell me about yourself.",
		Metadata: DiscoveryFixtureMetadata{
			StageComplete: false,
			NextStage:     "problem",
		},
	})

	ctx := context.Background()
	messages := []ClaudeMessage{
		{Role: "user", Content: "Hello"},
	}

	stream, err := mock.SendMessage(ctx, "discovery flow", messages)
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	var response strings.Builder
	for chunk := range stream.Chunks() {
		response.WriteString(chunk)
	}

	if stream.Err() != nil {
		t.Errorf("stream error: %v", stream.Err())
	}

	responseStr := response.String()
	if !strings.Contains(responseStr, "Welcome") {
		t.Errorf("expected response to contain 'Welcome', got '%s'", responseStr)
	}
}

func TestMockClaudeServiceStageProgression(t *testing.T) {
	mock := NewMockClaudeServiceSimple()

	// Verify initial stage
	if mock.GetCurrentStage() != StageWelcome {
		t.Errorf("expected initial stage to be welcome, got %s", mock.GetCurrentStage())
	}

	// Set stage manually
	mock.SetCurrentStage(StageProblem)
	if mock.GetCurrentStage() != StageProblem {
		t.Errorf("expected stage to be problem, got %s", mock.GetCurrentStage())
	}

	// Reset and verify
	mock.Reset()
	if mock.GetCurrentStage() != StageWelcome {
		t.Errorf("expected stage to reset to welcome, got %s", mock.GetCurrentStage())
	}
	if mock.GetMessageCount() != 0 {
		t.Errorf("expected message count to reset to 0, got %d", mock.GetMessageCount())
	}
}

func TestMockClaudeServiceCustomHandler(t *testing.T) {
	mock := NewMockClaudeServiceSimple()

	customResponse := "Custom handler response"
	mock.SetCustomHandler(func(ctx context.Context, systemPrompt string, messages []ClaudeMessage) (*ClaudeStream, error) {
		stream := &ClaudeStream{
			chunks: make(chan string, 1),
			done:   make(chan struct{}),
		}
		go func() {
			stream.chunks <- customResponse
			close(stream.chunks)
			close(stream.done)
		}()
		return stream, nil
	})

	ctx := context.Background()
	messages := []ClaudeMessage{
		{Role: "user", Content: "Test"},
	}

	stream, err := mock.SendMessage(ctx, "test", messages)
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	var response strings.Builder
	for chunk := range stream.Chunks() {
		response.WriteString(chunk)
	}

	if response.String() != customResponse {
		t.Errorf("expected '%s', got '%s'", customResponse, response.String())
	}
}

func TestMockClaudeServiceValidation(t *testing.T) {
	mock := NewMockClaudeServiceSimple()
	ctx := context.Background()

	// Test with invalid message (empty content)
	messages := []ClaudeMessage{
		{Role: "user", Content: ""},
	}

	_, err := mock.SendMessage(ctx, "test", messages)
	if err == nil {
		t.Error("expected error for empty message content")
	}

	// Test with invalid role
	messages = []ClaudeMessage{
		{Role: "invalid", Content: "Test"},
	}

	_, err = mock.SendMessage(ctx, "test", messages)
	if err == nil {
		t.Error("expected error for invalid role")
	}
}

func TestDiscoverySummaryCard(t *testing.T) {
	mock := NewMockClaudeServiceSimple()

	summaryCard := &DiscoverySummaryCard{
		ProjectName:     "Test Project",
		SolvesStatement: "Solves the test problem",
		Users: []DiscoverySummaryUser{
			{Role: "Admin", Access: "full access"},
			{Role: "User", Access: "limited access"},
		},
		MVPFeatures: []string{"Feature 1", "Feature 2"},
		FutureFeatures: []DiscoverySummaryFutureFeature{
			{Name: "Future Feature", Version: "V2"},
		},
	}

	mock.AddFixture("summary_response", &DiscoveryFixture{
		Stage:       StageSummary,
		Response:    "Here's the summary",
		SummaryCard: summaryCard,
	})

	card, err := mock.GetSummaryCard("summary_response")
	if err != nil {
		t.Fatalf("failed to get summary card: %v", err)
	}

	if card.ProjectName != "Test Project" {
		t.Errorf("expected project name 'Test Project', got '%s'", card.ProjectName)
	}
	if len(card.Users) != 2 {
		t.Errorf("expected 2 users, got %d", len(card.Users))
	}
	if len(card.MVPFeatures) != 2 {
		t.Errorf("expected 2 MVP features, got %d", len(card.MVPFeatures))
	}
}

func TestIsConfirmationMessage(t *testing.T) {
	mock := NewMockClaudeServiceSimple()

	tests := []struct {
		content  string
		expected bool
	}{
		{"yes, let's do it!", true},
		{"Yep, sounds good", true},
		{"I want to change something", false},
		{"No, let me think", false},
		{"Perfect!", true},
		{"okay", true},
		{"Start building now", true},
		{"I have a question", false},
	}

	for _, tt := range tests {
		messages := []ClaudeMessage{
			{Role: "user", Content: tt.content},
		}
		result := mock.isConfirmationMessage(messages)
		if result != tt.expected {
			t.Errorf("isConfirmationMessage(%q) = %v, expected %v", tt.content, result, tt.expected)
		}
	}
}
