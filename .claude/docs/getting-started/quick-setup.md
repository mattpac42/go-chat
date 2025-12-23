# Quick Setup

> Get Garden running in 2 minutes.

## Prerequisites

- Claude Code CLI installed
- Git repository

## Step 1: Verify Structure

```bash
ls .garden/
```

You should see:
- `PROTOCOLS.md`
- `PROJECT.md`
- `agents/`
- `skills/`
- `commands/`
- `work/`

## Step 2: Validate Setup

Run the Garden CLI:

```bash
.garden/scripts/garden validate
```

All checks should pass.

## Step 3: Create Project Context

If `PROJECT.md` is generic, customize it:

```markdown
# Project Context

## Overview
[Your project description]

## Tech Stack
[Languages, frameworks, tools]

## Key Files
[Important files to know about]
```

## Step 4: Start Working

Just ask Claude for help. The system handles delegation automatically.

Example:
```
"Add a login form to the application"
```

Garden will:
1. Understand the request
2. Delegate to the developer agent
3. Return the implementation

## Verification

Watch for the context bar at the end of responses:

```
Context: 🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 25%
```

This confirms Garden is active.

## Next Steps

- Read [Core Concepts](core-concepts.md) to understand how Garden works
- Try your [First Task](first-task.md) with guided walkthrough
- Explore [Agents Reference](../reference/agents.md) for capabilities

## Troubleshooting

**No .garden folder?**
Run `/onboard` to set up Garden.

**Validation fails?**
Check error messages and create missing files.

**Context bar not showing?**
Ensure PROTOCOLS.md includes context display rules.
