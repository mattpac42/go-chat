package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// PRDRepository defines the interface for PRD data access.
type PRDRepository interface {
	// CRUD operations
	Create(ctx context.Context, prd *model.PRD) (*model.PRD, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.PRD, error)
	Update(ctx context.Context, prd *model.PRD) (*model.PRD, error)
	Delete(ctx context.Context, id uuid.UUID) error

	// Query methods
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.PRD, error)
	GetByDiscoveryID(ctx context.Context, discoveryID uuid.UUID) ([]model.PRD, error)
	GetByFeatureID(ctx context.Context, featureID uuid.UUID) (*model.PRD, error)
	GetByStatus(ctx context.Context, projectID uuid.UUID, status model.PRDStatus) ([]model.PRD, error)
	GetByVersion(ctx context.Context, projectID uuid.UUID, version string) ([]model.PRD, error)

	// Status operations
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.PRDStatus) error
	IncrementGenerationAttempts(ctx context.Context, id uuid.UUID) error
	SetLastError(ctx context.Context, id uuid.UUID, err string) error

	// Timestamp operations
	SetGeneratedAt(ctx context.Context, id uuid.UUID) error
	SetApprovedAt(ctx context.Context, id uuid.UUID) error
	SetStartedAt(ctx context.Context, id uuid.UUID) error
	SetCompletedAt(ctx context.Context, id uuid.UUID) error
}

// PostgresPRDRepository implements PRDRepository using PostgreSQL.
type PostgresPRDRepository struct {
	db *sqlx.DB
}

// NewPostgresPRDRepository creates a new PostgresPRDRepository.
func NewPostgresPRDRepository(db *sqlx.DB) *PostgresPRDRepository {
	return &PostgresPRDRepository{db: db}
}

// prdColumns defines the common columns for PRD queries.
const prdColumns = `id, discovery_id, feature_id, project_id, title, overview, version, priority,
	user_stories, acceptance_criteria, technical_notes, status, generated_at, approved_at,
	started_at, completed_at, generation_attempts, last_error, created_at, updated_at`

// Create creates a new PRD record.
func (r *PostgresPRDRepository) Create(ctx context.Context, prd *model.PRD) (*model.PRD, error) {
	query := `
		INSERT INTO prds (discovery_id, feature_id, project_id, title, overview, version, priority,
			user_stories, acceptance_criteria, technical_notes, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING ` + prdColumns

	// Set default values for JSON fields if nil
	userStoriesJSON := prd.UserStoriesJSON
	if userStoriesJSON == nil {
		userStoriesJSON = []byte("[]")
	}

	acceptanceCriteriaJSON := prd.AcceptanceCriteriaJSON
	if acceptanceCriteriaJSON == nil {
		acceptanceCriteriaJSON = []byte("[]")
	}

	technicalNotesJSON := prd.TechnicalNotesJSON
	if technicalNotesJSON == nil {
		technicalNotesJSON = []byte("[]")
	}

	status := prd.Status
	if status == "" {
		status = model.PRDStatusPending
	}

	var created model.PRD
	if err := r.db.GetContext(ctx, &created, query,
		prd.DiscoveryID,
		prd.FeatureID,
		prd.ProjectID,
		prd.Title,
		prd.Overview,
		prd.Version,
		prd.Priority,
		userStoriesJSON,
		acceptanceCriteriaJSON,
		technicalNotesJSON,
		status,
	); err != nil {
		return nil, err
	}

	return &created, nil
}

// GetByID retrieves a PRD by its ID.
func (r *PostgresPRDRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PRD, error) {
	query := `SELECT ` + prdColumns + ` FROM prds WHERE id = $1`

	var prd model.PRD
	if err := r.db.GetContext(ctx, &prd, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &prd, nil
}

// Update updates a PRD record.
func (r *PostgresPRDRepository) Update(ctx context.Context, prd *model.PRD) (*model.PRD, error) {
	query := `
		UPDATE prds
		SET title = $2,
			overview = $3,
			version = $4,
			priority = $5,
			user_stories = $6,
			acceptance_criteria = $7,
			technical_notes = $8,
			status = $9,
			updated_at = NOW()
		WHERE id = $1
		RETURNING ` + prdColumns

	// Set default values for JSON fields if nil
	userStoriesJSON := prd.UserStoriesJSON
	if userStoriesJSON == nil {
		userStoriesJSON = []byte("[]")
	}

	acceptanceCriteriaJSON := prd.AcceptanceCriteriaJSON
	if acceptanceCriteriaJSON == nil {
		acceptanceCriteriaJSON = []byte("[]")
	}

	technicalNotesJSON := prd.TechnicalNotesJSON
	if technicalNotesJSON == nil {
		technicalNotesJSON = []byte("[]")
	}

	var updated model.PRD
	if err := r.db.GetContext(ctx, &updated, query,
		prd.ID,
		prd.Title,
		prd.Overview,
		prd.Version,
		prd.Priority,
		userStoriesJSON,
		acceptanceCriteriaJSON,
		technicalNotesJSON,
		prd.Status,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &updated, nil
}

// Delete removes a PRD record.
func (r *PostgresPRDRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM prds WHERE id = $1`

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

// GetByProjectID retrieves all PRDs for a project.
func (r *PostgresPRDRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.PRD, error) {
	query := `
		SELECT ` + prdColumns + `
		FROM prds
		WHERE project_id = $1
		ORDER BY version ASC, priority ASC, created_at ASC
	`

	var prds []model.PRD
	if err := r.db.SelectContext(ctx, &prds, query, projectID); err != nil {
		return nil, err
	}

	if prds == nil {
		prds = []model.PRD{}
	}

	return prds, nil
}

// GetByDiscoveryID retrieves all PRDs for a discovery.
func (r *PostgresPRDRepository) GetByDiscoveryID(ctx context.Context, discoveryID uuid.UUID) ([]model.PRD, error) {
	query := `
		SELECT ` + prdColumns + `
		FROM prds
		WHERE discovery_id = $1
		ORDER BY version ASC, priority ASC, created_at ASC
	`

	var prds []model.PRD
	if err := r.db.SelectContext(ctx, &prds, query, discoveryID); err != nil {
		return nil, err
	}

	if prds == nil {
		prds = []model.PRD{}
	}

	return prds, nil
}

// GetByFeatureID retrieves the PRD for a specific feature.
func (r *PostgresPRDRepository) GetByFeatureID(ctx context.Context, featureID uuid.UUID) (*model.PRD, error) {
	query := `SELECT ` + prdColumns + ` FROM prds WHERE feature_id = $1`

	var prd model.PRD
	if err := r.db.GetContext(ctx, &prd, query, featureID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &prd, nil
}

// GetByStatus retrieves all PRDs for a project with a specific status.
func (r *PostgresPRDRepository) GetByStatus(ctx context.Context, projectID uuid.UUID, status model.PRDStatus) ([]model.PRD, error) {
	query := `
		SELECT ` + prdColumns + `
		FROM prds
		WHERE project_id = $1 AND status = $2
		ORDER BY version ASC, priority ASC, created_at ASC
	`

	var prds []model.PRD
	if err := r.db.SelectContext(ctx, &prds, query, projectID, status); err != nil {
		return nil, err
	}

	if prds == nil {
		prds = []model.PRD{}
	}

	return prds, nil
}

// GetByVersion retrieves all PRDs for a project with a specific version.
func (r *PostgresPRDRepository) GetByVersion(ctx context.Context, projectID uuid.UUID, version string) ([]model.PRD, error) {
	query := `
		SELECT ` + prdColumns + `
		FROM prds
		WHERE project_id = $1 AND version = $2
		ORDER BY priority ASC, created_at ASC
	`

	var prds []model.PRD
	if err := r.db.SelectContext(ctx, &prds, query, projectID, version); err != nil {
		return nil, err
	}

	if prds == nil {
		prds = []model.PRD{}
	}

	return prds, nil
}

// UpdateStatus updates the status of a PRD.
func (r *PostgresPRDRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.PRDStatus) error {
	query := `UPDATE prds SET status = $2, updated_at = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, status)
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

// IncrementGenerationAttempts increments the generation attempts counter.
func (r *PostgresPRDRepository) IncrementGenerationAttempts(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE prds SET generation_attempts = generation_attempts + 1, updated_at = NOW() WHERE id = $1`

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

// SetLastError sets the last error message for a PRD.
func (r *PostgresPRDRepository) SetLastError(ctx context.Context, id uuid.UUID, errMsg string) error {
	query := `UPDATE prds SET last_error = $2, updated_at = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, errMsg)
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

// SetGeneratedAt sets the generated_at timestamp to now.
func (r *PostgresPRDRepository) SetGeneratedAt(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE prds SET generated_at = NOW(), updated_at = NOW() WHERE id = $1`

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

// SetApprovedAt sets the approved_at timestamp to now.
func (r *PostgresPRDRepository) SetApprovedAt(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE prds SET approved_at = NOW(), updated_at = NOW() WHERE id = $1`

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

// SetStartedAt sets the started_at timestamp to now.
func (r *PostgresPRDRepository) SetStartedAt(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE prds SET started_at = NOW(), updated_at = NOW() WHERE id = $1`

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

// SetCompletedAt sets the completed_at timestamp to now.
func (r *PostgresPRDRepository) SetCompletedAt(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE prds SET completed_at = NOW(), updated_at = NOW() WHERE id = $1`

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

// MockPRDRepository implements PRDRepository for testing.
type MockPRDRepository struct {
	prds        map[uuid.UUID]*model.PRD
	byProject   map[uuid.UUID][]uuid.UUID   // projectID -> []prdID
	byDiscovery map[uuid.UUID][]uuid.UUID   // discoveryID -> []prdID
	byFeature   map[uuid.UUID]uuid.UUID     // featureID -> prdID
}

// NewMockPRDRepository creates a new MockPRDRepository.
func NewMockPRDRepository() *MockPRDRepository {
	return &MockPRDRepository{
		prds:        make(map[uuid.UUID]*model.PRD),
		byProject:   make(map[uuid.UUID][]uuid.UUID),
		byDiscovery: make(map[uuid.UUID][]uuid.UUID),
		byFeature:   make(map[uuid.UUID]uuid.UUID),
	}
}

// Create creates a new PRD record.
func (r *MockPRDRepository) Create(ctx context.Context, prd *model.PRD) (*model.PRD, error) {
	// Check if feature already has a PRD
	if _, ok := r.byFeature[prd.FeatureID]; ok {
		return nil, errors.New("PRD already exists for feature")
	}

	now := time.Now().UTC()
	created := &model.PRD{
		ID:                     uuid.New(),
		DiscoveryID:            prd.DiscoveryID,
		FeatureID:              prd.FeatureID,
		ProjectID:              prd.ProjectID,
		Title:                  prd.Title,
		Overview:               prd.Overview,
		Version:                prd.Version,
		Priority:               prd.Priority,
		UserStoriesJSON:        prd.UserStoriesJSON,
		AcceptanceCriteriaJSON: prd.AcceptanceCriteriaJSON,
		TechnicalNotesJSON:     prd.TechnicalNotesJSON,
		Status:                 prd.Status,
		GenerationAttempts:     0,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	if created.Status == "" {
		created.Status = model.PRDStatusPending
	}
	if created.UserStoriesJSON == nil {
		created.UserStoriesJSON = []byte("[]")
	}
	if created.AcceptanceCriteriaJSON == nil {
		created.AcceptanceCriteriaJSON = []byte("[]")
	}
	if created.TechnicalNotesJSON == nil {
		created.TechnicalNotesJSON = []byte("[]")
	}

	r.prds[created.ID] = created
	r.byProject[created.ProjectID] = append(r.byProject[created.ProjectID], created.ID)
	r.byDiscovery[created.DiscoveryID] = append(r.byDiscovery[created.DiscoveryID], created.ID)
	r.byFeature[created.FeatureID] = created.ID

	// Return a copy
	copy := *created
	return &copy, nil
}

// GetByID retrieves a PRD by its ID.
func (r *MockPRDRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PRD, error) {
	prd, ok := r.prds[id]
	if !ok {
		return nil, ErrNotFound
	}
	// Return a copy
	copy := *prd
	return &copy, nil
}

// Update updates a PRD record.
func (r *MockPRDRepository) Update(ctx context.Context, prd *model.PRD) (*model.PRD, error) {
	existing, ok := r.prds[prd.ID]
	if !ok {
		return nil, ErrNotFound
	}

	existing.Title = prd.Title
	existing.Overview = prd.Overview
	existing.Version = prd.Version
	existing.Priority = prd.Priority
	existing.UserStoriesJSON = prd.UserStoriesJSON
	existing.AcceptanceCriteriaJSON = prd.AcceptanceCriteriaJSON
	existing.TechnicalNotesJSON = prd.TechnicalNotesJSON
	existing.Status = prd.Status
	existing.UpdatedAt = time.Now().UTC()

	// Return a copy
	copy := *existing
	return &copy, nil
}

// Delete removes a PRD record.
func (r *MockPRDRepository) Delete(ctx context.Context, id uuid.UUID) error {
	prd, ok := r.prds[id]
	if !ok {
		return ErrNotFound
	}

	// Remove from indexes
	delete(r.prds, id)
	delete(r.byFeature, prd.FeatureID)

	// Remove from project index
	projectPRDs := r.byProject[prd.ProjectID]
	for i, prdID := range projectPRDs {
		if prdID == id {
			r.byProject[prd.ProjectID] = append(projectPRDs[:i], projectPRDs[i+1:]...)
			break
		}
	}

	// Remove from discovery index
	discoveryPRDs := r.byDiscovery[prd.DiscoveryID]
	for i, prdID := range discoveryPRDs {
		if prdID == id {
			r.byDiscovery[prd.DiscoveryID] = append(discoveryPRDs[:i], discoveryPRDs[i+1:]...)
			break
		}
	}

	return nil
}

// GetByProjectID retrieves all PRDs for a project.
func (r *MockPRDRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.PRD, error) {
	prdIDs := r.byProject[projectID]
	prds := make([]model.PRD, 0, len(prdIDs))

	for _, prdID := range prdIDs {
		if prd, ok := r.prds[prdID]; ok {
			prds = append(prds, *prd)
		}
	}

	return prds, nil
}

// GetByDiscoveryID retrieves all PRDs for a discovery.
func (r *MockPRDRepository) GetByDiscoveryID(ctx context.Context, discoveryID uuid.UUID) ([]model.PRD, error) {
	prdIDs := r.byDiscovery[discoveryID]
	prds := make([]model.PRD, 0, len(prdIDs))

	for _, prdID := range prdIDs {
		if prd, ok := r.prds[prdID]; ok {
			prds = append(prds, *prd)
		}
	}

	return prds, nil
}

// GetByFeatureID retrieves the PRD for a specific feature.
func (r *MockPRDRepository) GetByFeatureID(ctx context.Context, featureID uuid.UUID) (*model.PRD, error) {
	prdID, ok := r.byFeature[featureID]
	if !ok {
		return nil, ErrNotFound
	}

	prd, ok := r.prds[prdID]
	if !ok {
		return nil, ErrNotFound
	}

	// Return a copy
	copy := *prd
	return &copy, nil
}

// GetByStatus retrieves all PRDs for a project with a specific status.
func (r *MockPRDRepository) GetByStatus(ctx context.Context, projectID uuid.UUID, status model.PRDStatus) ([]model.PRD, error) {
	prdIDs := r.byProject[projectID]
	prds := make([]model.PRD, 0)

	for _, prdID := range prdIDs {
		if prd, ok := r.prds[prdID]; ok {
			if prd.Status == status {
				prds = append(prds, *prd)
			}
		}
	}

	return prds, nil
}

// GetByVersion retrieves all PRDs for a project with a specific version.
func (r *MockPRDRepository) GetByVersion(ctx context.Context, projectID uuid.UUID, version string) ([]model.PRD, error) {
	prdIDs := r.byProject[projectID]
	prds := make([]model.PRD, 0)

	for _, prdID := range prdIDs {
		if prd, ok := r.prds[prdID]; ok {
			if prd.Version == version {
				prds = append(prds, *prd)
			}
		}
	}

	return prds, nil
}

// UpdateStatus updates the status of a PRD.
func (r *MockPRDRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.PRDStatus) error {
	prd, ok := r.prds[id]
	if !ok {
		return ErrNotFound
	}

	prd.Status = status
	prd.UpdatedAt = time.Now().UTC()
	return nil
}

// IncrementGenerationAttempts increments the generation attempts counter.
func (r *MockPRDRepository) IncrementGenerationAttempts(ctx context.Context, id uuid.UUID) error {
	prd, ok := r.prds[id]
	if !ok {
		return ErrNotFound
	}

	prd.GenerationAttempts++
	prd.UpdatedAt = time.Now().UTC()
	return nil
}

// SetLastError sets the last error message for a PRD.
func (r *MockPRDRepository) SetLastError(ctx context.Context, id uuid.UUID, errMsg string) error {
	prd, ok := r.prds[id]
	if !ok {
		return ErrNotFound
	}

	prd.LastError = &errMsg
	prd.UpdatedAt = time.Now().UTC()
	return nil
}

// SetGeneratedAt sets the generated_at timestamp to now.
func (r *MockPRDRepository) SetGeneratedAt(ctx context.Context, id uuid.UUID) error {
	prd, ok := r.prds[id]
	if !ok {
		return ErrNotFound
	}

	now := time.Now().UTC()
	prd.GeneratedAt = &now
	prd.UpdatedAt = now
	return nil
}

// SetApprovedAt sets the approved_at timestamp to now.
func (r *MockPRDRepository) SetApprovedAt(ctx context.Context, id uuid.UUID) error {
	prd, ok := r.prds[id]
	if !ok {
		return ErrNotFound
	}

	now := time.Now().UTC()
	prd.ApprovedAt = &now
	prd.UpdatedAt = now
	return nil
}

// SetStartedAt sets the started_at timestamp to now.
func (r *MockPRDRepository) SetStartedAt(ctx context.Context, id uuid.UUID) error {
	prd, ok := r.prds[id]
	if !ok {
		return ErrNotFound
	}

	now := time.Now().UTC()
	prd.StartedAt = &now
	prd.UpdatedAt = now
	return nil
}

// SetCompletedAt sets the completed_at timestamp to now.
func (r *MockPRDRepository) SetCompletedAt(ctx context.Context, id uuid.UUID) error {
	prd, ok := r.prds[id]
	if !ok {
		return ErrNotFound
	}

	now := time.Now().UTC()
	prd.CompletedAt = &now
	prd.UpdatedAt = now
	return nil
}
