package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// HealthHandler handles health check endpoints.
type HealthHandler struct {
	db *sqlx.DB
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db *sqlx.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status    string    `json:"status"`
	Database  string    `json:"database"`
	Timestamp time.Time `json:"timestamp"`
}

// Health returns the health status of the service.
func (h *HealthHandler) Health(c *gin.Context) {
	dbStatus := "connected"

	if h.db != nil {
		if err := h.db.Ping(); err != nil {
			dbStatus = "disconnected"
		}
	} else {
		dbStatus = "not configured"
	}

	status := "healthy"
	if dbStatus != "connected" {
		status = "degraded"
	}

	c.JSON(http.StatusOK, HealthResponse{
		Status:    status,
		Database:  dbStatus,
		Timestamp: time.Now().UTC(),
	})
}
