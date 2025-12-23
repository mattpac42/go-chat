# Quick Start Checklist

Use this checklist to quickly deploy your project using this template.

## Prerequisites Checklist

- [ ] Node.js 20+ installed (`node --version`)
- [ ] Git installed (`git --version`)
- [ ] GitHub account created
- [ ] Firebase account created
- [ ] Vercel account created
- [ ] Firebase CLI installed (`npm i -g firebase-tools`)
- [ ] Vercel CLI installed (`npm i -g vercel`)

## Setup Checklist

### 1. Copy Template Files

- [ ] Copy entire `deployment-template/` folder to your project
- [ ] Navigate to the copied folder

### 2. Gather Required Information

Fill in your values:

```bash
# Required values
PROJECT_NAME="[Your Project Name]"
PROJECT_DESCRIPTION="[Brief description]"
APP_DIRECTORY="[app directory name]"              # e.g., "application", "app", "src", or "."
FIREBASE_PROJECT_ID="[firebase-project-id]"       # e.g., "my-app-12345"
GITHUB_REPO_URL="[https://github.com/user/repo]"
GITHUB_USERNAME="[your-username]"
GITHUB_REPO_NAME="[repo-name]"

# Get these after creating Firebase/Vercel projects
VERCEL_PROJECT_ID="[prj_xxxxx]"                   # From .vercel/project.json
VERCEL_ORG_ID="[team_xxxxx]"                      # From .vercel/project.json

# Optional (can be set later)
CUSTOM_DOMAIN="[app.example.com]"                 # Optional
VERCEL_URL="[your-app.vercel.app]"               # After Vercel setup
FIREBASE_URL="[your-app.web.app]"                 # After Firebase setup
```

### 3. Create Firebase Project

- [ ] Go to [Firebase Console](https://console.firebase.google.com/)
- [ ] Click "Add project"
- [ ] Enter project name
- [ ] Complete setup wizard
- [ ] Copy Firebase Project ID
- [ ] Go to Project Settings → Service Accounts
- [ ] Generate new private key (save JSON file for later)

### 4. Create Vercel Project (Initial Link)

```bash
cd [your-app-directory]
vercel login
vercel link
# Copy orgId and projectId from .vercel/project.json
```

- [ ] Vercel project linked
- [ ] Copy `orgId` value → `VERCEL_ORG_ID`
- [ ] Copy `projectId` value → `VERCEL_PROJECT_ID`

### 5. Replace Placeholders

Choose **Option A** (Automated) or **Option B** (Manual):

#### Option A: Automated Replacement (Recommended)

```bash
# Set all environment variables (from step 2)
export PROJECT_NAME="Your Project Name"
export PROJECT_DESCRIPTION="Your description"
export APP_DIRECTORY="app"
# ... set all other variables ...

# Run replacement (macOS)
find . -type f \( -name "*.md" -o -name "*.json" -o -name "*.yml" \) -exec sed -i '' \
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

#### Option B: Manual Replacement

- [ ] Open each file and replace `{{PLACEHOLDER}}` values
- [ ] Check: `vercel.json`
- [ ] Check: `DEPLOY.md`
- [ ] Check: `.github/workflows/firebase-deploy.yml`
- [ ] Check: `.github/workflows/vercel-deploy.yml`
- [ ] Check: `app-config/firebase.json`
- [ ] Check: `app-config/.firebaserc`

### 6. Verify Replacements

```bash
# Check for any remaining placeholders
grep -r "{{" . --include="*.md" --include="*.json" --include="*.yml"

# Should only show {{PLACEHOLDER}} in README.md examples
```

- [ ] No placeholders remain (except in documentation examples)

### 7. Move Files to Project

```bash
# From your project root
cd /path/to/your/project

# Copy files
cp deployment-config/vercel.json ./
cp deployment-config/DEPLOY.md ./
mkdir -p .github/workflows
cp deployment-config/.github/workflows/*.yml .github/workflows/
cp deployment-config/app-config/firebase.json [APP_DIRECTORY]/
cp deployment-config/app-config/.firebaserc [APP_DIRECTORY]/
```

- [ ] `vercel.json` copied to project root
- [ ] `DEPLOY.md` copied to project root
- [ ] GitHub workflows copied to `.github/workflows/`
- [ ] `firebase.json` copied to app directory
- [ ] `.firebaserc` copied to app directory

### 8. Configure GitHub Secrets

Go to: `https://github.com/[USERNAME]/[REPO]/settings/secrets/actions`

Add these 5 secrets:

- [ ] `FIREBASE_SERVICE_ACCOUNT` (paste entire JSON from Firebase service account)
- [ ] `FIREBASE_PROJECT_ID` (your Firebase project ID)
- [ ] `VERCEL_TOKEN` (from [Vercel Account Settings](https://vercel.com/account/tokens))
- [ ] `VERCEL_ORG_ID` (from `.vercel/project.json`)
- [ ] `VERCEL_PROJECT_ID` (from `.vercel/project.json`)

### 9. Test Manual Deployments

#### Firebase

```bash
cd [APP_DIRECTORY]
npm install
npm run build
firebase login
firebase deploy --only hosting
```

- [ ] Firebase login successful
- [ ] Build completes without errors
- [ ] Firebase deploy successful
- [ ] Visit Firebase URL and verify app works

#### Vercel

```bash
cd [APP_DIRECTORY]
vercel --prod
```

- [ ] Vercel deploy successful
- [ ] Visit Vercel URL and verify app works

### 10. Commit and Push

```bash
cd /path/to/your/project

git add vercel.json DEPLOY.md .github/ [APP_DIRECTORY]/firebase.json [APP_DIRECTORY]/.firebaserc
git commit -m "chore: add deployment configuration for Firebase and Vercel"
git push origin main
```

- [ ] Files committed
- [ ] Pushed to GitHub

### 11. Verify Automated Deployments

Go to: `https://github.com/[USERNAME]/[REPO]/actions`

- [ ] "Deploy to Firebase Hosting" workflow runs successfully
- [ ] "Deploy to Vercel" workflow runs successfully
- [ ] Both deployments complete without errors

### 12. Configure Custom Domain (Optional)

#### Vercel

1. Go to Vercel Dashboard → Your Project → Settings → Domains
2. Add domain: `[CUSTOM_DOMAIN]`
3. Configure DNS as instructed
4. Wait for SSL certificate

- [ ] Domain added to Vercel
- [ ] DNS configured
- [ ] SSL certificate active

#### Firebase (if using)

1. Go to Firebase Console → Hosting
2. Add custom domain
3. Configure DNS as instructed
4. Wait for SSL certificate

- [ ] Domain added to Firebase
- [ ] DNS configured
- [ ] SSL certificate active

## Final Verification

- [ ] All manual deployments work
- [ ] All automated deployments work
- [ ] All URLs accessible:
  - [ ] Vercel production URL
  - [ ] Firebase production URL
  - [ ] Custom domain (if configured)
- [ ] Pull request preview deployments work
- [ ] No errors in GitHub Actions logs

## Troubleshooting

If something doesn't work:

1. Check GitHub Actions logs for specific errors
2. Verify all secrets are set correctly (case-sensitive)
3. Ensure placeholders are replaced in all files
4. Confirm `APP_DIRECTORY` matches your actual directory
5. Test builds locally: `npm run build`
6. Review `DEPLOY.md` troubleshooting section

## Success!

Once all checklist items are complete, your project is fully deployed with:

- Automated CI/CD pipelines
- Preview deployments for pull requests
- Production deployments on every push to main
- Dual hosting on Firebase and Vercel
- Optional custom domain support

## Next Steps

- [ ] Set up environment variables for different environments
- [ ] Enable monitoring and error tracking
- [ ] Configure caching strategies
- [ ] Set up staging environments
- [ ] Document any project-specific deployment steps
