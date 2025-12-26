package service

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/service/prompts"
)

// DiscoveryService orchestrates the discovery flow.
type DiscoveryService struct {
	repo          repository.DiscoveryRepository
	projectRepo   repository.ProjectRepository
	promptBuilder *prompts.DiscoveryPromptBuilder
	logger        zerolog.Logger
}

// NewDiscoveryService creates a new DiscoveryService.
func NewDiscoveryService(repo repository.DiscoveryRepository, projectRepo repository.ProjectRepository, logger zerolog.Logger) *DiscoveryService {
	return &DiscoveryService{
		repo:          repo,
		projectRepo:   projectRepo,
		promptBuilder: prompts.NewDiscoveryPromptBuilder(),
		logger:        logger,
	}
}

// ErrInvalidStageTransition is returned when attempting an invalid stage transition.
var ErrInvalidStageTransition = errors.New("invalid stage transition")

// ErrDiscoveryNotFound is returned when discovery is not found.
var ErrDiscoveryNotFound = errors.New("discovery not found")

// ErrDiscoveryAlreadyComplete is returned when trying to modify a completed discovery.
var ErrDiscoveryAlreadyComplete = errors.New("discovery is already complete")

// GetOrCreateDiscovery returns an existing discovery for the project or creates a new one.
func (s *DiscoveryService) GetOrCreateDiscovery(ctx context.Context, projectID uuid.UUID) (*model.ProjectDiscovery, error) {
	s.logger.Debug().
		Str("projectId", projectID.String()).
		Msg("getting or creating discovery")

	// Try to get existing discovery
	discovery, err := s.repo.GetByProjectID(ctx, projectID)
	if err == nil {
		return discovery, nil
	}

	// If not found, create a new one
	if errors.Is(err, repository.ErrNotFound) {
		s.logger.Info().
			Str("projectId", projectID.String()).
			Msg("creating new discovery for project")

		discovery, err = s.repo.Create(ctx, projectID)
		if err != nil {
			return nil, err
		}
		return discovery, nil
	}

	return nil, err
}

// GetDiscovery retrieves a discovery by project ID.
func (s *DiscoveryService) GetDiscovery(ctx context.Context, projectID uuid.UUID) (*model.ProjectDiscovery, error) {
	discovery, err := s.repo.GetByProjectID(ctx, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrDiscoveryNotFound
		}
		return nil, err
	}
	return discovery, nil
}

// GetDiscoveryByID retrieves a discovery by its ID.
func (s *DiscoveryService) GetDiscoveryByID(ctx context.Context, discoveryID uuid.UUID) (*model.ProjectDiscovery, error) {
	discovery, err := s.repo.GetByID(ctx, discoveryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrDiscoveryNotFound
		}
		return nil, err
	}
	return discovery, nil
}

// AdvanceStage moves the discovery to the next stage.
func (s *DiscoveryService) AdvanceStage(ctx context.Context, discoveryID uuid.UUID) (*model.ProjectDiscovery, error) {
	discovery, err := s.repo.GetByID(ctx, discoveryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrDiscoveryNotFound
		}
		return nil, err
	}

	// Check if already complete
	if discovery.Stage.IsComplete() {
		return nil, ErrDiscoveryAlreadyComplete
	}

	// Get next stage
	nextStage := discovery.Stage.NextStage()
	if nextStage == "" {
		return nil, ErrInvalidStageTransition
	}

	s.logger.Info().
		Str("discoveryId", discoveryID.String()).
		Str("currentStage", string(discovery.Stage)).
		Str("nextStage", string(nextStage)).
		Msg("advancing discovery stage")

	return s.repo.UpdateStage(ctx, discoveryID, nextStage)
}

// DiscoveryDataUpdate contains fields that can be updated on a discovery.
type DiscoveryDataUpdate struct {
	BusinessContext  *string
	ProblemStatement *string
	Goals            []string
	ProjectName      *string
	SolvesStatement  *string
}

// UpdateDiscoveryData updates the discovery data fields.
func (s *DiscoveryService) UpdateDiscoveryData(ctx context.Context, discoveryID uuid.UUID, data *DiscoveryDataUpdate) error {
	discovery, err := s.repo.GetByID(ctx, discoveryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrDiscoveryNotFound
		}
		return err
	}

	// Check if already complete
	if discovery.Stage.IsComplete() {
		return ErrDiscoveryAlreadyComplete
	}

	// Update fields if provided
	if data.BusinessContext != nil {
		discovery.BusinessContext = data.BusinessContext
	}
	if data.ProblemStatement != nil {
		discovery.ProblemStatement = data.ProblemStatement
	}
	if data.Goals != nil {
		if err := discovery.SetGoals(data.Goals); err != nil {
			return err
		}
	}
	if data.ProjectName != nil {
		discovery.ProjectName = data.ProjectName
	}
	if data.SolvesStatement != nil {
		discovery.SolvesStatement = data.SolvesStatement
	}

	_, err = s.repo.Update(ctx, discovery)
	return err
}

// AddUser adds a user persona to the discovery.
func (s *DiscoveryService) AddUser(ctx context.Context, discoveryID uuid.UUID, user *model.DiscoveryUser) (*model.DiscoveryUser, error) {
	// Verify discovery exists and is not complete
	discovery, err := s.repo.GetByID(ctx, discoveryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrDiscoveryNotFound
		}
		return nil, err
	}

	if discovery.Stage.IsComplete() {
		return nil, ErrDiscoveryAlreadyComplete
	}

	user.DiscoveryID = discoveryID
	return s.repo.AddUser(ctx, user)
}

// AddFeature adds a feature to the discovery.
func (s *DiscoveryService) AddFeature(ctx context.Context, discoveryID uuid.UUID, feature *model.DiscoveryFeature) (*model.DiscoveryFeature, error) {
	// Verify discovery exists and is not complete
	discovery, err := s.repo.GetByID(ctx, discoveryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrDiscoveryNotFound
		}
		return nil, err
	}

	if discovery.Stage.IsComplete() {
		return nil, ErrDiscoveryAlreadyComplete
	}

	feature.DiscoveryID = discoveryID
	return s.repo.AddFeature(ctx, feature)
}

// GetSummary returns the complete discovery summary.
func (s *DiscoveryService) GetSummary(ctx context.Context, discoveryID uuid.UUID) (*model.DiscoverySummary, error) {
	return s.repo.GetSummary(ctx, discoveryID)
}

// ConfirmDiscovery marks the discovery as complete.
func (s *DiscoveryService) ConfirmDiscovery(ctx context.Context, discoveryID uuid.UUID) (*model.ProjectDiscovery, error) {
	discovery, err := s.repo.GetByID(ctx, discoveryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrDiscoveryNotFound
		}
		return nil, err
	}

	if discovery.Stage.IsComplete() {
		return nil, ErrDiscoveryAlreadyComplete
	}

	// Must be in summary stage to confirm
	if discovery.Stage != model.StageSummary {
		return nil, ErrInvalidStageTransition
	}

	s.logger.Info().
		Str("discoveryId", discoveryID.String()).
		Msg("confirming discovery")

	return s.repo.MarkComplete(ctx, discoveryID)
}

// ResetDiscovery deletes the discovery and creates a new one for the same project.
func (s *DiscoveryService) ResetDiscovery(ctx context.Context, discoveryID uuid.UUID) (*model.ProjectDiscovery, error) {
	discovery, err := s.repo.GetByID(ctx, discoveryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrDiscoveryNotFound
		}
		return nil, err
	}

	projectID := discovery.ProjectID

	s.logger.Info().
		Str("discoveryId", discoveryID.String()).
		Str("projectId", projectID.String()).
		Msg("resetting discovery")

	// Delete the existing discovery
	if err := s.repo.Delete(ctx, discoveryID); err != nil {
		return nil, err
	}

	// Create a new discovery for the project
	return s.repo.Create(ctx, projectID)
}

// GetSystemPrompt returns the stage-appropriate system prompt for chat integration.
func (s *DiscoveryService) GetSystemPrompt(ctx context.Context, projectID uuid.UUID) (string, error) {
	discovery, err := s.repo.GetByProjectID(ctx, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// No discovery exists - return empty string (use default prompt)
			return "", nil
		}
		return "", err
	}

	// If discovery is complete, return empty string (use default prompt)
	if discovery.Stage.IsComplete() {
		return "", nil
	}

	// Build context for the prompt
	promptContext := s.buildPromptContext(ctx, discovery)

	return s.promptBuilder.Build(discovery.Stage, promptContext), nil
}

// buildPromptContext creates a DiscoveryContext from the current discovery state.
func (s *DiscoveryService) buildPromptContext(ctx context.Context, discovery *model.ProjectDiscovery) *prompts.DiscoveryContext {
	promptCtx := &prompts.DiscoveryContext{}

	// Copy basic fields
	if discovery.BusinessContext != nil {
		promptCtx.BusinessContext = *discovery.BusinessContext
	}
	if discovery.ProblemStatement != nil {
		promptCtx.ProblemStatement = *discovery.ProblemStatement
	}
	if goals, err := discovery.Goals(); err == nil {
		promptCtx.Goals = goals
	}
	if discovery.ProjectName != nil {
		promptCtx.ProjectName = *discovery.ProjectName
	}
	if discovery.SolvesStatement != nil {
		promptCtx.SolvesStatement = *discovery.SolvesStatement
	}
	if discovery.IsReturningUser != nil {
		promptCtx.IsReturningUser = *discovery.IsReturningUser
	}

	// Get users
	users, err := s.repo.GetUsers(ctx, discovery.ID)
	if err == nil {
		promptCtx.Users = users
	}

	// Get MVP features
	mvpFeatures, err := s.repo.GetMVPFeatures(ctx, discovery.ID)
	if err == nil {
		promptCtx.MVPFeatures = mvpFeatures
	}

	// Get future features
	futureFeatures, err := s.repo.GetFutureFeatures(ctx, discovery.ID)
	if err == nil {
		promptCtx.FutureFeatures = futureFeatures
	}

	return promptCtx
}

// DiscoveryMetadata represents metadata extracted from Claude's response.
type DiscoveryMetadata struct {
	StageComplete bool                   `json:"stage_complete"`
	Extracted     map[string]interface{} `json:"extracted"`
}

// discoveryDataRegex matches the metadata comment in Claude's response.
var discoveryDataRegex = regexp.MustCompile(`<!--DISCOVERY_DATA:(.+?)-->`)

// ExtractAndSaveData extracts structured data from Claude's response and saves it.
func (s *DiscoveryService) ExtractAndSaveData(ctx context.Context, discoveryID uuid.UUID, response string) error {
	discovery, err := s.repo.GetByID(ctx, discoveryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrDiscoveryNotFound
		}
		return err
	}

	// Extract metadata from response
	metadata, err := s.parseResponseMetadata(response)
	if err != nil {
		s.logger.Warn().
			Err(err).
			Str("discoveryId", discoveryID.String()).
			Msg("failed to parse discovery metadata from response")
		return nil // Not a fatal error - metadata is optional
	}

	if metadata == nil {
		return nil // No metadata found
	}

	s.logger.Debug().
		Str("discoveryId", discoveryID.String()).
		Bool("stageComplete", metadata.StageComplete).
		Msg("extracted discovery metadata from response")

	// Process extracted data based on current stage
	if err := s.processExtractedData(ctx, discovery, metadata.Extracted); err != nil {
		s.logger.Warn().
			Err(err).
			Str("discoveryId", discoveryID.String()).
			Msg("failed to process extracted data")
		return err
	}

	// If stage is complete, advance to next stage
	if metadata.StageComplete && !discovery.Stage.IsComplete() {
		nextStage := discovery.Stage.NextStage()
		if nextStage != "" {
			if _, err := s.repo.UpdateStage(ctx, discoveryID, nextStage); err != nil {
				s.logger.Warn().
					Err(err).
					Str("discoveryId", discoveryID.String()).
					Msg("failed to advance stage")
				return err
			}

			s.logger.Info().
				Str("discoveryId", discoveryID.String()).
				Str("newStage", string(nextStage)).
				Msg("advanced discovery stage based on metadata")

			// When discovery completes, rename the project using discovered name
			if nextStage == model.StageComplete && s.projectRepo != nil {
				s.renameProjectFromDiscovery(ctx, discovery)
			}
		}
	}

	return nil
}

// renameProjectFromDiscovery renames the project using the name discovered during the flow.
func (s *DiscoveryService) renameProjectFromDiscovery(ctx context.Context, discovery *model.ProjectDiscovery) {
	// Refresh discovery to get latest project_name
	updatedDiscovery, err := s.repo.GetByID(ctx, discovery.ID)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get discovery for project rename")
		return
	}

	if updatedDiscovery.ProjectName == nil || *updatedDiscovery.ProjectName == "" {
		s.logger.Debug().Msg("no project name in discovery, skipping rename")
		return
	}

	projectName := *updatedDiscovery.ProjectName

	// Update the project title
	if err := s.projectRepo.Update(ctx, discovery.ProjectID, projectName); err != nil {
		s.logger.Warn().
			Err(err).
			Str("projectId", discovery.ProjectID.String()).
			Str("newName", projectName).
			Msg("failed to rename project after discovery")
		return
	}

	s.logger.Info().
		Str("projectId", discovery.ProjectID.String()).
		Str("newName", projectName).
		Msg("renamed project after discovery complete")
}

// parseResponseMetadata extracts the metadata JSON from Claude's response.
func (s *DiscoveryService) parseResponseMetadata(response string) (*DiscoveryMetadata, error) {
	matches := discoveryDataRegex.FindStringSubmatch(response)
	if len(matches) < 2 {
		return nil, nil // No metadata found
	}

	jsonData := strings.TrimSpace(matches[1])
	var metadata DiscoveryMetadata
	if err := json.Unmarshal([]byte(jsonData), &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// StripMetadata removes the metadata comment from Claude's response for display.
func StripMetadata(response string) string {
	return discoveryDataRegex.ReplaceAllString(response, "")
}

// processExtractedData saves extracted data fields to the discovery.
func (s *DiscoveryService) processExtractedData(ctx context.Context, discovery *model.ProjectDiscovery, extracted map[string]interface{}) error {
	if extracted == nil {
		return nil
	}

	update := &DiscoveryDataUpdate{}
	hasUpdates := false

	// Extract business context
	if bc, ok := extracted["business_context"].(string); ok && bc != "" {
		update.BusinessContext = &bc
		hasUpdates = true
	}

	// Extract problem statement
	if ps, ok := extracted["problem_statement"].(string); ok && ps != "" {
		update.ProblemStatement = &ps
		hasUpdates = true
	}

	// Extract goals
	if goalsRaw, ok := extracted["goals"].([]interface{}); ok {
		goals := make([]string, 0, len(goalsRaw))
		for _, g := range goalsRaw {
			if gs, ok := g.(string); ok {
				goals = append(goals, gs)
			}
		}
		if len(goals) > 0 {
			update.Goals = goals
			hasUpdates = true
		}
	}

	// Extract project name (check top-level and nested under "summary")
	if pn, ok := extracted["project_name"].(string); ok && pn != "" {
		update.ProjectName = &pn
		hasUpdates = true
	} else if summary, ok := extracted["summary"].(map[string]interface{}); ok {
		if pn, ok := summary["project_name"].(string); ok && pn != "" {
			update.ProjectName = &pn
			hasUpdates = true
		}
	}

	// Extract solves statement (check top-level and nested under "summary")
	if ss, ok := extracted["solves_statement"].(string); ok && ss != "" {
		update.SolvesStatement = &ss
		hasUpdates = true
	} else if summary, ok := extracted["summary"].(map[string]interface{}); ok {
		if ss, ok := summary["solves_statement"].(string); ok && ss != "" {
			update.SolvesStatement = &ss
			hasUpdates = true
		}
	}

	// Update discovery data
	if hasUpdates {
		if err := s.UpdateDiscoveryData(ctx, discovery.ID, update); err != nil {
			return err
		}
	}

	// Extract and save users
	if usersRaw, ok := extracted["users"].([]interface{}); ok {
		for _, u := range usersRaw {
			if userMap, ok := u.(map[string]interface{}); ok {
				user := &model.DiscoveryUser{
					DiscoveryID: discovery.ID,
				}
				if desc, ok := userMap["description"].(string); ok {
					user.Description = desc
				}
				if count, ok := userMap["count"].(float64); ok {
					user.UserCount = int(count)
				}
				if hasPerms, ok := userMap["has_permissions"].(bool); ok {
					user.HasPermissions = hasPerms
				}
				if notes, ok := userMap["permission_notes"].(string); ok {
					user.PermissionNotes = &notes
				}

				if user.Description != "" {
					if _, err := s.repo.AddUser(ctx, user); err != nil {
						s.logger.Warn().
							Err(err).
							Str("description", user.Description).
							Msg("failed to add extracted user")
					}
				}
			}
		}
	}

	// Extract and save MVP features
	if mvpRaw, ok := extracted["mvp_features"].([]interface{}); ok {
		for i, f := range mvpRaw {
			if featureMap, ok := f.(map[string]interface{}); ok {
				feature := &model.DiscoveryFeature{
					DiscoveryID: discovery.ID,
					Version:     "v1",
					Priority:    i + 1,
				}
				if name, ok := featureMap["name"].(string); ok {
					feature.Name = name
				}
				if priority, ok := featureMap["priority"].(float64); ok {
					feature.Priority = int(priority)
				}

				if feature.Name != "" {
					if _, err := s.repo.AddFeature(ctx, feature); err != nil {
						s.logger.Warn().
							Err(err).
							Str("name", feature.Name).
							Msg("failed to add extracted MVP feature")
					}
				}
			}
		}
	}

	// Extract and save future features
	if futureRaw, ok := extracted["future_features"].([]interface{}); ok {
		for i, f := range futureRaw {
			if featureMap, ok := f.(map[string]interface{}); ok {
				feature := &model.DiscoveryFeature{
					DiscoveryID: discovery.ID,
					Version:     "v2",
					Priority:    i + 1,
				}
				if name, ok := featureMap["name"].(string); ok {
					feature.Name = name
				}
				if version, ok := featureMap["version"].(string); ok {
					feature.Version = version
				}

				if feature.Name != "" {
					if _, err := s.repo.AddFeature(ctx, feature); err != nil {
						s.logger.Warn().
							Err(err).
							Str("name", feature.Name).
							Msg("failed to add extracted future feature")
					}
				}
			}
		}
	}

	return nil
}

// IsDiscoveryMode returns true if the project is in discovery mode.
func (s *DiscoveryService) IsDiscoveryMode(ctx context.Context, projectID uuid.UUID) (bool, error) {
	discovery, err := s.repo.GetByProjectID(ctx, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return !discovery.Stage.IsComplete(), nil
}

// GetDiscoveryStage returns the current discovery stage for a project.
func (s *DiscoveryService) GetDiscoveryStage(ctx context.Context, projectID uuid.UUID) (model.DiscoveryStage, error) {
	discovery, err := s.repo.GetByProjectID(ctx, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrDiscoveryNotFound
		}
		return "", err
	}
	return discovery.Stage, nil
}
