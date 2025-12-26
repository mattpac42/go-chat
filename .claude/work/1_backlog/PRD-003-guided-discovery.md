# PRD-003: Guided Discovery Flow

**Version**: 1.0
**Created**: 2025-12-26
**Author**: System Architect
**Status**: Planning
**Phase**: 2B - Guided Discovery
**Dependencies**: PRD-001 (Chat UI), PRD-002 (Backend API)

---

## Overview

### Problem Statement

Users currently jump straight into code generation without clearly articulating what they want to build. This leads to:
- Misaligned expectations between user vision and generated output
- Excessive back-and-forth refinement iterations
- Unclear project scope making it hard to prioritize features
- No structured data to seed the App Map with meaningful context

### Proposed Solution

A 5-stage conversational discovery flow led by a "Product Guide" agent that helps non-technical users articulate their product vision, user personas, and MVP scope BEFORE any code generation begins. The discovery data is then used to seed the App Map and inform subsequent development conversations.

### Goals

1. Achieve 85%+ discovery completion rate (users reaching "Start Building")
2. Capture structured product data (vision, personas, features) through natural conversation
3. Complete discovery in 5-10 minutes average
4. Seamlessly transition from discovery to development phase
5. Seed App Map with discovery data for intelligent file organization

### Non-Goals

- Formal requirements documentation generation
- Technical architecture decisions during discovery
- Code generation during the discovery flow
- Complex branching conversation trees
- Integration with external project management tools

---

## Architecture Overview

### System Context

```
+-------------------+     WebSocket      +-------------------+
|   React Frontend  |<----------------->|    Go Backend     |
|                   |                    |                   |
| - DiscoveryFlow   |   discovery_stage  | - DiscoveryService|
| - ProgressIndicator|   discovery_data  | - DiscoveryRepo   |
| - SummaryCard     |                    | - ClaudeService   |
|                   |                    |   (with discovery |
|                   |                    |    prompts)       |
+-------------------+                    +-------------------+
                                                  |
                                                  v
                                         +-------------------+
                                         |    PostgreSQL     |
                                         |                   |
                                         | - project_discovery|
                                         | - discovery_users |
                                         | - discovery_features|
                                         +-------------------+
```

### Integration with Existing Chat

The discovery flow integrates as a **modal conversation state** within the existing chat infrastructure:

1. **Project Creation**: When a new project is created, check for discovery state
2. **Discovery Mode**: If `discovery_stage != 'complete'`, chat operates in discovery mode
3. **Prompt Switching**: ChatService uses discovery-specific system prompts based on stage
4. **Stage Progression**: Backend tracks stage transitions based on Claude's responses
5. **Transition to Development**: On "Start Building", switch to default development prompts

```
┌─────────────────────────────────────────────────────────────┐
│                    Project Lifecycle                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  CREATE PROJECT                                              │
│       │                                                      │
│       v                                                      │
│  ┌─────────────┐                                            │
│  │  Discovery  │◄──── New projects start here               │
│  │   Mode      │                                            │
│  └──────┬──────┘                                            │
│         │                                                    │
│         │ Stage progression: welcome → problem → personas   │
│         │                    → mvp → summary                 │
│         │                                                    │
│         v                                                    │
│  ┌─────────────┐                                            │
│  │  Summary    │◄──── User confirms with "Start Building"   │
│  │  Confirmed  │                                            │
│  └──────┬──────┘                                            │
│         │                                                    │
│         v                                                    │
│  ┌─────────────┐                                            │
│  │ Development │◄──── Normal chat with code generation      │
│  │   Mode      │                                            │
│  └─────────────┘                                            │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Backend Changes

### 1. Discovery State Model

**File**: `/workspace/backend/internal/model/discovery.go`

```go
package model

import (
    "time"
    "github.com/google/uuid"
)

// DiscoveryStage represents the current stage in the discovery flow
type DiscoveryStage string

const (
    StageWelcome  DiscoveryStage = "welcome"
    StageProblem  DiscoveryStage = "problem"
    StagePersonas DiscoveryStage = "personas"
    StageMVP      DiscoveryStage = "mvp"
    StageSummary  DiscoveryStage = "summary"
    StageComplete DiscoveryStage = "complete"
)

// ProjectDiscovery stores the discovery state and captured data for a project
type ProjectDiscovery struct {
    ID               uuid.UUID      `db:"id" json:"id"`
    ProjectID        uuid.UUID      `db:"project_id" json:"projectId"`
    Stage            DiscoveryStage `db:"stage" json:"stage"`
    StageStartedAt   time.Time      `db:"stage_started_at" json:"stageStartedAt"`

    // Captured data from conversation
    BusinessContext  *string `db:"business_context" json:"businessContext,omitempty"`
    ProblemStatement *string `db:"problem_statement" json:"problemStatement,omitempty"`
    Goals            []string `db:"-" json:"goals,omitempty"` // Stored as JSONB

    // Summary fields (populated in summary stage)
    ProjectName      *string `db:"project_name" json:"projectName,omitempty"`
    SolvesStatement  *string `db:"solves_statement" json:"solvesStatement,omitempty"`

    // Metadata
    IsReturningUser  bool       `db:"is_returning_user" json:"isReturningUser"`
    UsedTemplateID   *uuid.UUID `db:"used_template_id" json:"usedTemplateId,omitempty"`
    ConfirmedAt      *time.Time `db:"confirmed_at" json:"confirmedAt,omitempty"`

    CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
    UpdatedAt        time.Time  `db:"updated_at" json:"updatedAt"`
}

// DiscoveryUser represents a user persona defined during discovery
type DiscoveryUser struct {
    ID             uuid.UUID `db:"id" json:"id"`
    DiscoveryID    uuid.UUID `db:"discovery_id" json:"discoveryId"`
    Description    string    `db:"description" json:"description"`
    Count          int       `db:"user_count" json:"count"`
    HasPermissions bool      `db:"has_permissions" json:"hasPermissions"`
    PermissionNotes *string  `db:"permission_notes" json:"permissionNotes,omitempty"`
    CreatedAt      time.Time `db:"created_at" json:"createdAt"`
}

// DiscoveryFeature represents a feature captured during MVP scoping
type DiscoveryFeature struct {
    ID          uuid.UUID `db:"id" json:"id"`
    DiscoveryID uuid.UUID `db:"discovery_id" json:"discoveryId"`
    Name        string    `db:"name" json:"name"`
    Priority    int       `db:"priority" json:"priority"`
    Version     string    `db:"version" json:"version"` // "v1", "v2", etc.
    CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

// DiscoveryEditHistory tracks edits made to discovery data
type DiscoveryEditHistory struct {
    ID            uuid.UUID `db:"id" json:"id"`
    DiscoveryID   uuid.UUID `db:"discovery_id" json:"discoveryId"`
    Stage         string    `db:"stage" json:"stage"`
    FieldEdited   string    `db:"field_edited" json:"fieldEdited"`
    OriginalValue string    `db:"original_value" json:"originalValue"`
    NewValue      string    `db:"new_value" json:"newValue"`
    EditedAt      time.Time `db:"edited_at" json:"editedAt"`
}

// DiscoverySummary is the combined view shown to users
type DiscoverySummary struct {
    ProjectName     string             `json:"projectName"`
    SolvesStatement string             `json:"solvesStatement"`
    Users           []DiscoveryUser    `json:"users"`
    MVPFeatures     []DiscoveryFeature `json:"mvpFeatures"`
    FutureFeatures  []DiscoveryFeature `json:"futureFeatures"`
}
```

### 2. Database Schema

**File**: `/workspace/backend/migrations/004_discovery.sql`

```sql
-- 004_discovery.sql
-- Discovery flow tables for guided project setup

-- Main discovery state table
CREATE TABLE IF NOT EXISTS project_discovery (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    stage VARCHAR(20) NOT NULL DEFAULT 'welcome'
        CHECK (stage IN ('welcome', 'problem', 'personas', 'mvp', 'summary', 'complete')),
    stage_started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Captured conversation data
    business_context TEXT,
    problem_statement TEXT,
    goals JSONB DEFAULT '[]'::JSONB,

    -- Summary data
    project_name VARCHAR(255),
    solves_statement TEXT,

    -- Metadata
    is_returning_user BOOLEAN DEFAULT FALSE,
    used_template_id UUID,
    confirmed_at TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT unique_project_discovery UNIQUE (project_id)
);

-- Discovery users (personas)
CREATE TABLE IF NOT EXISTS discovery_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discovery_id UUID NOT NULL REFERENCES project_discovery(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    user_count INTEGER DEFAULT 1,
    has_permissions BOOLEAN DEFAULT FALSE,
    permission_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Discovery features (MVP and future)
CREATE TABLE IF NOT EXISTS discovery_features (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discovery_id UUID NOT NULL REFERENCES project_discovery(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    version VARCHAR(10) NOT NULL DEFAULT 'v1',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Edit history for tracking changes
CREATE TABLE IF NOT EXISTS discovery_edit_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discovery_id UUID NOT NULL REFERENCES project_discovery(id) ON DELETE CASCADE,
    stage VARCHAR(20) NOT NULL,
    field_edited VARCHAR(100) NOT NULL,
    original_value TEXT,
    new_value TEXT,
    edited_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_discovery_project ON project_discovery(project_id);
CREATE INDEX IF NOT EXISTS idx_discovery_stage ON project_discovery(stage);
CREATE INDEX IF NOT EXISTS idx_discovery_users_discovery ON discovery_users(discovery_id);
CREATE INDEX IF NOT EXISTS idx_discovery_features_discovery ON discovery_features(discovery_id);
CREATE INDEX IF NOT EXISTS idx_discovery_features_version ON discovery_features(version);

-- Comments
COMMENT ON TABLE project_discovery IS 'Stores discovery flow state and captured data for each project';
COMMENT ON TABLE discovery_users IS 'User personas identified during discovery';
COMMENT ON TABLE discovery_features IS 'Features identified during MVP scoping with version assignments';
```

### 3. API Endpoints

**New endpoints for discovery management:**

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/projects/:id/discovery` | Get current discovery state |
| PUT | `/api/projects/:id/discovery/stage` | Advance to next stage |
| PUT | `/api/projects/:id/discovery/data` | Update discovery data fields |
| POST | `/api/projects/:id/discovery/users` | Add a user persona |
| POST | `/api/projects/:id/discovery/features` | Add a feature |
| POST | `/api/projects/:id/discovery/confirm` | Confirm and complete discovery |
| DELETE | `/api/projects/:id/discovery` | Reset discovery (start over) |

**Handler structure** (`/workspace/backend/internal/handler/discovery.go`):

```go
package handler

// DiscoveryHandler handles discovery-related HTTP endpoints
type DiscoveryHandler struct {
    service *service.DiscoveryService
    logger  zerolog.Logger
}

// Routes to register:
// - GET    /projects/:id/discovery          -> GetDiscovery
// - PUT    /projects/:id/discovery/stage    -> UpdateStage
// - PUT    /projects/:id/discovery/data     -> UpdateData
// - POST   /projects/:id/discovery/users    -> AddUser
// - POST   /projects/:id/discovery/features -> AddFeature
// - POST   /projects/:id/discovery/confirm  -> ConfirmDiscovery
// - DELETE /projects/:id/discovery          -> ResetDiscovery
```

### 4. Discovery Service

**File**: `/workspace/backend/internal/service/discovery.go`

```go
package service

// DiscoveryService orchestrates the discovery flow
type DiscoveryService struct {
    repo            repository.DiscoveryRepository
    projectRepo     repository.ProjectRepository
    claudeService   *ClaudeService
    logger          zerolog.Logger
}

// Key methods:

// GetOrCreateDiscovery returns existing discovery or creates new one
func (s *DiscoveryService) GetOrCreateDiscovery(ctx context.Context, projectID uuid.UUID) (*model.ProjectDiscovery, error)

// AdvanceStage moves to the next discovery stage
func (s *DiscoveryService) AdvanceStage(ctx context.Context, discoveryID uuid.UUID) (*model.ProjectDiscovery, error)

// ExtractDiscoveryData parses Claude response for structured data
func (s *DiscoveryService) ExtractDiscoveryData(ctx context.Context, discoveryID uuid.UUID, response string) error

// GetSummary builds the complete summary view
func (s *DiscoveryService) GetSummary(ctx context.Context, discoveryID uuid.UUID) (*model.DiscoverySummary, error)

// ConfirmDiscovery marks discovery as complete and seeds App Map
func (s *DiscoveryService) ConfirmDiscovery(ctx context.Context, discoveryID uuid.UUID) error

// SeedAppMap creates initial functional groups from discovery data
func (s *DiscoveryService) SeedAppMap(ctx context.Context, projectID uuid.UUID, features []model.DiscoveryFeature) error
```

### 5. Claude Prompt Modifications

The system prompt changes based on discovery stage. A new prompt builder will be introduced:

**File**: `/workspace/backend/internal/service/prompts/discovery.go`

```go
package prompts

// DiscoveryPromptBuilder creates stage-appropriate system prompts
type DiscoveryPromptBuilder struct{}

func (b *DiscoveryPromptBuilder) Build(stage model.DiscoveryStage, context *DiscoveryContext) string {
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
        return DefaultSystemPrompt() // Fall back to development mode
    }
}

func (b *DiscoveryPromptBuilder) welcomePrompt() string {
    return `You are the Product Guide for Go Chat. Your role is to help users articulate
what they want to build through friendly conversation.

CURRENT STAGE: Welcome (1 of 5)

YOUR TASK:
1. Warmly greet the user
2. Set expectations that this will take "a few minutes"
3. Ask an open-ended question about what they do or their business

STYLE GUIDELINES:
- Use warm, encouraging language
- No technical jargon whatsoever
- Keep responses concise (2-4 sentences)
- End with an open-ended question

DO NOT:
- Generate any code
- Mention programming languages or frameworks
- Use bullet points in your greeting
- Ask yes/no questions

Example opening:
"Welcome! I'm here to help you turn your idea into a working application. Before we start
building, let's take a few minutes to understand exactly what you need. First, tell me a
bit about yourself - what do you do?"

After the user responds, I will advance the conversation to the next stage.`
}

// Additional prompt methods for each stage...
```

**Stage progression detection** - Claude's response will include metadata markers that the backend parses:

```go
// Response metadata format (hidden from user display):
// <!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"business_context":"bakery owner"}}-->

func (s *DiscoveryService) ParseResponseMetadata(response string) (*DiscoveryMetadata, string) {
    // Extract metadata, return clean response for display
}
```

### 6. Chat Service Integration

Modify `ChatService.ProcessMessage` to handle discovery mode:

```go
func (s *ChatService) ProcessMessage(ctx context.Context, projectID uuid.UUID, content string, onChunk func(string)) (*ChatResult, error) {
    // Check discovery state
    discovery, err := s.discoveryService.GetOrCreateDiscovery(ctx, projectID)
    if err != nil {
        return nil, err
    }

    // Select appropriate system prompt
    var systemPrompt string
    if discovery.Stage != model.StageComplete {
        systemPrompt = s.promptBuilder.Build(discovery.Stage, &DiscoveryContext{
            Discovery: discovery,
        })
    } else {
        systemPrompt = DefaultSystemPrompt()
    }

    // Continue with existing flow using appropriate prompt...
}
```

---

## Frontend Changes

### 1. New Components

#### DiscoveryProgress Component

**File**: `/workspace/frontend/src/components/discovery/DiscoveryProgress.tsx`

Purpose: Show the 5-stage progress indicator in the chat header.

```typescript
interface DiscoveryProgressProps {
  currentStage: DiscoveryStage;
  totalStages: number;
  onStageClick?: (stage: DiscoveryStage) => void; // For mobile drawer
  isMobile?: boolean;
}

// Desktop: Horizontal dots with current stage label
// Mobile: Compact dots with "X/5" indicator, tappable for drawer
```

**Visual states:**
- Completed stages: Filled dots (teal)
- Current stage: Filled dot + label
- Future stages: Empty dots (gray outline)

#### DiscoverySummaryCard Component

**File**: `/workspace/frontend/src/components/discovery/DiscoverySummaryCard.tsx`

Purpose: Display the final summary with edit/confirm actions.

```typescript
interface DiscoverySummaryCardProps {
  summary: DiscoverySummary;
  onEdit: () => void;
  onConfirm: () => void;
  isConfirming?: boolean;
}

// Renders:
// - Project name
// - "Solves" statement
// - User personas list
// - MVP features list
// - Future features list
// - Edit Details button (ghost)
// - Start Building button (primary)
```

**Responsive behavior:**
- Desktop: 2-column grid layout
- Mobile: Stacked single column

#### DiscoveryStageDrawer Component (Mobile)

**File**: `/workspace/frontend/src/components/discovery/DiscoveryStageDrawer.tsx`

Purpose: Bottom sheet showing detailed stage progress on mobile.

```typescript
interface DiscoveryStageDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  stages: StageInfo[];
  currentStage: DiscoveryStage;
}
```

### 2. State Management

#### Discovery Hook

**File**: `/workspace/frontend/src/hooks/useDiscovery.ts`

```typescript
interface UseDiscoveryReturn {
  discovery: ProjectDiscovery | null;
  isDiscoveryMode: boolean;
  currentStage: DiscoveryStage;
  stageProgress: number; // 1-5
  summary: DiscoverySummary | null;
  isLoading: boolean;
  error: string | null;

  // Actions
  advanceStage: () => Promise<void>;
  updateData: (data: Partial<DiscoveryData>) => Promise<void>;
  confirmDiscovery: () => Promise<void>;
  resetDiscovery: () => Promise<void>;
}

export function useDiscovery(projectId: string): UseDiscoveryReturn {
  // Fetch discovery state on mount
  // Subscribe to discovery updates via WebSocket
  // Provide action methods for stage management
}
```

#### Extended Chat State

Modify existing `ChatState` to include discovery awareness:

```typescript
interface ChatState {
  messages: Message[];
  isLoading: boolean;
  error: string | null;
  // New discovery fields
  discoveryStage: DiscoveryStage | null;
  isDiscoveryMode: boolean;
}
```

### 3. ChatContainer Integration

Modify `/workspace/frontend/src/components/chat/ChatContainer.tsx`:

```typescript
export function ChatContainer({ projectId, ... }) {
  const { discovery, isDiscoveryMode, currentStage } = useDiscovery(projectId);

  return (
    <div className="flex flex-col h-full bg-white">
      <header className="...">
        {/* Existing header content */}

        {/* Discovery progress - only show in discovery mode */}
        {isDiscoveryMode && (
          <DiscoveryProgress
            currentStage={currentStage}
            totalStages={5}
          />
        )}
      </header>

      {/* Messages - unchanged */}
      <MessageList messages={messages} isLoading={isLoading} />

      {/* Input - placeholder changes based on stage */}
      <ChatInput
        placeholder={getDiscoveryPlaceholder(currentStage)}
        ...
      />
    </div>
  );
}
```

### 4. Message Rendering Updates

The `MessageBubble` component needs to handle:

1. **Product Guide avatar**: Show distinct avatar for discovery messages
2. **Summary card injection**: Render `DiscoverySummaryCard` inline when stage is `summary`
3. **Transition animation**: "Discovery Complete" banner on confirmation

### 5. New Types

**File**: `/workspace/frontend/src/types/discovery.ts`

```typescript
export type DiscoveryStage =
  | 'welcome'
  | 'problem'
  | 'personas'
  | 'mvp'
  | 'summary'
  | 'complete';

export interface ProjectDiscovery {
  id: string;
  projectId: string;
  stage: DiscoveryStage;
  stageStartedAt: string;
  businessContext?: string;
  problemStatement?: string;
  goals?: string[];
  projectName?: string;
  solvesStatement?: string;
  isReturningUser: boolean;
  confirmedAt?: string;
}

export interface DiscoveryUser {
  id: string;
  description: string;
  count: number;
  hasPermissions: boolean;
  permissionNotes?: string;
}

export interface DiscoveryFeature {
  id: string;
  name: string;
  priority: number;
  version: string;
}

export interface DiscoverySummary {
  projectName: string;
  solvesStatement: string;
  users: DiscoveryUser[];
  mvpFeatures: DiscoveryFeature[];
  futureFeatures: DiscoveryFeature[];
}
```

---

## Implementation Phases

### Phase 1: Core Discovery Infrastructure (Backend)
**Estimated effort**: 3-4 days
**Priority**: P0

**Tasks**:
1. Create database migration `004_discovery.sql`
2. Implement `model/discovery.go` with all types
3. Create `repository/discovery.go` with CRUD operations
4. Implement `service/discovery.go` with stage management
5. Create `prompts/discovery.go` with stage-specific prompts
6. Add discovery check to `ChatService.ProcessMessage`
7. Unit tests for discovery service

**Deliverables**:
- Database tables created and migrated
- Discovery state tracked per project
- Stage-appropriate prompts sent to Claude
- Tests passing

**Exit Criteria**:
- New project triggers welcome stage
- Chat messages use discovery prompts when in discovery mode
- Stage advances manually via API call

---

### Phase 2: Discovery API & Data Extraction (Backend)
**Estimated effort**: 2-3 days
**Priority**: P0

**Tasks**:
1. Implement `handler/discovery.go` with all endpoints
2. Add response metadata parsing for structured data extraction
3. Implement automatic stage advancement logic
4. Add user and feature storage methods
5. Create summary generation logic
6. Wire up "confirm" endpoint to mark complete
7. Integration tests for full flow

**Deliverables**:
- All discovery API endpoints functional
- Claude responses parsed for structured data
- Stage auto-advances based on conversation
- Summary endpoint returns complete data

**Exit Criteria**:
- Can complete full discovery flow via API
- Structured data extracted and stored
- Summary matches captured data

---

### Phase 3: Frontend Discovery Components (Frontend)
**Estimated effort**: 3-4 days
**Priority**: P1

**Tasks**:
1. Create `types/discovery.ts` with TypeScript types
2. Implement `hooks/useDiscovery.ts` for state management
3. Build `DiscoveryProgress` component (desktop + mobile)
4. Build `DiscoverySummaryCard` component (responsive)
5. Build `DiscoveryStageDrawer` for mobile
6. Integrate progress indicator into `ChatContainer` header
7. Update `ChatInput` placeholder based on stage

**Deliverables**:
- Progress indicator shows in header during discovery
- Summary card renders inline in chat
- Mobile-optimized drawer for stage details
- Responsive layouts for all components

**Exit Criteria**:
- Visual progress visible on all viewports
- Summary card actions work (edit/confirm)
- Smooth transitions between stages

---

### Phase 4: Integration & Polish (Full Stack)
**Estimated effort**: 2-3 days
**Priority**: P1

**Tasks**:
1. Connect frontend discovery hook to backend API
2. Implement WebSocket events for discovery updates
3. Add "Discovery Complete" transition animation
4. Handle edit flow (tap section to modify)
5. Implement session recovery (continue incomplete discovery)
6. Add returning user detection and fast-track option
7. End-to-end testing

**Deliverables**:
- Complete discovery flow works end-to-end
- Real-time updates via WebSocket
- Session persistence for incomplete discoveries
- Returning user experience

**Exit Criteria**:
- User can complete discovery in 5-10 minutes
- Data seeds App Map functional groups
- Returning users see fast-track option

---

## Data Flow

### Discovery to App Map Seeding

```
DISCOVERY COMPLETION FLOW
=========================

1. User clicks "Start Building"
        │
        v
2. POST /api/projects/:id/discovery/confirm
        │
        v
3. DiscoveryService.ConfirmDiscovery()
        ├── Mark discovery.stage = 'complete'
        ├── Set confirmed_at timestamp
        └── Call SeedAppMap()
                │
                v
4. SeedAppMap() creates functional groups:
        │
        │  discovery_features (v1)     →    file_metadata.functional_group
        │  ├── "Order list"            →    "Order Management"
        │  ├── "Order form"            →    "Order Management"
        │  └── "Due date tracking"     →    "Scheduling"
        │
        v
5. Update project title from discovery.project_name
        │
        v
6. Return success → Frontend shows transition
        │
        v
7. ChatService now uses DefaultSystemPrompt()
   (development mode with functional groups context)
```

### Conversation Flow Data Capture

```
STAGE PROGRESSION WITH DATA EXTRACTION
======================================

STAGE 1: Welcome
├── User: "I run a bakery"
└── Extract: business_context = "bakery owner"
        │
        v
STAGE 2: Problem
├── User: "Order tracking is chaos"
├── Extract: problem_statement = "order tracking chaos with paper/WhatsApp"
└── Extract: goals = ["centralized order tracking", "reduce lost orders"]
        │
        v
STAGE 3: Personas
├── User: "Me and 2 employees"
├── Extract: users[] = [
│       { description: "owner/baker", count: 1, hasPermissions: true },
│       { description: "order takers", count: 2, hasPermissions: false }
│   ]
└── Extract: permission_notes = "employees: orders only; owner: full access"
        │
        v
STAGE 4: MVP
├── User: "Order list, order form, due dates"
├── Extract: features[] = [
│       { name: "Order list", priority: 1, version: "v1" },
│       { name: "Order form", priority: 2, version: "v1" },
│       { name: "Due date tracking", priority: 3, version: "v1" }
│   ]
├── User: "Calendar would be nice"
└── Extract: features[] += { name: "Calendar view", priority: 1, version: "v2" }
        │
        v
STAGE 5: Summary
├── Generate: project_name = "Cake Order Manager"
├── Generate: solves_statement = "Replaces paper and WhatsApp..."
└── Display: DiscoverySummaryCard with all data
        │
        v
CONFIRM → COMPLETE
└── Seed App Map with functional groups derived from features
```

---

## Success Metrics

| Metric | Target | Measurement Method |
|--------|--------|-------------------|
| Discovery completion rate | 85%+ | `confirmed_at IS NOT NULL` / total discoveries |
| Average completion time | 5-10 min | `confirmed_at - created_at` |
| Stage drop-off rate | <10% per stage | Count by stage where `stage != 'complete'` |
| Summary edit frequency | <30% | `discovery_edit_history` count |
| Returning user fast-track adoption | 60%+ | `used_template_id IS NOT NULL` |

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Claude doesn't follow stage prompts | Medium | High | Strict prompt engineering, fallback detection |
| Users skip or abandon discovery | Medium | Medium | Clear value proposition, time estimates |
| Data extraction fails | Medium | Medium | Manual fallback, edit capability |
| Mobile UX feels cramped | Low | Medium | Tested mobile-first design |
| Performance with large conversations | Low | Low | Limit context window, use summary |

---

## Open Questions

1. **Template sharing**: Should users be able to share discovery templates with others?
2. **AI summary generation**: Should Claude auto-generate project name and solves statement, or should user input these directly?
3. **Partial saves**: If user leaves mid-discovery, how long do we retain partial data?
4. **Skip option**: Should experienced users have a way to skip discovery entirely?

---

## Appendix: File Locations

### Backend Files to Create

```
/workspace/backend/
├── migrations/
│   └── 004_discovery.sql                    # Database schema
├── internal/
│   ├── model/
│   │   └── discovery.go                     # Discovery domain models
│   ├── repository/
│   │   ├── discovery.go                     # Discovery repository
│   │   └── discovery_test.go               # Repository tests
│   ├── service/
│   │   ├── discovery.go                     # Discovery service
│   │   ├── discovery_test.go               # Service tests
│   │   └── prompts/
│   │       └── discovery.go                 # Stage-specific prompts
│   └── handler/
│       ├── discovery.go                     # HTTP handlers
│       └── discovery_test.go               # Handler tests
```

### Frontend Files to Create

```
/workspace/frontend/src/
├── types/
│   └── discovery.ts                         # TypeScript types
├── hooks/
│   └── useDiscovery.ts                      # Discovery state hook
├── components/
│   └── discovery/
│       ├── DiscoveryProgress.tsx            # Progress indicator
│       ├── DiscoverySummaryCard.tsx         # Summary with actions
│       ├── DiscoveryStageDrawer.tsx         # Mobile stage drawer
│       └── DiscoveryTransition.tsx          # Completion animation
└── lib/
    └── discovery-api.ts                     # API client functions
```

### Files to Modify

```
/workspace/backend/
├── internal/service/
│   └── chat.go                              # Add discovery mode check
├── cmd/server/
│   └── main.go                              # Wire up discovery handler

/workspace/frontend/src/
├── components/chat/
│   ├── ChatContainer.tsx                    # Add discovery progress
│   ├── MessageBubble.tsx                    # Handle summary card
│   └── ChatInput.tsx                        # Dynamic placeholder
└── types/
    └── index.ts                             # Export discovery types
```

---

**Document History**

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2025-12-26 | 1.0 | Initial implementation plan | System Architect |
