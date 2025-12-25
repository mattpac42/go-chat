package model

import (
	"time"

	"github.com/google/uuid"
)

// FileMetadata represents metadata about a file for the App Map feature.
type FileMetadata struct {
	ID               uuid.UUID `db:"id" json:"id"`
	FileID           uuid.UUID `db:"file_id" json:"fileId"`
	ShortDescription string    `db:"short_description" json:"shortDescription,omitempty"`
	LongDescription  string    `db:"long_description" json:"longDescription,omitempty"`
	FunctionalGroup  string    `db:"functional_group" json:"functionalGroup,omitempty"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

// FileWithMetadata represents a file with its associated metadata for App Map display.
type FileWithMetadata struct {
	ID               uuid.UUID `db:"id" json:"id"`
	ProjectID        uuid.UUID `db:"project_id" json:"projectId"`
	Path             string    `db:"path" json:"path"`
	Filename         string    `db:"filename" json:"filename"`
	Language         string    `db:"language" json:"language,omitempty"`
	ShortDescription string    `db:"short_description" json:"shortDescription,omitempty"`
	LongDescription  string    `db:"long_description" json:"longDescription,omitempty"`
	FunctionalGroup  string    `db:"functional_group" json:"functionalGroup,omitempty"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
}

// CreateFileMetadataRequest represents the request payload for creating file metadata.
type CreateFileMetadataRequest struct {
	ShortDescription string `json:"shortDescription"`
	LongDescription  string `json:"longDescription"`
	FunctionalGroup  string `json:"functionalGroup"`
}

// UpdateFileMetadataRequest represents the request payload for updating file metadata.
type UpdateFileMetadataRequest struct {
	ShortDescription *string `json:"shortDescription,omitempty"`
	LongDescription  *string `json:"longDescription,omitempty"`
	FunctionalGroup  *string `json:"functionalGroup,omitempty"`
}

// FileMetadataResponse represents the response for file metadata operations.
type FileMetadataResponse struct {
	ID               uuid.UUID `json:"id"`
	FileID           uuid.UUID `json:"fileId"`
	ShortDescription string    `json:"shortDescription,omitempty"`
	LongDescription  string    `json:"longDescription,omitempty"`
	FunctionalGroup  string    `json:"functionalGroup,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// ListFileMetadataResponse represents the response for listing files with metadata.
type ListFileMetadataResponse struct {
	Files []FileWithMetadata `json:"files"`
}

// FunctionalGroupSummary represents a summary of files in a functional group.
type FunctionalGroupSummary struct {
	Name      string `json:"name"`
	FileCount int    `json:"fileCount"`
}

// ListFunctionalGroupsResponse represents the response for listing functional groups.
type ListFunctionalGroupsResponse struct {
	Groups []FunctionalGroupSummary `json:"groups"`
}
