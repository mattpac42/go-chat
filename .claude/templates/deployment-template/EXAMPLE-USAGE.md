# Example Usage - Deployment Template

This document shows a complete example of using the deployment template for a real project.

## Example Project: Task Manager App

Let's say you're deploying a task management application called "TaskFlow".

### Step 1: Gather Your Values

| Placeholder | Your Value |
|-------------|-----------|
| `{{PROJECT_NAME}}` | TaskFlow |
| `{{PROJECT_DESCRIPTION}}` | A modern task management application for teams |
| `{{APP_DIRECTORY}}` | app |
| `{{FIREBASE_PROJECT_ID}}` | taskflow-prod-12345 |
| `{{VERCEL_PROJECT_ID}}` | prj_BcD3f4G5h6I7j8K9 |
| `{{VERCEL_ORG_ID}}` | team_LmN0p1Q2r3S4t5 |
| `{{GITHUB_REPO_URL}}` | https://github.com/mycompany/taskflow |
| `{{GITHUB_USERNAME}}` | mycompany |
| `{{GITHUB_REPO_NAME}}` | taskflow |
| `{{CUSTOM_DOMAIN}}` | taskflow.mycompany.com |
| `{{VERCEL_URL}}` | taskflow.vercel.app |
| `{{FIREBASE_URL}}` | taskflow-prod-12345.web.app |

### Step 2: Copy Template to Your Project

```bash
# From the deployment-template directory
cd /path/to/deployment-template

# Copy to your project
cp -r . /path/to/taskflow-project/deployment-config/
cd /path/to/taskflow-project/deployment-config/
```

### Step 3: Replace Placeholders (Automated)

```bash
# Set environment variables
export PROJECT_NAME="TaskFlow"
export PROJECT_DESCRIPTION="A modern task management application for teams"
export APP_DIRECTORY="app"
export FIREBASE_PROJECT_ID="taskflow-prod-12345"
export VERCEL_PROJECT_ID="prj_BcD3f4G5h6I7j8K9"
export VERCEL_ORG_ID="team_LmN0p1Q2r3S4t5"
export GITHUB_REPO_URL="https://github.com/mycompany/taskflow"
export GITHUB_USERNAME="mycompany"
export GITHUB_REPO_NAME="taskflow"
export CUSTOM_DOMAIN="taskflow.mycompany.com"
export VERCEL_URL="taskflow.vercel.app"
export FIREBASE_URL="taskflow-prod-12345.web.app"

# Replace in all files (macOS)
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

# For Linux, use -i without ''
# find . -type f \( -name "*.md" -o -name "*.json" -o -name "*.yml" \) -exec sed -i ...
```

### Step 4: Move Files to Project

```bash
# From your project root
cd /path/to/taskflow-project

# Copy root-level files
cp deployment-config/vercel.json ./
cp deployment-config/DEPLOY.md ./

# Copy GitHub workflows
mkdir -p .github/workflows
cp deployment-config/.github/workflows/*.yml .github/workflows/

# Copy app config files
cp deployment-config/app-config/firebase.json app/
cp deployment-config/app-config/.firebaserc app/
```

### Step 5: Verify Replacements

Check that placeholders were replaced correctly:

```bash
# Check vercel.json
cat vercel.json
# Should show: "buildCommand": "cd app && npm install && npm run build"

# Check Firebase config
cat app/.firebaserc
# Should show: "default": "taskflow-prod-12345"

# Check GitHub workflow
cat .github/workflows/firebase-deploy.yml
# Should show: working-directory: ./app
```

### Step 6: Configure GitHub Secrets

Go to: `https://github.com/mycompany/taskflow/settings/secrets/actions`

Add these 5 secrets:
1. `FIREBASE_SERVICE_ACCOUNT` - JSON from Firebase Console
2. `FIREBASE_PROJECT_ID` - `taskflow-prod-12345`
3. `VERCEL_TOKEN` - From Vercel account settings
4. `VERCEL_ORG_ID` - `team_LmN0p1Q2r3S4t5`
5. `VERCEL_PROJECT_ID` - `prj_BcD3f4G5h6I7j8K9`

### Step 7: Test Manual Deployment

Before relying on CI/CD, test manual deployments:

```bash
# Test Firebase
cd app/
npm install
npm run build
firebase login
firebase deploy --only hosting

# Test Vercel
vercel login
vercel link
vercel --prod
```

### Step 8: Commit and Push

```bash
cd /path/to/taskflow-project

# Add deployment files
git add vercel.json DEPLOY.md .github/ app/firebase.json app/.firebaserc

# Commit
git commit -m "chore: add deployment configuration for Firebase and Vercel"

# Push to GitHub
git push origin main
```

### Step 9: Monitor Automated Deployment

1. Go to: `https://github.com/mycompany/taskflow/actions`
2. Watch both workflows run:
   - `Deploy to Firebase Hosting`
   - `Deploy to Vercel`
3. Verify deployments succeed
4. Visit deployment URLs to confirm

### Result

After successful deployment, you'll have:

- **Production URL (Vercel)**: https://taskflow.vercel.app
- **Production URL (Firebase)**: https://taskflow-prod-12345.web.app
- **Custom Domain**: https://taskflow.mycompany.com
- **Automated CI/CD**: Every push to `main` triggers deployments
- **PR Previews**: Every pull request gets preview deployments

## Troubleshooting Example Issues

### Issue: Workflow fails with "working-directory: ./{{APP_DIRECTORY}}"

**Cause**: Placeholders weren't replaced in workflow files.

**Solution**:
```bash
# Manually fix in .github/workflows/firebase-deploy.yml
sed -i '' 's|{{APP_DIRECTORY}}|app|g' .github/workflows/firebase-deploy.yml
sed -i '' 's|{{APP_DIRECTORY}}|app|g' .github/workflows/vercel-deploy.yml
git add .github/
git commit -m "fix: replace APP_DIRECTORY placeholder in workflows"
git push
```

### Issue: Firebase deploy fails with wrong project ID

**Cause**: `.firebaserc` still has placeholder.

**Solution**:
```bash
# Check current value
cat app/.firebaserc

# If shows {{FIREBASE_PROJECT_ID}}, replace it
sed -i '' 's|{{FIREBASE_PROJECT_ID}}|taskflow-prod-12345|g' app/.firebaserc
git add app/.firebaserc
git commit -m "fix: update Firebase project ID"
git push
```

## Tips for Success

1. **Double-check all placeholders**: Use `grep -r "{{" .` to find any remaining placeholders
2. **Test locally first**: Always test manual deployments before relying on CI/CD
3. **Verify secrets**: Make sure all 5 GitHub secrets are set correctly
4. **Check file paths**: Ensure `{{APP_DIRECTORY}}` matches your actual directory structure
5. **Review logs**: If CI/CD fails, check GitHub Actions logs for specific errors

## Next Steps

After successful deployment:

1. Set up environment variables for different environments
2. Configure custom domain DNS records
3. Enable monitoring and error tracking
4. Set up staging environments
5. Document any project-specific deployment steps
