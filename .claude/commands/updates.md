# Updates Command

Check for available updates from The Garden.

## Instructions for Claude

When this command is invoked:

### Step 1: Verify Garden Connection

Check if this is a Garden-rooted project:

```bash
cat .claude/lineage.json 2>/dev/null
```

**If no lineage.json**: Tell the user "This project is not connected to The Garden. Use `/plant` to create a new Garden-rooted project."

### Step 2: Run Version Check

Run the update checker script:

```bash
python3 .claude/skills/version-notify/scripts/check_updates.py --changelog
```

If the script doesn't exist, the project may need to sync first:
```bash
# Check if Garden path is accessible
jq -r '.garden.source_path' .claude/lineage.json
```

### Step 3: Display Results

**If up to date:**
```
Garden is up to date (v3.2.0)
  Source: local

Last checked: 2025-01-02 10:00:00
```

**If update available:**
Display the notification from the script output, which will show:
- Current vs latest version
- Breaking changes (if any)
- Added/changed/fixed items
- Instructions to sync

### Step 4: Offer Actions

Based on the result:

- **Update available**: "Run `/sync-baseline` to apply updates"
- **Breaking changes**: "Review the breaking changes above before syncing"
- **Garden not accessible**: "Configure a git remote with `git remote add garden <url>` or ensure the Garden path exists"

## Arguments

| Argument | Effect |
|----------|--------|
| (none) | Check and display summary with changelog |
| `--dismiss` | Dismiss current notification until next version |
| `--json` | Output as JSON for programmatic use |

## Examples

```bash
# Check for updates
/updates

# Dismiss notification for current version
/updates --dismiss

# Get JSON output
/updates --json
```

## Troubleshooting

### Garden path not accessible

If you're in a devcontainer and the Garden path points to the host machine:

1. Configure a git remote:
   ```bash
   git remote add garden https://github.com/your-org/the-garden.git
   ```

2. Or mount The Garden in your devcontainer configuration

### Script not found

If the version-notify skill is missing, sync from Garden:
```bash
/sync-baseline
```

This will copy the latest skills including version-notify.
