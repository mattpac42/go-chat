package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// MockFileMetadataRepository implements FileMetadataRepository for testing.
type MockFileMetadataRepository struct {
	mu       sync.RWMutex
	metadata map[uuid.UUID]*model.FileMetadata // keyed by metadata ID
	byFileID map[uuid.UUID]uuid.UUID           // fileID -> metadataID lookup
	files    map[uuid.UUID]*model.File         // file storage for join queries
}

// NewMockFileMetadataRepository creates a new MockFileMetadataRepository.
func NewMockFileMetadataRepository() *MockFileMetadataRepository {
	return &MockFileMetadataRepository{
		metadata: make(map[uuid.UUID]*model.FileMetadata),
		byFileID: make(map[uuid.UUID]uuid.UUID),
		files:    make(map[uuid.UUID]*model.File),
	}
}

// AddFile adds a file to the mock repository for testing join queries.
func (r *MockFileMetadataRepository) AddFile(file *model.File) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.files[file.ID] = file
}

// Create creates new metadata for a file.
func (r *MockFileMetadataRepository) Create(ctx context.Context, fileID uuid.UUID, shortDesc, longDesc, funcGroup string) (*model.FileMetadata, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if metadata already exists for this file
	if _, exists := r.byFileID[fileID]; exists {
		return nil, ErrNotFound // Simulates unique constraint violation
	}

	now := time.Now().UTC()
	meta := &model.FileMetadata{
		ID:               uuid.New(),
		FileID:           fileID,
		ShortDescription: shortDesc,
		LongDescription:  longDesc,
		FunctionalGroup:  funcGroup,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	r.metadata[meta.ID] = meta
	r.byFileID[fileID] = meta.ID

	return meta, nil
}

// GetByFileID retrieves metadata for a specific file.
func (r *MockFileMetadataRepository) GetByFileID(ctx context.Context, fileID uuid.UUID) (*model.FileMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metaID, exists := r.byFileID[fileID]
	if !exists {
		return nil, ErrNotFound
	}

	meta, ok := r.metadata[metaID]
	if !ok {
		return nil, ErrNotFound
	}

	return meta, nil
}

// GetByID retrieves metadata by its ID.
func (r *MockFileMetadataRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.FileMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	meta, ok := r.metadata[id]
	if !ok {
		return nil, ErrNotFound
	}

	return meta, nil
}

// Update updates existing file metadata.
func (r *MockFileMetadataRepository) Update(ctx context.Context, id uuid.UUID, shortDesc, longDesc, funcGroup *string) (*model.FileMetadata, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	meta, ok := r.metadata[id]
	if !ok {
		return nil, ErrNotFound
	}

	if shortDesc != nil {
		meta.ShortDescription = *shortDesc
	}
	if longDesc != nil {
		meta.LongDescription = *longDesc
	}
	if funcGroup != nil {
		meta.FunctionalGroup = *funcGroup
	}
	meta.UpdatedAt = time.Now().UTC()

	return meta, nil
}

// Upsert creates or updates metadata for a file.
func (r *MockFileMetadataRepository) Upsert(ctx context.Context, fileID uuid.UUID, shortDesc, longDesc, funcGroup string) (*model.FileMetadata, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()

	// Check if metadata exists for this file
	if metaID, exists := r.byFileID[fileID]; exists {
		meta := r.metadata[metaID]
		meta.ShortDescription = shortDesc
		meta.LongDescription = longDesc
		meta.FunctionalGroup = funcGroup
		meta.UpdatedAt = now
		return meta, nil
	}

	// Create new metadata
	meta := &model.FileMetadata{
		ID:               uuid.New(),
		FileID:           fileID,
		ShortDescription: shortDesc,
		LongDescription:  longDesc,
		FunctionalGroup:  funcGroup,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	r.metadata[meta.ID] = meta
	r.byFileID[fileID] = meta.ID

	return meta, nil
}

// Delete removes metadata by ID.
func (r *MockFileMetadataRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	meta, ok := r.metadata[id]
	if !ok {
		return ErrNotFound
	}

	delete(r.byFileID, meta.FileID)
	delete(r.metadata, id)

	return nil
}

// DeleteByFileID removes metadata by file ID.
func (r *MockFileMetadataRepository) DeleteByFileID(ctx context.Context, fileID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	metaID, exists := r.byFileID[fileID]
	if !exists {
		return ErrNotFound
	}

	delete(r.metadata, metaID)
	delete(r.byFileID, fileID)

	return nil
}

// GetFilesWithMetadata retrieves all files with their metadata for a project.
func (r *MockFileMetadataRepository) GetFilesWithMetadata(ctx context.Context, projectID uuid.UUID) ([]model.FileWithMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.FileWithMetadata

	for _, file := range r.files {
		if file.ProjectID != projectID {
			continue
		}

		fwm := model.FileWithMetadata{
			ID:        file.ID,
			ProjectID: file.ProjectID,
			Path:      file.Path,
			Filename:  file.Filename,
			Language:  file.Language,
			CreatedAt: file.CreatedAt,
		}

		// Join with metadata if it exists
		if metaID, exists := r.byFileID[file.ID]; exists {
			if meta, ok := r.metadata[metaID]; ok {
				fwm.ShortDescription = meta.ShortDescription
				fwm.LongDescription = meta.LongDescription
				fwm.FunctionalGroup = meta.FunctionalGroup
			}
		}

		result = append(result, fwm)
	}

	return result, nil
}

// GetFilesByFunctionalGroup retrieves files in a specific functional group.
func (r *MockFileMetadataRepository) GetFilesByFunctionalGroup(ctx context.Context, projectID uuid.UUID, funcGroup string) ([]model.FileWithMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.FileWithMetadata

	for _, file := range r.files {
		if file.ProjectID != projectID {
			continue
		}

		metaID, exists := r.byFileID[file.ID]
		if !exists {
			continue
		}

		meta, ok := r.metadata[metaID]
		if !ok || meta.FunctionalGroup != funcGroup {
			continue
		}

		fwm := model.FileWithMetadata{
			ID:               file.ID,
			ProjectID:        file.ProjectID,
			Path:             file.Path,
			Filename:         file.Filename,
			Language:         file.Language,
			ShortDescription: meta.ShortDescription,
			LongDescription:  meta.LongDescription,
			FunctionalGroup:  meta.FunctionalGroup,
			CreatedAt:        file.CreatedAt,
		}

		result = append(result, fwm)
	}

	return result, nil
}

// GetFunctionalGroups retrieves all functional groups for a project with file counts.
func (r *MockFileMetadataRepository) GetFunctionalGroups(ctx context.Context, projectID uuid.UUID) ([]model.FunctionalGroupSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	groupCounts := make(map[string]int)

	for _, file := range r.files {
		if file.ProjectID != projectID {
			continue
		}

		metaID, exists := r.byFileID[file.ID]
		if !exists {
			continue
		}

		meta, ok := r.metadata[metaID]
		if !ok || meta.FunctionalGroup == "" {
			continue
		}

		groupCounts[meta.FunctionalGroup]++
	}

	var result []model.FunctionalGroupSummary
	for name, count := range groupCounts {
		result = append(result, model.FunctionalGroupSummary{
			Name:      name,
			FileCount: count,
		})
	}

	return result, nil
}
