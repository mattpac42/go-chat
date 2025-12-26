package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// DiscoveryStage represents a stage in the discovery flow.
type DiscoveryStage string

const (
	DiscoveryStageWelcome  DiscoveryStage = "welcome"
	DiscoveryStageProblem  DiscoveryStage = "problem"
	DiscoveryStagePersonas DiscoveryStage = "personas"
	DiscoveryStageMVP      DiscoveryStage = "mvp"
	DiscoveryStageSummary  DiscoveryStage = "summary"
	DiscoveryStageComplete DiscoveryStage = "complete"
)

// DiscoveryUsers represents user information captured during discovery.
type DiscoveryUsers struct {
	Description     string `json:"description"`
	Count           int    `json:"count"`
	HasPermissions  bool   `json:"has_permissions"`
	PermissionNotes string `json:"permission_notes,omitempty"`
}

// DiscoveryMVPFeature represents an MVP feature with priority.
type DiscoveryMVPFeature struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

// DiscoveryFutureFeature represents a feature planned for future versions.
type DiscoveryFutureFeature struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// DiscoverySummary represents the final summary of the discovery session.
type DiscoverySummary struct {
	ProjectName     string    `json:"project_name"`
	SolvesStatement string    `json:"solves_statement"`
	ConfirmedAt     time.Time `json:"confirmed_at,omitempty"`
}

// DiscoveryState represents the complete state of a discovery session.
type DiscoveryState struct {
	ID               uuid.UUID                `json:"id"`
	ProjectID        uuid.UUID                `json:"project_id"`
	Stage            DiscoveryStage           `json:"stage"`
	StageStartedAt   time.Time                `json:"stage_started_at"`
	BusinessContext  string                   `json:"business_context,omitempty"`
	ProblemStatement string                   `json:"problem_statement,omitempty"`
	Goals            []string                 `json:"goals,omitempty"`
	Users            *DiscoveryUsers          `json:"users,omitempty"`
	MVPFeatures      []DiscoveryMVPFeature    `json:"mvp_features,omitempty"`
	FutureFeatures   []DiscoveryFutureFeature `json:"future_features,omitempty"`
	Summary          *DiscoverySummary        `json:"summary,omitempty"`
	IsReturningUser  bool                     `json:"is_returning_user"`
	UsedTemplate     string                   `json:"used_template,omitempty"`
	EditHistory      []DiscoveryEdit          `json:"edit_history,omitempty"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
}

// DiscoveryEdit represents an edit made during discovery.
type DiscoveryEdit struct {
	Stage         DiscoveryStage `json:"stage"`
	OriginalValue string         `json:"original_value"`
	NewValue      string         `json:"new_value"`
	EditedAt      time.Time      `json:"edited_at"`
}

// DiscoveryRepository defines the interface for discovery state persistence.
type DiscoveryRepository interface {
	// Create creates a new discovery state for a project.
	Create(ctx context.Context, projectID uuid.UUID) (*DiscoveryState, error)

	// GetByProjectID retrieves the discovery state for a project.
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*DiscoveryState, error)

	// GetByID retrieves a discovery state by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*DiscoveryState, error)

	// Update updates an existing discovery state.
	Update(ctx context.Context, state *DiscoveryState) error

	// UpdateStage updates the current stage and started timestamp.
	UpdateStage(ctx context.Context, id uuid.UUID, stage DiscoveryStage) error

	// SetBusinessContext sets the business context.
	SetBusinessContext(ctx context.Context, id uuid.UUID, context string) error

	// SetProblemStatement sets the problem statement.
	SetProblemStatement(ctx context.Context, id uuid.UUID, statement string) error

	// SetGoals sets the goals list.
	SetGoals(ctx context.Context, id uuid.UUID, goals []string) error

	// SetUsers sets the users information.
	SetUsers(ctx context.Context, id uuid.UUID, users *DiscoveryUsers) error

	// SetMVPFeatures sets the MVP features list.
	SetMVPFeatures(ctx context.Context, id uuid.UUID, features []DiscoveryMVPFeature) error

	// SetFutureFeatures sets the future features list.
	SetFutureFeatures(ctx context.Context, id uuid.UUID, features []DiscoveryFutureFeature) error

	// SetSummary sets the final summary.
	SetSummary(ctx context.Context, id uuid.UUID, summary *DiscoverySummary) error

	// AddEdit records an edit made during discovery.
	AddEdit(ctx context.Context, id uuid.UUID, edit DiscoveryEdit) error

	// MarkComplete marks the discovery as complete.
	MarkComplete(ctx context.Context, id uuid.UUID) error

	// Delete removes a discovery state.
	Delete(ctx context.Context, id uuid.UUID) error
}

// MockDiscoveryRepository implements DiscoveryRepository for testing.
type MockDiscoveryRepository struct {
	mu        sync.RWMutex
	states    map[uuid.UUID]*DiscoveryState // keyed by state ID
	byProject map[uuid.UUID]uuid.UUID       // projectID -> stateID lookup
}

// NewMockDiscoveryRepository creates a new MockDiscoveryRepository.
func NewMockDiscoveryRepository() *MockDiscoveryRepository {
	return &MockDiscoveryRepository{
		states:    make(map[uuid.UUID]*DiscoveryState),
		byProject: make(map[uuid.UUID]uuid.UUID),
	}
}

// Create creates a new discovery state for a project.
func (r *MockDiscoveryRepository) Create(ctx context.Context, projectID uuid.UUID) (*DiscoveryState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if discovery state already exists for this project
	if stateID, exists := r.byProject[projectID]; exists {
		return r.states[stateID], nil
	}

	now := time.Now().UTC()
	state := &DiscoveryState{
		ID:             uuid.New(),
		ProjectID:      projectID,
		Stage:          DiscoveryStageWelcome,
		StageStartedAt: now,
		Goals:          []string{},
		MVPFeatures:    []DiscoveryMVPFeature{},
		FutureFeatures: []DiscoveryFutureFeature{},
		EditHistory:    []DiscoveryEdit{},
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	r.states[state.ID] = state
	r.byProject[projectID] = state.ID

	return state, nil
}

// GetByProjectID retrieves the discovery state for a project.
func (r *MockDiscoveryRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*DiscoveryState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stateID, exists := r.byProject[projectID]
	if !exists {
		return nil, ErrNotFound
	}

	state, ok := r.states[stateID]
	if !ok {
		return nil, ErrNotFound
	}

	// Return a copy to prevent external modifications
	return r.copyState(state), nil
}

// GetByID retrieves a discovery state by its ID.
func (r *MockDiscoveryRepository) GetByID(ctx context.Context, id uuid.UUID) (*DiscoveryState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	state, ok := r.states[id]
	if !ok {
		return nil, ErrNotFound
	}

	return r.copyState(state), nil
}

// Update updates an existing discovery state.
func (r *MockDiscoveryRepository) Update(ctx context.Context, state *DiscoveryState) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.states[state.ID]; !ok {
		return ErrNotFound
	}

	state.UpdatedAt = time.Now().UTC()
	r.states[state.ID] = r.copyState(state)

	return nil
}

// UpdateStage updates the current stage and started timestamp.
func (r *MockDiscoveryRepository) UpdateStage(ctx context.Context, id uuid.UUID, stage DiscoveryStage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[id]
	if !ok {
		return ErrNotFound
	}

	now := time.Now().UTC()
	state.Stage = stage
	state.StageStartedAt = now
	state.UpdatedAt = now

	return nil
}

// SetBusinessContext sets the business context.
func (r *MockDiscoveryRepository) SetBusinessContext(ctx context.Context, id uuid.UUID, context string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[id]
	if !ok {
		return ErrNotFound
	}

	state.BusinessContext = context
	state.UpdatedAt = time.Now().UTC()

	return nil
}

// SetProblemStatement sets the problem statement.
func (r *MockDiscoveryRepository) SetProblemStatement(ctx context.Context, id uuid.UUID, statement string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[id]
	if !ok {
		return ErrNotFound
	}

	state.ProblemStatement = statement
	state.UpdatedAt = time.Now().UTC()

	return nil
}

// SetGoals sets the goals list.
func (r *MockDiscoveryRepository) SetGoals(ctx context.Context, id uuid.UUID, goals []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[id]
	if !ok {
		return ErrNotFound
	}

	state.Goals = make([]string, len(goals))
	copy(state.Goals, goals)
	state.UpdatedAt = time.Now().UTC()

	return nil
}

// SetUsers sets the users information.
func (r *MockDiscoveryRepository) SetUsers(ctx context.Context, id uuid.UUID, users *DiscoveryUsers) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[id]
	if !ok {
		return ErrNotFound
	}

	if users != nil {
		usersCopy := *users
		state.Users = &usersCopy
	} else {
		state.Users = nil
	}
	state.UpdatedAt = time.Now().UTC()

	return nil
}

// SetMVPFeatures sets the MVP features list.
func (r *MockDiscoveryRepository) SetMVPFeatures(ctx context.Context, id uuid.UUID, features []DiscoveryMVPFeature) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[id]
	if !ok {
		return ErrNotFound
	}

	state.MVPFeatures = make([]DiscoveryMVPFeature, len(features))
	copy(state.MVPFeatures, features)
	state.UpdatedAt = time.Now().UTC()

	return nil
}

// SetFutureFeatures sets the future features list.
func (r *MockDiscoveryRepository) SetFutureFeatures(ctx context.Context, id uuid.UUID, features []DiscoveryFutureFeature) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[id]
	if !ok {
		return ErrNotFound
	}

	state.FutureFeatures = make([]DiscoveryFutureFeature, len(features))
	copy(state.FutureFeatures, features)
	state.UpdatedAt = time.Now().UTC()

	return nil
}

// SetSummary sets the final summary.
func (r *MockDiscoveryRepository) SetSummary(ctx context.Context, id uuid.UUID, summary *DiscoverySummary) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[id]
	if !ok {
		return ErrNotFound
	}

	if summary != nil {
		summaryCopy := *summary
		state.Summary = &summaryCopy
	} else {
		state.Summary = nil
	}
	state.UpdatedAt = time.Now().UTC()

	return nil
}

// AddEdit records an edit made during discovery.
func (r *MockDiscoveryRepository) AddEdit(ctx context.Context, id uuid.UUID, edit DiscoveryEdit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[id]
	if !ok {
		return ErrNotFound
	}

	edit.EditedAt = time.Now().UTC()
	state.EditHistory = append(state.EditHistory, edit)
	state.UpdatedAt = time.Now().UTC()

	return nil
}

// MarkComplete marks the discovery as complete.
func (r *MockDiscoveryRepository) MarkComplete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[id]
	if !ok {
		return ErrNotFound
	}

	now := time.Now().UTC()
	state.Stage = DiscoveryStageComplete
	state.UpdatedAt = now

	if state.Summary != nil {
		state.Summary.ConfirmedAt = now
	}

	return nil
}

// Delete removes a discovery state.
func (r *MockDiscoveryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[id]
	if !ok {
		return ErrNotFound
	}

	delete(r.byProject, state.ProjectID)
	delete(r.states, id)

	return nil
}

// copyState creates a deep copy of a DiscoveryState.
func (r *MockDiscoveryRepository) copyState(state *DiscoveryState) *DiscoveryState {
	copy := *state

	// Deep copy slices
	if state.Goals != nil {
		copy.Goals = make([]string, len(state.Goals))
		for i, g := range state.Goals {
			copy.Goals[i] = g
		}
	}

	if state.MVPFeatures != nil {
		copy.MVPFeatures = make([]DiscoveryMVPFeature, len(state.MVPFeatures))
		for i, f := range state.MVPFeatures {
			copy.MVPFeatures[i] = f
		}
	}

	if state.FutureFeatures != nil {
		copy.FutureFeatures = make([]DiscoveryFutureFeature, len(state.FutureFeatures))
		for i, f := range state.FutureFeatures {
			copy.FutureFeatures[i] = f
		}
	}

	if state.EditHistory != nil {
		copy.EditHistory = make([]DiscoveryEdit, len(state.EditHistory))
		for i, e := range state.EditHistory {
			copy.EditHistory[i] = e
		}
	}

	// Deep copy pointers
	if state.Users != nil {
		usersCopy := *state.Users
		copy.Users = &usersCopy
	}

	if state.Summary != nil {
		summaryCopy := *state.Summary
		copy.Summary = &summaryCopy
	}

	return &copy
}

// Reset clears all data in the mock repository.
func (r *MockDiscoveryRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.states = make(map[uuid.UUID]*DiscoveryState)
	r.byProject = make(map[uuid.UUID]uuid.UUID)
}

// Count returns the number of discovery states stored.
func (r *MockDiscoveryRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.states)
}

// GetAll returns all discovery states (for testing purposes).
func (r *MockDiscoveryRepository) GetAll() []*DiscoveryState {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*DiscoveryState, 0, len(r.states))
	for _, state := range r.states {
		result = append(result, r.copyState(state))
	}
	return result
}
