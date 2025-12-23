# Project Context

> Project-specific information for Claude agents.

## Overview

**Project**: The Garden
**Description**: Agent orchestration system for Claude Code
**Version**: 2.0.0

## Purpose

Garden provides a structured approach to using Claude Code with:
- Specialized agents for different domains
- Skills for automatic behaviors
- Commands for user-triggered actions
- Plugin system for extensibility

## Tech Stack

- **Platform**: Claude Code CLI
- **Language**: Markdown (agent definitions), Python (hooks), JSON (config)
- **Version Control**: Git
- **CI/CD**: GitLab

## Key Files

| File | Purpose |
|------|---------|
| `.garden/PROTOCOLS.md` | Core operating rules |
| `.garden/PROJECT.md` | This file - project context |
| `.garden/config/plugins.json` | Plugin configuration |
| `CLAUDE.md` | Legacy v1 configuration (migration source) |

## Agents Available

| Agent | Domain |
|-------|--------|
| developer | Code implementation |
| architect | System design |
| product | Requirements |
| platform | Infrastructure |
| researcher | Analysis |

## Current Phase

**Phase 0**: Foundation
- Core structure created
- 5 agents defined
- 4 skills defined
- Plugin manifest initialized

## Migration Notes

This is Garden 2.0, redesigned from v1 (`.claude/`):
- Renamed folder: `.claude/` → `.garden/`
- Reduced protocol file: `CLAUDE.md` (945 lines) → `PROTOCOLS.md` (~200 lines)
- Consolidated agents: 12+ → 5 core
- Added skills system
- Added plugin architecture

## Conventions

- Agent files: `.garden/agents/[name].md`
- Skill files: `.garden/skills/[name].md`
- Command files: `.garden/commands/[name].md`
- Work items: `.garden/work/[stage]/[name].md`
- Session history: `.garden/work/history/[timestamp]-[agent]-[desc].md`

## Constraints

- Main agent orchestrates only - never implements
- TDD required for all code
- Context display required on every response
- Session handoff at 75% context usage
