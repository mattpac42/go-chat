package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/service"
)

// DiscoveryHandler handles discovery-related HTTP endpoints.
type DiscoveryHandler struct {
	service *service.DiscoveryService
	logger  zerolog.Logger
}

// NewDiscoveryHandler creates a new DiscoveryHandler.
func NewDiscoveryHandler(service *service.DiscoveryService, logger zerolog.Logger) *DiscoveryHandler {
	return &DiscoveryHandler{
		service: service,
		logger:  logger,
	}
}

// DiscoveryWithSummaryResponse is the response for GetDiscovery when stage >= summary.
type DiscoveryWithSummaryResponse struct {
	Discovery *model.DiscoveryResponse `json:"discovery"`
	Summary   *model.DiscoverySummary  `json:"summary,omitempty"`
}

// GetDiscovery returns the current discovery state for a project.
// GET /api/projects/:id/discovery
func (h *DiscoveryHandler) GetDiscovery(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	discovery, err := h.service.GetOrCreateDiscovery(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get discovery")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discovery"})
		return
	}

	response, err := discovery.ToResponse()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to convert discovery to response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process discovery data"})
		return
	}

	// Include summary if stage is summary or complete
	if discovery.Stage == model.StageSummary || discovery.Stage == model.StageComplete {
		summary, err := h.service.GetSummary(c.Request.Context(), discovery.ID)
		if err != nil {
			h.logger.Warn().Err(err).Msg("failed to get discovery summary")
			// Continue without summary
			c.JSON(http.StatusOK, DiscoveryWithSummaryResponse{
				Discovery: response,
			})
			return
		}

		c.JSON(http.StatusOK, DiscoveryWithSummaryResponse{
			Discovery: response,
			Summary:   summary,
		})
		return
	}

	c.JSON(http.StatusOK, DiscoveryWithSummaryResponse{
		Discovery: response,
	})
}

// AdvanceStage moves the discovery to the next stage.
// PUT /api/projects/:id/discovery/stage
func (h *DiscoveryHandler) AdvanceStage(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	// Get current discovery
	discovery, err := h.service.GetDiscovery(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, service.ErrDiscoveryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "discovery not found"})
			return
		}
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get discovery")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discovery"})
		return
	}

	// Advance to next stage
	updated, err := h.service.AdvanceStage(c.Request.Context(), discovery.ID)
	if err != nil {
		if errors.Is(err, service.ErrDiscoveryAlreadyComplete) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "discovery is already complete"})
			return
		}
		if errors.Is(err, service.ErrInvalidStageTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stage transition"})
			return
		}
		h.logger.Error().Err(err).Msg("failed to advance stage")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to advance stage"})
		return
	}

	response, err := updated.ToResponse()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to convert discovery to response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process discovery data"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateData updates the discovery data fields.
// PUT /api/projects/:id/discovery/data
func (h *DiscoveryHandler) UpdateData(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req model.UpdateDiscoveryDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Get current discovery
	discovery, err := h.service.GetDiscovery(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, service.ErrDiscoveryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "discovery not found"})
			return
		}
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get discovery")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discovery"})
		return
	}

	// Update data
	update := &service.DiscoveryDataUpdate{
		BusinessContext:  req.BusinessContext,
		ProblemStatement: req.ProblemStatement,
		Goals:            req.Goals,
		ProjectName:      req.ProjectName,
		SolvesStatement:  req.SolvesStatement,
	}

	if err := h.service.UpdateDiscoveryData(c.Request.Context(), discovery.ID, update); err != nil {
		if errors.Is(err, service.ErrDiscoveryAlreadyComplete) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "discovery is already complete"})
			return
		}
		h.logger.Error().Err(err).Msg("failed to update discovery data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update discovery data"})
		return
	}

	// Get updated discovery
	updated, err := h.service.GetDiscovery(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get updated discovery")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get updated discovery"})
		return
	}

	response, err := updated.ToResponse()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to convert discovery to response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process discovery data"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// AddUser adds a user persona to the discovery.
// POST /api/projects/:id/discovery/users
func (h *DiscoveryHandler) AddUser(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req model.AddDiscoveryUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "description is required"})
		return
	}

	// Get current discovery
	discovery, err := h.service.GetDiscovery(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, service.ErrDiscoveryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "discovery not found"})
			return
		}
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get discovery")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discovery"})
		return
	}

	// Create user
	user := &model.DiscoveryUser{
		Description:     req.Description,
		UserCount:       req.UserCount,
		HasPermissions:  req.HasPermissions,
		PermissionNotes: req.PermissionNotes,
	}

	created, err := h.service.AddUser(c.Request.Context(), discovery.ID, user)
	if err != nil {
		if errors.Is(err, service.ErrDiscoveryAlreadyComplete) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "discovery is already complete"})
			return
		}
		h.logger.Error().Err(err).Msg("failed to add user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add user"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// AddFeature adds a feature to the discovery.
// POST /api/projects/:id/discovery/features
func (h *DiscoveryHandler) AddFeature(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req model.AddDiscoveryFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	// Get current discovery
	discovery, err := h.service.GetDiscovery(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, service.ErrDiscoveryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "discovery not found"})
			return
		}
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get discovery")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discovery"})
		return
	}

	// Create feature
	version := req.Version
	if version == "" {
		version = "v1"
	}

	feature := &model.DiscoveryFeature{
		Name:     req.Name,
		Priority: req.Priority,
		Version:  version,
	}

	created, err := h.service.AddFeature(c.Request.Context(), discovery.ID, feature)
	if err != nil {
		if errors.Is(err, service.ErrDiscoveryAlreadyComplete) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "discovery is already complete"})
			return
		}
		h.logger.Error().Err(err).Msg("failed to add feature")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add feature"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// ConfirmDiscovery marks the discovery as complete.
// POST /api/projects/:id/discovery/confirm
func (h *DiscoveryHandler) ConfirmDiscovery(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	// Get current discovery
	discovery, err := h.service.GetDiscovery(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, service.ErrDiscoveryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "discovery not found"})
			return
		}
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get discovery")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discovery"})
		return
	}

	// Confirm discovery
	updated, err := h.service.ConfirmDiscovery(c.Request.Context(), discovery.ID)
	if err != nil {
		if errors.Is(err, service.ErrDiscoveryAlreadyComplete) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "discovery is already complete"})
			return
		}
		if errors.Is(err, service.ErrInvalidStageTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "must be in summary stage to confirm"})
			return
		}
		h.logger.Error().Err(err).Msg("failed to confirm discovery")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm discovery"})
		return
	}

	response, err := updated.ToResponse()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to convert discovery to response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process discovery data"})
		return
	}

	// Include summary in confirmation response
	summary, err := h.service.GetSummary(c.Request.Context(), discovery.ID)
	if err != nil {
		h.logger.Warn().Err(err).Msg("failed to get discovery summary")
		c.JSON(http.StatusOK, DiscoveryWithSummaryResponse{
			Discovery: response,
		})
		return
	}

	c.JSON(http.StatusOK, DiscoveryWithSummaryResponse{
		Discovery: response,
		Summary:   summary,
	})
}

// ResetDiscovery deletes the current discovery and starts over.
// DELETE /api/projects/:id/discovery
func (h *DiscoveryHandler) ResetDiscovery(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	// Get current discovery
	discovery, err := h.service.GetDiscovery(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, service.ErrDiscoveryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "discovery not found"})
			return
		}
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get discovery")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discovery"})
		return
	}

	// Reset discovery
	newDiscovery, err := h.service.ResetDiscovery(c.Request.Context(), discovery.ID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to reset discovery")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset discovery"})
		return
	}

	response, err := newDiscovery.ToResponse()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to convert discovery to response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process discovery data"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// parseProjectID extracts and validates the project ID from the URL.
func parseProjectID(c *gin.Context) (uuid.UUID, error) {
	idParam := c.Param("id")
	return uuid.Parse(idParam)
}
