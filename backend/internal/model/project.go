package model

import (
	"time"

	"github.com/google/uuid"
)

// Project represents a chat project/conversation.
type Project struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Title        string    `db:"title" json:"title"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
	MessageCount int       `db:"-" json:"messageCount,omitempty"`
	Messages     []Message `db:"-" json:"messages,omitempty"`
}

// ProjectListItem represents a project in list view.
type ProjectListItem struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Title        string    `db:"title" json:"title"`
	MessageCount int       `db:"message_count" json:"messageCount"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

// CreateProjectRequest represents the request body for creating a project.
type CreateProjectRequest struct {
	Title string `json:"title"`
}

// CreateProjectResponse represents the response after creating a project.
type CreateProjectResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ListProjectsResponse represents the response for listing projects.
type ListProjectsResponse struct {
	Projects []ProjectListItem `json:"projects"`
}

// GetProjectResponse represents the response for getting a project with messages.
type GetProjectResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Messages  []Message `json:"messages"`
}

// UpdateProjectRequest represents the request body for updating a project.
type UpdateProjectRequest struct {
	Title string `json:"title"`
}

// UpdateProjectResponse represents the response after updating a project.
type UpdateProjectResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
