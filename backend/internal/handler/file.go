package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

// FileHandler handles file-related endpoints.
type FileHandler struct {
	fileRepo    repository.FileRepository
	projectRepo repository.ProjectRepository
}

// NewFileHandler creates a new FileHandler.
func NewFileHandler(fileRepo repository.FileRepository, projectRepo repository.ProjectRepository) *FileHandler {
	return &FileHandler{
		fileRepo:    fileRepo,
		projectRepo: projectRepo,
	}
}

// ListFiles returns all files for a project.
// GET /api/projects/:id/files
func (h *FileHandler) ListFiles(c *gin.Context) {
	idParam := c.Param("id")
	projectID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	// Verify project exists
	_, err = h.projectRepo.GetByID(c.Request.Context(), projectID)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get project"})
		return
	}

	files, err := h.fileRepo.GetFilesByProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list files"})
		return
	}

	c.JSON(http.StatusOK, model.ListFilesResponse{
		Files: files,
	})
}

// GetFile returns a file by ID.
// GET /api/files/:id
func (h *FileHandler) GetFile(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}

	file, err := h.fileRepo.GetFile(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file"})
		return
	}

	c.JSON(http.StatusOK, model.GetFileResponse{
		ID:        file.ID,
		ProjectID: file.ProjectID,
		Path:      file.Path,
		Filename:  file.Filename,
		Language:  file.Language,
		Content:   file.Content,
		CreatedAt: file.CreatedAt,
	})
}
