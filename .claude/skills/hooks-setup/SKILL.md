---
name: hooks-setup
description: Configure Claude Code hooks for automated context tracking and completion notifications. One-time setup for UserPromptSubmit (context tracking) and Stop (audio + visual completion status) hooks.
---

# Hooks Setup

Configure Claude Code hooks for automated workflow enhancements.

## What This Configures

**Hook 1: UserPromptSubmit** (Context Tracking)
- Runs at the start of each user message
- Shows context usage with emoji visualization
- Triggers automatic handoff at 75% threshold
- Uses `.garden/scripts/context-tracker.py --show-ab`

**Hook 2: Stop** (Completion Notification)
- Runs when Claude finishes a response
- **Audio**: Speaks "The Garden is finished" (macOS `say` command)
- **Visual**: Displays rich status summary:
  ```
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  📊 Session complete | Context: XX% | Agents: X invocations
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ```
- Uses `.garden/scripts/completion-status.sh`

## Prerequisites

- macOS (for `say` command and audio output)
- `jq` installed (for JSON parsing in completion script)
- `.garden/scripts/context-tracker.py` exists (Garden 2.0 standard)
- `.garden/scripts/completion-status.sh` created by this skill

## Configuration Steps

### Option A: New Setup (No Existing Hooks)

Add this to `.claude/settings.json`:

```json
{
  "model": "opus",
  "hooks": {
    "UserPromptSubmit": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "./.garden/scripts/context-tracker.py --show-ab"
          }
        ]
      }
    ],
    "Stop": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "./.garden/scripts/completion-status.sh"
          }
        ]
      }
    ]
  },
  "permissions": {
    "allow": [
      "Bash(say:*)"
    ]
  }
}
```

### Option B: Update Existing Settings

If `.claude/settings.json` already exists with different hooks:

**For UserPromptSubmit:**
- Replace the command path to use `.garden/scripts/context-tracker.py --show-ab`
- Example old path: `./.claude/hooks/post-prompt-context.sh`
- Example new path: `./.garden/scripts/context-tracker.py --show-ab`

**For Stop:**
- Replace the command path to use `.garden/scripts/completion-status.sh`
- Example old path: `./.claude/hooks/completion-one-liners.sh`
- Example new path: `./.garden/scripts/completion-status.sh`

**Permissions:**
- Ensure `"Bash(say:*)"` is in the `permissions.allow` array
- This allows the `say` command for audio notifications

### Example Update

**Before:**
```json
"hooks": {
  "UserPromptSubmit": [
    {
      "matcher": "",
      "hooks": [
        {
          "type": "command",
          "command": "./.claude/hooks/post-prompt-context.sh"
        }
      ]
    }
  ],
  "Stop": [
    {
      "matcher": "",
      "hooks": [
        {
          "type": "command",
          "command": "./.claude/hooks/completion-one-liners.sh 2>/dev/null || true"
        }
      ]
    }
  ]
}
```

**After:**
```json
"hooks": {
  "UserPromptSubmit": [
    {
      "matcher": "",
      "hooks": [
        {
          "type": "command",
          "command": "./.garden/scripts/context-tracker.py --show-ab"
        }
      ]
    }
  ],
  "Stop": [
    {
      "matcher": "",
      "hooks": [
        {
          "type": "command",
          "command": "./.garden/scripts/completion-status.sh"
        }
      ]
    }
  ]
}
```

## Testing

**Test Context Tracking:**
```bash
# Submit any user message and check for context display at start
# Should see: Context: 🟩🟩🟩... XX% (XXk/200k)
```

**Test Completion Notification:**
```bash
# Wait for Claude to finish a response
# Should hear: "The Garden is finished"
# Should see: Session complete status line
```

**Manual Test Scripts:**
```bash
# Test context tracker directly
./.garden/scripts/context-tracker.py --show-ab

# Test completion status directly
./.garden/scripts/completion-status.sh
```

## Customization

**Disable Audio:**
- Comment out the `say` line in `.garden/scripts/completion-status.sh`

**Change Voice/Message:**
- Edit the `say` command in `.garden/scripts/completion-status.sh`
- Example: `say -v "Samantha" "Work complete"`

**Adjust Status Display:**
- Modify the echo statements in `.garden/scripts/completion-status.sh`
- Use different emojis or formatting

## Troubleshooting

**"say: command not found":**
- This setup requires macOS
- On Linux, replace with `espeak` or `festival`
- Or remove audio line for visual-only

**"jq: command not found":**
- Install: `brew install jq` (macOS)
- Or remove context percentage logic

**Hooks not running:**
- Check `.claude/settings.json` syntax (valid JSON)
- Verify script paths are correct (relative to repo root)
- Ensure scripts are executable: `chmod +x .garden/scripts/*.sh`

**Permissions denied:**
- Add `"Bash(say:*)"` to `permissions.allow` array
- Restart Claude Code if already running

## Files Created

- `.garden/scripts/completion-status.sh` - Completion notification script
- `.garden/.context-state.json` - Context state (auto-created by tracker)
- `.garden/.context-ab-test.jsonl` - A/B test logs (auto-created by tracker)

## Related

- **Context Display Skill**: See `.garden/skills/context-display/SKILL.md`
- **Context Tracker**: See `.garden/scripts/context-tracker.py`
- **Handoff Protocol**: See `.garden/PROTOCOLS.md` (Session Management)
