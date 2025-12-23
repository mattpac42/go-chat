# Garden 2.0 Documentation

> Complete guide to the Garden agent orchestration system.

## Quick Navigation

| Section | Description |
|---------|-------------|
| [Getting Started](getting-started/) | First-time setup and basics |
| [Guides](guides/) | How-to guides for common tasks |
| [Reference](reference/) | Detailed specifications |

## Getting Started

1. **[Quick Setup](getting-started/quick-setup.md)** - Get running in 2 minutes
2. **[Core Concepts](getting-started/core-concepts.md)** - Understand agents, skills, commands
3. **[First Task](getting-started/first-task.md)** - Complete your first task with Garden

## Guides

- **[Working with Agents](guides/working-with-agents.md)** - How delegation works
- **[Context Management](guides/context-management.md)** - Session handoffs and continuity
- **[Plugin System](guides/plugin-system.md)** - Installing and using plugins
- **[Customization](guides/customization.md)** - Adapting Garden to your workflow

## Reference

- **[PROTOCOLS.md](../PROTOCOLS.md)** - Core operating rules
- **[Agents Reference](reference/agents.md)** - All 5 core agents
- **[Skills Reference](reference/skills.md)** - Automatic behaviors
- **[Commands Reference](reference/commands.md)** - User-invoked actions
- **[Configuration](reference/configuration.md)** - Settings and plugins

## Key Files

| File | Purpose |
|------|---------|
| `.garden/PROTOCOLS.md` | Core rules (read first) |
| `.garden/PROJECT.md` | Project-specific context |
| `.garden/QUICKSTART.md` | 2-minute setup |
| `.garden/config/plugins.json` | Plugin configuration |

## Architecture

```
.garden/
├── PROTOCOLS.md      # Core rules
├── PROJECT.md        # Project context
├── QUICKSTART.md     # Setup guide
├── agents/           # 5 core agents
├── skills/           # Auto-invoked
├── commands/         # User-invoked
├── work/
│   ├── 0_vision/     # Strategic vision
│   ├── 1_backlog/    # Ready for development
│   ├── 2_active/     # In progress
│   ├── 3_done/       # Completed
│   └── history/      # Session logs
├── scripts/          # CLI tools
├── config/           # Configuration
├── plugins/          # Extensions
├── templates/        # Boilerplate
└── docs/             # This documentation
```

## Version

Garden 2.0 - Simplified, skills-first, plugin-ready.
