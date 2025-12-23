# Core Concepts

> Understanding how Garden works.

## The Orchestration Model

Garden uses an orchestration model where:

1. **Main Agent** = Orchestrator (understands, delegates, integrates)
2. **Specialized Agents** = Workers (implement specific domains)

```
You → Main Agent → Specialized Agent → Work Done → Results
         ↓                 ↓
    (orchestrates)    (implements)
```

## Agents

Agents are specialized workers for different domains.

| Agent | Domain | When Used |
|-------|--------|-----------|
| **developer** | Code | Writing, testing, debugging code |
| **architect** | Design | Architecture, patterns, decisions |
| **product** | Requirements | PRDs, features, user stories |
| **platform** | Infrastructure | DevOps, CI/CD, deployment |
| **researcher** | Analysis | Exploration, investigation |

### Invocation Pattern

Agents receive structured briefings:

```
Task: [What to do]
Context: [Background info]
Constraints: [Limitations]
Deliverables: [Expected outputs]
Success: [Completion criteria]
```

## Skills

Skills are **automatic behaviors** that run without user action.

| Skill | Trigger | Purpose |
|-------|---------|---------|
| **catch-up** | Session start | Restore previous context |
| **context-display** | Every response | Show token usage |
| **agent-session-summary** | Agent completion | Document work |
| **setup-validation** | Project init | Verify environment |

Skills are model-invoked (Claude decides when to run them).

## Commands

Commands are **user-triggered actions** via `/command`.

| Command | Purpose |
|---------|---------|
| `/handoff` | Create session transition files |
| `/commit` | Git commit with message |
| `/mr` | Create merge/pull request |
| `/onboard` | Full setup walkthrough |

Commands require explicit user action.

## Plugins

Plugins extend Garden with additional functionality.

- Installed via `garden install <plugin>`
- Configured in `.garden/config/plugins.json`
- Can add agents, skills, or commands

## Context Management

Garden tracks context window usage:

```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 40%
```

### Thresholds

| Level | Percentage | Meaning |
|-------|------------|---------|
| Normal | 0-60% | Continue working |
| Warning | 60-75% | Approaching limit |
| Critical | 75-85% | Handoff recommended |
| Emergency | 85%+ | New session needed |

### Session Continuity

At 75%, Garden creates handoff files:
- Session summary (what was done)
- Handoff file (how to continue)

Next session uses catch-up skill to restore context.

## Work Management

Tasks flow through stages:

```
backlog/ → active/ → done/
```

Each task gets a file tracking:
- Objective and context
- Subtasks with checkboxes
- Success criteria
- Completion notes

## Key Principles

1. **Main agent never implements** - Always delegates
2. **Skills run automatically** - No user action needed
3. **Commands need user action** - Explicit triggers
4. **Context is tracked** - Always visible
5. **Sessions hand off** - Work continues across limits
