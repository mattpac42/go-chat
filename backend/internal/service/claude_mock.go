package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// DiscoveryStage represents a stage in the discovery flow.
type DiscoveryStage string

const (
	StageWelcome  DiscoveryStage = "welcome"
	StageProblem  DiscoveryStage = "problem"
	StagePersonas DiscoveryStage = "personas"
	StageMVP      DiscoveryStage = "mvp"
	StageSummary  DiscoveryStage = "summary"
	StageComplete DiscoveryStage = "complete"
)

// DiscoveryFixtureMetadata contains metadata about a fixture response.
type DiscoveryFixtureMetadata struct {
	StageComplete        bool                   `json:"stage_complete"`
	NextStage            string                 `json:"next_stage,omitempty"`
	AwaitingConfirmation bool                   `json:"awaiting_confirmation,omitempty"`
	DiscoveryComplete    bool                   `json:"discovery_complete,omitempty"`
	Extracted            map[string]interface{} `json:"extracted,omitempty"`
}

// DiscoverySummaryUser represents a user role in the summary card.
type DiscoverySummaryUser struct {
	Role   string `json:"role"`
	Access string `json:"access"`
}

// DiscoverySummaryFutureFeature represents a future feature in the summary card.
type DiscoverySummaryFutureFeature struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// DiscoverySummaryCard contains the structured summary data for the discovery flow.
type DiscoverySummaryCard struct {
	ProjectName     string                          `json:"project_name"`
	SolvesStatement string                          `json:"solves_statement"`
	Users           []DiscoverySummaryUser          `json:"users"`
	MVPFeatures     []string                        `json:"mvp_features"`
	FutureFeatures  []DiscoverySummaryFutureFeature `json:"future_features"`
}

// DiscoveryFixture represents a fixture response for testing.
type DiscoveryFixture struct {
	Stage       DiscoveryStage           `json:"stage"`
	Response    string                   `json:"response"`
	Metadata    DiscoveryFixtureMetadata `json:"metadata"`
	SummaryCard *DiscoverySummaryCard    `json:"summary_card,omitempty"`
}

// MockClaudeService implements a mock Claude service for testing discovery flow.
type MockClaudeService struct {
	mu            sync.RWMutex
	fixtures      map[string]*DiscoveryFixture
	fixturesPath  string
	messageCount  int                          // Tracks messages to determine which fixture to return
	currentStage  DiscoveryStage               // Current discovery stage
	stageHistory  []DiscoveryStage             // History of stages
	customHandler func(ctx context.Context, systemPrompt string, messages []ClaudeMessage) (*ClaudeStream, error)
}

// NewMockClaudeService creates a new mock Claude service with fixtures loaded from the given path.
func NewMockClaudeService(fixturesPath string) (*MockClaudeService, error) {
	m := &MockClaudeService{
		fixtures:     make(map[string]*DiscoveryFixture),
		fixturesPath: fixturesPath,
		currentStage: StageWelcome,
		stageHistory: []DiscoveryStage{},
	}

	if err := m.loadFixtures(); err != nil {
		return nil, fmt.Errorf("failed to load fixtures: %w", err)
	}

	return m, nil
}

// NewMockClaudeServiceSimple creates a mock Claude service without loading fixtures.
// Useful for tests that provide custom handlers or inline fixtures.
func NewMockClaudeServiceSimple() *MockClaudeService {
	return &MockClaudeService{
		fixtures:     make(map[string]*DiscoveryFixture),
		currentStage: StageWelcome,
		stageHistory: []DiscoveryStage{},
	}
}

// loadFixtures loads all fixture JSON files from the fixtures path.
func (m *MockClaudeService) loadFixtures() error {
	entries, err := os.ReadDir(m.fixturesPath)
	if err != nil {
		return fmt.Errorf("failed to read fixtures directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(m.fixturesPath, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read fixture file %s: %w", entry.Name(), err)
		}

		var fixture DiscoveryFixture
		if err := json.Unmarshal(data, &fixture); err != nil {
			return fmt.Errorf("failed to parse fixture file %s: %w", entry.Name(), err)
		}

		// Use filename without extension as the key
		key := strings.TrimSuffix(entry.Name(), ".json")
		m.fixtures[key] = &fixture
	}

	return nil
}

// AddFixture adds or replaces a fixture programmatically.
func (m *MockClaudeService) AddFixture(key string, fixture *DiscoveryFixture) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fixtures[key] = fixture
}

// GetFixture retrieves a fixture by key.
func (m *MockClaudeService) GetFixture(key string) (*DiscoveryFixture, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	fixture, ok := m.fixtures[key]
	return fixture, ok
}

// SetCustomHandler sets a custom handler for processing messages.
// When set, this handler is used instead of fixture-based responses.
func (m *MockClaudeService) SetCustomHandler(handler func(ctx context.Context, systemPrompt string, messages []ClaudeMessage) (*ClaudeStream, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.customHandler = handler
}

// SetCurrentStage sets the current discovery stage.
func (m *MockClaudeService) SetCurrentStage(stage DiscoveryStage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentStage = stage
}

// GetCurrentStage returns the current discovery stage.
func (m *MockClaudeService) GetCurrentStage() DiscoveryStage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentStage
}

// Reset resets the mock service state.
func (m *MockClaudeService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messageCount = 0
	m.currentStage = StageWelcome
	m.stageHistory = []DiscoveryStage{}
}

// SendMessage implements the same interface as ClaudeService.
// It returns fixture responses based on the current discovery stage.
func (m *MockClaudeService) SendMessage(ctx context.Context, systemPrompt string, messages []ClaudeMessage) (*ClaudeStream, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Use custom handler if set
	if m.customHandler != nil {
		return m.customHandler(ctx, systemPrompt, messages)
	}

	// Validate messages
	for _, msg := range messages {
		if err := msg.Validate(); err != nil {
			return nil, fmt.Errorf("invalid message: %w", err)
		}
	}

	// Increment message count
	m.messageCount++

	// Determine which fixture to use based on current stage and message count
	fixtureKey := m.determineFixtureKey(systemPrompt, messages)

	fixture, ok := m.fixtures[fixtureKey]
	if !ok {
		// Fall back to a default response if fixture not found
		return m.createMockStream(fmt.Sprintf("Mock response for stage %s (fixture %s not found)", m.currentStage, fixtureKey)), nil
	}

	// Update stage history and current stage
	m.stageHistory = append(m.stageHistory, m.currentStage)
	if fixture.Metadata.NextStage != "" && fixture.Metadata.StageComplete {
		m.currentStage = DiscoveryStage(fixture.Metadata.NextStage)
	}

	// Build response with metadata comment (mimics real Claude response format)
	response := fixture.Response
	if fixture.Metadata.StageComplete || len(fixture.Metadata.Extracted) > 0 {
		metadataJSON, err := json.Marshal(fixture.Metadata)
		if err == nil {
			response = response + "\n\n<!--DISCOVERY_DATA:" + string(metadataJSON) + "-->"
		}
	}

	return m.createMockStream(response), nil
}

// determineFixtureKey determines which fixture to use based on context.
func (m *MockClaudeService) determineFixtureKey(systemPrompt string, messages []ClaudeMessage) string {
	// Check if system prompt indicates discovery mode
	isDiscovery := strings.Contains(strings.ToLower(systemPrompt), "discovery") ||
		strings.Contains(strings.ToLower(systemPrompt), "product guide")

	if !isDiscovery {
		// Return a generic mock response for non-discovery flows
		return "generic_response"
	}

	// Detect stage from system prompt - this is more reliable than internal state
	// because the mock service is a singleton and may have stale state from previous projects
	promptLower := strings.ToLower(systemPrompt)
	detectedStage := m.currentStage
	if strings.Contains(promptLower, "current stage: welcome") {
		detectedStage = StageWelcome
		// Only reset if we're coming from a later stage (stale state from previous project)
		// Don't reset if we're already at welcome stage to avoid duplicate welcome messages
		if m.currentStage != StageWelcome {
			m.currentStage = StageWelcome
			m.messageCount = 0
		}
	} else if strings.Contains(promptLower, "current stage: problem") {
		detectedStage = StageProblem
		m.currentStage = StageProblem
	} else if strings.Contains(promptLower, "current stage: personas") {
		detectedStage = StagePersonas
		m.currentStage = StagePersonas
	} else if strings.Contains(promptLower, "current stage: mvp") {
		detectedStage = StageMVP
		m.currentStage = StageMVP
	} else if strings.Contains(promptLower, "current stage: summary") {
		detectedStage = StageSummary
		m.currentStage = StageSummary
	}

	// Check if conversation already has assistant messages (welcome already sent)
	hasAssistantMessages := m.hasAssistantMessages(messages)

	// Count assistant messages in current stage to determine fixture progression
	stageResponseCount := m.countAssistantMessagesInStage(messages, detectedStage)

	// Determine fixture based on detected stage and how many responses we've given in this stage
	switch detectedStage {
	case StageWelcome:
		// Only return welcome if no assistant messages exist yet
		if !hasAssistantMessages {
			return "welcome_response"
		}
		// After welcome, user's response should trigger problem stage questions
		return "problem_response"

	case StageProblem:
		// First response in problem stage, or no problem responses yet
		if stageResponseCount == 0 {
			return "problem_response"
		}
		return "problem_followup_response"

	case StagePersonas:
		if stageResponseCount == 0 {
			return "personas_response"
		}
		return "personas_followup_response"

	case StageMVP:
		if stageResponseCount == 0 {
			return "mvp_response"
		} else if stageResponseCount == 1 {
			return "mvp_followup_response"
		}
		return "mvp_confirm_response"

	case StageSummary:
		if m.isConfirmationMessage(messages) {
			return "complete_response"
		}
		// First response in summary shows the summary
		if stageResponseCount == 0 {
			return "summary_response"
		}
		// Follow-up responses ask for confirmation
		return "summary_confirm_response"

	case StageComplete:
		return "complete_response"

	default:
		return fmt.Sprintf("%s_response", m.currentStage)
	}
}

// hasAssistantMessages checks if the conversation already has assistant messages.
func (m *MockClaudeService) hasAssistantMessages(messages []ClaudeMessage) bool {
	for _, msg := range messages {
		if msg.Role == "assistant" {
			return true
		}
	}
	return false
}

// countUserMessages counts the number of user messages in the conversation.
func (m *MockClaudeService) countUserMessages(messages []ClaudeMessage) int {
	count := 0
	for _, msg := range messages {
		if msg.Role == "user" {
			count++
		}
	}
	return count
}

// countAssistantMessagesInStage counts assistant messages that belong to the current stage.
// This looks at the content to determine which stage a response belongs to.
func (m *MockClaudeService) countAssistantMessagesInStage(messages []ClaudeMessage, stage DiscoveryStage) int {
	count := 0
	for _, msg := range messages {
		if msg.Role == "assistant" {
			// Check if this message belongs to the current stage based on content patterns
			content := strings.ToLower(msg.Content)
			switch stage {
			case StageWelcome:
				if strings.Contains(content, "welcome") || strings.Contains(content, "tell me a bit about yourself") {
					count++
				}
			case StageProblem:
				if strings.Contains(content, "challenge") || strings.Contains(content, "headache") ||
					strings.Contains(content, "perfect day") {
					count++
				}
			case StagePersonas:
				if strings.Contains(content, "who will actually use") || strings.Contains(content, "who else needs access") ||
					strings.Contains(content, "three people total") || strings.Contains(content, "employees") {
					count++
				}
			case StageMVP:
				if strings.Contains(content, "features") || strings.Contains(content, "priorities") ||
					strings.Contains(content, "essential") || strings.Contains(content, "three things") ||
					strings.Contains(content, "version 1") || strings.Contains(content, "version 2") {
					count++
				}
			case StageSummary:
				if strings.Contains(content, "here's what we're going to build") ||
					strings.Contains(content, "does this capture") {
					count++
				}
			}
		}
	}
	return count
}

// isConfirmationMessage checks if the last user message is a confirmation.
func (m *MockClaudeService) isConfirmationMessage(messages []ClaudeMessage) bool {
	if len(messages) == 0 {
		return false
	}

	lastMsg := messages[len(messages)-1]
	if lastMsg.Role != "user" {
		return false
	}

	content := strings.ToLower(lastMsg.Content)
	confirmationPhrases := []string{
		"yes", "yep", "yeah", "sure", "ok", "okay", "let's do it",
		"start building", "looks good", "that's right", "correct",
		"perfect", "sounds good", "confirmed", "go ahead",
	}

	for _, phrase := range confirmationPhrases {
		if strings.Contains(content, phrase) {
			return true
		}
	}

	return false
}

// createMockStream creates a mock ClaudeStream that returns the given response.
func (m *MockClaudeService) createMockStream(response string) *ClaudeStream {
	stream := &ClaudeStream{
		chunks: make(chan string, 100),
		done:   make(chan struct{}),
	}

	go func() {
		defer close(stream.chunks)
		defer close(stream.done)

		// Stream response preserving formatting (newlines, etc.)
		// Split into chunks of ~20 chars to simulate streaming while preserving structure
		chunkSize := 20
		for i := 0; i < len(response); i += chunkSize {
			end := i + chunkSize
			if end > len(response) {
				end = len(response)
			}
			stream.chunks <- response[i:end]

			// Small delay to simulate streaming
			time.Sleep(5 * time.Millisecond)
		}
	}()

	return stream
}

// CreateMockStreamFromFixture creates a mock stream from a specific fixture key.
func (m *MockClaudeService) CreateMockStreamFromFixture(fixtureKey string) (*ClaudeStream, error) {
	m.mu.RLock()
	fixture, ok := m.fixtures[fixtureKey]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("fixture not found: %s", fixtureKey)
	}

	return m.createMockStream(fixture.Response), nil
}

// GetSummaryCard returns the summary card from a fixture if available.
func (m *MockClaudeService) GetSummaryCard(fixtureKey string) (*DiscoverySummaryCard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fixture, ok := m.fixtures[fixtureKey]
	if !ok {
		return nil, fmt.Errorf("fixture not found: %s", fixtureKey)
	}

	if fixture.SummaryCard == nil {
		return nil, fmt.Errorf("fixture %s does not have a summary card", fixtureKey)
	}

	return fixture.SummaryCard, nil
}

// GetMessageCount returns the number of messages processed.
func (m *MockClaudeService) GetMessageCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.messageCount
}

// GetStageHistory returns the history of discovery stages.
func (m *MockClaudeService) GetStageHistory() []DiscoveryStage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	history := make([]DiscoveryStage, len(m.stageHistory))
	copy(history, m.stageHistory)
	return history
}

// SendMessageWithToolResults implements ClaudeMessenger interface for the mock service.
// It handles tool result continuation in the mock service.
func (m *MockClaudeService) SendMessageWithToolResults(
	ctx context.Context,
	systemPrompt string,
	messages []ClaudeMessage,
	assistantContent []ContentBlock,
	toolResults []ToolResult,
) (*ClaudeStream, error) {
	// For mock purposes, just return a simple continuation response
	// Real implementation would continue the conversation
	return m.createMockStream("I've processed the tool results and completed the operation."), nil
}

// AnalyzeImage implements ClaudeVision interface for the mock service.
// It returns a mock description of the image for testing purposes.
func (m *MockClaudeService) AnalyzeImage(ctx context.Context, imageData []byte, mimeType, prompt string) (string, error) {
	// Return a mock markdown description for testing
	return fmt.Sprintf(`## Image Description

This is a mock image analysis response for testing purposes.

### Image Details
- **Type**: %s
- **Size**: %d bytes

### Content
The image contains visual content that has been analyzed by Claude Vision.

### Transcribed Text
*No text was detected in this mock analysis.*

---
*This is a mock response from MockClaudeService.AnalyzeImage*
`, mimeType, len(imageData)), nil
}
