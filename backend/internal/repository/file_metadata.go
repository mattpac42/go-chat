package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// FileMetadataRepository defines the interface for file metadata data access.
type FileMetadataRepository interface {
	// Create creates new metadata for a file.
	Create(ctx context.Context, fileID uuid.UUID, shortDesc, longDesc, funcGroup string) (*model.FileMetadata, error)

	// GetByFileID retrieves metadata for a specific file.
	GetByFileID(ctx context.Context, fileID uuid.UUID) (*model.FileMetadata, error)

	// GetByID retrieves metadata by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*model.FileMetadata, error)

	// Update updates existing file metadata.
	Update(ctx context.Context, id uuid.UUID, shortDesc, longDesc, funcGroup *string) (*model.FileMetadata, error)

	// Upsert creates or updates metadata for a file.
	Upsert(ctx context.Context, fileID uuid.UUID, shortDesc, longDesc, funcGroup string) (*model.FileMetadata, error)

	// Delete removes metadata for a file.
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByFileID removes metadata by file ID.
	DeleteByFileID(ctx context.Context, fileID uuid.UUID) error

	// GetFilesWithMetadata retrieves all files with their metadata for a project.
	GetFilesWithMetadata(ctx context.Context, projectID uuid.UUID) ([]model.FileWithMetadata, error)

	// GetFilesByFunctionalGroup retrieves files in a specific functional group.
	GetFilesByFunctionalGroup(ctx context.Context, projectID uuid.UUID, funcGroup string) ([]model.FileWithMetadata, error)

	// GetFunctionalGroups retrieves all functional groups for a project with file counts.
	GetFunctionalGroups(ctx context.Context, projectID uuid.UUID) ([]model.FunctionalGroupSummary, error)
}

// PostgresFileMetadataRepository implements FileMetadataRepository using PostgreSQL.
type PostgresFileMetadataRepository struct {
	db *sqlx.DB
}

// NewPostgresFileMetadataRepository creates a new PostgresFileMetadataRepository.
func NewPostgresFileMetadataRepository(db *sqlx.DB) *PostgresFileMetadataRepository {
	return &PostgresFileMetadataRepository{db: db}
}

// Create creates new metadata for a file.
func (r *PostgresFileMetadataRepository) Create(ctx context.Context, fileID uuid.UUID, shortDesc, longDesc, funcGroup string) (*model.FileMetadata, error) {
	query := `
		INSERT INTO file_metadata (file_id, short_description, long_description, functional_group)
		VALUES ($1, $2, $3, $4)
		RETURNING id, file_id, short_description, long_description, functional_group, created_at, updated_at
	`

	var metadata model.FileMetadata
	if err := r.db.GetContext(ctx, &metadata, query, fileID, shortDesc, longDesc, funcGroup); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// GetByFileID retrieves metadata for a specific file.
func (r *PostgresFileMetadataRepository) GetByFileID(ctx context.Context, fileID uuid.UUID) (*model.FileMetadata, error) {
	query := `
		SELECT id, file_id, short_description, long_description, functional_group, created_at, updated_at
		FROM file_metadata
		WHERE file_id = $1
	`

	var metadata model.FileMetadata
	if err := r.db.GetContext(ctx, &metadata, query, fileID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &metadata, nil
}

// GetByID retrieves metadata by its ID.
func (r *PostgresFileMetadataRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.FileMetadata, error) {
	query := `
		SELECT id, file_id, short_description, long_description, functional_group, created_at, updated_at
		FROM file_metadata
		WHERE id = $1
	`

	var metadata model.FileMetadata
	if err := r.db.GetContext(ctx, &metadata, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &metadata, nil
}

// Update updates existing file metadata.
func (r *PostgresFileMetadataRepository) Update(ctx context.Context, id uuid.UUID, shortDesc, longDesc, funcGroup *string) (*model.FileMetadata, error) {
	query := `
		UPDATE file_metadata
		SET
			short_description = COALESCE($2, short_description),
			long_description = COALESCE($3, long_description),
			functional_group = COALESCE($4, functional_group),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, file_id, short_description, long_description, functional_group, created_at, updated_at
	`

	var metadata model.FileMetadata
	if err := r.db.GetContext(ctx, &metadata, query, id, shortDesc, longDesc, funcGroup); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &metadata, nil
}

// Upsert creates or updates metadata for a file.
func (r *PostgresFileMetadataRepository) Upsert(ctx context.Context, fileID uuid.UUID, shortDesc, longDesc, funcGroup string) (*model.FileMetadata, error) {
	query := `
		INSERT INTO file_metadata (file_id, short_description, long_description, functional_group)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (file_id)
		DO UPDATE SET
			short_description = EXCLUDED.short_description,
			long_description = EXCLUDED.long_description,
			functional_group = EXCLUDED.functional_group,
			updated_at = NOW()
		RETURNING id, file_id, short_description, long_description, functional_group, created_at, updated_at
	`

	var metadata model.FileMetadata
	if err := r.db.GetContext(ctx, &metadata, query, fileID, shortDesc, longDesc, funcGroup); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// Delete removes metadata by ID.
func (r *PostgresFileMetadataRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM file_metadata WHERE id = $1`

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

// DeleteByFileID removes metadata by file ID.
func (r *PostgresFileMetadataRepository) DeleteByFileID(ctx context.Context, fileID uuid.UUID) error {
	query := `DELETE FROM file_metadata WHERE file_id = $1`

	result, err := r.db.ExecContext(ctx, query, fileID)
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

// GetFilesWithMetadata retrieves all files with their metadata for a project.
func (r *PostgresFileMetadataRepository) GetFilesWithMetadata(ctx context.Context, projectID uuid.UUID) ([]model.FileWithMetadata, error) {
	query := `
		SELECT
			f.id,
			f.project_id,
			f.path,
			f.filename,
			f.language,
			COALESCE(fm.short_description, '') as short_description,
			COALESCE(fm.long_description, '') as long_description,
			COALESCE(fm.functional_group, '') as functional_group,
			f.created_at
		FROM files f
		LEFT JOIN file_metadata fm ON f.id = fm.file_id
		WHERE f.project_id = $1
		ORDER BY fm.functional_group NULLS LAST, f.path ASC
	`

	var files []model.FileWithMetadata
	if err := r.db.SelectContext(ctx, &files, query, projectID); err != nil {
		return nil, err
	}

	return files, nil
}

// GetFilesByFunctionalGroup retrieves files in a specific functional group.
func (r *PostgresFileMetadataRepository) GetFilesByFunctionalGroup(ctx context.Context, projectID uuid.UUID, funcGroup string) ([]model.FileWithMetadata, error) {
	query := `
		SELECT
			f.id,
			f.project_id,
			f.path,
			f.filename,
			f.language,
			COALESCE(fm.short_description, '') as short_description,
			COALESCE(fm.long_description, '') as long_description,
			COALESCE(fm.functional_group, '') as functional_group,
			f.created_at
		FROM files f
		INNER JOIN file_metadata fm ON f.id = fm.file_id
		WHERE f.project_id = $1 AND fm.functional_group = $2
		ORDER BY f.path ASC
	`

	var files []model.FileWithMetadata
	if err := r.db.SelectContext(ctx, &files, query, projectID, funcGroup); err != nil {
		return nil, err
	}

	return files, nil
}

// GetFunctionalGroups retrieves all functional groups for a project with file counts.
func (r *PostgresFileMetadataRepository) GetFunctionalGroups(ctx context.Context, projectID uuid.UUID) ([]model.FunctionalGroupSummary, error) {
	query := `
		SELECT
			fm.functional_group as name,
			COUNT(*) as file_count
		FROM files f
		INNER JOIN file_metadata fm ON f.id = fm.file_id
		WHERE f.project_id = $1 AND fm.functional_group IS NOT NULL AND fm.functional_group != ''
		GROUP BY fm.functional_group
		ORDER BY fm.functional_group ASC
	`

	var groups []model.FunctionalGroupSummary
	if err := r.db.SelectContext(ctx, &groups, query, projectID); err != nil {
		return nil, err
	}

	return groups, nil
}
