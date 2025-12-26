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
	fileRepo         repository.FileRepository
	projectRepo      repository.ProjectRepository
	fileMetadataRepo repository.FileMetadataRepository
}

// NewFileHandler creates a new FileHandler.
func NewFileHandler(fileRepo repository.FileRepository, projectRepo repository.ProjectRepository, fileMetadataRepo repository.FileMetadataRepository) *FileHandler {
	return &FileHandler{
		fileRepo:         fileRepo,
		projectRepo:      projectRepo,
		fileMetadataRepo: fileMetadataRepo,
	}
}

// ListFiles returns all files for a project with metadata.
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

	// Try to get files with metadata first
	if h.fileMetadataRepo != nil {
		filesWithMetadata, err := h.fileMetadataRepo.GetFilesWithMetadata(c.Request.Context(), projectID)
		if err == nil {
			// Convert to response format
			result := make([]model.FileListItemWithMetadata, len(filesWithMetadata))
			for i, f := range filesWithMetadata {
				result[i] = model.FileListItemWithMetadata{
					ID:               f.ID,
					Path:             f.Path,
					Filename:         f.Filename,
					Language:         f.Language,
					ShortDescription: f.ShortDescription,
					LongDescription:  f.LongDescription,
					FunctionalGroup:  f.FunctionalGroup,
					CreatedAt:        f.CreatedAt,
				}
			}
			c.JSON(http.StatusOK, model.ListFilesWithMetadataResponse{
				Files: result,
			})
			return
		}
		// Fall through to basic files if metadata query fails
	}

	// Fallback to basic file list (without metadata)
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
