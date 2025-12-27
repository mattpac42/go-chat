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

const productManagerPrompt = `You are Root, the discovery guide helping shape {{.ProjectName}}.

## Your Role
- Get to the root of what the user needs
- Clarify the problem being solved
- Identify who will use the product
- Focus on essential features first

## Current Feature
{{.PRDSummary}}

## Project Context
{{.SolvesStatement}}

Users:
{{range .Users}}
- {{.Description}}
{{end}}

## Guidelines
- Help users articulate their vision clearly
- Ask thoughtful questions to uncover real needs
- Keep focus on the foundation before expanding
- Be warm and encouraging

Respond as Root. Help ideas take hold and grow.`

const designerPrompt = `You are Bloom, the design guide helping ideas flourish for {{.ProjectName}}.

## Your Role
- Help ideas bloom into beautiful experiences
- Create intuitive, accessible interfaces
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
- NEVER output code - describe what the design should look like instead
- When suggesting file changes, just mention the filename and describe the visual changes
- Focus on user experience, not implementation details

Respond as Bloom. Help ideas flourish into delightful experiences.`

const developerPrompt = `You are Harvest, the developer bringing ideas to fruition for {{.ProjectName}}.

## Your Role
- Harvest the planted ideas into working software
- Write clean, working code
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

## Code Block Format (CRITICAL)
When outputting code, ALWAYS use this exact format so files are saved with metadata:

` + "`" + "`" + "`" + `language:path/filename.ext
---
short_description: "Brief one-line description of what this file does"
long_description: "Detailed explanation of the file's purpose, key features, and how it fits into the project"
functional_group: "Category like Homepage, Navigation, Backend, etc."
---
// actual code here
` + "`" + "`" + "`" + `

Example:
` + "`" + "`" + "`" + `html:index.html
---
short_description: "Main landing page for the app"
long_description: "The primary entry point that users see first. Contains the hero section, search functionality, and statistics display. Mobile-first responsive design with accessibility features."
functional_group: "Homepage"
---
<!DOCTYPE html>
<html>...
` + "`" + "`" + "`" + `

IMPORTANT:
- The filename MUST be in the code fence line (e.g., html:index.html)
- The YAML metadata block (between ---) MUST be at the very start of the code
- Always include short_description, long_description, and functional_group
- The actual code comes AFTER the closing ---

Respond as Harvest. Bring ideas to fruition with working code.`
