---
name: handoff
description: Create session handoff files for ad-hoc work without PRDs. Use when context reaches 75%, session is ending, or user requests handoff. For PRD-based work, prefer beads instead. Triggers on "handoff", "end session", "save progress".
---

# Session Handoff

Create handoff files for session continuity when not using beads/PRD workflow.

## When to Use

- **Ad-hoc work** without PRDs or formal task tracking
- **Quick sessions** that don't warrant beads setup
- **Context at 75%+** and no beads initialized
- **User explicitly requests** handoff

**For PRD-based work**: Use beads instead (`beads context`, `beads close`)

## Workflow

1. **Gather Context**:
   - What was being worked on
   - Key files modified
   - Decisions made
   - What's pending

2. **Create Session File**:
   - Location: `.claude/work/history/SESSION-[NNN].md`
   - Sequential numbering from existing sessions

3. **Create Handoff File**:
   - Location: `.claude/work/history/HANDOFF-SESSION.md`
   - Overwrites previous (only latest matters)

4. **Commit Changes**:
   - Include session files in commit
   - Note handoff in commit message

## Session File Template

```markdown
# Session [NNN]

**Date**: YYYY-MM-DD
**Duration**: ~X hours
**Context at Handoff**: XX%

## Work Completed
- [What was accomplished]

## Key Decisions
- [Decision]: [Rationale]

## Files Modified
- `path/to/file` - [what changed]

## Pending
- [What's left to do]

## Notes for Next Session
- [Important context]
```

## Handoff File Template

```markdown
# Session Handoff

**From Session**: [NNN]
**Date**: YYYY-MM-DD

## Continue With
[Immediate next step]

## Context
[Essential background]

## Key Files
- `path/to/file` - [why it matters]

## Pending Work
- [ ] [Item 1]
- [ ] [Item 2]
```

## File Locations

```
.claude/work/history/
├── SESSION-001.md      # Historical session
├── SESSION-002.md      # Historical session
└── HANDOFF-SESSION.md  # Current handoff (read by catch-up)
```

## Relationship to Beads

| Scenario | Use |
|----------|-----|
| PRD-based work with tasks | Beads |
| Ad-hoc work, quick fixes | Handoff |
| Exploratory sessions | Handoff |
| Long-running features | Beads |

Both can coexist - beads for structured work, handoff for everything else.
