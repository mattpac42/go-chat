package model

import (
	"time"

	"github.com/google/uuid"
)

// ReportStatus represents the overall status of a completeness check.
type ReportStatus string

const (
	StatusPass     ReportStatus = "pass"
	StatusWarning  ReportStatus = "warning"
	StatusCritical ReportStatus = "critical"
)

// Severity represents the severity of a completeness issue.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityWarning  Severity = "warning"
	SeverityInfo     Severity = "info"
)

// CompletenessReport represents the result of a completeness check.
type CompletenessReport struct {
	ProjectID    uuid.UUID           `json:"projectId"`
	CheckedAt    time.Time           `json:"checkedAt"`
	Status       ReportStatus        `json:"status"`
	Issues       []CompletenessIssue `json:"issues"`
	FilesChecked int                 `json:"filesChecked"`
	AutoFixable  int                 `json:"autoFixable"`
}

// CompletenessIssue represents a single issue found during completeness check.
type CompletenessIssue struct {
	ID            string   `json:"id"`
	Severity      Severity `json:"severity"`
	Type          string   `json:"type"` // "missing_file", "syntax_error", "broken_reference"
	MissingFile   string   `json:"missingFile,omitempty"`
	ReferencedBy  string   `json:"referencedBy,omitempty"`
	ReferenceType string   `json:"referenceType,omitempty"` // "script", "stylesheet", "import", "image"
	LineNumber    int      `json:"lineNumber,omitempty"`
	Context       string   `json:"context,omitempty"`
	AutoFixable   bool     `json:"autoFixable"`
	FixApplied    bool     `json:"fixApplied"`
}

// HasCriticalIssues returns true if there are any critical issues.
func (r *CompletenessReport) HasCriticalIssues() bool {
	return r.Status == StatusCritical
}

// GetCriticalIssues returns only critical issues.
func (r *CompletenessReport) GetCriticalIssues() []CompletenessIssue {
	var critical []CompletenessIssue
	for _, issue := range r.Issues {
		if issue.Severity == SeverityCritical {
			critical = append(critical, issue)
		}
	}
	return critical
}

// GetMissingFiles returns a list of missing file paths.
func (r *CompletenessReport) GetMissingFiles() []string {
	seen := make(map[string]bool)
	var missing []string
	for _, issue := range r.Issues {
		if issue.Type == "missing_file" && !seen[issue.MissingFile] {
			seen[issue.MissingFile] = true
			missing = append(missing, issue.MissingFile)
		}
	}
	return missing
}

// CompletenessFixRequest represents a request to fix completeness issues.
type CompletenessFixRequest struct {
	IssueIDs []string `json:"issueIds,omitempty"` // If empty, fix all auto-fixable
}

// CompletenessFixResponse represents the result of a fix attempt.
type CompletenessFixResponse struct {
	Fixed     []string            `json:"fixed"`
	Failed    []string            `json:"failed"`
	NewReport *CompletenessReport `json:"newReport"`
}
