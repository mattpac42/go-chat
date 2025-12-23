# Automated Context Tracking System

## Overview

The Garden implements an intelligent, self-correcting context tracking system that monitors token usage and automatically triggers session handoffs at 75% capacity.

## Architecture

### Three-Method A/B Testing Approach

The system tries three methods simultaneously and chooses the most accurate:

#### Method 1: System Warning Parsing (Primary)
- **How it works**: Parses `<system_warning>Token usage: X/Y</system_warning>` from Claude's responses
- **Accuracy**: High (directly from Claude)
- **Confidence**: High when available
- **Limitation**: Only present in actual Claude responses

#### Method 2: Environment Variables (Secondary)
- **How it works**: Checks `CLAUDE_TOKENS_USED` and `CLAUDE_TOKENS_TOTAL` environment variables
- **Accuracy**: High if exposed by Claude Code
- **Confidence**: High when available
- **Limitation**: Currently not exposed by Claude Code

#### Method 3: Estimation (Fallback)
- **How it works**: Tracks running estimate, adds ~5k per exchange
- **Accuracy**: Moderate (tends to underestimate)
- **Confidence**: Low
- **Limitation**: Drifts over time without correction

### Auto-Selection Logic

```
IF system_warning available:
    USE system_warning (most accurate)
ELSE IF environment variables available:
    USE environment (accurate)
ELSE:
    USE estimation (fallback)
END

Log all three results for A/B comparison
```

## Files

### `.claude/hooks/context-tracker.py`
Main tracking implementation with:
- Three tracking methods
- A/B test logging
- Automatic handoff detection
- State persistence

### `.claude/hooks/post-prompt-context.sh`
Hook that runs after every prompt:
```bash
#!/bin/bash
./.claude/hooks/context-tracker.py --show-ab
```

### `.claude/.context-state.json` (gitignored)
Persisted state tracking:
```json
{
  "estimated_tokens": 103198,
  "last_actual": 103198,
  "last_handoff_pct": 0,
  "preferred_method": "auto",
  "last_updated": "2025-12-13T10:30:00Z"
}
```

### `.claude/.context-ab-test.jsonl` (gitignored)
A/B test results log (JSONL format):
```json
{"timestamp": "2025-12-13T10:30:00Z", "all_methods": {...}, "chosen_method": "system_warning", ...}
{"timestamp": "2025-12-13T10:31:00Z", "all_methods": {...}, "chosen_method": "system_warning", ...}
```

## Usage

### Automatic Tracking

The hook runs automatically after every prompt. No user action needed.

**Output Example**:
```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 52% (103k/200k) [method: system_warning, confidence: high]

📊 A/B Test Results:
  • system_warning: 52% (103k) - confidence: high
  • estimation: 15% (30k) - confidence: low
```

### Automatic Handoff Trigger

When context reaches **75%**, the tracker outputs:

```
Context: 🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨⬛⬛⬛⬛⬛ 75% (150k/200k) 🔄 Session handoff recommended

🔔 AUTOMATIC HANDOFF TRIGGERED: 75% threshold reached
📝 Main agent should invoke /handoff command now
   Run: /handoff
```

**Main Agent Response**:
When the main Claude agent sees this trigger, it should:
1. Acknowledge the trigger
2. Immediately invoke `/handoff` command
3. Create session handoff files
4. Inform user that handoff has been created

### Manual Analysis

Check which tracking method is most accurate:

```bash
./.claude/hooks/context-tracker.py --analyze
```

**Output Example**:
```
📊 A/B Test Analysis Report:
Total samples: 45

  system_warning:
    - Used: 43 times
    - High confidence: 43 (100.0%)

  estimation:
    - Used: 45 times
    - High confidence: 0 (0.0%)

Recommendation:
  Prefer 'system_warning' (most reliable)
```

### Force Specific Method

If analysis shows one method is more reliable:

```bash
# Prefer system_warning (recommended)
./.claude/hooks/context-tracker.py --method system_warning

# Use estimation only
./.claude/hooks/context-tracker.py --method estimation

# Auto-select (default)
./.claude/hooks/context-tracker.py --method auto
```

Update `.claude/.context-state.json` to persist preference:
```json
{
  "preferred_method": "system_warning"
}
```

## Integration with CLAUDE.md

### Automatic Handoff Protocol

**CLAUDE.md Section**:
```markdown
#### Session Handoff Protocol (Automatic)

**Automatic Triggers**:
- **75% Context Usage**: Context tracker automatically detects and notifies
- Main agent responds by invoking /handoff command immediately
- No manual monitoring required

**Handoff Detection**:
The context-tracker.py hook monitors usage and outputs trigger message:
```
🔔 AUTOMATIC HANDOFF TRIGGERED: 75% threshold reached
📝 Main agent should invoke /handoff command now
   Run: /handoff
```

**Main Agent Responsibility**:
When seeing handoff trigger in hook output:
1. Acknowledge: "I've detected we've reached 75% context usage"
2. Execute: Immediately run /handoff command
3. Confirm: "Session handoff files created successfully"
4. Inform: User knows handoff is ready for next session
```

## Display Format

### Standard Display
```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 52% (103k/200k)
```

### With Method Info (--show-ab)
```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 52% (103k/200k) [method: system_warning, confidence: high]
```

### Color Coding
- 🟩 Green (0-60%): Safe zone
- 🟨 Yellow (60-75%): Warning zone, approaching handoff
- 🟧 Orange (75-85%): Handoff recommended
- 🟥 Red (85-100%): Critical, new session strongly recommended

### Status Messages
- **60-75%**: ⚠️ Approaching handoff
- **75-85%**: 🔄 Session handoff recommended
- **85-100%**: 🚨 New session recommended

## Benefits

### For Automated Orchestration
- ✅ No manual `/context` commands needed
- ✅ Automatic handoff detection
- ✅ Self-correcting estimates
- ✅ Works without user intervention

### For Accuracy
- ✅ Multiple methods compared (A/B testing)
- ✅ Always uses most accurate available method
- ✅ Logs results for continuous improvement
- ✅ High-confidence tracking via system warnings

### For Reliability
- ✅ Graceful degradation (three fallback levels)
- ✅ State persistence across sessions
- ✅ Prevents duplicate handoffs
- ✅ Observable via A/B test logs

## Troubleshooting

### Context tracker not running

**Check hook is enabled**:
```bash
cat .claude/settings.json | grep -A 5 UserPromptSubmit
```

Should show:
```json
"UserPromptSubmit": {
  "command": "./.claude/hooks/post-prompt-context.sh"
}
```

### Accuracy seems off

**Run analysis**:
```bash
./.claude/hooks/context-tracker.py --analyze
```

Check which method is being used and its confidence level.

### No automatic handoff at 75%

**Check state file**:
```bash
cat .claude/.context-state.json
```

If `last_handoff_pct` is already 75, handoff was previously triggered.
Reset by setting `last_handoff_pct: 0`.

### Want to see all three methods

**Use --show-ab flag**:
```bash
echo "test" | ./.claude/hooks/context-tracker.py --show-ab
```

This shows comparison of all three methods.

## Future Enhancements

### If Claude Code Exposes Token API

When Claude Code adds environment variables or API endpoints:

```python
# Method 2 would become primary
def method_environment(self):
    tokens_used = int(os.environ['CLAUDE_TOKENS_USED'])
    tokens_total = int(os.environ['CLAUDE_TOKENS_TOTAL'])
    return {
        'used': tokens_used,
        'total': tokens_total,
        'percentage': round((tokens_used / tokens_total) * 100),
        'method': 'environment',
        'confidence': 'high'
    }
```

A/B testing would automatically detect and prefer this method.

### Machine Learning Enhancement

With enough A/B test data:
- Train ML model to predict token usage
- Improve estimation accuracy
- Detect patterns in usage spikes

---

**Status**: Production Ready
**Version**: 1.0.0
**Last Updated**: 2025-12-13
