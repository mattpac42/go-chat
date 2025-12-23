# Garden 2.0 Quickstart

> Get started in 2 minutes.

## Navigation

| Document | Purpose |
|----------|---------|
| **This file** | Setup and getting started |
| [CLAUDE.md](../CLAUDE.md) | Quick reference (auto-loaded) |
| [PROTOCOLS.md](PROTOCOLS.md) | Implementation details |
| [docs/](docs/) | Deep dive documentation |

## What is Garden?

Garden is a "rooting" repository for Claude Code agent orchestration. Use it to birth new projects with the right agents, then sync updates as The Garden evolves. The main agent delegates work to specialized agents, tracks context usage, and manages session continuity.

## Core Concepts

| Concept | Description |
|---------|-------------|
| **Agents** | Specialized workers (developer, architect, product, platform, researcher) |
| **Skills** | Auto-invoked capabilities (catch-up, context-display) |
| **Commands** | User-invoked actions (`/handoff`, `/commit`) |
| **Plugins** | Installable extensions |

## Quick Setup

### Creating a New Project

From The Garden, run:
```bash
/plant
```

The wizard will guide you through:
1. Project type (webapp, api, cli, mobile, library, data, devops, business)
2. Project name and location
3. Agent selection based on your needs

**Quick mode**: `/plant my-app ~/projects webapp developer,architect`

### In Your New Project

1. **Run onboarding**
   ```bash
   /onboard
   ```

2. **Create PROJECT.md**
   ```markdown
   # Project Context

   ## Overview
   [Your project description]

   ## Tech Stack
   [Languages, frameworks, tools]

   ## Key Files
   [Important files to know about]
   ```

3. **Start Working**
   Just ask Claude to help. The system handles delegation automatically.

### Sync Updates

Pull updates from The Garden anytime:
```bash
/sync-baseline
```

## How It Works

```
You → Main Agent → Specialized Agent → Work Done → Results
         ↓                 ↓
    (orchestrates)    (implements)
```

1. **You ask** for something
2. **Main agent** understands and delegates
3. **Specialized agent** does the work
4. **Results** come back to you

## Context Management

Watch the context bar at the end of each response:
```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 40% (80k/200k)
```

- **Green (0-60%)**: Normal operation
- **Yellow (60-75%)**: Approaching limit
- **Orange (75-85%)**: Handoff recommended
- **Red (85%+)**: New session needed

## Common Commands

| Command | Purpose |
|---------|---------|
| `/plant` | Create a new project from The Garden |
| `/onboard` | Initialize workspace colors and hooks |
| `/sync-baseline` | Pull updates from The Garden |
| `/handoff` | Create session handoff files |
| `/commit` | Commit changes with message |
| `/mr` | Create merge request |

## Folder Structure

```
.claude/
├── PROTOCOLS.md      # Core rules (read this)
├── PROJECT.md        # Your project context
├── agents/           # Generic agents (quick tasks)
├── skills/           # Auto-invoked capabilities
├── commands/         # User-invoked commands
├── work/
│   ├── 0_vision/     # Strategic vision
│   ├── 1_backlog/    # Ready for development
│   ├── 2_active/     # In progress
│   ├── 3_done/       # Completed
│   └── history/      # Session logs
├── config/
│   └── plugins.json  # Plugin configuration
└── plugins/          # Installed extensions

gnomes/
└── agents/           # 130+ specialized agents (production work)
```

## Agent Libraries

| Library | When to Use |
|---------|-------------|
| `.claude/agents/` | Quick fixes, prototyping, simple tasks |
| `gnomes/agents/` | Production features, complex work, clear phases |

**Specialized agents** use tactical/strategic pairs for precise role separation.

## Next Steps

1. Read `PROTOCOLS.md` to understand the rules
2. Create your `PROJECT.md` with project context
3. Start working - delegation happens automatically

## Getting Help

- Check `.claude/docs/` for detailed documentation
- The system validates setup automatically
- Context is restored from previous sessions

---

Garden 2.0 - Simplified. Skills-first. Plugin-ready.
