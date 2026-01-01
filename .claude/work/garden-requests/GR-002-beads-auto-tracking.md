# Garden Request: GR-002

## Title
Automatic Beads Tracking Integration

## Priority
High

## Category
Enhancement

## Source Project
go-chat

## Date
2025-12-31

---

## Summary

Beads tracking requires explicit manual commands (`beads add`, `beads progress`, `beads close`). This friction means work often goes untracked. The main agent should automatically create and manage beads as part of its natural workflow.

## Current Behavior

1. User asks for work to be done
2. Main agent uses TodoWrite for ephemeral in-session tracking
3. Work gets done and committed
4. Beads are NOT updated unless explicitly commanded
5. Session ends with no persistent record in beads

**Result:** Beads only tracks work when the user explicitly asks about it or the agent remembers to use it.

## Proposed Behavior

### Option A: Integrate Beads with TodoWrite

When `TodoWrite` is used:
- Automatically create corresponding beads for non-trivial tasks
- Sync status changes (pending → in_progress → completed) to beads
- Add notes when tasks complete

```python
# Pseudo-code for TodoWrite hook
def on_todo_change(todo, old_status, new_status):
    if new_status == "in_progress" and not has_bead(todo):
        beads.add(todo.content, type="task")
    elif new_status == "completed":
        beads.close(todo.bead_id, note=todo.completion_note)
```

### Option B: Main Agent Protocol Update

Update CLAUDE.md / PROTOCOLS.md to instruct the main agent:

1. **Session Start**: Run `beads context` to see current state
2. **Before Implementation**: Create beads for planned work
3. **During Work**: Update beads status as tasks progress
4. **Session End**: Ensure all completed work has closed beads, commit `.beads/`

Add to delegation rules:
```
Before delegating to an agent:
1. Create a bead for the task if one doesn't exist
2. Mark the bead as in-progress
After agent completes:
3. Close the bead with a summary note
```

### Option C: Hooks Integration

Create a Claude Code hook that:
- On `Stop` event: Check if work was done, prompt to update beads
- On session start: Auto-run `beads context`

## Recommended Implementation

**Combine Options B and C:**

1. **Update PROTOCOLS.md** with explicit beads workflow instructions
2. **Add to CLAUDE.md** Quick Rules:
   ```
   6. **Track with beads**: Create beads for tasks, update as you work, commit .beads/
   ```
3. **Create a hook** that reminds about beads on session end if `.beads/` has uncommitted changes

## Files to Modify

| File | Change |
|------|--------|
| `.claude/PROTOCOLS.md` | Add "Beads Workflow" section with step-by-step instructions |
| `CLAUDE.md` | Add beads rule to Quick Rules |
| `.claude/skills/beads/SKILL.md` | Add "Automatic Integration" section |
| `.claude/hooks/` | Add post-session beads reminder hook (optional) |

## Benefits

1. **Persistent tracking** - All work is recorded, not just explicitly tagged items
2. **Session continuity** - New sessions can pick up where previous left off via `beads context`
3. **No extra effort** - Beads managed as natural part of workflow
4. **Git history alignment** - Beads provide semantic layer on top of commits

## Example Workflow (After Implementation)

```
User: "Add dark mode support"

Main Agent:
1. beads add "Add dark mode support" --type feature
2. beads progress bd-xxx
3. Delegate to developer agent
4. Developer completes work
5. beads close bd-xxx --note "Added ThemeProvider, toggle in settings"
6. git add .beads/ && git commit --amend (or separate commit)
```

## Risks

- **Over-tracking**: Every minor task getting a bead could be noisy
  - Mitigation: Only create beads for "non-trivial" work (>10 min, >1 file)
- **Sync issues**: TodoWrite and beads getting out of sync
  - Mitigation: Beads is source of truth, TodoWrite is ephemeral

## Notes

This request emerged from a session where significant work was done (migration, bug fixes, new features) but beads only tracked the first explicitly-planned feature. The TodoWrite tool was used extensively for in-session tracking, but none of that persisted to beads.
