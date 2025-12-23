---
name: catch-up
description: Restore session context from previous handoff. Use when a new conversation starts and HANDOFF-SESSION.md exists in .claude/work/history/, or when Claude needs to understand what was being worked on in a previous session. Triggers automatically at session start.
---

# Catch-Up

Restore context from previous session to enable seamless work continuation.

## Workflow

1. Check for `.claude/work/history/HANDOFF-SESSION.md`
2. If exists, read and parse the handoff file
3. Extract: last task, pending items, key decisions, critical files
4. Present summary to user
5. Confirm understanding before proceeding

## Output Format

```
Restored from Session [XXX]:

Last Working On:
- [task description]

Pending Items:
- [item 1]
- [item 2]

Ready to continue. What would you like to focus on?
```

## Fallback

If no handoff file exists:
- Check `.claude/work/2_active/` for active tasks
- Load `PROJECT_CONTEXT.md` for project context
- Offer to start fresh
