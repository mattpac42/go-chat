# Garden Protocols

> Implementation details for Claude agents. Reference CLAUDE.md for quick rules.

## Delegation Rules (YAML)

```yaml
delegate_always:
  - code_writing → developer
  - architecture → architect
  - requirements → product
  - infrastructure → platform
  - research → researcher

main_agent_only:
  - reading_context
  - clarifying_questions
  - invoking_agents
  - tracking_progress
  - presenting_results
```

## Agent Selection Logic

```yaml
use_generic:
  - quick_fixes
  - prototyping
  - small_tasks
  - blended_planning_and_doing

use_specialized:
  - production_features
  - clear_phase_separation
  - audit_trail_needed
  - complex_multi_step_work
```

## Agent Invocation Pattern

```
Task: [specific task]
Context: [minimal essential context]
Constraints: [limits and requirements]
Deliverables: [expected outputs]
Success: [completion criteria]
```

## Parallel Execution

```yaml
parallel_when:
  - independent_research
  - multi_domain_work
  - unrelated_files

sequential_when:
  - output_dependencies
  - same_file_edits
  - ordered_operations
```

## Skills (Auto-Invoked)

| Skill | Trigger | Purpose |
|-------|---------|---------|
| catch-up | Session start | Restore context |
| agent-session-summary | Agent exit | Document work |
| context-display | Every response | Show usage |
| setup-validation | Project init | Verify environment |

## Commands (User-Invoked)

| Command | Purpose |
|---------|---------|
| /handoff | Session transition |
| /commit | Git commit workflow |
| /mr | Merge request creation |

## Work Management

### Folder Structure
```
work/
├── 0_vision/   # Strategic vision
├── 1_backlog/  # Ready for development
├── 2_active/   # In progress
├── 3_done/     # Completed
└── history/    # Session logs
```

### Task Lifecycle
1. Create vision in `0_vision/` (if strategic)
2. Create PRD in `1_backlog/`
3. Move to `2_active/` when starting
4. Move to `3_done/` when complete

## Quality Standards

### TDD Required
1. Write failing test
2. Implement minimal code
3. Verify test passes
4. Refactor if needed
5. Run all tests
6. Commit only if green

### Code Principles
- No over-engineering
- No speculative features
- No unnecessary abstractions
- Delete unused code completely

## Error Handling

### Escalation Template
```
STOPPED - GUIDANCE NEEDED

Issue: [problem]
Expected: [plan said]
Actual: [what happened]
Impact: [consequences]

Options:
A) Continue as planned
B) Modify: [changes]
C) Alternative: [approach]

Recommendation: [analysis]
```

## Session Continuity

### Handoff Files
- `SESSION-[XXX].md` - Session summary
- `HANDOFF-SESSION.md` - Next session context

### Required at 75%
1. Create session summary
2. Create handoff file
3. Document pending work
4. Note key decisions

## Plugin System

### Manifest Location
`.claude/config/plugins.json`

## File Locations

```
.claude/
├── PROTOCOLS.md      # This file (implementation details)
├── PROJECT.md        # Project context
├── QUICKSTART.md     # Setup guide
├── agents/           # Generic agents
├── skills/           # Auto-invoked
├── commands/         # User-invoked
├── work/             # Task tracking
├── config/           # Configuration
├── templates/        # Boilerplate
└── docs/             # Documentation

gnomes/
└── agents/           # 130+ specialized agents
```

---

Garden 2.0 - See CLAUDE.md for quick reference.
