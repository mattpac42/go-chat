---
name: plant-project
description: Root a new project from The Garden. Use when user wants to create a new project, start a new codebase, or initialize a repository with Garden agents and workflows. Triggers on "/plant", "create new project", "start new repo", "initialize project".
---

# Plant Project

Root new projects from The Garden through an interactive wizard.

## Overview

This skill guides users through creating a new project that inherits The Garden's agents, skills, and workflows. The new project receives only the agents relevant to its needs and maintains a sync connection back to The Garden for updates.

## Workflow

### Phase 1: Discovery

Gather information about the new project through guided questions.

**Step 1.1: Project Type**

Ask the user what they're building:

| Type | Description | Default Agents |
|------|-------------|----------------|
| `webapp` | Web application (frontend/fullstack) | developer, architect, ux-tactical |
| `api` | Backend API service | developer, architect, platform-tactical |
| `cli` | Command-line tool | developer, architect |
| `mobile` | Mobile application | developer, architect, ux-tactical |
| `library` | Reusable package/library | developer, architect |
| `data` | Data/ML project | developer, researcher |
| `devops` | Infrastructure/DevOps | platform-tactical, platform-strategic, cicd-tactical |
| `business` | Business/startup planning | product, product-visionary, business-strategic |

**Step 1.1b: Tech Stack (Optional)**

For applicable project types, offer devcontainer generation:

| Project Type | Available Stacks |
|--------------|------------------|
| webapp | node-react, node-nextjs, python-flask |
| api | python-flask, python-fastapi, node-nextjs |
| library | python-flask, node-react |
| data | python-flask, python-fastapi |

Stacks include Claude Code CLI, language-specific dev tools, and VS Code extensions.

**Step 1.2: Project Name**

Ask for a project name:
- Must be kebab-case (lowercase, hyphens)
- Will become the directory name
- Will be title-cased in CLAUDE.md

**Step 1.3: Project Location**

Ask where to create the project:
- Always ask - no default assumed
- Validate the path exists and is writable
- Show the full resolved path before proceeding

**Step 1.4: Project Description**

Ask for a brief description:
- 1-2 sentences explaining the project's purpose
- Used in CLAUDE.md and lineage.json

### Phase 2: Domain Selection

Based on the project type, suggest relevant domains and let the user adjust.

**Available Domains:**

| Domain | Agents | When to Include |
|--------|--------|-----------------|
| Software Development | developer, architect | Almost always |
| Platform/Infrastructure | platform-tactical, platform-strategic | Cloud, deployment, infra |
| Product Management | product, product-visionary | Feature planning, PRDs |
| UX/Design | ux-tactical, ux-strategic | User interfaces |
| Research/Analysis | researcher | Data analysis, exploration |
| CI/CD | cicd-tactical, cicd-strategic | Pipelines, automation |
| Site Reliability | sre-tactical, sre-strategic | Reliability, monitoring |
| Business | business-strategic | Strategy, planning |

**Selection Logic:**

1. Start with default agents for the project type
2. Present the list to the user
3. Allow additions from marketplace/agents/ library
4. Allow removals if not needed
5. Confirm final selection

### Phase 3: Agent Manifest

Show the user exactly what will be copied:

```
📋 Agent Manifest for [project-name]

Core Agents:
  ✓ developer - Code implementation, testing, debugging
  ✓ architect - System design, patterns, technical strategy

Specialized Agents:
  ✓ ux-tactical - UI implementation and components

Skills to Include:
  ✓ catch-up, handoff, context-display, agent-session-summary
  ✓ setup-validation, hooks-setup, workspace-setup

Templates to Include:
  ✓ agent-template.md, prd.md, product-vision-template.md
  ✓ task.md, tdd-task-template.md, handoff.md

Commands to Include:
  ✓ /commit, /catch-up, /handoff, /mr, /onboard, /sync-baseline

Devcontainer: [stack name or "Not included"]
  ✓ .devcontainer/devcontainer.json
  ✓ .devcontainer/setup.sh

Proceed? [Y/n]
```

### Phase 4: Generation

Execute the scaffolding via init_project.py:

```bash
python3 .claude/skills/plant-project/scripts/init_project.py \
  <project-name> \
  --path <path> \
  --type <project-type> \
  --agents <agent1,agent2,...> \
  --description "<description>" \
  --stack <tech-stack>  # Optional
```

**What Gets Created:**

```
<project-name>/
├── CLAUDE.md              # Customized with project info and agents
├── .devcontainer/         # If tech stack selected
│   ├── devcontainer.json  # VS Code devcontainer config
│   └── setup.sh           # Environment setup script
├── .claude/
│   ├── lineage.json       # Garden connection for sync (includes tech_stack)
│   ├── PROTOCOLS.md       # Core protocols
│   ├── QUICKSTART.md      # Setup guide
│   ├── PROJECT.md         # Project context
│   ├── settings.json      # Claude settings
│   ├── agents/            # Selected agents only
│   ├── skills/            # Core skills
│   ├── commands/          # All commands
│   ├── templates/         # All templates
│   ├── config/            # Configuration
│   ├── docs/              # Documentation
│   └── work/              # Work folders
```

### Phase 5: Initialization

After generation, offer next steps:

1. **Navigate to project**: `cd <path>/<project-name>`
2. **Open in VS Code**: If devcontainer included, reopen in container
3. **Start Claude Code**: Launch in the new project
4. **Run onboarding**: `/onboard` to configure workspace colors and hooks
5. **Create vision** (optional): `/product-vision` for strategic planning
6. **Sync later**: `/sync-baseline` to pull Garden updates

## Quick Mode

For experienced users, support direct invocation:

```
/plant my-project ~/projects webapp developer,architect,ux-tactical
/plant my-project ~/projects webapp developer,architect --stack python-flask
```

Arguments:
1. Project name (required)
2. Path (required)
3. Project type (optional, default: general)
4. Agents (optional, comma-separated)
5. --stack (optional): python-flask, python-fastapi, node-react, node-nextjs

## Success Criteria

- [ ] Project directory created at specified path
- [ ] CLAUDE.md generated with correct project info and agent table
- [ ] .claude/lineage.json created with Garden connection (includes tech_stack)
- [ ] Selected agents copied to .claude/agents/
- [ ] Core skills copied to .claude/skills/
- [ ] All commands copied to .claude/commands/
- [ ] All templates copied to .claude/templates/
- [ ] Devcontainer generated if tech stack selected
- [ ] User informed of next steps

## Error Handling

| Error | Response |
|-------|----------|
| Directory exists | Warn and ask for new name or path |
| Invalid path | Show error, ask for valid path |
| Agent not found | Warn but continue with available agents |
| Permission denied | Show error, suggest alternative path |

## Resources

### Scripts

- `scripts/init_project.py` - Core scaffolding script that creates the project structure

### Related Commands

- `/onboard` - Configure workspace after planting
- `/sync-baseline` - Pull updates from Garden
- `/product-vision` - Create product vision for new project
