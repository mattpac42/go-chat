---
description: Prepare handoff summary for next agent (project, gitignored)
---

# Handoff Command

Invoke the handoff skill to create session transition files.

## What This Does

Creates two files in `.claude/work/history/`:
- `SESSION-[XXX].md` - Summary of current session
- `HANDOFF-SESSION.md` - Context for next session

## Workflow

This command invokes the **handoff skill** which:

1. Determines next session number (001, 002, etc.)
2. Gathers session context (tasks, decisions, modified files)
3. Creates session summary file using template
4. Creates/updates handoff file for next session
5. Confirms files were created

## Usage

Run `/handoff` when:
- Context reaches 75% (also triggers automatically)
- Ending a work session
- Handing off to another person/session

## Output

```
🔄 Session handoff created

Files created:
- .claude/work/history/SESSION-001.md
- .claude/work/history/HANDOFF-SESSION.md

Ready for new session. Run /catch-up to restore context.
```

## Related

- **catch-up skill**: Restores context from HANDOFF-SESSION.md
- **context-display skill**: Monitors threshold for auto-trigger
