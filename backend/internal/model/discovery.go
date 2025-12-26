package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DiscoveryStage represents the current stage in the discovery flow.
type DiscoveryStage string

const (
	StageWelcome  DiscoveryStage = "welcome"
	StageProblem  DiscoveryStage = "problem"
	StagePersonas DiscoveryStage = "personas"
	StageMVP      DiscoveryStage = "mvp"
	StageSummary  DiscoveryStage = "summary"
	StageComplete DiscoveryStage = "complete"
)

// ValidStages returns all valid discovery stages in order.
func ValidStages() []DiscoveryStage {
	return []DiscoveryStage{
		StageWelcome,
		StageProblem,
		StagePersonas,
		StageMVP,
		StageSummary,
		StageComplete,
	}
}

// IsValidStage checks if a stage string is a valid discovery stage.
func IsValidStage(s string) bool {
	for _, stage := range ValidStages() {
		if string(stage) == s {
			return true
		}
	}
	return false
}

// NextStage returns the next stage in the discovery flow.
// Returns empty string if already at complete or invalid stage.
func (s DiscoveryStage) NextStage() DiscoveryStage {
	stages := ValidStages()
	for i, stage := range stages {
		if stage == s && i < len(stages)-1 {
			return stages[i+1]
		}
	}
	return ""
}

// StageNumber returns the 1-based stage number (1-6).
func (s DiscoveryStage) StageNumber() int {
	stages := ValidStages()
	for i, stage := range stages {
		if stage == s {
			return i + 1
		}
	}
	return 0
}

// IsComplete returns true if the discovery stage is complete.
func (s DiscoveryStage) IsComplete() bool {
	return s == StageComplete
}

// ProjectDiscovery stores the discovery state and captured data for a project.
type ProjectDiscovery struct {
	ID             uuid.UUID      `db:"id" json:"id"`
	ProjectID      uuid.UUID      `db:"project_id" json:"projectId"`
	Stage          DiscoveryStage `db:"stage" json:"stage"`
	StageStartedAt time.Time      `db:"stage_started_at" json:"stageStartedAt"`

	// Captured data from conversation
	BusinessContext  *string `db:"business_context" json:"businessContext,omitempty"`
	ProblemStatement *string `db:"problem_statement" json:"problemStatement,omitempty"`
	GoalsJSON        []byte  `db:"goals" json:"-"` // Raw JSONB from database

	// Summary fields (populated in summary stage)
	ProjectName     *string `db:"project_name" json:"projectName,omitempty"`
	SolvesStatement *string `db:"solves_statement" json:"solvesStatement,omitempty"`

	// Metadata
	IsReturningUser *bool      `db:"is_returning_user" json:"isReturningUser,omitempty"`
	UsedTemplateID  *uuid.UUID `db:"used_template_id" json:"usedTemplateId,omitempty"`
	ConfirmedAt     *time.Time `db:"confirmed_at" json:"confirmedAt,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// Goals returns the goals as a string slice, parsing from JSON.
func (d *ProjectDiscovery) Goals() ([]string, error) {
	if d.GoalsJSON == nil || len(d.GoalsJSON) == 0 {
		return []string{}, nil
	}
	var goals []string
	if err := json.Unmarshal(d.GoalsJSON, &goals); err != nil {
		return nil, err
	}
	return goals, nil
}

// SetGoals sets the goals from a string slice, converting to JSON.
func (d *ProjectDiscovery) SetGoals(goals []string) error {
	if goals == nil {
		goals = []string{}
	}
	data, err := json.Marshal(goals)
	if err != nil {
		return err
	}
	d.GoalsJSON = data
	return nil
}

// DiscoveryUser represents a user persona defined during discovery.
type DiscoveryUser struct {
	ID              uuid.UUID `db:"id" json:"id"`
	DiscoveryID     uuid.UUID `db:"discovery_id" json:"discoveryId"`
	Description     string    `db:"description" json:"description"`
	UserCount       int       `db:"user_count" json:"count"`
	HasPermissions  bool      `db:"has_permissions" json:"hasPermissions"`
	PermissionNotes *string   `db:"permission_notes" json:"permissionNotes,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
}

// DiscoveryFeature represents a feature captured during MVP scoping.
type DiscoveryFeature struct {
	ID          uuid.UUID `db:"id" json:"id"`
	DiscoveryID uuid.UUID `db:"discovery_id" json:"discoveryId"`
	Name        string    `db:"name" json:"name"`
	Priority    int       `db:"priority" json:"priority"`
	Version     string    `db:"version" json:"version"` // "v1", "v2", etc.
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

// IsMVP returns true if this feature is part of the MVP (v1).
func (f *DiscoveryFeature) IsMVP() bool {
	return f.Version == "v1"
}

// DiscoveryEditHistory tracks edits made to discovery data.
type DiscoveryEditHistory struct {
	ID            uuid.UUID `db:"id" json:"id"`
	DiscoveryID   uuid.UUID `db:"discovery_id" json:"discoveryId"`
	Stage         string    `db:"stage" json:"stage"`
	FieldEdited   string    `db:"field_edited" json:"fieldEdited"`
	OriginalValue string    `db:"original_value" json:"originalValue"`
	NewValue      string    `db:"new_value" json:"newValue"`
	EditedAt      time.Time `db:"edited_at" json:"editedAt"`
}

// DiscoverySummary is the combined view shown to users before confirmation.
type DiscoverySummary struct {
	ProjectName     string             `json:"projectName"`
	SolvesStatement string             `json:"solvesStatement"`
	Users           []DiscoveryUser    `json:"users"`
	MVPFeatures     []DiscoveryFeature `json:"mvpFeatures"`
	FutureFeatures  []DiscoveryFeature `json:"futureFeatures"`
}

// Request/Response types for API

// CreateDiscoveryRequest represents the request payload for creating discovery.
type CreateDiscoveryRequest struct {
	ProjectID uuid.UUID `json:"projectId" validate:"required"`
}

// UpdateDiscoveryDataRequest represents the request payload for updating discovery data.
type UpdateDiscoveryDataRequest struct {
	BusinessContext  *string  `json:"businessContext,omitempty"`
	ProblemStatement *string  `json:"problemStatement,omitempty"`
	Goals            []string `json:"goals,omitempty"`
	ProjectName      *string  `json:"projectName,omitempty"`
	SolvesStatement  *string  `json:"solvesStatement,omitempty"`
}

// UpdateDiscoveryStageRequest represents the request payload for updating stage.
type UpdateDiscoveryStageRequest struct {
	Stage DiscoveryStage `json:"stage" validate:"required"`
}

// AddDiscoveryUserRequest represents the request payload for adding a user persona.
type AddDiscoveryUserRequest struct {
	Description     string  `json:"description" validate:"required"`
	UserCount       int     `json:"count"`
	HasPermissions  bool    `json:"hasPermissions"`
	PermissionNotes *string `json:"permissionNotes,omitempty"`
}

// AddDiscoveryFeatureRequest represents the request payload for adding a feature.
type AddDiscoveryFeatureRequest struct {
	Name     string `json:"name" validate:"required"`
	Priority int    `json:"priority"`
	Version  string `json:"version"` // defaults to "v1" if empty
}

// DiscoveryResponse represents the API response for discovery state.
type DiscoveryResponse struct {
	ID               uuid.UUID      `json:"id"`
	ProjectID        uuid.UUID      `json:"projectId"`
	Stage            DiscoveryStage `json:"stage"`
	StageNumber      int            `json:"stageNumber"`
	TotalStages      int            `json:"totalStages"`
	StageStartedAt   time.Time      `json:"stageStartedAt"`
	BusinessContext  *string        `json:"businessContext,omitempty"`
	ProblemStatement *string        `json:"problemStatement,omitempty"`
	Goals            []string       `json:"goals,omitempty"`
	ProjectName      *string        `json:"projectName,omitempty"`
	SolvesStatement  *string        `json:"solvesStatement,omitempty"`
	IsReturningUser  bool           `json:"isReturningUser"`
	ConfirmedAt      *time.Time     `json:"confirmedAt,omitempty"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}

// ToResponse converts a ProjectDiscovery to a DiscoveryResponse.
func (d *ProjectDiscovery) ToResponse() (*DiscoveryResponse, error) {
	goals, err := d.Goals()
	if err != nil {
		return nil, err
	}

	isReturning := false
	if d.IsReturningUser != nil {
		isReturning = *d.IsReturningUser
	}

	return &DiscoveryResponse{
		ID:               d.ID,
		ProjectID:        d.ProjectID,
		Stage:            d.Stage,
		StageNumber:      d.Stage.StageNumber(),
		TotalStages:      len(ValidStages()) - 1, // Exclude 'complete' from visible stages
		StageStartedAt:   d.StageStartedAt,
		BusinessContext:  d.BusinessContext,
		ProblemStatement: d.ProblemStatement,
		Goals:            goals,
		ProjectName:      d.ProjectName,
		SolvesStatement:  d.SolvesStatement,
		IsReturningUser:  isReturning,
		ConfirmedAt:      d.ConfirmedAt,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}, nil
}

// DiscoverySummaryResponse represents the full summary response.
type DiscoverySummaryResponse struct {
	Discovery      DiscoveryResponse  `json:"discovery"`
	Users          []DiscoveryUser    `json:"users"`
	MVPFeatures    []DiscoveryFeature `json:"mvpFeatures"`
	FutureFeatures []DiscoveryFeature `json:"futureFeatures"`
}
