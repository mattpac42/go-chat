# The Garden - Complete System Context

## Purpose

**The Garden** is a baseline template for Claude Code agent orchestration. It provides:
- Agent orchestration framework (main agent delegates, specialized agents implement)
- Automatic context management with session handoff
- PRD-driven development workflow
- 130+ specialized agents for production work

## Core Philosophy

### 1. Main Agent as Orchestrator
The main agent **delegates 100% of implementation** to specialized agents:
- Expert-level quality through domain specialization
- Context efficiency through focused agent scopes
- Parallel execution of independent tasks

### 2. Documentation-First Development
All features begin with PRDs moving through a structured lifecycle:
- `0_vision/` → Strategic vision
- `1_backlog/` → PRDs and task planning
- `2_active/` → Implementation in progress
- `3_done/` → Completed with retrospectives

### 3. Test-Driven Development
Mandatory TDD for all code:
1. Write failing test
2. Implement minimal code
3. Verify tests pass
4. Refactor if needed
5. Commit only when green

## System Architecture

### Agent Orchestration Model

```
┌─────────────────────────────────────────────────┐
│              Main Claude Agent                   │
│             (Orchestrator Only)                  │
│                                                  │
│  • Understands user requests                     │
│  • Delegates to specialized agents               │
│  • Runs agents in parallel when possible         │
│  • Integrates outputs and tracks progress        │
│  • NEVER implements directly                     │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│           Specialized Agents                     │
├─────────────────────────────────────────────────┤
│                                                  │
│  Generic (.claude/agents/) - 9 agents:          │
│  • developer, architect, product, platform      │
│  • researcher, garden-guide, project-navigator  │
│  • product-visionary, prompt-optimizer          │
│                                                  │
│  Specialized (marketplace/agents/) - 130+ agents:    │
│  • software-tactical / software-strategic       │
│  • platform-tactical / platform-strategic       │
│  • product-manager-tactical / strategic         │
│  • ux-tactical / ux-strategic                   │
│  • And many more domain experts                 │
│                                                  │
└─────────────────────────────────────────────────┘
```

### Agent Selection

| Library | Location | Use When |
|---------|----------|----------|
| Generic | `.claude/agents/` | Quick tasks, prototyping, simple projects |
| Specialized | `marketplace/agents/` | Production work, complex features, clear phases |

## Work Management

### Folder Structure
```
.claude/work/
├── 0_vision/    # Strategic vision documents
├── 1_backlog/   # PRDs ready for development
├── 2_active/    # Work in progress
├── 3_done/      # Completed features
└── history/     # Session logs and handoffs
    ├── SESSION-001.md
    ├── SESSION-002.md
    └── HANDOFF-SESSION.md
```

### PRD Lifecycle Workflow

```
┌─────────────────────────────────────┐
│     1. PRD CREATION (1_backlog/)    │
│                                     │
│  • User describes feature           │
│  • Main agent asks questions        │
│  • Creates feature folder           │
│  • Generates PRD document           │
└──────────────┬──────────────────────┘
               ▼
┌─────────────────────────────────────┐
│     2. TASK GENERATION              │
│                                     │
│  • Break into 3-7 parent tasks      │
│  • Assign specialized agents        │
│  • Identify files to create/modify  │
└──────────────┬──────────────────────┘
               ▼
┌─────────────────────────────────────┐
│     3. IMPLEMENTATION (2_active/)   │
│                                     │
│  • Move folder to 2_active/         │
│  • Delegate tasks to agents         │
│  • Follow TDD workflow              │
│  • Validate quality gates           │
└──────────────┬──────────────────────┘
               ▼
┌─────────────────────────────────────┐
│     4. COMPLETION (3_done/)         │
│                                     │
│  • All tasks complete, tests pass   │
│  • Create retrospective             │
│  • Move folder to 3_done/           │
└─────────────────────────────────────┘
```

## Context Management

### Thresholds
| Level | Action |
|-------|--------|
| 50%+ | Context bar displayed |
| 60% | Warning - approaching limit |
| 75% | Handoff triggered automatically |
| 85% | New session recommended |

### Display Format
```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 50% (100k/200k)
```
- Blocks 1-12: Green (0-60%)
- Blocks 13-15: Yellow (60-75%)
- Blocks 16-17: Orange (75-85%)
- Blocks 18-20: Red (85-100%)

### Session Handoff

At 75% context, the system creates:
1. `SESSION-[XXX].md` - Complete session summary
2. `HANDOFF-SESSION.md` - Context for next session

The `catch-up` skill restores context at session start.

## Skills and Commands

### Skills (Auto-Invoked)
| Skill | Trigger | Purpose |
|-------|---------|---------|
| catch-up | Session start | Restore context from HANDOFF-SESSION.md |
| agent-session-summary | Agent exit | Document agent work |
| context-display | 50%+ or agent completion | Show usage bar |
| handoff | 75% threshold | Create session files |
| setup-validation | Project init | Verify environment |

### Commands (User-Invoked)
| Command | Purpose |
|---------|---------|
| /handoff | Create session handoff files |
| /commit | Git commit workflow |
| /mr | Create merge request |
| /onboard | Initialize workspace |
| /catch-up | Restore session context |

## File Locations

```
.claude/
├── PROTOCOLS.md      # Implementation details
├── QUICKSTART.md     # Setup guide
├── PROJECT.md        # Project context template
├── agents/           # 9 generic agents
├── skills/           # 13 workflow skills
├── commands/         # User-invoked commands
├── templates/        # Workflow templates
├── work/             # PRD lifecycle folders
├── config/           # Plugin configuration
├── docs/             # Technical documentation
└── settings.json     # Claude Code configuration

marketplace/
├── agents/           # 130+ specialized agents
├── plugins/          # Extension plugins
├── skills/           # Community skills
└── knowledge-base/   # Reference materials
```

## Quality Standards

### Delegation Enforcement
Main agent NEVER implements - 100% delegation required:
- Reading context: ✅ Allowed
- Clarifying questions: ✅ Allowed
- Invoking agents: ✅ Allowed
- Writing code: ❌ Delegate to developer
- Creating files: ❌ Delegate to appropriate agent

### TDD Quality Gate
Before marking coding tasks complete:
- Unit tests exist (>80% coverage)
- All tests pass
- Test output documented
- Tests committed with code

## Downstream Project Setup

### Option 1: Plant from The Garden (Recommended)

From The Garden repository:
```bash
/plant
```

The wizard guides you through:
1. Project type (webapp, api, cli, mobile, library, data, devops, business)
2. Project name and location
3. Agent selection (only copies what you need)
4. Generation with `lineage.json` for sync

Then in your new project:
```bash
/onboard
```

### Option 2: Manual Setup
```bash
# Copy system to your project
cp -r .claude/ /path/to/project/
cp CLAUDE.md /path/to/project/

# Run onboarding
cd /path/to/project/
/onboard

# Create project context
# Edit .claude/PROJECT.md with your details
```

### Sync Updates from The Garden

Projects created with `/plant` can pull updates:
```bash
/sync-baseline
```

This compares your project against The Garden and offers selective updates for agents, skills, commands, and templates.

### What Gets Created at Runtime
- `.claude/work/history/` - Session files
- `.vscode/settings.json` - Workspace colors (from /onboard)
- `.claude/lineage.json` - Garden connection (from /plant)

## Benefits

### For Developers
- Structured workflow: idea → PRD → tasks → implementation → done
- Quality enforcement: TDD ensures working code
- Context efficiency: Specialized agents = focused work
- Session continuity: Never lose context

### For Teams
- Consistent process: Everyone follows same workflow
- Knowledge preservation: Agent history captures decisions
- Faster onboarding: Read history to understand project

---

**Version**: Garden 2.0
**Status**: Production Ready
