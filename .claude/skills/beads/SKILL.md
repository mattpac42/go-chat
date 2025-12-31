---
name: beads
description: Git-backed issue tracking for AI agent workflows. Replaces handoff/catch-up with persistent execution state. Use for session context, task tracking, and dependency management. Triggers on "beads", "what should I work on", "session context", "track issue".
---

# Beads - Execution State Management

Persistent, git-backed issue tracking that replaces handoff/catch-up workflows.

## Core Concept

Beads provides the **execution layer** for Garden workflows:
- PRDs define *what* to build (strategic)
- Beads track *where we are* (execution)
- Git preserves history (no separate handoff files needed)

## Quick Reference

| Command | Purpose |
|---------|---------|
| `beads init` | Initialize .beads/ in project |
| `beads context` | Session start - what to work on |
| `beads add "title"` | Create new bead |
| `beads progress <id>` | Mark as in-progress |
| `beads close <id>` | Complete a bead |
| `beads list --ready` | Show unblocked work |
| `beads import <file>` | Import from PRD tasks |

## Session Workflow

### Starting a Session

```bash
python .claude/skills/beads/scripts/beads.py context
```

Shows:
- In-progress beads (continue these)
- Ready beads (unblocked, can start)
- Blocked beads (waiting on dependencies)
- Recently closed (context)

### During Work

```bash
# Start working on something
python .claude/skills/beads/scripts/beads.py progress bd-a1b2

# Discover a new issue
python .claude/skills/beads/scripts/beads.py add "Handle edge case" --parent bd-a1b2

# Create a blocker
python .claude/skills/beads/scripts/beads.py add "Need API key" --type bug
python .claude/skills/beads/scripts/beads.py link bd-a1b2 blocked-by bd-c3d4
```

### Ending a Session

```bash
# Close completed work
python .claude/skills/beads/scripts/beads.py close bd-a1b2 --note "Implemented with tests"

# Commit (beads state is in .beads/issues.jsonl)
git add .beads/ && git commit -m "Update beads state"
```

No handoff file needed - beads IS the state.

## Integration with PRDs

### From task-generate

When generating tasks from a PRD, import them as beads:

```bash
python .claude/skills/beads/scripts/beads.py import .claude/work/2_active/001-feature/tasks-feature.md
```

This creates beads with:
- Parent-child relationships
- Agent assignments preserved
- File references attached

### Workflow

```
PRD Created → task-generate → beads import → work on beads → close beads → PRD done
```

## Bead Structure

```json
{
  "id": "bd-a1b2c3",
  "title": "Implement login endpoint",
  "type": "task",
  "status": "open",
  "agent": "developer",
  "priority": "high",
  "tags": ["auth", "api"],
  "files": ["src/auth/login.ts", "tests/auth/login.test.ts"],
  "relationships": {
    "parent": "bd-parent",
    "blocked-by": ["bd-other"],
    "children": []
  },
  "notes": []
}
```

## Statuses

| Status | Meaning |
|--------|---------|
| `open` | Ready to work (if unblocked) |
| `in-progress` | Currently being worked on |
| `blocked` | Waiting on dependencies |
| `closed` | Completed |

## Relationship Types

| Relationship | Meaning |
|--------------|---------|
| `blocks` | This bead blocks another |
| `blocked-by` | This bead is blocked by another |
| `parent` | Parent task |
| `child` | Child/subtask |
| `related` | Related context |
| `discovered-from` | Found while working on another |

## File Locations

```
.beads/
├── issues.jsonl    # All beads (JSONL format)
└── config.json     # Beads configuration
```

## Replaces

This skill replaces:
- `handoff` skill - No longer needed, beads state persists
- `catch-up` skill - Use `beads context` instead
- `agent-session-summary` skill - Bead updates replace summaries
- `.claude/work/history/` files - Git history + beads = full record

## Agent Usage

Agents should:
1. **Session start**: Run `beads context` to see what to work on
2. **Before work**: Run `beads progress <id>` on the bead they're starting
3. **During work**: Create child beads for discovered tasks
4. **After work**: Run `beads close <id>` when done
5. **Session end**: Commit `.beads/` changes

## Example: Full Workflow

```bash
# Initialize (once per project)
python .claude/skills/beads/scripts/beads.py init

# Import tasks from PRD
python .claude/skills/beads/scripts/beads.py import tasks-auth.md

# Start session - what should I work on?
python .claude/skills/beads/scripts/beads.py context

# Start working
python .claude/skills/beads/scripts/beads.py progress bd-a1b2

# Found an issue while working
python .claude/skills/beads/scripts/beads.py add "Edge case: empty password" \
  --type bug --parent bd-a1b2

# Completed the work
python .claude/skills/beads/scripts/beads.py close bd-a1b2 --note "Done with tests"

# Commit
git add .beads/ && git commit -m "feat: implement login - closes bd-a1b2"
```
