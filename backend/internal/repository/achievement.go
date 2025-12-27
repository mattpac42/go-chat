package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// AchievementRepository defines data access for achievements.
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
	UnlockAchievement(ctx context.Context, projectID uuid.UUID, achievementID uuid.UUID, triggerContext []byte) (*model.UserAchievement, error)
	MarkAchievementSeen(ctx context.Context, id uuid.UUID) error

	// Nudge history
	GetRecentNudges(ctx context.Context, projectID uuid.UUID, limit int) ([]model.NudgeHistory, error)
	HasSeenNudge(ctx context.Context, projectID uuid.UUID, nudgeType model.NudgeType) (bool, error)
	RecordNudgeShown(ctx context.Context, projectID uuid.UUID, nudgeType model.NudgeType) error
	RecordNudgeDismissed(ctx context.Context, id uuid.UUID) error
	RecordNudgeClicked(ctx context.Context, id uuid.UUID) error
}

// achievementRepository implements AchievementRepository using PostgreSQL.
type achievementRepository struct {
	db *sqlx.DB
}

// NewAchievementRepository creates a new AchievementRepository.
func NewAchievementRepository(db *sqlx.DB) AchievementRepository {
	return &achievementRepository{db: db}
}

// ============================================================================
// Achievement Definitions
// ============================================================================

// GetAllActive returns all active achievements.
func (r *achievementRepository) GetAllActive(ctx context.Context) ([]model.Achievement, error) {
	query := `
		SELECT id, code, name, description, category, icon, points, trigger_type,
		       trigger_config, prerequisites, is_active, created_at
		FROM achievements
		WHERE is_active = TRUE
		ORDER BY category, points ASC
	`

	var achievements []model.Achievement
	if err := r.db.SelectContext(ctx, &achievements, query); err != nil {
		return nil, err
	}

	if achievements == nil {
		achievements = []model.Achievement{}
	}

	return achievements, nil
}

// GetByCode returns an achievement by its unique code.
func (r *achievementRepository) GetByCode(ctx context.Context, code string) (*model.Achievement, error) {
	query := `
		SELECT id, code, name, description, category, icon, points, trigger_type,
		       trigger_config, prerequisites, is_active, created_at
		FROM achievements
		WHERE code = $1
	`

	var achievement model.Achievement
	if err := r.db.GetContext(ctx, &achievement, query, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &achievement, nil
}

// GetByCategory returns all achievements in a category.
func (r *achievementRepository) GetByCategory(ctx context.Context, category model.AchievementCategory) ([]model.Achievement, error) {
	query := `
		SELECT id, code, name, description, category, icon, points, trigger_type,
		       trigger_config, prerequisites, is_active, created_at
		FROM achievements
		WHERE category = $1 AND is_active = TRUE
		ORDER BY points ASC
	`

	var achievements []model.Achievement
	if err := r.db.SelectContext(ctx, &achievements, query, category); err != nil {
		return nil, err
	}

	if achievements == nil {
		achievements = []model.Achievement{}
	}

	return achievements, nil
}

// ============================================================================
// User Progress
// ============================================================================

// GetProgress returns the user progress for a project.
func (r *achievementRepository) GetProgress(ctx context.Context, projectID uuid.UUID) (*model.UserProgress, error) {
	query := `
		SELECT id, project_id, current_level, total_points, files_viewed_count,
		       code_views_count, tree_expansions_count, level_changes_count,
		       first_code_view_at, first_level_up_at, last_activity_at, created_at, updated_at
		FROM user_progress
		WHERE project_id = $1
	`

	var progress model.UserProgress
	if err := r.db.GetContext(ctx, &progress, query, projectID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &progress, nil
}

// CreateProgress creates a new user progress record for a project.
func (r *achievementRepository) CreateProgress(ctx context.Context, projectID uuid.UUID) (*model.UserProgress, error) {
	query := `
		INSERT INTO user_progress (project_id, current_level, total_points,
		                           files_viewed_count, code_views_count,
		                           tree_expansions_count, level_changes_count)
		VALUES ($1, 1, 0, 0, 0, 0, 0)
		RETURNING id, project_id, current_level, total_points, files_viewed_count,
		          code_views_count, tree_expansions_count, level_changes_count,
		          first_code_view_at, first_level_up_at, last_activity_at, created_at, updated_at
	`

	var progress model.UserProgress
	if err := r.db.GetContext(ctx, &progress, query, projectID); err != nil {
		return nil, err
	}

	return &progress, nil
}

// UpdateProgress updates an existing user progress record.
func (r *achievementRepository) UpdateProgress(ctx context.Context, progress *model.UserProgress) error {
	query := `
		UPDATE user_progress
		SET current_level = $2,
		    total_points = $3,
		    files_viewed_count = $4,
		    code_views_count = $5,
		    tree_expansions_count = $6,
		    level_changes_count = $7,
		    first_code_view_at = $8,
		    first_level_up_at = $9,
		    last_activity_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		progress.ID,
		progress.CurrentLevel,
		progress.TotalPoints,
		progress.FilesViewedCount,
		progress.CodeViewsCount,
		progress.TreeExpansionsCount,
		progress.LevelChangesCount,
		progress.FirstCodeViewAt,
		progress.FirstLevelUpAt,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// IncrementStat increments a specific stat counter for a project's progress.
// Valid stat names: files_viewed_count, code_views_count, tree_expansions_count, level_changes_count
func (r *achievementRepository) IncrementStat(ctx context.Context, projectID uuid.UUID, stat string, delta int) error {
	// Validate stat name to prevent SQL injection
	validStats := map[string]bool{
		"files_viewed_count":    true,
		"code_views_count":      true,
		"tree_expansions_count": true,
		"level_changes_count":   true,
	}

	if !validStats[stat] {
		return errors.New("invalid stat name")
	}

	// Use a parameterized query with the stat name embedded (safe due to whitelist)
	query := `
		UPDATE user_progress
		SET ` + stat + ` = ` + stat + ` + $2,
		    last_activity_at = NOW(),
		    updated_at = NOW()
		WHERE project_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, projectID, delta)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// ============================================================================
// User Achievements
// ============================================================================

// GetUserAchievements returns all achievements unlocked by a project.
func (r *achievementRepository) GetUserAchievements(ctx context.Context, projectID uuid.UUID) ([]model.UserAchievement, error) {
	query := `
		SELECT id, project_id, achievement_id, unlocked_at, trigger_context, is_seen, seen_at
		FROM user_achievements
		WHERE project_id = $1
		ORDER BY unlocked_at DESC
	`

	var achievements []model.UserAchievement
	if err := r.db.SelectContext(ctx, &achievements, query, projectID); err != nil {
		return nil, err
	}

	if achievements == nil {
		achievements = []model.UserAchievement{}
	}

	return achievements, nil
}

// GetUnseenAchievements returns achievements that haven't been seen by the user.
func (r *achievementRepository) GetUnseenAchievements(ctx context.Context, projectID uuid.UUID) ([]model.UserAchievement, error) {
	query := `
		SELECT id, project_id, achievement_id, unlocked_at, trigger_context, is_seen, seen_at
		FROM user_achievements
		WHERE project_id = $1 AND is_seen = FALSE
		ORDER BY unlocked_at ASC
	`

	var achievements []model.UserAchievement
	if err := r.db.SelectContext(ctx, &achievements, query, projectID); err != nil {
		return nil, err
	}

	if achievements == nil {
		achievements = []model.UserAchievement{}
	}

	return achievements, nil
}

// HasAchievement checks if a project has unlocked a specific achievement.
func (r *achievementRepository) HasAchievement(ctx context.Context, projectID uuid.UUID, achievementCode string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM user_achievements ua
			JOIN achievements a ON ua.achievement_id = a.id
			WHERE ua.project_id = $1 AND a.code = $2
		)
	`

	var exists bool
	if err := r.db.GetContext(ctx, &exists, query, projectID, achievementCode); err != nil {
		return false, err
	}

	return exists, nil
}

// UnlockAchievement records a new achievement unlock for a project.
func (r *achievementRepository) UnlockAchievement(ctx context.Context, projectID uuid.UUID, achievementID uuid.UUID, triggerContext []byte) (*model.UserAchievement, error) {
	// Default to empty JSON object if no context provided
	if triggerContext == nil {
		triggerContext = []byte("{}")
	}

	query := `
		INSERT INTO user_achievements (project_id, achievement_id, trigger_context)
		VALUES ($1, $2, $3)
		ON CONFLICT (project_id, achievement_id) DO NOTHING
		RETURNING id, project_id, achievement_id, unlocked_at, trigger_context, is_seen, seen_at
	`

	var achievement model.UserAchievement
	if err := r.db.GetContext(ctx, &achievement, query, projectID, achievementID, triggerContext); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Achievement already exists, return nil without error
			return nil, nil
		}
		return nil, err
	}

	return &achievement, nil
}

// MarkAchievementSeen marks an achievement as seen by the user.
func (r *achievementRepository) MarkAchievementSeen(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE user_achievements
		SET is_seen = TRUE, seen_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// ============================================================================
// Nudge History
// ============================================================================

// GetRecentNudges returns the most recent nudges shown to a project.
func (r *achievementRepository) GetRecentNudges(ctx context.Context, projectID uuid.UUID, limit int) ([]model.NudgeHistory, error) {
	query := `
		SELECT id, project_id, nudge_type, shown_at, dismissed_at, clicked_at, context
		FROM nudge_history
		WHERE project_id = $1
		ORDER BY shown_at DESC
		LIMIT $2
	`

	var nudges []model.NudgeHistory
	if err := r.db.SelectContext(ctx, &nudges, query, projectID, limit); err != nil {
		return nil, err
	}

	if nudges == nil {
		nudges = []model.NudgeHistory{}
	}

	return nudges, nil
}

// HasSeenNudge checks if a specific nudge type has been shown to a project.
func (r *achievementRepository) HasSeenNudge(ctx context.Context, projectID uuid.UUID, nudgeType model.NudgeType) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM nudge_history
			WHERE project_id = $1 AND nudge_type = $2
		)
	`

	var exists bool
	if err := r.db.GetContext(ctx, &exists, query, projectID, nudgeType); err != nil {
		return false, err
	}

	return exists, nil
}

// RecordNudgeShown records that a nudge was shown to the user.
func (r *achievementRepository) RecordNudgeShown(ctx context.Context, projectID uuid.UUID, nudgeType model.NudgeType) error {
	query := `
		INSERT INTO nudge_history (project_id, nudge_type)
		VALUES ($1, $2)
	`

	_, err := r.db.ExecContext(ctx, query, projectID, nudgeType)
	return err
}

// RecordNudgeDismissed records that a nudge was dismissed by the user.
func (r *achievementRepository) RecordNudgeDismissed(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE nudge_history
		SET dismissed_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// RecordNudgeClicked records that a nudge was clicked/followed by the user.
func (r *achievementRepository) RecordNudgeClicked(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE nudge_history
		SET clicked_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
