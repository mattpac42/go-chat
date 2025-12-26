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

// PRDHandler handles PRD-related HTTP endpoints.
type PRDHandler struct {
	prdService *service.PRDService
	logger     zerolog.Logger
}

// NewPRDHandler creates a new PRDHandler.
func NewPRDHandler(prdService *service.PRDService, logger zerolog.Logger) *PRDHandler {
	return &PRDHandler{
		prdService: prdService,
		logger:     logger,
	}
}

// PRDListResponse represents the list of PRDs for a project.
type PRDListResponse struct {
	PRDs       []model.PRDResponse `json:"prds"`
	TotalCount int                 `json:"totalCount"`
	MVPCount   int                 `json:"mvpCount"`
}

// SetActivePRDRequest represents the request to set the active PRD.
type SetActivePRDRequest struct {
	PRDID uuid.UUID `json:"prdId" binding:"required"`
}

// ListPRDs returns all PRDs for a project.
// GET /api/projects/:id/prds
func (h *PRDHandler) ListPRDs(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	prds, err := h.prdService.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to list PRDs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list PRDs"})
		return
	}

	// Convert to response format
	var responses []model.PRDResponse
	mvpCount := 0
	for _, prd := range prds {
		response, err := prd.ToResponse()
		if err != nil {
			h.logger.Warn().Err(err).Str("prdId", prd.ID.String()).Msg("failed to convert PRD to response")
			continue
		}
		responses = append(responses, *response)
		if prd.IsMVP() {
			mvpCount++
		}
	}

	// Ensure we return an empty array, not null
	if responses == nil {
		responses = []model.PRDResponse{}
	}

	c.JSON(http.StatusOK, PRDListResponse{
		PRDs:       responses,
		TotalCount: len(responses),
		MVPCount:   mvpCount,
	})
}

// GetPRD returns a single PRD by ID.
// GET /api/prds/:id
func (h *PRDHandler) GetPRD(c *gin.Context) {
	prdID, err := parsePRDID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid PRD id"})
		return
	}

	prd, err := h.prdService.GetByID(c.Request.Context(), prdID)
	if err != nil {
		if errors.Is(err, service.ErrPRDNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "PRD not found"})
			return
		}
		h.logger.Error().Err(err).Str("prdId", prdID.String()).Msg("failed to get PRD")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get PRD"})
		return
	}

	response, err := prd.ToResponse()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to convert PRD to response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process PRD data"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdatePRDStatus updates the status of a PRD.
// PUT /api/prds/:id/status
func (h *PRDHandler) UpdatePRDStatus(c *gin.Context) {
	prdID, err := parsePRDID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid PRD id"})
		return
	}

	var req model.UpdatePRDStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validate status value
	if !model.IsValidPRDStatus(string(req.Status)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status value"})
		return
	}

	// Only allow specific status transitions via this endpoint
	// ready, in_progress, complete are the user-controllable statuses
	allowedStatuses := map[model.PRDStatus]bool{
		model.PRDStatusReady:      true,
		model.PRDStatusInProgress: true,
		model.PRDStatusComplete:   true,
	}

	if !allowedStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be one of: ready, in_progress, complete"})
		return
	}

	// Handle status-specific logic
	var updateErr error
	switch req.Status {
	case model.PRDStatusReady:
		updateErr = h.prdService.MarkAsReady(c.Request.Context(), prdID)
	case model.PRDStatusInProgress:
		updateErr = h.prdService.StartImplementation(c.Request.Context(), prdID)
	case model.PRDStatusComplete:
		updateErr = h.prdService.CompleteImplementation(c.Request.Context(), prdID)
	}

	if updateErr != nil {
		if errors.Is(updateErr, service.ErrPRDNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "PRD not found"})
			return
		}
		if errors.Is(updateErr, service.ErrInvalidStatusChange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": updateErr.Error()})
			return
		}
		h.logger.Error().Err(updateErr).Str("prdId", prdID.String()).Msg("failed to update PRD status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update PRD status"})
		return
	}

	// Return updated PRD
	prd, err := h.prdService.GetByID(c.Request.Context(), prdID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get updated PRD")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get updated PRD"})
		return
	}

	response, err := prd.ToResponse()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to convert PRD to response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process PRD data"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RetryPRDGeneration retries generation for a failed PRD.
// POST /api/prds/:id/retry
func (h *PRDHandler) RetryPRDGeneration(c *gin.Context) {
	prdID, err := parsePRDID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid PRD id"})
		return
	}

	prd, err := h.prdService.RetryGeneration(c.Request.Context(), prdID)
	if err != nil {
		if errors.Is(err, service.ErrPRDNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "PRD not found"})
			return
		}
		if errors.Is(err, service.ErrMaxRetriesExceeded) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "maximum retries exceeded"})
			return
		}
		if errors.Is(err, service.ErrInvalidStatusChange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can only retry failed or pending PRDs"})
			return
		}
		h.logger.Error().Err(err).Str("prdId", prdID.String()).Msg("failed to retry PRD generation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retry PRD generation"})
		return
	}

	response, err := prd.ToResponse()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to convert PRD to response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process PRD data"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SetActivePRD sets a PRD as the active one for a project.
// PUT /api/projects/:id/active-prd
func (h *PRDHandler) SetActivePRD(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req SetActivePRDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body, prdId is required"})
		return
	}

	err = h.prdService.SetActivePRD(c.Request.Context(), projectID, req.PRDID)
	if err != nil {
		if errors.Is(err, service.ErrPRDNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "PRD not found"})
			return
		}
		if errors.Is(err, service.ErrInvalidStatusChange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error().Err(err).
			Str("projectId", projectID.String()).
			Str("prdId", req.PRDID.String()).
			Msg("failed to set active PRD")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set active PRD"})
		return
	}

	// Return the now-active PRD
	prd, err := h.prdService.GetByID(c.Request.Context(), req.PRDID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get active PRD")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get active PRD"})
		return
	}

	response, err := prd.ToResponse()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to convert PRD to response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process PRD data"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ClearActivePRD clears the active PRD for a project.
// DELETE /api/projects/:id/active-prd
func (h *PRDHandler) ClearActivePRD(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	err = h.prdService.ClearActivePRD(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to clear active PRD")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear active PRD"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetActivePRD returns the currently active PRD for a project.
// GET /api/projects/:id/active-prd
func (h *PRDHandler) GetActivePRD(c *gin.Context) {
	projectID, err := parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	prd, err := h.prdService.GetActivePRD(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, service.ErrNoActivePRD) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no active PRD"})
			return
		}
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get active PRD")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get active PRD"})
		return
	}

	response, err := prd.ToResponse()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to convert PRD to response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process PRD data"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// parsePRDID extracts and validates the PRD ID from the URL.
func parsePRDID(c *gin.Context) (uuid.UUID, error) {
	idParam := c.Param("id")
	return uuid.Parse(idParam)
}
