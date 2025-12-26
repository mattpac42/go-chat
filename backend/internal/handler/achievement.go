package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/service"
)

// AchievementHandler handles learning journey related endpoints.
type AchievementHandler struct {
	achievementSvc *service.AchievementService
	nudgeSvc       *service.NudgeService
	logger         zerolog.Logger
}

// NewAchievementHandler creates a new AchievementHandler.
func NewAchievementHandler(achievementSvc *service.AchievementService, nudgeSvc *service.NudgeService, logger zerolog.Logger) *AchievementHandler {
	return &AchievementHandler{
		achievementSvc: achievementSvc,
		nudgeSvc:       nudgeSvc,
		logger:         logger,
	}
}

// RegisterRoutes registers achievement routes on the given router group.
func (h *AchievementHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/projects/:id/progress", h.GetProgress)
	router.GET("/projects/:id/achievements", h.GetAchievements)
	router.GET("/projects/:id/achievements/unseen", h.GetUnseenAchievements)
	router.POST("/projects/:id/achievements/:achievementId/seen", h.MarkAchievementSeen)
	router.POST("/projects/:id/events", h.RecordEvent)
	router.PUT("/projects/:id/level", h.UpdateLevel)
	router.GET("/projects/:id/nudge", h.GetNextNudge)
	router.POST("/projects/:id/nudge/:type/action", h.RecordNudgeAction)
}

// GetProgress returns the learning progress for a project.
// GET /api/projects/:id/progress
func (h *AchievementHandler) GetProgress(c *gin.Context) {
	projectID, err := h.parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	progress, err := h.achievementSvc.GetProgress(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get progress")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// AchievementsResponse contains achievements with summary stats.
type AchievementsResponse struct {
	Achievements     []model.UserAchievement `json:"achievements"`
	TotalPoints      int                     `json:"totalPoints"`
	AchievementCount int                     `json:"achievementCount"`
}

// GetAchievements returns all unlocked achievements for a project.
// GET /api/projects/:id/achievements
func (h *AchievementHandler) GetAchievements(c *gin.Context) {
	projectID, err := h.parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	achievements, err := h.achievementSvc.GetAchievements(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get achievements")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get achievements"})
		return
	}

	// Calculate total points from achievements
	totalPoints := 0
	for _, a := range achievements {
		if a.Achievement != nil {
			totalPoints += a.Achievement.Points
		}
	}

	c.JSON(http.StatusOK, AchievementsResponse{
		Achievements:     achievements,
		TotalPoints:      totalPoints,
		AchievementCount: len(achievements),
	})
}

// GetUnseenAchievements returns achievements that haven't been seen (for notifications).
// GET /api/projects/:id/achievements/unseen
func (h *AchievementHandler) GetUnseenAchievements(c *gin.Context) {
	projectID, err := h.parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	achievements, err := h.achievementSvc.GetUnseenAchievements(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get unseen achievements")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get unseen achievements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"achievements": achievements,
		"count":        len(achievements),
	})
}

// MarkAchievementSeen marks an achievement as seen.
// POST /api/projects/:id/achievements/:achievementId/seen
func (h *AchievementHandler) MarkAchievementSeen(c *gin.Context) {
	_, err := h.parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	achievementIDParam := c.Param("achievementId")
	achievementID, err := uuid.Parse(achievementIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid achievement id"})
		return
	}

	if err := h.achievementSvc.MarkSeen(c.Request.Context(), achievementID); err != nil {
		h.logger.Error().Err(err).Str("achievementId", achievementID.String()).Msg("failed to mark achievement as seen")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark achievement as seen"})
		return
	}

	c.Status(http.StatusNoContent)
}

// RecordEventRequest is the request body for recording learning events.
type RecordEventRequest struct {
	Type    string                 `json:"type" binding:"required"`
	Context map[string]interface{} `json:"context"`
}

// RecordEventResponse is the response after recording an event.
type RecordEventResponse struct {
	Unlocked []UnlockedAchievement `json:"unlocked"`
	Progress ProgressSummary       `json:"progress"`
}

// UnlockedAchievement represents a newly unlocked achievement.
type UnlockedAchievement struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Points int    `json:"points"`
}

// ProgressSummary provides a summary of current progress.
type ProgressSummary struct {
	CurrentLevel int `json:"currentLevel"`
	TotalPoints  int `json:"totalPoints"`
}

// RecordEvent records a learning event which may trigger achievements.
// POST /api/projects/:id/events
func (h *AchievementHandler) RecordEvent(c *gin.Context) {
	projectID, err := h.parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req RecordEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	event := service.LearningEvent{
		Type:      req.Type,
		ProjectID: projectID,
		Context:   req.Context,
	}

	unlocked, err := h.achievementSvc.ProcessEvent(c.Request.Context(), event)
	if err != nil {
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Str("eventType", req.Type).Msg("failed to process event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process event"})
		return
	}

	// Get current progress for the response
	progress, err := h.achievementSvc.GetProgress(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get progress after event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get progress"})
		return
	}

	// Convert unlocked achievements to response format
	unlockedList := make([]UnlockedAchievement, 0, len(unlocked))
	for _, a := range unlocked {
		if a.Achievement != nil {
			unlockedList = append(unlockedList, UnlockedAchievement{
				Code:   a.Achievement.Code,
				Name:   a.Achievement.Name,
				Points: a.Achievement.Points,
			})
		}
	}

	c.JSON(http.StatusOK, RecordEventResponse{
		Unlocked: unlockedList,
		Progress: ProgressSummary{
			CurrentLevel: int(progress.CurrentLevel),
			TotalPoints:  progress.TotalPoints,
		},
	})
}

// UpdateLevelRequest is the request body for updating learning level.
type UpdateLevelRequest struct {
	Level int `json:"level" binding:"required,min=1,max=4"`
}

// UpdateLevel updates the user's learning level for a project.
// PUT /api/projects/:id/level
func (h *AchievementHandler) UpdateLevel(c *gin.Context) {
	projectID, err := h.parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req UpdateLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: level must be between 1 and 4"})
		return
	}

	newLevel := model.LearningLevel(req.Level)
	if err := h.achievementSvc.UpdateLevel(c.Request.Context(), projectID, newLevel); err != nil {
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Int("level", req.Level).Msg("failed to update level")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update level"})
		return
	}

	// Get updated progress
	progress, err := h.achievementSvc.GetProgress(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get progress after level update")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetNextNudge returns the next contextual nudge for a project.
// GET /api/projects/:id/nudge
func (h *AchievementHandler) GetNextNudge(c *gin.Context) {
	projectID, err := h.parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	nudge, err := h.nudgeSvc.GetNextNudge(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error().Err(err).Str("projectId", projectID.String()).Msg("failed to get next nudge")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get nudge"})
		return
	}

	if nudge == nil {
		c.JSON(http.StatusOK, gin.H{"nudge": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nudge": nudge})
}

// RecordNudgeActionRequest is the request body for recording nudge interactions.
type RecordNudgeActionRequest struct {
	Action string `json:"action" binding:"required,oneof=shown dismissed clicked"`
}

// RecordNudgeAction records user interaction with a nudge.
// POST /api/projects/:id/nudge/:type/action
func (h *AchievementHandler) RecordNudgeAction(c *gin.Context) {
	projectID, err := h.parseProjectID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	nudgeTypeParam := c.Param("type")
	nudgeType := model.NudgeType(nudgeTypeParam)

	// Validate nudge type
	validNudgeTypes := map[model.NudgeType]bool{
		model.NudgeExploreCode:   true,
		model.NudgeTryTreeView:   true,
		model.NudgeLevelUp:       true,
		model.NudgeViewRelations: true,
		model.NudgeExport:        true,
	}
	if !validNudgeTypes[nudgeType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nudge type"})
		return
	}

	var req RecordNudgeActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: action must be shown, dismissed, or clicked"})
		return
	}

	if err := h.nudgeSvc.RecordNudgeAction(c.Request.Context(), projectID, nudgeType, req.Action); err != nil {
		h.logger.Error().Err(err).
			Str("projectId", projectID.String()).
			Str("nudgeType", string(nudgeType)).
			Str("action", req.Action).
			Msg("failed to record nudge action")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record nudge action"})
		return
	}

	c.Status(http.StatusNoContent)
}

// parseProjectID extracts and validates the project ID from the URL.
func (h *AchievementHandler) parseProjectID(c *gin.Context) (uuid.UUID, error) {
	idParam := c.Param("id")
	return uuid.Parse(idParam)
}
