---
name: version-notify
description: Check for Garden updates and notify when new versions are available. Auto-triggers at session start for planted projects. Use /updates command for manual check or to see changelog.
---

# Version Notification

Detect and display available updates from The Garden.

## Auto-Detection

At session start, if `.claude/lineage.json` exists:
1. Check Garden source path accessibility (or git remote fallback)
2. Compare current vs latest version
3. Display one-line notification if update available
4. Silent skip if Garden not accessible (common in devcontainers)

## Manual Commands

| Command | Action |
|---------|--------|
| `/updates` | Check for updates, show summary |
| `/updates --changelog` | Show detailed changes |
| `/updates --dismiss` | Hide notification until next version |

## Notification Display

### Minor Update (no breaking changes)
```
============================================================
Garden UPDATE AVAILABLE: 3.1.0 -> 3.2.0
============================================================

ADDED:
   - Version notification system
   - New marketplace agents

CHANGED:
   - Improved sync-baseline performance

Run `/sync-baseline` to update or `/updates --dismiss` to hide
============================================================
```

### Breaking Update
```
============================================================
Garden BREAKING CHANGES: 3.1.0 -> 4.0.0
============================================================

BREAKING:
   - lineage.json schema changed to v3.0
   - Removed deprecated handoff skill

Run `/sync-baseline` to update (review breaking changes first)
============================================================
```

## Version Sources

The skill checks for Garden version in order:
1. Local path from `lineage.json` -> reads `VERSION` file directly
2. Git remote named `garden` -> `git show garden/main:VERSION`

## Configuration

In `lineage.json`:
```json
{
  "sync": {
    "auto_notify": true
  },
  "notifications": {
    "dismissed_versions": ["3.1.5"],
    "last_check": "2025-01-02T10:00:00Z"
  }
}
```

## Suppressing Notifications

1. **Dismiss specific version**: `/updates --dismiss 3.2.0`
2. **Disable auto-check**: Set `sync.auto_notify: false` in lineage.json
