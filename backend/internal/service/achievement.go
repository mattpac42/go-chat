package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

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

			// Handle nil return (achievement was already unlocked)
			if userAchievement == nil {
				continue
			}

			// Add points
			progress.TotalPoints += achievement.Points
			if err := s.repo.UpdateProgress(ctx, progress); err != nil {
				s.logger.Warn().Err(err).Msg("failed to update progress with points")
			}

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

// evaluateMilestoneTrigger checks milestone-based triggers (complex conditions)
func (s *AchievementService) evaluateMilestoneTrigger(event LearningEvent, progress *model.UserProgress, config map[string]interface{}) bool {
	// Check for required level
	if requiredLevel, ok := config["required_level"].(float64); ok {
		if int(progress.CurrentLevel) < int(requiredLevel) {
			return false
		}
	}

	// Check for minimum points
	if minPoints, ok := config["min_points"].(float64); ok {
		if progress.TotalPoints < int(minPoints) {
			return false
		}
	}

	// Check for required event type
	if requiredEvent, ok := config["event"].(string); ok {
		if event.Type != requiredEvent {
			return false
		}
	}

	// Check for multiple metrics combined
	if conditions, ok := config["conditions"].([]interface{}); ok {
		for _, cond := range conditions {
			if condMap, ok := cond.(map[string]interface{}); ok {
				if !s.evaluateCountTrigger(progress, condMap) {
					return false
				}
			}
		}
	}

	return true
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

	// Track first level up time
	if progress.FirstLevelUpAt == nil && newLevel > model.LevelFunctional {
		now := time.Now()
		progress.FirstLevelUpAt = &now
	}

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

	_, _ = s.ProcessEvent(ctx, event)
	return nil
}

// getOrCreateProgress retrieves or creates a UserProgress for a project
func (s *AchievementService) getOrCreateProgress(ctx context.Context, projectID uuid.UUID) (*model.UserProgress, error) {
	progress, err := s.repo.GetProgress(ctx, projectID)
	if err == nil {
		return progress, nil
	}

	// If not found, create new progress
	if errors.Is(err, repository.ErrNotFound) {
		progress, err = s.repo.CreateProgress(ctx, projectID)
		if err != nil {
			return nil, err
		}
		return progress, nil
	}

	return nil, err
}

// updateStats updates progress statistics based on the event type
func (s *AchievementService) updateStats(ctx context.Context, event LearningEvent, progress *model.UserProgress) error {
	switch event.Type {
	case EventCodeViewed:
		progress.CodeViewsCount++
		// Track first code view
		if progress.FirstCodeViewAt == nil {
			now := time.Now()
			progress.FirstCodeViewAt = &now
		}
		// Track unique files viewed
		if fileID, ok := event.Context["fileId"].(string); ok && fileID != "" {
			progress.FilesViewedCount++
		}

	case EventTreeExpanded:
		progress.TreeExpansionsCount++

	case EventFullTreeView:
		progress.TreeExpansionsCount++

	case EventLevelChanged:
		// Level changes are tracked in UpdateLevel method
		return nil

	case EventProjectExport:
		// No specific stat for exports currently

	case EventCodeEdited:
		// Could track edits in the future
	}

	return s.repo.UpdateProgress(ctx, progress)
}

// checkPrerequisites verifies all prerequisite achievements are unlocked
func (s *AchievementService) checkPrerequisites(ctx context.Context, projectID uuid.UUID, prerequisites []string) bool {
	if len(prerequisites) == 0 {
		return true
	}

	for _, prereqCode := range prerequisites {
		hasIt, err := s.repo.HasAchievement(ctx, projectID, prereqCode)
		if err != nil {
			s.logger.Warn().Err(err).Str("prereq", prereqCode).Msg("failed to check prerequisite")
			return false
		}
		if !hasIt {
			return false
		}
	}

	return true
}
