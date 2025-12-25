package markdown

import (
	"path/filepath"
	"regexp"
	"strings"
)

// CodeBlock represents a code block extracted from markdown.
type CodeBlock struct {
	Language   string `json:"language"`
	Filename   string `json:"filename,omitempty"`
	Code       string `json:"code"`
	StartIndex int    `json:"startIndex"`
	EndIndex   int    `json:"endIndex"`
}

// codeBlockPattern matches markdown fenced code blocks.
// It captures the language/filename info (optional) and the code content.
// Supports formats like:
//   - ```typescript
//   - ```typescript:src/components/Button.tsx
//   - ```src/components/Button.tsx
var codeBlockPattern = regexp.MustCompile("(?s)```([^\\n`]*)\\n?(.*?)```")

// ExtractCodeBlocks parses markdown content and extracts all fenced code blocks.
// It returns a slice of CodeBlock structs with language, code, and position information.
func ExtractCodeBlocks(content string) []CodeBlock {
	matches := codeBlockPattern.FindAllStringSubmatchIndex(content, -1)
	if matches == nil {
		return []CodeBlock{}
	}

	blocks := make([]CodeBlock, 0, len(matches))

	for _, match := range matches {
		// match[0] and match[1] are the start and end of the entire match
		// match[2] and match[3] are the start and end of the language/filename group
		// match[4] and match[5] are the start and end of the code group

		startIndex := match[0]
		endIndex := match[1]

		var language, filename string
		if match[2] != -1 && match[3] != -1 {
			langInfo := strings.TrimSpace(content[match[2]:match[3]])
			language, filename = parseLanguageAndFilename(langInfo)
		}

		var code string
		if match[4] != -1 && match[5] != -1 {
			code = content[match[4]:match[5]]
			// Trim trailing newline if present
			code = strings.TrimSuffix(code, "\n")
		}

		blocks = append(blocks, CodeBlock{
			Language:   language,
			Filename:   filename,
			Code:       code,
			StartIndex: startIndex,
			EndIndex:   endIndex,
		})
	}

	return blocks
}

// parseLanguageAndFilename extracts language and filename from code block header.
// Supports formats:
//   - "typescript" -> language="typescript", filename=""
//   - "typescript:src/Button.tsx" -> language="typescript", filename="src/Button.tsx"
//   - "src/Button.tsx" -> language="tsx", filename="src/Button.tsx"
//   - ":src/Button.tsx" -> language="tsx", filename="src/Button.tsx"
func parseLanguageAndFilename(info string) (language, filename string) {
	if info == "" {
		return "", ""
	}

	// Check for colon separator (language:filename format)
	if colonIdx := strings.Index(info, ":"); colonIdx != -1 {
		language = strings.TrimSpace(info[:colonIdx])
		filename = strings.TrimSpace(info[colonIdx+1:])
		// If language is empty but we have a filename, infer language from extension
		if language == "" && filename != "" {
			language = inferLanguageFromFilename(filename)
		}
		return language, filename
	}

	// Check if info looks like a file path (contains / or . with extension)
	if looksLikeFilePath(info) {
		filename = info
		language = inferLanguageFromFilename(filename)
		return language, filename
	}

	// Just a language identifier
	return info, ""
}

// looksLikeFilePath checks if a string appears to be a file path.
func looksLikeFilePath(s string) bool {
	// Contains path separator
	if strings.Contains(s, "/") {
		return true
	}
	// Has a file extension (e.g., "Button.tsx", "main.go")
	ext := filepath.Ext(s)
	return ext != "" && len(ext) > 1 && len(ext) <= 10
}

// inferLanguageFromFilename returns a language identifier based on file extension.
func inferLanguageFromFilename(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".ts":
		return "typescript"
	case ".tsx":
		return "tsx"
	case ".js":
		return "javascript"
	case ".jsx":
		return "jsx"
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".rb":
		return "ruby"
	case ".rs":
		return "rust"
	case ".java":
		return "java"
	case ".kt", ".kts":
		return "kotlin"
	case ".swift":
		return "swift"
	case ".c":
		return "c"
	case ".cpp", ".cc", ".cxx":
		return "cpp"
	case ".h", ".hpp":
		return "cpp"
	case ".cs":
		return "csharp"
	case ".php":
		return "php"
	case ".sql":
		return "sql"
	case ".html", ".htm":
		return "html"
	case ".css":
		return "css"
	case ".scss", ".sass":
		return "scss"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".xml":
		return "xml"
	case ".md", ".markdown":
		return "markdown"
	case ".sh", ".bash":
		return "bash"
	case ".dockerfile":
		return "dockerfile"
	default:
		// Check for Dockerfile without extension
		base := filepath.Base(filename)
		if strings.ToLower(base) == "dockerfile" {
			return "dockerfile"
		}
		return ""
	}
}
