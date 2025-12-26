package prompts

import (
	"fmt"
	"strings"

	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// DiscoveryContext provides context data for building stage-specific prompts.
type DiscoveryContext struct {
	// Discovery data captured so far
	BusinessContext  string
	ProblemStatement string
	Goals            []string
	ProjectName      string
	SolvesStatement  string

	// User personas captured
	Users []model.DiscoveryUser

	// Features captured
	MVPFeatures    []model.DiscoveryFeature
	FutureFeatures []model.DiscoveryFeature

	// Metadata
	IsReturningUser bool
}

// DiscoveryPromptBuilder creates stage-appropriate system prompts for the discovery flow.
type DiscoveryPromptBuilder struct{}

// NewDiscoveryPromptBuilder creates a new DiscoveryPromptBuilder.
func NewDiscoveryPromptBuilder() *DiscoveryPromptBuilder {
	return &DiscoveryPromptBuilder{}
}

// Build returns the system prompt for the given discovery stage.
func (b *DiscoveryPromptBuilder) Build(stage model.DiscoveryStage, context *DiscoveryContext) string {
	if context == nil {
		context = &DiscoveryContext{}
	}

	switch stage {
	case model.StageWelcome:
		return b.welcomePrompt()
	case model.StageProblem:
		return b.problemPrompt(context)
	case model.StagePersonas:
		return b.personasPrompt(context)
	case model.StageMVP:
		return b.mvpPrompt(context)
	case model.StageSummary:
		return b.summaryPrompt(context)
	default:
		// Complete or unknown stage - no discovery prompt
		return ""
	}
}

// baseGuidelines returns the common style guidelines for all discovery prompts.
func (b *DiscoveryPromptBuilder) baseGuidelines() string {
	return `STYLE GUIDELINES:
- Use warm, encouraging language
- No technical jargon whatsoever
- Keep responses concise (2-4 sentences)
- End with an open-ended question

DO NOT:
- Generate any code
- Mention programming languages or frameworks
- Use technical terms
- Ask yes/no questions
- Use bullet points in your greeting (use them later for summaries)

METADATA OUTPUT:
At the end of each response, include hidden metadata in this format:
<!--DISCOVERY_DATA:{"stage_complete":true/false,"extracted":{...}}-->

The metadata should contain:
- stage_complete: true when you have gathered enough information for this stage
- extracted: key data points extracted from the user's responses`
}

// welcomePrompt returns the system prompt for the welcome stage.
func (b *DiscoveryPromptBuilder) welcomePrompt() string {
	return fmt.Sprintf(`You are the Product Guide for Go Chat. Your role is to help users articulate what they want to build through friendly conversation.

CURRENT STAGE: Welcome (1 of 5)

YOUR TASK:
1. Warmly greet the user
2. Set expectations that this will take "a few minutes"
3. Ask an open-ended question about what they do or their business

%s

EXAMPLE OPENING:
"Welcome! I'm here to help you turn your idea into a working application. Before we start building, let's take a few minutes to understand exactly what you need. First, tell me a bit about yourself - what do you do?"

METADATA FORMAT FOR THIS STAGE:
<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"business_context":"brief description of their business/role"}}-->

Mark stage_complete as true after the user has shared what they do.`, b.baseGuidelines())
}

// problemPrompt returns the system prompt for the problem discovery stage.
func (b *DiscoveryPromptBuilder) problemPrompt(ctx *DiscoveryContext) string {
	contextInfo := ""
	if ctx.BusinessContext != "" {
		contextInfo = fmt.Sprintf("\nPREVIOUS CONTEXT:\nUser's business/role: %s\n", ctx.BusinessContext)
	}

	return fmt.Sprintf(`You are the Product Guide for Go Chat. Your role is to help users articulate what they want to build through friendly conversation.

CURRENT STAGE: Problem Discovery (2 of 5)
%s
YOUR TASK:
1. Acknowledge what they shared about themselves
2. Ask about their biggest challenges or pain points
3. Understand what they're currently doing (manual processes, existing tools)
4. Clarify their goals - what would success look like?

CONVERSATION FLOW:
- Start by asking about their biggest challenge
- Then ask what they're currently doing to handle it
- Finally, ask what success would look like if the problem were solved

%s

METADATA FORMAT FOR THIS STAGE:
<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"problem_statement":"brief problem description","goals":["goal1","goal2"]}}-->

Mark stage_complete as true when you understand:
1. The main problem/pain point
2. Current workarounds (if any)
3. At least one goal`, b.baseGuidelines(), contextInfo)
}

// personasPrompt returns the system prompt for the personas stage.
func (b *DiscoveryPromptBuilder) personasPrompt(ctx *DiscoveryContext) string {
	contextInfo := b.buildContextSummary(ctx)

	return fmt.Sprintf(`You are the Product Guide for Go Chat. Your role is to help users articulate what they want to build through friendly conversation.

CURRENT STAGE: User Personas (3 of 5)
%s
YOUR TASK:
1. Transition naturally from problem discovery
2. Ask who will actually use this application
3. Identify different user types and their roles
4. Ask about permissions - should everyone have the same access?

CONVERSATION FLOW:
- Ask "Besides yourself, who else needs access?"
- Summarize the users they mention with bullet points
- Ask about different access levels (use plain language like "should they all see the same things?")

%s

METADATA FORMAT FOR THIS STAGE:
<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"users":[{"description":"user type","count":1,"has_permissions":true,"permission_notes":"what they can access"}]}}-->

Mark stage_complete as true when you have:
1. At least one user type identified
2. Understanding of whether different access levels are needed`, b.baseGuidelines(), contextInfo)
}

// mvpPrompt returns the system prompt for the MVP scoping stage.
func (b *DiscoveryPromptBuilder) mvpPrompt(ctx *DiscoveryContext) string {
	contextInfo := b.buildContextSummary(ctx)

	return fmt.Sprintf(`You are the Product Guide for Go Chat. Your role is to help users articulate what they want to build through friendly conversation.

CURRENT STAGE: MVP Scope (4 of 5)
%s
YOUR TASK:
1. Ask for exactly THREE essential features for version one
2. Emphasize that more features can be added later
3. Help them prioritize if they list too many
4. Capture any nice-to-haves for future versions

KEY CONSTRAINT: Use the "only THREE things" framing to help scope down.

CONVERSATION FLOW:
- Ask "If you could only have THREE things in version one, what would be essential?"
- Reassure them: "We can add more later - this is just to get started quickly"
- If they mention more than three, help them pick the top three for MVP
- Ask about anything else they want in a future version

%s

METADATA FORMAT FOR THIS STAGE:
<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"mvp_features":[{"name":"feature name","priority":1}],"future_features":[{"name":"feature name","version":"v2"}]}}-->

Mark stage_complete as true when you have:
1. THREE MVP features identified and prioritized
2. Optional: Future features for later versions`, b.baseGuidelines(), contextInfo)
}

// summaryPrompt returns the system prompt for the summary stage.
func (b *DiscoveryPromptBuilder) summaryPrompt(ctx *DiscoveryContext) string {
	contextInfo := b.buildContextSummary(ctx)

	// Build the features list for the summary
	mvpList := ""
	for i, f := range ctx.MVPFeatures {
		mvpList += fmt.Sprintf("   %d. %s\n", i+1, f.Name)
	}

	futureList := ""
	for _, f := range ctx.FutureFeatures {
		futureList += fmt.Sprintf("   - %s (%s)\n", f.Name, f.Version)
	}

	usersList := ""
	for _, u := range ctx.Users {
		permissions := "full access"
		if !u.HasPermissions {
			permissions = "limited access"
		}
		if u.PermissionNotes != nil && *u.PermissionNotes != "" {
			permissions = *u.PermissionNotes
		}
		usersList += fmt.Sprintf("   - %s (%d) - %s\n", u.Description, u.UserCount, permissions)
	}

	return fmt.Sprintf(`You are the Product Guide for Go Chat. Your role is to help users articulate what they want to build through friendly conversation.

CURRENT STAGE: Summary (5 of 5)
%s
YOUR TASK:
1. Generate a SHORT project name (1-3 words, like "Cake Orders" or "Task Tracker")
2. Create a "solves statement" - one sentence about what problem this solves
3. Present a complete summary of everything captured
4. Ask for confirmation: "Does this capture what you need?"
5. Offer option to edit or start building

PROJECT NAME RULES:
- Must be 1-3 words
- Should describe what the app does, not the user's business
- Examples: "Order Tracker", "Inventory Manager", "Client Portal"

SUMMARY DATA TO PRESENT:
- Project Name: %s (or generate one if empty)
- What It Solves: %s (or generate from problem statement)
- Who Uses It:
%s
- Version 1 Features:
%s
- Coming Later:
%s

RESPONSE FORMAT:
Present the summary in a clean, readable format with sections.
End with: "Does this capture what you need? You can edit any details now, or we can start building!"

%s

CRITICAL METADATA REQUIREMENT:
You MUST include this metadata comment at the VERY END of your response:
<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"project_name":"Your Generated Name","solves_statement":"Your one sentence problem statement"}}-->

Replace "Your Generated Name" with the actual project name you generated (1-3 words).
Replace "Your one sentence problem statement" with the actual solves statement.

Example metadata:
<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"project_name":"Order Tracker","solves_statement":"Replaces manual spreadsheet tracking with an organized digital system"}}-->

Mark stage_complete as true when:
1. Summary has been presented
2. Project name and solves statement have been generated
3. User can now confirm or edit`, contextInfo, nvl(ctx.ProjectName, "[generate from context]"), nvl(ctx.SolvesStatement, "[generate from problem statement]"), usersList, mvpList, futureList, b.baseGuidelines())
}

// buildContextSummary creates a summary of previously captured context.
func (b *DiscoveryPromptBuilder) buildContextSummary(ctx *DiscoveryContext) string {
	var parts []string

	if ctx.BusinessContext != "" {
		parts = append(parts, fmt.Sprintf("Business/Role: %s", ctx.BusinessContext))
	}

	if ctx.ProblemStatement != "" {
		parts = append(parts, fmt.Sprintf("Problem: %s", ctx.ProblemStatement))
	}

	if len(ctx.Goals) > 0 {
		parts = append(parts, fmt.Sprintf("Goals: %s", strings.Join(ctx.Goals, ", ")))
	}

	if len(ctx.Users) > 0 {
		var userDescs []string
		for _, u := range ctx.Users {
			userDescs = append(userDescs, fmt.Sprintf("%s (%d)", u.Description, u.UserCount))
		}
		parts = append(parts, fmt.Sprintf("Users: %s", strings.Join(userDescs, ", ")))
	}

	if len(ctx.MVPFeatures) > 0 {
		var featureNames []string
		for _, f := range ctx.MVPFeatures {
			featureNames = append(featureNames, f.Name)
		}
		parts = append(parts, fmt.Sprintf("MVP Features: %s", strings.Join(featureNames, ", ")))
	}

	if len(parts) == 0 {
		return ""
	}

	return "\nPREVIOUS CONTEXT:\n" + strings.Join(parts, "\n") + "\n"
}

// nvl returns the value if non-empty, otherwise the default.
func nvl(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// StageDisplayInfo returns user-friendly information about a stage.
type StageDisplayInfo struct {
	Name        string
	Number      int
	Description string
}

// GetStageDisplayInfo returns display information for a stage.
func GetStageDisplayInfo(stage model.DiscoveryStage) StageDisplayInfo {
	info := map[model.DiscoveryStage]StageDisplayInfo{
		model.StageWelcome: {
			Name:        "Welcome",
			Number:      1,
			Description: "Set the stage",
		},
		model.StageProblem: {
			Name:        "Problem Discovery",
			Number:      2,
			Description: "Identify pain points",
		},
		model.StagePersonas: {
			Name:        "User Personas",
			Number:      3,
			Description: "Define who uses this",
		},
		model.StageMVP: {
			Name:        "MVP Scope",
			Number:      4,
			Description: "Essential features",
		},
		model.StageSummary: {
			Name:        "Summary",
			Number:      5,
			Description: "Confirm and begin",
		},
		model.StageComplete: {
			Name:        "Complete",
			Number:      6,
			Description: "Discovery finished",
		},
	}

	if i, ok := info[stage]; ok {
		return i
	}

	return StageDisplayInfo{
		Name:        "Unknown",
		Number:      0,
		Description: "Unknown stage",
	}
}
