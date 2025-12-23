# Go Chat

A full-stack chat wrapper for creating projects in local GitLab with devcontainer support

**Main agent orchestrates. Agents implement. No exceptions.**

## Agents

| Agent | Domain |
|-------|--------|
| developer | Code implementation, testing, debugging |
| architect | System design, patterns, technical strategy |
| ux-tactical | UI implementation and components |
| platform-tactical | Infrastructure hands-on work |
| researcher | Analysis, exploration, information gathering |

## Context Thresholds

| Level | Action |
|-------|--------|
| 60% | Warning - approaching limit |
| 75% | Handoff triggered automatically |
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
5. **Session handoff**: Auto-created at 75% context

## Key Files

- `.claude/PROTOCOLS.md` - Detailed rules and workflows
- `.claude/QUICKSTART.md` - 2-minute setup guide
- `.claude/skills/` - Workflow skills
- `.claude/agents/` - Project agents
- `.claude/lineage.json` - Garden connection for sync

## Delegation Decision

Ask three questions:
1. Is this specialized work? → Delegate
2. Will this use >10k tokens? → Delegate
3. Is this my third attempt? → Delegate

If NO to all three → Handle directly

---

## Garden Lineage

This project was rooted from The Garden on 2025-12-23.

To sync updates from The Garden:
```
/sync-baseline
```

For detailed protocols: Read `.claude/PROTOCOLS.md`
