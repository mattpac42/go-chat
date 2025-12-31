#!/bin/bash
# Test script for context-tracker.py with context_command method

echo "Testing context-tracker.py with new context_command method"
echo "=========================================================="
echo ""

echo "Test 1: Parse /context output via command line argument"
python3 .claude/hooks/context-tracker.py --context-output "claude-sonnet-4-5-20250929 · 44k/200k tokens (22%)" --show-ab
echo ""

echo "Test 2: Parse /context output via environment variable"
CLAUDE_CONTEXT_OUTPUT="claude-sonnet-4-5-20250929 · 75k/200k tokens (37%)" python3 .claude/hooks/context-tracker.py --show-ab
echo ""

echo "Test 3: Fallback to estimation when no context output available"
python3 .claude/hooks/context-tracker.py --show-ab
echo ""

echo "Test 4: Priority test - context_command wins over system_warning"
echo "<system_warning>Token usage: 50000/200000; 150000 remaining</system_warning>" | python3 .claude/hooks/context-tracker.py --context-output "claude-sonnet-4-5-20250929 · 30k/200k tokens (15%)" --show-ab
echo ""

echo "Test 5: Handoff trigger at 75% threshold"
rm -f .claude/.context-state.json
python3 .claude/hooks/context-tracker.py --context-output "claude-sonnet-4-5-20250929 · 150k/200k tokens (75%)"
echo ""

echo "Test 6: Invalid context output - graceful fallback"
python3 .claude/hooks/context-tracker.py --context-output "invalid format" --show-ab
echo ""

echo "Test 7: A/B test analysis"
python3 .claude/hooks/context-tracker.py --analyze
echo ""

echo "All tests completed!"
