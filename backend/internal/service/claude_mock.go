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

	return m.createMockStream(fixture.Response), nil
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

	// Determine fixture based on current stage and conversation context
	switch m.currentStage {
	case StageWelcome:
		return "welcome_response"

	case StageProblem:
		// Check if this is initial or follow-up
		if m.messageCount <= 2 {
			return "problem_response"
		}
		return "problem_followup_response"

	case StagePersonas:
		if m.messageCount <= 4 {
			return "personas_response"
		}
		return "personas_followup_response"

	case StageMVP:
		if m.messageCount <= 6 {
			return "mvp_response"
		} else if m.messageCount <= 8 {
			return "mvp_followup_response"
		}
		return "mvp_confirm_response"

	case StageSummary:
		if m.isConfirmationMessage(messages) {
			return "complete_response"
		}
		if m.messageCount <= 10 {
			return "summary_response"
		}
		return "summary_confirm_response"

	case StageComplete:
		return "complete_response"

	default:
		return fmt.Sprintf("%s_response", m.currentStage)
	}
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

		// Simulate streaming by sending response in chunks
		words := strings.Fields(response)
		for i, word := range words {
			if i > 0 {
				stream.chunks <- " "
			}
			stream.chunks <- word

			// Small delay to simulate streaming (optional, can be removed for faster tests)
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
