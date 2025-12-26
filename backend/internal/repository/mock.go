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

// CreateMessage creates a new message without agent type.
func (r *MockProjectRepository) CreateMessage(ctx context.Context, projectID uuid.UUID, role model.Role, content string) (*model.Message, error) {
	return r.CreateMessageWithAgent(ctx, projectID, role, content, nil)
}

// CreateMessageWithAgent creates a new message with optional agent type.
func (r *MockProjectRepository) CreateMessageWithAgent(ctx context.Context, projectID uuid.UUID, role model.Role, content string, agentType *string) (*model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	message := model.Message{
		ID:        uuid.New(),
		ProjectID: projectID,
		Role:      role,
		Content:   content,
		AgentType: agentType,
		CreatedAt: time.Now().UTC(),
	}

	r.messages[projectID] = append(r.messages[projectID], message)

	// Update project timestamp
	if project, ok := r.projects[projectID]; ok {
		project.UpdatedAt = message.CreatedAt
	}

	return &message, nil
}

// MockFileRepository implements FileRepository for testing.
type MockFileRepository struct {
	mu    sync.RWMutex
	files map[uuid.UUID]*model.File       // keyed by file ID
	byPath map[string]uuid.UUID           // projectID:path -> fileID lookup
}

// NewMockFileRepository creates a new MockFileRepository.
func NewMockFileRepository() *MockFileRepository {
	return &MockFileRepository{
		files:  make(map[uuid.UUID]*model.File),
		byPath: make(map[string]uuid.UUID),
	}
}

// makePathKey creates a unique key for project+path combination.
func makePathKey(projectID uuid.UUID, path string) string {
	return projectID.String() + ":" + path
}

// SaveFile saves or updates a file for a project (upsert by project_id + path).
func (r *MockFileRepository) SaveFile(ctx context.Context, projectID uuid.UUID, path, language, content string) (*model.File, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pathKey := makePathKey(projectID, path)
	now := time.Now().UTC()

	// Check if file already exists
	if fileID, exists := r.byPath[pathKey]; exists {
		file := r.files[fileID]
		file.Language = language
		file.Content = content
		file.CreatedAt = now // matches upsert behavior in real repo
		return file, nil
	}

	// Create new file
	file := &model.File{
		ID:        uuid.New(),
		ProjectID: projectID,
		Path:      path,
		Filename:  path, // simplified for mock
		Language:  language,
		Content:   content,
		CreatedAt: now,
	}

	r.files[file.ID] = file
	r.byPath[pathKey] = file.ID

	return file, nil
}

// GetFilesByProject returns all files for a project (without content).
func (r *MockFileRepository) GetFilesByProject(ctx context.Context, projectID uuid.UUID) ([]model.FileListItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.FileListItem
	for _, file := range r.files {
		if file.ProjectID == projectID {
			result = append(result, model.FileListItem{
				ID:        file.ID,
				Path:      file.Path,
				Filename:  file.Filename,
				Language:  file.Language,
				CreatedAt: file.CreatedAt,
			})
		}
	}

	return result, nil
}

// GetFilesWithContentByProject returns all files for a project with content.
func (r *MockFileRepository) GetFilesWithContentByProject(ctx context.Context, projectID uuid.UUID) ([]model.File, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.File
	for _, file := range r.files {
		if file.ProjectID == projectID {
			result = append(result, *file)
		}
	}

	return result, nil
}

// GetFile returns a file by ID.
func (r *MockFileRepository) GetFile(ctx context.Context, id uuid.UUID) (*model.File, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, ok := r.files[id]
	if !ok {
		return nil, ErrNotFound
	}

	return file, nil
}

// GetFileByPath returns a file by project ID and path.
func (r *MockFileRepository) GetFileByPath(ctx context.Context, projectID uuid.UUID, path string) (*model.File, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pathKey := makePathKey(projectID, path)
	fileID, exists := r.byPath[pathKey]
	if !exists {
		return nil, ErrNotFound
	}

	file, ok := r.files[fileID]
	if !ok {
		return nil, ErrNotFound
	}

	return file, nil
}
