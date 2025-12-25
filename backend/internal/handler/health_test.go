package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealthHandler_Health(t *testing.T) {
	t.Run("returns healthy status when no database", func(t *testing.T) {
		// Arrange
		handler := NewHealthHandler(nil)
		router := gin.New()
		router.GET("/health", handler.Health)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "degraded", response.Status)
		assert.Equal(t, "not configured", response.Database)
		assert.NotZero(t, response.Timestamp)
	})
}
