# Devcontainer Templates

Stack-specific devcontainer configurations for projects planted from The Garden.

## Available Stacks

| Stack | Folder | Base Image | Default Port |
|-------|--------|------------|--------------|
| Python Flask | `python-flask/` | python:3.11-bookworm | 5000 |
| Python FastAPI | `python-fastapi/` | python:3.11-bookworm | 8000 |
| Node React | `node-react/` | node:18-bookworm | 3000 |
| Node Next.js | `node-nextjs/` | node:18-bookworm | 3000 |

## Common Features

All templates include:
- Claude Code CLI (installed via npm)
- Git configuration (main branch, no rebase)
- VS Code extensions for the stack
- Format on save enabled

## Placeholders

Templates use these placeholders that get replaced during project creation:

| Placeholder | Replaced With |
|-------------|---------------|
| `{{PROJECT_NAME}}` | Title case project name (e.g., "My Web App") |
| `{{project_name}}` | Kebab case project name (e.g., "my-web-app") |

## Usage

These templates are automatically copied by `init_project.py` when a tech stack is selected during `/plant`.

```bash
# Quick mode with stack
/plant my-app ~/projects webapp --stack python-flask

# Interactive mode asks about tech stack in Phase 1
/plant
```

## Adding New Stacks

1. Create a new folder: `devcontainer-templates/<stack-name>/`
2. Add `devcontainer.json` and `setup.sh`
3. Update `TECH_STACKS` in `init_project.py`
4. Update `PROJECT_TYPE_STACKS` to map project types to the new stack
