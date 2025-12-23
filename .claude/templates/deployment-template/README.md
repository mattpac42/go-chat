# Deployment Template - Setup Guide

This template provides a complete deployment configuration for GitHub, Firebase Hosting, and Vercel.

## Quick Start

1. Copy this entire `deployment-template/` folder to your project root
2. Replace all placeholders marked with `{{PLACEHOLDER}}` syntax
3. Follow the setup instructions below for each platform
4. Configure GitHub secrets for CI/CD automation

## Placeholder Reference Table

Replace these placeholders throughout all template files:

| Placeholder | Description | Example |
|-------------|-------------|---------|
| `{{PROJECT_NAME}}` | Display name for your project | "My App" or "AAC Communication Board" |
| `{{PROJECT_DESCRIPTION}}` | Brief description of your app | "Web application for augmentative communication" |
| `{{APP_DIRECTORY}}` | Application subdirectory path | "application" or "src" or "." (for root) |
| `{{FIREBASE_PROJECT_ID}}` | Firebase project ID | "my-app-12345" |
| `{{VERCEL_PROJECT_ID}}` | Vercel project ID | "prj_xxxxx" |
| `{{VERCEL_ORG_ID}}` | Vercel organization/team ID | "team_xxxxx" |
| `{{GITHUB_REPO_URL}}` | Full GitHub repository URL | "https://github.com/username/repo" |
| `{{GITHUB_USERNAME}}` | Your GitHub username | "username" |
| `{{GITHUB_REPO_NAME}}` | GitHub repository name | "repo" |
| `{{CUSTOM_DOMAIN}}` | Custom domain (optional) | "app.example.com" |
| `{{VERCEL_URL}}` | Vercel deployment URL | "my-app.vercel.app" |
| `{{FIREBASE_URL}}` | Firebase Hosting URL | "my-app.web.app" |

## File Setup Instructions

### Step 1: Copy Template Files to Project

```bash
# From your deployment-template/ directory, copy files to project root:

# Copy root-level files
cp vercel.json /path/to/your/project/
cp DEPLOY.md /path/to/your/project/

# Copy GitHub workflows
cp -r .github/ /path/to/your/project/

# Copy app config files to your application directory
# Replace {{APP_DIRECTORY}} with your actual app directory name
cp app-config/firebase.json /path/to/your/project/{{APP_DIRECTORY}}/
cp app-config/.firebaserc /path/to/your/project/{{APP_DIRECTORY}}/
```

### Step 2: Replace Placeholders

**Option A: Manual Replacement**

Open each file and replace `{{PLACEHOLDER}}` values with your actual values.

**Option B: Automated Replacement (Linux/macOS)**

```bash
# Set your values
export PROJECT_NAME="My App"
export PROJECT_DESCRIPTION="Web application for..."
export APP_DIRECTORY="application"
export FIREBASE_PROJECT_ID="my-app-12345"
export VERCEL_PROJECT_ID="prj_xxxxx"
export VERCEL_ORG_ID="team_xxxxx"
export GITHUB_REPO_URL="https://github.com/username/repo"
export GITHUB_USERNAME="username"
export GITHUB_REPO_NAME="repo"
export CUSTOM_DOMAIN="app.example.com"
export VERCEL_URL="my-app.vercel.app"
export FIREBASE_URL="my-app.web.app"

# Replace in all files
find . -type f -exec sed -i '' \
  -e "s|{{PROJECT_NAME}}|$PROJECT_NAME|g" \
  -e "s|{{PROJECT_DESCRIPTION}}|$PROJECT_DESCRIPTION|g" \
  -e "s|{{APP_DIRECTORY}}|$APP_DIRECTORY|g" \
  -e "s|{{FIREBASE_PROJECT_ID}}|$FIREBASE_PROJECT_ID|g" \
  -e "s|{{VERCEL_PROJECT_ID}}|$VERCEL_PROJECT_ID|g" \
  -e "s|{{VERCEL_ORG_ID}}|$VERCEL_ORG_ID|g" \
  -e "s|{{GITHUB_REPO_URL}}|$GITHUB_REPO_URL|g" \
  -e "s|{{GITHUB_USERNAME}}|$GITHUB_USERNAME|g" \
  -e "s|{{GITHUB_REPO_NAME}}|$GITHUB_REPO_NAME|g" \
  -e "s|{{CUSTOM_DOMAIN}}|$CUSTOM_DOMAIN|g" \
  -e "s|{{VERCEL_URL}}|$VERCEL_URL|g" \
  -e "s|{{FIREBASE_URL}}|$FIREBASE_URL|g" \
  {} \;
```

### Step 3: Configure GitHub Secrets

Before CI/CD automation works, add these secrets to your GitHub repository:

**Navigate to**: Repository → Settings → Secrets and variables → Actions → New repository secret

#### Required Secrets

1. **FIREBASE_SERVICE_ACCOUNT**
   - Go to [Firebase Console](https://console.firebase.google.com/)
   - Project Settings → Service accounts
   - Generate new private key
   - Copy entire JSON contents as secret value

2. **FIREBASE_PROJECT_ID**
   - Your Firebase project ID (e.g., `my-app-12345`)

3. **VERCEL_TOKEN**
   - Go to [Vercel Account Settings](https://vercel.com/account/tokens)
   - Create new token
   - Copy token as secret value

4. **VERCEL_ORG_ID**
   - Run `vercel link` in your app directory
   - Copy `orgId` from `.vercel/project.json`

5. **VERCEL_PROJECT_ID**
   - Copy `projectId` from `.vercel/project.json`

### Step 4: Platform Setup

#### Firebase Setup

```bash
# Install Firebase CLI
npm install -g firebase-tools

# Login to Firebase
firebase login

# Navigate to your app directory
cd {{APP_DIRECTORY}}

# Deploy
firebase deploy --only hosting
```

#### Vercel Setup

```bash
# Install Vercel CLI
npm install -g vercel

# Login to Vercel
vercel login

# Navigate to your app directory
cd {{APP_DIRECTORY}}

# Link project
vercel link

# Deploy to production
vercel --prod
```

## Template Files Included

```
deployment-template/
├── README.md                           # This file - setup instructions
├── DEPLOY.md                           # Full deployment guide (templated)
├── vercel.json                         # Root-level Vercel configuration
├── .github/
│   └── workflows/
│       ├── firebase-deploy.yml         # Firebase CI/CD workflow
│       └── vercel-deploy.yml           # Vercel CI/CD workflow
└── app-config/                         # Files to copy to {{APP_DIRECTORY}}/
    ├── firebase.json                   # Firebase hosting configuration
    └── .firebaserc                     # Firebase project reference
```

## Validation Checklist

After setup, verify:

- [ ] All `{{PLACEHOLDER}}` values replaced in all files
- [ ] `vercel.json` in project root
- [ ] `DEPLOY.md` in project root
- [ ] `.github/workflows/` directory with both workflow files
- [ ] `firebase.json` and `.firebaserc` in `{{APP_DIRECTORY}}/`
- [ ] All 5 GitHub secrets configured
- [ ] Firebase CLI authenticated
- [ ] Vercel CLI authenticated and project linked
- [ ] Test deployments work manually before relying on CI/CD

## Troubleshooting

**Problem**: Placeholder values still showing after replacement

**Solution**: Ensure you've replaced ALL occurrences in ALL files, including DEPLOY.md

---

**Problem**: GitHub Actions failing with "secret not found"

**Solution**: Verify secret names match exactly (case-sensitive) and are added to the correct repository

---

**Problem**: Build failing with "directory not found"

**Solution**: Verify `{{APP_DIRECTORY}}` matches your actual application directory name

---

**Problem**: Firebase deploy fails with "project not found"

**Solution**: Ensure `.firebaserc` has correct project ID and you've run `firebase login`

## Getting Help

- **Firebase**: [Firebase Hosting Docs](https://firebase.google.com/docs/hosting)
- **Vercel**: [Vercel Documentation](https://vercel.com/docs)
- **GitHub Actions**: [GitHub Actions Docs](https://docs.github.com/en/actions)

## Next Steps

After successful setup:

1. Review and customize `DEPLOY.md` for your project
2. Test manual deployments first
3. Make a test commit to trigger automated deployments
4. Set up custom domains (optional)
5. Configure environment variables if needed
6. Enable monitoring and analytics
