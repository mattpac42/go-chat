package model

import (
	"time"

	"github.com/google/uuid"
)

// Role represents the sender of a message.
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message represents a chat message in a project.
type Message struct {
	ID         uuid.UUID   `db:"id" json:"id"`
	ProjectID  uuid.UUID   `db:"project_id" json:"projectId,omitempty"`
	Role       Role        `db:"role" json:"role"`
	Content    string      `db:"content" json:"content"`
	AgentType  *string     `db:"agent_type" json:"agentType,omitempty"` // "product_manager", "designer", "developer", or null for user messages
	CreatedAt  time.Time   `db:"created_at" json:"createdAt"`
	CodeBlocks []CodeBlock `db:"-" json:"codeBlocks,omitempty"`
}

// CodeBlock represents a code block extracted from a message.
type CodeBlock struct {
	Language   string `json:"language"`
	Code       string `json:"code"`
	StartIndex int    `json:"startIndex"`
	EndIndex   int    `json:"endIndex"`
}

// CreateMessageRequest represents the request body for creating a message.
type CreateMessageRequest struct {
	ProjectID uuid.UUID `json:"projectId"`
	Role      Role      `json:"role"`
	Content   string    `json:"content"`
}
