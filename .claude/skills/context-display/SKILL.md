---
name: context-display
description: Show token usage visualization when context exceeds 50% or after agent task completion. Reduces visual noise for shorter sessions while ensuring users are warned before handoff is needed.
---

# Context Display

Show token usage when context exceeds 50% or after agent completion.

## When to Display

Display context bar when:
- Context usage is **50% or higher**
- An agent task has just completed
- User explicitly requests context status

Skip display when:
- Context is below 50% AND no agent just completed
- Reduces noise for short sessions

## Output Format

```
Context: [emoji-bar] XX% (XXXk/200k) [status]
```

## Visualization

20 blocks total, each = 5%:
- Blocks 1-12 (0-60%): 🟩 Green
- Blocks 13-15 (60-75%): 🟨 Yellow
- Blocks 16-17 (75-85%): 🟧 Orange
- Blocks 18-20 (85-100%): 🟥 Red
- Empty blocks: ⬛

## Status Messages

| Percentage | Message |
|------------|---------|
| 60-75% | ⚠️ Approaching handoff |
| 75-85% | 🔄 Session handoff recommended |
| 85%+ | 🚨 New session recommended |

## Example

```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟨🟨🟨⬛⬛⬛⬛⬛ 75% (150k/200k) 🔄 Session handoff recommended
```

## Data Source

Use `/context` command output for accurate values. Fall back to estimation if unavailable.
