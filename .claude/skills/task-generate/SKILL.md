---
name: task-generate
description: Generate task breakdown from PRD with agent assignments. Use when user has an approved PRD and wants to create implementation tasks, break down work, or assign agents. Triggers on "generate tasks", "break down PRD", "create task list", "what needs to be done".
---

# Task Generation

Convert PRD into actionable tasks with agent assignments.

## Prerequisites

- Approved PRD exists
- PRD in `.claude/work/1_backlog/[NNN]-[feature]/`

## Workflow

1. **Load PRD**:
   - Read the approved PRD
   - Identify functional requirements
   - Note dependencies and constraints

2. **Create Parent Tasks** (3-7 total):
   - Group related requirements
   - Define clear objectives
   - Sequence by dependencies

3. **Create Subtasks** (max 8 per parent):
   - Specific, actionable items
   - 1-4 hours of work each
   - Clear completion criteria

4. **Assign Agents**:
   - Match expertise to task domain
   - Balance workload (max 10 per agent)
   - Document assignment rationale

5. **Identify Files**:
   - List files to create/modify
   - Note test files needed

## Agent Selection

| Domain | Agent |
|--------|-------|
| Code, APIs, logic | developer |
| Architecture, design | architect |
| Requirements, stories | product |
| Infrastructure, DevOps | platform |
| Analysis, research | researcher |

## Output: tasks-[feature].md

```markdown
# Tasks: [Feature Name]

**PRD**: prd-[feature].md
**Status**: Ready for implementation

## Task 1: [Parent Task Name]
**Agent**: [agent]
**Objective**: [what to accomplish]

### Subtasks
- [ ] 1.1 [Subtask] - Files: [file paths]
- [ ] 1.2 [Subtask] - Files: [file paths]

## Task 2: [Parent Task Name]
...

## Agent Workload
| Agent | Tasks | Subtasks |
|-------|-------|----------|
```

## File Location

Same folder as PRD: `.claude/work/1_backlog/[NNN]-[feature]/tasks-[feature].md`

## Beads Integration

After generating tasks, import them into beads for persistent tracking:

```bash
python .claude/skills/beads/scripts/beads.py import .claude/work/1_backlog/[NNN]-[feature]/tasks-[feature].md
```

This creates beads with:
- Parent-child relationships from task hierarchy
- Agent assignments preserved
- File references attached
- Dependency tracking enabled

The tasks markdown remains the reference, beads track execution state.

## Next Step

After task approval:
1. Import tasks to beads: `beads import tasks-[feature].md`
2. Move PRD to `2_active/`
3. Start working: `beads context` → `beads progress <id>`
