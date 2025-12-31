#!/bin/bash
# Completion notification with audio + visual status
# Called when Claude finishes a response (Stop hook)

# Audio: Speak random one-liner (uses project-themed lines from /onboard)
ONE_LINER_SCRIPT="${CLAUDE_PROJECT_DIR:-.}/.claude/hooks/completion-one-liners.sh"
if [ -x "$ONE_LINER_SCRIPT" ]; then
    "$ONE_LINER_SCRIPT" &
else
    # Fallback if one-liners not configured
    say "Task complete." &
fi

# Visual: Display rich status summary
# Read context state from .claude/.context-state.json
STATE_FILE=".claude/.context-state.json"

# Default values
CONTEXT_PCT="--"
AGENTS="0"

# Try to read actual state
if [ -f "$STATE_FILE" ]; then
    # Extract percentage from last_actual tokens
    LAST_ACTUAL=$(jq -r '.last_actual // 30000' "$STATE_FILE" 2>/dev/null)
    if [ -n "$LAST_ACTUAL" ] && [ "$LAST_ACTUAL" != "null" ]; then
        CONTEXT_PCT=$((LAST_ACTUAL * 100 / 200000))
    fi
fi

# Count agent session history files if they exist
if [ -d ".claude/work/history" ]; then
    AGENT_COUNT=$(find .claude/work/history -name "*.md" -type f 2>/dev/null | wc -l | tr -d ' ')
    if [ -n "$AGENT_COUNT" ] && [ "$AGENT_COUNT" -gt 0 ]; then
        AGENTS=$AGENT_COUNT
    fi
fi

# Output rich status line
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Session complete | Context: ${CONTEXT_PCT}% | Agents: ${AGENTS} invocations"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
