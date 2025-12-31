# The Garden - Documentation Index

Welcome to The Garden documentation. This directory contains high-level documentation for understanding and using The Garden baseline template.

## Navigation

| Document | Purpose |
|----------|---------|
| [QUICKSTART.md](../.claude/QUICKSTART.md) | 2-minute setup guide |
| [CLAUDE.md](../CLAUDE.md) | Quick reference (auto-loaded) |
| [PROTOCOLS.md](../.claude/PROTOCOLS.md) | Implementation details |
| This folder | Deep dive documentation |

## Getting Started

### New to The Garden?
1. **[QUICKSTART.md](../.claude/QUICKSTART.md)** - Fast 2-minute setup
2. **[GARDEN_CONTEXT.md](GARDEN_CONTEXT.md)** - Comprehensive system overview

### Understanding The System
- **[Architecture Overview](GARDEN_CONTEXT.md#system-architecture)** - Agent orchestration model
- **[PRD Lifecycle](GARDEN_CONTEXT.md#prd-lifecycle-workflow)** - From design to delivery
- **[Context Management](GARDEN_CONTEXT.md#context-management)** - Automated tracking and handoff

## Core Documentation

### Agent Guides
- **[garden-guide-agent.md](garden-guide-agent.md)** - System setup and optimization coach
- **[project-navigator-agent.md](project-navigator-agent.md)** - Project knowledge and memory

### Technical Docs
Located in `.claude/docs/`:
- Context management strategies
- TDD workflow requirements
- CI/CD best practices
- Agent invocation patterns

## Key Concepts

### Agent Orchestration
- **Main Agent** = Orchestrator only (delegates all implementation)
- **Generic Agents** = Quick tasks (`.claude/agents/`)
- **Specialized Agents** = Production work (`marketplace/agents/`)

### Work Management
```
.claude/work/
├── 0_vision/    # Strategic vision documents
├── 1_backlog/   # PRDs ready for development
├── 2_active/    # Work in progress
├── 3_done/      # Completed features
└── history/     # Session logs and handoffs
```

### Context Thresholds
| Level | Action |
|-------|--------|
| 50%+ | Context bar displayed |
| 60% | Warning - approaching limit |
| 75% | Handoff triggered automatically |
| 85% | New session recommended |

## Repository Structure

```
the_garden/
├── docs/                    # High-level docs (you are here)
├── marketplace/agents/           # 130+ specialized agents
├── .claude/
│   ├── PROTOCOLS.md         # Implementation details
│   ├── QUICKSTART.md        # Setup guide
│   ├── PROJECT.md           # Project context template
│   ├── agents/              # 9 generic agents
│   ├── skills/              # 13 workflow skills
│   ├── commands/            # User-invoked commands
│   ├── templates/           # Workflow templates
│   ├── work/                # PRD lifecycle folders
│   ├── config/              # Plugin configuration
│   └── docs/                # Technical documentation
├── CLAUDE.md                # Quick reference (auto-loaded)
└── README.md                # Main entry point
```

## Common Tasks

### Setting Up a New Project

**Option 1: Plant from The Garden (Recommended)**
1. From The Garden, run: `/plant`
2. Follow the wizard to select project type, name, location, and agents
3. In new project, run: `/onboard` to configure workspace
4. Start working - delegation happens automatically

**Option 2: Manual Copy**
1. Copy: `cp -r .claude/ /path/to/project/` and `cp CLAUDE.md /path/to/project/`
2. Run onboarding: `/onboard`
3. Create PROJECT.md with project context
4. Start working - delegation happens automatically

**Sync Updates**: Run `/sync-baseline` in planted projects to pull Garden updates

### Creating a PRD
1. Tell Claude: "I want to create a PRD for [feature]"
2. Answer clarifying questions
3. Review generated PRD in `1_backlog/`
4. Generate task breakdown
5. Move to `2_active/` and implement

## Philosophy

The Garden is built on three core principles:

1. **Main Agent as Orchestrator** - Delegates 100% of implementation work
2. **Documentation-First Development** - PRD-driven workflow
3. **Test-Driven Development** - Mandatory TDD with quality gates

---

**Status**: Production Ready
**Version**: Garden 2.0
