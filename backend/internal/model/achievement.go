package model

import (
	"time"

	"github.com/google/uuid"
)

// AchievementCategory represents the type of achievement.
type AchievementCategory string

const (
	CategoryExploration   AchievementCategory = "exploration"
	CategoryUnderstanding AchievementCategory = "understanding"
	CategoryMastery       AchievementCategory = "mastery"
	CategoryGraduation    AchievementCategory = "graduation"
)

// TriggerType defines how an achievement is triggered.
type TriggerType string

const (
	TriggerEvent     TriggerType = "event"     // Single action triggers
	TriggerCount     TriggerType = "count"     // Accumulation threshold
	TriggerMilestone TriggerType = "milestone" // Complex conditions
)

// Achievement defines a learnable milestone.
type Achievement struct {
	ID            uuid.UUID           `db:"id" json:"id"`
	Code          string              `db:"code" json:"code"`
	Name          string              `db:"name" json:"name"`
	Description   string              `db:"description" json:"description"`
	Category      AchievementCategory `db:"category" json:"category"`
	Icon          string              `db:"icon" json:"icon"`
	Points        int                 `db:"points" json:"points"`
	TriggerType   TriggerType         `db:"trigger_type" json:"triggerType"`
	TriggerConfig []byte              `db:"trigger_config" json:"-"`
	Prerequisites []string            `db:"prerequisites" json:"prerequisites"`
	IsActive      bool                `db:"is_active" json:"isActive"`
	CreatedAt     time.Time           `db:"created_at" json:"createdAt"`
}

// LearningLevel represents progressive understanding levels.
type LearningLevel int

const (
	LevelFunctional LearningLevel = 1 // Understand what files DO
	LevelTree       LearningLevel = 2 // See relationships
	LevelTechnical  LearningLevel = 3 // Read code with annotations
	LevelDeveloper  LearningLevel = 4 // Full developer view
)

// UserProgress tracks learning advancement per project.
type UserProgress struct {
	ID                  uuid.UUID     `db:"id" json:"id"`
	ProjectID           uuid.UUID     `db:"project_id" json:"projectId"`
	CurrentLevel        LearningLevel `db:"current_level" json:"currentLevel"`
	TotalPoints         int           `db:"total_points" json:"totalPoints"`
	FilesViewedCount    int           `db:"files_viewed_count" json:"filesViewedCount"`
	CodeViewsCount      int           `db:"code_views_count" json:"codeViewsCount"`
	TreeExpansionsCount int           `db:"tree_expansions_count" json:"treeExpansionsCount"`
	LevelChangesCount   int           `db:"level_changes_count" json:"levelChangesCount"`
	FirstCodeViewAt     *time.Time    `db:"first_code_view_at" json:"firstCodeViewAt,omitempty"`
	FirstLevelUpAt      *time.Time    `db:"first_level_up_at" json:"firstLevelUpAt,omitempty"`
	LastActivityAt      time.Time     `db:"last_activity_at" json:"lastActivityAt"`
	CreatedAt           time.Time     `db:"created_at" json:"createdAt"`
	UpdatedAt           time.Time     `db:"updated_at" json:"updatedAt"`
}

// UserAchievement records an unlocked achievement.
type UserAchievement struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	ProjectID      uuid.UUID  `db:"project_id" json:"projectId"`
	AchievementID  uuid.UUID  `db:"achievement_id" json:"achievementId"`
	UnlockedAt     time.Time  `db:"unlocked_at" json:"unlockedAt"`
	TriggerContext []byte     `db:"trigger_context" json:"-"`
	IsSeen         bool       `db:"is_seen" json:"isSeen"`
	SeenAt         *time.Time `db:"seen_at" json:"seenAt,omitempty"`

	// Joined fields for API responses
	Achievement *Achievement `db:"-" json:"achievement,omitempty"`
}

// NudgeType represents contextual suggestions.
type NudgeType string

const (
	NudgeExploreCode   NudgeType = "explore_code"
	NudgeTryTreeView   NudgeType = "try_tree_view"
	NudgeLevelUp       NudgeType = "level_up"
	NudgeViewRelations NudgeType = "view_relations"
	NudgeExport        NudgeType = "export_project"
)

// NudgeHistory tracks which nudges have been shown.
type NudgeHistory struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	ProjectID   uuid.UUID  `db:"project_id" json:"projectId"`
	NudgeType   NudgeType  `db:"nudge_type" json:"nudgeType"`
	ShownAt     time.Time  `db:"shown_at" json:"shownAt"`
	DismissedAt *time.Time `db:"dismissed_at" json:"dismissedAt,omitempty"`
	ClickedAt   *time.Time `db:"clicked_at" json:"clickedAt,omitempty"`
	Context     []byte     `db:"context" json:"-"`
}
