---
name: catch-up
description: Restore session context from previous handoff. Use at session start when HANDOFF-SESSION.md exists. For PRD-based work, use beads context instead. Triggers on "catch up", "continue session", "what was I working on".
---

# Session Catch-Up

Restore context from previous handoff session.

## When to Use

- **Session start** when `HANDOFF-SESSION.md` exists
- **Resuming ad-hoc work** not tracked in beads
- **User asks** "what was I working on?"

**For PRD-based work**: Use `beads context` instead

## Workflow

1. **Check for Handoff File**:
   ```
   .claude/work/history/HANDOFF-SESSION.md
   ```

2. **If Found**:
   - Read handoff file
   - Summarize key context
   - List pending work
   - Identify files to review

3. **If Not Found**:
   - Check for beads: `beads context`
   - If no beads, start fresh session

4. **Present to User**:
   - Brief summary of previous session
   - What needs to continue
   - Suggested first action

## Context Restoration

```markdown
## Previous Session Summary

**From**: Session [NNN] on YYYY-MM-DD

### Continue With
[Immediate next step from handoff]

### Key Context
[Essential background]

### Pending Items
- [ ] [Item 1]
- [ ] [Item 2]

### Files to Review
- `path/to/file` - [why it matters]
```

## Decision Tree

```
Session Start:
  ├── HANDOFF-SESSION.md exists?
  │   └── YES → Read and summarize (catch-up)
  │   └── NO → Check beads
  │       ├── .beads/issues.jsonl exists?
  │       │   └── YES → Run beads context
  │       │   └── NO → Fresh session
```

## After Catch-Up

Once context is restored:
1. User confirms understanding
2. Continue with pending work
3. Use handoff again at session end (if not switching to beads)

## Relationship to Beads

| Check | Result |
|-------|--------|
| Handoff file exists | Use catch-up |
| Beads initialized | Use `beads context` |
| Neither | Fresh session |

Both systems can coexist. Catch-up handles ad-hoc work, beads handles PRD-tracked work.
