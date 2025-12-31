# Beads Skill

Git-backed issue tracking for AI agent workflows. Inspired by [Steve Yegge's beads](https://github.com/steveyegge/beads).

## Overview

Beads provides persistent execution state that survives sessions:
- **JSONL storage** in `.beads/issues.jsonl`
- **Hash-based IDs** (bd-a1b2) prevent merge conflicts
- **Dependency graph** tracks blockers and relationships
- **Git-native** - versions with your code

## Why Beads?

Replaces the handoff/catch-up workflow with something simpler:

| Before | After |
|--------|-------|
| Handoff files in history/ | Beads state in .beads/ |
| Catch-up reads session files | `beads context` shows state |
| Manual PRD checkbox tracking | Beads track completion |
| Session summaries | Bead notes and git history |

## Installation

Beads is included in the Garden baseline. For existing projects:

```bash
# Copy the skill
cp -r .claude/skills/beads /path/to/project/.claude/skills/

# Initialize beads
python .claude/skills/beads/scripts/beads.py init
```

## Commands

### Initialize
```bash
python .claude/skills/beads/scripts/beads.py init
```

### Add a bead
```bash
python .claude/skills/beads/scripts/beads.py add "Fix login bug" \
  --type bug \
  --agent developer \
  --priority high \
  --tags "auth,urgent"
```

### List beads
```bash
# All open beads
python .claude/skills/beads/scripts/beads.py list

# Ready to work (unblocked)
python .claude/skills/beads/scripts/beads.py list --ready

# By agent
python .claude/skills/beads/scripts/beads.py list --agent developer

# Include closed
python .claude/skills/beads/scripts/beads.py list --all
```

### Session context
```bash
python .claude/skills/beads/scripts/beads.py context
```

### Update bead status
```bash
# Start working
python .claude/skills/beads/scripts/beads.py progress bd-a1b2

# Close completed
python .claude/skills/beads/scripts/beads.py close bd-a1b2 --note "Done"
```

### Link beads
```bash
# Create a blocker
python .claude/skills/beads/scripts/beads.py link bd-a1b2 blocked-by bd-c3d4

# Parent-child
python .claude/skills/beads/scripts/beads.py link bd-child parent bd-parent
```

### Import from PRD tasks
```bash
python .claude/skills/beads/scripts/beads.py import tasks-feature.md
```

## Integration

### With PRD Workflow
1. Create PRD with prd-create skill
2. Generate tasks with task-generate skill
3. Import tasks: `beads import tasks-feature.md`
4. Work through beads, closing as complete
5. PRD moves to done when all beads closed

### With Git
```bash
# After work session
git add .beads/
git commit -m "Update beads: closed bd-a1b2, bd-c3d4"
```

## File Format

`.beads/issues.jsonl`:
```jsonl
{"id":"bd-a1b2c3","title":"Task one","status":"open",...}
{"id":"bd-d4e5f6","title":"Task two","status":"closed",...}
```

Each line is a complete JSON object. Easy to:
- Parse with any language
- Merge across branches (hash IDs prevent conflicts)
- Search with grep
- Track in git history

## Credits

Inspired by [Steve Yegge's beads system](https://github.com/steveyegge/beads) - a memory upgrade for coding agents. This implementation is a lightweight, Python-based alternative optimized for Garden workflows.
