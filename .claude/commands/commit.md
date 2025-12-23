# Commit Command

Automate the full commit workflow with changelog awareness, conventional commits, and push.

## Usage

```bash
/commit                     # Auto-generate message, full workflow
/commit "message"           # Use provided message
/commit --no-push           # Commit only, don't push
/commit --amend             # Amend previous commit
/commit --skip-changelog    # Skip changelog check
```

## What This Command Does

1. Shows current git status (branch, changes)
2. Auto-stages all changes with `git add -A`
3. Checks if changelog update needed (baseline repos only)
4. Generates conventional commit message (if not provided)
5. Creates commit with standard footer
6. Pushes to current branch (unless --no-push)
7. Displays summary with commit hash and next steps

## Instructions for Claude

### Step 1: Parse Arguments

Extract flags and custom message from user input:
- `--no-push`: Skip push after commit
- `--amend`: Amend previous commit instead of new commit
- `--skip-changelog`: Skip changelog requirement check
- Quoted text: Custom commit message (e.g., `/commit "fix: resolve auth bug"`)

### Step 2: Check Git Status

Run these commands in parallel:
```bash
git status --short
git branch --show-current
git diff --staged --name-only
git diff --name-only
```

Display current state:
```
# Current Status

Branch: [branch-name]
Staged: [count] files
Unstaged: [count] files
```

**Safety Check**: If current branch is `main` or `master`:
```
⚠️  WARNING: You are on the [branch] branch.
Are you sure you want to commit directly to [branch]? (y/n)
```

Wait for user confirmation before proceeding.

### Step 3: Auto-Stage Changes

Run: `git add -A`

Confirm: "Staged all changes for commit."

### Step 4: Changelog Check (Baseline Repos Only)

**Detect if baseline repo**:
```bash
# Check for TEMPLATE_UPDATES.md existence
[ -f "TEMPLATE_UPDATES.md" ] && echo "baseline" || echo "downstream"
```

**If baseline repo AND --skip-changelog NOT set**:

1. Check if trigger files were modified:
```bash
git diff --staged --name-only | grep -E "(CLAUDE\.md|\.claude/templates/|\.claude/docs/|\.claude/commands/)"
```

2. If trigger files found, check if `TEMPLATE_UPDATES.md` was modified:
```bash
git diff --staged --name-only | grep "TEMPLATE_UPDATES.md"
```

3. If trigger files changed but `TEMPLATE_UPDATES.md` NOT modified:
```
📝 CHANGELOG REMINDER

You modified framework files that require changelog entries:
[list trigger files]

Please update TEMPLATE_UPDATES.md with:
- What changed
- Why it changed
- Impact on downstream repos

Options:
1. Add changelog entry and re-run /commit
2. Use /commit --skip-changelog to bypass (not recommended)
```

Stop and wait for user action. Do NOT proceed with commit.

**If downstream repo**: Skip changelog check automatically (no message needed).

### Step 5: Generate Commit Message (if not provided)

**If user provided custom message**: Use it as-is (validate conventional format optional).

**If no message provided**: Auto-generate using this logic:

1. Get list of changed files:
```bash
git diff --staged --name-status
```

2. Apply commit type detection rules:

| Condition | Type | Example Message |
|-----------|------|-----------------|
| New files in features/, new functionality | `feat:` | "feat: add user authentication" |
| Files in fix/, bug/, error corrections | `fix:` | "fix: resolve login timeout" |
| Only .md files in docs/ or README | `docs:` | "docs: update installation guide" |
| .claude/templates/, .claude/commands/, configs | `chore:` | "chore: update agent templates" |
| Code moved/renamed, no logic change | `refactor:` | "refactor: restructure auth module" |
| Only test files (test_*, *_test.py, *.test.js) | `test:` | "test: add auth unit tests" |
| CLAUDE.md, .claude/docs/ | `docs:` | "docs: update agent protocols" |
| Multiple types mixed | `chore:` | "chore: multiple improvements" |

3. Generate concise summary (50 chars max for subject line)

4. If large commit (>10 files), ask for confirmation:
```
📦 LARGE COMMIT DETECTED

Files changed: [count]
Types: [list file types]

Generated message: "[type]: [summary]"

Proceed with this commit? (y/n)
Or provide custom message: /commit "your message here"
```

### Step 6: Create Commit

**Standard commit format**:
```bash
git commit -m "[type]: [summary]

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

**If --amend flag set**:

1. Check authorship first:
```bash
git log -1 --format='%an %ae'
```

2. If author is NOT "Claude <noreply@anthropic.com>":
```
⚠️  WARNING: Amending commit by [author name]

Last commit: [commit message]
Author: [author name] <[email]>

Are you sure you want to amend someone else's commit? (y/n)
```

3. If confirmed or Claude's commit:
```bash
git commit --amend --no-edit
# OR if new message provided:
git commit --amend -m "[new message]"
```

**Capture commit hash**:
```bash
git log -1 --format='%h'
```

### Step 7: Push to Remote (unless --no-push)

**If --no-push flag**: Skip push, show local commit summary only.

**Otherwise**:

1. Get current branch and remote:
```bash
BRANCH=$(git branch --show-current)
REMOTE=$(git config branch.$BRANCH.remote || echo "origin")
```

2. Push to remote:
```bash
git push $REMOTE $BRANCH
```

3. If push fails (no upstream, permissions, conflicts):
```
🚨 PUSH FAILED

Error: [error message]

Commit was created locally but not pushed.

Options:
1. Set upstream: git push -u origin [branch]
2. Force push: git push --force (NOT RECOMMENDED)
3. Pull first: git pull origin [branch]

Run one of these commands manually or re-run /commit after resolving.
```

### Step 8: Display Summary

**Success output format**:
```
# ✅ Commit Successful

Branch: [branch-name]
Files changed: [count] ([added] added, [modified] modified, [deleted] deleted)

Commit: [hash]
Message: [commit message]

📝 Changelog: [TEMPLATE_UPDATES.md updated ✓ | Skipped (downstream repo) | Skipped (--skip-changelog)]
🚀 [Pushed to origin/[branch] | Local only (--no-push)]

Next steps:
- Create PR: /pr
- Continue development
- View changes: git show [hash]
```

## Commit Type Detection Rules

Reference table for auto-detection:

| Files Changed | Commit Type | Keywords to Detect |
|--------------|-------------|-------------------|
| New features, functionality | `feat:` | new, add, create, implement |
| Bug fixes, corrections | `fix:` | fix, resolve, correct, patch |
| Documentation only | `docs:` | *.md in docs/, README |
| Templates, configs, tools | `chore:` | .claude/, config, setup |
| Code restructure | `refactor:` | rename, move, restructure |
| Tests only | `test:` | test_*, *.test.*, spec |
| Performance improvements | `perf:` | optimize, performance, speed |
| CI/CD changes | `ci:` | .gitlab-ci.yml, .github/ |

## Baseline Detection Logic

Determine if this is The Garden baseline repo:

```bash
# Method 1: Check for TEMPLATE_UPDATES.md
if [ -f "TEMPLATE_UPDATES.md" ]; then
  echo "baseline"
fi

# Method 2: Check git remote URL
git config --get remote.origin.url | grep -q "the_garden" && echo "baseline"

# Method 3: Check for marker in CLAUDE.md
grep -q "# GARDEN BASELINE REPO" CLAUDE.md && echo "baseline"
```

If ANY method returns "baseline" → require changelog check.
If ALL methods fail → assume downstream repo, skip changelog.

## Safety Checks

Implement these safety protocols:

1. **Main/Master Branch Protection**
   - Warn before committing to `main` or `master`
   - Require explicit confirmation (y/n)

2. **Large Commit Confirmation**
   - Trigger on >10 files changed
   - Show file count and types
   - Allow cancel or custom message

3. **Amend Authorship**
   - Check if amending someone else's commit
   - Warn and require confirmation

4. **No Force Push Ever**
   - Never use `git push --force`
   - Suggest alternatives if push fails

5. **Pre-Commit Hooks**
   - Respect existing git hooks
   - Don't skip hooks with --no-verify

## Error Handling

Handle these common scenarios:

| Error | Detection | Response |
|-------|-----------|----------|
| No changes to commit | `git status` shows clean | "No changes to commit. Working tree clean." |
| Merge conflict | `git status` shows conflict | "Resolve merge conflicts before committing." |
| Detached HEAD | `git branch` shows detached | "Create branch before committing: git checkout -b [name]" |
| Push rejected | Push returns non-zero | Show error, suggest pull or upstream setup |
| Invalid message format | User message doesn't match conventional | Warn but allow (don't enforce strictly) |

## Examples

**Example 1: Basic auto-commit**
```
User: /commit
Claude:
  - Shows status (3 files staged)
  - Auto-generates: "feat: add commit automation command"
  - Creates commit with standard footer
  - Pushes to origin/feature-branch
  - Shows summary with commit hash
```

**Example 2: Custom message, no push**
```
User: /commit "fix: resolve auth timeout bug" --no-push
Claude:
  - Shows status
  - Uses provided message
  - Creates commit locally
  - Skips push
  - Shows local-only summary
```

**Example 3: Baseline repo with changelog**
```
User: /commit
Claude:
  - Detects baseline repo
  - Sees .claude/templates/agent.md changed
  - Checks TEMPLATE_UPDATES.md NOT modified
  - Shows changelog reminder
  - STOPS and waits for user to add changelog
```

**Example 4: Amend previous commit**
```
User: /commit --amend "fix: resolve auth bug (updated)"
Claude:
  - Checks last commit author (Claude)
  - Amends commit with new message
  - Pushes amended commit
  - Shows updated summary
```

## Notes for Claude

- Always run git commands from repository root
- Use `git status --short` for concise output
- Capture command output for error handling
- Don't proceed if safety checks fail without confirmation
- Be helpful with error messages and suggest fixes
- Keep output concise but informative
- Validate branch exists before pushing
- Handle both SSH and HTTPS git remotes
- Respect user's git config (signing, hooks, etc.)
