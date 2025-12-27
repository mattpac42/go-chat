package model

import (
	"time"

	"github.com/google/uuid"
)

// FileSource represents the original source of an uploaded file that was converted.
type FileSource struct {
	ID               uuid.UUID `db:"id" json:"id"`
	FileID           uuid.UUID `db:"file_id" json:"fileId"`
	OriginalFilename string    `db:"original_filename" json:"originalFilename"`
	OriginalMimeType string    `db:"original_mime_type" json:"originalMimeType"`
	OriginalSizeBytes int64    `db:"original_size_bytes" json:"originalSizeBytes"`
	ConversionStatus string    `db:"conversion_status" json:"conversionStatus"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
}

// UploadResponse represents the response from a successful file upload.
type UploadResponse struct {
	File   UploadedFile   `json:"file"`
	Source UploadedSource `json:"source"`
}

// UploadedFile represents the converted file information in the upload response.
type UploadedFile struct {
	ID              uuid.UUID `json:"id"`
	Path            string    `json:"path"`
	Content         string    `json:"content"`
	ShortDescription string   `json:"shortDescription"`
	FunctionalGroup string    `json:"functionalGroup"`
}

// UploadedSource represents the original source file information in the upload response.
type UploadedSource struct {
	OriginalFilename  string `json:"originalFilename"`
	OriginalMimeType  string `json:"originalMimeType"`
	OriginalSizeBytes int64  `json:"originalSizeBytes"`
}
