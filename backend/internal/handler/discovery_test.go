package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/service"
)

func setupDiscoveryTestRouter(discoveryService *service.DiscoveryService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger := zerolog.Nop()
	handler := NewDiscoveryHandler(discoveryService, logger)

	projects := router.Group("/api/projects")
	{
		projects.GET("/:id/discovery", handler.GetDiscovery)
		projects.PUT("/:id/discovery/stage", handler.AdvanceStage)
		projects.PUT("/:id/discovery/data", handler.UpdateData)
		projects.POST("/:id/discovery/users", handler.AddUser)
		projects.POST("/:id/discovery/features", handler.AddFeature)
		projects.POST("/:id/discovery/confirm", handler.ConfirmDiscovery)
		projects.DELETE("/:id/discovery", handler.ResetDiscovery)
	}

	return router
}

func TestGetDiscovery_CreatesNew(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()

	req, _ := http.NewRequest("GET", "/api/projects/"+projectID.String()+"/discovery", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response DiscoveryWithSummaryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, projectID, response.Discovery.ProjectID)
	assert.Equal(t, model.StageWelcome, response.Discovery.Stage)
	assert.Equal(t, 1, response.Discovery.StageNumber)
	assert.Nil(t, response.Summary) // Summary should not be included for welcome stage
}

func TestGetDiscovery_InvalidProjectID(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	req, _ := http.NewRequest("GET", "/api/projects/invalid-uuid/discovery", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDiscovery_WithSummary(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	// Create a discovery and advance it to summary stage
	projectID := uuid.New()
	discovery, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	// Advance to summary stage
	discovery, err = mockRepo.UpdateStage(nil, discovery.ID, model.StageSummary)
	require.NoError(t, err)

	// Set project name and solves statement
	projectName := "Test Project"
	solvesStatement := "Test Solves Statement"
	discovery.ProjectName = &projectName
	discovery.SolvesStatement = &solvesStatement
	_, err = mockRepo.Update(nil, discovery)
	require.NoError(t, err)

	// Add a user
	_, err = mockRepo.AddUser(nil, &model.DiscoveryUser{
		DiscoveryID: discovery.ID,
		Description: "Test User",
		UserCount:   1,
	})
	require.NoError(t, err)

	// Add a feature
	_, err = mockRepo.AddFeature(nil, &model.DiscoveryFeature{
		DiscoveryID: discovery.ID,
		Name:        "Test Feature",
		Priority:    1,
		Version:     "v1",
	})
	require.NoError(t, err)

	req, _ := http.NewRequest("GET", "/api/projects/"+projectID.String()+"/discovery", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response DiscoveryWithSummaryResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, model.StageSummary, response.Discovery.Stage)
	assert.NotNil(t, response.Summary)
	assert.Equal(t, "Test Project", response.Summary.ProjectName)
	assert.Equal(t, "Test Solves Statement", response.Summary.SolvesStatement)
	assert.Len(t, response.Summary.Users, 1)
	assert.Len(t, response.Summary.MVPFeatures, 1)
}

func TestAdvanceStage_Success(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	_, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	req, _ := http.NewRequest("PUT", "/api/projects/"+projectID.String()+"/discovery/stage", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.DiscoveryResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, model.StageProblem, response.Stage) // Moved from welcome to problem
	assert.Equal(t, 2, response.StageNumber)
}

func TestAdvanceStage_NotFound(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()

	req, _ := http.NewRequest("PUT", "/api/projects/"+projectID.String()+"/discovery/stage", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdvanceStage_AlreadyComplete(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	discovery, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	// Mark as complete
	_, err = mockRepo.MarkComplete(nil, discovery.ID)
	require.NoError(t, err)

	req, _ := http.NewRequest("PUT", "/api/projects/"+projectID.String()+"/discovery/stage", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "discovery is already complete", errResp["error"])
}

func TestUpdateData_Success(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	_, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	businessContext := "I run a bakery"
	problemStatement := "Order tracking is chaos"
	requestBody := model.UpdateDiscoveryDataRequest{
		BusinessContext:  &businessContext,
		ProblemStatement: &problemStatement,
		Goals:            []string{"centralized order tracking", "reduce lost orders"},
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/api/projects/"+projectID.String()+"/discovery/data", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.DiscoveryResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotNil(t, response.BusinessContext)
	assert.Equal(t, "I run a bakery", *response.BusinessContext)
	assert.NotNil(t, response.ProblemStatement)
	assert.Equal(t, "Order tracking is chaos", *response.ProblemStatement)
	assert.Len(t, response.Goals, 2)
}

func TestUpdateData_AlreadyComplete(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	discovery, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	// Mark as complete
	_, err = mockRepo.MarkComplete(nil, discovery.ID)
	require.NoError(t, err)

	businessContext := "I run a bakery"
	requestBody := model.UpdateDiscoveryDataRequest{
		BusinessContext: &businessContext,
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/api/projects/"+projectID.String()+"/discovery/data", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddUser_Success(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	_, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	permissionNotes := "employees: orders only"
	requestBody := model.AddDiscoveryUserRequest{
		Description:     "order takers",
		UserCount:       2,
		HasPermissions:  false,
		PermissionNotes: &permissionNotes,
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/projects/"+projectID.String()+"/discovery/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.DiscoveryUser
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "order takers", response.Description)
	assert.Equal(t, 2, response.UserCount)
	assert.False(t, response.HasPermissions)
	assert.NotNil(t, response.PermissionNotes)
	assert.Equal(t, "employees: orders only", *response.PermissionNotes)
}

func TestAddUser_MissingDescription(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	_, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	requestBody := model.AddDiscoveryUserRequest{
		Description: "", // Missing description
		UserCount:   2,
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/projects/"+projectID.String()+"/discovery/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "description is required", errResp["error"])
}

func TestAddFeature_Success(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	_, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	requestBody := model.AddDiscoveryFeatureRequest{
		Name:     "Order list",
		Priority: 1,
		Version:  "v1",
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/projects/"+projectID.String()+"/discovery/features", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.DiscoveryFeature
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Order list", response.Name)
	assert.Equal(t, 1, response.Priority)
	assert.Equal(t, "v1", response.Version)
}

func TestAddFeature_DefaultVersion(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	_, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	requestBody := model.AddDiscoveryFeatureRequest{
		Name:     "Order list",
		Priority: 1,
		// Version not specified
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/projects/"+projectID.String()+"/discovery/features", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.DiscoveryFeature
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "v1", response.Version) // Default to v1
}

func TestAddFeature_MissingName(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	_, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	requestBody := model.AddDiscoveryFeatureRequest{
		Name:     "", // Missing name
		Priority: 1,
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/projects/"+projectID.String()+"/discovery/features", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "name is required", errResp["error"])
}

func TestConfirmDiscovery_Success(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	discovery, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	// Advance to summary stage
	_, err = mockRepo.UpdateStage(nil, discovery.ID, model.StageSummary)
	require.NoError(t, err)

	req, _ := http.NewRequest("POST", "/api/projects/"+projectID.String()+"/discovery/confirm", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response DiscoveryWithSummaryResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, model.StageComplete, response.Discovery.Stage)
	assert.NotNil(t, response.Discovery.ConfirmedAt)
}

func TestConfirmDiscovery_NotInSummaryStage(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	_, err := mockRepo.Create(nil, projectID) // Still in welcome stage
	require.NoError(t, err)

	req, _ := http.NewRequest("POST", "/api/projects/"+projectID.String()+"/discovery/confirm", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "must be in summary stage to confirm", errResp["error"])
}

func TestResetDiscovery_Success(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()
	discovery, err := mockRepo.Create(nil, projectID)
	require.NoError(t, err)

	// Advance the discovery
	_, err = mockRepo.UpdateStage(nil, discovery.ID, model.StageProblem)
	require.NoError(t, err)

	// Add some data
	businessContext := "Test business"
	discovery.BusinessContext = &businessContext
	_, err = mockRepo.Update(nil, discovery)
	require.NoError(t, err)

	req, _ := http.NewRequest("DELETE", "/api/projects/"+projectID.String()+"/discovery", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.DiscoveryResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should be a fresh discovery
	assert.Equal(t, model.StageWelcome, response.Stage)
	assert.Nil(t, response.BusinessContext) // Data should be cleared
}

func TestResetDiscovery_NotFound(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()

	req, _ := http.NewRequest("DELETE", "/api/projects/"+projectID.String()+"/discovery", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFullDiscoveryFlow(t *testing.T) {
	mockRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	discoveryService := service.NewDiscoveryService(mockRepo, logger)
	router := setupDiscoveryTestRouter(discoveryService)

	projectID := uuid.New()

	// 1. Get discovery (creates new one at welcome stage)
	req, _ := http.NewRequest("GET", "/api/projects/"+projectID.String()+"/discovery", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var discovery DiscoveryWithSummaryResponse
	json.Unmarshal(w.Body.Bytes(), &discovery)
	assert.Equal(t, model.StageWelcome, discovery.Discovery.Stage)

	// 2. Update data with business context
	businessContext := "I run a bakery"
	updateReq := model.UpdateDiscoveryDataRequest{BusinessContext: &businessContext}
	body, _ := json.Marshal(updateReq)
	req, _ = http.NewRequest("PUT", "/api/projects/"+projectID.String()+"/discovery/data", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. Advance through stages
	stages := []model.DiscoveryStage{model.StageProblem, model.StagePersonas, model.StageMVP, model.StageSummary}
	for _, expectedStage := range stages {
		req, _ = http.NewRequest("PUT", "/api/projects/"+projectID.String()+"/discovery/stage", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var stageResp model.DiscoveryResponse
		json.Unmarshal(w.Body.Bytes(), &stageResp)
		assert.Equal(t, expectedStage, stageResp.Stage)
	}

	// 4. Add a user persona
	userReq := model.AddDiscoveryUserRequest{Description: "owner/baker", UserCount: 1, HasPermissions: true}
	body, _ = json.Marshal(userReq)
	req, _ = http.NewRequest("POST", "/api/projects/"+projectID.String()+"/discovery/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 5. Add MVP features
	featureReq := model.AddDiscoveryFeatureRequest{Name: "Order list", Priority: 1, Version: "v1"}
	body, _ = json.Marshal(featureReq)
	req, _ = http.NewRequest("POST", "/api/projects/"+projectID.String()+"/discovery/features", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 6. Add future feature
	featureReq = model.AddDiscoveryFeatureRequest{Name: "Calendar view", Priority: 1, Version: "v2"}
	body, _ = json.Marshal(featureReq)
	req, _ = http.NewRequest("POST", "/api/projects/"+projectID.String()+"/discovery/features", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 7. Confirm discovery
	req, _ = http.NewRequest("POST", "/api/projects/"+projectID.String()+"/discovery/confirm", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var confirmResp DiscoveryWithSummaryResponse
	json.Unmarshal(w.Body.Bytes(), &confirmResp)
	assert.Equal(t, model.StageComplete, confirmResp.Discovery.Stage)
	assert.NotNil(t, confirmResp.Discovery.ConfirmedAt)
	assert.NotNil(t, confirmResp.Summary)
	assert.Len(t, confirmResp.Summary.Users, 1)
	assert.Len(t, confirmResp.Summary.MVPFeatures, 1)
	assert.Len(t, confirmResp.Summary.FutureFeatures, 1)
}
