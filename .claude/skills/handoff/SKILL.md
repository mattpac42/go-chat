---
name: handoff
description: Auto-create session handoff files when context reaches 75%. Use when context usage hits 75% threshold, when user explicitly requests handoff, or when ending a work session. Creates session summary and handoff files in .claude/work/history/ for seamless continuation in next session.
---

# Handoff

Auto-create session transition files at 75% context.

## Trigger

Automatically when:
- Context usage reaches 75%
- User requests session transition
- Work session is ending

## Workflow

1. Detect 75% threshold (from context-display)
2. Gather session context:
   - Current task status
   - Pending work items
   - Key decisions made
   - Modified files
   - Git status
3. Create session summary file
4. Create handoff file
5. Notify user files are ready

## Output Files

### Session Summary
`SESSION-[XXX].md` in `.claude/work/history/`

Sequential numbering (001, 002, 003...).

### Handoff File
`HANDOFF-SESSION.md` in `.claude/work/history/`

Overwrites previous handoff. Contains:
- Immediate context
- Critical files to read
- Last task in progress
- Pending items
- Suggested first action

## Notification

```
🔄 Session handoff created (75% context)

Files created:
- .claude/work/history/SESSION-001.md
- .claude/work/history/HANDOFF-SESSION.md

Ready for new session. Run /catch-up to restore context.
```

## Integration

Works with:
- **context-display**: Monitors threshold
- **catch-up**: Restores from handoff file
