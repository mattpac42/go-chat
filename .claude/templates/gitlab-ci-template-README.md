# GitLab CI/CD Template - Usage Guide

## Overview

This template provides a production-ready GitLab CI/CD pipeline for Python/Flask applications with comprehensive linting, security scanning, testing, Docker builds, and automated deployment.

## Quick Start

### 1. Copy and Customize Template

```bash
# Copy template to your project root
cp .claude/templates/gitlab-ci-template.yml .gitlab-ci.yml

# Replace placeholders (use your IDE's Find & Replace)
# Find: {{APP_NAME}}          Replace with: my-flask-app
# Find: {{APP_DOMAIN}}        Replace with: app.example.com
# Find: {{DEPLOY_PORT}}       Replace with: 5000
# Find: {{CI_REGISTRY_IMAGE}} Replace with: gitlab.example.com:5050/group/project
# Find: {{PYTHON_VERSION}}    Replace with: 3.12
```

### 2. Create Required Files

```bash
# Create version file
echo "1.0.0" > VERSION

# Create docker-compose.yml (see example below)
touch docker-compose.yml

# Create Dockerfile (see example below)
touch Dockerfile

# Optional: Create Dockerfile.dind for home lab CA support
touch Dockerfile.dind
```

### 3. Configure GitLab Variables

Navigate to your GitLab project: **Settings > CI/CD > Variables**

**Required Variables:**

| Variable | Type | Value | Protected | Masked |
|----------|------|-------|-----------|--------|
| `DEPLOY_SSH_KEY` | File | SSH private key contents | ✅ | ❌ |
| `DEPLOY_HOST` | Variable | `192.168.0.200` | ✅ | ❌ |
| `DEPLOY_USER` | Variable | `deployer` | ✅ | ❌ |
| `DEPLOY_PATH` | Variable | `/opt/my-flask-app` | ✅ | ❌ |

**Optional Variables (Home Lab CA):**

| Variable | Type | Value | Protected | Masked |
|----------|------|-------|-----------|--------|
| `LOCAL_CA_CERT` | Variable | Base64-encoded CA certificate | ✅ | ❌ |

### 4. Push and Deploy

```bash
git add .gitlab-ci.yml VERSION docker-compose.yml Dockerfile
git commit -m "Add CI/CD pipeline configuration"
git push origin main
```

## Detailed Setup Instructions

### Docker Configuration

**Dockerfile Example:**

```dockerfile
FROM python:3.12-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Expose application port
EXPOSE 5000

# Run application
CMD ["python", "app.py"]
```

**docker-compose.yml Example:**

```yaml
version: '3.8'

services:
  my-flask-app:
    image: ${DOCKER_IMAGE}
    container_name: my-flask-app
    restart: unless-stopped
    ports:
      - "5000:5000"
    environment:
      - FLASK_ENV=production
      - DATABASE_URL=${DATABASE_URL}
    labels:
      # Traefik reverse proxy configuration (optional)
      - "traefik.enable=true"
      - "traefik.http.routers.my-flask-app.rule=Host(`app.example.com`)"
      - "traefik.http.routers.my-flask-app.entrypoints=websecure"
      - "traefik.http.routers.my-flask-app.tls.certresolver=letsencrypt"
```

### Home Lab CA Certificate Setup (Optional)

If you're using a self-signed certificate for your GitLab registry (common in home labs):

**1. Export your root CA certificate:**

```bash
# macOS
base64 -i /path/to/root-ca.crt | pbcopy

# Linux
base64 -w 0 /path/to/root-ca.crt
```

**2. Add to GitLab Variables:**
- Name: `LOCAL_CA_CERT`
- Type: Variable
- Value: Paste the base64 output
- Protected: Yes
- Masked: No (certificate is too long)

**3. Create Dockerfile.dind:**

```dockerfile
FROM docker:27-dind

ARG LOCAL_CA_CERT

RUN apk add --no-cache ca-certificates && \
    if [ -n "$LOCAL_CA_CERT" ]; then \
        echo "$LOCAL_CA_CERT" | base64 -d > /usr/local/share/ca-certificates/caddy-root-ca.crt && \
        update-ca-certificates; \
    fi

# Configure Docker registry trust
RUN mkdir -p /etc/docker/certs.d/gitlab.yuki.lan:5050 && \
    if [ -n "$LOCAL_CA_CERT" ]; then \
        echo "$LOCAL_CA_CERT" | base64 -d > /etc/docker/certs.d/gitlab.yuki.lan:5050/ca.crt; \
    fi
```

### SSH Deployment Setup

**1. Generate SSH Key Pair:**

```bash
ssh-keygen -t ed25519 -C "gitlab-ci-deploy" -f ~/.ssh/gitlab-ci-deploy
```

**2. Add Public Key to Target Host:**

```bash
# Copy public key to deployment server
ssh-copy-id -i ~/.ssh/gitlab-ci-deploy.pub deployer@192.168.0.200

# Verify access
ssh -i ~/.ssh/gitlab-ci-deploy deployer@192.168.0.200
```

**3. Add Private Key to GitLab:**

- Navigate to: **Settings > CI/CD > Variables**
- Click: **Add Variable**
- Key: `DEPLOY_SSH_KEY`
- Type: **File**
- Value: Paste entire contents of `~/.ssh/gitlab-ci-deploy` (private key)
- Protected: ✅ Yes
- Masked: ❌ No (too long)

**4. Create Deployment Directory on Target:**

```bash
ssh deployer@192.168.0.200
sudo mkdir -p /opt/my-flask-app
sudo chown deployer:deployer /opt/my-flask-app
```

### GitLab Runner Configuration

**Ensure your GitLab Runner has the 'docker' tag and privileged mode:**

```bash
# Register runner with docker executor
sudo gitlab-runner register \
  --tag-list docker \
  --executor docker \
  --docker-image docker:27 \
  --docker-privileged

# Or modify existing runner config
sudo nano /etc/gitlab-runner/config.toml
```

**config.toml Example:**

```toml
[[runners]]
  name = "docker-runner"
  url = "https://gitlab.yuki.lan/"
  token = "YOUR_RUNNER_TOKEN"
  executor = "docker"
  [runners.docker]
    image = "docker:27"
    privileged = true
    volumes = ["/var/run/docker.sock:/var/run/docker.sock", "/cache"]
```

## Pipeline Stages Explained

### 1. Prepare Stage
- **Purpose**: Builds custom Docker-in-Docker image with CA certificate
- **When**: Only when `Dockerfile.dind` changes or manually triggered
- **Output**: Custom DinD image pushed to registry

### 2. Lint Stage
- **Purpose**: Code quality and style checks
- **Tools**: Ruff (replaces flake8, black, isort)
- **Skip**: Set `QUICK_DEPLOY=true` to skip

### 3. Test Stage
- **Purpose**: Runs unit tests and security scans
- **Jobs**:
  - `test`: pytest with coverage reporting
  - `security_scan`: bandit, safety, pip-audit, detect-secrets
- **Skip**: Set `QUICK_DEPLOY=true` to skip

### 4. Build Stage
- **Purpose**: Build Docker image and push to registry
- **Output**: Multiple image tags (version, sha, latest)
- **Artifacts**: `build.env` with version info

### 5. Scan Stage
- **Purpose**: Container vulnerability scanning
- **Tool**: Trivy (scans OS packages and dependencies)
- **Skip**: Set `QUICK_DEPLOY=true` to skip

### 6. Release Stage
- **Purpose**: Validates build quality before deployment
- **Checks**: Vulnerability scan results, build artifacts
- **Output**: Release manifest

### 7. Deploy Stage
- **Purpose**: Deploys to production via SSH + docker-compose
- **Method**: Pulls image, updates compose, runs health check
- **Trigger**: Auto-deploy on `main` branch or manual

## Configuration Options

### Quick Deploy Mode

Speed up development by skipping heavy stages:

```yaml
variables:
  QUICK_DEPLOY: "true"  # Skip lint, test, scan stages
```

**Warning**: Never use `QUICK_DEPLOY=true` for production releases!

### Deployment Rules

**Auto-deploy on main branch:**
```yaml
rules:
  - if: '$CI_COMMIT_BRANCH == "main"'
    when: on_success
```

**Manual deployment only:**
```yaml
rules:
  - when: manual
```

**Deploy on tags:**
```yaml
rules:
  - if: '$CI_COMMIT_TAG'
    when: on_success
```

### Security Scanning Configuration

**Adjust vulnerability severity:**

```yaml
# Only block on CRITICAL (default: HIGH,CRITICAL)
trivy image --severity CRITICAL --exit-code 1 $DOCKER_IMAGE

# Scan all severities
trivy image --severity UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL $DOCKER_IMAGE
```

**Block deployment on vulnerabilities:**

```yaml
scan_image:
  script:
    # Add exit code to fail job on vulnerabilities
    - trivy image --severity HIGH,CRITICAL --exit-code 1 $DOCKER_IMAGE
  allow_failure: false  # Change from true to false
```

### Python Version Override

**Use different Python version for specific jobs:**

```yaml
test:
  image: python:3.11-slim  # Override global PYTHON_VERSION
```

## Troubleshooting

### Common Issues

**❌ "mapping values are not allowed here"**
- **Cause**: Colons in echo statements confuse YAML parser
- **Fix**: Use `--` instead of `:` in echo statements
- **Example**: `echo "Version--${VERSION}"` not `echo "Version: ${VERSION}"`

**❌ Docker login fails with SSL certificate error**
- **Cause**: GitLab registry uses self-signed certificate
- **Fix**: Configure `LOCAL_CA_CERT` variable (see Home Lab CA setup)

**❌ Deployment fails: "Permission denied (publickey)"**
- **Cause**: SSH key not configured or permissions incorrect
- **Fix**:
  1. Verify `DEPLOY_SSH_KEY` variable contains private key
  2. Confirm public key is in `~/.ssh/authorized_keys` on target
  3. Test SSH access manually: `ssh -i deploy_key user@host`

**❌ Tests not found**
- **Cause**: Test files not following naming convention
- **Fix**: Name test files as `test_*.py` or `*_test.py`

**❌ Docker-in-Docker service not starting**
- **Cause**: GitLab Runner not configured with privileged mode
- **Fix**: Enable `privileged = true` in `/etc/gitlab-runner/config.toml`

**❌ Health check fails after deployment**
- **Cause**: Application not ready or wrong port
- **Fix**:
  1. Increase sleep time before health check
  2. Verify `DEPLOY_PORT` matches container's exposed port
  3. Check application logs: `docker-compose logs`

### Debug Pipeline Issues

**Enable verbose logging:**

```yaml
script:
  - set -x  # Enable bash debug mode
  - docker build -t myapp .
  - set +x  # Disable debug mode
```

**View Docker logs:**

```yaml
script:
  - docker build -t myapp . || (docker logs myapp && exit 1)
```

**Test SSH connection:**

```yaml
script:
  - ssh -i ~/.ssh/deploy_key -v ${DEPLOY_USER}@${DEPLOY_HOST} "echo 'SSH works!'"
```

## Customization Examples

### Add Database Migrations

```yaml
deploy_to_platform:
  script:
    # ... existing deployment steps ...
    - ssh -i ~/.ssh/deploy_key ${DEPLOY_USER}@${DEPLOY_HOST} "cd ${DEPLOY_PATH} && docker-compose exec -T app flask db upgrade"
```

### Add Slack Notifications

```yaml
.notify_slack:
  after_script:
    - |
      curl -X POST -H 'Content-type: application/json' \
        --data "{\"text\":\"Deployment ${CI_JOB_STATUS}: ${CI_PROJECT_NAME}\"}" \
        ${SLACK_WEBHOOK_URL}

deploy_to_platform:
  extends: .notify_slack
```

### Multi-Environment Deployment

```yaml
deploy_staging:
  extends: .deploy
  variables:
    DEPLOY_HOST: "staging.example.com"
    DEPLOY_PATH: "/opt/app-staging"
  environment:
    name: staging
    url: https://staging.example.com
  rules:
    - if: '$CI_COMMIT_BRANCH == "develop"'

deploy_production:
  extends: .deploy
  variables:
    DEPLOY_HOST: "production.example.com"
    DEPLOY_PATH: "/opt/app"
  environment:
    name: production
    url: https://app.example.com
  rules:
    - if: '$CI_COMMIT_BRANCH == "main"'
```

## Best Practices

### Version Management
- Use semantic versioning: `MAJOR.MINOR.PATCH`
- Update `VERSION` file before releases
- Tag releases: `git tag v1.0.0 && git push --tags`

### Security
- Never commit secrets to git
- Use GitLab Protected Variables for production credentials
- Review security scan results regularly
- Keep base images updated
- Rotate SSH keys periodically

### Testing
- Maintain >80% test coverage
- Write integration tests for critical paths
- Mock external dependencies
- Use fixtures for test data

### Documentation
- Document environment variables in README
- Keep deployment runbook updated
- Add comments for complex logic
- Track breaking changes in CHANGELOG

## Additional Resources

- [GitLab CI/CD Documentation](https://docs.gitlab.com/ee/ci/)
- [Docker Documentation](https://docs.docker.com/)
- [Trivy Security Scanner](https://aquasecurity.github.io/trivy/)
- [pytest Documentation](https://docs.pytest.org/)
- [Ruff Linter](https://docs.astral.sh/ruff/)

## Support

For issues or questions:
1. Check troubleshooting section above
2. Review GitLab CI/CD logs in pipeline view
3. Verify all required variables are configured
4. Test SSH/Docker access manually
5. Consult `.claude/docs/gitlab-cicd-guide.md` for YAML formatting rules
