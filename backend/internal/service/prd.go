package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

// PRDService manages PRD lifecycle and generation.
type PRDService struct {
	prdRepo       PRDRepository
	discoveryRepo repository.DiscoveryRepository
	claudeService ClaudeMessenger
	logger        zerolog.Logger
}

// PRDRepository defines the interface for PRD data access.
type PRDRepository interface {
	// CRUD
	Create(ctx context.Context, prd *model.PRD) (*model.PRD, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.PRD, error)
	Update(ctx context.Context, prd *model.PRD) (*model.PRD, error)
	Delete(ctx context.Context, id uuid.UUID) error

	// Queries
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.PRD, error)
	GetByDiscoveryID(ctx context.Context, discoveryID uuid.UUID) ([]model.PRD, error)
	GetByFeatureID(ctx context.Context, featureID uuid.UUID) (*model.PRD, error)
	GetByStatus(ctx context.Context, projectID uuid.UUID, status model.PRDStatus) ([]model.PRD, error)
	GetByVersion(ctx context.Context, projectID uuid.UUID, version string) ([]model.PRD, error)

	// Status Operations
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.PRDStatus) error
	IncrementGenerationAttempts(ctx context.Context, id uuid.UUID) error
	SetLastError(ctx context.Context, id uuid.UUID, err string) error

	// Timestamps
	SetGeneratedAt(ctx context.Context, id uuid.UUID) error
	SetApprovedAt(ctx context.Context, id uuid.UUID) error
	SetStartedAt(ctx context.Context, id uuid.UUID) error
	SetCompletedAt(ctx context.Context, id uuid.UUID) error
}

// PRD generation errors
var (
	ErrPRDNotFound           = errors.New("prd not found")
	ErrFeatureNotFound       = errors.New("feature not found")
	ErrPRDGenerationFailed   = errors.New("prd generation failed")
	ErrInvalidStatusChange   = errors.New("invalid status change")
	ErrNoActivePRD           = errors.New("no active prd")
	ErrNoReadyPRD            = errors.New("no ready prd available")
	ErrPRDAlreadyExists      = errors.New("prd already exists for this feature")
	ErrMaxRetriesExceeded    = errors.New("maximum generation retries exceeded")
)

const (
	maxGenerationAttempts = 3
)

// NewPRDService creates a new PRDService.
func NewPRDService(
	prdRepo PRDRepository,
	discoveryRepo repository.DiscoveryRepository,
	claudeService ClaudeMessenger,
	logger zerolog.Logger,
) *PRDService {
	return &PRDService{
		prdRepo:       prdRepo,
		discoveryRepo: discoveryRepo,
		claudeService: claudeService,
		logger:        logger,
	}
}

// GenerateAllPRDs creates PRDs for all features in a discovery.
// MVP features are processed in parallel, future features sequentially.
func (s *PRDService) GenerateAllPRDs(ctx context.Context, discoveryID uuid.UUID) error {
	s.logger.Info().
		Str("discoveryId", discoveryID.String()).
		Msg("starting PRD generation for all features")

	// Get discovery summary
	summary, err := s.discoveryRepo.GetSummary(ctx, discoveryID)
	if err != nil {
		return fmt.Errorf("failed to get discovery summary: %w", err)
	}

	// Get the discovery to get project ID
	discovery, err := s.discoveryRepo.GetByID(ctx, discoveryID)
	if err != nil {
		return fmt.Errorf("failed to get discovery: %w", err)
	}

	// Create pending PRD records for all features
	var mvpPRDIDs []uuid.UUID
	var futurePRDIDs []uuid.UUID

	// Create PRDs for MVP features
	for _, feature := range summary.MVPFeatures {
		prd, err := s.createPendingPRD(ctx, discovery.ProjectID, discoveryID, &feature)
		if err != nil {
			if errors.Is(err, ErrPRDAlreadyExists) {
				s.logger.Debug().
					Str("featureId", feature.ID.String()).
					Msg("PRD already exists for feature, skipping")
				continue
			}
			s.logger.Error().Err(err).
				Str("featureId", feature.ID.String()).
				Msg("failed to create pending PRD")
			continue
		}
		mvpPRDIDs = append(mvpPRDIDs, prd.ID)
	}

	// Create PRDs for future features
	for _, feature := range summary.FutureFeatures {
		prd, err := s.createPendingPRD(ctx, discovery.ProjectID, discoveryID, &feature)
		if err != nil {
			if errors.Is(err, ErrPRDAlreadyExists) {
				continue
			}
			s.logger.Error().Err(err).
				Str("featureId", feature.ID.String()).
				Msg("failed to create pending PRD")
			continue
		}
		futurePRDIDs = append(futurePRDIDs, prd.ID)
	}

	// Generate MVP PRDs in parallel
	errChan := make(chan error, len(mvpPRDIDs))
	for _, prdID := range mvpPRDIDs {
		go func(id uuid.UUID) {
			_, genErr := s.GeneratePRD(ctx, id)
			errChan <- genErr
		}(prdID)
	}

	// Collect MVP generation results
	var genErrors []error
	for range mvpPRDIDs {
		if genErr := <-errChan; genErr != nil {
			genErrors = append(genErrors, genErr)
		}
	}

	// Generate future PRDs sequentially (lower priority)
	for _, prdID := range futurePRDIDs {
		if _, genErr := s.GeneratePRD(ctx, prdID); genErr != nil {
			genErrors = append(genErrors, genErr)
		}
	}

	if len(genErrors) > 0 {
		s.logger.Warn().
			Int("errorCount", len(genErrors)).
			Msg("some PRDs failed to generate")
	}

	s.logger.Info().
		Str("discoveryId", discoveryID.String()).
		Int("mvpCount", len(mvpPRDIDs)).
		Int("futureCount", len(futurePRDIDs)).
		Int("errors", len(genErrors)).
		Msg("completed PRD generation")

	return nil
}

// createPendingPRD creates a new PRD record in pending status.
func (s *PRDService) createPendingPRD(ctx context.Context, projectID, discoveryID uuid.UUID, feature *model.DiscoveryFeature) (*model.PRD, error) {
	// Check if PRD already exists for this feature
	existing, err := s.prdRepo.GetByFeatureID(ctx, feature.ID)
	if err == nil && existing != nil {
		return nil, ErrPRDAlreadyExists
	}

	prd := &model.PRD{
		DiscoveryID: discoveryID,
		FeatureID:   feature.ID,
		ProjectID:   projectID,
		Title:       feature.Name,
		Version:     feature.Version,
		Priority:    feature.Priority,
		Status:      model.PRDStatusPending,
	}

	// Initialize empty JSON arrays
	if err := prd.SetUserStories([]model.UserStory{}); err != nil {
		return nil, err
	}
	if err := prd.SetAcceptanceCriteria([]model.AcceptanceCriterion{}); err != nil {
		return nil, err
	}
	if err := prd.SetTechnicalNotes([]model.TechnicalNote{}); err != nil {
		return nil, err
	}

	return s.prdRepo.Create(ctx, prd)
}

// GeneratePRD generates content for a single PRD using Claude.
func (s *PRDService) GeneratePRD(ctx context.Context, prdID uuid.UUID) (*model.PRD, error) {
	prd, err := s.prdRepo.GetByID(ctx, prdID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PRD: %w", err)
	}

	// Check if already generated
	if prd.Status == model.PRDStatusDraft || prd.Status == model.PRDStatusReady {
		return prd, nil
	}

	// Check max retries
	if prd.GenerationAttempts >= maxGenerationAttempts {
		return nil, ErrMaxRetriesExceeded
	}

	// Get discovery data for context
	summary, err := s.discoveryRepo.GetSummary(ctx, prd.DiscoveryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get discovery summary: %w", err)
	}

	// Get the feature for this PRD
	feature, err := s.findFeature(summary, prd.FeatureID)
	if err != nil {
		return nil, err
	}

	// Update status to generating
	if err := s.prdRepo.UpdateStatus(ctx, prdID, model.PRDStatusGenerating); err != nil {
		return nil, err
	}
	if err := s.prdRepo.IncrementGenerationAttempts(ctx, prdID); err != nil {
		s.logger.Warn().Err(err).Msg("failed to increment generation attempts")
	}

	// Generate based on feature type
	var genErr error
	if feature.IsMVP() {
		prd, genErr = s.generateFullPRD(ctx, prd, feature, summary)
	} else {
		prd, genErr = s.generateLightweightPRD(ctx, prd, feature, summary)
	}

	if genErr != nil {
		// Mark as failed and record error
		if err := s.prdRepo.UpdateStatus(ctx, prdID, model.PRDStatusFailed); err != nil {
			s.logger.Error().Err(err).Msg("failed to update PRD status to failed")
		}
		if err := s.prdRepo.SetLastError(ctx, prdID, genErr.Error()); err != nil {
			s.logger.Error().Err(err).Msg("failed to set last error")
		}
		return nil, fmt.Errorf("%w: %v", ErrPRDGenerationFailed, genErr)
	}

	// Update status to draft and set generated timestamp
	if err := s.prdRepo.UpdateStatus(ctx, prdID, model.PRDStatusDraft); err != nil {
		return nil, err
	}
	if err := s.prdRepo.SetGeneratedAt(ctx, prdID); err != nil {
		s.logger.Warn().Err(err).Msg("failed to set generated timestamp")
	}

	// Refresh PRD from database
	return s.prdRepo.GetByID(ctx, prdID)
}

// findFeature looks up a feature by ID in the discovery summary.
func (s *PRDService) findFeature(summary *model.DiscoverySummary, featureID uuid.UUID) (*model.DiscoveryFeature, error) {
	for _, f := range summary.MVPFeatures {
		if f.ID == featureID {
			return &f, nil
		}
	}
	for _, f := range summary.FutureFeatures {
		if f.ID == featureID {
			return &f, nil
		}
	}
	return nil, ErrFeatureNotFound
}

// PRDGenerationContext contains data for the PRD generation template.
type PRDGenerationContext struct {
	ProjectName      string
	ProblemStatement string
	Users            []model.DiscoveryUser
	FeatureName      string
	Version          string
	Priority         int
	RelatedFeatures  []model.DiscoveryFeature
}

// prdGenerationPrompt is the template for generating PRDs.
const prdGenerationPrompt = `You are a Product Manager creating a PRD for a specific feature.

## Project Context
Project: {{.ProjectName}}
Problem Statement: {{.ProblemStatement}}
Target Users:
{{range .Users}}
- {{.Description}} ({{.UserCount}} users{{if .HasPermissions}}, elevated permissions{{end}})
{{end}}

## Feature to Document
Feature: {{.FeatureName}}
Version: {{.Version}}
Priority: {{.Priority}}

## Related Features (for context)
{{range .RelatedFeatures}}
- {{.Name}} ({{.Version}}, Priority {{.Priority}})
{{end}}

## Instructions
Create a PRD for this feature with the following structure:

1. **Overview** (2-3 sentences)
   - What this feature does
   - Why it matters to users

2. **User Stories** (3-5 stories in format)
   For each story, provide:
   - ID (US-001, US-002, etc.)
   - As a [user type], I want [action], so that [benefit]
   - Priority: must/should/could
   - Complexity: low/medium/high

3. **Acceptance Criteria** (Gherkin format)
   For each user story, 2-3 criteria:
   - ID (AC-001, AC-002, etc.)
   - Given [precondition]
   - When [action]
   - Then [expected result]

4. **Technical Notes** (implementation guidance)
   - Data considerations
   - UI/UX notes
   - Integration points

Output as JSON matching this structure:
{
  "overview": "string",
  "userStories": [
    {
      "id": "US-001",
      "asA": "user type",
      "iWant": "action",
      "soThat": "benefit",
      "priority": "must",
      "complexity": "low"
    }
  ],
  "acceptanceCriteria": [
    {
      "id": "AC-001",
      "given": "precondition",
      "when": "action",
      "then": "expected result",
      "userStoryId": "US-001"
    }
  ],
  "technicalNotes": [
    {
      "category": "data",
      "title": "title",
      "description": "description",
      "suggestions": ["suggestion1", "suggestion2"]
    }
  ]
}

IMPORTANT: Output ONLY the JSON object. No markdown code blocks, no additional text.`

// lightweightPRDPrompt is the template for generating lightweight PRDs for future features.
const lightweightPRDPrompt = `You are a Product Manager creating a brief PRD for a future feature.

## Project Context
Project: {{.ProjectName}}
Problem Statement: {{.ProblemStatement}}

## Feature to Document
Feature: {{.FeatureName}}
Version: {{.Version}}
Priority: {{.Priority}}

## Instructions
Create a brief overview for this future feature. Just 2-3 sentences describing:
- What this feature will do
- Why it matters to users
- Key considerations for future implementation

Output as JSON matching this structure:
{
  "overview": "string",
  "userStories": [],
  "acceptanceCriteria": [],
  "technicalNotes": []
}

IMPORTANT: Output ONLY the JSON object. No markdown code blocks, no additional text.`

// generateFullPRD generates a complete PRD for MVP features.
func (s *PRDService) generateFullPRD(ctx context.Context, prd *model.PRD, feature *model.DiscoveryFeature, summary *model.DiscoverySummary) (*model.PRD, error) {
	s.logger.Debug().
		Str("prdId", prd.ID.String()).
		Str("feature", feature.Name).
		Msg("generating full PRD")

	// Build template context
	genCtx := PRDGenerationContext{
		ProjectName:      summary.ProjectName,
		ProblemStatement: summary.SolvesStatement,
		Users:            summary.Users,
		FeatureName:      feature.Name,
		Version:          feature.Version,
		Priority:         feature.Priority,
		RelatedFeatures:  s.getRelatedFeatures(summary, feature.ID),
	}

	// Parse and execute template
	tmpl, err := template.New("prd").Parse(prdGenerationPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var promptBuf bytes.Buffer
	if err := tmpl.Execute(&promptBuf, genCtx); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Call Claude for generation
	content, err := s.callClaude(ctx, promptBuf.String())
	if err != nil {
		return nil, err
	}

	// Parse response and update PRD
	return s.parsePRDResponse(ctx, prd, content)
}

// generateLightweightPRD generates a simple overview for future features.
func (s *PRDService) generateLightweightPRD(ctx context.Context, prd *model.PRD, feature *model.DiscoveryFeature, summary *model.DiscoverySummary) (*model.PRD, error) {
	s.logger.Debug().
		Str("prdId", prd.ID.String()).
		Str("feature", feature.Name).
		Msg("generating lightweight PRD")

	// Build template context
	genCtx := PRDGenerationContext{
		ProjectName:      summary.ProjectName,
		ProblemStatement: summary.SolvesStatement,
		FeatureName:      feature.Name,
		Version:          feature.Version,
		Priority:         feature.Priority,
	}

	// Parse and execute template
	tmpl, err := template.New("prd-light").Parse(lightweightPRDPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var promptBuf bytes.Buffer
	if err := tmpl.Execute(&promptBuf, genCtx); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Call Claude for generation
	content, err := s.callClaude(ctx, promptBuf.String())
	if err != nil {
		return nil, err
	}

	// Parse response and update PRD
	return s.parsePRDResponse(ctx, prd, content)
}

// getRelatedFeatures returns features related to the current one (same version, different ID).
func (s *PRDService) getRelatedFeatures(summary *model.DiscoverySummary, excludeID uuid.UUID) []model.DiscoveryFeature {
	var related []model.DiscoveryFeature
	for _, f := range summary.MVPFeatures {
		if f.ID != excludeID {
			related = append(related, f)
		}
	}
	return related
}

// callClaude sends a message to Claude and collects the streaming response.
func (s *PRDService) callClaude(ctx context.Context, prompt string) (string, error) {
	messages := []ClaudeMessage{
		{Role: "user", Content: prompt},
	}

	stream, err := s.claudeService.SendMessage(ctx, "You are a Product Manager. Output only valid JSON.", messages)
	if err != nil {
		return "", fmt.Errorf("failed to send message to Claude: %w", err)
	}
	defer stream.Close()

	// Collect all chunks
	var content strings.Builder
	for chunk := range stream.Chunks() {
		content.WriteString(chunk)
	}

	if stream.Err() != nil {
		return "", fmt.Errorf("stream error: %w", stream.Err())
	}

	return content.String(), nil
}

// PRDResponseJSON represents the JSON structure returned by Claude.
type PRDResponseJSON struct {
	Overview           string                     `json:"overview"`
	UserStories        []model.UserStory          `json:"userStories"`
	AcceptanceCriteria []model.AcceptanceCriterion `json:"acceptanceCriteria"`
	TechnicalNotes     []model.TechnicalNote      `json:"technicalNotes"`
}

// parsePRDResponse parses Claude's response and updates the PRD.
func (s *PRDService) parsePRDResponse(ctx context.Context, prd *model.PRD, content string) (*model.PRD, error) {
	// Clean up the response - remove any markdown code blocks
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var response PRDResponseJSON
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		s.logger.Error().
			Err(err).
			Str("content", content[:min(200, len(content))]).
			Msg("failed to parse PRD response JSON")
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Update PRD fields
	prd.Overview = response.Overview

	if err := prd.SetUserStories(response.UserStories); err != nil {
		return nil, fmt.Errorf("failed to set user stories: %w", err)
	}

	if err := prd.SetAcceptanceCriteria(response.AcceptanceCriteria); err != nil {
		return nil, fmt.Errorf("failed to set acceptance criteria: %w", err)
	}

	if err := prd.SetTechnicalNotes(response.TechnicalNotes); err != nil {
		return nil, fmt.Errorf("failed to set technical notes: %w", err)
	}

	// Save to database
	return s.prdRepo.Update(ctx, prd)
}

// RetryGeneration retries PRD generation for a failed PRD.
func (s *PRDService) RetryGeneration(ctx context.Context, prdID uuid.UUID) (*model.PRD, error) {
	prd, err := s.prdRepo.GetByID(ctx, prdID)
	if err != nil {
		return nil, ErrPRDNotFound
	}

	// Only retry failed PRDs
	if prd.Status != model.PRDStatusFailed && prd.Status != model.PRDStatusPending {
		return nil, fmt.Errorf("%w: can only retry failed or pending PRDs", ErrInvalidStatusChange)
	}

	// Reset status to pending
	if err := s.prdRepo.UpdateStatus(ctx, prdID, model.PRDStatusPending); err != nil {
		return nil, err
	}

	return s.GeneratePRD(ctx, prdID)
}

// GetByID retrieves a PRD by its ID.
func (s *PRDService) GetByID(ctx context.Context, prdID uuid.UUID) (*model.PRD, error) {
	prd, err := s.prdRepo.GetByID(ctx, prdID)
	if err != nil {
		return nil, ErrPRDNotFound
	}
	return prd, nil
}

// GetByProjectID retrieves all PRDs for a project.
func (s *PRDService) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.PRD, error) {
	return s.prdRepo.GetByProjectID(ctx, projectID)
}

// GetByDiscoveryID retrieves all PRDs generated from a discovery.
func (s *PRDService) GetByDiscoveryID(ctx context.Context, discoveryID uuid.UUID) ([]model.PRD, error) {
	return s.prdRepo.GetByDiscoveryID(ctx, discoveryID)
}

// GetMVPPRDs retrieves all MVP (v1) PRDs for a project.
func (s *PRDService) GetMVPPRDs(ctx context.Context, projectID uuid.UUID) ([]model.PRD, error) {
	return s.prdRepo.GetByVersion(ctx, projectID, "v1")
}

// UpdateStatus updates the status of a PRD.
func (s *PRDService) UpdateStatus(ctx context.Context, prdID uuid.UUID, status model.PRDStatus) error {
	prd, err := s.prdRepo.GetByID(ctx, prdID)
	if err != nil {
		return ErrPRDNotFound
	}

	// Validate status transition
	if !s.isValidStatusTransition(prd.Status, status) {
		return fmt.Errorf("%w: cannot transition from %s to %s", ErrInvalidStatusChange, prd.Status, status)
	}

	return s.prdRepo.UpdateStatus(ctx, prdID, status)
}

// isValidStatusTransition checks if a status transition is valid.
func (s *PRDService) isValidStatusTransition(from, to model.PRDStatus) bool {
	validTransitions := map[model.PRDStatus][]model.PRDStatus{
		model.PRDStatusPending:    {model.PRDStatusGenerating, model.PRDStatusFailed},
		model.PRDStatusGenerating: {model.PRDStatusDraft, model.PRDStatusFailed},
		model.PRDStatusDraft:      {model.PRDStatusReady},
		model.PRDStatusReady:      {model.PRDStatusInProgress},
		model.PRDStatusInProgress: {model.PRDStatusComplete, model.PRDStatusReady}, // Can go back to ready if paused
		model.PRDStatusComplete:   {},                                               // Terminal state
		model.PRDStatusFailed:     {model.PRDStatusPending},                         // Can retry
	}

	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// MarkAsReady marks a draft PRD as ready for implementation.
func (s *PRDService) MarkAsReady(ctx context.Context, prdID uuid.UUID) error {
	if err := s.UpdateStatus(ctx, prdID, model.PRDStatusReady); err != nil {
		return err
	}
	return s.prdRepo.SetApprovedAt(ctx, prdID)
}

// StartImplementation marks a PRD as in progress.
func (s *PRDService) StartImplementation(ctx context.Context, prdID uuid.UUID) error {
	if err := s.UpdateStatus(ctx, prdID, model.PRDStatusInProgress); err != nil {
		return err
	}
	return s.prdRepo.SetStartedAt(ctx, prdID)
}

// CompleteImplementation marks a PRD as complete.
func (s *PRDService) CompleteImplementation(ctx context.Context, prdID uuid.UUID) error {
	if err := s.UpdateStatus(ctx, prdID, model.PRDStatusComplete); err != nil {
		return err
	}
	return s.prdRepo.SetCompletedAt(ctx, prdID)
}

// GetActivePRD retrieves the currently active PRD for a project.
// The active PRD is the one currently being worked on.
func (s *PRDService) GetActivePRD(ctx context.Context, projectID uuid.UUID) (*model.PRD, error) {
	prds, err := s.prdRepo.GetByStatus(ctx, projectID, model.PRDStatusInProgress)
	if err != nil {
		return nil, err
	}

	if len(prds) == 0 {
		return nil, ErrNoActivePRD
	}

	// Return the first in-progress PRD (should only be one)
	return &prds[0], nil
}

// SetActivePRD sets a PRD as the active one by starting its implementation.
func (s *PRDService) SetActivePRD(ctx context.Context, projectID uuid.UUID, prdID uuid.UUID) error {
	// Verify PRD belongs to project
	prd, err := s.prdRepo.GetByID(ctx, prdID)
	if err != nil {
		return ErrPRDNotFound
	}

	if prd.ProjectID != projectID {
		return fmt.Errorf("PRD does not belong to this project")
	}

	// PRD must be ready to become active
	if prd.Status != model.PRDStatusReady {
		return fmt.Errorf("%w: PRD must be in ready status", ErrInvalidStatusChange)
	}

	return s.StartImplementation(ctx, prdID)
}

// ClearActivePRD pauses the current active PRD, moving it back to ready.
func (s *PRDService) ClearActivePRD(ctx context.Context, projectID uuid.UUID) error {
	activePRD, err := s.GetActivePRD(ctx, projectID)
	if err != nil {
		if errors.Is(err, ErrNoActivePRD) {
			return nil // Nothing to clear
		}
		return err
	}

	// Move back to ready status
	return s.prdRepo.UpdateStatus(ctx, activePRD.ID, model.PRDStatusReady)
}

// GetNextPRD returns the next PRD that should be worked on.
// Returns the highest priority ready PRD.
func (s *PRDService) GetNextPRD(ctx context.Context, projectID uuid.UUID) (*model.PRD, error) {
	prds, err := s.prdRepo.GetByStatus(ctx, projectID, model.PRDStatusReady)
	if err != nil {
		return nil, err
	}

	if len(prds) == 0 {
		return nil, ErrNoReadyPRD
	}

	// Find highest priority (lowest number = highest priority)
	var next *model.PRD
	for i := range prds {
		if next == nil || prds[i].Priority < next.Priority {
			next = &prds[i]
		}
	}

	return next, nil
}

// UpdateOverview updates the overview section of a PRD.
func (s *PRDService) UpdateOverview(ctx context.Context, prdID uuid.UUID, overview string) error {
	prd, err := s.prdRepo.GetByID(ctx, prdID)
	if err != nil {
		return ErrPRDNotFound
	}

	// Only allow editing draft PRDs
	if prd.Status != model.PRDStatusDraft {
		return fmt.Errorf("%w: can only edit draft PRDs", ErrInvalidStatusChange)
	}

	prd.Overview = overview
	_, err = s.prdRepo.Update(ctx, prd)
	return err
}

// AddUserStory adds a user story to a PRD.
func (s *PRDService) AddUserStory(ctx context.Context, prdID uuid.UUID, story *model.UserStory) error {
	prd, err := s.prdRepo.GetByID(ctx, prdID)
	if err != nil {
		return ErrPRDNotFound
	}

	if prd.Status != model.PRDStatusDraft {
		return fmt.Errorf("%w: can only edit draft PRDs", ErrInvalidStatusChange)
	}

	stories, err := prd.UserStories()
	if err != nil {
		return err
	}

	// Generate ID if not provided
	if story.ID == "" {
		story.ID = fmt.Sprintf("US-%03d", len(stories)+1)
	}

	stories = append(stories, *story)
	if err := prd.SetUserStories(stories); err != nil {
		return err
	}

	_, err = s.prdRepo.Update(ctx, prd)
	return err
}

// UpdateUserStory updates a user story in a PRD.
func (s *PRDService) UpdateUserStory(ctx context.Context, prdID uuid.UUID, storyID string, story *model.UserStory) error {
	prd, err := s.prdRepo.GetByID(ctx, prdID)
	if err != nil {
		return ErrPRDNotFound
	}

	if prd.Status != model.PRDStatusDraft {
		return fmt.Errorf("%w: can only edit draft PRDs", ErrInvalidStatusChange)
	}

	stories, err := prd.UserStories()
	if err != nil {
		return err
	}

	found := false
	for i, s := range stories {
		if s.ID == storyID {
			story.ID = storyID // Preserve ID
			stories[i] = *story
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("user story %s not found", storyID)
	}

	if err := prd.SetUserStories(stories); err != nil {
		return err
	}

	_, err = s.prdRepo.Update(ctx, prd)
	return err
}

// DeleteUserStory removes a user story from a PRD.
func (s *PRDService) DeleteUserStory(ctx context.Context, prdID uuid.UUID, storyID string) error {
	prd, err := s.prdRepo.GetByID(ctx, prdID)
	if err != nil {
		return ErrPRDNotFound
	}

	if prd.Status != model.PRDStatusDraft {
		return fmt.Errorf("%w: can only edit draft PRDs", ErrInvalidStatusChange)
	}

	stories, err := prd.UserStories()
	if err != nil {
		return err
	}

	var newStories []model.UserStory
	found := false
	for _, s := range stories {
		if s.ID == storyID {
			found = true
			continue
		}
		newStories = append(newStories, s)
	}

	if !found {
		return fmt.Errorf("user story %s not found", storyID)
	}

	if err := prd.SetUserStories(newStories); err != nil {
		return err
	}

	_, err = s.prdRepo.Update(ctx, prd)
	return err
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
