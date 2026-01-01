# Go Chat

A full-stack chat wrapper for creating projects in local GitLab with devcontainer support

**Main agent orchestrates. Agents implement. No exceptions.**

## Core Agents

| Agent | Domain |
|-------|--------|
| developer | Code implementation, testing, debugging |
| architect | System design, patterns, technical strategy |
| product | Requirements, PRDs, feature planning |
| platform | Infrastructure, DevOps, CI/CD |
| researcher | Analysis, exploration, information gathering |

Additional: garden-guide, project-navigator, prompt-optimizer, product-visionary

## Agent Libraries

| Library | Location | Use When |
|---------|----------|----------|
| **Generic** | `.claude/agents/` | Quick tasks, prototyping, simple projects |
| **Specialized** | `marketplace/agents/` | Production work, complex features, clear phases |

**Specialized agents use tactical/strategic pairs:**
- `software-tactical` + `software-strategic` (code + architecture)
- `platform-tactical` + `platform-strategic` (infra hands-on + cloud strategy)
- `product-manager-tactical` + `product-manager-strategic` (sprints + roadmap)
- `ux-tactical` + `ux-strategic` (UI implementation + design systems)

Use generic for speed, specialized for precision.

## Context Thresholds

| Level | Action |
|-------|--------|
| 60% | Warning - approaching limit |
| 75% | Update beads, commit, consider new session |
| 85% | New session recommended |

Display context bar at 50%+ or after agent completion:
```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 50% (100k/200k)
```

## Quick Rules

1. **Delegate all implementation** to specialized agents
2. **Main agent only**: reads context, asks questions, invokes agents, tracks progress
3. **Use skills** for workflows (PRD creation, task management, etc.)
4. **Parallel execution**: Run independent agent tasks simultaneously
5. **Beads for state**: Track execution in `.beads/` - no handoff files needed
6. **Track with beads**: Create beads for non-trivial tasks, update as you work

## Beads Workflow

Beads tracks work persistently across sessions. **Always use beads for non-trivial work.**

| When | Command | Example |
|------|---------|---------|
| Session start | `beads context` | See what to work on |
| Starting task | `beads add "title"` | `beads add "Add dark mode" --type feature` |
| Working on it | `beads progress <id>` | `beads progress bd-a1b2` |
| Task done | `beads close <id>` | `beads close bd-a1b2 --note "Added ThemeProvider"` |
| Session end | Commit `.beads/` | `git add .beads/ && git commit` |

**Script location:** `python3 .claude/skills/beads/scripts/beads.py`

**What to track:** Features, bug fixes, refactors, investigations taking >10 min or touching >1 file.
**What to skip:** Quick answers, single-line fixes, research questions.

## Key Files

- `.claude/PROTOCOLS.md` - Detailed rules and workflows
- `.claude/QUICKSTART.md` - 2-minute setup guide
- `.claude/skills/beads/` - Execution state tracking (replaces handoff)
- `.claude/agents/` - 9 active agents
- `marketplace/agents/` - 130+ specialized agents library

## Delegation Decision

Ask four questions:
1. Is this specialized work? → Delegate
2. Will this use >10k tokens? → Delegate
3. Is this my third attempt? → Delegate
4. Are there 2+ independent tasks? → Delegate in **PARALLEL**

If NO to all four → Handle directly

## Parallel Execution

**Default to parallel when possible.** Each agent has its own 200k context window.

**Parallelize when:**
- Multiple independent research tasks
- Different domains (e.g., frontend + backend + infra)
- Unrelated file changes
- Reviews or analysis of separate components

**Invoke parallel agents in a SINGLE message with multiple Task tool calls.**

Example: User asks for "add auth and update docs"
```
→ Invoke developer agent (auth implementation)
→ Invoke researcher agent (docs update)
Both in ONE message, running simultaneously
```

---

For detailed protocols: Read `.claude/PROTOCOLS.md`
