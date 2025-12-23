---
name: workspace-setup
description: Configure VSCode workspace colors for visual project identification. User-invoked skill that creates .vscode/settings.json with customizable color themes. Helps distinguish multiple VSCode windows at a glance.
---

# Workspace Setup

Configure VSCode workspace colors to visually identify this project.

## Purpose

Create `.vscode/settings.json` with workspace color customizations that make this project window visually distinct from other VSCode windows. Useful when working on multiple projects simultaneously.

## Workflow

1. Check if `.vscode/settings.json` already exists
2. If user specified a theme, use that color palette
3. Otherwise, use default Garden theme (Green + Brown)
4. Create `.vscode/` directory if needed
5. Write `settings.json` with color customizations
6. Confirm workspace colors applied

## Available Color Themes

| Theme | Primary | Secondary | Description |
|-------|---------|-----------|-------------|
| Garden (default) | Green (#22c55e) | Brown (#78350f, #d97706) | Emerald green with earthy brown accents |
| Nature | Teal (#14b8a6) | Sage (#365314, #84cc16) | Teal with sage green accents |
| Ocean | Blue (#3b82f6) | Teal (#0f766e, #5eead4) | Bright blue with teal accents |
| Fire | Red (#ef4444) | Orange (#9a3412, #fdba74) | Red with orange accents |
| Royal | Purple (#a855f7) | Blue (#1e40af, #93c5fd) | Purple with blue accents |
| Sunset | Pink (#ec4899) | Orange (#9a3412, #fdba74) | Pink with orange accents |

## Color Application

Colors are applied to these VSCode elements:
- **Title Bar**: Primary accent background with bright foreground
- **Status Bar**: Primary accent background with bright foreground
- **Activity Bar**: Primary accent icons with secondary border
- **Borders**: Secondary accent for tabs, panels, sidebar
- **Sidebar**: Blended background with section headers

## Usage Examples

**Default (Garden theme)**:
```
User: "Run workspace setup"
or
User: "Configure workspace colors"
```

**Specific theme**:
```
User: "Set up workspace with Ocean theme"
or
User: "Configure workspace colors using Royal palette"
```

## Output Format

```
Workspace Colors Configured: [Theme Name]

Created: .vscode/settings.json

Colors Applied:
- Title Bar: [primary color]
- Status Bar: [primary color]
- Activity Bar: [primary color]
- Borders: [secondary color]

Restart VSCode or reload window to see changes.
```

## Technical Details

**File Created**: `.vscode/settings.json`

**Settings Applied**:
- `workbench.colorCustomizations` with theme-specific values
- Title bar, status bar, activity bar colors
- Tab borders, panel borders, sidebar borders
- Sidebar background and section headers

**Reference Template**: `.claude/templates/vscode-settings-template.json`

## Fallback Behavior

- If `.vscode/settings.json` exists with custom colors, ask user to confirm overwrite
- If user doesn't specify theme, default to Garden
- If theme name not recognized, show available themes and ask again
