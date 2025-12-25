package model

import (
	"time"

	"github.com/google/uuid"
)

// File represents an extracted code file from an AI response.
type File struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ProjectID uuid.UUID `db:"project_id" json:"projectId"`
	Path      string    `db:"path" json:"path"`
	Filename  string    `db:"filename" json:"filename"`
	Language  string    `db:"language" json:"language,omitempty"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// FileListItem represents a file in list view (without content).
type FileListItem struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Path      string    `db:"path" json:"path"`
	Filename  string    `db:"filename" json:"filename"`
	Language  string    `db:"language" json:"language,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// ListFilesResponse represents the response for listing files.
type ListFilesResponse struct {
	Files []FileListItem `json:"files"`
}

// GetFileResponse represents the response for getting a single file.
type GetFileResponse struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"projectId"`
	Path      string    `json:"path"`
	Filename  string    `json:"filename"`
	Language  string    `json:"language,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
