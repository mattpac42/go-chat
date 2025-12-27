package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// MockClaudeVision implements ClaudeVision for testing.
type MockClaudeVision struct {
	mu               sync.RWMutex
	responses        map[string]string // mimeType -> response
	defaultResponse  string
	errorToReturn    error
	analyzeCallCount int
	lastImageData    []byte
	lastMimeType     string
	lastPrompt       string
}

// NewMockClaudeVision creates a new MockClaudeVision.
func NewMockClaudeVision() *MockClaudeVision {
	return &MockClaudeVision{
		responses: make(map[string]string),
		defaultResponse: `## Image Description

This is a mock image analysis for testing.

### Content
The image contains visual content that was analyzed.

### Extracted Text
*Sample text from the image.*
`,
	}
}

// SetResponse sets a specific response for a MIME type.
func (m *MockClaudeVision) SetResponse(mimeType, response string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses[mimeType] = response
}

// SetDefaultResponse sets the default response for all MIME types.
func (m *MockClaudeVision) SetDefaultResponse(response string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.defaultResponse = response
}

// SetError sets an error to be returned from AnalyzeImage.
func (m *MockClaudeVision) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorToReturn = err
}

// AnalyzeImage implements the ClaudeVision interface.
func (m *MockClaudeVision) AnalyzeImage(ctx context.Context, imageData []byte, mimeType, prompt string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.analyzeCallCount++
	m.lastImageData = imageData
	m.lastMimeType = mimeType
	m.lastPrompt = prompt

	if m.errorToReturn != nil {
		return "", m.errorToReturn
	}

	if response, ok := m.responses[mimeType]; ok {
		return response, nil
	}

	return m.defaultResponse, nil
}

// GetAnalyzeCallCount returns how many times AnalyzeImage was called.
func (m *MockClaudeVision) GetAnalyzeCallCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.analyzeCallCount
}

// GetLastCall returns the parameters from the last AnalyzeImage call.
func (m *MockClaudeVision) GetLastCall() (imageData []byte, mimeType, prompt string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastImageData, m.lastMimeType, m.lastPrompt
}

// Reset resets all state.
func (m *MockClaudeVision) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses = make(map[string]string)
	m.analyzeCallCount = 0
	m.lastImageData = nil
	m.lastMimeType = ""
	m.lastPrompt = ""
	m.errorToReturn = nil
}

// ErrVisionAnalysisFailed is a sample error for testing error cases.
var ErrVisionAnalysisFailed = errors.New("vision analysis failed")

// MockClaudeVisionWithError creates a mock that returns an error.
func MockClaudeVisionWithError(errMsg string) *MockClaudeVision {
	m := NewMockClaudeVision()
	m.SetError(fmt.Errorf("%s", errMsg))
	return m
}
