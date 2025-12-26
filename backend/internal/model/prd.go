package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PRDStatus represents the lifecycle state of a PRD.
type PRDStatus string

const (
	PRDStatusPending    PRDStatus = "pending"     // Queued for generation
	PRDStatusGenerating PRDStatus = "generating"  // Claude is generating content
	PRDStatusDraft      PRDStatus = "draft"       // Generated, awaiting review
	PRDStatusReady      PRDStatus = "ready"       // Approved, ready to build
	PRDStatusInProgress PRDStatus = "in_progress" // Currently being implemented
	PRDStatusComplete   PRDStatus = "complete"    // Feature implemented
	PRDStatusFailed     PRDStatus = "failed"      // Generation failed
)

// ValidPRDStatuses returns all valid PRD statuses.
func ValidPRDStatuses() []PRDStatus {
	return []PRDStatus{
		PRDStatusPending,
		PRDStatusGenerating,
		PRDStatusDraft,
		PRDStatusReady,
		PRDStatusInProgress,
		PRDStatusComplete,
		PRDStatusFailed,
	}
}

// IsValidPRDStatus checks if a status string is a valid PRD status.
func IsValidPRDStatus(s string) bool {
	for _, status := range ValidPRDStatuses() {
		if string(status) == s {
			return true
		}
	}
	return false
}

// PRD represents a Product Requirements Document for a single feature.
type PRD struct {
	ID          uuid.UUID `db:"id" json:"id"`
	DiscoveryID uuid.UUID `db:"discovery_id" json:"discoveryId"`
	FeatureID   uuid.UUID `db:"feature_id" json:"featureId"`
	ProjectID   uuid.UUID `db:"project_id" json:"projectId"`

	// Core PRD Content
	Title    string `db:"title" json:"title"`
	Overview string `db:"overview" json:"overview"`
	Version  string `db:"version" json:"version"` // "v1", "v2", etc.
	Priority int    `db:"priority" json:"priority"`

	// Detailed Sections (JSONB for flexibility)
	UserStoriesJSON        []byte `db:"user_stories" json:"-"`
	AcceptanceCriteriaJSON []byte `db:"acceptance_criteria" json:"-"`
	TechnicalNotesJSON     []byte `db:"technical_notes" json:"-"`

	// Status Tracking
	Status      PRDStatus  `db:"status" json:"status"`
	GeneratedAt *time.Time `db:"generated_at" json:"generatedAt,omitempty"`
	ApprovedAt  *time.Time `db:"approved_at" json:"approvedAt,omitempty"`
	StartedAt   *time.Time `db:"started_at" json:"startedAt,omitempty"`
	CompletedAt *time.Time `db:"completed_at" json:"completedAt,omitempty"`

	// Metadata
	GenerationAttempts int     `db:"generation_attempts" json:"generationAttempts"`
	LastError          *string `db:"last_error" json:"lastError,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// UserStory represents a single user story within a PRD.
type UserStory struct {
	ID         string `json:"id"`         // e.g., "US-001"
	AsA        string `json:"asA"`        // User persona
	IWant      string `json:"iWant"`      // Desired action
	SoThat     string `json:"soThat"`     // Expected benefit
	Priority   string `json:"priority"`   // "must", "should", "could"
	Complexity string `json:"complexity"` // "low", "medium", "high"
}

// AcceptanceCriterion represents a single acceptance criterion.
type AcceptanceCriterion struct {
	ID          string `json:"id"`          // e.g., "AC-001"
	Given       string `json:"given"`       // Precondition
	When        string `json:"when"`        // Action
	Then        string `json:"then"`        // Expected outcome
	UserStoryID string `json:"userStoryId"` // Links to parent story
}

// TechnicalNote captures implementation guidance.
type TechnicalNote struct {
	Category    string   `json:"category"` // "architecture", "data", "ui", "integration"
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// AgentType identifies the specialized agent.
type AgentType string

const (
	AgentProductManager AgentType = "product_manager"
	AgentDesigner       AgentType = "designer"
	AgentDeveloper      AgentType = "developer"
)

// AgentContext represents the context provided to an agent.
type AgentContext struct {
	Agent     AgentType         `json:"agent"`
	PRD       *PRD              `json:"prd,omitempty"`
	Discovery *DiscoverySummary `json:"discovery"`

	// Condensed PRD for prompt
	PRDSummary string `json:"prdSummary,omitempty"`

	// Related PRDs (for cross-feature awareness)
	RelatedPRDs []PRDReference `json:"relatedPrds,omitempty"`
}

// PRDReference is a lightweight PRD summary for context.
type PRDReference struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Status   PRDStatus `json:"status"`
	Priority int       `json:"priority"`
}

// UserStories returns the user stories as a slice, parsing from JSON.
func (p *PRD) UserStories() ([]UserStory, error) {
	if p.UserStoriesJSON == nil || len(p.UserStoriesJSON) == 0 {
		return []UserStory{}, nil
	}
	var stories []UserStory
	if err := json.Unmarshal(p.UserStoriesJSON, &stories); err != nil {
		return nil, err
	}
	return stories, nil
}

// SetUserStories sets the user stories from a slice, converting to JSON.
func (p *PRD) SetUserStories(stories []UserStory) error {
	if stories == nil {
		stories = []UserStory{}
	}
	data, err := json.Marshal(stories)
	if err != nil {
		return err
	}
	p.UserStoriesJSON = data
	return nil
}

// AcceptanceCriteria returns the acceptance criteria as a slice, parsing from JSON.
func (p *PRD) AcceptanceCriteria() ([]AcceptanceCriterion, error) {
	if p.AcceptanceCriteriaJSON == nil || len(p.AcceptanceCriteriaJSON) == 0 {
		return []AcceptanceCriterion{}, nil
	}
	var criteria []AcceptanceCriterion
	if err := json.Unmarshal(p.AcceptanceCriteriaJSON, &criteria); err != nil {
		return nil, err
	}
	return criteria, nil
}

// SetAcceptanceCriteria sets the acceptance criteria from a slice, converting to JSON.
func (p *PRD) SetAcceptanceCriteria(criteria []AcceptanceCriterion) error {
	if criteria == nil {
		criteria = []AcceptanceCriterion{}
	}
	data, err := json.Marshal(criteria)
	if err != nil {
		return err
	}
	p.AcceptanceCriteriaJSON = data
	return nil
}

// TechnicalNotes returns the technical notes as a slice, parsing from JSON.
func (p *PRD) TechnicalNotes() ([]TechnicalNote, error) {
	if p.TechnicalNotesJSON == nil || len(p.TechnicalNotesJSON) == 0 {
		return []TechnicalNote{}, nil
	}
	var notes []TechnicalNote
	if err := json.Unmarshal(p.TechnicalNotesJSON, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

// SetTechnicalNotes sets the technical notes from a slice, converting to JSON.
func (p *PRD) SetTechnicalNotes(notes []TechnicalNote) error {
	if notes == nil {
		notes = []TechnicalNote{}
	}
	data, err := json.Marshal(notes)
	if err != nil {
		return err
	}
	p.TechnicalNotesJSON = data
	return nil
}

// IsMVP returns true if this PRD is for an MVP (v1) feature.
func (p *PRD) IsMVP() bool {
	return p.Version == "v1"
}

// ToReference creates a lightweight PRDReference from this PRD.
func (p *PRD) ToReference() PRDReference {
	return PRDReference{
		ID:       p.ID,
		Title:    p.Title,
		Status:   p.Status,
		Priority: p.Priority,
	}
}

// PRDResponse represents the API response for a PRD.
type PRDResponse struct {
	ID          uuid.UUID `json:"id"`
	DiscoveryID uuid.UUID `json:"discoveryId"`
	FeatureID   uuid.UUID `json:"featureId"`
	ProjectID   uuid.UUID `json:"projectId"`

	Title    string `json:"title"`
	Overview string `json:"overview"`
	Version  string `json:"version"`
	Priority int    `json:"priority"`

	UserStories        []UserStory           `json:"userStories"`
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptanceCriteria"`
	TechnicalNotes     []TechnicalNote       `json:"technicalNotes"`

	Status      PRDStatus  `json:"status"`
	GeneratedAt *time.Time `json:"generatedAt,omitempty"`
	ApprovedAt  *time.Time `json:"approvedAt,omitempty"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`

	GenerationAttempts int     `json:"generationAttempts"`
	LastError          *string `json:"lastError,omitempty"`

	IsMVP bool `json:"isMvp"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToResponse converts a PRD to a PRDResponse.
func (p *PRD) ToResponse() (*PRDResponse, error) {
	stories, err := p.UserStories()
	if err != nil {
		return nil, err
	}

	criteria, err := p.AcceptanceCriteria()
	if err != nil {
		return nil, err
	}

	notes, err := p.TechnicalNotes()
	if err != nil {
		return nil, err
	}

	return &PRDResponse{
		ID:                 p.ID,
		DiscoveryID:        p.DiscoveryID,
		FeatureID:          p.FeatureID,
		ProjectID:          p.ProjectID,
		Title:              p.Title,
		Overview:           p.Overview,
		Version:            p.Version,
		Priority:           p.Priority,
		UserStories:        stories,
		AcceptanceCriteria: criteria,
		TechnicalNotes:     notes,
		Status:             p.Status,
		GeneratedAt:        p.GeneratedAt,
		ApprovedAt:         p.ApprovedAt,
		StartedAt:          p.StartedAt,
		CompletedAt:        p.CompletedAt,
		GenerationAttempts: p.GenerationAttempts,
		LastError:          p.LastError,
		IsMVP:              p.IsMVP(),
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}, nil
}

// Request types for API

// UpdatePRDStatusRequest represents the request payload for updating PRD status.
type UpdatePRDStatusRequest struct {
	Status PRDStatus `json:"status" validate:"required"`
}

// UpdatePRDContentRequest represents the request payload for updating PRD content.
type UpdatePRDContentRequest struct {
	Title              *string               `json:"title,omitempty"`
	Overview           *string               `json:"overview,omitempty"`
	UserStories        []UserStory           `json:"userStories,omitempty"`
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptanceCriteria,omitempty"`
	TechnicalNotes     []TechnicalNote       `json:"technicalNotes,omitempty"`
}

// PRDListResponse represents a list of PRDs for API response.
type PRDListResponse struct {
	PRDs       []PRDReference `json:"prds"`
	TotalCount int            `json:"totalCount"`
	MVPCount   int            `json:"mvpCount"`
}
