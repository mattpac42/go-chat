package service

import (
	"context"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

// CompletenessChecker validates that all file references in a project resolve correctly.
type CompletenessChecker struct {
	fileRepo repository.FileRepository
	logger   zerolog.Logger
}

// NewCompletenessChecker creates a new completeness checker.
func NewCompletenessChecker(fileRepo repository.FileRepository, logger zerolog.Logger) *CompletenessChecker {
	return &CompletenessChecker{
		fileRepo: fileRepo,
		logger:   logger.With().Str("component", "completeness_checker").Logger(),
	}
}

// Regex patterns for extracting file references
var (
	// HTML patterns
	htmlScriptPattern = regexp.MustCompile(`<script[^>]+src=["']([^"']+)["']`)
	htmlLinkPattern   = regexp.MustCompile(`<link[^>]+href=["']([^"']+\.css)["']`)
	htmlImgPattern    = regexp.MustCompile(`<img[^>]+src=["']([^"']+)["']`)

	// JavaScript/TypeScript patterns
	jsImportPattern  = regexp.MustCompile(`import\s+(?:.*\s+from\s+)?["']([^"']+)["']`)
	jsRequirePattern = regexp.MustCompile(`require\s*\(\s*["']([^"']+)["']\s*\)`)

	// CSS patterns
	cssImportPattern = regexp.MustCompile(`@import\s+(?:url\s*\(\s*)?["']?([^"')]+)["']?\s*\)?`)
	cssUrlPattern    = regexp.MustCompile(`url\s*\(\s*["']?([^"')]+)["']?\s*\)`)
)

// Check validates all files in a project and returns a completeness report.
func (c *CompletenessChecker) Check(ctx context.Context, projectID uuid.UUID) (*model.CompletenessReport, error) {
	c.logger.Debug().Str("projectId", projectID.String()).Msg("starting completeness check")

	// Get all files in project
	files, err := c.fileRepo.GetByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Build a set of existing file paths for quick lookup
	existingFiles := make(map[string]bool)
	for _, f := range files {
		// Store both with and without leading slash for flexible matching
		existingFiles[f.Path] = true
		existingFiles[strings.TrimPrefix(f.Path, "/")] = true
		existingFiles[f.Filename] = true
	}

	var issues []model.CompletenessIssue
	issueID := 0

	// Check each file for references
	for _, file := range files {
		// Get file content
		fileWithContent, err := c.fileRepo.GetByID(ctx, file.ID)
		if err != nil {
			c.logger.Warn().Err(err).Str("fileId", file.ID.String()).Msg("failed to get file content")
			continue
		}

		content := fileWithContent.Content
		if content == "" {
			continue
		}

		// Extract and validate references based on file type
		refs := c.extractReferences(file.Filename, content)

		for _, ref := range refs {
			// Skip external URLs
			if isExternalURL(ref.Path) {
				continue
			}

			// Resolve relative path
			resolvedPath := c.resolvePath(file.Path, ref.Path)

			// Check if file exists
			if !c.fileExists(existingFiles, resolvedPath, ref.Path) {
				issueID++
				severity := c.getSeverity(ref.Path, file.Filename, ref.Type)

				issues = append(issues, model.CompletenessIssue{
					ID:            generateIssueID(issueID),
					Severity:      severity,
					Type:          "missing_file",
					MissingFile:   ref.Path,
					ReferencedBy:  file.Filename,
					ReferenceType: ref.Type,
					LineNumber:    ref.LineNumber,
					Context:       ref.Context,
					AutoFixable:   severity == model.SeverityCritical,
				})
			}
		}
	}

	// Determine overall status
	status := model.StatusPass
	autoFixable := 0
	for _, issue := range issues {
		if issue.Severity == model.SeverityCritical {
			status = model.StatusCritical
		} else if issue.Severity == model.SeverityWarning && status != model.StatusCritical {
			status = model.StatusWarning
		}
		if issue.AutoFixable {
			autoFixable++
		}
	}

	report := &model.CompletenessReport{
		ProjectID:    projectID,
		CheckedAt:    time.Now(),
		Status:       status,
		Issues:       issues,
		FilesChecked: len(files),
		AutoFixable:  autoFixable,
	}

	c.logger.Info().
		Str("projectId", projectID.String()).
		Str("status", string(status)).
		Int("issues", len(issues)).
		Int("filesChecked", len(files)).
		Msg("completeness check completed")

	return report, nil
}

// FileReference represents a reference found in a file.
type FileReference struct {
	Path       string
	Type       string // "script", "stylesheet", "import", "image"
	LineNumber int
	Context    string
}

// extractReferences extracts file references from content based on file type.
func (c *CompletenessChecker) extractReferences(filename, content string) []FileReference {
	ext := strings.ToLower(filepath.Ext(filename))
	var refs []FileReference

	switch ext {
	case ".html", ".htm":
		refs = append(refs, c.extractHTMLReferences(content)...)
	case ".js", ".jsx", ".ts", ".tsx", ".mjs":
		refs = append(refs, c.extractJSReferences(content)...)
	case ".css", ".scss", ".sass", ".less":
		refs = append(refs, c.extractCSSReferences(content)...)
	}

	return refs
}

// extractHTMLReferences extracts script, stylesheet, and image references from HTML.
func (c *CompletenessChecker) extractHTMLReferences(content string) []FileReference {
	var refs []FileReference
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		// Scripts
		for _, match := range htmlScriptPattern.FindAllStringSubmatch(line, -1) {
			if len(match) >= 2 {
				refs = append(refs, FileReference{
					Path:       match[1],
					Type:       "script",
					LineNumber: lineNum + 1,
					Context:    strings.TrimSpace(line),
				})
			}
		}

		// Stylesheets
		for _, match := range htmlLinkPattern.FindAllStringSubmatch(line, -1) {
			if len(match) >= 2 {
				refs = append(refs, FileReference{
					Path:       match[1],
					Type:       "stylesheet",
					LineNumber: lineNum + 1,
					Context:    strings.TrimSpace(line),
				})
			}
		}

		// Images (warning level only)
		for _, match := range htmlImgPattern.FindAllStringSubmatch(line, -1) {
			if len(match) >= 2 {
				refs = append(refs, FileReference{
					Path:       match[1],
					Type:       "image",
					LineNumber: lineNum + 1,
					Context:    strings.TrimSpace(line),
				})
			}
		}
	}

	return refs
}

// extractJSReferences extracts import and require references from JavaScript.
func (c *CompletenessChecker) extractJSReferences(content string) []FileReference {
	var refs []FileReference
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		// ES6 imports
		for _, match := range jsImportPattern.FindAllStringSubmatch(line, -1) {
			if len(match) >= 2 {
				path := match[1]
				// Only check relative imports
				if strings.HasPrefix(path, ".") {
					refs = append(refs, FileReference{
						Path:       path,
						Type:       "import",
						LineNumber: lineNum + 1,
						Context:    strings.TrimSpace(line),
					})
				}
			}
		}

		// CommonJS require
		for _, match := range jsRequirePattern.FindAllStringSubmatch(line, -1) {
			if len(match) >= 2 {
				path := match[1]
				// Only check relative requires
				if strings.HasPrefix(path, ".") {
					refs = append(refs, FileReference{
						Path:       path,
						Type:       "import",
						LineNumber: lineNum + 1,
						Context:    strings.TrimSpace(line),
					})
				}
			}
		}
	}

	return refs
}

// extractCSSReferences extracts @import and url() references from CSS.
func (c *CompletenessChecker) extractCSSReferences(content string) []FileReference {
	var refs []FileReference
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		// @import
		for _, match := range cssImportPattern.FindAllStringSubmatch(line, -1) {
			if len(match) >= 2 && !isExternalURL(match[1]) {
				refs = append(refs, FileReference{
					Path:       match[1],
					Type:       "stylesheet",
					LineNumber: lineNum + 1,
					Context:    strings.TrimSpace(line),
				})
			}
		}

		// url() - but not for data: URIs or external URLs
		for _, match := range cssUrlPattern.FindAllStringSubmatch(line, -1) {
			if len(match) >= 2 {
				path := match[1]
				if !isExternalURL(path) && !strings.HasPrefix(path, "data:") {
					refs = append(refs, FileReference{
						Path:       path,
						Type:       "image",
						LineNumber: lineNum + 1,
						Context:    strings.TrimSpace(line),
					})
				}
			}
		}
	}

	return refs
}

// resolvePath resolves a relative path against a base file path.
func (c *CompletenessChecker) resolvePath(basePath, refPath string) string {
	if filepath.IsAbs(refPath) || strings.HasPrefix(refPath, "/") {
		return strings.TrimPrefix(refPath, "/")
	}

	baseDir := filepath.Dir(basePath)
	resolved := filepath.Join(baseDir, refPath)
	return filepath.Clean(resolved)
}

// fileExists checks if a file exists in the project.
func (c *CompletenessChecker) fileExists(existingFiles map[string]bool, resolvedPath, originalPath string) bool {
	// Try resolved path
	if existingFiles[resolvedPath] {
		return true
	}

	// Try original path
	if existingFiles[originalPath] {
		return true
	}

	// Try with common extensions for JS imports without extension
	if !strings.Contains(originalPath, ".") {
		extensions := []string{".js", ".jsx", ".ts", ".tsx", ".mjs"}
		for _, ext := range extensions {
			if existingFiles[resolvedPath+ext] || existingFiles[originalPath+ext] {
				return true
			}
		}
		// Try index files
		for _, ext := range extensions {
			indexPath := filepath.Join(resolvedPath, "index"+ext)
			if existingFiles[indexPath] {
				return true
			}
		}
	}

	return false
}

// getSeverity determines the severity of a missing file.
func (c *CompletenessChecker) getSeverity(missingFile, referencedBy, refType string) model.Severity {
	ext := strings.ToLower(filepath.Ext(missingFile))
	refExt := strings.ToLower(filepath.Ext(referencedBy))

	// Critical: Missing JS/TS referenced by HTML (app won't work)
	if (ext == ".js" || ext == ".ts" || ext == ".jsx" || ext == ".tsx") && (refExt == ".html" || refExt == ".htm") {
		return model.SeverityCritical
	}

	// Critical: Missing CSS referenced by HTML (app will look broken)
	if ext == ".css" && (refExt == ".html" || refExt == ".htm") {
		return model.SeverityCritical
	}

	// Critical: Missing JS module imports (app won't work)
	if refType == "import" && (ext == ".js" || ext == ".ts" || ext == ".jsx" || ext == ".tsx" || ext == "") {
		return model.SeverityCritical
	}

	// Warning: Missing images
	if refType == "image" {
		return model.SeverityWarning
	}

	// Warning: Missing CSS imports
	if refType == "stylesheet" && ext == ".css" {
		return model.SeverityWarning
	}

	return model.SeverityInfo
}

// isExternalURL checks if a path is an external URL.
func isExternalURL(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "//") ||
		strings.HasPrefix(path, "data:")
}

// generateIssueID generates a unique issue ID.
func generateIssueID(num int) string {
	return "issue-" + uuid.New().String()[:8]
}
