package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/service"
)

func TestUploadHandler_Upload(t *testing.T) {
	logger := zerolog.Nop()

	t.Run("successfully uploads and converts PNG image to markdown with smart filename", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()
		fileMetadataRepo := repository.NewMockFileMetadataRepository()
		fileSourceRepo := repository.NewMockFileSourceRepository()
		mockVision := service.NewMockClaudeVision()

		// Response with FILENAME prefix (smart filename feature)
		mockVision.SetDefaultResponse(`FILENAME: webapp-screenshot

## Screenshot Analysis

This is a screenshot of a web application.

### Extracted Text
- Button: "Submit"
- Header: "Welcome to the App"
`)

		project, _ := projectRepo.Create(nil, "Test Project")

		handler := NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, mockVision, logger)
		router := gin.New()
		router.POST("/api/projects/:id/upload", handler.Upload)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "screenshot.png")
		part.Write([]byte("fake PNG data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID.String()+"/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response model.UploadResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.File.ID)
		assert.Contains(t, response.File.Path, "sources/")
		assert.Contains(t, response.File.Path, "webapp-screenshot") // Smart filename used
		assert.Contains(t, response.File.Path, ".md")
		assert.Contains(t, response.File.Content, "Screenshot Analysis")
		assert.NotContains(t, response.File.Content, "FILENAME:") // FILENAME line stripped from content
		assert.Equal(t, "Source Materials", response.File.FunctionalGroup)

		assert.Equal(t, "screenshot.png", response.Source.OriginalFilename)
		assert.Equal(t, "image/png", response.Source.OriginalMimeType)
		assert.Greater(t, response.Source.OriginalSizeBytes, int64(0))

		// Verify Claude Vision was called
		assert.Equal(t, 1, mockVision.GetAnalyzeCallCount())
		_, mimeType, prompt := mockVision.GetLastCall()
		assert.Equal(t, "image/png", mimeType)
		assert.Contains(t, prompt, "FILENAME:")
	})

	t.Run("falls back to default filename when vision response has no FILENAME prefix", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()
		fileMetadataRepo := repository.NewMockFileMetadataRepository()
		fileSourceRepo := repository.NewMockFileSourceRepository()
		mockVision := service.NewMockClaudeVision()

		// Response without FILENAME prefix (legacy format)
		mockVision.SetDefaultResponse(`## Image Analysis

This is an image without a filename prefix.
`)

		project, _ := projectRepo.Create(nil, "Test Project")

		handler := NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, mockVision, logger)
		router := gin.New()
		router.POST("/api/projects/:id/upload", handler.Upload)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "random-file.png")
		part.Write([]byte("fake PNG data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID.String()+"/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response model.UploadResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Falls back to "image-upload" when no FILENAME prefix
		assert.Contains(t, response.File.Path, "sources/image-upload")
		assert.Contains(t, response.File.Path, ".md")
		// Content should be preserved as-is when no FILENAME prefix
		assert.Contains(t, response.File.Content, "## Image Analysis")
	})

	t.Run("successfully uploads JPEG image", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()
		fileMetadataRepo := repository.NewMockFileMetadataRepository()
		fileSourceRepo := repository.NewMockFileSourceRepository()
		mockVision := service.NewMockClaudeVision()

		project, _ := projectRepo.Create(nil, "Test Project")

		handler := NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, mockVision, logger)
		router := gin.New()
		router.POST("/api/projects/:id/upload", handler.Upload)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "photo.jpg")
		part.Write([]byte("fake JPEG data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID.String()+"/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response model.UploadResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "photo.jpg", response.Source.OriginalFilename)
	})

	t.Run("returns 400 for unsupported file type", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()
		fileMetadataRepo := repository.NewMockFileMetadataRepository()
		fileSourceRepo := repository.NewMockFileSourceRepository()
		mockVision := service.NewMockClaudeVision()

		project, _ := projectRepo.Create(nil, "Test Project")

		handler := NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, mockVision, logger)
		router := gin.New()
		router.POST("/api/projects/:id/upload", handler.Upload)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		h := make(map[string][]string)
		h["Content-Disposition"] = []string{`form-data; name="file"; filename="document.pdf"`}
		h["Content-Type"] = []string{"application/pdf"}
		part, _ := writer.CreatePart(h)
		part.Write([]byte("fake PDF data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID.String()+"/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "unsupported file type")
	})

	t.Run("returns 404 for non-existent project", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()
		fileMetadataRepo := repository.NewMockFileMetadataRepository()
		fileSourceRepo := repository.NewMockFileSourceRepository()
		mockVision := service.NewMockClaudeVision()

		handler := NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, mockVision, logger)
		router := gin.New()
		router.POST("/api/projects/:id/upload", handler.Upload)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "screenshot.png")
		part.Write([]byte("fake PNG data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/projects/00000000-0000-0000-0000-000000000000/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 400 for invalid project ID", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()
		fileMetadataRepo := repository.NewMockFileMetadataRepository()
		fileSourceRepo := repository.NewMockFileSourceRepository()
		mockVision := service.NewMockClaudeVision()

		handler := NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, mockVision, logger)
		router := gin.New()
		router.POST("/api/projects/:id/upload", handler.Upload)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "screenshot.png")
		part.Write([]byte("fake PNG data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/projects/invalid-uuid/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 400 when no file provided", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()
		fileMetadataRepo := repository.NewMockFileMetadataRepository()
		fileSourceRepo := repository.NewMockFileSourceRepository()
		mockVision := service.NewMockClaudeVision()

		project, _ := projectRepo.Create(nil, "Test Project")

		handler := NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, mockVision, logger)
		router := gin.New()
		router.POST("/api/projects/:id/upload", handler.Upload)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID.String()+"/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "no file provided")
	})

	t.Run("returns 500 when Claude Vision fails", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()
		fileMetadataRepo := repository.NewMockFileMetadataRepository()
		fileSourceRepo := repository.NewMockFileSourceRepository()
		mockVision := service.MockClaudeVisionWithError("API error")

		project, _ := projectRepo.Create(nil, "Test Project")

		handler := NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, mockVision, logger)
		router := gin.New()
		router.POST("/api/projects/:id/upload", handler.Upload)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "screenshot.png")
		part.Write([]byte("fake PNG data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID.String()+"/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "failed to analyze image")
	})

	t.Run("saves file to sources folder with smart filename from vision", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()
		fileMetadataRepo := repository.NewMockFileMetadataRepository()
		fileSourceRepo := repository.NewMockFileSourceRepository()
		mockVision := service.NewMockClaudeVision()

		// Vision returns a smart filename regardless of original filename
		mockVision.SetDefaultResponse(`FILENAME: dashboard-ui-design

## Dashboard Design

A dashboard interface with charts and widgets.
`)

		project, _ := projectRepo.Create(nil, "Test Project")

		handler := NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, mockVision, logger)
		router := gin.New()
		router.POST("/api/projects/:id/upload", handler.Upload)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		// Original filename is "My Screenshot.png" but smart filename from vision will be used
		part, _ := writer.CreateFormFile("file", "My Screenshot.png")
		part.Write([]byte("fake PNG data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID.String()+"/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response model.UploadResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.True(t, len(response.File.Path) > 0)
		assert.Contains(t, response.File.Path, "sources/")
		// Uses smart filename from vision, not original filename
		assert.Contains(t, response.File.Path, "dashboard-ui-design")
		assert.NotContains(t, response.File.Path, "my-screenshot")
		assert.Contains(t, response.File.Path, ".md")
	})

	t.Run("extracts short description from markdown content", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()
		fileMetadataRepo := repository.NewMockFileMetadataRepository()
		fileSourceRepo := repository.NewMockFileSourceRepository()
		mockVision := service.NewMockClaudeVision()

		mockVision.SetDefaultResponse(`FILENAME: login-screen

## Login Screen

The image shows a login form with username and password fields.

### Elements
- Username input
- Password input
- Submit button
`)

		project, _ := projectRepo.Create(nil, "Test Project")

		handler := NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, mockVision, logger)
		router := gin.New()
		router.POST("/api/projects/:id/upload", handler.Upload)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "screenshot.png")
		part.Write([]byte("fake PNG data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID.String()+"/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response model.UploadResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		// Short description should be extracted from first meaningful line (without ## prefix)
		assert.Equal(t, "Login Screen", response.File.ShortDescription)
	})

	t.Run("supports all allowed image types", func(t *testing.T) {
		testCases := []struct {
			filename string
			mimeType string
		}{
			{"image.png", "image/png"},
			{"image.jpg", "image/jpeg"},
			{"image.gif", "image/gif"},
			{"image.webp", "image/webp"},
		}

		for _, tc := range testCases {
			t.Run(tc.mimeType, func(t *testing.T) {
				projectRepo := repository.NewMockProjectRepository()
				fileRepo := repository.NewMockFileRepository()
				fileMetadataRepo := repository.NewMockFileMetadataRepository()
				fileSourceRepo := repository.NewMockFileSourceRepository()
				mockVision := service.NewMockClaudeVision()

				project, _ := projectRepo.Create(nil, "Test Project")

				handler := NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, mockVision, logger)
				router := gin.New()
				router.POST("/api/projects/:id/upload", handler.Upload)

				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				h := make(map[string][]string)
				h["Content-Disposition"] = []string{`form-data; name="file"; filename="` + tc.filename + `"`}
				h["Content-Type"] = []string{tc.mimeType}
				part, _ := writer.CreatePart(h)
				part.Write([]byte("fake image data"))
				writer.Close()

				req := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID.String()+"/upload", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK for %s", tc.mimeType)
			})
		}
	})
}

func TestSanitizeUploadFilename(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"screenshot.png", "screenshot"},
		{"My Screenshot.png", "my-screenshot"},
		{"Image With <Special> Chars!.jpg", "image-with-special-chars"},
		{"a-b_c.gif", "a-b_c"},
		{"", "upload"},
		{".png", "upload"},
		{"A" + string(make([]byte, 100)) + ".png", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}, // truncated to 40 chars
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := sanitizeUploadFilename(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParseVisionResponse(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedFilename string
		expectedContent  string
	}{
		{
			name:             "parses valid filename and content",
			input:            "FILENAME: bakery-menu-mockup\n\n## Menu Design\n\nThis is a bakery menu.",
			expectedFilename: "bakery-menu-mockup",
			expectedContent:  "## Menu Design\n\nThis is a bakery menu.",
		},
		{
			name:             "handles filename with extra whitespace",
			input:            "FILENAME:   login-screen   \n\n## Login Screen\n\nA login form.",
			expectedFilename: "login-screen",
			expectedContent:  "## Login Screen\n\nA login form.",
		},
		{
			name:             "falls back to default when no FILENAME prefix",
			input:            "## Screenshot Analysis\n\nThis is a screenshot.",
			expectedFilename: "image-upload",
			expectedContent:  "## Screenshot Analysis\n\nThis is a screenshot.",
		},
		{
			name:             "falls back when FILENAME line is empty",
			input:            "FILENAME:\n\n## Content\n\nSome content.",
			expectedFilename: "image-upload",
			expectedContent:  "## Content\n\nSome content.",
		},
		{
			name:             "converts to lowercase",
			input:            "FILENAME: Login-Screen-Design\n\n## Login Screen",
			expectedFilename: "login-screen-design",
			expectedContent:  "## Login Screen",
		},
		{
			name:             "removes invalid characters",
			input:            "FILENAME: my_image@#$test!.png\n\n## Content",
			expectedFilename: "my-imagetest",
			expectedContent:  "## Content",
		},
		{
			name:             "truncates long filenames to 40 chars",
			input:            "FILENAME: this-is-a-very-long-filename-that-should-be-truncated-to-forty-characters\n\n## Content",
			expectedFilename: "this-is-a-very-long-filename-that-shoul",
			expectedContent:  "## Content",
		},
		{
			name:             "replaces spaces with dashes",
			input:            "FILENAME: my image name\n\n## Content",
			expectedFilename: "my-image-name",
			expectedContent:  "## Content",
		},
		{
			name:             "handles underscores",
			input:            "FILENAME: my_image_name\n\n## Content",
			expectedFilename: "my-image-name",
			expectedContent:  "## Content",
		},
		{
			name:             "handles empty input",
			input:            "",
			expectedFilename: "image-upload",
			expectedContent:  "",
		},
		{
			name:             "handles FILENAME only with no content after",
			input:            "FILENAME: test-image",
			expectedFilename: "test-image",
			expectedContent:  "",
		},
		{
			name:             "preserves content with leading newlines trimmed",
			input:            "FILENAME: test\n\n\n## Heading\n\nParagraph",
			expectedFilename: "test",
			expectedContent:  "## Heading\n\nParagraph",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename, content := parseVisionResponse(tc.input)
			assert.Equal(t, tc.expectedFilename, filename)
			assert.Equal(t, tc.expectedContent, content)
		})
	}
}

func TestExtractShortDescription(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "extracts from header",
			input:    "## Login Screen\n\nSome content here.",
			expected: "Login Screen",
		},
		{
			name:     "skips empty lines",
			input:    "\n\n\n## Dashboard View\n\nContent",
			expected: "Dashboard View",
		},
		{
			name:     "truncates long descriptions",
			input:    "## " + string(make([]byte, 150)),
			expected: string(make([]byte, 97)) + "...",
		},
		{
			name:     "handles plain text",
			input:    "This is a simple description.",
			expected: "This is a simple description.",
		},
		{
			name:     "handles empty content",
			input:    "",
			expected: "Uploaded image content",
		},
		{
			name:     "handles only whitespace",
			input:    "   \n\n   ",
			expected: "Uploaded image content",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := extractShortDescription(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
