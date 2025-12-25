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

// ErrNotFound is returned when a resource is not found.
var ErrNotFound = errors.New("not found")

// ProjectRepository defines the interface for project data access.
type ProjectRepository interface {
	List(ctx context.Context) ([]model.ProjectListItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error)
	Create(ctx context.Context, title string) (*model.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateTimestamp(ctx context.Context, id uuid.UUID, timestamp time.Time) error
	UpdateTitle(ctx context.Context, id uuid.UUID, title string) (*model.Project, error)
	GetMessages(ctx context.Context, projectID uuid.UUID) ([]model.Message, error)
	CreateMessage(ctx context.Context, projectID uuid.UUID, role model.Role, content string) (*model.Message, error)
}

// PostgresProjectRepository implements ProjectRepository using PostgreSQL.
type PostgresProjectRepository struct {
	db *sqlx.DB
}

// NewPostgresProjectRepository creates a new PostgresProjectRepository.
func NewPostgresProjectRepository(db *sqlx.DB) *PostgresProjectRepository {
	return &PostgresProjectRepository{db: db}
}

// List returns all projects with message counts.
func (r *PostgresProjectRepository) List(ctx context.Context) ([]model.ProjectListItem, error) {
	query := `
		SELECT
			p.id,
			p.title,
			p.created_at,
			p.updated_at,
			COUNT(m.id) as message_count
		FROM projects p
		LEFT JOIN messages m ON p.id = m.project_id
		GROUP BY p.id
		ORDER BY p.updated_at DESC
	`

	var projects []model.ProjectListItem
	if err := r.db.SelectContext(ctx, &projects, query); err != nil {
		return nil, err
	}

	return projects, nil
}

// GetByID returns a project by ID.
func (r *PostgresProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	query := `SELECT id, title, created_at, updated_at FROM projects WHERE id = $1`

	var project model.Project
	if err := r.db.GetContext(ctx, &project, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &project, nil
}

// Create creates a new project.
func (r *PostgresProjectRepository) Create(ctx context.Context, title string) (*model.Project, error) {
	query := `
		INSERT INTO projects (title)
		VALUES ($1)
		RETURNING id, title, created_at, updated_at
	`

	var project model.Project
	if err := r.db.GetContext(ctx, &project, query, title); err != nil {
		return nil, err
	}

	return &project, nil
}

// Delete removes a project by ID.
func (r *PostgresProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`

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

// UpdateTimestamp updates a project's updated_at timestamp.
func (r *PostgresProjectRepository) UpdateTimestamp(ctx context.Context, id uuid.UUID, timestamp time.Time) error {
	query := `UPDATE projects SET updated_at = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, timestamp, id)
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

// UpdateTitle updates a project's title and updated_at timestamp.
func (r *PostgresProjectRepository) UpdateTitle(ctx context.Context, id uuid.UUID, title string) (*model.Project, error) {
	query := `
		UPDATE projects
		SET title = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, title, created_at, updated_at
	`

	var project model.Project
	if err := r.db.GetContext(ctx, &project, query, title, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &project, nil
}

// GetMessages returns all messages for a project.
func (r *PostgresProjectRepository) GetMessages(ctx context.Context, projectID uuid.UUID) ([]model.Message, error) {
	query := `
		SELECT id, project_id, role, content, created_at
		FROM messages
		WHERE project_id = $1
		ORDER BY created_at ASC
	`

	var messages []model.Message
	if err := r.db.SelectContext(ctx, &messages, query, projectID); err != nil {
		return nil, err
	}

	return messages, nil
}

// CreateMessage creates a new message.
func (r *PostgresProjectRepository) CreateMessage(ctx context.Context, projectID uuid.UUID, role model.Role, content string) (*model.Message, error) {
	query := `
		INSERT INTO messages (project_id, role, content)
		VALUES ($1, $2, $3)
		RETURNING id, project_id, role, content, created_at
	`

	var message model.Message
	if err := r.db.GetContext(ctx, &message, query, projectID, role, content); err != nil {
		return nil, err
	}

	return &message, nil
}
