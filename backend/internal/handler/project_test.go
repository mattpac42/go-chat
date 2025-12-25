package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

func TestProjectHandler_Create(t *testing.T) {
	t.Run("creates project with default title", func(t *testing.T) {
		// Arrange
		repo := repository.NewMockProjectRepository()
		handler := NewProjectHandler(repo)
		router := gin.New()
		router.POST("/api/projects", handler.Create)

		req := httptest.NewRequest(http.MethodPost, "/api/projects", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		var response model.CreateProjectResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "New Project", response.Title)
		assert.NotEmpty(t, response.ID)
		assert.NotZero(t, response.CreatedAt)
		assert.NotZero(t, response.UpdatedAt)
	})

	t.Run("creates project with custom title", func(t *testing.T) {
		// Arrange
		repo := repository.NewMockProjectRepository()
		handler := NewProjectHandler(repo)
		router := gin.New()
		router.POST("/api/projects", handler.Create)

		body := bytes.NewBufferString(`{"title": "My Custom Project"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/projects", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		var response model.CreateProjectResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "My Custom Project", response.Title)
	})
}

func TestProjectHandler_List(t *testing.T) {
	t.Run("returns empty list when no projects", func(t *testing.T) {
		// Arrange
		repo := repository.NewMockProjectRepository()
		handler := NewProjectHandler(repo)
		router := gin.New()
		router.GET("/api/projects", handler.List)

		req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response model.ListProjectsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Empty(t, response.Projects)
	})

	t.Run("returns list of projects", func(t *testing.T) {
		// Arrange
		repo := repository.NewMockProjectRepository()
		_, _ = repo.Create(nil, "Project 1")
		_, _ = repo.Create(nil, "Project 2")

		handler := NewProjectHandler(repo)
		router := gin.New()
		router.GET("/api/projects", handler.List)

		req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response model.ListProjectsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Projects, 2)
	})
}

func TestProjectHandler_Get(t *testing.T) {
	t.Run("returns project with messages", func(t *testing.T) {
		// Arrange
		repo := repository.NewMockProjectRepository()
		project, _ := repo.Create(nil, "Test Project")
		_, _ = repo.CreateMessage(nil, project.ID, model.RoleUser, "Hello")
		_, _ = repo.CreateMessage(nil, project.ID, model.RoleAssistant, "Hi there!")

		handler := NewProjectHandler(repo)
		router := gin.New()
		router.GET("/api/projects/:id", handler.Get)

		req := httptest.NewRequest(http.MethodGet, "/api/projects/"+project.ID.String(), nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response model.GetProjectResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, project.ID, response.ID)
		assert.Equal(t, "Test Project", response.Title)
		assert.Len(t, response.Messages, 2)
	})

	t.Run("returns 404 for non-existent project", func(t *testing.T) {
		// Arrange
		repo := repository.NewMockProjectRepository()
		handler := NewProjectHandler(repo)
		router := gin.New()
		router.GET("/api/projects/:id", handler.Get)

		req := httptest.NewRequest(http.MethodGet, "/api/projects/00000000-0000-0000-0000-000000000000", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 400 for invalid uuid", func(t *testing.T) {
		// Arrange
		repo := repository.NewMockProjectRepository()
		handler := NewProjectHandler(repo)
		router := gin.New()
		router.GET("/api/projects/:id", handler.Get)

		req := httptest.NewRequest(http.MethodGet, "/api/projects/invalid-uuid", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProjectHandler_Delete(t *testing.T) {
	t.Run("deletes existing project", func(t *testing.T) {
		// Arrange
		repo := repository.NewMockProjectRepository()
		project, _ := repo.Create(nil, "Project to Delete")

		handler := NewProjectHandler(repo)
		router := gin.New()
		router.DELETE("/api/projects/:id", handler.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/api/projects/"+project.ID.String(), nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify project is deleted
		_, err := repo.GetByID(nil, project.ID)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("returns 404 for non-existent project", func(t *testing.T) {
		// Arrange
		repo := repository.NewMockProjectRepository()
		handler := NewProjectHandler(repo)
		router := gin.New()
		router.DELETE("/api/projects/:id", handler.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/api/projects/00000000-0000-0000-0000-000000000000", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
