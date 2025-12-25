package markdown

import (
	"testing"
)

func TestParseMetadataFromContent_ValidFrontMatter(t *testing.T) {
	content := `---
short_description: Main homepage structure
long_description: This HTML file defines the landing page structure.
functional_group: Homepage
---
<!DOCTYPE html>
<html>
<body>Hello</body>
</html>`

	metadata, cleanContent, err := ParseMetadataFromContent(content)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata == nil {
		t.Fatal("expected metadata to be parsed, got nil")
	}

	if metadata.ShortDescription != "Main homepage structure" {
		t.Errorf("expected short_description 'Main homepage structure', got '%s'", metadata.ShortDescription)
	}

	if metadata.LongDescription != "This HTML file defines the landing page structure." {
		t.Errorf("expected long_description, got '%s'", metadata.LongDescription)
	}

	if metadata.FunctionalGroup != "Homepage" {
		t.Errorf("expected functional_group 'Homepage', got '%s'", metadata.FunctionalGroup)
	}

	expectedContent := `<!DOCTYPE html>
<html>
<body>Hello</body>
</html>`
	if cleanContent != expectedContent {
		t.Errorf("expected clean content:\n%s\n\ngot:\n%s", expectedContent, cleanContent)
	}
}

func TestParseMetadataFromContent_NoFrontMatter(t *testing.T) {
	content := `<!DOCTYPE html>
<html>
<body>Hello</body>
</html>`

	metadata, cleanContent, err := ParseMetadataFromContent(content)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata != nil {
		t.Errorf("expected no metadata, got %+v", metadata)
	}

	if cleanContent != content {
		t.Errorf("expected original content returned, got '%s'", cleanContent)
	}
}

func TestParseMetadataFromContent_InvalidYAML(t *testing.T) {
	content := `---
short_description: [invalid yaml
long_description: unclosed bracket
---
<html></html>`

	metadata, cleanContent, err := ParseMetadataFromContent(content)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should gracefully return nil metadata for invalid YAML
	if metadata != nil {
		t.Errorf("expected nil metadata for invalid YAML, got %+v", metadata)
	}

	// Should return original content
	if cleanContent != content {
		t.Errorf("expected original content returned for invalid YAML")
	}
}

func TestParseMetadataFromContent_EmptyFrontMatter(t *testing.T) {
	content := `---
---
<html></html>`

	metadata, cleanContent, err := ParseMetadataFromContent(content)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty metadata should be treated as no metadata
	if metadata != nil {
		t.Errorf("expected nil metadata for empty front matter, got %+v", metadata)
	}

	if cleanContent != content {
		t.Errorf("expected original content returned for empty front matter")
	}
}

func TestParseMetadataFromContent_PartialMetadata(t *testing.T) {
	content := `---
short_description: Just a short description
---
<html></html>`

	metadata, cleanContent, err := ParseMetadataFromContent(content)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata == nil {
		t.Fatal("expected metadata to be parsed for partial front matter")
	}

	if metadata.ShortDescription != "Just a short description" {
		t.Errorf("expected short_description, got '%s'", metadata.ShortDescription)
	}

	if metadata.LongDescription != "" {
		t.Errorf("expected empty long_description, got '%s'", metadata.LongDescription)
	}

	expectedContent := `<html></html>`
	if cleanContent != expectedContent {
		t.Errorf("expected clean content '%s', got '%s'", expectedContent, cleanContent)
	}
}

func TestParseMetadataFromContent_LeadingWhitespace(t *testing.T) {
	content := `
---
short_description: Description with leading whitespace
functional_group: Test
---
<html></html>`

	metadata, cleanContent, err := ParseMetadataFromContent(content)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata == nil {
		t.Fatal("expected metadata to be parsed even with leading whitespace")
	}

	if cleanContent != "<html></html>" {
		t.Errorf("expected clean content, got '%s'", cleanContent)
	}

	if metadata.ShortDescription != "Description with leading whitespace" {
		t.Errorf("expected short_description, got '%s'", metadata.ShortDescription)
	}
}

func TestParseMetadataFromContent_WindowsLineEndings(t *testing.T) {
	content := "---\r\nshort_description: Windows line endings\r\nfunctional_group: Test\r\n---\r\n<html></html>"

	metadata, cleanContent, err := ParseMetadataFromContent(content)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata == nil {
		t.Fatal("expected metadata to be parsed with Windows line endings")
	}

	if metadata.ShortDescription != "Windows line endings" {
		t.Errorf("expected short_description 'Windows line endings', got '%s'", metadata.ShortDescription)
	}

	if cleanContent != "<html></html>" {
		t.Errorf("expected clean content without front matter, got '%s'", cleanContent)
	}
}

func TestExtractCodeBlocksWithMetadata_FullExample(t *testing.T) {
	content := `Here is an HTML file:

` + "```html:index.html" + `
---
short_description: Main landing page
long_description: The homepage with navigation and hero section.
functional_group: Homepage
---
<!DOCTYPE html>
<html></html>
` + "```" + `

And here is a CSS file:

` + "```css:styles.css" + `
---
short_description: Homepage styles
functional_group: Homepage
---
body { margin: 0; }
` + "```" + `

And a file without metadata:

` + "```javascript:app.js" + `
console.log('hello');
` + "```"

	blocks := ExtractCodeBlocksWithMetadata(content)

	if len(blocks) != 3 {
		t.Fatalf("expected 3 code blocks, got %d", len(blocks))
	}

	// First block: HTML with full metadata
	if blocks[0].Language != "html" {
		t.Errorf("expected first block language 'html', got '%s'", blocks[0].Language)
	}
	if blocks[0].Filename != "index.html" {
		t.Errorf("expected first block filename 'index.html', got '%s'", blocks[0].Filename)
	}
	if blocks[0].Metadata == nil {
		t.Fatal("expected first block to have metadata")
	}
	if blocks[0].Metadata.ShortDescription != "Main landing page" {
		t.Errorf("expected short_description 'Main landing page', got '%s'", blocks[0].Metadata.ShortDescription)
	}
	if blocks[0].Metadata.FunctionalGroup != "Homepage" {
		t.Errorf("expected functional_group 'Homepage', got '%s'", blocks[0].Metadata.FunctionalGroup)
	}
	expectedCode := `<!DOCTYPE html>
<html></html>`
	if blocks[0].Code != expectedCode {
		t.Errorf("expected clean code:\n%s\n\ngot:\n%s", expectedCode, blocks[0].Code)
	}

	// Second block: CSS with partial metadata
	if blocks[1].Language != "css" {
		t.Errorf("expected second block language 'css', got '%s'", blocks[1].Language)
	}
	if blocks[1].Filename != "styles.css" {
		t.Errorf("expected second block filename 'styles.css', got '%s'", blocks[1].Filename)
	}
	if blocks[1].Metadata == nil {
		t.Fatal("expected second block to have metadata")
	}
	if blocks[1].Metadata.ShortDescription != "Homepage styles" {
		t.Errorf("expected short_description 'Homepage styles', got '%s'", blocks[1].Metadata.ShortDescription)
	}
	if blocks[1].Code != "body { margin: 0; }" {
		t.Errorf("expected clean code 'body { margin: 0; }', got '%s'", blocks[1].Code)
	}

	// Third block: JavaScript without metadata
	if blocks[2].Language != "javascript" {
		t.Errorf("expected third block language 'javascript', got '%s'", blocks[2].Language)
	}
	if blocks[2].Filename != "app.js" {
		t.Errorf("expected third block filename 'app.js', got '%s'", blocks[2].Filename)
	}
	if blocks[2].Metadata != nil {
		t.Errorf("expected third block to have no metadata, got %+v", blocks[2].Metadata)
	}
	if blocks[2].Code != "console.log('hello');" {
		t.Errorf("expected code \"console.log('hello');\", got '%s'", blocks[2].Code)
	}
}

func TestExtractCodeBlocksWithMetadata_EmptyMarkdown(t *testing.T) {
	blocks := ExtractCodeBlocksWithMetadata("")

	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks for empty markdown, got %d", len(blocks))
	}
}

func TestExtractCodeBlocksWithMetadata_NoCodeBlocks(t *testing.T) {
	content := "This is just plain text without any code blocks."
	blocks := ExtractCodeBlocksWithMetadata(content)

	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks, got %d", len(blocks))
	}
}

func TestFileMetadata_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		metadata FileMetadata
		expected bool
	}{
		{
			name:     "all empty",
			metadata: FileMetadata{},
			expected: true,
		},
		{
			name:     "only short_description",
			metadata: FileMetadata{ShortDescription: "test"},
			expected: false,
		},
		{
			name:     "only long_description",
			metadata: FileMetadata{LongDescription: "test"},
			expected: false,
		},
		{
			name:     "only functional_group",
			metadata: FileMetadata{FunctionalGroup: "test"},
			expected: false,
		},
		{
			name: "all populated",
			metadata: FileMetadata{
				ShortDescription: "short",
				LongDescription:  "long",
				FunctionalGroup:  "group",
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.metadata.IsEmpty()
			if result != tc.expected {
				t.Errorf("expected IsEmpty() = %v, got %v", tc.expected, result)
			}
		})
	}
}
