package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// mockPRDRepo is a mock implementation of PRDRepository for testing.
type mockPRDRepo struct {
	prds      map[uuid.UUID]*model.PRD
	projectID uuid.UUID
}

func newMockPRDRepo() *mockPRDRepo {
	return &mockPRDRepo{
		prds: make(map[uuid.UUID]*model.PRD),
	}
}

func (m *mockPRDRepo) Create(ctx context.Context, prd *model.PRD) (*model.PRD, error) {
	prd.ID = uuid.New()
	m.prds[prd.ID] = prd
	return prd, nil
}

func (m *mockPRDRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.PRD, error) {
	prd, ok := m.prds[id]
	if !ok {
		return nil, ErrPRDNotFound
	}
	return prd, nil
}

func (m *mockPRDRepo) Update(ctx context.Context, prd *model.PRD) (*model.PRD, error) {
	m.prds[prd.ID] = prd
	return prd, nil
}

func (m *mockPRDRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.prds, id)
	return nil
}

func (m *mockPRDRepo) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.PRD, error) {
	var result []model.PRD
	for _, prd := range m.prds {
		if prd.ProjectID == projectID {
			result = append(result, *prd)
		}
	}
	return result, nil
}

func (m *mockPRDRepo) GetByDiscoveryID(ctx context.Context, discoveryID uuid.UUID) ([]model.PRD, error) {
	var result []model.PRD
	for _, prd := range m.prds {
		if prd.DiscoveryID == discoveryID {
			result = append(result, *prd)
		}
	}
	return result, nil
}

func (m *mockPRDRepo) GetByFeatureID(ctx context.Context, featureID uuid.UUID) (*model.PRD, error) {
	for _, prd := range m.prds {
		if prd.FeatureID == featureID {
			return prd, nil
		}
	}
	return nil, ErrPRDNotFound
}

func (m *mockPRDRepo) GetByStatus(ctx context.Context, projectID uuid.UUID, status model.PRDStatus) ([]model.PRD, error) {
	var result []model.PRD
	for _, prd := range m.prds {
		if prd.ProjectID == projectID && prd.Status == status {
			result = append(result, *prd)
		}
	}
	return result, nil
}

func (m *mockPRDRepo) GetByVersion(ctx context.Context, projectID uuid.UUID, version string) ([]model.PRD, error) {
	var result []model.PRD
	for _, prd := range m.prds {
		if prd.ProjectID == projectID && prd.Version == version {
			result = append(result, *prd)
		}
	}
	return result, nil
}

func (m *mockPRDRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status model.PRDStatus) error {
	prd, ok := m.prds[id]
	if !ok {
		return ErrPRDNotFound
	}
	prd.Status = status
	return nil
}

func (m *mockPRDRepo) IncrementGenerationAttempts(ctx context.Context, id uuid.UUID) error {
	prd, ok := m.prds[id]
	if !ok {
		return ErrPRDNotFound
	}
	prd.GenerationAttempts++
	return nil
}

func (m *mockPRDRepo) SetLastError(ctx context.Context, id uuid.UUID, err string) error {
	prd, ok := m.prds[id]
	if !ok {
		return ErrPRDNotFound
	}
	prd.LastError = &err
	return nil
}

func (m *mockPRDRepo) SetGeneratedAt(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockPRDRepo) SetApprovedAt(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockPRDRepo) SetStartedAt(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockPRDRepo) SetCompletedAt(ctx context.Context, id uuid.UUID) error {
	return nil
}

// mockProjectRepo is a mock implementation of ProjectRepository for testing.
type mockProjectRepo struct {
	projects map[uuid.UUID]*model.Project
}

func newMockProjectRepo() *mockProjectRepo {
	return &mockProjectRepo{
		projects: make(map[uuid.UUID]*model.Project),
	}
}

func (m *mockProjectRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	project, ok := m.projects[id]
	if !ok {
		return nil, ErrPRDNotFound
	}
	return project, nil
}

func (m *mockProjectRepo) Create(ctx context.Context, userID uuid.UUID) (*model.Project, error) {
	project := &model.Project{ID: uuid.New()}
	m.projects[project.ID] = project
	return project, nil
}

func (m *mockProjectRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.projects, id)
	return nil
}

func (m *mockProjectRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Project, error) {
	return nil, nil
}

func (m *mockProjectRepo) UpdateTitle(ctx context.Context, id uuid.UUID, title string) (*model.Project, error) {
	project, ok := m.projects[id]
	if !ok {
		return nil, ErrPRDNotFound
	}
	project.Title = title
	return project, nil
}

// mockDiscoveryRepoForAgent is a mock implementation for agent context tests.
type mockDiscoveryRepoForAgent struct {
	discoveries map[uuid.UUID]*model.ProjectDiscovery
	summaries   map[uuid.UUID]*model.DiscoverySummary
}

func newMockDiscoveryRepoForAgent() *mockDiscoveryRepoForAgent {
	return &mockDiscoveryRepoForAgent{
		discoveries: make(map[uuid.UUID]*model.ProjectDiscovery),
		summaries:   make(map[uuid.UUID]*model.DiscoverySummary),
	}
}

func (m *mockDiscoveryRepoForAgent) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*model.ProjectDiscovery, error) {
	for _, d := range m.discoveries {
		if d.ProjectID == projectID {
			return d, nil
		}
	}
	return nil, ErrDiscoveryNotFound
}

func (m *mockDiscoveryRepoForAgent) GetByID(ctx context.Context, id uuid.UUID) (*model.ProjectDiscovery, error) {
	d, ok := m.discoveries[id]
	if !ok {
		return nil, ErrDiscoveryNotFound
	}
	return d, nil
}

func (m *mockDiscoveryRepoForAgent) GetSummary(ctx context.Context, id uuid.UUID) (*model.DiscoverySummary, error) {
	s, ok := m.summaries[id]
	if !ok {
		return nil, ErrDiscoveryNotFound
	}
	return s, nil
}

func (m *mockDiscoveryRepoForAgent) Create(ctx context.Context, projectID uuid.UUID) (*model.ProjectDiscovery, error) {
	return nil, nil
}

func (m *mockDiscoveryRepoForAgent) Update(ctx context.Context, discovery *model.ProjectDiscovery) (*model.ProjectDiscovery, error) {
	return nil, nil
}

func (m *mockDiscoveryRepoForAgent) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockDiscoveryRepoForAgent) UpdateStage(ctx context.Context, id uuid.UUID, stage model.DiscoveryStage) (*model.ProjectDiscovery, error) {
	return nil, nil
}

func (m *mockDiscoveryRepoForAgent) MarkComplete(ctx context.Context, id uuid.UUID) (*model.ProjectDiscovery, error) {
	return nil, nil
}

func (m *mockDiscoveryRepoForAgent) AddUser(ctx context.Context, user *model.DiscoveryUser) (*model.DiscoveryUser, error) {
	return nil, nil
}

func (m *mockDiscoveryRepoForAgent) GetUsers(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryUser, error) {
	return nil, nil
}

func (m *mockDiscoveryRepoForAgent) AddFeature(ctx context.Context, feature *model.DiscoveryFeature) (*model.DiscoveryFeature, error) {
	return nil, nil
}

func (m *mockDiscoveryRepoForAgent) GetMVPFeatures(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryFeature, error) {
	return nil, nil
}

func (m *mockDiscoveryRepoForAgent) GetFutureFeatures(ctx context.Context, discoveryID uuid.UUID) ([]model.DiscoveryFeature, error) {
	return nil, nil
}

// Helper to create a test service
func newTestAgentContextService() (*AgentContextService, *mockPRDRepo, *mockProjectRepo, *mockDiscoveryRepoForAgent) {
	prdRepo := newMockPRDRepo()
	projectRepo := newMockProjectRepo()
	discoveryRepo := newMockDiscoveryRepoForAgent()
	logger := zerolog.Nop()

	svc := NewAgentContextService(prdRepo, projectRepo, discoveryRepo, logger)
	return svc, prdRepo, projectRepo, discoveryRepo
}

func TestSelectAgent_ProductManager(t *testing.T) {
	svc, _, _, _ := newTestAgentContextService()

	tests := []struct {
		message string
		want    model.AgentType
	}{
		{"What should we include in the scope?", model.AgentProductManager},
		{"Can you explain this feature?", model.AgentProductManager},
		{"What are the requirements for login?", model.AgentProductManager},
		{"Let me tell you a user story", model.AgentProductManager},
		{"What's the priority here?", model.AgentProductManager},
		{"Show me the roadmap", model.AgentProductManager},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got := svc.SelectAgent(tt.message, nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSelectAgent_Designer(t *testing.T) {
	svc, _, _, _ := newTestAgentContextService()

	tests := []struct {
		message string
		want    model.AgentType
	}{
		{"How should we design the login page?", model.AgentDesigner},
		{"What's the best layout for the dashboard?", model.AgentDesigner},
		{"Can you help with the UI?", model.AgentDesigner},
		{"I need help with the UX flow", model.AgentDesigner},
		{"What should the look and feel be?", model.AgentDesigner},
		{"Let's pick a color scheme", model.AgentDesigner},
		{"Create a wireframe for this", model.AgentDesigner},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got := svc.SelectAgent(tt.message, nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSelectAgent_Developer(t *testing.T) {
	svc, _, _, _ := newTestAgentContextService()

	tests := []struct {
		message string
		want    model.AgentType
	}{
		{"Build the login form", model.AgentDeveloper},
		{"Write the API endpoint", model.AgentDeveloper},
		{"Implement the database schema", model.AgentDeveloper},
		{"Add a button that submits", model.AgentDeveloper},
		{"Create the user model", model.AgentDeveloper},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got := svc.SelectAgent(tt.message, nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMatchPRDByKeywords(t *testing.T) {
	svc, _, _, _ := newTestAgentContextService()

	projectID := uuid.New()
	prds := []model.PRD{
		{ID: uuid.New(), ProjectID: projectID, Title: "User Authentication", Status: model.PRDStatusReady},
		{ID: uuid.New(), ProjectID: projectID, Title: "Order Management System", Status: model.PRDStatusReady},
		{ID: uuid.New(), ProjectID: projectID, Title: "Product Catalog", Status: model.PRDStatusReady},
	}

	tests := []struct {
		message   string
		wantTitle string
	}{
		{"Add authentication for users", "User Authentication"},
		{"Build the order management feature", "Order Management System"},
		{"Show me the product catalog", "Product Catalog"},
		{"I need to manage orders", "Order Management System"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got := svc.matchPRDByKeywords(tt.message, prds)
			require.NotNil(t, got, "expected a PRD match for: %s", tt.message)
			assert.Equal(t, tt.wantTitle, got.Title)
		})
	}
}

func TestMatchPRDByKeywords_NoMatch(t *testing.T) {
	svc, _, _, _ := newTestAgentContextService()

	prds := []model.PRD{
		{ID: uuid.New(), Title: "User Authentication"},
		{ID: uuid.New(), Title: "Order Management"},
	}

	got := svc.matchPRDByKeywords("something completely unrelated", prds)
	assert.Nil(t, got)
}

func TestGetNextReadyPRD(t *testing.T) {
	svc, _, _, _ := newTestAgentContextService()

	prds := []model.PRD{
		{ID: uuid.New(), Title: "Feature A", Status: model.PRDStatusComplete, Priority: 1},
		{ID: uuid.New(), Title: "Feature B", Status: model.PRDStatusReady, Priority: 2},
		{ID: uuid.New(), Title: "Feature C", Status: model.PRDStatusReady, Priority: 1},
		{ID: uuid.New(), Title: "Feature D", Status: model.PRDStatusDraft, Priority: 0},
	}

	got := svc.getNextReadyPRD(prds)
	require.NotNil(t, got)
	assert.Equal(t, "Feature C", got.Title) // Highest priority (lowest number) among ready PRDs
}

func TestCondensePRD(t *testing.T) {
	svc, _, _, _ := newTestAgentContextService()

	prd := &model.PRD{
		Title:    "User Login",
		Overview: "Allows users to authenticate with the system.",
	}

	// Add user stories
	stories := []model.UserStory{
		{ID: "US-001", AsA: "user", IWant: "to log in with email", SoThat: "I can access my account"},
		{ID: "US-002", AsA: "admin", IWant: "to reset passwords", SoThat: "I can help locked out users"},
	}
	err := prd.SetUserStories(stories)
	require.NoError(t, err)

	// Add acceptance criteria
	criteria := []model.AcceptanceCriterion{
		{ID: "AC-001", Given: "valid credentials", When: "user submits login form", Then: "user is logged in"},
		{ID: "AC-002", Given: "invalid password", When: "user submits login form", Then: "error is shown"},
	}
	err = prd.SetAcceptanceCriteria(criteria)
	require.NoError(t, err)

	summary := svc.CondensePRD(prd)

	assert.Contains(t, summary, "## User Login")
	assert.Contains(t, summary, "Allows users to authenticate")
	assert.Contains(t, summary, "US-001")
	assert.Contains(t, summary, "to log in with email")
	assert.Contains(t, summary, "AC-001")
	assert.Contains(t, summary, "valid credentials")
}

func TestGetContextForMessage(t *testing.T) {
	svc, prdRepo, projectRepo, discoveryRepo := newTestAgentContextService()
	ctx := context.Background()

	// Setup project
	projectID := uuid.New()
	project := &model.Project{ID: projectID, Title: "Test Project"}
	projectRepo.projects[projectID] = project

	// Setup discovery
	discoveryID := uuid.New()
	discovery := &model.ProjectDiscovery{ID: discoveryID, ProjectID: projectID, Stage: model.StageComplete}
	discoveryRepo.discoveries[discoveryID] = discovery
	discoveryRepo.summaries[discoveryID] = &model.DiscoverySummary{
		ProjectName:     "Test Project",
		SolvesStatement: "Solves inventory tracking",
	}

	// Setup PRDs
	prd1 := &model.PRD{
		ID:        uuid.New(),
		ProjectID: projectID,
		Title:     "Order Management",
		Overview:  "Manage customer orders",
		Status:    model.PRDStatusReady,
		Priority:  1,
	}
	prdRepo.prds[prd1.ID] = prd1

	// Test getting context
	agentCtx, err := svc.GetContextForMessage(ctx, projectID, "Build the order management feature")
	require.NoError(t, err)
	require.NotNil(t, agentCtx)

	assert.Equal(t, model.AgentDeveloper, agentCtx.Agent)
	require.NotNil(t, agentCtx.PRD)
	assert.Equal(t, "Order Management", agentCtx.PRD.Title)
	require.NotNil(t, agentCtx.Discovery)
	assert.Equal(t, "Test Project", agentCtx.Discovery.ProjectName)
}

func TestGetSystemPrompt_ProductManager(t *testing.T) {
	svc, _, _, _ := newTestAgentContextService()
	ctx := context.Background()

	prd := &model.PRD{
		Title:    "User Login",
		Overview: "Authentication feature",
	}
	err := prd.SetUserStories([]model.UserStory{
		{ID: "US-001", AsA: "user", IWant: "to login", SoThat: "I can access the app"},
	})
	require.NoError(t, err)

	agentCtx := &model.AgentContext{
		Agent:      model.AgentProductManager,
		PRD:        prd,
		PRDSummary: svc.CondensePRD(prd),
		Discovery: &model.DiscoverySummary{
			ProjectName:     "MyApp",
			SolvesStatement: "Helps users manage tasks",
			Users: []model.DiscoveryUser{
				{Description: "Regular users"},
			},
		},
	}

	prompt, err := svc.GetSystemPrompt(ctx, agentCtx)
	require.NoError(t, err)

	assert.Contains(t, prompt, "Product Manager")
	assert.Contains(t, prompt, "MyApp")
	assert.Contains(t, prompt, "User Login")
	assert.Contains(t, prompt, "Regular users")
}

func TestGetSystemPrompt_Designer(t *testing.T) {
	svc, _, _, _ := newTestAgentContextService()
	ctx := context.Background()

	agentCtx := &model.AgentContext{
		Agent:      model.AgentDesigner,
		PRDSummary: "## Dashboard\n\nMain user dashboard.",
		Discovery: &model.DiscoverySummary{
			ProjectName: "DesignApp",
			Users: []model.DiscoveryUser{
				{Description: "Designers", UserCount: 50},
			},
		},
	}

	prompt, err := svc.GetSystemPrompt(ctx, agentCtx)
	require.NoError(t, err)

	assert.Contains(t, prompt, "UI/UX Designer")
	assert.Contains(t, prompt, "DesignApp")
	assert.Contains(t, prompt, "Dashboard")
	assert.Contains(t, prompt, "Designers")
}

func TestGetSystemPrompt_Developer(t *testing.T) {
	svc, _, _, _ := newTestAgentContextService()
	ctx := context.Background()

	prd := &model.PRD{
		Title:    "API Gateway",
		Overview: "Central API routing",
	}
	err := prd.SetTechnicalNotes([]model.TechnicalNote{
		{Category: "architecture", Description: "Use REST conventions"},
	})
	require.NoError(t, err)

	agentCtx := &model.AgentContext{
		Agent:      model.AgentDeveloper,
		PRD:        prd,
		PRDSummary: svc.CondensePRD(prd),
		Discovery: &model.DiscoverySummary{
			ProjectName: "DevProject",
		},
	}

	prompt, err := svc.GetSystemPrompt(ctx, agentCtx)
	require.NoError(t, err)

	assert.Contains(t, prompt, "Developer")
	assert.Contains(t, prompt, "DevProject")
	assert.Contains(t, prompt, "API Gateway")
	assert.Contains(t, prompt, "architecture")
}
