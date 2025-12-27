package handler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/service"
)

const (
	// MaxUploadSize is the maximum allowed file size for uploads (10MB).
	MaxUploadSize = 10 * 1024 * 1024

	// VisionPrompt is the prompt used for Claude Vision to describe images.
	VisionPrompt = "Describe this image in detail. If there is any text, transcribe it exactly. Format as markdown."

	// SourceMaterialsGroup is the functional group for uploaded source materials.
	SourceMaterialsGroup = "Source Materials"
)

// AllowedMimeTypes defines the allowed MIME types for image uploads.
var AllowedMimeTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/gif":  true,
	"image/webp": true,
}

// UploadHandler handles file upload endpoints.
type UploadHandler struct {
	projectRepo      repository.ProjectRepository
	fileRepo         repository.FileRepository
	fileMetadataRepo repository.FileMetadataRepository
	fileSourceRepo   repository.FileSourceRepository
	claudeVision     service.ClaudeVision
	logger           zerolog.Logger
}

// NewUploadHandler creates a new UploadHandler.
func NewUploadHandler(
	projectRepo repository.ProjectRepository,
	fileRepo repository.FileRepository,
	fileMetadataRepo repository.FileMetadataRepository,
	fileSourceRepo repository.FileSourceRepository,
	claudeVision service.ClaudeVision,
	logger zerolog.Logger,
) *UploadHandler {
	return &UploadHandler{
		projectRepo:      projectRepo,
		fileRepo:         fileRepo,
		fileMetadataRepo: fileMetadataRepo,
		fileSourceRepo:   fileSourceRepo,
		claudeVision:     claudeVision,
		logger:           logger,
	}
}

// Upload handles multipart file uploads.
// POST /api/projects/:id/upload
func (h *UploadHandler) Upload(c *gin.Context) {
	// Parse project ID
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
		h.logger.Error().Err(err).Msg("failed to get project")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get project"})
		return
	}

	// Parse multipart form with size limit
	if err := c.Request.ParseMultipartForm(MaxUploadSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large or invalid form"})
		return
	}

	// Get the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	// Validate file size
	if fileHeader.Size > MaxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file too large, maximum size is %d MB", MaxUploadSize/(1024*1024))})
		return
	}

	// Validate MIME type
	mimeType := fileHeader.Header.Get("Content-Type")
	if !AllowedMimeTypes[mimeType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("unsupported file type: %s. Allowed types: image/png, image/jpeg, image/gif, image/webp", mimeType),
		})
		return
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to open uploaded file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process upload"})
		return
	}
	defer file.Close()

	// Read file contents
	imageData, err := io.ReadAll(file)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to read uploaded file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	h.logger.Info().
		Str("projectId", projectID.String()).
		Str("filename", fileHeader.Filename).
		Str("mimeType", mimeType).
		Int64("size", fileHeader.Size).
		Msg("processing image upload")

	// Call Claude Vision to analyze the image
	markdownContent, err := h.claudeVision.AnalyzeImage(c.Request.Context(), imageData, mimeType, VisionPrompt)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to analyze image with Claude Vision")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to analyze image"})
		return
	}

	// Generate file path in sources/ folder
	originalFilename := fileHeader.Filename
	baseFilename := sanitizeUploadFilename(originalFilename)
	timestamp := time.Now().Format("2006-01-02")
	mdFilename := fmt.Sprintf("%s-%s.md", baseFilename, timestamp)
	filePath := fmt.Sprintf("sources/%s", mdFilename)

	// Save the markdown file
	savedFile, err := h.fileRepo.SaveFile(c.Request.Context(), projectID, filePath, "markdown", markdownContent)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to save markdown file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Create a short description from the first line of markdown content
	shortDesc := extractShortDescription(markdownContent)

	// Save file metadata
	_, err = h.fileMetadataRepo.Upsert(c.Request.Context(), savedFile.ID, shortDesc, "Converted from uploaded image: "+originalFilename, SourceMaterialsGroup)
	if err != nil {
		h.logger.Warn().Err(err).Msg("failed to save file metadata")
		// Continue even if metadata save fails
	}

	// Save file source record
	_, err = h.fileSourceRepo.Create(c.Request.Context(), savedFile.ID, originalFilename, mimeType, fileHeader.Size)
	if err != nil {
		h.logger.Warn().Err(err).Msg("failed to save file source record")
		// Continue even if source record save fails
	}

	h.logger.Info().
		Str("projectId", projectID.String()).
		Str("fileId", savedFile.ID.String()).
		Str("path", filePath).
		Msg("image upload processed successfully")

	// Return response
	c.JSON(http.StatusOK, model.UploadResponse{
		File: model.UploadedFile{
			ID:               savedFile.ID,
			Path:             filePath,
			Content:          markdownContent,
			ShortDescription: shortDesc,
			FunctionalGroup:  SourceMaterialsGroup,
		},
		Source: model.UploadedSource{
			OriginalFilename:  originalFilename,
			OriginalMimeType:  mimeType,
			OriginalSizeBytes: fileHeader.Size,
		},
	})
}

// sanitizeUploadFilename removes the extension and sanitizes the filename for use in paths.
func sanitizeUploadFilename(filename string) string {
	// Remove extension
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	// Replace spaces with dashes
	base = strings.ReplaceAll(base, " ", "-")

	// Remove or replace unsafe characters
	var result strings.Builder
	for _, r := range base {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			result.WriteRune(r)
		}
	}

	sanitized := result.String()

	// Limit filename length
	if len(sanitized) > 40 {
		sanitized = sanitized[:40]
	}

	// Ensure non-empty
	if sanitized == "" {
		sanitized = "upload"
	}

	return strings.ToLower(sanitized)
}

// extractShortDescription creates a short description from the first meaningful line of content.
func extractShortDescription(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and markdown headers
		if line == "" {
			continue
		}
		// Remove markdown header prefixes
		line = strings.TrimLeft(line, "#")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Truncate to reasonable length
		if len(line) > 100 {
			line = line[:97] + "..."
		}
		return line
	}
	return "Uploaded image content"
}
