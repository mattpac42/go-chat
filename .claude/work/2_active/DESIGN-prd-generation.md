# PRD Auto-Generation and Agent Context System

**Date**: 2025-12-26
**Author**: Architect Agent
**Status**: Design Specification
**Version**: 1.0

---

## Executive Summary

This document defines the architecture for automatic PRD generation from completed discovery and the agent context system that provides relevant PRDs to specialized agents (Product Manager, Designer, Developer) during the build phase.

**Key Architectural Decisions**:
- One PRD per feature (not per phase) for granular tracking and focused agent context
- Claude generates PRD content using discovery data as structured input
- PRDs are generated asynchronously upon discovery confirmation
- Agent context is determined by matching user intent to relevant PRDs
- "Active PRD" concept tracks which feature the user is currently working on

---

## 1. Architecture Overview

### 1.1 System Flow Diagram

```
+==============================================================================+
|                         PRD GENERATION FLOW                                   |
+==============================================================================+

     DISCOVERY COMPLETE
           |
           v
+----------------------+     +----------------------+     +------------------+
|  Discovery Service   |---->|  PRD Generator       |---->|  PRD Repository  |
|                      |     |  Service             |     |                  |
|  - discovery data    |     |  - template context  |     |  - CRUD ops      |
|  - features list     |     |  - Claude generation |     |  - status track  |
|  - user personas     |     |  - async processing  |     |  - version mgmt  |
+----------------------+     +----------------------+     +------------------+
                                      |
                                      v
                             +------------------+
                             |   Claude API     |
                             |                  |
                             |  - PRD content   |
                             |  - user stories  |
                             |  - acceptance    |
                             +------------------+

+==============================================================================+
|                         AGENT CONTEXT FLOW                                    |
+==============================================================================+

     USER MESSAGE (Build Phase)
           |
           v
+----------------------+     +----------------------+     +------------------+
|  Chat Service        |---->|  Agent Context       |---->|  System Prompt   |
|                      |     |  Service             |     |  Builder         |
|  - user message      |     |  - intent detection  |     |                  |
|  - project context   |     |  - PRD matching      |     |  - agent persona |
|  - active PRD        |     |  - context assembly  |     |  - PRD context   |
+----------------------+     +----------------------+     +------------------+
                                      |
                                      v
                             +------------------+
                             |  Selected Agent  |
                             |                  |
                             |  - Product PM    |
                             |  - Designer      |
                             |  - Developer     |
                             +------------------+
```

### 1.2 Component Relationships

```
+------------------+          +------------------+          +------------------+
|   Discovery      |   1:N    |   PRD            |   N:1    |   Project        |
|                  |<-------->|                  |<-------->|                  |
|  - project_id    |          |  - discovery_id  |          |  - id            |
|  - features[]    |          |  - feature_id    |          |  - active_prd_id |
|  - personas[]    |          |  - status        |          |                  |
+------------------+          +------------------+          +------------------+
         |                           |
         |                           |
         v                           v
+------------------+          +------------------+
|  Discovery       |          |  PRD Sections    |
|  Feature         |          |                  |
|                  |          |  - overview      |
|  - name          |          |  - user_stories  |
|  - priority      |          |  - acceptance    |
|  - version       |          |  - tech_notes    |
+------------------+          +------------------+
```

---

## 2. Data Model Design

### 2.1 PRD Structure Decision

**Decision**: One PRD per feature (not per phase/version)

**Rationale**:
- **Granular tracking**: Each feature can have independent status (draft, ready, in_progress, complete)
- **Focused context**: Agents receive only relevant PRD for current work, reducing token usage
- **Flexible prioritization**: Features can be reordered without restructuring PRDs
- **Incremental delivery**: Users see tangible progress as each PRD completes
- **Clear handoffs**: Developer knows exactly which feature to implement next

**Alternative Considered**: One PRD per version (v1, v2, etc.)
- Rejected because: Harder to track individual feature progress; larger context documents; less clear agent focus

### 2.2 PRD Model

```go
// PRD represents a Product Requirements Document for a single feature.
type PRD struct {
    ID           uuid.UUID   `db:"id" json:"id"`
    DiscoveryID  uuid.UUID   `db:"discovery_id" json:"discoveryId"`
    FeatureID    uuid.UUID   `db:"feature_id" json:"featureId"`
    ProjectID    uuid.UUID   `db:"project_id" json:"projectId"`

    // Core PRD Content
    Title        string      `db:"title" json:"title"`
    Overview     string      `db:"overview" json:"overview"`
    Version      string      `db:"version" json:"version"` // "v1", "v2", etc.
    Priority     int         `db:"priority" json:"priority"`

    // Detailed Sections (JSONB for flexibility)
    UserStoriesJSON    []byte `db:"user_stories" json:"-"`
    AcceptanceCriteria []byte `db:"acceptance_criteria" json:"-"`
    TechnicalNotesJSON []byte `db:"technical_notes" json:"-"`

    // Status Tracking
    Status       PRDStatus   `db:"status" json:"status"`
    GeneratedAt  *time.Time  `db:"generated_at" json:"generatedAt,omitempty"`
    ApprovedAt   *time.Time  `db:"approved_at" json:"approvedAt,omitempty"`
    StartedAt    *time.Time  `db:"started_at" json:"startedAt,omitempty"`
    CompletedAt  *time.Time  `db:"completed_at" json:"completedAt,omitempty"`

    // Metadata
    GenerationAttempts int       `db:"generation_attempts" json:"generationAttempts"`
    LastError          *string   `db:"last_error" json:"lastError,omitempty"`

    CreatedAt    time.Time   `db:"created_at" json:"createdAt"`
    UpdatedAt    time.Time   `db:"updated_at" json:"updatedAt"`
}

// PRDStatus represents the lifecycle state of a PRD.
type PRDStatus string

const (
    PRDStatusPending    PRDStatus = "pending"     // Queued for generation
    PRDStatusGenerating PRDStatus = "generating"  // Claude is generating content
    PRDStatusDraft      PRDStatus = "draft"       // Generated, awaiting review
    PRDStatusReady      PRDStatus = "ready"       // Approved, ready to build
    PRDStatusInProgress PRDStatus = "in_progress" // Currently being implemented
    PRDStatusComplete   PRDStatus = "complete"    // Feature implemented
    PRDStatusFailed     PRDStatus = "failed"      // Generation failed
)

// UserStory represents a single user story within a PRD.
type UserStory struct {
    ID          string `json:"id"`          // e.g., "US-001"
    AsA         string `json:"asA"`         // User persona
    IWant       string `json:"iWant"`       // Desired action
    SoThat      string `json:"soThat"`      // Expected benefit
    Priority    string `json:"priority"`    // "must", "should", "could"
    Complexity  string `json:"complexity"`  // "low", "medium", "high"
}

// AcceptanceCriterion represents a single acceptance criterion.
type AcceptanceCriterion struct {
    ID          string `json:"id"`          // e.g., "AC-001"
    Given       string `json:"given"`       // Precondition
    When        string `json:"when"`        // Action
    Then        string `json:"then"`        // Expected outcome
    UserStoryID string `json:"userStoryId"` // Links to parent story
}

// TechnicalNote captures implementation guidance.
type TechnicalNote struct {
    Category    string   `json:"category"`    // "architecture", "data", "ui", "integration"
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Suggestions []string `json:"suggestions,omitempty"`
}
```

### 2.3 Active PRD Tracking

```go
// ProjectActivePRD tracks which PRD is currently being worked on.
// This is stored on the Project model to avoid separate table.
type Project struct {
    // ... existing fields ...

    // Active PRD for agent context
    ActivePRDID *uuid.UUID `db:"active_prd_id" json:"activePrdId,omitempty"`
}
```

### 2.4 Agent Context Model

```go
// AgentContext represents the context provided to an agent.
type AgentContext struct {
    Agent       AgentType       `json:"agent"`
    PRD         *PRD            `json:"prd,omitempty"`
    Discovery   *DiscoverySummary `json:"discovery"`

    // Condensed PRD for prompt
    PRDSummary  string          `json:"prdSummary,omitempty"`

    // Related PRDs (for cross-feature awareness)
    RelatedPRDs []PRDReference  `json:"relatedPrds,omitempty"`
}

// AgentType identifies the specialized agent.
type AgentType string

const (
    AgentProductManager AgentType = "product_manager"
    AgentDesigner       AgentType = "designer"
    AgentDeveloper      AgentType = "developer"
)

// PRDReference is a lightweight PRD summary for context.
type PRDReference struct {
    ID       uuid.UUID `json:"id"`
    Title    string    `json:"title"`
    Status   PRDStatus `json:"status"`
    Priority int       `json:"priority"`
}
```

---

## 3. PRD Generation

### 3.1 Generation Trigger

**Decision**: Generate PRDs asynchronously upon discovery confirmation

**Trigger Point**:
```go
// In DiscoveryService.ConfirmDiscovery()
func (s *DiscoveryService) ConfirmDiscovery(ctx context.Context, discoveryID uuid.UUID) (*ProjectDiscovery, error) {
    // ... existing confirmation logic ...

    // Trigger PRD generation asynchronously
    go s.prdGenerator.GenerateAllPRDs(context.Background(), discoveryID)

    return discovery, nil
}
```

**Why Async**:
- User doesn't wait for generation (can proceed to chat)
- Multiple PRDs can be generated in parallel
- Failures don't block the confirmation
- User sees PRDs appear progressively

### 3.2 Generation Process

```
PRD GENERATION PIPELINE
=======================

1. INITIALIZATION
   - Fetch discovery summary (problem, personas, features)
   - Create PRD records with status "pending" for each feature
   - Sort by priority (MVP features first)

2. FOR EACH FEATURE (parallel for MVP, sequential for future):
   a. Set status to "generating"
   b. Build generation context:
      - Discovery summary
      - Feature details
      - User personas
      - Related features (for awareness)
   c. Call Claude API with PRD generation prompt
   d. Parse response into structured PRD sections
   e. Validate completeness
   f. Save to database
   g. Set status to "draft"

3. COMPLETION
   - Notify frontend via WebSocket (optional)
   - First MVP PRD auto-set as "active"
   - Log generation metrics
```

### 3.3 Generation Prompt Template

```go
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
  "userStories": [...],
  "acceptanceCriteria": [...],
  "technicalNotes": [...]
}
`
```

### 3.4 PRD Detail Levels

**MVP Features (v1)**: Full PRD with all sections
**Future Features (v2+)**: Lightweight PRD (overview only, user stories deferred)

```go
func (g *PRDGenerator) GeneratePRD(ctx context.Context, feature *DiscoveryFeature, discovery *DiscoverySummary) (*PRD, error) {
    if feature.IsMVP() {
        return g.generateFullPRD(ctx, feature, discovery)
    }
    return g.generateLightweightPRD(ctx, feature, discovery)
}
```

---

## 4. Agent Context System

### 4.1 Context Determination Flow

```
USER MESSAGE
     |
     v
+--------------------+
| Intent Detection   |
| - feature mention? |
| - action type?     |
| - clarification?   |
+--------------------+
     |
     v
+--------------------+
| PRD Matching       |
| - active PRD       |
| - keyword match    |
| - fallback: oldest |
|   ready PRD        |
+--------------------+
     |
     v
+--------------------+
| Agent Selection    |
| - product: scope   |
| - designer: UI/UX  |
| - developer: code  |
+--------------------+
     |
     v
+--------------------+
| Context Assembly   |
| - agent persona    |
| - PRD context      |
| - discovery summary|
+--------------------+
```

### 4.2 PRD Matching Logic

```go
// AgentContextService determines context for agent responses.
type AgentContextService struct {
    prdRepo     repository.PRDRepository
    projectRepo repository.ProjectRepository
    logger      zerolog.Logger
}

// GetContextForMessage determines the appropriate agent and PRD context.
func (s *AgentContextService) GetContextForMessage(
    ctx context.Context,
    projectID uuid.UUID,
    message string,
) (*AgentContext, error) {
    // 1. Check for explicit active PRD
    project, err := s.projectRepo.GetByID(ctx, projectID)
    if err != nil {
        return nil, err
    }

    var prd *model.PRD

    // 2. If active PRD exists and is in progress, use it
    if project.ActivePRDID != nil {
        prd, err = s.prdRepo.GetByID(ctx, *project.ActivePRDID)
        if err == nil && prd.Status == model.PRDStatusInProgress {
            return s.buildContext(ctx, project, prd, message)
        }
    }

    // 3. Try to match message to a PRD by keywords
    prds, _ := s.prdRepo.GetByProjectID(ctx, projectID)
    prd = s.matchPRDByKeywords(message, prds)

    // 4. Fallback: oldest "ready" PRD (next feature to build)
    if prd == nil {
        prd = s.getNextReadyPRD(prds)
    }

    return s.buildContext(ctx, project, prd, message)
}

// matchPRDByKeywords finds a PRD that matches keywords in the message.
func (s *AgentContextService) matchPRDByKeywords(message string, prds []model.PRD) *model.PRD {
    messageLower := strings.ToLower(message)

    for _, prd := range prds {
        titleLower := strings.ToLower(prd.Title)
        // Check if PRD title words appear in message
        titleWords := strings.Fields(titleLower)
        matchCount := 0
        for _, word := range titleWords {
            if len(word) > 3 && strings.Contains(messageLower, word) {
                matchCount++
            }
        }
        if matchCount >= 2 || strings.Contains(messageLower, titleLower) {
            return &prd
        }
    }

    return nil
}
```

### 4.3 Agent Selection

```go
// SelectAgent determines which agent should respond based on user intent.
func (s *AgentContextService) SelectAgent(message string, prd *model.PRD) AgentType {
    messageLower := strings.ToLower(message)

    // Product Manager triggers
    if containsAny(messageLower, []string{
        "scope", "feature", "requirement", "user story",
        "priority", "roadmap", "why", "what should",
    }) {
        return AgentProductManager
    }

    // Designer triggers
    if containsAny(messageLower, []string{
        "design", "layout", "ui", "ux", "look", "feel",
        "color", "style", "wireframe", "mockup",
    }) {
        return AgentDesigner
    }

    // Developer is default for build phase
    return AgentDeveloper
}
```

### 4.4 Context Assembly

```go
// buildContext assembles the full agent context.
func (s *AgentContextService) buildContext(
    ctx context.Context,
    project *model.Project,
    prd *model.PRD,
    message string,
) (*AgentContext, error) {
    agent := s.SelectAgent(message, prd)

    context := &AgentContext{
        Agent: agent,
        PRD:   prd,
    }

    // Get discovery summary for project context
    if discovery, err := s.discoveryRepo.GetSummary(ctx, project.ID); err == nil {
        context.Discovery = discovery
    }

    // Build condensed PRD summary for token efficiency
    if prd != nil {
        context.PRDSummary = s.condensePRD(prd)

        // Get related PRDs for cross-feature awareness
        context.RelatedPRDs = s.getRelatedPRDs(ctx, prd)
    }

    return context, nil
}

// condensePRD creates a token-efficient PRD summary.
func (s *AgentContextService) condensePRD(prd *model.PRD) string {
    var sb strings.Builder

    sb.WriteString(fmt.Sprintf("## %s\n\n", prd.Title))
    sb.WriteString(fmt.Sprintf("%s\n\n", prd.Overview))

    // Include user stories in condensed format
    stories, _ := prd.UserStories()
    sb.WriteString("### User Stories\n")
    for _, story := range stories {
        sb.WriteString(fmt.Sprintf("- %s: As a %s, I want %s, so that %s\n",
            story.ID, story.AsA, story.IWant, story.SoThat))
    }

    // Include key acceptance criteria
    criteria, _ := prd.AcceptanceCriteria()
    sb.WriteString("\n### Key Acceptance Criteria\n")
    for i, ac := range criteria {
        if i >= 5 { // Limit for token efficiency
            sb.WriteString(fmt.Sprintf("- ... and %d more\n", len(criteria)-5))
            break
        }
        sb.WriteString(fmt.Sprintf("- %s: Given %s, When %s, Then %s\n",
            ac.ID, ac.Given, ac.When, ac.Then))
    }

    return sb.String()
}
```

---

## 5. Database Schema

### 5.1 SQL Migration (005_prds.sql)

```sql
-- 005_prds.sql
-- PRD (Product Requirements Document) tables for feature tracking

-- PRD status enum (using CHECK constraint for PostgreSQL)
CREATE TYPE prd_status AS ENUM (
    'pending',
    'generating',
    'draft',
    'ready',
    'in_progress',
    'complete',
    'failed'
);

-- Main PRD table
CREATE TABLE IF NOT EXISTS prds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discovery_id UUID NOT NULL REFERENCES project_discovery(id) ON DELETE CASCADE,
    feature_id UUID NOT NULL REFERENCES discovery_features(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

    -- Core Content
    title VARCHAR(255) NOT NULL,
    overview TEXT,
    version VARCHAR(10) NOT NULL DEFAULT 'v1',
    priority INTEGER NOT NULL DEFAULT 0,

    -- Structured Sections (JSONB)
    user_stories JSONB DEFAULT '[]'::JSONB,
    acceptance_criteria JSONB DEFAULT '[]'::JSONB,
    technical_notes JSONB DEFAULT '[]'::JSONB,

    -- Status Tracking
    status prd_status NOT NULL DEFAULT 'pending',
    generated_at TIMESTAMP WITH TIME ZONE,
    approved_at TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,

    -- Generation Metadata
    generation_attempts INTEGER DEFAULT 0,
    last_error TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT unique_feature_prd UNIQUE (feature_id),
    CONSTRAINT valid_priority CHECK (priority >= 0)
);

-- Add active PRD reference to projects
ALTER TABLE projects
ADD COLUMN IF NOT EXISTS active_prd_id UUID REFERENCES prds(id) ON DELETE SET NULL;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_prds_discovery ON prds(discovery_id);
CREATE INDEX IF NOT EXISTS idx_prds_project ON prds(project_id);
CREATE INDEX IF NOT EXISTS idx_prds_status ON prds(status);
CREATE INDEX IF NOT EXISTS idx_prds_version ON prds(version);
CREATE INDEX IF NOT EXISTS idx_prds_priority ON prds(priority);
CREATE INDEX IF NOT EXISTS idx_projects_active_prd ON projects(active_prd_id);

-- Updated_at trigger
CREATE OR REPLACE FUNCTION update_prds_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_prds_updated_at
    BEFORE UPDATE ON prds
    FOR EACH ROW
    EXECUTE FUNCTION update_prds_updated_at();

-- Comments
COMMENT ON TABLE prds IS 'Product Requirements Documents auto-generated from discovery features';
COMMENT ON COLUMN prds.discovery_id IS 'Source discovery that generated this PRD';
COMMENT ON COLUMN prds.feature_id IS 'Feature this PRD documents (one PRD per feature)';
COMMENT ON COLUMN prds.status IS 'Lifecycle status: pending->generating->draft->ready->in_progress->complete';
COMMENT ON COLUMN prds.user_stories IS 'JSON array of user stories with As-a/I-want/So-that format';
COMMENT ON COLUMN prds.acceptance_criteria IS 'JSON array of Given/When/Then acceptance criteria';
COMMENT ON COLUMN prds.technical_notes IS 'JSON array of implementation notes by category';
COMMENT ON COLUMN prds.generation_attempts IS 'Number of times Claude generation was attempted';
COMMENT ON COLUMN prds.last_error IS 'Error message from last failed generation attempt';
COMMENT ON COLUMN projects.active_prd_id IS 'Currently active PRD for focused agent context';
```

---

## 6. Service Layer

### 6.1 PRDService Interface

```go
package service

import (
    "context"
    "github.com/google/uuid"
    "gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// PRDService manages PRD lifecycle and generation.
type PRDService interface {
    // Generation
    GenerateAllPRDs(ctx context.Context, discoveryID uuid.UUID) error
    GeneratePRD(ctx context.Context, featureID uuid.UUID) (*model.PRD, error)
    RetryGeneration(ctx context.Context, prdID uuid.UUID) (*model.PRD, error)

    // Retrieval
    GetByID(ctx context.Context, prdID uuid.UUID) (*model.PRD, error)
    GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.PRD, error)
    GetByDiscoveryID(ctx context.Context, discoveryID uuid.UUID) ([]model.PRD, error)
    GetMVPPRDs(ctx context.Context, projectID uuid.UUID) ([]model.PRD, error)

    // Status Management
    UpdateStatus(ctx context.Context, prdID uuid.UUID, status model.PRDStatus) error
    MarkAsReady(ctx context.Context, prdID uuid.UUID) error
    StartImplementation(ctx context.Context, prdID uuid.UUID) error
    CompleteImplementation(ctx context.Context, prdID uuid.UUID) error

    // Active PRD
    GetActivePRD(ctx context.Context, projectID uuid.UUID) (*model.PRD, error)
    SetActivePRD(ctx context.Context, projectID uuid.UUID, prdID uuid.UUID) error
    ClearActivePRD(ctx context.Context, projectID uuid.UUID) error
    GetNextPRD(ctx context.Context, projectID uuid.UUID) (*model.PRD, error)

    // Content
    UpdateOverview(ctx context.Context, prdID uuid.UUID, overview string) error
    AddUserStory(ctx context.Context, prdID uuid.UUID, story *model.UserStory) error
    UpdateUserStory(ctx context.Context, prdID uuid.UUID, storyID string, story *model.UserStory) error
    DeleteUserStory(ctx context.Context, prdID uuid.UUID, storyID string) error
}

// AgentContextService provides context for agent responses.
type AgentContextService interface {
    // GetContextForMessage determines agent and PRD context for a message.
    GetContextForMessage(ctx context.Context, projectID uuid.UUID, message string) (*model.AgentContext, error)

    // GetSystemPrompt builds the system prompt for the selected agent with PRD context.
    GetSystemPrompt(ctx context.Context, agentContext *model.AgentContext) (string, error)

    // SelectAgent determines which agent should respond.
    SelectAgent(message string, prd *model.PRD) model.AgentType

    // CondensePRD creates a token-efficient PRD summary for prompts.
    CondensePRD(prd *model.PRD) string
}
```

### 6.2 PRDRepository Interface

```go
package repository

import (
    "context"
    "github.com/google/uuid"
    "gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

// PRDRepository defines database operations for PRDs.
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
```

---

## 7. Integration Points

### 7.1 Chat Service Integration

```go
// In ChatService.ProcessMessage()
func (s *ChatService) ProcessMessage(ctx context.Context, projectID uuid.UUID, content string, onChunk func(string)) (*ChatResult, error) {
    // ... existing discovery check ...

    // If discovery is complete, get agent context
    if discovery != nil && discovery.Stage.IsComplete() {
        agentCtx, err := s.agentContextService.GetContextForMessage(ctx, projectID, content)
        if err != nil {
            s.logger.Warn().Err(err).Msg("failed to get agent context")
        } else {
            // Use agent-specific system prompt with PRD context
            systemPrompt = s.agentContextService.GetSystemPrompt(ctx, agentCtx)

            // Update active PRD if matched to a different one
            if agentCtx.PRD != nil {
                s.prdService.SetActivePRD(ctx, projectID, agentCtx.PRD.ID)
            }
        }
    }

    // ... rest of message processing ...
}
```

### 7.2 Discovery Service Integration

```go
// In DiscoveryService.ConfirmDiscovery()
func (s *DiscoveryService) ConfirmDiscovery(ctx context.Context, discoveryID uuid.UUID) (*model.ProjectDiscovery, error) {
    // ... existing logic ...

    // Mark complete
    discovery, err := s.repo.MarkComplete(ctx, discoveryID)
    if err != nil {
        return nil, err
    }

    // Trigger async PRD generation
    go func() {
        if err := s.prdService.GenerateAllPRDs(context.Background(), discoveryID); err != nil {
            s.logger.Error().
                Err(err).
                Str("discoveryId", discoveryID.String()).
                Msg("failed to generate PRDs")
        }
    }()

    return discovery, nil
}
```

### 7.3 WebSocket Notifications (Optional Enhancement)

```go
// PRD generation events pushed to frontend
type PRDEvent struct {
    Type     string       `json:"type"`     // "prd.generating", "prd.ready", "prd.failed"
    PRDID    uuid.UUID    `json:"prdId"`
    Title    string       `json:"title"`
    Status   string       `json:"status"`
    Progress int          `json:"progress"` // 0-100 for multi-PRD generation
}
```

---

## 8. Agent System Prompts

### 8.1 Product Manager Agent

```go
const productManagerPrompt = `You are a Product Manager helping guide the development of {{.ProjectName}}.

## Your Role
- Clarify requirements and scope questions
- Break down features into actionable items
- Ensure user needs are addressed
- Maintain focus on MVP priorities

## Current Feature
{{.PRDSummary}}

## Project Context
{{.Discovery.SolvesStatement}}

Users:
{{range .Discovery.Users}}
- {{.Description}}
{{end}}

## Guidelines
- Always refer back to user stories when discussing scope
- Be concise and actionable
- Flag scope creep gently
- Suggest breaking large requests into phases

Respond as the Product Manager. Be helpful and focused on delivering value.`
```

### 8.2 Designer Agent

```go
const designerPrompt = `You are a UI/UX Designer working on {{.ProjectName}}.

## Your Role
- Create user-friendly interface designs
- Ensure accessibility and usability
- Design for the target users
- Follow mobile-first principles

## Current Feature
{{.PRDSummary}}

## Target Users
{{range .Discovery.Users}}
- {{.Description}} ({{.UserCount}} users)
{{end}}

## Guidelines
- Start with mobile layouts
- Use simple, clear language in UI copy
- Consider the technical skill level of users
- Describe layouts in plain terms (no design jargon)

Respond as the Designer. Create interfaces that delight users.`
```

### 8.3 Developer Agent

```go
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
```

---

## 9. Sequence Diagrams

### 9.1 PRD Generation Sequence

```
User          DiscoveryService      PRDService       ClaudeAPI        Database
  |                  |                   |                |                |
  |--Confirm-------->|                   |                |                |
  |                  |--MarkComplete---->|                |                |-->UPDATE
  |                  |                   |                |                |<--
  |                  |--GenerateAll----->|                |                |
  |                  |                   |--GetFeatures---|                |-->SELECT
  |<--OK-------------|                   |<---------------|                |<--
  |                  |                   |                |                |
  |                  |                   |--CreatePRD-----|                |-->INSERT
  |                  |                   |<---------------|                |<--
  |                  |                   |                |                |
  |                  |                   |--Generate----->|                |
  |                  |                   |<--Content------|                |
  |                  |                   |                |                |
  |                  |                   |--UpdatePRD-----|                |-->UPDATE
  |                  |                   |<---------------|                |<--
  |                  |                   |                |                |
  |<--[Optional WS notification]---------|                |                |
```

### 9.2 Agent Context Sequence

```
User          ChatService       AgentContextSvc    PRDService       Database
  |                |                   |                |                |
  |--Message------>|                   |                |                |
  |                |--GetContext------>|                |                |
  |                |                   |--GetActive---->|                |-->SELECT
  |                |                   |<---------------|                |<--
  |                |                   |                |                |
  |                |                   |--MatchByKeyword|                |
  |                |                   |                |                |
  |                |                   |--SelectAgent---|                |
  |                |                   |                |                |
  |                |<--AgentContext----|                |                |
  |                |                   |                |                |
  |                |--BuildPrompt----->|                |                |
  |                |                   |                |                |
  |                |--SendToClaud------|                |                |
  |                |                   |                |                |
  |<--Response-----|                   |                |                |
```

---

## 10. Migration Path

### 10.1 Implementation Phases

**Phase 1: Data Model (Day 1)**
- Create migration 005_prds.sql
- Add PRD model to backend
- Implement PRDRepository

**Phase 2: PRD Generation (Day 2-3)**
- Implement PRDService
- Create generation prompts
- Add Claude integration for PRD content
- Hook into DiscoveryService.ConfirmDiscovery()

**Phase 3: Agent Context (Day 4-5)**
- Implement AgentContextService
- Create agent persona prompts
- Integrate with ChatService
- Add active PRD tracking

**Phase 4: Frontend (Day 6-7)**
- PRD list view in project
- PRD status indicators
- Active PRD badge in chat
- PRD detail view (optional)

### 10.2 Testing Strategy

```go
// Key test scenarios
func TestPRDGeneration(t *testing.T) {
    // Given a completed discovery with 3 MVP features
    // When GenerateAllPRDs is called
    // Then 3 PRDs should be created with status "generating"
    // And Claude should be called 3 times
    // And PRDs should transition to "draft" status
}

func TestAgentContextMatching(t *testing.T) {
    // Given an active PRD for "Order List"
    // When user says "add a search box to the order list"
    // Then the agent context should include the Order List PRD
    // And the Developer agent should be selected
}

func TestPRDStatusTransitions(t *testing.T) {
    // Test valid transitions: pending->generating->draft->ready->in_progress->complete
    // Test invalid transitions are rejected
}
```

---

## 11. Future Considerations

### 11.1 PRD Versioning
When users request changes to completed features, consider:
- PRD version history
- Change tracking
- Re-generation with context of original

### 11.2 Multi-Agent Collaboration
Future agents (QA, DevOps) may need:
- Different PRD views
- Specialized context extraction
- Cross-agent handoff protocols

### 11.3 PRD Templates
Industry-specific templates could accelerate generation:
- E-commerce features
- Booking systems
- Inventory management

---

## Document History

| Date | Version | Author | Changes |
|------|---------|--------|---------|
| 2025-12-26 | 1.0 | Architect Agent | Initial design specification |
