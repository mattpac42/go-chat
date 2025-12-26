package handler

import (
	"bytes"
	"context"
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

func setupPRDTestRouter(prdService *service.PRDService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger := zerolog.Nop()
	handler := NewPRDHandler(prdService, logger)

	projects := router.Group("/api/projects")
	{
		projects.GET("/:id/prds", handler.ListPRDs)
		projects.GET("/:id/active-prd", handler.GetActivePRD)
		projects.PUT("/:id/active-prd", handler.SetActivePRD)
		projects.DELETE("/:id/active-prd", handler.ClearActivePRD)
	}

	prds := router.Group("/api/prds")
	{
		prds.GET("/:id", handler.GetPRD)
		prds.PUT("/:id/status", handler.UpdatePRDStatus)
		prds.POST("/:id/retry", handler.RetryPRDGeneration)
	}

	return router
}

// createTestPRD creates a PRD directly in the mock repository for testing.
func createTestPRD(t *testing.T, prdRepo *repository.MockPRDRepository, projectID, discoveryID, featureID uuid.UUID, status model.PRDStatus, version string) *model.PRD {
	prd := &model.PRD{
		ProjectID:   projectID,
		DiscoveryID: discoveryID,
		FeatureID:   featureID,
		Title:       "Test Feature",
		Overview:    "Test overview",
		Version:     version,
		Priority:    1,
		Status:      status,
	}
	prd.SetUserStories([]model.UserStory{})
	prd.SetAcceptanceCriteria([]model.AcceptanceCriterion{})
	prd.SetTechnicalNotes([]model.TechnicalNote{})

	created, err := prdRepo.Create(context.Background(), prd)
	require.NoError(t, err)
	return created
}

func TestListPRDs_Empty(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()

	req, _ := http.NewRequest("GET", "/api/projects/"+projectID.String()+"/prds", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PRDListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 0, response.TotalCount)
	assert.Equal(t, 0, response.MVPCount)
	assert.NotNil(t, response.PRDs)
	assert.Len(t, response.PRDs, 0)
}

func TestListPRDs_WithPRDs(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()

	// Create test PRDs
	createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusDraft, "v1")
	createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusReady, "v1")
	createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusPending, "v2")

	req, _ := http.NewRequest("GET", "/api/projects/"+projectID.String()+"/prds", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PRDListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 2, response.MVPCount) // Two v1 PRDs
	assert.Len(t, response.PRDs, 3)
}

func TestListPRDs_InvalidProjectID(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	req, _ := http.NewRequest("GET", "/api/projects/invalid-uuid/prds", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPRD_Success(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()
	prd := createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusDraft, "v1")

	req, _ := http.NewRequest("GET", "/api/prds/"+prd.ID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.PRDResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, prd.ID, response.ID)
	assert.Equal(t, "Test Feature", response.Title)
	assert.Equal(t, model.PRDStatusDraft, response.Status)
	assert.True(t, response.IsMVP)
}

func TestGetPRD_NotFound(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	prdID := uuid.New()

	req, _ := http.NewRequest("GET", "/api/prds/"+prdID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPRD_InvalidID(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	req, _ := http.NewRequest("GET", "/api/prds/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePRDStatus_ToReady(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()
	prd := createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusDraft, "v1")

	requestBody := model.UpdatePRDStatusRequest{Status: model.PRDStatusReady}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/prds/"+prd.ID.String()+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.PRDResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, model.PRDStatusReady, response.Status)
	assert.NotNil(t, response.ApprovedAt)
}

func TestUpdatePRDStatus_ToInProgress(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()
	prd := createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusReady, "v1")

	requestBody := model.UpdatePRDStatusRequest{Status: model.PRDStatusInProgress}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/prds/"+prd.ID.String()+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.PRDResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, model.PRDStatusInProgress, response.Status)
	assert.NotNil(t, response.StartedAt)
}

func TestUpdatePRDStatus_ToComplete(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()
	prd := createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusInProgress, "v1")

	requestBody := model.UpdatePRDStatusRequest{Status: model.PRDStatusComplete}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/prds/"+prd.ID.String()+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.PRDResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, model.PRDStatusComplete, response.Status)
	assert.NotNil(t, response.CompletedAt)
}

func TestUpdatePRDStatus_InvalidTransition(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()
	// Create a PRD in pending status - cannot go directly to ready
	prd := createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusPending, "v1")

	requestBody := model.UpdatePRDStatusRequest{Status: model.PRDStatusReady}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/prds/"+prd.ID.String()+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePRDStatus_InvalidStatus(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()
	prd := createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusDraft, "v1")

	// Try to set to pending (not allowed via API)
	requestBody := model.UpdatePRDStatusRequest{Status: model.PRDStatusPending}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/prds/"+prd.ID.String()+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "status must be one of: ready, in_progress, complete", errResp["error"])
}

func TestUpdatePRDStatus_NotFound(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	prdID := uuid.New()

	requestBody := model.UpdatePRDStatusRequest{Status: model.PRDStatusReady}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/prds/"+prdID.String()+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSetActivePRD_Success(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()
	prd := createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusReady, "v1")

	requestBody := SetActivePRDRequest{PRDID: prd.ID}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/projects/"+projectID.String()+"/active-prd", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.PRDResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, model.PRDStatusInProgress, response.Status)
	assert.NotNil(t, response.StartedAt)
}

func TestSetActivePRD_PRDNotReady(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()
	// Create a PRD in draft status - cannot be made active
	prd := createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusDraft, "v1")

	requestBody := SetActivePRDRequest{PRDID: prd.ID}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/projects/"+projectID.String()+"/active-prd", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSetActivePRD_PRDNotFound(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()

	requestBody := SetActivePRDRequest{PRDID: uuid.New()}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/projects/"+projectID.String()+"/active-prd", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetActivePRD_Success(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()
	prd := createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusInProgress, "v1")

	req, _ := http.NewRequest("GET", "/api/projects/"+projectID.String()+"/active-prd", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.PRDResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, prd.ID, response.ID)
	assert.Equal(t, model.PRDStatusInProgress, response.Status)
}

func TestGetActivePRD_NoActive(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()

	req, _ := http.NewRequest("GET", "/api/projects/"+projectID.String()+"/active-prd", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "no active PRD", errResp["error"])
}

func TestClearActivePRD_Success(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()
	createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusInProgress, "v1")

	req, _ := http.NewRequest("DELETE", "/api/projects/"+projectID.String()+"/active-prd", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestClearActivePRD_NoActive(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()

	// Should succeed even if no active PRD
	req, _ := http.NewRequest("DELETE", "/api/projects/"+projectID.String()+"/active-prd", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestRetryPRDGeneration_NotFailedPRD(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	projectID := uuid.New()
	discoveryID := uuid.New()
	// Create a PRD in draft status - cannot retry
	prd := createTestPRD(t, prdRepo, projectID, discoveryID, uuid.New(), model.PRDStatusDraft, "v1")

	req, _ := http.NewRequest("POST", "/api/prds/"+prd.ID.String()+"/retry", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "can only retry failed or pending PRDs", errResp["error"])
}

func TestRetryPRDGeneration_NotFound(t *testing.T) {
	prdRepo := repository.NewMockPRDRepository()
	discoveryRepo := repository.NewMockDiscoveryRepository()
	logger := zerolog.Nop()
	prdService := service.NewPRDService(prdRepo, discoveryRepo, nil, logger)
	router := setupPRDTestRouter(prdService)

	prdID := uuid.New()

	req, _ := http.NewRequest("POST", "/api/prds/"+prdID.String()+"/retry", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
