# Garden Request: GR-004

## Title
Persist Claude Code sessions across devcontainer rebuilds

## Priority
High

## Category
Enhancement

## Source Project
Go Chat

## Date
2026-01-01

---

## Summary

Claude Code stores session data, credentials, history, and settings in `~/.claude/` inside the container. When a devcontainer is rebuilt, this data is lost, requiring re-authentication and losing conversation context. All devcontainer templates should include a mount to persist this data on the host.

## Current Behavior

Devcontainer templates have no mount for `~/.claude/`. On rebuild:
- User must re-authenticate with Claude
- Session history is lost
- Todo lists are lost
- Project-specific settings are lost
- File history is lost

This creates friction, especially during rapid iteration where rebuilds are common.

## Proposed Behavior

All devcontainer templates should include:

```json
"mounts": [
  "source=${localWorkspaceFolder}/.devcontainer/.claude-data,target=/home/{{USER}}/.claude,type=bind,consistency=cached"
],
```

Where `{{USER}}` matches the `remoteUser` for each template:
- `node` for Node.js templates
- `vscode` for Python templates

Additionally:
- `.devcontainer/.claude-data/` added to template `.gitignore`
- `postCreateCommand` creates the directory if it doesn't exist

## Implementation Notes

**Template-specific home directories:**

| Template | remoteUser | Mount Target |
|----------|------------|--------------|
| node-nextjs | node | /home/node/.claude |
| node-react | node | /home/node/.claude |
| python-fastapi | vscode | /home/vscode/.claude |
| python-flask | vscode | /home/vscode/.claude |

**Setup script addition:**
```bash
# Ensure Claude data directory exists on host
mkdir -p /workspace/.devcontainer/.claude-data
```

## Files to Modify

| File | Change |
|------|--------|
| `.claude/templates/devcontainer-templates/*/devcontainer.json` | Add mounts array |
| `.claude/templates/devcontainer-templates/*/.gitignore` (or template) | Add `.devcontainer/.claude-data/` |
| `.claude/templates/devcontainer-templates/*/setup.sh` | Create directory if missing |

## Testing

1. Create new project with `/plant`
2. Run Claude session, create todos, authenticate
3. Rebuild devcontainer
4. Verify session persists - history, credentials, todos intact
5. Test with both Node.js and Python templates

## Notes

- High priority: rebuilds are common during development, losing session is disruptive
- Security: `.claude-data/` contains credentials, must be gitignored
- The mount directory lives in `.devcontainer/` to keep project root clean
- Backwards compatible: existing projects can add this mount manually

## Design Decision: Per-Project Isolation

Each project gets its own `.devcontainer/.claude-data/` directory. This means:

| Aspect | Behavior |
|--------|----------|
| Credentials | Auth once per project (acceptable) |
| History | Isolated per project (no mixing) |
| Todos | Isolated per project (no mixing) |
| Sessions | Isolated per project (no mixing) |

**Why not share credentials across projects?**
- Adds complexity (split mounts, host-level directories)
- One-time auth per project is minimal friction
- Clean isolation is more predictable
- Avoids cross-project data leakage
