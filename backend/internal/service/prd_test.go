package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

// MockPRDRepository implements PRDRepository for testing.
type MockPRDRepository struct {
	mu   sync.RWMutex
	prds map[uuid.UUID]*model.PRD
}

// NewMockPRDRepository creates a new MockPRDRepository.
func NewMockPRDRepository() *MockPRDRepository {
	return &MockPRDRepository{
		prds: make(map[uuid.UUID]*model.PRD),
	}
}

func (r *MockPRDRepository) Create(ctx context.Context, prd *model.PRD) (*model.PRD, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	prd.ID = uuid.New()
	prd.CreatedAt = time.Now()
	prd.UpdatedAt = time.Now()

	copy := *prd
	r.prds[prd.ID] = &copy
	return &copy, nil
}

func (r *MockPRDRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PRD, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prd, ok := r.prds[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	copy := *prd
	return &copy, nil
}

func (r *MockPRDRepository) Update(ctx context.Context, prd *model.PRD) (*model.PRD, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.prds[prd.ID]; !ok {
		return nil, repository.ErrNotFound
	}
	prd.UpdatedAt = time.Now()
	copy := *prd
	r.prds[prd.ID] = &copy
	return &copy, nil
}

func (r *MockPRDRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.prds[id]; !ok {
		return repository.ErrNotFound
	}
	delete(r.prds, id)
	return nil
}

func (r *MockPRDRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.PRD, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.PRD
	for _, prd := range r.prds {
		if prd.ProjectID == projectID {
			result = append(result, *prd)
		}
	}
	return result, nil
}

func (r *MockPRDRepository) GetByDiscoveryID(ctx context.Context, discoveryID uuid.UUID) ([]model.PRD, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.PRD
	for _, prd := range r.prds {
		if prd.DiscoveryID == discoveryID {
			result = append(result, *prd)
		}
	}
	return result, nil
}

func (r *MockPRDRepository) GetByFeatureID(ctx context.Context, featureID uuid.UUID) (*model.PRD, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, prd := range r.prds {
		if prd.FeatureID == featureID {
			copy := *prd
			return &copy, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *MockPRDRepository) GetByStatus(ctx context.Context, projectID uuid.UUID, status model.PRDStatus) ([]model.PRD, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.PRD
	for _, prd := range r.prds {
		if prd.ProjectID == projectID && prd.Status == status {
			result = append(result, *prd)
		}
	}
	return result, nil
}

func (r *MockPRDRepository) GetByVersion(ctx context.Context, projectID uuid.UUID, version string) ([]model.PRD, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.PRD
	for _, prd := range r.prds {
		if prd.ProjectID == projectID && prd.Version == version {
			result = append(result, *prd)
		}
	}
	return result, nil
}

func (r *MockPRDRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.PRDStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	prd, ok := r.prds[id]
	if !ok {
		return repository.ErrNotFound
	}
	prd.Status = status
	prd.UpdatedAt = time.Now()
	return nil
}

func (r *MockPRDRepository) IncrementGenerationAttempts(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	prd, ok := r.prds[id]
	if !ok {
		return repository.ErrNotFound
	}
	prd.GenerationAttempts++
	prd.UpdatedAt = time.Now()
	return nil
}

func (r *MockPRDRepository) SetLastError(ctx context.Context, id uuid.UUID, err string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	prd, ok := r.prds[id]
	if !ok {
		return repository.ErrNotFound
	}
	prd.LastError = &err
	prd.UpdatedAt = time.Now()
	return nil
}

func (r *MockPRDRepository) SetGeneratedAt(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	prd, ok := r.prds[id]
	if !ok {
		return repository.ErrNotFound
	}
	now := time.Now()
	prd.GeneratedAt = &now
	prd.UpdatedAt = now
	return nil
}

func (r *MockPRDRepository) SetApprovedAt(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	prd, ok := r.prds[id]
	if !ok {
		return repository.ErrNotFound
	}
	now := time.Now()
	prd.ApprovedAt = &now
	prd.UpdatedAt = now
	return nil
}

func (r *MockPRDRepository) SetStartedAt(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	prd, ok := r.prds[id]
	if !ok {
		return repository.ErrNotFound
	}
	now := time.Now()
	prd.StartedAt = &now
	prd.UpdatedAt = now
	return nil
}

func (r *MockPRDRepository) SetCompletedAt(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	prd, ok := r.prds[id]
	if !ok {
		return repository.ErrNotFound
	}
	now := time.Now()
	prd.CompletedAt = &now
	prd.UpdatedAt = now
	return nil
}

// MockClaudeMessengerForPRD is a mock Claude messenger for PRD generation tests.
type MockClaudeMessengerForPRD struct {
	Response string
	Err      error
}

func (m *MockClaudeMessengerForPRD) SendMessage(ctx context.Context, systemPrompt string, messages []ClaudeMessage) (*ClaudeStream, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	stream := &ClaudeStream{
		chunks: make(chan string, 1),
		done:   make(chan struct{}),
	}

	go func() {
		stream.chunks <- m.Response
		close(stream.chunks)
		close(stream.done)
	}()

	return stream, nil
}

func newTestPRDService() (*PRDService, *MockPRDRepository, *repository.MockDiscoveryRepository, *MockClaudeMessengerForPRD) {
	prdRepo := NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	claudeMock := &MockClaudeMessengerForPRD{
		Response: `{
			"overview": "Test feature overview",
			"userStories": [
				{
					"id": "US-001",
					"asA": "user",
					"iWant": "to test",
					"soThat": "I can verify",
					"priority": "must",
					"complexity": "low"
				}
			],
			"acceptanceCriteria": [
				{
					"id": "AC-001",
					"given": "a test setup",
					"when": "I run tests",
					"then": "they pass",
					"userStoryId": "US-001"
				}
			],
			"technicalNotes": [
				{
					"category": "data",
					"title": "Test note",
					"description": "Test description",
					"suggestions": ["suggestion1"]
				}
			]
		}`,
	}
	logger := zerolog.Nop()
	service := NewPRDService(prdRepo, discoveryRepo, claudeMock, logger)
	return service, prdRepo, discoveryRepo, claudeMock
}

func setupDiscoveryWithFeatures(ctx context.Context, discoveryRepo *repository.MockDiscoveryRepository) (uuid.UUID, uuid.UUID) {
	projectID := uuid.New()
	discovery, _ := discoveryRepo.Create(ctx, projectID)

	// Add project name
	pn := "Test Project"
	ss := "Solves test problems"
	discovery.ProjectName = &pn
	discovery.SolvesStatement = &ss
	discoveryRepo.Update(ctx, discovery)

	// Add MVP feature
	discoveryRepo.AddFeature(ctx, &model.DiscoveryFeature{
		DiscoveryID: discovery.ID,
		Name:        "Test Feature 1",
		Priority:    1,
		Version:     "v1",
	})

	// Add future feature
	discoveryRepo.AddFeature(ctx, &model.DiscoveryFeature{
		DiscoveryID: discovery.ID,
		Name:        "Future Feature",
		Priority:    1,
		Version:     "v2",
	})

	// Add user
	discoveryRepo.AddUser(ctx, &model.DiscoveryUser{
		DiscoveryID:    discovery.ID,
		Description:    "Test User",
		UserCount:      5,
		HasPermissions: false,
	})

	return discovery.ID, projectID
}

func TestNewPRDService(t *testing.T) {
	service, _, _, _ := newTestPRDService()
	assert.NotNil(t, service)
}

func TestGeneratePRD_Success(t *testing.T) {
	service, prdRepo, discoveryRepo, _ := newTestPRDService()
	ctx := context.Background()

	discoveryID, projectID := setupDiscoveryWithFeatures(ctx, discoveryRepo)

	// Get MVP feature
	mvpFeatures, _ := discoveryRepo.GetMVPFeatures(ctx, discoveryID)
	require.Len(t, mvpFeatures, 1)

	// Create pending PRD
	prd := &model.PRD{
		DiscoveryID: discoveryID,
		FeatureID:   mvpFeatures[0].ID,
		ProjectID:   projectID,
		Title:       mvpFeatures[0].Name,
		Version:     "v1",
		Priority:    1,
		Status:      model.PRDStatusPending,
	}
	prd.SetUserStories([]model.UserStory{})
	prd.SetAcceptanceCriteria([]model.AcceptanceCriterion{})
	prd.SetTechnicalNotes([]model.TechnicalNote{})

	created, err := prdRepo.Create(ctx, prd)
	require.NoError(t, err)

	// Generate PRD
	result, err := service.GeneratePRD(ctx, created.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, model.PRDStatusDraft, result.Status)
	assert.Equal(t, "Test feature overview", result.Overview)

	// Verify user stories were saved
	stories, err := result.UserStories()
	require.NoError(t, err)
	assert.Len(t, stories, 1)
	assert.Equal(t, "US-001", stories[0].ID)
}

func TestGeneratePRD_AlreadyGenerated(t *testing.T) {
	service, prdRepo, discoveryRepo, _ := newTestPRDService()
	ctx := context.Background()

	discoveryID, projectID := setupDiscoveryWithFeatures(ctx, discoveryRepo)
	mvpFeatures, _ := discoveryRepo.GetMVPFeatures(ctx, discoveryID)

	// Create PRD already in draft status
	prd := &model.PRD{
		DiscoveryID: discoveryID,
		FeatureID:   mvpFeatures[0].ID,
		ProjectID:   projectID,
		Title:       mvpFeatures[0].Name,
		Version:     "v1",
		Priority:    1,
		Status:      model.PRDStatusDraft,
		Overview:    "Existing overview",
	}
	prd.SetUserStories([]model.UserStory{})
	prd.SetAcceptanceCriteria([]model.AcceptanceCriterion{})
	prd.SetTechnicalNotes([]model.TechnicalNote{})

	created, _ := prdRepo.Create(ctx, prd)

	// Generate should return existing without calling Claude
	result, err := service.GeneratePRD(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Existing overview", result.Overview)
}

func TestGetByID_Success(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Status:      model.PRDStatusDraft,
	}
	created, _ := prdRepo.Create(ctx, prd)

	result, err := service.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, "Test PRD", result.Title)
}

func TestGetByID_NotFound(t *testing.T) {
	service, _, _, _ := newTestPRDService()
	ctx := context.Background()

	_, err := service.GetByID(ctx, uuid.New())
	assert.Equal(t, ErrPRDNotFound, err)
}

func TestGetByProjectID(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	projectID := uuid.New()

	// Create two PRDs for the project
	prd1 := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "PRD 1",
		Status:      model.PRDStatusDraft,
	}
	prd2 := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "PRD 2",
		Status:      model.PRDStatusDraft,
	}
	prdRepo.Create(ctx, prd1)
	prdRepo.Create(ctx, prd2)

	// Create one for another project
	prd3 := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "PRD 3",
		Status:      model.PRDStatusDraft,
	}
	prdRepo.Create(ctx, prd3)

	prds, err := service.GetByProjectID(ctx, projectID)
	require.NoError(t, err)
	assert.Len(t, prds, 2)
}

func TestGetMVPPRDs(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	projectID := uuid.New()

	// Create MVP PRD
	mvp := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "MVP Feature",
		Version:     "v1",
		Status:      model.PRDStatusDraft,
	}
	prdRepo.Create(ctx, mvp)

	// Create future PRD
	future := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "Future Feature",
		Version:     "v2",
		Status:      model.PRDStatusDraft,
	}
	prdRepo.Create(ctx, future)

	prds, err := service.GetMVPPRDs(ctx, projectID)
	require.NoError(t, err)
	assert.Len(t, prds, 1)
	assert.Equal(t, "v1", prds[0].Version)
}

func TestUpdateStatus_ValidTransition(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Status:      model.PRDStatusDraft,
	}
	created, _ := prdRepo.Create(ctx, prd)

	// Draft -> Ready is valid
	err := service.UpdateStatus(ctx, created.ID, model.PRDStatusReady)
	require.NoError(t, err)

	updated, _ := prdRepo.GetByID(ctx, created.ID)
	assert.Equal(t, model.PRDStatusReady, updated.Status)
}

func TestUpdateStatus_InvalidTransition(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Status:      model.PRDStatusDraft,
	}
	created, _ := prdRepo.Create(ctx, prd)

	// Draft -> Complete is invalid (must go through ready and in_progress)
	err := service.UpdateStatus(ctx, created.ID, model.PRDStatusComplete)
	assert.ErrorIs(t, err, ErrInvalidStatusChange)
}

func TestMarkAsReady(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Status:      model.PRDStatusDraft,
	}
	created, _ := prdRepo.Create(ctx, prd)

	err := service.MarkAsReady(ctx, created.ID)
	require.NoError(t, err)

	updated, _ := prdRepo.GetByID(ctx, created.ID)
	assert.Equal(t, model.PRDStatusReady, updated.Status)
	assert.NotNil(t, updated.ApprovedAt)
}

func TestStartImplementation(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Status:      model.PRDStatusReady,
	}
	created, _ := prdRepo.Create(ctx, prd)

	err := service.StartImplementation(ctx, created.ID)
	require.NoError(t, err)

	updated, _ := prdRepo.GetByID(ctx, created.ID)
	assert.Equal(t, model.PRDStatusInProgress, updated.Status)
	assert.NotNil(t, updated.StartedAt)
}

func TestCompleteImplementation(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Status:      model.PRDStatusInProgress,
	}
	created, _ := prdRepo.Create(ctx, prd)

	err := service.CompleteImplementation(ctx, created.ID)
	require.NoError(t, err)

	updated, _ := prdRepo.GetByID(ctx, created.ID)
	assert.Equal(t, model.PRDStatusComplete, updated.Status)
	assert.NotNil(t, updated.CompletedAt)
}

func TestGetActivePRD(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	projectID := uuid.New()

	// Create in-progress PRD
	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "Active PRD",
		Status:      model.PRDStatusInProgress,
	}
	prdRepo.Create(ctx, prd)

	result, err := service.GetActivePRD(ctx, projectID)
	require.NoError(t, err)
	assert.Equal(t, "Active PRD", result.Title)
}

func TestGetActivePRD_NoActive(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	projectID := uuid.New()

	// Create ready PRD (not in progress)
	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "Ready PRD",
		Status:      model.PRDStatusReady,
	}
	prdRepo.Create(ctx, prd)

	_, err := service.GetActivePRD(ctx, projectID)
	assert.Equal(t, ErrNoActivePRD, err)
}

func TestSetActivePRD(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	projectID := uuid.New()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "Ready PRD",
		Status:      model.PRDStatusReady,
	}
	created, _ := prdRepo.Create(ctx, prd)

	err := service.SetActivePRD(ctx, projectID, created.ID)
	require.NoError(t, err)

	updated, _ := prdRepo.GetByID(ctx, created.ID)
	assert.Equal(t, model.PRDStatusInProgress, updated.Status)
}

func TestSetActivePRD_NotReady(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	projectID := uuid.New()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "Draft PRD",
		Status:      model.PRDStatusDraft,
	}
	created, _ := prdRepo.Create(ctx, prd)

	err := service.SetActivePRD(ctx, projectID, created.ID)
	assert.ErrorIs(t, err, ErrInvalidStatusChange)
}

func TestClearActivePRD(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	projectID := uuid.New()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "Active PRD",
		Status:      model.PRDStatusInProgress,
	}
	prdRepo.Create(ctx, prd)

	err := service.ClearActivePRD(ctx, projectID)
	require.NoError(t, err)

	// Should not have any active PRDs now
	_, err = service.GetActivePRD(ctx, projectID)
	assert.Equal(t, ErrNoActivePRD, err)
}

func TestGetNextPRD(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	projectID := uuid.New()

	// Create multiple ready PRDs with different priorities
	prd1 := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "Lower priority",
		Priority:    2,
		Status:      model.PRDStatusReady,
	}
	prd2 := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "Higher priority",
		Priority:    1,
		Status:      model.PRDStatusReady,
	}
	prdRepo.Create(ctx, prd1)
	prdRepo.Create(ctx, prd2)

	result, err := service.GetNextPRD(ctx, projectID)
	require.NoError(t, err)
	assert.Equal(t, "Higher priority", result.Title)
}

func TestGetNextPRD_NoReady(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	projectID := uuid.New()

	// Create only draft PRDs
	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   projectID,
		Title:       "Draft PRD",
		Status:      model.PRDStatusDraft,
	}
	prdRepo.Create(ctx, prd)

	_, err := service.GetNextPRD(ctx, projectID)
	assert.Equal(t, ErrNoReadyPRD, err)
}

func TestUpdateOverview(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Overview:    "Original overview",
		Status:      model.PRDStatusDraft,
	}
	created, _ := prdRepo.Create(ctx, prd)

	err := service.UpdateOverview(ctx, created.ID, "Updated overview")
	require.NoError(t, err)

	updated, _ := prdRepo.GetByID(ctx, created.ID)
	assert.Equal(t, "Updated overview", updated.Overview)
}

func TestUpdateOverview_NotDraft(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Overview:    "Original overview",
		Status:      model.PRDStatusReady,
	}
	created, _ := prdRepo.Create(ctx, prd)

	err := service.UpdateOverview(ctx, created.ID, "Updated overview")
	assert.ErrorIs(t, err, ErrInvalidStatusChange)
}

func TestAddUserStory(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Status:      model.PRDStatusDraft,
	}
	prd.SetUserStories([]model.UserStory{})
	created, _ := prdRepo.Create(ctx, prd)

	story := &model.UserStory{
		AsA:        "user",
		IWant:      "to test",
		SoThat:     "I verify",
		Priority:   "must",
		Complexity: "low",
	}

	err := service.AddUserStory(ctx, created.ID, story)
	require.NoError(t, err)

	updated, _ := prdRepo.GetByID(ctx, created.ID)
	stories, _ := updated.UserStories()
	assert.Len(t, stories, 1)
	assert.Equal(t, "US-001", stories[0].ID)
}

func TestUpdateUserStory(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Status:      model.PRDStatusDraft,
	}
	prd.SetUserStories([]model.UserStory{
		{ID: "US-001", AsA: "user", IWant: "original", SoThat: "I test", Priority: "must", Complexity: "low"},
	})
	created, _ := prdRepo.Create(ctx, prd)

	updatedStory := &model.UserStory{
		AsA:        "admin",
		IWant:      "updated",
		SoThat:     "I verify",
		Priority:   "should",
		Complexity: "medium",
	}

	err := service.UpdateUserStory(ctx, created.ID, "US-001", updatedStory)
	require.NoError(t, err)

	updated, _ := prdRepo.GetByID(ctx, created.ID)
	stories, _ := updated.UserStories()
	assert.Len(t, stories, 1)
	assert.Equal(t, "US-001", stories[0].ID) // ID preserved
	assert.Equal(t, "admin", stories[0].AsA)
	assert.Equal(t, "updated", stories[0].IWant)
}

func TestDeleteUserStory(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Status:      model.PRDStatusDraft,
	}
	prd.SetUserStories([]model.UserStory{
		{ID: "US-001", AsA: "user", IWant: "first", SoThat: "I test", Priority: "must", Complexity: "low"},
		{ID: "US-002", AsA: "admin", IWant: "second", SoThat: "I verify", Priority: "should", Complexity: "medium"},
	})
	created, _ := prdRepo.Create(ctx, prd)

	err := service.DeleteUserStory(ctx, created.ID, "US-001")
	require.NoError(t, err)

	updated, _ := prdRepo.GetByID(ctx, created.ID)
	stories, _ := updated.UserStories()
	assert.Len(t, stories, 1)
	assert.Equal(t, "US-002", stories[0].ID)
}

func TestRetryGeneration(t *testing.T) {
	service, prdRepo, discoveryRepo, _ := newTestPRDService()
	ctx := context.Background()

	discoveryID, projectID := setupDiscoveryWithFeatures(ctx, discoveryRepo)
	mvpFeatures, _ := discoveryRepo.GetMVPFeatures(ctx, discoveryID)

	// Create failed PRD
	prd := &model.PRD{
		DiscoveryID: discoveryID,
		FeatureID:   mvpFeatures[0].ID,
		ProjectID:   projectID,
		Title:       mvpFeatures[0].Name,
		Version:     "v1",
		Priority:    1,
		Status:      model.PRDStatusFailed,
	}
	prd.SetUserStories([]model.UserStory{})
	prd.SetAcceptanceCriteria([]model.AcceptanceCriterion{})
	prd.SetTechnicalNotes([]model.TechnicalNote{})
	created, _ := prdRepo.Create(ctx, prd)

	result, err := service.RetryGeneration(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, model.PRDStatusDraft, result.Status)
}

func TestRetryGeneration_InvalidStatus(t *testing.T) {
	service, prdRepo, _, _ := newTestPRDService()
	ctx := context.Background()

	// Create ready PRD (not failed)
	prd := &model.PRD{
		DiscoveryID: uuid.New(),
		FeatureID:   uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test PRD",
		Status:      model.PRDStatusReady,
	}
	created, _ := prdRepo.Create(ctx, prd)

	_, err := service.RetryGeneration(ctx, created.ID)
	assert.ErrorIs(t, err, ErrInvalidStatusChange)
}

func TestIsValidStatusTransition(t *testing.T) {
	service, _, _, _ := newTestPRDService()

	tests := []struct {
		from    model.PRDStatus
		to      model.PRDStatus
		isValid bool
	}{
		{model.PRDStatusPending, model.PRDStatusGenerating, true},
		{model.PRDStatusPending, model.PRDStatusFailed, true},
		{model.PRDStatusPending, model.PRDStatusDraft, false},
		{model.PRDStatusGenerating, model.PRDStatusDraft, true},
		{model.PRDStatusGenerating, model.PRDStatusFailed, true},
		{model.PRDStatusDraft, model.PRDStatusReady, true},
		{model.PRDStatusDraft, model.PRDStatusComplete, false},
		{model.PRDStatusReady, model.PRDStatusInProgress, true},
		{model.PRDStatusInProgress, model.PRDStatusComplete, true},
		{model.PRDStatusInProgress, model.PRDStatusReady, true},
		{model.PRDStatusComplete, model.PRDStatusReady, false},
		{model.PRDStatusFailed, model.PRDStatusPending, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"_to_"+string(tt.to), func(t *testing.T) {
			result := service.isValidStatusTransition(tt.from, tt.to)
			assert.Equal(t, tt.isValid, result)
		})
	}
}
