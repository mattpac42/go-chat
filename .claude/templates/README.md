# Templates

Templates used by commands and agents.

## Core Templates

| Template | Used By | Purpose |
|----------|---------|---------|
| `handoff.md` | `/handoff` command | Session handoff for next session |
| `session-summary-template.md` | `/handoff` command | Archive of completed session |
| `agent-template.md` | Agent creation | Standard agent structure |
| `prd.md` | Product agent | PRD document structure |
| `task.md` | Task workflows | Simple task structure |
| `tdd-task-template.md` | Developer agent | TDD-specific task structure |
| `product-vision-template.md` | `/onboard` command | Product vision document |
| `vscode-settings-template.json` | `/onboard` command | Workspace color customization |

## Specialized Templates

| Folder | Purpose |
|--------|---------|
| `deployment-template/` | Vercel deployment configuration |
| `secure-pipeline/` | GitLab CI/CD secure pipeline |

## Notes

- Skills (product-vision, feature-roadmap, prd-create, task-generate) have **inline templates** and don't use external files
- Commands reference external templates via the Read tool
- Agent template is the current standard for creating new agents
