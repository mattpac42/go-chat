package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// MockFileSourceRepository implements FileSourceRepository for testing.
type MockFileSourceRepository struct {
	mu       sync.RWMutex
	sources  map[uuid.UUID]*model.FileSource // keyed by source ID
	byFileID map[uuid.UUID]uuid.UUID         // fileID -> sourceID lookup
}

// NewMockFileSourceRepository creates a new MockFileSourceRepository.
func NewMockFileSourceRepository() *MockFileSourceRepository {
	return &MockFileSourceRepository{
		sources:  make(map[uuid.UUID]*model.FileSource),
		byFileID: make(map[uuid.UUID]uuid.UUID),
	}
}

// Create creates a new file source record.
func (r *MockFileSourceRepository) Create(ctx context.Context, fileID uuid.UUID, originalFilename, originalMimeType string, originalSizeBytes int64) (*model.FileSource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	source := &model.FileSource{
		ID:                uuid.New(),
		FileID:            fileID,
		OriginalFilename:  originalFilename,
		OriginalMimeType:  originalMimeType,
		OriginalSizeBytes: originalSizeBytes,
		ConversionStatus:  "completed",
		CreatedAt:         now,
	}

	r.sources[source.ID] = source
	r.byFileID[fileID] = source.ID

	return source, nil
}

// GetByFileID retrieves a file source by the associated file ID.
func (r *MockFileSourceRepository) GetByFileID(ctx context.Context, fileID uuid.UUID) (*model.FileSource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sourceID, exists := r.byFileID[fileID]
	if !exists {
		return nil, ErrNotFound
	}

	source, ok := r.sources[sourceID]
	if !ok {
		return nil, ErrNotFound
	}

	return source, nil
}

// GetByID retrieves a file source by its ID.
func (r *MockFileSourceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.FileSource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	source, ok := r.sources[id]
	if !ok {
		return nil, ErrNotFound
	}

	return source, nil
}

// UpdateStatus updates the conversion status of a file source.
func (r *MockFileSourceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	source, ok := r.sources[id]
	if !ok {
		return ErrNotFound
	}

	source.ConversionStatus = status
	return nil
}

// Delete removes a file source record.
func (r *MockFileSourceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	source, ok := r.sources[id]
	if !ok {
		return ErrNotFound
	}

	delete(r.byFileID, source.FileID)
	delete(r.sources, id)

	return nil
}

// GetAll returns all file sources (for testing).
func (r *MockFileSourceRepository) GetAll() []*model.FileSource {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.FileSource
	for _, source := range r.sources {
		result = append(result, source)
	}
	return result
}
