# Plant Command

**Task**: Root a new project from The Garden or add Garden capabilities to an existing repository. Creates project structure with selected agents, skills, and workflows, maintaining a sync connection back to The Garden.

## Usage

```bash
/plant [project-name] [path] [type] [agents] [--stack tech-stack] [--existing]
```

**Arguments** (all optional - if omitted, enters interactive mode):
- `[project-name]`: Kebab-case project identifier (e.g., `my-web-app`)
- `[path]`: Directory where project will be created (or existing repo path)
- `[type]`: Project type (webapp, api, cli, mobile, library, data, devops, business)
- `[agents]`: Comma-separated list of agents to include
- `[--stack]`: Tech stack for devcontainer (python-flask, python-fastapi, node-react, node-nextjs)
- `[--existing]`: Add Garden to an existing repository (won't create new directory)

**Examples**:
- `/plant` - Interactive wizard mode
- `/plant my-app ~/projects webapp` - Quick mode with defaults
- `/plant my-api ~/projects api developer,architect,platform-tactical` - Full specification
- `/plant my-app ~/projects webapp --stack python-flask` - With devcontainer
- `/plant my-repo ~/existing/repo webapp --existing` - Add Garden to existing repo

## Steps

### 1. Parse Command Arguments

- Check if arguments were provided
- **No arguments**: Enter full interactive wizard (Phase 1-5)
- **Arguments provided**: Use provided values, prompt only for missing required fields

### 2. Phase 1: Discovery (Interactive Mode)

**Step 2.1: New or Existing Project**

First, determine if this is a new project or existing repository:

```
What would you like to do?

A) Create a new project - Start fresh with a new directory
B) Add Garden to existing repo - Enhance an existing codebase

Select [A-B]:
```

**If "Existing Repo" selected**: Skip to Step 2.1a (Existing Repo Flow)
**If "New Project" selected**: Continue to Step 2.1b (Project Type)

**Step 2.1a: Existing Repository Flow**

```
Enter the path to your existing repository:

Examples:
- ~/projects/my-existing-app
- /Users/username/work/legacy-api
- ../sibling-project

Path:
```

Validate:
- Path exists and is a directory
- Has write permissions
- Optionally: Check if it's a git repository

Then ask for project type (for agent selection):

```
What type of project is this? (Used for agent recommendations)

A) Web App - Frontend/fullstack web application
B) API - Backend service or REST/GraphQL API
C) CLI - Command-line tool or utility
D) Mobile - iOS/Android mobile application
E) Library - Reusable package or library
F) Data - Data pipeline, ML, or analytics project
G) DevOps - Infrastructure or platform project
H) Business - Business planning or strategy project

Select [A-H]:
```

**Existing Repo Behavior**:
- Will NOT overwrite existing CLAUDE.md (backs up to CLAUDE.md.backup if exists)
- Will NOT overwrite existing .claude/ contents (merges new files only)
- Creates .claude/lineage.json for Garden sync
- Adds agents to .claude/agents/ (skips existing)
- Adds skills, commands, templates (skips existing)

**Step 2.1b: Project Type (New Projects)**

Present project type options:

```
What are you building?

A) Web App - Frontend/fullstack web application
B) API - Backend service or REST/GraphQL API
C) CLI - Command-line tool or utility
D) Mobile - iOS/Android mobile application
E) Library - Reusable package or library
F) Data - Data pipeline, ML, or analytics project
G) DevOps - Infrastructure or platform project
H) Business - Business planning or strategy project

Select [A-H]:
```

Default agents per type:

| Type | Agents |
|------|--------|
| webapp | developer, architect, ux-tactical |
| api | developer, architect, platform-tactical |
| cli | developer, architect |
| mobile | developer, architect, ux-tactical |
| library | developer, architect |
| data | developer, researcher |
| devops | platform-tactical, platform-strategic, cicd-tactical |
| business | product, product-visionary, business-strategic |

**Step 2.1b: Tech Stack (Optional)**

For applicable project types, offer devcontainer generation:

```
Would you like a devcontainer for local development?

A) Python - Flask
B) Python - FastAPI
C) Node.js - React
D) Node.js - Next.js
N) No devcontainer

Select [A-D/N]:
```

Stack mapping by project type:

| Project Type | Available Stacks |
|--------------|------------------|
| webapp | node-react, node-nextjs, python-flask |
| api | python-flask, python-fastapi, node-nextjs |
| cli | (no devcontainer) |
| mobile | (no devcontainer) |
| library | python-flask, node-react |
| data | python-flask, python-fastapi |
| devops | (no devcontainer) |
| business | (no devcontainer) |

If user selects N or project type doesn't support devcontainers, skip devcontainer generation.

**Step 2.2: Project Name**

```
What should the project be called?

Requirements:
- Lowercase letters, numbers, and hyphens only
- Will become the directory name
- Example: my-awesome-app

Project name:
```

Validate:
- Kebab-case format
- No special characters except hyphens
- Not empty

**Step 2.3: Project Location**

```
Where should the project be created?

Enter the full path to the parent directory.
The project folder will be created inside this path.

Examples:
- ~/projects
- /Users/username/work
- ../

Path:
```

Validate:
- Path exists or can be created
- Have write permissions
- Show resolved full path for confirmation

**Step 2.4: Project Description**

```
Briefly describe the project (1-2 sentences):
```

Default to: "A [type] project rooted from The Garden"

### 3. Phase 2: Domain Selection

Show default agents for the selected type and allow adjustments:

```
Based on your project type, here are the recommended agents:

Default Agents:
  [x] developer - Code implementation, testing, debugging
  [x] architect - System design, patterns, technical strategy
  [x] ux-tactical - UI implementation and components

Additional agents available:
  [ ] platform-tactical - Infrastructure hands-on work
  [ ] platform-strategic - Cloud strategy and planning
  [ ] product - Requirements, PRDs, feature planning
  [ ] researcher - Analysis, exploration
  [ ] cicd-tactical - Pipeline implementation
  [ ] (more available in gnomes/agents/)

Options:
A) Accept recommended agents
B) Add more agents
C) Remove agents
D) Browse gnomes/agents/ library

Select [A-D]:
```

If user selects B or D, show available agents from gnomes/agents/ directory.

### 4. Phase 3: Manifest Confirmation

Display exactly what will be created:

```
📋 Project Manifest

Project: [project-name]
Location: [full-path]/[project-name]
Type: [type]
Tech Stack: [stack or "None"]
Description: [description]

Agents to copy:
  ✓ developer
  ✓ architect
  ✓ ux-tactical

Structure to create:
  ✓ CLAUDE.md (customized)
  ✓ .claude/lineage.json (Garden connection)
  ✓ .claude/agents/ (3 agents)
  ✓ .claude/skills/ (7 core skills)
  ✓ .claude/commands/ (6 commands)
  ✓ .claude/templates/ (all templates)
  ✓ .claude/work/ (work folders)

Devcontainer: [stack name or "Not included"]
  ✓ .devcontainer/devcontainer.json
  ✓ .devcontainer/setup.sh

Proceed? [Y/n]:
```

### 5. Phase 4: Generation

Execute the scaffolding:

**For new projects:**
```bash
python3 .claude/skills/plant-project/scripts/init_project.py \
  "[project-name]" \
  --path "[path]" \
  --type "[type]" \
  --agents "[agent1,agent2,...]" \
  --description "[description]" \
  --stack "[tech-stack]"  # Optional, omit if no devcontainer
```

**For existing repositories:**
```bash
python3 .claude/skills/plant-project/scripts/init_project.py \
  "[repo-name]" \
  --path "[existing-repo-path]" \
  --type "[type]" \
  --agents "[agent1,agent2,...]" \
  --description "[description]" \
  --existing  # Enables merge mode
```

Display progress as files are created.

### 6. Phase 5: Next Steps

After successful creation:

```
🌱 Project '[project-name]' planted successfully!

Location: [full-path]
Devcontainer: [stack name or "Not included"]

Next steps:
1. cd [full-path]
2. Open in VS Code and reopen in container (if devcontainer included)
   OR Start Claude Code in the new project
3. Run /onboard to configure workspace
4. (Optional) Run /product-vision to create strategic vision

Your project is connected to The Garden.
Run /sync-baseline anytime to pull updates.
```

## Quick Mode

When all arguments are provided, skip interactive prompts:

```bash
/plant my-app ~/projects webapp developer,architect
/plant my-app ~/projects webapp developer,architect --stack python-flask
```

1. Validate all arguments
2. Show brief manifest
3. Generate project (with devcontainer if --stack provided)
4. Display success message

## Error Handling

| Scenario | Response |
|----------|----------|
| Directory already exists (new mode) | "Directory [path]/[name] already exists. Choose a different name or path." |
| Directory doesn't exist (existing mode) | "Directory [path] doesn't exist. Did you mean to create a new project?" |
| Invalid project name | "Project name must be kebab-case (lowercase, hyphens). Example: my-project" |
| Path doesn't exist | "Path [path] doesn't exist. Create it? [Y/n]" |
| No write permission | "Cannot write to [path]. Check permissions or choose another location." |
| Agent not found | "Warning: Agent '[name]' not found, skipping. Continuing with available agents." |
| Existing CLAUDE.md found | "Found existing CLAUDE.md. Backed up to CLAUDE.md.backup" |
| Existing .claude/ found | "Found existing .claude/ folder. Merging new files (won't overwrite existing)." |
| Script error | Display error message, suggest checking paths and permissions |

## Project Types Reference

| Type | Description | Use Case |
|------|-------------|----------|
| `webapp` | Web application | React, Vue, Next.js, frontend projects |
| `api` | Backend service | REST API, GraphQL, microservices |
| `cli` | Command-line tool | CLI utilities, scripts, automation |
| `mobile` | Mobile app | iOS, Android, React Native, Flutter |
| `library` | Package/library | npm packages, Python modules, shared code |
| `data` | Data/ML | Data pipelines, analytics, ML models |
| `devops` | Infrastructure | Terraform, Kubernetes, CI/CD |
| `business` | Business planning | Startup planning, strategy, product vision |

## Success Criteria

- [ ] User guided through project setup (or quick mode completes)
- [ ] Project directory created at specified path
- [ ] CLAUDE.md generated with project info and agent table
- [ ] lineage.json created with Garden connection (includes tech_stack)
- [ ] Selected agents copied successfully
- [ ] Core skills and commands copied
- [ ] Devcontainer generated if tech stack selected
- [ ] User informed of next steps
- [ ] No errors during generation
