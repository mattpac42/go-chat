# Quick Start Guide - Secure Pipeline Template

This guide provides the minimal steps to get the secure pipeline running in a new project.

## Prerequisites

- GitLab project with container registry enabled
- GitLab Runner with `docker` tag and privileged mode enabled
- Deployment target server with Docker and Docker Compose installed
- SSH access to deployment target

## 5-Minute Setup

### Step 1: Copy Template Files (1 min)

```bash
# From your project root
cp .claude/templates/secure-pipeline/gitlab-ci.yml.template .gitlab-ci.yml
cp .claude/templates/secure-pipeline/Dockerfile.template Dockerfile
cp .claude/templates/secure-pipeline/Dockerfile.dind.template Dockerfile.dind
cp .claude/templates/secure-pipeline/docker-compose.yml.template docker-compose.yml
```

### Step 2: Create VERSION File (10 sec)

```bash
echo "1.0.0" > VERSION
```

### Step 3: Find and Replace Variables (2 min)

Use your editor's find and replace to update these placeholders in ALL copied files:

| Find | Replace With | Example |
|------|--------------|---------|
| `{{APP_NAME}}` | Your app name | `myapp` |
| `{{APP_DOMAIN}}` | Your domain | `myapp.example.com` |
| `{{DEPLOY_PORT}}` | Your port | `4000` |
| `{{CI_REGISTRY}}` | Your registry | `gitlab.example.com:5050` |
| `{{CI_REGISTRY_IMAGE}}` | Full registry path | `gitlab.example.com:5050/group/project` |
| `{{VERSION_FILE}}` | Version file | `VERSION` |
| `{{APP_USER}}` | App user | `appuser` |
| `{{FLASK_APP}}` | Flask app file | `app.py` |
| `{{WSGI_MODULE}}` | WSGI module | `wsgi:application` |
| `{{WORKERS}}` | Worker count | `1` |
| `{{DEPLOY_HOST}}` | Deploy IP | `192.168.0.100` |
| `{{DEPLOY_USER}}` | Deploy user | `deployer` |
| `{{DEPLOY_PATH}}` | Deploy path | `/srv/myapp` |
| `{{DATA_PATH}}` | Data path | `/srv/myapp/data` |

**Quick Bash Script**:
```bash
# Replace all in one command (edit values first!)
find . -name "*.yml" -o -name "Dockerfile*" -o -name "docker-compose.yml" | xargs sed -i '' \
  -e 's/{{APP_NAME}}/myapp/g' \
  -e 's/{{APP_DOMAIN}}/myapp.example.com/g' \
  -e 's/{{DEPLOY_PORT}}/4000/g' \
  -e 's|{{CI_REGISTRY}}|gitlab.example.com:5050|g' \
  -e 's|{{CI_REGISTRY_IMAGE}}|gitlab.example.com:5050/group/project|g' \
  -e 's/{{VERSION_FILE}}/VERSION/g' \
  -e 's/{{APP_USER}}/appuser/g' \
  -e 's/{{FLASK_APP}}/app.py/g' \
  -e 's/{{WSGI_MODULE}}/wsgi:application/g' \
  -e 's/{{WORKERS}}/1/g' \
  -e 's/{{DEPLOY_HOST}}/192.168.0.100/g' \
  -e 's/{{DEPLOY_USER}}/deployer/g' \
  -e 's|{{DEPLOY_PATH}}|/srv/myapp|g' \
  -e 's|{{DATA_PATH}}|/srv/myapp/data|g'
```

### Step 4: Configure GitLab Variables (1 min)

In GitLab: **Settings > CI/CD > Variables**

Add this **required** variable:

1. **DEPLOY_SSH_KEY** (Type: File)
   - Generate: `ssh-keygen -t ed25519 -C "gitlab-deploy"`
   - Copy entire private key content
   - Add public key to deployment server's `~/.ssh/authorized_keys`

### Step 5: Commit and Push (30 sec)

```bash
git add .gitlab-ci.yml Dockerfile Dockerfile.dind docker-compose.yml VERSION
git commit -m "Add secure CI/CD pipeline"
git push origin main
```

Pipeline will run automatically!

## What Happens Next?

1. **Prepare Stage**: Builds custom Docker-in-Docker image (first run only)
2. **Lint Stage**: Checks code quality
3. **Test Stage**: Runs tests
4. **Build Stage**: Creates container image with multiple tags
5. **Scan Stage**: Scans for vulnerabilities
6. **Release Stage**: Validates security checks
7. **Deploy Stage**: Deploys to your server

## First Run Expected Behavior

- Build stage: ~3-5 minutes (downloads base images)
- Scan stage: ~2-3 minutes (downloads vulnerability database)
- Deploy stage: ~1-2 minutes

Total first run: ~10-15 minutes

Subsequent runs: ~5-7 minutes (cached images)

## Verify Deployment

After pipeline completes:

```bash
# On deployment server
docker ps | grep myapp

# Test locally
curl http://localhost:4000/

# Test via domain
curl https://myapp.example.com/
```

## Fast Development Mode

For quick testing during development:

```yaml
# In .gitlab-ci.yml
variables:
  QUICK_DEPLOY: "true"
```

This skips lint, test, and scan stages (deploys in ~2 minutes).

**WARNING**: Only for development! Set to `"false"` for production.

## Troubleshooting

### Pipeline Fails at Build Stage

**Check**: GitLab Runner configuration
```bash
# On runner server
sudo gitlab-runner verify
sudo gitlab-runner list
```

**Fix**: Ensure runner has privileged mode enabled

### Pipeline Fails at Deploy Stage

**Check**: SSH key and deployment server access
```bash
# Test SSH connection
ssh -i deploy_key deployer@192.168.0.100

# Check Docker on deployment server
docker --version
docker-compose --version
```

**Fix**: Verify DEPLOY_SSH_KEY variable is set correctly

### Health Check Fails

**Check**: Application logs
```bash
# On deployment server
docker logs myapp-app
```

**Fix**: Ensure your application:
- Listens on correct port
- Has `/healthz` or `/` endpoint
- Returns HTTP 200 when healthy

## Need Custom CA Certificate?

If your GitLab uses self-signed certificates:

```bash
# Encode certificate
base64 -w 0 /path/to/ca.crt

# Add to GitLab Variables
# Name: LOCAL_CA_CERT
# Type: File
# Value: <base64 output>

# Rebuild DinD image (manually trigger job)
```

## Next Steps

1. Review `README.md` for detailed documentation
2. Check `variables.env.example` for all configuration options
3. Review `.gitlab-ci.yml.backup` for advanced features (SBOM, signing)
4. Customize health checks for your application
5. Set up monitoring and alerting

## Common Customizations

### Change Python Version

```dockerfile
# In Dockerfile
FROM python:3.11-slim as builder
# ...
FROM python:3.11-slim
```

### Add Build Arguments

```yaml
# In .gitlab-ci.yml build_image job
script:
  - docker build --build-arg VERSION=$VERSION -t temp-build:${CI_PIPELINE_ID} .
```

### Change Deployment Method

Replace deploy stage with your preferred method (Kubernetes, Nomad, etc.)

### Add More Test Stages

```yaml
# In .gitlab-ci.yml
integration_test:
  stage: test
  script:
    - pytest tests/integration/
```

## Support

- Full documentation: `README.md`
- Variable reference: `variables.env.example`
- Advanced features: `.gitlab-ci.yml.backup`
- GitLab CI/CD docs: https://docs.gitlab.com/ee/ci/
