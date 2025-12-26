package handler

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestFileHandler_DownloadProjectZip(t *testing.T) {
	t.Run("downloads zip with all project files", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()

		project, _ := projectRepo.Create(nil, "Test Project")
		_, _ = fileRepo.SaveFile(nil, project.ID, "src/main.go", "go", "package main\n\nfunc main() {}")
		_, _ = fileRepo.SaveFile(nil, project.ID, "public/index.html", "html", "<html><body>Hello</body></html>")

		handler := NewFileHandler(fileRepo, projectRepo, nil)
		router := gin.New()
		router.GET("/api/projects/:id/download", handler.DownloadProjectZip)

		req := httptest.NewRequest(http.MethodGet, "/api/projects/"+project.ID.String()+"/download", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
		assert.Contains(t, w.Header().Get("Content-Disposition"), "Test-Project.zip")

		// Verify zip contents
		zipReader, err := zip.NewReader(bytes.NewReader(w.Body.Bytes()), int64(w.Body.Len()))
		require.NoError(t, err)
		assert.Len(t, zipReader.File, 2)

		// Check file paths are preserved
		filePaths := make(map[string]bool)
		for _, f := range zipReader.File {
			filePaths[f.Name] = true
		}
		assert.True(t, filePaths["src/main.go"])
		assert.True(t, filePaths["public/index.html"])

		// Verify file content
		for _, f := range zipReader.File {
			rc, err := f.Open()
			require.NoError(t, err)
			content, err := io.ReadAll(rc)
			require.NoError(t, err)
			rc.Close()

			if f.Name == "src/main.go" {
				assert.Equal(t, "package main\n\nfunc main() {}", string(content))
			}
		}
	})

	t.Run("returns 404 when project not found", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()

		handler := NewFileHandler(fileRepo, projectRepo, nil)
		router := gin.New()
		router.GET("/api/projects/:id/download", handler.DownloadProjectZip)

		req := httptest.NewRequest(http.MethodGet, "/api/projects/00000000-0000-0000-0000-000000000000/download", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 404 when project has no files", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()

		project, _ := projectRepo.Create(nil, "Empty Project")

		handler := NewFileHandler(fileRepo, projectRepo, nil)
		router := gin.New()
		router.GET("/api/projects/:id/download", handler.DownloadProjectZip)

		req := httptest.NewRequest(http.MethodGet, "/api/projects/"+project.ID.String()+"/download", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "no files found in project", response["error"])
	})

	t.Run("returns 400 for invalid project ID", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()

		handler := NewFileHandler(fileRepo, projectRepo, nil)
		router := gin.New()
		router.GET("/api/projects/:id/download", handler.DownloadProjectZip)

		req := httptest.NewRequest(http.MethodGet, "/api/projects/invalid-uuid/download", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("sanitizes project title for filename", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()

		project, _ := projectRepo.Create(nil, "My Project <With> Special/Chars!")
		_, _ = fileRepo.SaveFile(nil, project.ID, "file.txt", "text", "content")

		handler := NewFileHandler(fileRepo, projectRepo, nil)
		router := gin.New()
		router.GET("/api/projects/:id/download", handler.DownloadProjectZip)

		req := httptest.NewRequest(http.MethodGet, "/api/projects/"+project.ID.String()+"/download", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		// Filename should only contain safe characters
		assert.Contains(t, w.Header().Get("Content-Disposition"), "My-Project-With-SpecialChars.zip")
	})
}

func TestFileHandler_DownloadFile(t *testing.T) {
	t.Run("downloads single file with correct headers", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()

		project, _ := projectRepo.Create(nil, "Test Project")
		file, _ := fileRepo.SaveFile(nil, project.ID, "src/main.go", "go", "package main\n\nfunc main() {}")

		handler := NewFileHandler(fileRepo, projectRepo, nil)
		router := gin.New()
		router.GET("/api/files/:id/download", handler.DownloadFile)

		req := httptest.NewRequest(http.MethodGet, "/api/files/"+file.ID.String()+"/download", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
		assert.Contains(t, w.Header().Get("Content-Disposition"), "src/main.go")
		assert.Equal(t, "package main\n\nfunc main() {}", w.Body.String())
	})

	t.Run("returns 404 for non-existent file", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()

		handler := NewFileHandler(fileRepo, projectRepo, nil)
		router := gin.New()
		router.GET("/api/files/:id/download", handler.DownloadFile)

		req := httptest.NewRequest(http.MethodGet, "/api/files/00000000-0000-0000-0000-000000000000/download", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 400 for invalid file ID", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()

		handler := NewFileHandler(fileRepo, projectRepo, nil)
		router := gin.New()
		router.GET("/api/files/:id/download", handler.DownloadFile)

		req := httptest.NewRequest(http.MethodGet, "/api/files/invalid-uuid/download", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("sets correct content type for various file extensions", func(t *testing.T) {
		testCases := []struct {
			filename    string
			expected    string
		}{
			{"main.go", "text/x-go"},
			{"app.ts", "text/typescript"},
			{"style.css", "text/css; charset=utf-8"},
			{"script.js", "text/javascript; charset=utf-8"},
			{"data.json", "application/json"},
			{"readme.md", "text/markdown"},
			{"unknown.xyz", "application/octet-stream"},
		}

		for _, tc := range testCases {
			t.Run(tc.filename, func(t *testing.T) {
				projectRepo := repository.NewMockProjectRepository()
				fileRepo := repository.NewMockFileRepository()

				project, _ := projectRepo.Create(nil, "Test Project")
				file, _ := fileRepo.SaveFile(nil, project.ID, tc.filename, "", "content")

				handler := NewFileHandler(fileRepo, projectRepo, nil)
				router := gin.New()
				router.GET("/api/files/:id/download", handler.DownloadFile)

				req := httptest.NewRequest(http.MethodGet, "/api/files/"+file.ID.String()+"/download", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, tc.expected, w.Header().Get("Content-Type"))
			})
		}
	})
}

func TestFileHandler_GetFile(t *testing.T) {
	t.Run("returns file by ID", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()

		project, _ := projectRepo.Create(nil, "Test Project")
		file, _ := fileRepo.SaveFile(nil, project.ID, "src/main.go", "go", "package main")

		handler := NewFileHandler(fileRepo, projectRepo, nil)
		router := gin.New()
		router.GET("/api/files/:id", handler.GetFile)

		req := httptest.NewRequest(http.MethodGet, "/api/files/"+file.ID.String(), nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response model.GetFileResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, file.ID, response.ID)
		assert.Equal(t, "src/main.go", response.Path)
		assert.Equal(t, "package main", response.Content)
	})

	t.Run("returns 404 for non-existent file", func(t *testing.T) {
		// Arrange
		projectRepo := repository.NewMockProjectRepository()
		fileRepo := repository.NewMockFileRepository()

		handler := NewFileHandler(fileRepo, projectRepo, nil)
		router := gin.New()
		router.GET("/api/files/:id", handler.GetFile)

		req := httptest.NewRequest(http.MethodGet, "/api/files/00000000-0000-0000-0000-000000000000", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSanitizeFilename(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Simple Project", "Simple-Project"},
		{"My <Project>", "My-Project"},
		{"Test/Project\\Name", "TestProjectName"},
		{"Project: The Sequel", "Project-The-Sequel"},
		{"a-b_c-d_e", "a-b_c-d_e"},
		{"", ""},
		{"A" + string(make([]byte, 100)), "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := sanitizeFilename(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetContentType(t *testing.T) {
	testCases := []struct {
		filename string
		expected string
	}{
		{"main.go", "text/x-go"},
		{"app.rs", "text/x-rust"},
		{"index.ts", "text/typescript"},
		{"component.tsx", "text/typescript-jsx"},
		{"app.vue", "text/x-vue"},
		{"component.svelte", "text/x-svelte"},
		{"config.yaml", "text/yaml"},
		{"config.yml", "text/yaml"},
		{"Cargo.toml", "text/toml"},
		{"readme.md", "text/markdown"},
		{"schema.sql", "application/sql"},
		{"build.sh", "application/x-sh"},
		{"noextension", "application/octet-stream"},
		{"unknown.xyz", "application/octet-stream"},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			result := getContentType(tc.filename)
			assert.Equal(t, tc.expected, result)
		})
	}
}
