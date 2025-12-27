package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// FileSourceRepository defines the interface for file source data access.
type FileSourceRepository interface {
	// Create creates a new file source record.
	Create(ctx context.Context, fileID uuid.UUID, originalFilename, originalMimeType string, originalSizeBytes int64) (*model.FileSource, error)

	// GetByFileID retrieves a file source by the associated file ID.
	GetByFileID(ctx context.Context, fileID uuid.UUID) (*model.FileSource, error)

	// GetByID retrieves a file source by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*model.FileSource, error)

	// UpdateStatus updates the conversion status of a file source.
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error

	// Delete removes a file source record.
	Delete(ctx context.Context, id uuid.UUID) error
}

// PostgresFileSourceRepository implements FileSourceRepository using PostgreSQL.
type PostgresFileSourceRepository struct {
	db *sqlx.DB
}

// NewPostgresFileSourceRepository creates a new PostgresFileSourceRepository.
func NewPostgresFileSourceRepository(db *sqlx.DB) *PostgresFileSourceRepository {
	return &PostgresFileSourceRepository{db: db}
}

// Create creates a new file source record.
func (r *PostgresFileSourceRepository) Create(ctx context.Context, fileID uuid.UUID, originalFilename, originalMimeType string, originalSizeBytes int64) (*model.FileSource, error) {
	query := `
		INSERT INTO file_sources (file_id, original_filename, original_mime_type, original_size_bytes, conversion_status)
		VALUES ($1, $2, $3, $4, 'completed')
		RETURNING id, file_id, original_filename, original_mime_type, original_size_bytes, conversion_status, created_at
	`

	var source model.FileSource
	if err := r.db.GetContext(ctx, &source, query, fileID, originalFilename, originalMimeType, originalSizeBytes); err != nil {
		return nil, err
	}

	return &source, nil
}

// GetByFileID retrieves a file source by the associated file ID.
func (r *PostgresFileSourceRepository) GetByFileID(ctx context.Context, fileID uuid.UUID) (*model.FileSource, error) {
	query := `
		SELECT id, file_id, original_filename, original_mime_type, original_size_bytes, conversion_status, created_at
		FROM file_sources
		WHERE file_id = $1
	`

	var source model.FileSource
	if err := r.db.GetContext(ctx, &source, query, fileID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &source, nil
}

// GetByID retrieves a file source by its ID.
func (r *PostgresFileSourceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.FileSource, error) {
	query := `
		SELECT id, file_id, original_filename, original_mime_type, original_size_bytes, conversion_status, created_at
		FROM file_sources
		WHERE id = $1
	`

	var source model.FileSource
	if err := r.db.GetContext(ctx, &source, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &source, nil
}

// UpdateStatus updates the conversion status of a file source.
func (r *PostgresFileSourceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE file_sources SET conversion_status = $2 WHERE id = $1`

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

// Delete removes a file source record.
func (r *PostgresFileSourceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM file_sources WHERE id = $1`

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
