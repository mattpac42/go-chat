package repository

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// FileRepository defines the interface for file data access.
type FileRepository interface {
	SaveFile(ctx context.Context, projectID uuid.UUID, path, language, content string) (*model.File, error)
	GetFilesByProject(ctx context.Context, projectID uuid.UUID) ([]model.FileListItem, error)
	GetFile(ctx context.Context, id uuid.UUID) (*model.File, error)
	GetFileByPath(ctx context.Context, projectID uuid.UUID, path string) (*model.File, error)
}

// PostgresFileRepository implements FileRepository using PostgreSQL.
type PostgresFileRepository struct {
	db *sqlx.DB
}

// NewPostgresFileRepository creates a new PostgresFileRepository.
func NewPostgresFileRepository(db *sqlx.DB) *PostgresFileRepository {
	return &PostgresFileRepository{db: db}
}

// SaveFile saves or updates a file for a project (upsert by project_id + path).
func (r *PostgresFileRepository) SaveFile(ctx context.Context, projectID uuid.UUID, path, language, content string) (*model.File, error) {
	filename := filepath.Base(path)

	query := `
		INSERT INTO files (project_id, path, filename, language, content)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (project_id, path)
		DO UPDATE SET
			language = EXCLUDED.language,
			content = EXCLUDED.content,
			created_at = NOW()
		RETURNING id, project_id, path, filename, language, content, created_at
	`

	var file model.File
	if err := r.db.GetContext(ctx, &file, query, projectID, path, filename, language, content); err != nil {
		return nil, err
	}

	return &file, nil
}

// GetFilesByProject returns all files for a project (without content).
func (r *PostgresFileRepository) GetFilesByProject(ctx context.Context, projectID uuid.UUID) ([]model.FileListItem, error) {
	query := `
		SELECT id, path, filename, language, created_at
		FROM files
		WHERE project_id = $1
		ORDER BY path ASC
	`

	var files []model.FileListItem
	if err := r.db.SelectContext(ctx, &files, query, projectID); err != nil {
		return nil, err
	}

	return files, nil
}

// GetFile returns a file by ID.
func (r *PostgresFileRepository) GetFile(ctx context.Context, id uuid.UUID) (*model.File, error) {
	query := `
		SELECT id, project_id, path, filename, language, content, created_at
		FROM files
		WHERE id = $1
	`

	var file model.File
	if err := r.db.GetContext(ctx, &file, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &file, nil
}

// GetFileByPath returns a file by project ID and path.
func (r *PostgresFileRepository) GetFileByPath(ctx context.Context, projectID uuid.UUID, path string) (*model.File, error) {
	query := `
		SELECT id, project_id, path, filename, language, content, created_at
		FROM files
		WHERE project_id = $1 AND path = $2
	`

	var file model.File
	if err := r.db.GetContext(ctx, &file, query, projectID, path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &file, nil
}
