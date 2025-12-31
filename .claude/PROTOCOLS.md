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

**Priority: PARALLEL BY DEFAULT.** Each agent runs in its own 200k context window.

### When to Parallelize

```yaml
parallel_when:
  - independent_research_tasks
  - multi_domain_work (frontend + backend + infra)
  - unrelated_file_changes
  - separate_component_reviews
  - multi_agent_analysis (security + performance + architecture)

sequential_when:
  - output_dependencies (agent B needs agent A's output)
  - same_file_edits
  - ordered_operations (build before deploy)
```

### How to Invoke Parallel Agents

**CRITICAL**: Use a SINGLE message with multiple Task tool calls.

```
User: "Add authentication and update the API docs"

Main Agent Response (SINGLE message):
  → Task tool: developer agent (implement auth)
  → Task tool: researcher agent (update docs)

Both agents run SIMULTANEOUSLY in separate contexts.
```

### Parallel Patterns

| Pattern | Agents | Use Case |
|---------|--------|----------|
| **Star** | 3-5 independent | Reviews, research, parallel features |
| **Pipeline** | 2-3 sequential | Build → test → deploy |
| **Hybrid** | Mix | Strategic planning + tactical execution |

### Swarm Execution (Advanced)

For complex multi-faceted work, invoke 3-5 agents simultaneously:

```
Complex Feature Request:
  → architect (design)
  → developer (implementation plan)
  → platform (infra requirements)
  → researcher (prior art)

All in ONE message. Main agent synthesizes results.
```

### Anti-Patterns

- Invoking agents one-at-a-time when independent
- Waiting for agent A before invoking unrelated agent B
- Main agent doing research that could parallelize

## Skills

| Skill | Trigger | Purpose |
|-------|---------|---------|
| beads | Session start/end | Execution state tracking |
| context-display | Every response | Show usage |
| setup-validation | Project init | Verify environment |
| task-generate | PRD approved | Create tasks from PRD |
| prd-create | Feature planning | Create PRD documents |

## Commands (User-Invoked)

| Command | Purpose |
|---------|---------|
| /commit | Git commit workflow |
| /mr | Merge request creation |

## Work Management

### Folder Structure
```
work/
├── 0_vision/   # Strategic vision
├── 1_backlog/  # Ready for development
├── 2_active/   # In progress
└── 3_done/     # Completed

.beads/
└── issues.jsonl  # Execution state (replaces history/)
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

## Session Continuity (Beads)

Session state is tracked via beads, not handoff files.

### Session Start
```bash
python .claude/skills/beads/scripts/beads.py context
```
Shows in-progress, ready, and blocked beads.

### Session End
```bash
# Close completed work
python .claude/skills/beads/scripts/beads.py close <id> --note "Done"

# Commit beads state
git add .beads/ && git commit -m "Update beads state"
```

### At 75% Context
1. Close/update relevant beads
2. Commit `.beads/` changes
3. Context bar shows usage - start new session if needed

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

marketplace/
└── agents/           # 130+ specialized agents
```

---

Garden 2.0 - See CLAUDE.md for quick reference.
