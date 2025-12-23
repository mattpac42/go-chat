# Deployment Guide - {{PROJECT_NAME}}

This guide provides step-by-step instructions for deploying {{PROJECT_NAME}} to GitHub, Firebase Hosting, and Vercel.

## Quick Reference

| Resource | URL |
|----------|-----|
| **Custom Domain** | https://{{CUSTOM_DOMAIN}} |
| **Vercel (Production)** | https://{{VERCEL_URL}} |
| **Firebase Hosting** | https://{{FIREBASE_URL}} |
| **GitHub Repository** | {{GITHUB_REPO_URL}} |
| **Firebase Console** | https://console.firebase.google.com/project/{{FIREBASE_PROJECT_ID}} |
| **Vercel Dashboard** | https://vercel.com/{{VERCEL_ORG_ID}}/{{GITHUB_REPO_NAME}} |

### Project IDs & Secrets Reference

| Secret | Value |
|--------|-------|
| `FIREBASE_PROJECT_ID` | `{{FIREBASE_PROJECT_ID}}` |
| `VERCEL_PROJECT_ID` | `{{VERCEL_PROJECT_ID}}` |
| `VERCEL_ORG_ID` | `{{VERCEL_ORG_ID}}` |

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [GitHub Setup](#2-github-setup)
3. [Firebase Hosting Setup](#3-firebase-hosting-setup)
4. [Vercel Setup](#4-vercel-setup)
5. [CI/CD Automation](#5-cicd-automation)
6. [Troubleshooting](#6-troubleshooting)

---

## 1. Prerequisites

### Required Software

- **Node.js**: Version 20.x or higher (LTS recommended)
  ```bash
  node --version  # Should show v20.x.x or higher
  ```

- **npm**: Version 10.x or higher (comes with Node.js)
  ```bash
  npm --version
  ```

- **Git**: Latest version
  ```bash
  git --version
  ```

### Required Accounts

- **GitHub Account**: [Sign up at github.com](https://github.com/join)
- **Firebase Account**: [Sign up at firebase.google.com](https://firebase.google.com/)
- **Vercel Account**: [Sign up at vercel.com](https://vercel.com/signup)

### Install Required CLI Tools

#### Firebase CLI
```bash
npm install -g firebase-tools

# Verify installation
firebase --version
```

#### Vercel CLI
```bash
npm install -g vercel

# Verify installation
vercel --version
```

---

## 2. GitHub Setup

### Step 1: Create a New GitHub Repository

1. Go to [github.com/new](https://github.com/new)
2. Enter repository name: `{{GITHUB_REPO_NAME}}`
3. Choose visibility: **Public** or **Private**
4. **DO NOT** initialize with README, .gitignore, or license (we already have these)
5. Click **Create repository**

### Step 2: Push Existing Code to GitHub

From your project root directory:

```bash
# Add the remote repository
git remote add origin {{GITHUB_REPO_URL}}.git

# Verify remote was added
git remote -v

# Push to GitHub
git push -u origin main
```

### Step 3: Verify Repository

1. Visit your repository at {{GITHUB_REPO_URL}}
2. Verify all files are present
3. Ensure the `{{APP_DIRECTORY}}/` directory contains your application code

---

## 3. Firebase Hosting Setup

### Step 1: Create Firebase Project

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Click **Add project**
3. Enter project name (e.g., `{{PROJECT_NAME}}`)
4. Accept terms and click **Continue**
5. Disable Google Analytics (optional) and click **Create project**
6. Wait for project creation, then click **Continue**

### Step 2: Update Firebase Configuration

The `.firebaserc` file should be configured with your Firebase project ID:

```json
{
  "projects": {
    "default": "{{FIREBASE_PROJECT_ID}}"
  }
}
```

**Firebase Project ID**: `{{FIREBASE_PROJECT_ID}}`

### Step 3: Authenticate Firebase CLI

```bash
# Login to Firebase
firebase login

# You'll be redirected to browser to authenticate
# Follow the prompts and allow access
```

### Step 4: Verify Firebase Configuration

The project should have `firebase.json` configured:

```bash
cd {{APP_DIRECTORY}}/
cat firebase.json
```

### Step 5: Build and Deploy to Firebase

```bash
# From the {{APP_DIRECTORY}}/ directory
cd {{APP_DIRECTORY}}

# Install dependencies if not already done
npm install

# Build the application
npm run build

# Deploy to Firebase Hosting
firebase deploy --only hosting
```

### Step 6: Verify Deployment

After successful deployment, you'll see output like:

```
✔  Deploy complete!

Project Console: https://console.firebase.google.com/project/{{FIREBASE_PROJECT_ID}}/overview
Hosting URL: https://{{FIREBASE_URL}}
```

Visit the Hosting URL to verify your app is live.

### Optional: Custom Domain Setup

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Select your project
3. Navigate to **Hosting** in left sidebar
4. Click **Add custom domain**
5. Follow the wizard to add your domain
6. Update DNS records as instructed
7. Wait for SSL certificate provisioning (can take up to 24 hours)

---

## 4. Vercel Setup

### Step 1: Create Vercel Account and Connect GitHub

1. Go to [vercel.com/signup](https://vercel.com/signup)
2. Choose **Continue with GitHub**
3. Authorize Vercel to access your GitHub account

### Step 2: Import Project from GitHub

#### Option A: Via Vercel Dashboard (Recommended)

1. Click **Add New Project** on Vercel dashboard
2. Find your `{{GITHUB_REPO_NAME}}` repository in the list
3. Click **Import**
4. Configure project settings:
   - **Framework Preset**: Vite
   - **Root Directory**: `{{APP_DIRECTORY}}`
   - **Build Command**: `npm run build`
   - **Output Directory**: `build`
   - **Install Command**: `npm install`
5. Click **Deploy**

#### Option B: Via Vercel CLI

```bash
# From the {{APP_DIRECTORY}}/ directory
cd {{APP_DIRECTORY}}

# Login to Vercel
vercel login

# Deploy to Vercel
vercel

# Follow the prompts:
# - Set up and deploy? Y
# - Which scope? (select your account)
# - Link to existing project? N
# - What's your project's name? {{GITHUB_REPO_NAME}}
# - In which directory is your code located? ./
# - Want to override settings? N

# Deploy to production
vercel --prod
```

### Step 3: Configure Environment Variables (if needed)

If your app uses environment variables:

1. Go to your project in Vercel Dashboard
2. Navigate to **Settings** > **Environment Variables**
3. Add each variable with appropriate scope (Production, Preview, Development)
4. Redeploy if needed

### Step 4: Verify Deployment

After deployment, Vercel will provide URLs:

```
Preview: https://{{GITHUB_REPO_NAME}}-hash.vercel.app
Production: https://{{VERCEL_URL}}
```

Visit these URLs to verify deployment.

### Custom Domain Setup

**Important**: Vercel must be configured to accept traffic for your custom domain. Without this step, your DNS configuration alone won't work.

#### Step 1: Add Domain to Vercel Project

1. Go to [Vercel Dashboard](https://vercel.com/{{VERCEL_ORG_ID}}/{{GITHUB_REPO_NAME}}/settings/domains)
2. Navigate to **Settings** > **Domains**
3. Click **Add Domain**
4. Enter your domain name (e.g., `{{CUSTOM_DOMAIN}}`)
5. Vercel will then accept traffic for that hostname

#### Step 2: Configure DNS

After adding the domain in Vercel, configure your DNS:

**For subdomain (e.g., app.example.com)**:
```
Type: CNAME
Name: app
Value: cname.vercel-dns.com
```

**For root domain (e.g., example.com)**:
```
Type: A
Name: @
Value: 76.76.21.21
```

#### Step 3: Verify and Wait for SSL

1. Vercel will automatically provision an SSL certificate
2. This can take a few minutes to 24 hours
3. Check status in Vercel Dashboard > Domains

**Current Custom Domain**: `{{CUSTOM_DOMAIN}}`

---

## 5. CI/CD Automation

The project includes GitHub Actions workflows for automated deployments.

### Required Repository Secrets

Before automated deployments work, configure these secrets in your GitHub repository:

#### Step 1: Navigate to Repository Settings

1. Go to your GitHub repository: {{GITHUB_REPO_URL}}
2. Click **Settings** tab
3. Navigate to **Secrets and variables** > **Actions**
4. Click **New repository secret**

#### Step 2: Add Firebase Secrets

**FIREBASE_SERVICE_ACCOUNT**

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Select your project: `{{FIREBASE_PROJECT_ID}}`
3. Click gear icon > **Project settings**
4. Navigate to **Service accounts** tab
5. Click **Generate new private key**
6. Download the JSON file
7. Copy entire contents of the JSON file
8. In GitHub, create secret named `FIREBASE_SERVICE_ACCOUNT`
9. Paste the JSON contents as the value

**FIREBASE_PROJECT_ID**

1. In Firebase Console, copy your Project ID: `{{FIREBASE_PROJECT_ID}}`
2. In GitHub, create secret named `FIREBASE_PROJECT_ID`
3. Paste the Project ID as the value

#### Step 3: Add Vercel Secrets

**VERCEL_TOKEN**

1. Go to [Vercel Account Settings](https://vercel.com/account/tokens)
2. Click **Create Token**
3. Name it `GitHub Actions` with appropriate scope
4. Copy the token
5. In GitHub, create secret named `VERCEL_TOKEN`
6. Paste the token as the value

**VERCEL_ORG_ID**

1. Install Vercel CLI: `npm i -g vercel`
2. Run `vercel login` if not already logged in
3. In your project directory, run:
   ```bash
   cd {{APP_DIRECTORY}}/
   vercel link
   ```
4. Open `.vercel/project.json`
5. Copy the `orgId` value: `{{VERCEL_ORG_ID}}`
6. In GitHub, create secret named `VERCEL_ORG_ID`
7. Paste the orgId as the value

**VERCEL_PROJECT_ID**

1. From the same `.vercel/project.json` file
2. Copy the `projectId` value: `{{VERCEL_PROJECT_ID}}`
3. In GitHub, create secret named `VERCEL_PROJECT_ID`
4. Paste the projectId as the value

### How Automated Deployments Work

#### Firebase Deployment Workflow

**Trigger**: Push to `main` branch or pull request

**Process**:
1. Checks out code
2. Installs dependencies
3. Builds the application
4. Deploys to Firebase Hosting
5. Creates preview deployments for PRs

**File**: `.github/workflows/firebase-deploy.yml`

#### Vercel Deployment Workflow

**Trigger**: Push to `main` branch or pull request

**Process**:
1. Checks out code
2. Deploys to Vercel using CLI
3. Creates preview deployments for pull requests
4. Promotes to production on main branch
5. Comments preview URL on PRs

**File**: `.github/workflows/vercel-deploy.yml`

### Preview Deployments

Both workflows support preview deployments for pull requests:

1. **Create a pull request** against `main` branch
2. Workflows automatically deploy preview versions
3. Preview URLs are commented on the PR
4. Each commit to the PR updates the preview
5. Merging to `main` triggers production deployment

### Monitoring Deployments

#### Via GitHub Actions

1. Go to your repository: {{GITHUB_REPO_URL}}
2. Click **Actions** tab
3. View workflow runs and their status
4. Click on a run to see detailed logs

#### Via Hosting Platforms

- **Firebase**: [Firebase Console](https://console.firebase.google.com/project/{{FIREBASE_PROJECT_ID}}) > Hosting > View History
- **Vercel**: [Vercel Dashboard](https://vercel.com/{{VERCEL_ORG_ID}}/{{GITHUB_REPO_NAME}}) > Deployments

---

## 6. Troubleshooting

### Common Issues and Solutions

#### Vercel Build Fails - "Could not resolve entry module index.html"

**Error**: `Could not resolve entry module "index.html"`

**Cause**: The `index.html` file is not tracked in git (likely ignored by `.gitignore`).

**Solution**:
```bash
# Check if index.html is being ignored
git check-ignore {{APP_DIRECTORY}}/index.html

# If ignored, update .gitignore to allow it
# Add this line to .gitignore:
!{{APP_DIRECTORY}}/index.html

# Then add and commit
git add -f {{APP_DIRECTORY}}/index.html
git commit -m "fix: add index.html for Vite build"
git push
```

#### Build Fails - Node Version Mismatch

**Error**: `The engine "node" is incompatible with this module`

**Solution**:
```bash
# Update Node.js to version 20 or higher
nvm install 20
nvm use 20

# Or download from nodejs.org
```

#### Firebase Deploy Fails - Not Logged In

**Error**: `Error: Failed to authenticate`

**Solution**:
```bash
firebase login --reauth
```

#### Firebase Deploy Fails - Wrong Project

**Error**: `Error: HTTP Error: 404, Project not found`

**Solution**:
```bash
# Verify project ID in .firebaserc matches your Firebase project
cat {{APP_DIRECTORY}}/.firebaserc

# Or reinitialize
cd {{APP_DIRECTORY}}/
firebase use --add
```

#### Vercel Deploy Fails - Build Error

**Error**: `Error: Build failed`

**Solution**:
```bash
# Test build locally first
cd {{APP_DIRECTORY}}/
npm run build

# Check for errors and fix them
# Then deploy again
```

#### Vercel Deploy Fails - Wrong Root Directory

**Error**: `No package.json found`

**Solution**:
1. In Vercel Dashboard, go to **Settings** > **General**
2. Set **Root Directory** to `{{APP_DIRECTORY}}`
3. Redeploy

#### GitHub Actions Fails - Missing Secrets

**Error**: `Error: Input required and not supplied: VERCEL_TOKEN`

**Solution**:
1. Verify all required secrets are added (see section 5.2)
2. Check secret names match exactly (case-sensitive)
3. Re-run the workflow

#### Build Output Directory Not Found

**Error**: `Error: Could not find build directory`

**Solution**:
```bash
# Verify vite.config.ts has correct outDir
cat {{APP_DIRECTORY}}/vite.config.ts | grep outDir

# Should show: outDir: 'build'
```

### Getting Help

- **Firebase**: [Firebase Support](https://firebase.google.com/support)
- **Vercel**: [Vercel Support](https://vercel.com/support)
- **GitHub Actions**: [GitHub Actions Documentation](https://docs.github.com/en/actions)

### Useful Commands

```bash
# Check build output
cd {{APP_DIRECTORY}}/
npm run build
ls -la build/

# Preview build locally
npm run preview

# Check Firebase project
firebase projects:list

# Check Vercel project
vercel ls

# View deployment logs
firebase hosting:channel:list
vercel logs
```

---

## Next Steps

After successful deployment:

1. Set up custom domains (optional)
2. Configure environment variables for production
3. Enable monitoring and analytics
4. Set up error tracking (e.g., Sentry)
5. Configure CDN and caching rules
6. Set up SSL certificates (handled automatically by both platforms)

## Additional Resources

- [Vite Documentation](https://vite.dev/)
- [Firebase Hosting Documentation](https://firebase.google.com/docs/hosting)
- [Vercel Documentation](https://vercel.com/docs)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
