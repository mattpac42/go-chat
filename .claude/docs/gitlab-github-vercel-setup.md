# GitLab → GitHub → Vercel Pipeline Setup Guide

This guide explains how to set up a CI/CD pipeline that:
1. Runs quality gates (lint, test, security scan) on GitLab
2. Mirrors passing code to GitHub
3. Auto-deploys to Vercel from GitHub

## Prerequisites

- GitLab repository (your primary origin)
- GitHub repository (mirror target)
- Vercel account
- Node.js project (this guide assumes Next.js, but adaptable)

## Architecture

```
Developer pushes to GitLab
         │
         ▼
┌─────────────────────────┐
│     GitLab CI/CD        │
│  ┌─────┐ ┌─────┐ ┌────┐ │
│  │lint │→│test │→│scan│ │
│  └─────┘ └─────┘ └────┘ │
│           │             │
│           ▼             │
│       ┌────────┐        │
│       │ mirror │        │
│       └────────┘        │
└───────────│─────────────┘
            │
            ▼
┌─────────────────────────┐
│        GitHub           │
│   (receives mirror)     │
└───────────│─────────────┘
            │
            ▼
┌─────────────────────────┐
│        Vercel           │
│    (auto-deploys)       │
└─────────────────────────┘
```

## Step 1: Create GitHub Repository

If you don't have one already:

1. Go to https://github.com/new
2. Create repository (can be empty)
3. Note the URL: `https://github.com/USERNAME/REPO.git`

## Step 2: Link Vercel to GitHub

1. Go to https://vercel.com
2. Click "Add New Project"
3. Import your GitHub repository
4. Configure build settings:
   - Framework: Auto-detect (or select yours)
   - Root Directory: `./` (or your app directory)
5. Deploy

Vercel will now auto-deploy whenever GitHub receives a push.

## Step 3: Set Up GitLab Runner

You need a GitLab runner with Docker executor. If you don't have one:

### Option A: Use Shared Runners
GitLab.com provides shared runners. Skip to Step 4.

### Option B: Register a New Runner

On your runner host:

```bash
# Install GitLab Runner (if not installed)
curl -L https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh | sudo bash
sudo apt install gitlab-runner

# Register new runner
sudo gitlab-runner register
```

When prompted:
- **URL**: Your GitLab instance URL
- **Token**: Get from GitLab → Project → Settings → CI/CD → Runners
- **Description**: `project-name-nodejs`
- **Tags**: `nodejs`
- **Executor**: `docker`
- **Default image**: `node:20-alpine`

Edit `/etc/gitlab-runner/config.toml` and ensure clean volumes:
```toml
[runners.docker]
  volumes = ["/cache"]
```

Restart: `sudo gitlab-runner restart`

## Step 4: Generate SSH Key for GitHub

On your local machine:

```bash
# Generate key pair (no passphrase)
ssh-keygen -t ed25519 -C "gitlab-ci@yourproject.app" -f ~/.ssh/gitlab_to_github -N ""

# View private key (for GitLab)
cat ~/.ssh/gitlab_to_github

# View public key (for GitHub)
cat ~/.ssh/gitlab_to_github.pub
```

## Step 5: Add Deploy Key to GitHub

1. Go to GitHub → Your Repo → Settings → Deploy keys
2. Click "Add deploy key"
3. Title: `gitlab-ci-mirror`
4. Key: Paste the **public key** content (from `.pub` file)
5. **Check "Allow write access"** ← Critical!
6. Click "Add key"

## Step 6: Add SSH Key to GitLab CI Variables

1. Go to GitLab → Your Project → Settings → CI/CD → Variables
2. Click "Add variable"
3. Configure:
   - **Key**: `GITHUB_SSH_KEY`
   - **Value**: Paste the **private key** content (entire file including BEGIN/END lines)
   - **Type**: `File`
   - **Protected**: ✅ Yes
   - **Masked**: Leave unchecked (can't mask multiline)
   - **Expand variable reference**: ❌ **UNCHECK THIS** ← Critical!
4. Click "Add variable"

## Step 7: Create GitLab CI Pipeline

Create `.gitlab-ci.yml` in your project root:

```yaml
# GitLab CI/CD Pipeline
# Workflow: GitLab → Quality Gates → GitHub Mirror → Vercel Auto-Deploy

stages:
  - lint
  - test
  - scan
  - mirror

# Use your runner tag (or remove for shared runners)
default:
  tags:
    - nodejs

cache:
  key:
    files:
      - package-lock.json
  paths:
    - node_modules/
  policy: pull-push

# Stage 1: Lint
lint:
  stage: lint
  image: node:20-alpine
  script:
    - npm ci --prefer-offline --no-audit
    - npm run lint
  only:
    - merge_requests
    - main

# Stage 2: Test
test:
  stage: test
  image: node:20-alpine
  script:
    - npm ci --prefer-offline --no-audit
    - npm run test -- --ci --coverage --maxWorkers=2
  cache:
    key:
      files:
        - package-lock.json
    paths:
      - node_modules/
    policy: pull
  artifacts:
    paths:
      - coverage/
    expire_in: 30 days
  only:
    - merge_requests
    - main

# Stage 3: Security Scan
scan:
  stage: scan
  image:
    name: aquasec/trivy:latest
    entrypoint: [""]
  script:
    - trivy fs --exit-code 0 --severity HIGH,CRITICAL --format table .
    - trivy fs --exit-code 0 --severity HIGH,CRITICAL --format json --output trivy-report.json .
  artifacts:
    paths:
      - trivy-report.json
    expire_in: 30 days
  allow_failure: true
  only:
    - merge_requests
    - main

# Stage 4: Mirror to GitHub
mirror:
  stage: mirror
  image: alpine:latest
  before_script:
    - apk add --no-cache git openssh-client
    - mkdir -p ~/.ssh && chmod 700 ~/.ssh
    - cp "$GITHUB_SSH_KEY" ~/.ssh/id_rsa && chmod 600 ~/.ssh/id_rsa
    - ssh-keyscan -H github.com >> ~/.ssh/known_hosts 2>/dev/null
    - git config --global user.email "gitlab-ci@yourproject.app"
    - git config --global user.name "GitLab CI"
  script:
    - git remote add github git@github.com:USERNAME/REPO.git || true
    - git push github HEAD:main --force
    - echo "✅ Mirrored to GitHub → Vercel will auto-deploy"
  only:
    - main
  when: on_success
```

**Replace these values:**
- `USERNAME/REPO` → Your GitHub username/repo
- `gitlab-ci@yourproject.app` → Your project email
- `nodejs` tag → Your runner tag (or remove `tags:` section for shared runners)

## Step 8: Customize for Your Project

### If Not Using Next.js

Update the lint/test commands:

```yaml
# Python example
lint:
  image: python:3.12-slim
  script:
    - pip install ruff
    - ruff check .

test:
  image: python:3.12-slim
  script:
    - pip install pytest
    - pip install -r requirements.txt
    - pytest
```

### If No Tests Yet

Remove or comment out the test stage, or make it optional:

```yaml
test:
  # ... existing config ...
  allow_failure: true  # Won't block pipeline
```

### If Using Different Branch

Change `main` to your branch name:

```yaml
only:
  - merge_requests
  - your-branch-name
```

## Step 9: Commit and Test

```bash
git add .gitlab-ci.yml
git commit -m "ci: add GitLab pipeline with GitHub mirror"
git push origin main
```

Watch the pipeline at: `https://your-gitlab.com/project/-/pipelines`

## Troubleshooting

### SSH Key Error: "error in libcrypto"

**Cause**: "Expand variable reference" is checked in GitLab
**Fix**: Edit the variable, uncheck "Expand variable reference"

### SSH Key Error: "Permission denied"

**Cause**: Deploy key doesn't have write access
**Fix**: Delete and re-add the GitHub deploy key with "Allow write access" checked

### Pipeline Not Running on Runner

**Cause**: Runner tags don't match
**Fix**: Either remove `tags:` from pipeline or ensure runner has matching tag

### Mirror Fails: "Repository not found"

**Cause**: Wrong GitHub URL
**Fix**: Verify the git@github.com:USERNAME/REPO.git URL is correct

### Vercel Not Deploying

**Cause**: Vercel not connected to GitHub repo
**Fix**: Go to Vercel dashboard, ensure project is linked to the correct GitHub repo

## Verification Checklist

- [ ] GitLab runner registered and online
- [ ] SSH key pair generated
- [ ] Public key added to GitHub with write access
- [ ] Private key added to GitLab (File type, Expand unchecked)
- [ ] `.gitlab-ci.yml` committed
- [ ] Pipeline runs successfully
- [ ] Mirror stage pushes to GitHub
- [ ] Vercel deploys automatically

## Security Notes

1. **Protected Variables**: Keep `GITHUB_SSH_KEY` protected so it's only available on protected branches
2. **Deploy Keys**: Use deploy keys (not personal access tokens) for least-privilege access
3. **Scan Results**: Review Trivy security scan results regularly
4. **Force Push**: The mirror uses `--force` which overwrites GitHub history. This is intentional for mirroring but be aware.

## Support

If you encounter issues:
1. Check GitLab pipeline logs for specific errors
2. Verify all variables are set correctly
3. Test SSH connection manually if needed
4. Ensure GitHub deploy key has write access
