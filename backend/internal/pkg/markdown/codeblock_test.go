package markdown

import (
	"testing"
)

func TestExtractCodeBlocks_NoCodeBlocks(t *testing.T) {
	content := "This is plain text without any code blocks."
	blocks := ExtractCodeBlocks(content)

	if len(blocks) != 0 {
		t.Errorf("expected 0 code blocks, got %d", len(blocks))
	}
}

func TestExtractCodeBlocks_SingleBlock(t *testing.T) {
	content := `Here is some code:

` + "```python" + `
def hello():
    print("Hello, World!")
` + "```" + `

That's all!`

	blocks := ExtractCodeBlocks(content)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 code block, got %d", len(blocks))
	}

	block := blocks[0]
	if block.Language != "python" {
		t.Errorf("expected language 'python', got '%s'", block.Language)
	}

	expectedCode := `def hello():
    print("Hello, World!")`
	if block.Code != expectedCode {
		t.Errorf("expected code:\n%s\n\ngot:\n%s", expectedCode, block.Code)
	}
}

func TestExtractCodeBlocks_MultipleBlocks(t *testing.T) {
	content := `First, create the HTML:

` + "```html" + `
<div id="app"></div>
` + "```" + `

Then add the JavaScript:

` + "```javascript" + `
const app = document.getElementById('app');
app.textContent = 'Hello';
` + "```" + `

And finally some CSS:

` + "```css" + `
#app {
    color: blue;
}
` + "```"

	blocks := ExtractCodeBlocks(content)

	if len(blocks) != 3 {
		t.Fatalf("expected 3 code blocks, got %d", len(blocks))
	}

	if blocks[0].Language != "html" {
		t.Errorf("expected first block language 'html', got '%s'", blocks[0].Language)
	}
	if blocks[1].Language != "javascript" {
		t.Errorf("expected second block language 'javascript', got '%s'", blocks[1].Language)
	}
	if blocks[2].Language != "css" {
		t.Errorf("expected third block language 'css', got '%s'", blocks[2].Language)
	}
}

func TestExtractCodeBlocks_NoLanguage(t *testing.T) {
	content := "```\nsome code without language\n```"

	blocks := ExtractCodeBlocks(content)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 code block, got %d", len(blocks))
	}

	if blocks[0].Language != "" {
		t.Errorf("expected empty language, got '%s'", blocks[0].Language)
	}
}

func TestExtractCodeBlocks_IndicesCorrect(t *testing.T) {
	content := "Start```python\ncode\n```End"

	blocks := ExtractCodeBlocks(content)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 code block, got %d", len(blocks))
	}

	block := blocks[0]
	// The start index should be where the code block starts (after "Start")
	if block.StartIndex != 5 {
		t.Errorf("expected StartIndex 5, got %d", block.StartIndex)
	}

	// The end index should be where the code block ends (before "End")
	expectedEnd := 5 + len("```python\ncode\n```")
	if block.EndIndex != expectedEnd {
		t.Errorf("expected EndIndex %d, got %d", expectedEnd, block.EndIndex)
	}
}

func TestExtractCodeBlocks_EmptyCodeBlock(t *testing.T) {
	content := "```python\n```"

	blocks := ExtractCodeBlocks(content)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 code block, got %d", len(blocks))
	}

	if blocks[0].Code != "" {
		t.Errorf("expected empty code, got '%s'", blocks[0].Code)
	}
}

func TestExtractCodeBlocks_NestedBackticks(t *testing.T) {
	content := "```python\nprint(\"Use `backticks` for inline code\")\n```"

	blocks := ExtractCodeBlocks(content)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 code block, got %d", len(blocks))
	}

	expectedCode := "print(\"Use `backticks` for inline code\")"
	if blocks[0].Code != expectedCode {
		t.Errorf("expected code:\n%s\n\ngot:\n%s", expectedCode, blocks[0].Code)
	}
}
