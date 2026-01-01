#!/bin/bash
# Quick version check for session hooks
# Outputs nothing if up-to-date or inaccessible, one-line notification if update available
#
# Usage: Run from project root at session start
# Exit codes: 0 = no update or inaccessible, 1 = update available

set -e

LINEAGE=".claude/lineage.json"

# Check if lineage.json exists
if [[ ! -f "$LINEAGE" ]]; then
    exit 0
fi

# Check if jq is available
if ! command -v jq &> /dev/null; then
    exit 0
fi

# Read current version and garden path
CURRENT_VERSION=$(jq -r '.garden.version // "unknown"' "$LINEAGE" 2>/dev/null)
GARDEN_PATH=$(jq -r '.garden.source_path // empty' "$LINEAGE" 2>/dev/null)
GARDEN_REMOTE=$(jq -r '.garden.remote_url // empty' "$LINEAGE" 2>/dev/null)
AUTO_NOTIFY=$(jq -r '.sync.auto_notify // true' "$LINEAGE" 2>/dev/null)

# Check if auto-notify is disabled
if [[ "$AUTO_NOTIFY" == "false" ]]; then
    exit 0
fi

# Function to get latest version from local path
get_local_version() {
    local path="$1"
    if [[ -d "$path" && -f "$path/VERSION" ]]; then
        cat "$path/VERSION" | tr -d '[:space:]'
        return 0
    fi
    return 1
}

# Function to get latest version from git remote
get_remote_version() {
    # Check if 'garden' remote exists
    if git remote get-url garden &>/dev/null 2>&1; then
        # Fetch quietly and get VERSION from remote
        if git fetch garden --quiet 2>/dev/null; then
            git show garden/main:VERSION 2>/dev/null | tr -d '[:space:]'
            return 0
        fi
    fi
    return 1
}

# Try to get latest version
LATEST_VERSION=""

# First try local path
if [[ -n "$GARDEN_PATH" ]]; then
    LATEST_VERSION=$(get_local_version "$GARDEN_PATH" 2>/dev/null || true)
fi

# Fall back to git remote if local path didn't work
if [[ -z "$LATEST_VERSION" ]]; then
    LATEST_VERSION=$(get_remote_version 2>/dev/null || true)
fi

# If we couldn't get a version, exit silently
if [[ -z "$LATEST_VERSION" ]]; then
    exit 0
fi

# Compare versions (simple string comparison works for semver)
if [[ "$CURRENT_VERSION" == "unknown" || "$CURRENT_VERSION" != "$LATEST_VERSION" ]]; then
    # Check if this version was dismissed
    DISMISSED=$(jq -r ".notifications.dismissed_versions // [] | index(\"$LATEST_VERSION\")" "$LINEAGE" 2>/dev/null)
    if [[ "$DISMISSED" != "null" && -n "$DISMISSED" ]]; then
        exit 0
    fi

    echo "Garden update: $CURRENT_VERSION -> $LATEST_VERSION (run /updates)"
    exit 1
fi

exit 0
