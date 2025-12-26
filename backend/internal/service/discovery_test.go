package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

func newTestDiscoveryService() (*DiscoveryService, *repository.MockDiscoveryRepository) {
	repo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	service := NewDiscoveryService(repo, logger)
	return service, repo
}

func TestGetOrCreateDiscovery_CreatesNew(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// First call should create a new discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)
	assert.NotNil(t, discovery)
	assert.Equal(t, projectID, discovery.ProjectID)
	assert.Equal(t, model.StageWelcome, discovery.Stage)
}

func TestGetOrCreateDiscovery_ReturnsExisting(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create first discovery
	discovery1, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Second call should return the same discovery
	discovery2, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)
	assert.Equal(t, discovery1.ID, discovery2.ID)
	assert.Equal(t, discovery1.ProjectID, discovery2.ProjectID)
}

func TestAdvanceStage_ProgressesThroughStages(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery (starts at welcome)
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)
	assert.Equal(t, model.StageWelcome, discovery.Stage)

	// Advance to problem
	discovery, err = service.AdvanceStage(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Equal(t, model.StageProblem, discovery.Stage)

	// Advance to personas
	discovery, err = service.AdvanceStage(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Equal(t, model.StagePersonas, discovery.Stage)

	// Advance to mvp
	discovery, err = service.AdvanceStage(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Equal(t, model.StageMVP, discovery.Stage)

	// Advance to summary
	discovery, err = service.AdvanceStage(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Equal(t, model.StageSummary, discovery.Stage)
}

func TestAdvanceStage_ErrorsOnComplete(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create and complete discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Advance through all stages
	for i := 0; i < 4; i++ {
		discovery, err = service.AdvanceStage(ctx, discovery.ID)
		require.NoError(t, err)
	}

	// Now at summary, confirm it
	discovery, err = service.ConfirmDiscovery(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Equal(t, model.StageComplete, discovery.Stage)

	// Try to advance - should error
	_, err = service.AdvanceStage(ctx, discovery.ID)
	assert.Equal(t, ErrDiscoveryAlreadyComplete, err)
}

func TestAdvanceStage_ErrorsOnNotFound(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()

	_, err := service.AdvanceStage(ctx, uuid.New())
	assert.Equal(t, ErrDiscoveryNotFound, err)
}

func TestGetSystemPrompt_ReturnsStageAppropriatePrompt(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery
	_, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Get welcome stage prompt
	prompt, err := service.GetSystemPrompt(ctx, projectID)
	require.NoError(t, err)
	assert.Contains(t, prompt, "Welcome (1 of 5)")
	assert.Contains(t, prompt, "Product Guide")
}

func TestGetSystemPrompt_ReturnsEmptyForComplete(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create and complete discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Advance through all stages and confirm
	for i := 0; i < 4; i++ {
		discovery, err = service.AdvanceStage(ctx, discovery.ID)
		require.NoError(t, err)
	}
	_, err = service.ConfirmDiscovery(ctx, discovery.ID)
	require.NoError(t, err)

	// Get prompt - should be empty for complete
	prompt, err := service.GetSystemPrompt(ctx, projectID)
	require.NoError(t, err)
	assert.Empty(t, prompt)
}

func TestGetSystemPrompt_ReturnsEmptyForNoDiscovery(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// No discovery exists - should return empty string
	prompt, err := service.GetSystemPrompt(ctx, projectID)
	require.NoError(t, err)
	assert.Empty(t, prompt)
}

func TestConfirmDiscovery_MarksComplete(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery and advance to summary
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	for i := 0; i < 4; i++ {
		discovery, err = service.AdvanceStage(ctx, discovery.ID)
		require.NoError(t, err)
	}
	assert.Equal(t, model.StageSummary, discovery.Stage)

	// Confirm
	discovery, err = service.ConfirmDiscovery(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Equal(t, model.StageComplete, discovery.Stage)
	assert.NotNil(t, discovery.ConfirmedAt)
}

func TestConfirmDiscovery_ErrorsOnNonSummaryStage(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery at welcome stage
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Try to confirm - should error because not at summary
	_, err = service.ConfirmDiscovery(ctx, discovery.ID)
	assert.Equal(t, ErrInvalidStageTransition, err)
}

func TestConfirmDiscovery_ErrorsOnNotFound(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()

	_, err := service.ConfirmDiscovery(ctx, uuid.New())
	assert.Equal(t, ErrDiscoveryNotFound, err)
}

func TestUpdateDiscoveryData(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Update data
	bc := "bakery owner"
	ps := "order tracking chaos"
	err = service.UpdateDiscoveryData(ctx, discovery.ID, &DiscoveryDataUpdate{
		BusinessContext:  &bc,
		ProblemStatement: &ps,
		Goals:            []string{"centralized orders", "reduce lost orders"},
	})
	require.NoError(t, err)

	// Verify updates
	updated, err := service.GetDiscoveryByID(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Equal(t, bc, *updated.BusinessContext)
	assert.Equal(t, ps, *updated.ProblemStatement)
	goals, _ := updated.Goals()
	assert.Equal(t, []string{"centralized orders", "reduce lost orders"}, goals)
}

func TestUpdateDiscoveryData_ErrorsOnComplete(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create and complete discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	for i := 0; i < 4; i++ {
		discovery, err = service.AdvanceStage(ctx, discovery.ID)
		require.NoError(t, err)
	}
	_, err = service.ConfirmDiscovery(ctx, discovery.ID)
	require.NoError(t, err)

	// Try to update - should error
	bc := "test"
	err = service.UpdateDiscoveryData(ctx, discovery.ID, &DiscoveryDataUpdate{
		BusinessContext: &bc,
	})
	assert.Equal(t, ErrDiscoveryAlreadyComplete, err)
}

func TestAddUser(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Add user
	permNotes := "orders only"
	user, err := service.AddUser(ctx, discovery.ID, &model.DiscoveryUser{
		Description:     "employee",
		UserCount:       2,
		HasPermissions:  false,
		PermissionNotes: &permNotes,
	})
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "employee", user.Description)
	assert.Equal(t, 2, user.UserCount)
	assert.Equal(t, discovery.ID, user.DiscoveryID)
}

func TestAddFeature(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Add MVP feature
	feature, err := service.AddFeature(ctx, discovery.ID, &model.DiscoveryFeature{
		Name:     "Order list",
		Priority: 1,
		Version:  "v1",
	})
	require.NoError(t, err)
	assert.NotNil(t, feature)
	assert.Equal(t, "Order list", feature.Name)
	assert.Equal(t, "v1", feature.Version)
	assert.Equal(t, discovery.ID, feature.DiscoveryID)
}

func TestGetSummary(t *testing.T) {
	service, repo := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Add data
	pn := "Cake Order Manager"
	ss := "Replaces paper tracking"
	err = service.UpdateDiscoveryData(ctx, discovery.ID, &DiscoveryDataUpdate{
		ProjectName:     &pn,
		SolvesStatement: &ss,
	})
	require.NoError(t, err)

	// Add users
	_, err = repo.AddUser(ctx, &model.DiscoveryUser{
		DiscoveryID:    discovery.ID,
		Description:    "owner",
		UserCount:      1,
		HasPermissions: true,
	})
	require.NoError(t, err)

	// Add features
	_, err = repo.AddFeature(ctx, &model.DiscoveryFeature{
		DiscoveryID: discovery.ID,
		Name:        "Order list",
		Priority:    1,
		Version:     "v1",
	})
	require.NoError(t, err)

	_, err = repo.AddFeature(ctx, &model.DiscoveryFeature{
		DiscoveryID: discovery.ID,
		Name:        "Calendar view",
		Priority:    1,
		Version:     "v2",
	})
	require.NoError(t, err)

	// Get summary
	summary, err := service.GetSummary(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Equal(t, "Cake Order Manager", summary.ProjectName)
	assert.Equal(t, "Replaces paper tracking", summary.SolvesStatement)
	assert.Len(t, summary.Users, 1)
	assert.Len(t, summary.MVPFeatures, 1)
	assert.Len(t, summary.FutureFeatures, 1)
}

func TestResetDiscovery(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery and advance
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)
	originalID := discovery.ID

	discovery, err = service.AdvanceStage(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Equal(t, model.StageProblem, discovery.Stage)

	// Reset
	newDiscovery, err := service.ResetDiscovery(ctx, discovery.ID)
	require.NoError(t, err)
	assert.NotEqual(t, originalID, newDiscovery.ID)
	assert.Equal(t, projectID, newDiscovery.ProjectID)
	assert.Equal(t, model.StageWelcome, newDiscovery.Stage)
}

func TestIsDiscoveryMode(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// No discovery - not in discovery mode
	isMode, err := service.IsDiscoveryMode(ctx, projectID)
	require.NoError(t, err)
	assert.False(t, isMode)

	// Create discovery - in discovery mode
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	isMode, err = service.IsDiscoveryMode(ctx, projectID)
	require.NoError(t, err)
	assert.True(t, isMode)

	// Complete discovery - not in discovery mode
	for i := 0; i < 4; i++ {
		discovery, err = service.AdvanceStage(ctx, discovery.ID)
		require.NoError(t, err)
	}
	_, err = service.ConfirmDiscovery(ctx, discovery.ID)
	require.NoError(t, err)

	isMode, err = service.IsDiscoveryMode(ctx, projectID)
	require.NoError(t, err)
	assert.False(t, isMode)
}

func TestExtractAndSaveData(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Simulate Claude response with metadata
	response := `Welcome! I'm here to help you turn your idea into a working application.

Tell me a bit about yourself - what do you do?

<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"business_context":"bakery owner"}}-->`

	// Extract and save
	err = service.ExtractAndSaveData(ctx, discovery.ID, response)
	require.NoError(t, err)

	// Verify data was saved and stage advanced
	updated, err := service.GetDiscoveryByID(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Equal(t, model.StageProblem, updated.Stage)
	assert.Equal(t, "bakery owner", *updated.BusinessContext)
}

func TestExtractAndSaveData_WithUsers(t *testing.T) {
	service, repo := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Simulate Claude response with users
	response := `Great! Let me note your users.
<!--DISCOVERY_DATA:{"stage_complete":false,"extracted":{"users":[{"description":"owner","count":1,"has_permissions":true},{"description":"employees","count":2,"has_permissions":false,"permission_notes":"orders only"}]}}-->`

	// Extract and save
	err = service.ExtractAndSaveData(ctx, discovery.ID, response)
	require.NoError(t, err)

	// Verify users were saved
	users, err := repo.GetUsers(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestExtractAndSaveData_WithFeatures(t *testing.T) {
	service, repo := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Simulate Claude response with features
	response := `Here are your MVP features!
<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"mvp_features":[{"name":"Order list","priority":1},{"name":"Order form","priority":2}],"future_features":[{"name":"Calendar view","version":"v2"}]}}-->`

	// Extract and save
	err = service.ExtractAndSaveData(ctx, discovery.ID, response)
	require.NoError(t, err)

	// Verify features were saved
	mvpFeatures, err := repo.GetMVPFeatures(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Len(t, mvpFeatures, 2)

	futureFeatures, err := repo.GetFutureFeatures(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Len(t, futureFeatures, 1)
}

func TestExtractAndSaveData_NoMetadata(t *testing.T) {
	service, _ := newTestDiscoveryService()
	ctx := context.Background()
	projectID := uuid.New()

	// Create discovery
	discovery, err := service.GetOrCreateDiscovery(ctx, projectID)
	require.NoError(t, err)

	// Response without metadata
	response := `Just a regular response without any metadata.`

	// Should not error
	err = service.ExtractAndSaveData(ctx, discovery.ID, response)
	require.NoError(t, err)

	// Stage should not have changed
	updated, err := service.GetDiscoveryByID(ctx, discovery.ID)
	require.NoError(t, err)
	assert.Equal(t, model.StageWelcome, updated.Stage)
}

func TestStripMetadata(t *testing.T) {
	input := `Hello world!
<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"business_context":"test"}}-->
More text.`

	expected := `Hello world!

More text.`

	result := StripMetadata(input)
	assert.Equal(t, expected, result)
}

func TestStripMetadata_NoMetadata(t *testing.T) {
	input := `Hello world! No metadata here.`
	result := StripMetadata(input)
	assert.Equal(t, input, result)
}
