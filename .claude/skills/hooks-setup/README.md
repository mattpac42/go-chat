# Hooks Setup Skill

Quick reference for configuring Claude Code hooks in Garden 2.0.

## What It Does

Configures two automated hooks:

1. **UserPromptSubmit**: Shows context usage at start of each exchange
2. **Stop**: Audio + visual notification when Claude finishes

## Usage

Run this skill once to set up hooks, then manually update `.claude/settings.json` with the provided configuration.

## Quick Setup

1. Read `.claude/skills/hooks-setup/SKILL.md`
2. Copy JSON configuration from skill document
3. Update `.claude/settings.json`
4. Ensure scripts are executable (already done by skill)
5. Test by submitting a message and waiting for completion

## Files

- `SKILL.md` - Complete setup instructions
- `../../scripts/completion-status.sh` - Completion notification script

## Platform

- **macOS only** (uses `say` command)
- Requires `jq` for JSON parsing

## Related

- Context Display: `.claude/skills/context-display/`
- Context Tracker: `.claude/scripts/context-tracker.py`
