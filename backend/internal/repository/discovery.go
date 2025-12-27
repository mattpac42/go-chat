package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// DiscoveryRepository defines the interface for discovery data access.
type DiscoveryRepository interface {
	// GetByProjectID retrieves the discovery state for a project.
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*model.ProjectDiscovery, error)

	// GetByID retrieves a discovery by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*model.ProjectDiscovery, error)

	// Create creates a new discovery for a project.
	Create(ctx context.Context, projectID uuid.UUID) (*model.ProjectDiscovery, error)

	// Update updates discovery data fields.
	Update(ctx context.Context, discovery *model.ProjectDiscovery) (*model.ProjectDiscovery, error)

	// UpdateStage advances the discovery to a new stage.
	UpdateStage(ctx context.Context, discoveryID uuid.UUID, stage model.DiscoveryStage) (*model.ProjectDiscovery, error)

	// MarkComplete marks the discovery as complete with confirmation timestamp.
	MarkComplete(ctx context.Context, discoveryID uuid.UUID) (*model.ProjectDiscovery, error)

	// Delete removes a discovery and all associated data.
	Delete(ctx context.Context, discoveryID uuid.UUID) error

	// User persona methods
	AddUser(ctx context.Context, user *model.DiscoveryUser) (*model.DiscoveryUser, error)
	GetUsers(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryUser, error)
	UpdateUser(ctx context.Context, user *model.DiscoveryUser) (*model.DiscoveryUser, error)
	ClearUsers(ctx context.Context, discoveryID uuid.UUID) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	// Feature methods
	AddFeature(ctx context.Context, feature *model.DiscoveryFeature) (*model.DiscoveryFeature, error)
	GetFeatures(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryFeature, error)
	GetMVPFeatures(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryFeature, error)
	GetFutureFeatures(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryFeature, error)
	UpdateFeature(ctx context.Context, feature *model.DiscoveryFeature) (*model.DiscoveryFeature, error)
	DeleteFeature(ctx context.Context, featureID uuid.UUID) error

	// Edit history methods
	AddEditHistory(ctx context.Context, history *model.DiscoveryEditHistory) (*model.DiscoveryEditHistory, error)
	GetEditHistory(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryEditHistory, error)

	// Summary method
	GetSummary(ctx context.Context, discoveryID uuid.UUID) (*model.DiscoverySummary, error)
}

// PostgresDiscoveryRepository implements DiscoveryRepository using PostgreSQL.
type PostgresDiscoveryRepository struct {
	db *sqlx.DB
}

// NewPostgresDiscoveryRepository creates a new PostgresDiscoveryRepository.
func NewPostgresDiscoveryRepository(db *sqlx.DB) *PostgresDiscoveryRepository {
	return &PostgresDiscoveryRepository{db: db}
}

// GetByProjectID retrieves the discovery state for a project.
func (r *PostgresDiscoveryRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*model.ProjectDiscovery, error) {
	query := `
		SELECT id, project_id, stage, stage_started_at, business_context, problem_statement,
		       goals, project_name, solves_statement, is_returning_user, used_template_id,
		       confirmed_at, created_at, updated_at
		FROM project_discovery
		WHERE project_id = $1
	`

	var discovery model.ProjectDiscovery
	if err := r.db.GetContext(ctx, &discovery, query, projectID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &discovery, nil
}

// GetByID retrieves a discovery by its ID.
func (r *PostgresDiscoveryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ProjectDiscovery, error) {
	query := `
		SELECT id, project_id, stage, stage_started_at, business_context, problem_statement,
		       goals, project_name, solves_statement, is_returning_user, used_template_id,
		       confirmed_at, created_at, updated_at
		FROM project_discovery
		WHERE id = $1
	`

	var discovery model.ProjectDiscovery
	if err := r.db.GetContext(ctx, &discovery, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &discovery, nil
}

// Create creates a new discovery for a project.
func (r *PostgresDiscoveryRepository) Create(ctx context.Context, projectID uuid.UUID) (*model.ProjectDiscovery, error) {
	query := `
		INSERT INTO project_discovery (project_id, stage, stage_started_at)
		VALUES ($1, $2, NOW())
		RETURNING id, project_id, stage, stage_started_at, business_context, problem_statement,
		          goals, project_name, solves_statement, is_returning_user, used_template_id,
		          confirmed_at, created_at, updated_at
	`

	var discovery model.ProjectDiscovery
	if err := r.db.GetContext(ctx, &discovery, query, projectID, model.StageWelcome); err != nil {
		return nil, err
	}

	return &discovery, nil
}

// Update updates discovery data fields.
func (r *PostgresDiscoveryRepository) Update(ctx context.Context, discovery *model.ProjectDiscovery) (*model.ProjectDiscovery, error) {
	// Marshal goals to JSON
	goalsJSON, err := json.Marshal([]string{})
	if discovery.GoalsJSON != nil && len(discovery.GoalsJSON) > 0 {
		goalsJSON = discovery.GoalsJSON
	}

	query := `
		UPDATE project_discovery
		SET business_context = $2,
		    problem_statement = $3,
		    goals = $4,
		    project_name = $5,
		    solves_statement = $6,
		    is_returning_user = $7,
		    used_template_id = $8,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, project_id, stage, stage_started_at, business_context, problem_statement,
		          goals, project_name, solves_statement, is_returning_user, used_template_id,
		          confirmed_at, created_at, updated_at
	`

	var updated model.ProjectDiscovery
	if err = r.db.GetContext(ctx, &updated, query,
		discovery.ID,
		discovery.BusinessContext,
		discovery.ProblemStatement,
		goalsJSON,
		discovery.ProjectName,
		discovery.SolvesStatement,
		discovery.IsReturningUser,
		discovery.UsedTemplateID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &updated, nil
}

// UpdateStage advances the discovery to a new stage.
func (r *PostgresDiscoveryRepository) UpdateStage(ctx context.Context, discoveryID uuid.UUID, stage model.DiscoveryStage) (*model.ProjectDiscovery, error) {
	query := `
		UPDATE project_discovery
		SET stage = $2,
		    stage_started_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, project_id, stage, stage_started_at, business_context, problem_statement,
		          goals, project_name, solves_statement, is_returning_user, used_template_id,
		          confirmed_at, created_at, updated_at
	`

	var discovery model.ProjectDiscovery
	if err := r.db.GetContext(ctx, &discovery, query, discoveryID, stage); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &discovery, nil
}

// MarkComplete marks the discovery as complete with confirmation timestamp.
func (r *PostgresDiscoveryRepository) MarkComplete(ctx context.Context, discoveryID uuid.UUID) (*model.ProjectDiscovery, error) {
	query := `
		UPDATE project_discovery
		SET stage = $2,
		    confirmed_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, project_id, stage, stage_started_at, business_context, problem_statement,
		          goals, project_name, solves_statement, is_returning_user, used_template_id,
		          confirmed_at, created_at, updated_at
	`

	var discovery model.ProjectDiscovery
	if err := r.db.GetContext(ctx, &discovery, query, discoveryID, model.StageComplete); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &discovery, nil
}

// Delete removes a discovery and all associated data.
func (r *PostgresDiscoveryRepository) Delete(ctx context.Context, discoveryID uuid.UUID) error {
	query := `DELETE FROM project_discovery WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, discoveryID)
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

// AddUser adds a new user persona to a discovery.
func (r *PostgresDiscoveryRepository) AddUser(ctx context.Context, user *model.DiscoveryUser) (*model.DiscoveryUser, error) {
	query := `
		INSERT INTO discovery_users (discovery_id, description, user_count, has_permissions, permission_notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, discovery_id, description, user_count, has_permissions, permission_notes, created_at
	`

	var created model.DiscoveryUser
	if err := r.db.GetContext(ctx, &created, query,
		user.DiscoveryID,
		user.Description,
		user.UserCount,
		user.HasPermissions,
		user.PermissionNotes,
	); err != nil {
		return nil, err
	}

	return &created, nil
}

// GetUsers retrieves all user personas for a discovery.
func (r *PostgresDiscoveryRepository) GetUsers(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryUser, error) {
	query := `
		SELECT id, discovery_id, description, user_count, has_permissions, permission_notes, created_at
		FROM discovery_users
		WHERE discovery_id = $1
		ORDER BY created_at ASC
	`

	var users []model.DiscoveryUser
	if err := r.db.SelectContext(ctx, &users, query, discoveryID); err != nil {
		return nil, err
	}

	if users == nil {
		users = []model.DiscoveryUser{}
	}

	return users, nil
}

// UpdateUser updates a user persona.
func (r *PostgresDiscoveryRepository) UpdateUser(ctx context.Context, user *model.DiscoveryUser) (*model.DiscoveryUser, error) {
	query := `
		UPDATE discovery_users
		SET description = $2,
		    user_count = $3,
		    has_permissions = $4,
		    permission_notes = $5
		WHERE id = $1
		RETURNING id, discovery_id, description, user_count, has_permissions, permission_notes, created_at
	`

	var updated model.DiscoveryUser
	if err := r.db.GetContext(ctx, &updated, query,
		user.ID,
		user.Description,
		user.UserCount,
		user.HasPermissions,
		user.PermissionNotes,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &updated, nil
}

// DeleteUser removes a user persona.
func (r *PostgresDiscoveryRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM discovery_users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
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

// ClearUsers removes all users for a discovery.
func (r *PostgresDiscoveryRepository) ClearUsers(ctx context.Context, discoveryID uuid.UUID) error {
	query := `DELETE FROM discovery_users WHERE discovery_id = $1`
	_, err := r.db.ExecContext(ctx, query, discoveryID)
	return err
}

// AddFeature adds a new feature to a discovery.
func (r *PostgresDiscoveryRepository) AddFeature(ctx context.Context, feature *model.DiscoveryFeature) (*model.DiscoveryFeature, error) {
	version := feature.Version
	if version == "" {
		version = "v1"
	}

	query := `
		INSERT INTO discovery_features (discovery_id, name, priority, version)
		VALUES ($1, $2, $3, $4)
		RETURNING id, discovery_id, name, priority, version, created_at
	`

	var created model.DiscoveryFeature
	if err := r.db.GetContext(ctx, &created, query,
		feature.DiscoveryID,
		feature.Name,
		feature.Priority,
		version,
	); err != nil {
		return nil, err
	}

	return &created, nil
}

// GetFeatures retrieves all features for a discovery.
func (r *PostgresDiscoveryRepository) GetFeatures(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryFeature, error) {
	query := `
		SELECT id, discovery_id, name, priority, version, created_at
		FROM discovery_features
		WHERE discovery_id = $1
		ORDER BY version ASC, priority ASC
	`

	var features []model.DiscoveryFeature
	if err := r.db.SelectContext(ctx, &features, query, discoveryID); err != nil {
		return nil, err
	}

	if features == nil {
		features = []model.DiscoveryFeature{}
	}

	return features, nil
}

// GetMVPFeatures retrieves MVP (v1) features for a discovery.
func (r *PostgresDiscoveryRepository) GetMVPFeatures(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryFeature, error) {
	query := `
		SELECT id, discovery_id, name, priority, version, created_at
		FROM discovery_features
		WHERE discovery_id = $1 AND version = 'v1'
		ORDER BY priority ASC
	`

	var features []model.DiscoveryFeature
	if err := r.db.SelectContext(ctx, &features, query, discoveryID); err != nil {
		return nil, err
	}

	if features == nil {
		features = []model.DiscoveryFeature{}
	}

	return features, nil
}

// GetFutureFeatures retrieves future (non-v1) features for a discovery.
func (r *PostgresDiscoveryRepository) GetFutureFeatures(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryFeature, error) {
	query := `
		SELECT id, discovery_id, name, priority, version, created_at
		FROM discovery_features
		WHERE discovery_id = $1 AND version != 'v1'
		ORDER BY version ASC, priority ASC
	`

	var features []model.DiscoveryFeature
	if err := r.db.SelectContext(ctx, &features, query, discoveryID); err != nil {
		return nil, err
	}

	if features == nil {
		features = []model.DiscoveryFeature{}
	}

	return features, nil
}

// UpdateFeature updates a feature.
func (r *PostgresDiscoveryRepository) UpdateFeature(ctx context.Context, feature *model.DiscoveryFeature) (*model.DiscoveryFeature, error) {
	query := `
		UPDATE discovery_features
		SET name = $2,
		    priority = $3,
		    version = $4
		WHERE id = $1
		RETURNING id, discovery_id, name, priority, version, created_at
	`

	var updated model.DiscoveryFeature
	if err := r.db.GetContext(ctx, &updated, query,
		feature.ID,
		feature.Name,
		feature.Priority,
		feature.Version,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &updated, nil
}

// DeleteFeature removes a feature.
func (r *PostgresDiscoveryRepository) DeleteFeature(ctx context.Context, featureID uuid.UUID) error {
	query := `DELETE FROM discovery_features WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, featureID)
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

// AddEditHistory adds an edit history entry.
func (r *PostgresDiscoveryRepository) AddEditHistory(ctx context.Context, history *model.DiscoveryEditHistory) (*model.DiscoveryEditHistory, error) {
	query := `
		INSERT INTO discovery_edit_history (discovery_id, stage, field_edited, original_value, new_value, edited_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, discovery_id, stage, field_edited, original_value, new_value, edited_at
	`

	var created model.DiscoveryEditHistory
	if err := r.db.GetContext(ctx, &created, query,
		history.DiscoveryID,
		history.Stage,
		history.FieldEdited,
		history.OriginalValue,
		history.NewValue,
	); err != nil {
		return nil, err
	}

	return &created, nil
}

// GetEditHistory retrieves edit history for a discovery.
func (r *PostgresDiscoveryRepository) GetEditHistory(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryEditHistory, error) {
	query := `
		SELECT id, discovery_id, stage, field_edited, original_value, new_value, edited_at
		FROM discovery_edit_history
		WHERE discovery_id = $1
		ORDER BY edited_at DESC
	`

	var history []model.DiscoveryEditHistory
	if err := r.db.SelectContext(ctx, &history, query, discoveryID); err != nil {
		return nil, err
	}

	if history == nil {
		history = []model.DiscoveryEditHistory{}
	}

	return history, nil
}

// GetSummary builds the complete summary view for a discovery.
func (r *PostgresDiscoveryRepository) GetSummary(ctx context.Context, discoveryID uuid.UUID) (*model.DiscoverySummary, error) {
	// Get the discovery
	discovery, err := r.GetByID(ctx, discoveryID)
	if err != nil {
		return nil, err
	}

	// Get users
	users, err := r.GetUsers(ctx, discoveryID)
	if err != nil {
		return nil, err
	}

	// Get MVP features
	mvpFeatures, err := r.GetMVPFeatures(ctx, discoveryID)
	if err != nil {
		return nil, err
	}

	// Get future features
	futureFeatures, err := r.GetFutureFeatures(ctx, discoveryID)
	if err != nil {
		return nil, err
	}

	// Build summary
	projectName := ""
	if discovery.ProjectName != nil {
		projectName = *discovery.ProjectName
	}

	solvesStatement := ""
	if discovery.SolvesStatement != nil {
		solvesStatement = *discovery.SolvesStatement
	}

	return &model.DiscoverySummary{
		ProjectName:     projectName,
		SolvesStatement: solvesStatement,
		Users:           users,
		MVPFeatures:     mvpFeatures,
		FutureFeatures:  futureFeatures,
	}, nil
}

// MockDiscoveryRepository implements DiscoveryRepository for testing.
type MockDiscoveryRepository struct {
	discoveries map[uuid.UUID]*model.ProjectDiscovery
	byProject   map[uuid.UUID]uuid.UUID // projectID -> discoveryID
	users       map[uuid.UUID][]model.DiscoveryUser
	features    map[uuid.UUID][]model.DiscoveryFeature
	editHistory map[uuid.UUID][]model.DiscoveryEditHistory
}

// NewMockDiscoveryRepository creates a new MockDiscoveryRepository.
func NewMockDiscoveryRepository() *MockDiscoveryRepository {
	return &MockDiscoveryRepository{
		discoveries: make(map[uuid.UUID]*model.ProjectDiscovery),
		byProject:   make(map[uuid.UUID]uuid.UUID),
		users:       make(map[uuid.UUID][]model.DiscoveryUser),
		features:    make(map[uuid.UUID][]model.DiscoveryFeature),
		editHistory: make(map[uuid.UUID][]model.DiscoveryEditHistory),
	}
}

// GetByProjectID retrieves the discovery state for a project.
func (r *MockDiscoveryRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*model.ProjectDiscovery, error) {
	discoveryID, ok := r.byProject[projectID]
	if !ok {
		return nil, ErrNotFound
	}
	return r.GetByID(ctx, discoveryID)
}

// GetByID retrieves a discovery by its ID.
func (r *MockDiscoveryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ProjectDiscovery, error) {
	discovery, ok := r.discoveries[id]
	if !ok {
		return nil, ErrNotFound
	}
	// Return a copy
	copy := *discovery
	return &copy, nil
}

// Create creates a new discovery for a project.
func (r *MockDiscoveryRepository) Create(ctx context.Context, projectID uuid.UUID) (*model.ProjectDiscovery, error) {
	// Check if discovery already exists for project
	if _, ok := r.byProject[projectID]; ok {
		return nil, errors.New("discovery already exists for project")
	}

	now := time.Now().UTC()
	discovery := &model.ProjectDiscovery{
		ID:             uuid.New(),
		ProjectID:      projectID,
		Stage:          model.StageWelcome,
		StageStartedAt: now,
		GoalsJSON:      []byte("[]"),
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	r.discoveries[discovery.ID] = discovery
	r.byProject[projectID] = discovery.ID

	copy := *discovery
	return &copy, nil
}

// Update updates discovery data fields.
func (r *MockDiscoveryRepository) Update(ctx context.Context, discovery *model.ProjectDiscovery) (*model.ProjectDiscovery, error) {
	existing, ok := r.discoveries[discovery.ID]
	if !ok {
		return nil, ErrNotFound
	}

	existing.BusinessContext = discovery.BusinessContext
	existing.ProblemStatement = discovery.ProblemStatement
	existing.GoalsJSON = discovery.GoalsJSON
	existing.ProjectName = discovery.ProjectName
	existing.SolvesStatement = discovery.SolvesStatement
	existing.IsReturningUser = discovery.IsReturningUser
	existing.UsedTemplateID = discovery.UsedTemplateID
	existing.UpdatedAt = time.Now().UTC()

	copy := *existing
	return &copy, nil
}

// UpdateStage advances the discovery to a new stage.
func (r *MockDiscoveryRepository) UpdateStage(ctx context.Context, discoveryID uuid.UUID, stage model.DiscoveryStage) (*model.ProjectDiscovery, error) {
	discovery, ok := r.discoveries[discoveryID]
	if !ok {
		return nil, ErrNotFound
	}

	discovery.Stage = stage
	discovery.StageStartedAt = time.Now().UTC()
	discovery.UpdatedAt = time.Now().UTC()

	copy := *discovery
	return &copy, nil
}

// MarkComplete marks the discovery as complete with confirmation timestamp.
func (r *MockDiscoveryRepository) MarkComplete(ctx context.Context, discoveryID uuid.UUID) (*model.ProjectDiscovery, error) {
	discovery, ok := r.discoveries[discoveryID]
	if !ok {
		return nil, ErrNotFound
	}

	now := time.Now().UTC()
	discovery.Stage = model.StageComplete
	discovery.ConfirmedAt = &now
	discovery.UpdatedAt = now

	copy := *discovery
	return &copy, nil
}

// Delete removes a discovery and all associated data.
func (r *MockDiscoveryRepository) Delete(ctx context.Context, discoveryID uuid.UUID) error {
	discovery, ok := r.discoveries[discoveryID]
	if !ok {
		return ErrNotFound
	}

	delete(r.discoveries, discoveryID)
	delete(r.byProject, discovery.ProjectID)
	delete(r.users, discoveryID)
	delete(r.features, discoveryID)
	delete(r.editHistory, discoveryID)

	return nil
}

// AddUser adds a new user persona to a discovery.
func (r *MockDiscoveryRepository) AddUser(ctx context.Context, user *model.DiscoveryUser) (*model.DiscoveryUser, error) {
	if _, ok := r.discoveries[user.DiscoveryID]; !ok {
		return nil, ErrNotFound
	}

	created := model.DiscoveryUser{
		ID:              uuid.New(),
		DiscoveryID:     user.DiscoveryID,
		Description:     user.Description,
		UserCount:       user.UserCount,
		HasPermissions:  user.HasPermissions,
		PermissionNotes: user.PermissionNotes,
		CreatedAt:       time.Now().UTC(),
	}

	r.users[user.DiscoveryID] = append(r.users[user.DiscoveryID], created)
	return &created, nil
}

// GetUsers retrieves all user personas for a discovery.
func (r *MockDiscoveryRepository) GetUsers(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryUser, error) {
	users := r.users[discoveryID]
	if users == nil {
		return []model.DiscoveryUser{}, nil
	}
	// Return copy
	result := make([]model.DiscoveryUser, len(users))
	copy(result, users)
	return result, nil
}

// UpdateUser updates a user persona.
func (r *MockDiscoveryRepository) UpdateUser(ctx context.Context, user *model.DiscoveryUser) (*model.DiscoveryUser, error) {
	users := r.users[user.DiscoveryID]
	for i, u := range users {
		if u.ID == user.ID {
			users[i].Description = user.Description
			users[i].UserCount = user.UserCount
			users[i].HasPermissions = user.HasPermissions
			users[i].PermissionNotes = user.PermissionNotes
			return &users[i], nil
		}
	}
	return nil, ErrNotFound
}

// DeleteUser removes a user persona.
func (r *MockDiscoveryRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	for discoveryID, users := range r.users {
		for i, u := range users {
			if u.ID == userID {
				r.users[discoveryID] = append(users[:i], users[i+1:]...)
				return nil
			}
		}
	}
	return ErrNotFound
}

// ClearUsers removes all users for a discovery.
func (r *MockDiscoveryRepository) ClearUsers(ctx context.Context, discoveryID uuid.UUID) error {
	delete(r.users, discoveryID)
	return nil
}

// AddFeature adds a new feature to a discovery.
func (r *MockDiscoveryRepository) AddFeature(ctx context.Context, feature *model.DiscoveryFeature) (*model.DiscoveryFeature, error) {
	if _, ok := r.discoveries[feature.DiscoveryID]; !ok {
		return nil, ErrNotFound
	}

	version := feature.Version
	if version == "" {
		version = "v1"
	}

	created := model.DiscoveryFeature{
		ID:          uuid.New(),
		DiscoveryID: feature.DiscoveryID,
		Name:        feature.Name,
		Priority:    feature.Priority,
		Version:     version,
		CreatedAt:   time.Now().UTC(),
	}

	r.features[feature.DiscoveryID] = append(r.features[feature.DiscoveryID], created)
	return &created, nil
}

// GetFeatures retrieves all features for a discovery.
func (r *MockDiscoveryRepository) GetFeatures(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryFeature, error) {
	features := r.features[discoveryID]
	if features == nil {
		return []model.DiscoveryFeature{}, nil
	}
	result := make([]model.DiscoveryFeature, len(features))
	copy(result, features)
	return result, nil
}

// GetMVPFeatures retrieves MVP (v1) features for a discovery.
func (r *MockDiscoveryRepository) GetMVPFeatures(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryFeature, error) {
	features := r.features[discoveryID]
	var mvp []model.DiscoveryFeature
	for _, f := range features {
		if f.Version == "v1" {
			mvp = append(mvp, f)
		}
	}
	if mvp == nil {
		mvp = []model.DiscoveryFeature{}
	}
	return mvp, nil
}

// GetFutureFeatures retrieves future (non-v1) features for a discovery.
func (r *MockDiscoveryRepository) GetFutureFeatures(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryFeature, error) {
	features := r.features[discoveryID]
	var future []model.DiscoveryFeature
	for _, f := range features {
		if f.Version != "v1" {
			future = append(future, f)
		}
	}
	if future == nil {
		future = []model.DiscoveryFeature{}
	}
	return future, nil
}

// UpdateFeature updates a feature.
func (r *MockDiscoveryRepository) UpdateFeature(ctx context.Context, feature *model.DiscoveryFeature) (*model.DiscoveryFeature, error) {
	features := r.features[feature.DiscoveryID]
	for i, f := range features {
		if f.ID == feature.ID {
			features[i].Name = feature.Name
			features[i].Priority = feature.Priority
			features[i].Version = feature.Version
			return &features[i], nil
		}
	}
	return nil, ErrNotFound
}

// DeleteFeature removes a feature.
func (r *MockDiscoveryRepository) DeleteFeature(ctx context.Context, featureID uuid.UUID) error {
	for discoveryID, features := range r.features {
		for i, f := range features {
			if f.ID == featureID {
				r.features[discoveryID] = append(features[:i], features[i+1:]...)
				return nil
			}
		}
	}
	return ErrNotFound
}

// AddEditHistory adds an edit history entry.
func (r *MockDiscoveryRepository) AddEditHistory(ctx context.Context, history *model.DiscoveryEditHistory) (*model.DiscoveryEditHistory, error) {
	if _, ok := r.discoveries[history.DiscoveryID]; !ok {
		return nil, ErrNotFound
	}

	created := model.DiscoveryEditHistory{
		ID:            uuid.New(),
		DiscoveryID:   history.DiscoveryID,
		Stage:         history.Stage,
		FieldEdited:   history.FieldEdited,
		OriginalValue: history.OriginalValue,
		NewValue:      history.NewValue,
		EditedAt:      time.Now().UTC(),
	}

	r.editHistory[history.DiscoveryID] = append(r.editHistory[history.DiscoveryID], created)
	return &created, nil
}

// GetEditHistory retrieves edit history for a discovery.
func (r *MockDiscoveryRepository) GetEditHistory(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryEditHistory, error) {
	history := r.editHistory[discoveryID]
	if history == nil {
		return []model.DiscoveryEditHistory{}, nil
	}
	result := make([]model.DiscoveryEditHistory, len(history))
	copy(result, history)
	return result, nil
}

// GetSummary builds the complete summary view for a discovery.
func (r *MockDiscoveryRepository) GetSummary(ctx context.Context, discoveryID uuid.UUID) (*model.DiscoverySummary, error) {
	discovery, err := r.GetByID(ctx, discoveryID)
	if err != nil {
		return nil, err
	}

	users, _ := r.GetUsers(ctx, discoveryID)
	mvpFeatures, _ := r.GetMVPFeatures(ctx, discoveryID)
	futureFeatures, _ := r.GetFutureFeatures(ctx, discoveryID)

	projectName := ""
	if discovery.ProjectName != nil {
		projectName = *discovery.ProjectName
	}

	solvesStatement := ""
	if discovery.SolvesStatement != nil {
		solvesStatement = *discovery.SolvesStatement
	}

	return &model.DiscoverySummary{
		ProjectName:     projectName,
		SolvesStatement: solvesStatement,
		Users:           users,
		MVPFeatures:     mvpFeatures,
		FutureFeatures:  futureFeatures,
	}, nil
}
