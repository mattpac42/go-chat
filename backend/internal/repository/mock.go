package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// MockProjectRepository implements ProjectRepository for testing.
type MockProjectRepository struct {
	mu       sync.RWMutex
	projects map[uuid.UUID]*model.Project
	messages map[uuid.UUID][]model.Message
}

// NewMockProjectRepository creates a new MockProjectRepository.
func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{
		projects: make(map[uuid.UUID]*model.Project),
		messages: make(map[uuid.UUID][]model.Message),
	}
}

// List returns all projects with message counts.
func (r *MockProjectRepository) List(ctx context.Context) ([]model.ProjectListItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]model.ProjectListItem, 0, len(r.projects))
	for _, p := range r.projects {
		items = append(items, model.ProjectListItem{
			ID:           p.ID,
			Title:        p.Title,
			MessageCount: len(r.messages[p.ID]),
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
		})
	}

	return items, nil
}

// GetByID returns a project by ID.
func (r *MockProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	project, ok := r.projects[id]
	if !ok {
		return nil, ErrNotFound
	}

	return project, nil
}

// Create creates a new project.
func (r *MockProjectRepository) Create(ctx context.Context, title string) (*model.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	project := &model.Project{
		ID:        uuid.New(),
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.projects[project.ID] = project
	r.messages[project.ID] = []model.Message{}

	return project, nil
}

// Delete removes a project by ID.
func (r *MockProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.projects[id]; !ok {
		return ErrNotFound
	}

	delete(r.projects, id)
	delete(r.messages, id)

	return nil
}

// UpdateTimestamp updates a project's updated_at timestamp.
func (r *MockProjectRepository) UpdateTimestamp(ctx context.Context, id uuid.UUID, timestamp time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	project, ok := r.projects[id]
	if !ok {
		return ErrNotFound
	}

	project.UpdatedAt = timestamp

	return nil
}

// UpdateTitle updates a project's title and updated_at timestamp.
func (r *MockProjectRepository) UpdateTitle(ctx context.Context, id uuid.UUID, title string) (*model.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	project, ok := r.projects[id]
	if !ok {
		return nil, ErrNotFound
	}

	project.Title = title
	project.UpdatedAt = time.Now().UTC()

	return project, nil
}

// GetMessages returns all messages for a project.
func (r *MockProjectRepository) GetMessages(ctx context.Context, projectID uuid.UUID) ([]model.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	messages, ok := r.messages[projectID]
	if !ok {
		return []model.Message{}, nil
	}

	return messages, nil
}

// CreateMessage creates a new message.
func (r *MockProjectRepository) CreateMessage(ctx context.Context, projectID uuid.UUID, role model.Role, content string) (*model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	message := model.Message{
		ID:        uuid.New(),
		ProjectID: projectID,
		Role:      role,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	r.messages[projectID] = append(r.messages[projectID], message)

	// Update project timestamp
	if project, ok := r.projects[projectID]; ok {
		project.UpdatedAt = message.CreatedAt
	}

	return &message, nil
}
