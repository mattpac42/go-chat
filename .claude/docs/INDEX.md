# Garden 2.0 Documentation

> Complete guide to the Garden agent orchestration system.

## Quick Navigation

| Section | Description |
|---------|-------------|
| [Getting Started](getting-started/) | First-time setup and basics |
| [PROTOCOLS.md](../PROTOCOLS.md) | Core operating rules |
| [Skills Manifest](../config/skills-manifest.md) | Workflow skills reference |

## Getting Started

1. **[Quick Setup](getting-started/quick-setup.md)** - Get running in 2 minutes
2. **[Core Concepts](getting-started/core-concepts.md)** - Understand agents, skills, commands
3. **[First Task](getting-started/first-task.md)** - Complete your first task with Garden

## Guides

- **[TDD Workflow](tdd-workflow.md)** - Test-driven development guide
- **[Context Display](context-display-guide.md)** - Token usage visualization
- **[Vision Workflow](vision-workflow-guide.md)** - Product discovery process
- **[GitLab CI/CD](gitlab-cicd-guide.md)** - Pipeline configuration
- **[Strategic Agents](strategic-agents-quick-start.md)** - Using tactical/strategic agent pairs
- **[Devcontainer Audio](devcontainer-audio-setup.md)** - Audio notifications in devcontainers

## Key Files

| File | Purpose |
|------|---------|
| `.claude/PROTOCOLS.md` | Core rules (read first) |
| `.claude/PROJECT.md` | Project-specific context |
| `.claude/QUICKSTART.md` | 2-minute setup |
| `.claude/config/plugins.json` | Plugin configuration |
| `.claude/config/skills-manifest.md` | Skills workflow reference |

## Architecture

```
.claude/
├── PROTOCOLS.md      # Core rules
├── PROJECT.md        # Project context
├── QUICKSTART.md     # Setup guide
├── agents/           # 5 core agents
├── skills/           # Auto and user-invoked
├── work/
│   ├── 0_vision/     # Strategic vision
│   ├── 1_backlog/    # Ready for development
│   ├── 2_active/     # In progress
│   ├── 3_done/       # Completed
│   └── history/      # Session logs (ad-hoc work)
├── scripts/          # CLI tools
├── config/           # Configuration
├── templates/        # Boilerplate
└── docs/             # This documentation

.beads/
└── issues.jsonl      # Execution state (PRD-based work)
```

## Session Management

| Scenario | Use |
|----------|-----|
| PRD-based work with tasks | beads |
| Ad-hoc work, quick fixes | handoff/catch-up |
| Exploratory sessions | handoff/catch-up |
| Long-running features | beads |

## Version

Garden 2.0 - Simplified, skills-first, plugin-ready.
