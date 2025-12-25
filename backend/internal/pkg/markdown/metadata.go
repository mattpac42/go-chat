package markdown

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileMetadata represents App Map metadata extracted from code blocks.
type FileMetadata struct {
	ShortDescription string `yaml:"short_description"`
	LongDescription  string `yaml:"long_description"`
	FunctionalGroup  string `yaml:"functional_group"`
}

// IsEmpty returns true if no metadata fields are populated.
func (m *FileMetadata) IsEmpty() bool {
	return m.ShortDescription == "" && m.LongDescription == "" && m.FunctionalGroup == ""
}

// CodeBlockWithMetadata extends CodeBlock with App Map metadata.
type CodeBlockWithMetadata struct {
	Language   string        `json:"language"`
	Filename   string        `json:"filename,omitempty"`
	Code       string        `json:"code"`
	StartIndex int           `json:"startIndex"`
	EndIndex   int           `json:"endIndex"`
	Metadata   *FileMetadata `json:"metadata,omitempty"`
}

// yamlFrontMatterPattern matches YAML front matter at the start of content.
// It matches: ---\n(yaml content)\n---\n(rest of content)
var yamlFrontMatterPattern = regexp.MustCompile(`(?s)^---\r?\n(.+?)\r?\n---\r?\n(.*)$`)

// ParseMetadataFromContent extracts YAML front matter from file content.
// Returns the parsed metadata (or nil if not found), the clean content without
// the front matter, and any parsing error.
func ParseMetadataFromContent(content string) (*FileMetadata, string, error) {
	// Trim leading whitespace/newlines before checking for front matter
	trimmedContent := strings.TrimLeft(content, " \t\r\n")

	matches := yamlFrontMatterPattern.FindStringSubmatch(trimmedContent)
	if len(matches) != 3 {
		// No front matter found, return original content
		return nil, content, nil
	}

	yamlContent := matches[1]
	fileContent := matches[2]

	var metadata FileMetadata
	if err := yaml.Unmarshal([]byte(yamlContent), &metadata); err != nil {
		// Invalid YAML - return original content without error
		// (allows graceful degradation if metadata is malformed)
		return nil, content, nil
	}

	// Only return metadata if at least one field is populated
	if metadata.IsEmpty() {
		return nil, content, nil
	}

	return &metadata, fileContent, nil
}

// ExtractCodeBlocksWithMetadata parses markdown and extracts code blocks with metadata.
// It first uses the standard ExtractCodeBlocks function, then parses YAML front matter
// from each block's content.
func ExtractCodeBlocksWithMetadata(markdown string) []CodeBlockWithMetadata {
	blocks := ExtractCodeBlocks(markdown)
	result := make([]CodeBlockWithMetadata, 0, len(blocks))

	for _, block := range blocks {
		metadata, cleanContent, _ := ParseMetadataFromContent(block.Code)

		result = append(result, CodeBlockWithMetadata{
			Language:   block.Language,
			Filename:   block.Filename,
			Code:       cleanContent,
			StartIndex: block.StartIndex,
			EndIndex:   block.EndIndex,
			Metadata:   metadata,
		})
	}

	return result
}
