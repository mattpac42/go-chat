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
	// The response format includes a smart filename on the first line.
	VisionPrompt = `First, provide a short descriptive filename (1-3 words, kebab-case) for this image on the FIRST line, prefixed with "FILENAME:".

Then describe this image in detail. If there is any text, transcribe it exactly. Format as markdown.

Example response format:
FILENAME: login-screen-mockup

## Login Screen Design
...`

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
	visionResponse, err := h.claudeVision.AnalyzeImage(c.Request.Context(), imageData, mimeType, VisionPrompt)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to analyze image with Claude Vision")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to analyze image"})
		return
	}

	// Parse vision response to extract smart filename and content
	smartFilename, markdownContent := parseVisionResponse(visionResponse)

	// Generate file path in sources/ folder using smart filename
	originalFilename := fileHeader.Filename
	timestamp := time.Now().Format("2006-01-02")
	mdFilename := fmt.Sprintf("%s-%s.md", smartFilename, timestamp)
	filePath := fmt.Sprintf("sources/%s", mdFilename)

	h.logger.Debug().
		Str("originalFilename", originalFilename).
		Str("smartFilename", smartFilename).
		Str("filePath", filePath).
		Msg("generated smart filename from vision response")

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

// parseVisionResponse extracts a smart filename and content from Claude Vision's response.
// It expects the response to have "FILENAME: <name>" on the first line, followed by the content.
// Returns the sanitized filename and the remaining content.
// Falls back to "image-upload" if parsing fails.
func parseVisionResponse(response string) (filename, content string) {
	const defaultFilename = "image-upload"

	if response == "" {
		return defaultFilename, ""
	}

	// Find the first newline to extract the first line
	firstNewline := strings.Index(response, "\n")
	var firstLine string
	var restOfContent string

	if firstNewline == -1 {
		// No newline, entire response is first line
		firstLine = response
		restOfContent = ""
	} else {
		firstLine = response[:firstNewline]
		restOfContent = response[firstNewline+1:]
	}

	// Check if first line starts with "FILENAME:"
	const prefix = "FILENAME:"
	if !strings.HasPrefix(firstLine, prefix) {
		// No FILENAME prefix, return default filename and full response as content
		return defaultFilename, response
	}

	// Extract the filename after the prefix
	rawFilename := strings.TrimSpace(strings.TrimPrefix(firstLine, prefix))
	if rawFilename == "" {
		// Empty filename after prefix, use default
		return defaultFilename, strings.TrimLeft(restOfContent, "\n")
	}

	// Sanitize the filename
	sanitizedFilename := sanitizeSmartFilename(rawFilename)

	// Trim leading newlines from content
	content = strings.TrimLeft(restOfContent, "\n")

	return sanitizedFilename, content
}

// sanitizeSmartFilename sanitizes a filename for use in file paths.
// - Converts to lowercase
// - Replaces spaces and underscores with dashes
// - Removes invalid characters (keeps only alphanumeric and dashes)
// - Truncates to 40 characters
// - Falls back to "image-upload" if result is empty
func sanitizeSmartFilename(filename string) string {
	// Convert to lowercase
	filename = strings.ToLower(filename)

	// Replace spaces and underscores with dashes
	filename = strings.ReplaceAll(filename, " ", "-")
	filename = strings.ReplaceAll(filename, "_", "-")

	// Keep only alphanumeric characters and dashes
	var result strings.Builder
	for _, r := range filename {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	sanitized := result.String()

	// Truncate to 40 characters
	if len(sanitized) > 40 {
		sanitized = sanitized[:40]
	}

	// Fall back to default if empty
	if sanitized == "" {
		sanitized = "image-upload"
	}

	return sanitized
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
