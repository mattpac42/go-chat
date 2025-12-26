package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

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

	// Try to get metadata for the file
	var shortDesc, longDesc, funcGroup string
	if h.fileMetadataRepo != nil {
		metadata, err := h.fileMetadataRepo.GetByFileID(c.Request.Context(), id)
		if err == nil && metadata != nil {
			shortDesc = metadata.ShortDescription
			longDesc = metadata.LongDescription
			funcGroup = metadata.FunctionalGroup
		}
	}

	c.JSON(http.StatusOK, model.GetFileResponse{
		ID:               file.ID,
		ProjectID:        file.ProjectID,
		Path:             file.Path,
		Filename:         file.Filename,
		Language:         file.Language,
		Content:          file.Content,
		ShortDescription: shortDesc,
		LongDescription:  longDesc,
		FunctionalGroup:  funcGroup,
		CreatedAt:        file.CreatedAt,
	})
}

// DownloadProjectZip creates and returns a zip archive of all project files.
// GET /api/projects/:id/download
func (h *FileHandler) DownloadProjectZip(c *gin.Context) {
	idParam := c.Param("id")
	projectID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	// Verify project exists and get project details
	project, err := h.projectRepo.GetByID(c.Request.Context(), projectID)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get project"})
		return
	}

	// Get all files with content
	files, err := h.fileRepo.GetFilesWithContentByProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get project files"})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no files found in project"})
		return
	}

	// Create zip archive in memory
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, file := range files {
		// Create file in zip with preserved path
		writer, err := zipWriter.Create(file.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create zip entry"})
			return
		}

		// Write file content
		if _, err := writer.Write([]byte(file.Content)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write file content"})
			return
		}
	}

	// Close the zip writer to finalize the archive
	if err := zipWriter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to finalize zip archive"})
		return
	}

	// Generate filename from project title or ID
	filename := sanitizeFilename(project.Title)
	if filename == "" {
		filename = fmt.Sprintf("project-%s", projectID.String()[:8])
	}
	filename += ".zip"

	// Set response headers for file download
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))

	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}

// DownloadFile returns a single file as a download.
// GET /api/files/:id/download
func (h *FileHandler) DownloadFile(c *gin.Context) {
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

	// Determine content type from file extension
	contentType := getContentType(file.Filename)

	// Set response headers for file download
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.Filename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(file.Content)))

	c.Data(http.StatusOK, contentType, []byte(file.Content))
}

// sanitizeFilename removes or replaces characters that are unsafe for filenames.
func sanitizeFilename(name string) string {
	// Replace spaces with dashes
	name = strings.ReplaceAll(name, " ", "-")

	// Remove or replace unsafe characters
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			result.WriteRune(r)
		}
	}

	sanitized := result.String()

	// Limit filename length
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}

	return sanitized
}

// getContentType returns the MIME type for a file based on its extension.
func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return "application/octet-stream"
	}

	// Try standard mime type lookup
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		return mimeType
	}

	// Common programming file types that may not be in the standard mime database
	switch strings.ToLower(ext) {
	case ".go":
		return "text/x-go"
	case ".rs":
		return "text/x-rust"
	case ".ts":
		return "text/typescript"
	case ".tsx":
		return "text/typescript-jsx"
	case ".jsx":
		return "text/javascript-jsx"
	case ".vue":
		return "text/x-vue"
	case ".svelte":
		return "text/x-svelte"
	case ".yaml", ".yml":
		return "text/yaml"
	case ".toml":
		return "text/toml"
	case ".md":
		return "text/markdown"
	case ".sql":
		return "application/sql"
	case ".sh":
		return "application/x-sh"
	case ".dockerfile":
		return "text/x-dockerfile"
	default:
		return "application/octet-stream"
	}
}
