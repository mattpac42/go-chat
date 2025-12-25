package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

// ProjectHandler handles project-related endpoints.
type ProjectHandler struct {
	repo repository.ProjectRepository
}

// NewProjectHandler creates a new ProjectHandler.
func NewProjectHandler(repo repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{repo: repo}
}

// List returns all projects.
func (h *ProjectHandler) List(c *gin.Context) {
	projects, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list projects"})
		return
	}

	c.JSON(http.StatusOK, model.ListProjectsResponse{
		Projects: projects,
	})
}

// Create creates a new project.
func (h *ProjectHandler) Create(c *gin.Context) {
	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If no body provided, use defaults
		req = model.CreateProjectRequest{}
	}

	title := req.Title
	if title == "" {
		title = "New Project"
	}

	project, err := h.repo.Create(c.Request.Context(), title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, model.CreateProjectResponse{
		ID:        project.ID,
		Title:     project.Title,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	})
}

// Get returns a project with its messages.
func (h *ProjectHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	project, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get project"})
		return
	}

	messages, err := h.repo.GetMessages(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, model.GetProjectResponse{
		ID:        project.ID,
		Title:     project.Title,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
		Messages:  messages,
	})
}

// Delete removes a project and all its messages.
func (h *ProjectHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	err = h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete project"})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateTimestamp updates the project's updated_at timestamp.
func (h *ProjectHandler) UpdateTimestamp(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	err = h.repo.UpdateTimestamp(c.Request.Context(), id, time.Now().UTC())
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update project"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Update updates a project's title (partial update).
func (h *ProjectHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req model.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	project, err := h.repo.UpdateTitle(c.Request.Context(), id, req.Title)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update project"})
		return
	}

	c.JSON(http.StatusOK, model.UpdateProjectResponse{
		ID:        project.ID,
		Title:     project.Title,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	})
}
