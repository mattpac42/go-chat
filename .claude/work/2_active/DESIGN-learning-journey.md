# Phase 3: Learning Journey & Progression - Architecture Design

**Created**: 2025-12-26
**Author**: Architect Agent
**Status**: Draft
**Phase**: Phase 3 (Learning Journey)

---

## Executive Summary

This document defines the architecture for the Learning Journey & Progression system, enabling gamified learning milestones that track user advancement through progressive understanding of their applications. The system aligns with Vision Theme 5: "Teaching, Not Just Building."

### Design Goals

1. **Event-Driven Achievement System** - Trigger achievements from user actions across the application
2. **Progressive Learning Levels** - Track user advancement through understanding tiers
3. **Contextual Nudges** - Suggest next learning steps based on current state
4. **Non-Intrusive UX** - Celebrate achievements without disrupting workflow
5. **Extensible Design** - Easy to add new achievements and progression paths

---

## Architecture Overview

```
+------------------+     +-------------------+     +------------------+
|   Frontend       |     |   Backend         |     |   Database       |
|                  |     |                   |     |                  |
| AchievementToast |<--->| AchievementSvc    |<--->| user_progress    |
| ProgressBadge    |     | ProgressionSvc    |     | achievements     |
| NudgePopover     |     | NudgeSvc          |     | user_achievements|
| LevelIndicator   |     | EventBus          |     | nudge_history    |
+------------------+     +-------------------+     +------------------+
         ^                       ^
         |                       |
         v                       v
+------------------+     +-------------------+
| useAchievements  |     | Event Triggers    |
| useProgress      |     | (FileViewed,      |
| useNudges        |     |  LevelChanged,    |
+------------------+     |  ExportCompleted) |
                         +-------------------+
```

### Key Components

| Layer | Component | Responsibility |
|-------|-----------|----------------|
| Database | Schema | Store achievements, progress, nudge state |
| Repository | AchievementRepo, ProgressRepo | Data access patterns |
| Service | AchievementService, ProgressionService, NudgeService | Business logic |
| Handler | AchievementHandler | REST API endpoints |
| Frontend | Components, Hooks, Types | UI and state management |

---

## Database Schema

### Migration: 007_learning_journey.sql

```sql
-- 007_learning_journey.sql
-- Learning Journey & Progression System
-- Phase 3: Gamified learning milestones and user advancement tracking

-- ============================================================================
-- ACHIEVEMENTS DEFINITION TABLE
-- Stores all available achievements (seeded at startup, rarely changes)
-- ============================================================================
CREATE TABLE IF NOT EXISTS achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,          -- e.g., 'first_look', 'explorer'
    name VARCHAR(100) NOT NULL,                -- Display name: "First Look"
    description TEXT NOT NULL,                 -- What this achievement means
    category VARCHAR(30) NOT NULL              -- 'exploration', 'understanding', 'mastery'
        CHECK (category IN ('exploration', 'understanding', 'mastery', 'graduation')),
    icon VARCHAR(50) NOT NULL,                 -- Icon identifier for frontend
    points INTEGER NOT NULL DEFAULT 10,        -- Points value for gamification
    trigger_type VARCHAR(30) NOT NULL          -- How it's triggered
        CHECK (trigger_type IN ('event', 'count', 'milestone')),
    trigger_config JSONB DEFAULT '{}'::JSONB,  -- Configuration for trigger logic
    prerequisites TEXT[] DEFAULT '{}',         -- Array of achievement codes required first
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================================
-- USER PROGRESS TABLE
-- Tracks overall learning progression per project
-- ============================================================================
CREATE TABLE IF NOT EXISTS user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

    -- Current learning level (1=Functional, 2=Tree, 3=Technical, 4=Developer)
    current_level INTEGER NOT NULL DEFAULT 1
        CHECK (current_level BETWEEN 1 AND 4),

    -- Total points accumulated
    total_points INTEGER NOT NULL DEFAULT 0,

    -- Stats for trigger evaluation
    files_viewed_count INTEGER NOT NULL DEFAULT 0,
    code_views_count INTEGER NOT NULL DEFAULT 0,
    tree_expansions_count INTEGER NOT NULL DEFAULT 0,
    level_changes_count INTEGER NOT NULL DEFAULT 0,

    -- Timestamps
    first_code_view_at TIMESTAMP WITH TIME ZONE,
    first_level_up_at TIMESTAMP WITH TIME ZONE,
    last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT unique_project_progress UNIQUE (project_id)
);

-- ============================================================================
-- USER ACHIEVEMENTS TABLE
-- Junction table: which achievements each project has unlocked
-- ============================================================================
CREATE TABLE IF NOT EXISTS user_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,

    -- Context when achievement was earned
    unlocked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    trigger_context JSONB DEFAULT '{}'::JSONB,  -- What triggered it (file, action, etc.)

    -- Whether user has seen the notification
    is_seen BOOLEAN DEFAULT FALSE,
    seen_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_user_achievement UNIQUE (project_id, achievement_id)
);

-- ============================================================================
-- NUDGE HISTORY TABLE
-- Tracks which nudges have been shown to avoid repetition
-- ============================================================================
CREATE TABLE IF NOT EXISTS nudge_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    nudge_type VARCHAR(50) NOT NULL,           -- 'explore_code', 'try_tree_view', etc.

    -- When nudge was shown and user response
    shown_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    dismissed_at TIMESTAMP WITH TIME ZONE,
    clicked_at TIMESTAMP WITH TIME ZONE,       -- If user followed the nudge

    -- Context
    context JSONB DEFAULT '{}'::JSONB
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Achievements lookup
CREATE INDEX IF NOT EXISTS idx_achievements_code ON achievements(code);
CREATE INDEX IF NOT EXISTS idx_achievements_category ON achievements(category);
CREATE INDEX IF NOT EXISTS idx_achievements_active ON achievements(is_active) WHERE is_active = TRUE;

-- User progress
CREATE INDEX IF NOT EXISTS idx_user_progress_project ON user_progress(project_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_level ON user_progress(current_level);

-- User achievements
CREATE INDEX IF NOT EXISTS idx_user_achievements_project ON user_achievements(project_id);
CREATE INDEX IF NOT EXISTS idx_user_achievements_unseen ON user_achievements(project_id, is_seen)
    WHERE is_seen = FALSE;
CREATE INDEX IF NOT EXISTS idx_user_achievements_recent ON user_achievements(unlocked_at DESC);

-- Nudge history
CREATE INDEX IF NOT EXISTS idx_nudge_history_project ON nudge_history(project_id);
CREATE INDEX IF NOT EXISTS idx_nudge_history_type ON nudge_history(nudge_type);

-- ============================================================================
-- SEED DATA: Initial Achievements
-- ============================================================================
INSERT INTO achievements (code, name, description, category, icon, points, trigger_type, trigger_config) VALUES
-- Exploration achievements
('first_look', 'First Look', 'Viewed code for the first time', 'exploration', 'eye', 10,
 'event', '{"event": "code_viewed", "count": 1}'::JSONB),

('curious_mind', 'Curious Mind', 'Viewed code for 5 different files', 'exploration', 'search', 20,
 'count', '{"metric": "files_viewed_count", "threshold": 5}'::JSONB),

('deep_diver', 'Deep Diver', 'Viewed code for 10 different files', 'exploration', 'layers', 30,
 'count', '{"metric": "files_viewed_count", "threshold": 10}'::JSONB),

-- Understanding achievements
('connection_maker', 'Connection Maker', 'Expanded a relationship in the tree view', 'understanding', 'git-branch', 15,
 'event', '{"event": "tree_expanded", "count": 1}'::JSONB),

('big_picture', 'Big Picture', 'Viewed the full application tree', 'understanding', 'network', 25,
 'event', '{"event": "full_tree_viewed"}'::JSONB),

-- Mastery achievements
('level_up', 'Level Up', 'Advanced from Level 1 to Level 2 view', 'mastery', 'trending-up', 25,
 'event', '{"event": "level_changed", "from": 1, "to": 2}'::JSONB),

('explorer', 'Explorer', 'Viewed the full technical tree (Level 3)', 'mastery', 'compass', 35,
 'event', '{"event": "level_changed", "to": 3}'::JSONB),

('technologist', 'Technologist', 'Spent time in developer view (Level 4)', 'mastery', 'code', 50,
 'event', '{"event": "level_changed", "to": 4}'::JSONB),

-- Graduation achievements
('graduate', 'Graduate', 'Exported project to VS Code or local development', 'graduation', 'graduation-cap', 100,
 'event', '{"event": "project_exported"}'::JSONB),

('self_sufficient', 'Self-Sufficient', 'Made a direct code edit after learning', 'graduation', 'edit', 75,
 'event', '{"event": "code_edited"}'::JSONB)

ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- COMMENTS
-- ============================================================================
COMMENT ON TABLE achievements IS 'Defines all available achievements users can unlock';
COMMENT ON TABLE user_progress IS 'Tracks learning progression and stats per project';
COMMENT ON TABLE user_achievements IS 'Records which achievements each project has unlocked';
COMMENT ON TABLE nudge_history IS 'Tracks nudge interactions to avoid repetition';

COMMENT ON COLUMN achievements.trigger_type IS 'event=single action, count=accumulation threshold, milestone=complex condition';
COMMENT ON COLUMN achievements.trigger_config IS 'JSON configuration for trigger logic evaluation';
COMMENT ON COLUMN user_progress.current_level IS '1=Functional, 2=Tree, 3=Technical, 4=Developer';
```

### Entity Relationship Diagram

```
+----------------+       +-------------------+       +----------------+
|   projects     |       |   user_progress   |       |  achievements  |
|----------------|       |-------------------|       |----------------|
| id (PK)        |<---+  | id (PK)           |       | id (PK)        |
| ...            |    |  | project_id (FK)   |---+   | code (UNIQUE)  |
+----------------+    |  | current_level     |   |   | name           |
                      |  | total_points      |   |   | category       |
                      |  | files_viewed_count|   |   | trigger_type   |
                      |  | ...               |   |   | trigger_config |
                      |  +-------------------+   |   +----------------+
                      |                          |          ^
                      |  +-------------------+   |          |
                      +->| user_achievements |<--+          |
                         |-------------------|              |
                         | id (PK)           |              |
                         | project_id (FK)   |--------------+
                         | achievement_id(FK)|
                         | unlocked_at       |
                         | is_seen           |
                         +-------------------+
```

---

## Backend Service Layer

### Models: `/backend/internal/model/achievement.go`

```go
package model

import (
    "time"
    "github.com/google/uuid"
)

// AchievementCategory represents the type of achievement
type AchievementCategory string

const (
    CategoryExploration   AchievementCategory = "exploration"
    CategoryUnderstanding AchievementCategory = "understanding"
    CategoryMastery       AchievementCategory = "mastery"
    CategoryGraduation    AchievementCategory = "graduation"
)

// TriggerType defines how an achievement is triggered
type TriggerType string

const (
    TriggerEvent     TriggerType = "event"     // Single action triggers
    TriggerCount     TriggerType = "count"     // Accumulation threshold
    TriggerMilestone TriggerType = "milestone" // Complex conditions
)

// Achievement defines a learnable milestone
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

// LearningLevel represents progressive understanding levels
type LearningLevel int

const (
    LevelFunctional LearningLevel = 1 // Understand what files DO
    LevelTree       LearningLevel = 2 // See relationships
    LevelTechnical  LearningLevel = 3 // Read code with annotations
    LevelDeveloper  LearningLevel = 4 // Full developer view
)

// UserProgress tracks learning advancement per project
type UserProgress struct {
    ID                   uuid.UUID     `db:"id" json:"id"`
    ProjectID            uuid.UUID     `db:"project_id" json:"projectId"`
    CurrentLevel         LearningLevel `db:"current_level" json:"currentLevel"`
    TotalPoints          int           `db:"total_points" json:"totalPoints"`
    FilesViewedCount     int           `db:"files_viewed_count" json:"filesViewedCount"`
    CodeViewsCount       int           `db:"code_views_count" json:"codeViewsCount"`
    TreeExpansionsCount  int           `db:"tree_expansions_count" json:"treeExpansionsCount"`
    LevelChangesCount    int           `db:"level_changes_count" json:"levelChangesCount"`
    FirstCodeViewAt      *time.Time    `db:"first_code_view_at" json:"firstCodeViewAt,omitempty"`
    FirstLevelUpAt       *time.Time    `db:"first_level_up_at" json:"firstLevelUpAt,omitempty"`
    LastActivityAt       time.Time     `db:"last_activity_at" json:"lastActivityAt"`
    CreatedAt            time.Time     `db:"created_at" json:"createdAt"`
    UpdatedAt            time.Time     `db:"updated_at" json:"updatedAt"`
}

// UserAchievement records an unlocked achievement
type UserAchievement struct {
    ID             uuid.UUID   `db:"id" json:"id"`
    ProjectID      uuid.UUID   `db:"project_id" json:"projectId"`
    AchievementID  uuid.UUID   `db:"achievement_id" json:"achievementId"`
    UnlockedAt     time.Time   `db:"unlocked_at" json:"unlockedAt"`
    TriggerContext []byte      `db:"trigger_context" json:"-"`
    IsSeen         bool        `db:"is_seen" json:"isSeen"`
    SeenAt         *time.Time  `db:"seen_at" json:"seenAt,omitempty"`

    // Joined fields for API responses
    Achievement    *Achievement `db:"-" json:"achievement,omitempty"`
}

// NudgeType represents contextual suggestions
type NudgeType string

const (
    NudgeExploreCode   NudgeType = "explore_code"
    NudgeTryTreeView   NudgeType = "try_tree_view"
    NudgeLevelUp       NudgeType = "level_up"
    NudgeViewRelations NudgeType = "view_relations"
    NudgeExport        NudgeType = "export_project"
)

// NudgeHistory tracks which nudges have been shown
type NudgeHistory struct {
    ID          uuid.UUID  `db:"id" json:"id"`
    ProjectID   uuid.UUID  `db:"project_id" json:"projectId"`
    NudgeType   NudgeType  `db:"nudge_type" json:"nudgeType"`
    ShownAt     time.Time  `db:"shown_at" json:"shownAt"`
    DismissedAt *time.Time `db:"dismissed_at" json:"dismissedAt,omitempty"`
    ClickedAt   *time.Time `db:"clicked_at" json:"clickedAt,omitempty"`
    Context     []byte     `db:"context" json:"-"`
}
```

### Repository Interface: `/backend/internal/repository/achievement.go`

```go
package repository

import (
    "context"
    "github.com/google/uuid"
    "gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// AchievementRepository defines data access for achievements
type AchievementRepository interface {
    // Achievement definitions
    GetAllActive(ctx context.Context) ([]model.Achievement, error)
    GetByCode(ctx context.Context, code string) (*model.Achievement, error)
    GetByCategory(ctx context.Context, category model.AchievementCategory) ([]model.Achievement, error)

    // User progress
    GetProgress(ctx context.Context, projectID uuid.UUID) (*model.UserProgress, error)
    CreateProgress(ctx context.Context, projectID uuid.UUID) (*model.UserProgress, error)
    UpdateProgress(ctx context.Context, progress *model.UserProgress) error
    IncrementStat(ctx context.Context, projectID uuid.UUID, stat string, delta int) error

    // User achievements
    GetUserAchievements(ctx context.Context, projectID uuid.UUID) ([]model.UserAchievement, error)
    GetUnseenAchievements(ctx context.Context, projectID uuid.UUID) ([]model.UserAchievement, error)
    HasAchievement(ctx context.Context, projectID uuid.UUID, achievementCode string) (bool, error)
    UnlockAchievement(ctx context.Context, projectID uuid.UUID, achievementID uuid.UUID, context []byte) (*model.UserAchievement, error)
    MarkAchievementSeen(ctx context.Context, id uuid.UUID) error

    // Nudge history
    GetRecentNudges(ctx context.Context, projectID uuid.UUID, limit int) ([]model.NudgeHistory, error)
    HasSeenNudge(ctx context.Context, projectID uuid.UUID, nudgeType model.NudgeType) (bool, error)
    RecordNudgeShown(ctx context.Context, projectID uuid.UUID, nudgeType model.NudgeType) error
    RecordNudgeDismissed(ctx context.Context, id uuid.UUID) error
    RecordNudgeClicked(ctx context.Context, id uuid.UUID) error
}
```

### Achievement Service: `/backend/internal/service/achievement.go`

```go
package service

import (
    "context"
    "encoding/json"

    "github.com/google/uuid"
    "github.com/rs/zerolog"
    "gitlab.yuki.lan/goodies/gochat/backend/internal/model"
    "gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

// LearningEvent represents an action that may trigger achievements
type LearningEvent struct {
    Type      string                 `json:"type"`
    ProjectID uuid.UUID              `json:"projectId"`
    Context   map[string]interface{} `json:"context"`
}

// Event types
const (
    EventCodeViewed    = "code_viewed"
    EventTreeExpanded  = "tree_expanded"
    EventFullTreeView  = "full_tree_viewed"
    EventLevelChanged  = "level_changed"
    EventProjectExport = "project_exported"
    EventCodeEdited    = "code_edited"
)

// AchievementService handles achievement logic
type AchievementService struct {
    repo   repository.AchievementRepository
    logger zerolog.Logger
}

// NewAchievementService creates a new AchievementService
func NewAchievementService(repo repository.AchievementRepository, logger zerolog.Logger) *AchievementService {
    return &AchievementService{
        repo:   repo,
        logger: logger,
    }
}

// ProcessEvent evaluates an event against all achievement triggers
func (s *AchievementService) ProcessEvent(ctx context.Context, event LearningEvent) ([]model.UserAchievement, error) {
    var unlocked []model.UserAchievement

    // Get or create progress
    progress, err := s.getOrCreateProgress(ctx, event.ProjectID)
    if err != nil {
        return nil, err
    }

    // Update stats based on event type
    if err := s.updateStats(ctx, event, progress); err != nil {
        s.logger.Warn().Err(err).Msg("failed to update stats")
    }

    // Get all active achievements
    achievements, err := s.repo.GetAllActive(ctx)
    if err != nil {
        return nil, err
    }

    // Evaluate each achievement
    for _, achievement := range achievements {
        // Check if already unlocked
        hasIt, err := s.repo.HasAchievement(ctx, event.ProjectID, achievement.Code)
        if err != nil {
            s.logger.Warn().Err(err).Str("code", achievement.Code).Msg("failed to check achievement")
            continue
        }
        if hasIt {
            continue
        }

        // Check prerequisites
        if !s.checkPrerequisites(ctx, event.ProjectID, achievement.Prerequisites) {
            continue
        }

        // Evaluate trigger
        if s.evaluateTrigger(event, progress, achievement) {
            contextBytes, _ := json.Marshal(event.Context)
            userAchievement, err := s.repo.UnlockAchievement(ctx, event.ProjectID, achievement.ID, contextBytes)
            if err != nil {
                s.logger.Warn().Err(err).Str("code", achievement.Code).Msg("failed to unlock achievement")
                continue
            }

            // Add points
            progress.TotalPoints += achievement.Points
            s.repo.UpdateProgress(ctx, progress)

            userAchievement.Achievement = &achievement
            unlocked = append(unlocked, *userAchievement)

            s.logger.Info().
                Str("projectId", event.ProjectID.String()).
                Str("achievement", achievement.Code).
                Int("points", achievement.Points).
                Msg("achievement unlocked")
        }
    }

    return unlocked, nil
}

// evaluateTrigger checks if an achievement's trigger conditions are met
func (s *AchievementService) evaluateTrigger(event LearningEvent, progress *model.UserProgress, achievement model.Achievement) bool {
    var config map[string]interface{}
    if err := json.Unmarshal(achievement.TriggerConfig, &config); err != nil {
        return false
    }

    switch achievement.TriggerType {
    case model.TriggerEvent:
        return s.evaluateEventTrigger(event, config)
    case model.TriggerCount:
        return s.evaluateCountTrigger(progress, config)
    case model.TriggerMilestone:
        return s.evaluateMilestoneTrigger(event, progress, config)
    }

    return false
}

// evaluateEventTrigger checks event-based triggers
func (s *AchievementService) evaluateEventTrigger(event LearningEvent, config map[string]interface{}) bool {
    expectedEvent, ok := config["event"].(string)
    if !ok || expectedEvent != event.Type {
        return false
    }

    // Check additional conditions (e.g., level transitions)
    if from, ok := config["from"].(float64); ok {
        if eventFrom, ok := event.Context["from"].(int); ok {
            if float64(eventFrom) != from {
                return false
            }
        }
    }

    if to, ok := config["to"].(float64); ok {
        if eventTo, ok := event.Context["to"].(int); ok {
            if float64(eventTo) != to {
                return false
            }
        }
    }

    return true
}

// evaluateCountTrigger checks accumulation-based triggers
func (s *AchievementService) evaluateCountTrigger(progress *model.UserProgress, config map[string]interface{}) bool {
    metric, ok := config["metric"].(string)
    if !ok {
        return false
    }

    threshold, ok := config["threshold"].(float64)
    if !ok {
        return false
    }

    var value int
    switch metric {
    case "files_viewed_count":
        value = progress.FilesViewedCount
    case "code_views_count":
        value = progress.CodeViewsCount
    case "tree_expansions_count":
        value = progress.TreeExpansionsCount
    case "level_changes_count":
        value = progress.LevelChangesCount
    }

    return value >= int(threshold)
}

// GetProgress returns current progress for a project
func (s *AchievementService) GetProgress(ctx context.Context, projectID uuid.UUID) (*model.UserProgress, error) {
    return s.getOrCreateProgress(ctx, projectID)
}

// GetAchievements returns all achievements for a project
func (s *AchievementService) GetAchievements(ctx context.Context, projectID uuid.UUID) ([]model.UserAchievement, error) {
    return s.repo.GetUserAchievements(ctx, projectID)
}

// GetUnseenAchievements returns achievements not yet seen
func (s *AchievementService) GetUnseenAchievements(ctx context.Context, projectID uuid.UUID) ([]model.UserAchievement, error) {
    return s.repo.GetUnseenAchievements(ctx, projectID)
}

// MarkSeen marks an achievement as seen
func (s *AchievementService) MarkSeen(ctx context.Context, achievementID uuid.UUID) error {
    return s.repo.MarkAchievementSeen(ctx, achievementID)
}

// UpdateLevel updates the user's learning level
func (s *AchievementService) UpdateLevel(ctx context.Context, projectID uuid.UUID, newLevel model.LearningLevel) error {
    progress, err := s.getOrCreateProgress(ctx, projectID)
    if err != nil {
        return err
    }

    oldLevel := progress.CurrentLevel
    progress.CurrentLevel = newLevel
    progress.LevelChangesCount++

    if err := s.repo.UpdateProgress(ctx, progress); err != nil {
        return err
    }

    // Trigger level change event
    event := LearningEvent{
        Type:      EventLevelChanged,
        ProjectID: projectID,
        Context: map[string]interface{}{
            "from": int(oldLevel),
            "to":   int(newLevel),
        },
    }

    s.ProcessEvent(ctx, event)
    return nil
}
```

### Nudge Service: `/backend/internal/service/nudge.go`

```go
package service

import (
    "context"

    "github.com/google/uuid"
    "github.com/rs/zerolog"
    "gitlab.yuki.lan/goodies/gochat/backend/internal/model"
    "gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

// Nudge represents a contextual suggestion
type Nudge struct {
    Type    model.NudgeType `json:"type"`
    Title   string          `json:"title"`
    Message string          `json:"message"`
    Action  string          `json:"action"`
    Icon    string          `json:"icon"`
}

// NudgeService manages contextual suggestions
type NudgeService struct {
    repo            repository.AchievementRepository
    achievementSvc  *AchievementService
    logger          zerolog.Logger
}

// NewNudgeService creates a new NudgeService
func NewNudgeService(repo repository.AchievementRepository, achievementSvc *AchievementService, logger zerolog.Logger) *NudgeService {
    return &NudgeService{
        repo:           repo,
        achievementSvc: achievementSvc,
        logger:         logger,
    }
}

// GetNextNudge returns the most appropriate nudge for current context
func (s *NudgeService) GetNextNudge(ctx context.Context, projectID uuid.UUID) (*Nudge, error) {
    progress, err := s.achievementSvc.GetProgress(ctx, projectID)
    if err != nil {
        return nil, err
    }

    // Nudge selection logic based on progress

    // Haven't viewed any code yet
    if progress.CodeViewsCount == 0 {
        if hasSeenNudge, _ := s.repo.HasSeenNudge(ctx, projectID, model.NudgeExploreCode); !hasSeenNudge {
            return &Nudge{
                Type:    model.NudgeExploreCode,
                Title:   "Peek Behind the Curtain",
                Message: "Curious what makes your app tick? Tap any file to see how it works.",
                Action:  "Show me",
                Icon:    "eye",
            }, nil
        }
    }

    // Viewed code but never expanded tree
    if progress.CodeViewsCount >= 3 && progress.TreeExpansionsCount == 0 {
        if hasSeenNudge, _ := s.repo.HasSeenNudge(ctx, projectID, model.NudgeTryTreeView); !hasSeenNudge {
            return &Nudge{
                Type:    model.NudgeTryTreeView,
                Title:   "See the Big Picture",
                Message: "Want to see how your files connect? Try the tree view!",
                Action:  "Show connections",
                Icon:    "git-branch",
            }, nil
        }
    }

    // At level 1 for a while, suggest level up
    if progress.CurrentLevel == model.LevelFunctional && progress.FilesViewedCount >= 5 {
        if hasSeenNudge, _ := s.repo.HasSeenNudge(ctx, projectID, model.NudgeLevelUp); !hasSeenNudge {
            return &Nudge{
                Type:    model.NudgeLevelUp,
                Title:   "Ready for More?",
                Message: "You're getting the hang of it! Try Level 2 to see more details.",
                Action:  "Level up",
                Icon:    "trending-up",
            }, nil
        }
    }

    return nil, nil
}

// RecordNudgeAction records user interaction with a nudge
func (s *NudgeService) RecordNudgeAction(ctx context.Context, projectID uuid.UUID, nudgeType model.NudgeType, action string) error {
    switch action {
    case "shown":
        return s.repo.RecordNudgeShown(ctx, projectID, nudgeType)
    case "dismissed":
        // Find the most recent nudge of this type and mark dismissed
        nudges, err := s.repo.GetRecentNudges(ctx, projectID, 10)
        if err != nil {
            return err
        }
        for _, n := range nudges {
            if n.NudgeType == nudgeType && n.DismissedAt == nil {
                return s.repo.RecordNudgeDismissed(ctx, n.ID)
            }
        }
    case "clicked":
        nudges, err := s.repo.GetRecentNudges(ctx, projectID, 10)
        if err != nil {
            return err
        }
        for _, n := range nudges {
            if n.NudgeType == nudgeType && n.ClickedAt == nil {
                return s.repo.RecordNudgeClicked(ctx, n.ID)
            }
        }
    }
    return nil
}
```

---

## REST API Endpoints

### Handler: `/backend/internal/handler/achievement.go`

```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/rs/zerolog"
    "gitlab.yuki.lan/goodies/gochat/backend/internal/model"
    "gitlab.yuki.lan/goodies/gochat/backend/internal/service"
)

type AchievementHandler struct {
    achievementSvc *service.AchievementService
    nudgeSvc       *service.NudgeService
    logger         zerolog.Logger
}

func NewAchievementHandler(achievementSvc *service.AchievementService, nudgeSvc *service.NudgeService, logger zerolog.Logger) *AchievementHandler {
    return &AchievementHandler{
        achievementSvc: achievementSvc,
        nudgeSvc:       nudgeSvc,
        logger:         logger,
    }
}

// RegisterRoutes registers achievement routes
func (h *AchievementHandler) RegisterRoutes(router *gin.RouterGroup) {
    router.GET("/projects/:id/progress", h.GetProgress)
    router.GET("/projects/:id/achievements", h.GetAchievements)
    router.GET("/projects/:id/achievements/unseen", h.GetUnseenAchievements)
    router.POST("/projects/:id/achievements/:achievementId/seen", h.MarkAchievementSeen)
    router.POST("/projects/:id/events", h.RecordEvent)
    router.PUT("/projects/:id/level", h.UpdateLevel)
    router.GET("/projects/:id/nudge", h.GetNextNudge)
    router.POST("/projects/:id/nudge/:type/action", h.RecordNudgeAction)
}
```

### API Endpoint Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/projects/:id/progress` | Get user progress and stats |
| GET | `/api/projects/:id/achievements` | Get all unlocked achievements |
| GET | `/api/projects/:id/achievements/unseen` | Get unseen achievements for notifications |
| POST | `/api/projects/:id/achievements/:achievementId/seen` | Mark achievement as seen |
| POST | `/api/projects/:id/events` | Record learning event (triggers achievement check) |
| PUT | `/api/projects/:id/level` | Update learning level |
| GET | `/api/projects/:id/nudge` | Get next contextual nudge |
| POST | `/api/projects/:id/nudge/:type/action` | Record nudge interaction |

### Response Examples

**GET /api/projects/:id/progress**
```json
{
  "id": "uuid",
  "projectId": "uuid",
  "currentLevel": 2,
  "totalPoints": 85,
  "filesViewedCount": 7,
  "codeViewsCount": 12,
  "treeExpansionsCount": 3,
  "levelChangesCount": 1,
  "firstCodeViewAt": "2025-12-26T10:30:00Z",
  "firstLevelUpAt": "2025-12-26T11:00:00Z",
  "lastActivityAt": "2025-12-26T14:22:00Z"
}
```

**GET /api/projects/:id/achievements**
```json
{
  "achievements": [
    {
      "id": "uuid",
      "achievementId": "uuid",
      "unlockedAt": "2025-12-26T10:35:00Z",
      "isSeen": true,
      "achievement": {
        "code": "first_look",
        "name": "First Look",
        "description": "Viewed code for the first time",
        "category": "exploration",
        "icon": "eye",
        "points": 10
      }
    }
  ],
  "totalPoints": 85,
  "achievementCount": 5
}
```

**POST /api/projects/:id/events**
```json
// Request
{
  "type": "code_viewed",
  "context": {
    "fileId": "uuid",
    "fileName": "api.go"
  }
}

// Response (new achievements unlocked)
{
  "unlocked": [
    {
      "code": "first_look",
      "name": "First Look",
      "points": 10
    }
  ],
  "progress": {
    "currentLevel": 1,
    "totalPoints": 10
  }
}
```

---

## Frontend Components

### Types: `/frontend/src/types/achievements.ts`

```typescript
export type AchievementCategory = 'exploration' | 'understanding' | 'mastery' | 'graduation';

export interface Achievement {
  id: string;
  code: string;
  name: string;
  description: string;
  category: AchievementCategory;
  icon: string;
  points: number;
}

export interface UserAchievement {
  id: string;
  achievementId: string;
  unlockedAt: string;
  isSeen: boolean;
  achievement?: Achievement;
}

export type LearningLevel = 1 | 2 | 3 | 4;

export interface UserProgress {
  id: string;
  projectId: string;
  currentLevel: LearningLevel;
  totalPoints: number;
  filesViewedCount: number;
  codeViewsCount: number;
  treeExpansionsCount: number;
  levelChangesCount: number;
  firstCodeViewAt?: string;
  firstLevelUpAt?: string;
  lastActivityAt: string;
}

export interface Nudge {
  type: string;
  title: string;
  message: string;
  action: string;
  icon: string;
}

export interface LearningEvent {
  type: string;
  context?: Record<string, unknown>;
}
```

### Hook: `/frontend/src/hooks/useAchievements.ts`

```typescript
'use client';

import { useState, useCallback, useEffect } from 'react';
import { API_BASE_URL } from '@/lib/api';
import { UserProgress, UserAchievement, Nudge, LearningEvent } from '@/types/achievements';

interface UseAchievementsReturn {
  progress: UserProgress | null;
  achievements: UserAchievement[];
  unseenAchievements: UserAchievement[];
  currentNudge: Nudge | null;
  isLoading: boolean;
  error: string | null;
  recordEvent: (event: LearningEvent) => Promise<UserAchievement[]>;
  markAchievementSeen: (id: string) => Promise<void>;
  dismissNudge: () => Promise<void>;
  acceptNudge: () => Promise<void>;
  refetch: () => Promise<void>;
}

export function useAchievements(projectId: string | null): UseAchievementsReturn {
  const [progress, setProgress] = useState<UserProgress | null>(null);
  const [achievements, setAchievements] = useState<UserAchievement[]>([]);
  const [unseenAchievements, setUnseenAchievements] = useState<UserAchievement[]>([]);
  const [currentNudge, setCurrentNudge] = useState<Nudge | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchProgress = useCallback(async () => {
    if (!projectId) return;

    try {
      const response = await fetch(`${API_BASE_URL}/api/projects/${projectId}/progress`);
      if (response.ok) {
        const data = await response.json();
        setProgress(data);
      }
    } catch (err) {
      console.error('Failed to fetch progress:', err);
    }
  }, [projectId]);

  const fetchAchievements = useCallback(async () => {
    if (!projectId) return;

    try {
      const [allRes, unseenRes] = await Promise.all([
        fetch(`${API_BASE_URL}/api/projects/${projectId}/achievements`),
        fetch(`${API_BASE_URL}/api/projects/${projectId}/achievements/unseen`)
      ]);

      if (allRes.ok) {
        const data = await allRes.json();
        setAchievements(data.achievements || []);
      }

      if (unseenRes.ok) {
        const data = await unseenRes.json();
        setUnseenAchievements(data.achievements || []);
      }
    } catch (err) {
      console.error('Failed to fetch achievements:', err);
    }
  }, [projectId]);

  const fetchNudge = useCallback(async () => {
    if (!projectId) return;

    try {
      const response = await fetch(`${API_BASE_URL}/api/projects/${projectId}/nudge`);
      if (response.ok) {
        const data = await response.json();
        setCurrentNudge(data);
      }
    } catch (err) {
      console.error('Failed to fetch nudge:', err);
    }
  }, [projectId]);

  const recordEvent = useCallback(async (event: LearningEvent): Promise<UserAchievement[]> => {
    if (!projectId) return [];

    try {
      const response = await fetch(`${API_BASE_URL}/api/projects/${projectId}/events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(event)
      });

      if (response.ok) {
        const data = await response.json();

        // Refresh progress and achievements
        await Promise.all([fetchProgress(), fetchAchievements()]);

        return data.unlocked || [];
      }
    } catch (err) {
      console.error('Failed to record event:', err);
    }

    return [];
  }, [projectId, fetchProgress, fetchAchievements]);

  const markAchievementSeen = useCallback(async (id: string) => {
    if (!projectId) return;

    try {
      await fetch(`${API_BASE_URL}/api/projects/${projectId}/achievements/${id}/seen`, {
        method: 'POST'
      });

      setUnseenAchievements(prev => prev.filter(a => a.id !== id));
    } catch (err) {
      console.error('Failed to mark achievement seen:', err);
    }
  }, [projectId]);

  // Initial fetch
  useEffect(() => {
    if (!projectId) {
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    Promise.all([fetchProgress(), fetchAchievements(), fetchNudge()])
      .finally(() => setIsLoading(false));
  }, [projectId, fetchProgress, fetchAchievements, fetchNudge]);

  return {
    progress,
    achievements,
    unseenAchievements,
    currentNudge,
    isLoading,
    error,
    recordEvent,
    markAchievementSeen,
    dismissNudge: async () => { /* ... */ },
    acceptNudge: async () => { /* ... */ },
    refetch: async () => {
      await Promise.all([fetchProgress(), fetchAchievements(), fetchNudge()]);
    }
  };
}
```

### Component: AchievementToast

```tsx
// /frontend/src/components/achievements/AchievementToast.tsx
'use client';

import { useEffect, useState } from 'react';
import { UserAchievement } from '@/types/achievements';
import { Trophy, X, Star } from 'lucide-react';

interface AchievementToastProps {
  achievement: UserAchievement;
  onDismiss: () => void;
  autoDismissMs?: number;
}

export function AchievementToast({
  achievement,
  onDismiss,
  autoDismissMs = 5000
}: AchievementToastProps) {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    // Animate in
    requestAnimationFrame(() => setIsVisible(true));

    // Auto-dismiss
    const timer = setTimeout(() => {
      setIsVisible(false);
      setTimeout(onDismiss, 300);
    }, autoDismissMs);

    return () => clearTimeout(timer);
  }, [autoDismissMs, onDismiss]);

  const { achievement: ach } = achievement;
  if (!ach) return null;

  return (
    <div
      className={`fixed bottom-4 right-4 z-50 transform transition-all duration-300 ${
        isVisible ? 'translate-y-0 opacity-100' : 'translate-y-8 opacity-0'
      }`}
    >
      <div className="bg-gradient-to-r from-amber-500 to-orange-500 text-white rounded-lg shadow-lg p-4 max-w-sm">
        <div className="flex items-start gap-3">
          <div className="flex-shrink-0 bg-white/20 rounded-full p-2">
            <Trophy className="w-6 h-6" />
          </div>

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <span className="font-semibold">Achievement Unlocked!</span>
              <span className="text-amber-200 text-sm flex items-center gap-1">
                <Star className="w-3 h-3" />
                +{ach.points}
              </span>
            </div>
            <p className="font-bold text-lg">{ach.name}</p>
            <p className="text-white/80 text-sm">{ach.description}</p>
          </div>

          <button
            onClick={() => {
              setIsVisible(false);
              setTimeout(onDismiss, 300);
            }}
            className="flex-shrink-0 text-white/60 hover:text-white"
          >
            <X className="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  );
}
```

### Component: ProgressBadge

```tsx
// /frontend/src/components/achievements/ProgressBadge.tsx
'use client';

import { UserProgress, LearningLevel } from '@/types/achievements';
import { Star, TrendingUp } from 'lucide-react';

interface ProgressBadgeProps {
  progress: UserProgress;
  compact?: boolean;
  onClick?: () => void;
}

const LEVEL_NAMES: Record<LearningLevel, string> = {
  1: 'Explorer',
  2: 'Navigator',
  3: 'Technologist',
  4: 'Developer'
};

const LEVEL_COLORS: Record<LearningLevel, string> = {
  1: 'bg-emerald-500',
  2: 'bg-blue-500',
  3: 'bg-purple-500',
  4: 'bg-amber-500'
};

export function ProgressBadge({ progress, compact = false, onClick }: ProgressBadgeProps) {
  if (compact) {
    return (
      <button
        onClick={onClick}
        className="flex items-center gap-1.5 px-2 py-1 rounded-full bg-gray-100 hover:bg-gray-200 transition-colors"
      >
        <div className={`w-2 h-2 rounded-full ${LEVEL_COLORS[progress.currentLevel]}`} />
        <span className="text-xs font-medium text-gray-600">
          Lvl {progress.currentLevel}
        </span>
        <span className="text-xs text-gray-400 flex items-center gap-0.5">
          <Star className="w-3 h-3" />
          {progress.totalPoints}
        </span>
      </button>
    );
  }

  return (
    <div
      className="bg-white rounded-lg shadow-sm border p-4 cursor-pointer hover:shadow-md transition-shadow"
      onClick={onClick}
    >
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <div className={`w-3 h-3 rounded-full ${LEVEL_COLORS[progress.currentLevel]}`} />
          <span className="font-medium">
            Level {progress.currentLevel}: {LEVEL_NAMES[progress.currentLevel]}
          </span>
        </div>
        <div className="flex items-center gap-1 text-amber-500">
          <Star className="w-4 h-4" />
          <span className="font-semibold">{progress.totalPoints}</span>
        </div>
      </div>

      {/* Progress to next level */}
      {progress.currentLevel < 4 && (
        <div className="space-y-1">
          <div className="flex justify-between text-xs text-gray-500">
            <span>Progress to Level {progress.currentLevel + 1}</span>
            <span>{Math.min(progress.filesViewedCount, 10)}/10 files explored</span>
          </div>
          <div className="h-2 bg-gray-100 rounded-full overflow-hidden">
            <div
              className={`h-full ${LEVEL_COLORS[progress.currentLevel]} transition-all`}
              style={{ width: `${Math.min((progress.filesViewedCount / 10) * 100, 100)}%` }}
            />
          </div>
        </div>
      )}
    </div>
  );
}
```

### Component: NudgePopover

```tsx
// /frontend/src/components/achievements/NudgePopover.tsx
'use client';

import { Nudge } from '@/types/achievements';
import { X, ChevronRight } from 'lucide-react';
import * as Icons from 'lucide-react';

interface NudgePopoverProps {
  nudge: Nudge;
  onAccept: () => void;
  onDismiss: () => void;
}

export function NudgePopover({ nudge, onAccept, onDismiss }: NudgePopoverProps) {
  // Dynamically get icon component
  const IconComponent = (Icons as Record<string, React.FC<{ className?: string }>>)[
    nudge.icon.charAt(0).toUpperCase() + nudge.icon.slice(1).replace(/-([a-z])/g, g => g[1].toUpperCase())
  ] || Icons.Lightbulb;

  return (
    <div className="fixed bottom-20 left-4 right-4 md:left-auto md:right-4 md:w-80 z-40 animate-slide-up">
      <div className="bg-teal-600 text-white rounded-lg shadow-lg p-4">
        <div className="flex items-start gap-3">
          <div className="flex-shrink-0 bg-white/20 rounded-full p-2">
            <IconComponent className="w-5 h-5" />
          </div>

          <div className="flex-1 min-w-0">
            <p className="font-semibold">{nudge.title}</p>
            <p className="text-white/80 text-sm mt-1">{nudge.message}</p>

            <div className="flex items-center gap-2 mt-3">
              <button
                onClick={onAccept}
                className="flex items-center gap-1 bg-white text-teal-600 px-3 py-1.5 rounded-md text-sm font-medium hover:bg-white/90 transition-colors"
              >
                {nudge.action}
                <ChevronRight className="w-4 h-4" />
              </button>
              <button
                onClick={onDismiss}
                className="text-white/60 hover:text-white text-sm"
              >
                Maybe later
              </button>
            </div>
          </div>

          <button
            onClick={onDismiss}
            className="flex-shrink-0 text-white/60 hover:text-white"
          >
            <X className="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  );
}
```

---

## Event Trigger Integration Points

### Where Events Are Triggered

| Event | Trigger Location | Component/Action |
|-------|------------------|------------------|
| `code_viewed` | FilePreviewModal opens | FileExplorer.tsx when user taps "View Code" |
| `tree_expanded` | Tree node expanded | FileTree.tsx when user expands a branch |
| `full_tree_viewed` | All nodes visible | FileTree.tsx when tree is fully expanded |
| `level_changed` | Level selector changed | ViewLevelSelector.tsx on level change |
| `project_exported` | Export completed | Export button in project settings |
| `code_edited` | Direct code edit saved | Future: inline code editor |

### Example Integration: FilePreviewModal

```tsx
// In FilePreviewModal.tsx
import { useAchievements } from '@/hooks/useAchievements';

export function FilePreviewModal({ projectId, file, ... }) {
  const { recordEvent } = useAchievements(projectId);

  useEffect(() => {
    // Record code view event when modal opens with code visible
    if (showCode && file) {
      recordEvent({
        type: 'code_viewed',
        context: {
          fileId: file.id,
          fileName: file.name
        }
      });
    }
  }, [showCode, file, recordEvent]);

  // ... rest of component
}
```

---

## Implementation Order

### Phase 3.1: Database & Models (1 day)
1. Create migration `007_learning_journey.sql`
2. Run migration
3. Implement `model/achievement.go`

### Phase 3.2: Repository Layer (1 day)
1. Implement `repository/achievement.go`
2. Add tests for repository methods

### Phase 3.3: Service Layer (2 days)
1. Implement `AchievementService`
2. Implement `NudgeService`
3. Add unit tests for trigger evaluation logic

### Phase 3.4: REST API (1 day)
1. Implement `handler/achievement.go`
2. Register routes in main.go
3. Add API tests

### Phase 3.5: Frontend Types & Hook (1 day)
1. Create `types/achievements.ts`
2. Implement `useAchievements` hook
3. Test API integration

### Phase 3.6: Frontend Components (2 days)
1. Implement `AchievementToast`
2. Implement `ProgressBadge`
3. Implement `NudgePopover`
4. Integrate with existing components

### Phase 3.7: Event Integration (1 day)
1. Add event triggers to FilePreviewModal
2. Add event triggers to FileTree
3. Add event triggers to level selector
4. End-to-end testing

**Total Estimated Time: 9 days**

---

## Architecture Decision Records

### ADR-001: Event-Driven Achievement Triggers

**Context**: Achievements need to be triggered by various user actions across the application.

**Decision**: Use an event-driven pattern where UI components call `recordEvent()` with typed events. The backend evaluates all active achievement triggers on each event.

**Consequences**:
- (+) Decoupled: UI doesn't need to know achievement logic
- (+) Extensible: New achievements can be added without UI changes
- (+) Centralized: All trigger logic in one service
- (-) Latency: Each event requires API call and evaluation
- (-) Complexity: Need to track which events map to which triggers

**Mitigation**: Events are fire-and-forget (optimistic UI), evaluation is fast (in-memory trigger configs).

### ADR-002: Project-Scoped Progress

**Context**: Should progress be per-user globally or per-project?

**Decision**: Progress is per-project. Each project has its own progress state and achievements.

**Rationale**:
- Users may work on multiple projects with different learning needs
- Avoids complexity of user authentication (current MVP has no auth)
- Matches mental model: "learning about THIS project"
- Easy to aggregate across projects later if needed

### ADR-003: Nudge Frequency Control

**Context**: How to prevent nudges from becoming annoying?

**Decision**:
- Each nudge type shown at most once per project
- Nudges only appear when specific conditions are met
- User can dismiss nudges, which records the dismissal
- No nudges for 5 minutes after any nudge interaction

**Consequences**:
- (+) Non-intrusive experience
- (+) Nudges feel helpful, not nagging
- (-) Some users may miss useful nudges if dismissed early

---

## Success Criteria

### Must Have for Phase 3 Complete
- [ ] Database migration runs successfully
- [ ] At least 5 achievements defined and triggering correctly
- [ ] Progress tracking shows accurate stats
- [ ] Achievement toast displays on unlock
- [ ] Progress badge visible in UI
- [ ] Events triggered from file viewing

### Should Have
- [ ] Nudge system delivering contextual suggestions
- [ ] Level progression tracked and displayed
- [ ] All 10 seed achievements implemented

### Could Have
- [ ] Achievement history/gallery view
- [ ] Points leaderboard across projects
- [ ] Export achievement (VS Code graduation)

---

## Files to Create/Modify

### New Files
| Path | Description |
|------|-------------|
| `/backend/migrations/007_learning_journey.sql` | Database schema |
| `/backend/internal/model/achievement.go` | Data models |
| `/backend/internal/repository/achievement.go` | Data access |
| `/backend/internal/service/achievement.go` | Achievement logic |
| `/backend/internal/service/nudge.go` | Nudge logic |
| `/backend/internal/handler/achievement.go` | REST endpoints |
| `/frontend/src/types/achievements.ts` | TypeScript types |
| `/frontend/src/hooks/useAchievements.ts` | React hook |
| `/frontend/src/components/achievements/AchievementToast.tsx` | Toast component |
| `/frontend/src/components/achievements/ProgressBadge.tsx` | Badge component |
| `/frontend/src/components/achievements/NudgePopover.tsx` | Nudge component |

### Modified Files
| Path | Change |
|------|--------|
| `/backend/cmd/server/main.go` | Wire achievement services and routes |
| `/frontend/src/components/projects/FilePreviewModal.tsx` | Add event triggers |
| `/frontend/src/components/projects/FileTree.tsx` | Add event triggers |
| `/frontend/src/components/projects/ProjectPageClient.tsx` | Add progress badge |

---

**Document End**
