# Skills Manifest

> Single source of truth for all available skills. Agents should reference this file to discover workflows.

## Quick Reference

| Skill | Trigger Phrases | Next Skill |
|-------|-----------------|------------|
| product-vision | "create a vision", "define product strategy" | feature-roadmap |
| feature-roadmap | "create roadmap", "break down the vision" | prd-create |
| prd-create | "create PRD", "write requirements" | task-generate |
| task-generate | "generate tasks", "break down PRD" | task-process |
| task-process | "start implementing", "work on tasks" | - |

## Workflow Skills (User-Invoked)

### Product Discovery Flow

```
product-vision → feature-roadmap → prd-create → task-generate → task-process
```

#### product-vision
- **Purpose**: Create product vision through structured discovery
- **Triggers**: "create a vision", "define product strategy", "what should we build"
- **Output**: `.claude/work/0_vision/product-vision.md`, `.claude/work/0_vision/strategic-themes.md`
- **Next**: feature-roadmap

#### feature-roadmap
- **Purpose**: Decompose product vision into prioritized feature roadmap
- **Triggers**: "create roadmap", "break down the vision", "plan features"
- **Requires**: Approved product vision
- **Output**: `.claude/work/0_vision/feature-roadmap.md`
- **Next**: prd-create

#### prd-create
- **Purpose**: Create Product Requirements Document for a feature
- **Triggers**: "create PRD", "write requirements", "document this feature", "PRD for [feature]"
- **Output**: `.claude/work/1_backlog/[NNN]-[feature-name]/prd-[feature-name].md`
- **Next**: task-generate

#### task-generate
- **Purpose**: Generate task breakdown from PRD with agent assignments
- **Triggers**: "generate tasks", "break down PRD", "create task list", "what needs to be done"
- **Requires**: Approved PRD
- **Output**: `.claude/work/1_backlog/[NNN]-[feature-name]/tasks-[feature-name].md`
- **Next**: task-process

#### task-process
- **Purpose**: Process task list through implementation workflow
- **Triggers**: "start implementing", "work on tasks", "execute task list", "let's build this"
- **Requires**: Approved task list
- **Action**: Moves work to `.claude/work/2_active/`, assigns to agents

### Project Setup

#### plant-project
- **Purpose**: Root a new project from The Garden
- **Triggers**: "/plant", "create new project", "start new repo", "initialize project"
- **Output**: New project directory with Garden structure

#### setup-validation
- **Purpose**: Verify Garden environment is correctly configured
- **Triggers**: Project init, setup issues, "check setup status"
- **Output**: Validation report with remediation steps

#### workspace-setup
- **Purpose**: Configure VSCode workspace colors for visual project identification
- **Triggers**: User request for workspace customization
- **Output**: `.vscode/settings.json` with color theme

#### hooks-setup
- **Purpose**: Configure Claude Code hooks for context tracking and notifications
- **Triggers**: One-time setup request
- **Output**: Hook configuration in settings

### Execution State

#### beads
- **Purpose**: Git-backed issue tracking for PRD-based work
- **Triggers**: "beads", "what should I work on", "session context", "track issue"
- **Commands**:
  - `beads init` - Initialize .beads/ directory
  - `beads context` - Session start - shows in-progress, ready, blocked beads
  - `beads add "title"` - Create new bead
  - `beads progress <id>` - Mark as in-progress
  - `beads close <id>` - Complete a bead
  - `beads list --ready` - Show unblocked work
  - `beads import <file>` - Import from PRD tasks
- **Storage**: `.beads/issues.jsonl`
- **Use for**: PRD-based work, task tracking, long-running features

#### handoff
- **Purpose**: Create session handoff files for ad-hoc work
- **Triggers**: "handoff", "end session", "save progress", context at 75%+
- **Output**: `.claude/work/history/SESSION-[NNN].md`, `.claude/work/history/HANDOFF-SESSION.md`
- **Use for**: Ad-hoc work, quick sessions, exploratory work without PRDs

#### catch-up
- **Purpose**: Restore session context from previous handoff
- **Triggers**: Session start with HANDOFF-SESSION.md, "catch up", "what was I working on"
- **Reads**: `.claude/work/history/HANDOFF-SESSION.md`
- **Use for**: Resuming ad-hoc sessions

#### context-display
- **Purpose**: Show token usage visualization
- **Triggers**: Context exceeds 50%, agent task completion
- **Output**: Visual context bar in response

## Session Management Decision

| Scenario | Use |
|----------|-----|
| PRD-based work with tasks | beads |
| Ad-hoc work, quick fixes | handoff/catch-up |
| Exploratory sessions | handoff/catch-up |
| Long-running features | beads |

## Auto-Invoked Skills

These skills trigger automatically without user action:

| Skill | Trigger Condition |
|-------|-------------------|
| beads context | Session start with .beads/ (PRD work) |
| catch-up | Session start with HANDOFF-SESSION.md (ad-hoc work) |
| context-display | Every response when context > 50% |
| setup-validation | Project initialization |

## Work Directory Structure

All skills use this folder structure:

```
.claude/work/
├── 0_vision/    # Strategic vision documents
├── 1_backlog/   # PRDs and task lists ready for development
├── 2_active/    # Work in progress
├── 3_done/      # Completed work
└── history/     # Session summaries and handoff files (ad-hoc work)

.beads/
└── issues.jsonl # Execution state (PRD-based work)
```

## Agent-Skill Mapping

| Agent | Relevant Skills |
|-------|-----------------|
| product | product-vision, feature-roadmap, prd-create |
| product-visionary | product-vision, feature-roadmap |
| developer | task-process |
| architect | task-generate (reviews), task-process |
| platform | task-process |
| researcher | task-generate (analysis) |

## Skill Invocation

Skills can be invoked by:
1. **User triggers**: Phrases matching skill description
2. **Skill chaining**: Previous skill suggests next step
3. **Agent reference**: Agents consult this manifest for available workflows
4. **Direct command**: `/skill-name` syntax

## Adding New Skills

When creating skills, update this manifest:
1. Add to appropriate section
2. Include triggers, output, and next skill
3. Update agent-skill mapping if relevant
