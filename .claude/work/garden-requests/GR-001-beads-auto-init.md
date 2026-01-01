# Garden Request: GR-001

## Title
Auto-initialize Beads in Plant and Onboard Workflows

## Priority
Medium

## Category
Workflow Enhancement

## Source Project
go-chat

## Date
2025-12-31

---

## Summary

The beads system (execution state management) requires manual initialization (`beads init`) after planting a new project. This should be automated as part of the standard Garden workflows.

## Current Behavior

1. User runs `/plant` to create new project
2. Project is created with `.claude/` structure
3. User must manually run `python .claude/skills/beads/scripts/beads.py init` to use beads
4. If forgotten, `beads context` and other commands fail

## Proposed Behavior

### Option A: Initialize in `/plant` (Recommended)

Update `init_project.py` to create `.beads/` folder with initial files:

```python
# In init_project.py, add after creating .claude/ structure:
def init_beads(project_path):
    beads_dir = project_path / ".beads"
    beads_dir.mkdir(exist_ok=True)

    # Create empty issues file
    (beads_dir / "issues.jsonl").touch()

    # Create config
    config = {
        "version": "1.0",
        "project": project_name,
        "created": datetime.now().isoformat()
    }
    (beads_dir / "config.json").write_text(json.dumps(config, indent=2))
```

### Option B: Initialize in `/onboard` (Fallback)

Update onboard to check for `.beads/` and initialize if missing:

```python
# In onboard workflow:
if not (project_path / ".beads").exists():
    # Run beads init
    subprocess.run(["python", ".claude/skills/beads/scripts/beads.py", "init"])
```

## Recommended Implementation

1. **Primary**: Add beads init to `/plant` - all new projects get beads automatically
2. **Secondary**: Add beads check to `/onboard` - catches existing projects and manual setups

## Files to Modify

| File | Change |
|------|--------|
| `.claude/skills/plant-project/scripts/init_project.py` | Add `init_beads()` function, call during scaffolding |
| `.claude/skills/plant-project/SKILL.md` | Document beads initialization |
| `.claude/commands/onboard.md` | Add beads check/init step |

## Testing

1. Run `/plant test-project ~/tmp webapp developer`
2. Verify `.beads/` exists with `issues.jsonl` and `config.json`
3. Run `beads context` - should work without manual init

## Notes

- Beads replaces the old handoff/catch-up workflow
- Projects without beads can still use old workflow, but beads is preferred
- Backwards compatible - existing projects unaffected until they run `/onboard`
