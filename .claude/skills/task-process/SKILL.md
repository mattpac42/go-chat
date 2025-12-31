---
name: task-process
description: Process task list through implementation workflow. Use when user wants to start implementing tasks, work through a task list, or execute on a feature. Triggers on "start implementing", "work on tasks", "execute task list", "let's build this".
---

# Task Processing

Guide implementation through task list with quality gates.

## Prerequisites

- Approved task list exists
- Move feature folder to active: `1_backlog/` → `2_active/`

## Workflow

1. **Activate Feature**:
   ```bash
   mv .claude/work/1_backlog/[NNN]-[feature]/ .claude/work/2_active/
   ```

2. **Process Subtasks** (one at a time):
   - Present next uncompleted subtask
   - Get user approval to proceed
   - Delegate to assigned agent
   - Verify completion
   - Mark complete immediately

3. **Quality Gates** per subtask:
   - [ ] Code written with tests (TDD)
   - [ ] Tests pass
   - [ ] No regressions
   - [ ] Documentation updated if needed

4. **Parent Task Completion**:
   - All subtasks complete → mark parent complete
   - Verify integration between subtasks

5. **Feature Completion**:
   - All parent tasks complete
   - Final integration test
   - Move to done: `2_active/` → `3_done/`

## Processing Pattern

```
For each uncompleted subtask:
1. "Next: [subtask description]. Proceed?"
2. [User confirms]
3. Delegate to [agent]
4. [Agent completes work]
5. Mark subtask complete
6. Repeat
```

## Completion

When all tasks done:
```bash
mv .claude/work/2_active/[NNN]-[feature]/ .claude/work/3_done/
```

Update status in task file:
```markdown
**Status**: ✅ Complete
**Completed**: [date]
```

## Rules

- **One subtask at a time** - no batching
- **User approval** before each subtask
- **Immediate marking** - complete right after finishing
- **TDD required** - tests before implementation
- **Agent delegation** - never implement directly
