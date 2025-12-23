# Merge Request Command

Create a GitLab merge request with auto-generated title, description, and proper formatting.

## Usage

```bash
/mr                        # Auto-generate title/description from commits
/mr "Title here"           # Custom title, auto-generate description
/mr --draft                # Create as draft/WIP MR
/mr --target develop       # Target specific branch (default: main)
/mr --no-push              # Don't push before creating MR
```

## What This Command Does

1. Ensures current branch is pushed to remote
2. Detects target branch (main/master/develop)
3. Generates MR title from branch name or commits
4. Generates description from commit history
5. Adds standard checklist and formatting
6. Creates MR via `gh` or `glab` CLI
7. Returns MR URL for review

## Instructions for Claude

### Step 1: Parse Arguments

Extract flags and custom title from user input:
- `--draft`: Create as draft/WIP merge request
- `--target [branch]`: Target branch for merge (default: main)
- `--no-push`: Skip push step (assume already pushed)
- Quoted text: Custom MR title (e.g., `/mr "Add user authentication"`)

### Step 2: Validate Git State

Run these commands in parallel:
```bash
git branch --show-current
git remote -v
git status --short
```

**Validations**:
- Must be on a feature branch (not main/master)
- Must have a remote configured
- Working tree should be clean (warn if dirty)

```
# Git State

Branch: [branch-name]
Remote: [remote-name] ([url])
Status: [clean | X uncommitted changes]
```

**If on main/master**:
```
Cannot create MR from [branch] branch.
Please checkout a feature branch first.
```

**If uncommitted changes**:
```
You have uncommitted changes. Options:
1. Run /commit first to commit changes
2. Proceed anyway (changes won't be in MR)
3. Cancel and handle manually

Proceed with MR creation? (y/n)
```

### Step 3: Ensure Branch is Pushed

**If --no-push NOT set**:

1. Check if remote tracking exists:
```bash
git rev-parse --abbrev-ref @{upstream} 2>/dev/null || echo "no-upstream"
```

2. If no upstream, push with -u:
```bash
REMOTE=$(git remote | head -1)
BRANCH=$(git branch --show-current)
git push -u $REMOTE $BRANCH
```

3. If upstream exists, push any new commits:
```bash
git push
```

Confirm: "Branch pushed to remote."

### Step 4: Detect Target Branch

**If --target specified**: Use provided branch.

**Otherwise, auto-detect**:
```bash
# Check for common default branches
git show-ref --verify refs/remotes/origin/main 2>/dev/null && echo "main"
git show-ref --verify refs/remotes/origin/master 2>/dev/null && echo "master"
git show-ref --verify refs/remotes/origin/develop 2>/dev/null && echo "develop"
```

Priority: `main` > `master` > `develop`

If none found, ask user:
```
Could not detect target branch. Which branch should this MR target?
1. main
2. master
3. develop
4. Other: [specify]
```

### Step 5: Generate MR Title

**If user provided custom title**: Use it directly.

**Otherwise, auto-generate from branch name**:

1. Get branch name:
```bash
git branch --show-current
```

2. Parse branch name patterns:

| Pattern | Generated Title |
|---------|-----------------|
| `feature/add-auth` | "Add auth" |
| `fix/login-bug` | "Fix login bug" |
| `123-add-feature` | "Add feature (#123)" |
| `feat/user-profile` | "feat: user profile" |
| `bugfix/timeout-issue` | "fix: timeout issue" |

3. Transformation rules:
   - Remove prefix (`feature/`, `fix/`, `bugfix/`, `feat/`, `hotfix/`)
   - Replace hyphens/underscores with spaces
   - Capitalize first letter
   - Extract issue numbers and append as `(#123)`
   - Keep conventional commit prefix if present

### Step 6: Generate MR Description

Build description from commit history:

```bash
# Get commits unique to this branch
TARGET_BRANCH="main"
CURRENT_BRANCH=$(git branch --show-current)
git log $TARGET_BRANCH..$CURRENT_BRANCH --pretty=format:"- %s" --reverse
```

**Description template**:
```markdown
## Summary

[Auto-generated from branch name or first commit message]

## Changes

[List of commits as bullet points]

## Checklist

- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] TEMPLATE_UPDATES.md updated (if baseline repo)
- [ ] Self-review completed
- [ ] Ready for review

---
Generated with [Claude Code](https://claude.com/claude-code)
```

**If draft MR**: Add `[DRAFT]` prefix to title or use draft flag.

### Step 7: Detect CLI Tool

Check which CLI tool is available:

```bash
# Check for GitLab CLI
which glab 2>/dev/null && echo "glab"

# Check for GitHub CLI (also works with GitLab)
which gh 2>/dev/null && echo "gh"
```

**If neither found**:
```
No CLI tool found. Please install one of:
- glab (GitLab CLI): https://gitlab.com/gitlab-org/cli
- gh (GitHub CLI): https://cli.github.com

Or create MR manually at:
[remote-url]/-/merge_requests/new?merge_request[source_branch]=[branch]
```

Provide the manual URL and stop.

### Step 8: Create Merge Request

**Using glab (GitLab CLI)**:
```bash
glab mr create \
  --title "[title]" \
  --description "[description]" \
  --target-branch "[target]" \
  [--draft if draft flag set]
```

**Using gh with GitLab** (if configured):
```bash
gh mr create \
  --title "[title]" \
  --body "[description]" \
  --base "[target]" \
  [--draft if draft flag set]
```

**If no CLI available, provide manual instructions**:
```
# Manual MR Creation

Please create MR at:
[URL]

Title: [title]

Description:
[description]

Target branch: [target]
```

### Step 9: Display Summary

**Success output format**:
```
# Merge Request Created

Branch: [source-branch] -> [target-branch]
Title: [title]
Status: [Draft | Ready for review]

URL: [mr-url]

Next steps:
- Review changes in GitLab
- Request reviewers
- Address feedback
- Merge when approved
```

**If CLI not available**:
```
# Manual MR Required

Copy this information to create your MR:

URL: [manual-create-url]
Title: [title]
Target: [target-branch]

Description:
---
[full description]
---
```

## Title Generation Rules

| Branch Pattern | Example | Generated Title |
|----------------|---------|-----------------|
| `feature/*` | `feature/add-oauth` | "Add oauth" |
| `fix/*` | `fix/auth-timeout` | "Fix auth timeout" |
| `bugfix/*` | `bugfix/memory-leak` | "Fix memory leak" |
| `hotfix/*` | `hotfix/critical-bug` | "Hotfix: critical bug" |
| `feat/*` | `feat/dark-mode` | "feat: dark mode" |
| `chore/*` | `chore/update-deps` | "chore: update deps" |
| `docs/*` | `docs/api-guide` | "docs: api guide" |
| `[number]-*` | `123-add-feature` | "Add feature (#123)" |
| `[number]/*` | `PROJ-456/add-auth` | "Add auth (PROJ-456)" |

## Description Template

```markdown
## Summary

[Brief description - from first commit or branch name]

## Changes

- [Commit 1 message]
- [Commit 2 message]
- [Commit 3 message]

## Type of Change

- [ ] Bug fix (non-breaking change fixing an issue)
- [ ] New feature (non-breaking change adding functionality)
- [ ] Breaking change (fix or feature causing existing functionality to change)
- [ ] Documentation update
- [ ] Refactoring (no functional changes)

## Checklist

- [ ] My code follows the project style guidelines
- [ ] I have performed a self-review
- [ ] I have added tests that prove my fix/feature works
- [ ] New and existing tests pass locally
- [ ] I have updated documentation as needed
- [ ] TEMPLATE_UPDATES.md updated (baseline repo only)

---
Generated with [Claude Code](https://claude.com/claude-code)
```

## Error Handling

| Error | Detection | Response |
|-------|-----------|----------|
| On main/master branch | Branch name check | "Cannot create MR from main. Checkout feature branch first." |
| No remote configured | `git remote` empty | "No remote configured. Add remote first: git remote add origin [url]" |
| Push failed | Non-zero exit | Show error, suggest `git pull` or force push |
| CLI not installed | `which` returns empty | Provide manual URL and instructions |
| MR already exists | CLI returns error | Show existing MR URL |
| Target branch doesn't exist | `git show-ref` fails | Ask user to specify valid target |

## Examples

**Example 1: Basic MR creation**
```
User: /mr
Claude:
  - Validates on feature/add-auth branch
  - Pushes to origin
  - Detects target: main
  - Generates title: "Add auth"
  - Generates description from 3 commits
  - Creates MR via glab
  - Returns: https://gitlab.com/repo/-/merge_requests/42
```

**Example 2: Custom title with draft**
```
User: /mr "feat: implement OAuth2 authentication" --draft
Claude:
  - Uses provided title
  - Creates as draft MR
  - Returns draft MR URL
```

**Example 3: Different target branch**
```
User: /mr --target develop
Claude:
  - Validates develop branch exists
  - Creates MR targeting develop instead of main
  - Returns MR URL
```

**Example 4: No CLI available**
```
User: /mr
Claude:
  - Detects no glab/gh installed
  - Generates title and description
  - Provides manual creation URL
  - Shows formatted description to copy
```

## GitLab-Specific Features

When using `glab`:
- Supports labels: `--label "priority::high"`
- Supports assignees: `--assignee @username`
- Supports milestones: `--milestone "v1.0"`
- Supports reviewers: `--reviewer @username`

Future enhancement: Add flags for these options.

## Notes for Claude

- Always check branch state before proceeding
- Push branch before creating MR (unless --no-push)
- Handle both GitLab (glab) and GitHub (gh) CLIs
- Provide manual fallback when CLI unavailable
- Keep descriptions concise but informative
- Include checklist for reviewer guidance
- Respect draft status for WIP work
- Extract issue numbers from branch names when possible
- Use conventional commit prefixes when detected in branch name
