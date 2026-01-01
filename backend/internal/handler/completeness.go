package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/service"
)

// CompletenessHandler handles completeness check HTTP endpoints.
type CompletenessHandler struct {
	checker *service.CompletenessChecker
	logger  zerolog.Logger
}

// NewCompletenessHandler creates a new CompletenessHandler.
func NewCompletenessHandler(checker *service.CompletenessChecker, logger zerolog.Logger) *CompletenessHandler {
	return &CompletenessHandler{
		checker: checker,
		logger:  logger,
	}
}

// GetCompleteness returns the completeness report for a project.
// GET /api/projects/:id/completeness
func (h *CompletenessHandler) GetCompleteness(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	report, err := h.checker.Check(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error().Err(err).Str("projectId", projectIDStr).Msg("failed to check completeness")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check completeness"})
		return
	}

	c.JSON(http.StatusOK, report)
}
