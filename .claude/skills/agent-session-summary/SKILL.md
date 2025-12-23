---
name: agent-session-summary
description: Auto-document agent work before returning control to main agent. Use when any specialized agent (developer, architect, product, platform, researcher) completes its assigned task. Creates session history file documenting work performed, decisions made, and deliverables.
---

# Agent Session Summary

Document agent work before returning control to main agent.

## Workflow

1. Gather work summary: task, actions, decisions, files modified, test results
2. Format into structured summary
3. Write to `.claude/work/history/[YYYYMMDD-HHMMSS]-[agent]-[description].md`
4. Return control to main agent with summary

## File Naming

Pattern: `[YYYYMMDD-HHMMSS]-[agent]-[description].md`

- Description: 1-4 words, kebab-case, max 25 chars
- Example: `20251221-143022-developer-add-auth-flow.md`

## Output Format

```markdown
# [Agent] Session: [Description]

**Date**: [timestamp]
**Agent**: [agent name]
**Task**: [assigned task]

## Work Completed
[Summary of actions]

## Decisions Made
- [decision]: [rationale]

## Files Modified
- [path]: [change]

## Recommendations
[Next steps]
```
