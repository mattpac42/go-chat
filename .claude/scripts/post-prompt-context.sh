#!/bin/bash
# Automated context tracking via A/B testing
# Tries multiple methods and compares results

# Get the directory where this script lives
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Run the Python context tracker
# It will automatically try all methods and choose the best one
"$SCRIPT_DIR/context-tracker.py" --show-ab

# Note: The tracker saves state and logs A/B test results
# Run with --analyze flag to see which method is most accurate:
#   ./.claude/hooks/context-tracker.py --analyze
