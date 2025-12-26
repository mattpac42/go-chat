package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

// AgentContextService determines context for agent responses.
type AgentContextService struct {
	prdRepo       PRDRepository
	projectRepo   repository.ProjectRepository
	discoveryRepo repository.DiscoveryRepository
	logger        zerolog.Logger
}

// NewAgentContextService creates a new AgentContextService.
func NewAgentContextService(
	prdRepo PRDRepository,
	projectRepo repository.ProjectRepository,
	discoveryRepo repository.DiscoveryRepository,
	logger zerolog.Logger,
) *AgentContextService {
	return &AgentContextService{
		prdRepo:       prdRepo,
		projectRepo:   projectRepo,
		discoveryRepo: discoveryRepo,
		logger:        logger,
	}
}

// GetContextForMessage determines the appropriate agent and PRD context for a message.
func (s *AgentContextService) GetContextForMessage(
	ctx context.Context,
	projectID uuid.UUID,
	message string,
) (*model.AgentContext, error) {
	// 1. Get the project
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("projectId", projectID.String()).
			Msg("failed to get project for agent context")
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	var prd *model.PRD

	// 2. If active PRD exists and is in progress, use it
	if project.ActivePRDID != nil {
		activePRD, err := s.prdRepo.GetByID(ctx, *project.ActivePRDID)
		if err == nil && activePRD.Status == model.PRDStatusInProgress {
			prd = activePRD
			s.logger.Debug().
				Str("prdId", prd.ID.String()).
				Str("title", prd.Title).
				Msg("using active PRD for context")
		}
	}

	// 3. Try to match message to a PRD by keywords
	if prd == nil {
		prds, err := s.prdRepo.GetByProjectID(ctx, projectID)
		if err == nil && len(prds) > 0 {
			prd = s.matchPRDByKeywords(message, prds)
			if prd != nil {
				s.logger.Debug().
					Str("prdId", prd.ID.String()).
					Str("title", prd.Title).
					Msg("matched PRD by keywords")
			}
		}

		// 4. Fallback: oldest "ready" PRD (next feature to build)
		if prd == nil && len(prds) > 0 {
			prd = s.getNextReadyPRD(prds)
			if prd != nil {
				s.logger.Debug().
					Str("prdId", prd.ID.String()).
					Str("title", prd.Title).
					Msg("using fallback next ready PRD")
			}
		}
	}

	return s.buildContext(ctx, project, prd, message)
}

// matchPRDByKeywords finds a PRD that matches keywords in the message.
func (s *AgentContextService) matchPRDByKeywords(message string, prds []model.PRD) *model.PRD {
	messageLower := strings.ToLower(message)

	for i := range prds {
		prd := &prds[i]
		titleLower := strings.ToLower(prd.Title)

		// Check if PRD title words appear in message
		titleWords := strings.Fields(titleLower)
		matchCount := 0
		for _, word := range titleWords {
			// Only check words longer than 3 characters to avoid noise
			if len(word) > 3 && strings.Contains(messageLower, word) {
				matchCount++
			}
		}

		// Match if 2+ significant words match or full title appears
		if matchCount >= 2 || strings.Contains(messageLower, titleLower) {
			return prd
		}
	}

	return nil
}

// getNextReadyPRD returns the oldest PRD with "ready" status, ordered by priority.
func (s *AgentContextService) getNextReadyPRD(prds []model.PRD) *model.PRD {
	var nextPRD *model.PRD

	for i := range prds {
		prd := &prds[i]
		if prd.Status != model.PRDStatusReady {
			continue
		}

		if nextPRD == nil || prd.Priority < nextPRD.Priority {
			nextPRD = prd
		}
	}

	return nextPRD
}

// SelectAgent determines which agent should respond based on user intent.
func (s *AgentContextService) SelectAgent(message string, prd *model.PRD) model.AgentType {
	messageLower := strings.ToLower(message)

	// Product Manager triggers
	productManagerKeywords := []string{
		"scope", "feature", "requirement", "user story",
		"priority", "roadmap", "why", "what should",
	}
	if containsAny(messageLower, productManagerKeywords) {
		return model.AgentProductManager
	}

	// Designer triggers
	designerKeywords := []string{
		"design", "layout", "ui", "ux", "look", "feel",
		"color", "style", "wireframe", "mockup",
	}
	if containsAny(messageLower, designerKeywords) {
		return model.AgentDesigner
	}

	// Developer is default for build phase
	return model.AgentDeveloper
}

// containsAny checks if the text contains any of the keywords.
func containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

// buildContext assembles the full agent context.
func (s *AgentContextService) buildContext(
	ctx context.Context,
	project *model.Project,
	prd *model.PRD,
	message string,
) (*model.AgentContext, error) {
	agent := s.SelectAgent(message, prd)

	agentContext := &model.AgentContext{
		Agent: agent,
		PRD:   prd,
	}

	// Get discovery summary for project context
	discovery, err := s.discoveryRepo.GetByProjectID(ctx, project.ID)
	if err == nil {
		summary, err := s.discoveryRepo.GetSummary(ctx, discovery.ID)
		if err == nil {
			agentContext.Discovery = summary
		}
	}

	// Build condensed PRD summary for token efficiency
	if prd != nil {
		agentContext.PRDSummary = s.CondensePRD(prd)

		// Get related PRDs for cross-feature awareness
		agentContext.RelatedPRDs = s.getRelatedPRDs(ctx, prd)
	}

	s.logger.Debug().
		Str("agent", string(agent)).
		Bool("hasPRD", prd != nil).
		Bool("hasDiscovery", agentContext.Discovery != nil).
		Msg("built agent context")

	return agentContext, nil
}

// CondensePRD creates a token-efficient PRD summary for prompts.
func (s *AgentContextService) CondensePRD(prd *model.PRD) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## %s\n\n", prd.Title))
	sb.WriteString(fmt.Sprintf("%s\n\n", prd.Overview))

	// Include user stories in condensed format
	stories, err := prd.UserStories()
	if err == nil && len(stories) > 0 {
		sb.WriteString("### User Stories\n")
		for _, story := range stories {
			sb.WriteString(fmt.Sprintf("- %s: As a %s, I want %s, so that %s\n",
				story.ID, story.AsA, story.IWant, story.SoThat))
		}
		sb.WriteString("\n")
	}

	// Include key acceptance criteria (limit to 5 for token efficiency)
	criteria, err := prd.AcceptanceCriteria()
	if err == nil && len(criteria) > 0 {
		sb.WriteString("### Key Acceptance Criteria\n")
		for i, ac := range criteria {
			if i >= 5 {
				sb.WriteString(fmt.Sprintf("- ... and %d more\n", len(criteria)-5))
				break
			}
			sb.WriteString(fmt.Sprintf("- %s: Given %s, When %s, Then %s\n",
				ac.ID, ac.Given, ac.When, ac.Then))
		}
	}

	return sb.String()
}

// getRelatedPRDs returns lightweight references to other PRDs in the same project.
func (s *AgentContextService) getRelatedPRDs(ctx context.Context, currentPRD *model.PRD) []model.PRDReference {
	prds, err := s.prdRepo.GetByProjectID(ctx, currentPRD.ProjectID)
	if err != nil {
		return nil
	}

	var related []model.PRDReference
	for _, prd := range prds {
		if prd.ID == currentPRD.ID {
			continue
		}
		related = append(related, prd.ToReference())
	}

	return related
}

// GetSystemPrompt builds the system prompt for the selected agent with PRD context.
func (s *AgentContextService) GetSystemPrompt(ctx context.Context, agentContext *model.AgentContext) (string, error) {
	if agentContext == nil {
		return "", nil
	}

	var promptTemplate string
	switch agentContext.Agent {
	case model.AgentProductManager:
		promptTemplate = productManagerPrompt
	case model.AgentDesigner:
		promptTemplate = designerPrompt
	case model.AgentDeveloper:
		promptTemplate = developerPrompt
	default:
		promptTemplate = developerPrompt
	}

	// Build template context
	templateCtx := s.buildTemplateContext(agentContext)

	// Parse and execute template
	tmpl, err := template.New("agent-prompt").Parse(promptTemplate)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to parse agent prompt template")
		return "", fmt.Errorf("failed to parse prompt template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateCtx); err != nil {
		s.logger.Error().Err(err).Msg("failed to execute agent prompt template")
		return "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return buf.String(), nil
}

// AgentPromptContext contains the data for agent prompt templates.
type AgentPromptContext struct {
	ProjectName    string
	PRDSummary     string
	SolvesStatement string
	Users          []model.DiscoveryUser
	TechnicalNotes []model.TechnicalNote
}

// buildTemplateContext creates the context for agent prompt templates.
func (s *AgentContextService) buildTemplateContext(agentContext *model.AgentContext) *AgentPromptContext {
	ctx := &AgentPromptContext{}

	// Set project name from discovery
	if agentContext.Discovery != nil {
		ctx.ProjectName = agentContext.Discovery.ProjectName
		ctx.SolvesStatement = agentContext.Discovery.SolvesStatement
		ctx.Users = agentContext.Discovery.Users
	}

	// Set PRD summary
	ctx.PRDSummary = agentContext.PRDSummary

	// Get technical notes if PRD is available
	if agentContext.PRD != nil {
		notes, err := agentContext.PRD.TechnicalNotes()
		if err == nil {
			ctx.TechnicalNotes = notes
		}
	}

	return ctx
}

// Agent prompt templates

const productManagerPrompt = `You are a Product Manager helping guide the development of {{.ProjectName}}.

## Your Role
- Clarify requirements and scope questions
- Break down features into actionable items
- Ensure user needs are addressed
- Maintain focus on MVP priorities

## Current Feature
{{.PRDSummary}}

## Project Context
{{.SolvesStatement}}

Users:
{{range .Users}}
- {{.Description}}
{{end}}

## Guidelines
- Always refer back to user stories when discussing scope
- Be concise and actionable
- Flag scope creep gently
- Suggest breaking large requests into phases

Respond as the Product Manager. Be helpful and focused on delivering value.`

const designerPrompt = `You are a UI/UX Designer working on {{.ProjectName}}.

## Your Role
- Create user-friendly interface designs
- Ensure accessibility and usability
- Design for the target users
- Follow mobile-first principles

## Current Feature
{{.PRDSummary}}

## Target Users
{{range .Users}}
- {{.Description}} ({{.UserCount}} users)
{{end}}

## Guidelines
- Start with mobile layouts
- Use simple, clear language in UI copy
- Consider the technical skill level of users
- Describe layouts in plain terms (no design jargon)

Respond as the Designer. Create interfaces that delight users.`

const developerPrompt = `You are a Developer building {{.ProjectName}}.

## Your Role
- Write clean, working code
- Follow the acceptance criteria
- Create files that work together
- Explain what you're building

## Current Feature
{{.PRDSummary}}

## Technical Notes
{{range .TechnicalNotes}}
**{{.Category}}**: {{.Description}}
{{end}}

## Guidelines
- Generate complete, working files
- Include helpful comments
- Follow the acceptance criteria exactly
- Explain choices in plain language

Respond as the Developer. Build features that meet the requirements.`
