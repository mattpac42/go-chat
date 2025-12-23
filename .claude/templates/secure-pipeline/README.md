# Secure GitLab CI/CD Release Pipeline Template

This template provides a production-ready, security-focused GitLab CI/CD pipeline for containerized applications. It implements industry best practices for building, scanning, and deploying Docker containers with comprehensive security controls.

## Features

### Core Pipeline Capabilities
- **Multi-stage Docker builds** for optimized, minimal production images
- **Security vulnerability scanning** with Trivy
- **Release gates** to prevent deployment of vulnerable images
- **Build metadata tracking** for full auditability
- **Automated deployment** via SSH and Docker Compose
- **Health checks** to validate successful deployments

### Security Features
- **Non-root container execution** for defense in depth
- **Custom CA certificate support** for internal/self-signed registries
- **Multiple image tagging strategies** (version, SHA, latest)
- **Vulnerability scanning** with configurable severity thresholds
- **Optional image signing** with Cosign (see advanced features)
- **SBOM generation** with Syft (see advanced features)

### Pipeline Optimization
- **Quick deploy mode** to skip heavy stages during development
- **Conditional job execution** based on file changes
- **Artifact caching** to minimize redundant work
- **Parallel job execution** where possible

## Quick Start

### 1. Copy Template Files

Copy all template files to your project root:

```bash
cp .claude/templates/secure-pipeline/gitlab-ci.yml.template .gitlab-ci.yml
cp .claude/templates/secure-pipeline/Dockerfile.template Dockerfile
cp .claude/templates/secure-pipeline/Dockerfile.dind.template Dockerfile.dind
cp .claude/templates/secure-pipeline/docker-compose.yml.template docker-compose.yml
cp .claude/templates/secure-pipeline/variables.env.example variables.env
```

### 2. Create VERSION File

Create a VERSION file in your project root containing your application version:

```bash
echo "1.0.0" > VERSION
```

### 3. Configure Variables

Edit the copied template files and replace all `{{VARIABLE}}` placeholders:

#### Required Variables

| Variable | Example | Description |
|----------|---------|-------------|
| `{{APP_NAME}}` | `myapp` | Application name (lowercase, no spaces) |
| `{{APP_DOMAIN}}` | `myapp.example.com` | Public domain for your application |
| `{{DEPLOY_PORT}}` | `4000` | Application port inside container |
| `{{CI_REGISTRY}}` | `gitlab.example.com:5050` | GitLab registry URL |
| `{{CI_REGISTRY_IMAGE}}` | `gitlab.example.com:5050/group/project` | Full registry path |
| `{{VERSION_FILE}}` | `VERSION` | File containing version number |
| `{{APP_USER}}` | `appuser` | Non-root user to run application |
| `{{FLASK_APP}}` | `app.py` | Main Flask application file |
| `{{WSGI_MODULE}}` | `wsgi:application` | WSGI module for Gunicorn |
| `{{WORKERS}}` | `1` | Number of Gunicorn workers |
| `{{DEPLOY_HOST}}` | `192.168.0.100` | Deployment target IP/hostname |
| `{{DEPLOY_USER}}` | `deployer` | SSH user for deployment |
| `{{DEPLOY_PATH}}` | `/srv/myapp` | Deployment directory on target |
| `{{DATA_PATH}}` | `/srv/myapp/data` | Persistent data storage path |

### 4. Configure GitLab CI/CD Variables

In your GitLab project, go to **Settings > CI/CD > Variables** and add:

#### Required Variables

| Variable | Type | Protected | Masked | Description |
|----------|------|-----------|--------|-------------|
| `DEPLOY_SSH_KEY` | File | Yes | No | SSH private key for deployment |
| `CI_REGISTRY_USER` | Variable | No | No | Auto-provided by GitLab |
| `CI_REGISTRY_PASSWORD` | Variable | No | Yes | Auto-provided by GitLab |

#### Optional Variables

| Variable | Type | Protected | Masked | Description |
|----------|------|-----------|--------|-------------|
| `LOCAL_CA_CERT` | File | No | No | Base64-encoded custom CA certificate |
| `QUICK_DEPLOY` | Variable | No | No | Set to `true` for fast testing |

### 5. Configure GitLab Runner

Ensure you have a GitLab runner with:
- **Tag**: `docker`
- **Executor**: `docker`
- **Privileged mode**: Enabled (required for Docker-in-Docker)

Example runner configuration in `/etc/gitlab-runner/config.toml`:

```toml
[[runners]]
  name = "docker-runner"
  url = "https://gitlab.example.com"
  token = "YOUR_RUNNER_TOKEN"
  executor = "docker"
  [runners.docker]
    image = "docker:27"
    privileged = true
    volumes = ["/cache", "/var/run/docker.sock:/var/run/docker.sock"]
  [runners.cache]
    Type = "volume"
```

### 6. Set Up Deployment Target

On your deployment target server:

```bash
# Create deployment directory
sudo mkdir -p /srv/myapp/data
sudo chown deployer:deployer /srv/myapp

# Install Docker and Docker Compose
sudo apt update
sudo apt install docker.io docker-compose

# Add deployer user to docker group
sudo usermod -aG docker deployer

# Generate SSH key for deployment (if needed)
ssh-keygen -t ed25519 -C "gitlab-deploy"
# Add the public key to ~/.ssh/authorized_keys on deployment target
```

### 7. Commit and Push

```bash
git add .gitlab-ci.yml Dockerfile Dockerfile.dind docker-compose.yml VERSION
git commit -m "Add secure CI/CD pipeline"
git push origin main
```

The pipeline will automatically run on push to `main` branch.

## Pipeline Stages

### 1. Prepare Stage
- Builds custom Docker-in-Docker image with CA certificates
- Only runs when `Dockerfile.dind` changes or manually triggered
- Caches image to speed up subsequent builds

### 2. Lint Stage
- Runs code quality checks with Ruff (Python)
- Auto-fixes common issues
- Can be skipped with `QUICK_DEPLOY=true`

### 3. Test Stage
- Runs pytest tests
- Automatically detects test files
- Can be skipped with `QUICK_DEPLOY=true`

### 4. Build Stage
- Creates multi-stage Docker image
- Generates build metadata (`.buildinfo`)
- Tags image with multiple strategies:
  - `VERSION-SHA` (e.g., `1.0.0-a1b2c3d`)
  - `sha-FULL_SHA` (e.g., `sha-a1b2c3d4e5f6...`)
  - `VERSION` (e.g., `1.0.0`)
  - `latest`
- Pushes to GitLab container registry

### 5. Scan Stage
- Scans image for vulnerabilities with Trivy
- Reports HIGH and CRITICAL severity issues
- Generates JSON report artifact
- Can be skipped with `QUICK_DEPLOY=true`

### 6. Release Stage
- Validates scan results
- Blocks deployment if HIGH/CRITICAL vulnerabilities found
- Creates release artifact with image details
- Can be bypassed with `QUICK_DEPLOY=true`

### 7. Deploy Stage
- Copies `docker-compose.yml` to target server
- Pulls latest image
- Deploys with Docker Compose
- Runs health check
- Only runs on `main` branch or manually

## Usage Patterns

### Development Workflow

For fast iteration during development:

```yaml
# In .gitlab-ci.yml, set:
QUICK_DEPLOY: "true"
```

This skips lint, test, and scan stages, deploying directly after build.

**WARNING**: Only use for development! Always run full pipeline before production release.

### Production Release

1. Ensure `QUICK_DEPLOY` is set to `false`
2. Push to `main` branch
3. Pipeline runs all stages
4. Review scan results in job logs
5. If vulnerabilities found, fix and retry
6. Deploy stage runs automatically after successful release gate

### Manual Deployment

To deploy a specific commit:

1. Go to **CI/CD > Pipelines** in GitLab
2. Find the pipeline for your commit
3. Click the manual play button on `deploy_to_platform` job

### Rebuilding DinD Image

The custom Docker-in-Docker image only rebuilds when `Dockerfile.dind` changes. To force a rebuild (e.g., for certificate rotation):

1. Go to **CI/CD > Pipelines**
2. Click **Run pipeline**
3. Manually trigger the `build_dind_image` job

## Custom CA Certificates

If your GitLab registry uses self-signed certificates or an internal CA:

### 1. Encode Your CA Certificate

```bash
base64 -w 0 /path/to/your/ca.crt > ca-base64.txt
```

### 2. Add to GitLab CI/CD Variables

1. Go to **Settings > CI/CD > Variables**
2. Add variable:
   - **Key**: `LOCAL_CA_CERT`
   - **Type**: File
   - **Value**: Contents of `ca-base64.txt`
   - **Protected**: No
   - **Masked**: No

### 3. Rebuild DinD Image

Trigger the `build_dind_image` job manually to incorporate the new certificate.

## Health Checks

The pipeline expects a health check endpoint in your application. Implement one of:

### Option 1: Dedicated Health Endpoint

```python
@app.route('/healthz')
def health_check():
    return {'status': 'healthy'}, 200
```

### Option 2: Root Endpoint

```python
@app.route('/')
def index():
    return 'OK', 200
```

Update the health check URL in templates if using a different endpoint:
- `Dockerfile.template`: HEALTHCHECK directive
- `docker-compose.yml.template`: healthcheck.test
- `gitlab-ci.yml.template`: deploy stage curl command

## Troubleshooting

### Pipeline Fails at Build Stage

**Symptom**: Docker build fails with registry certificate errors

**Solution**:
1. Verify `LOCAL_CA_CERT` is set correctly
2. Rebuild DinD image: manually trigger `build_dind_image`
3. Check GitLab registry URL in `CI_REGISTRY` variable

### Pipeline Fails at Scan Stage

**Symptom**: Trivy cannot download vulnerability database

**Solution**:
1. Check runner has internet access
2. Verify firewall allows HTTPS to `ghcr.io`
3. Consider running Trivy in offline mode (advanced)

### Deployment Fails at Health Check

**Symptom**: Deployment completes but health check fails

**Solution**:
1. Check application logs: `docker logs <container>`
2. Verify port mapping: `docker ps`
3. Test health endpoint manually: `curl http://localhost:PORT/healthz`
4. Check firewall rules on deployment target

### Release Gate Blocks Deployment

**Symptom**: HIGH/CRITICAL vulnerabilities prevent deployment

**Solution**:
1. Review vulnerability report in scan job logs
2. Update base image in Dockerfile
3. Update dependencies in requirements.txt
4. For false positives, add exceptions to scan job (advanced)

## Advanced Features

The `.gitlab-ci.yml.backup` file in this project includes advanced features not included in the basic template:

### SBOM Generation
- Software Bill of Materials with Syft
- SPDX JSON format
- Included in release artifacts

### Image Signing
- Cosign-based cryptographic signing
- Signature verification in release gate
- Requires `COSIGN_PRIVATE_KEY` variable

### Advanced Security Scanning
- Multiple scanners (Trivy + Grype)
- Custom vulnerability policies
- Detailed reporting

To enable these features, refer to `.gitlab-ci.yml.backup` and merge desired sections into your pipeline.

## Security Considerations

### Non-Root Containers
All containers run as non-root users for defense in depth. The application user is created during image build.

### Least Privilege
- Deployment SSH key should only have access to deployment directory
- Service users should not have sudo privileges
- Container networks are isolated

### Secrets Management
- Never commit secrets to Git
- Use GitLab CI/CD variables for all sensitive data
- Mark variables as "Masked" where possible
- Use "Protected" for production secrets

### Supply Chain Security
- Multi-stage builds minimize attack surface
- Vulnerability scanning catches known CVEs
- Image signing ensures provenance (advanced)
- SBOM provides transparency (advanced)

## Maintenance

### Regular Updates
- Update base images monthly: `FROM python:3.12-slim`
- Update dependencies: `pip list --outdated`
- Rebuild DinD image quarterly for security patches
- Review and update GitLab Runner version

### Monitoring
- Check pipeline success rates
- Review vulnerability scan trends
- Monitor deployment health checks
- Track build times and optimize as needed

### Version Management
- Update `VERSION` file for each release
- Use semantic versioning (MAJOR.MINOR.PATCH)
- Tag releases in Git: `git tag v1.0.0`

## Customization

### Different Tech Stacks

This template is designed for Python/Flask applications. To adapt for other stacks:

1. **Node.js**: Replace Python base image with `node:20-slim`, adjust build commands
2. **Go**: Use `golang:1.21-alpine` for build, `scratch` or `alpine` for runtime
3. **Java**: Use `maven:3-openjdk-17` for build, `openjdk:17-jre-slim` for runtime

### Different Deployment Methods

This template uses Docker Compose. Alternatives:

1. **Kubernetes**: Replace deploy stage with `kubectl apply`
2. **Docker Swarm**: Use `docker stack deploy`
3. **Nomad**: Use `nomad job run`

### Custom Build Tools

To use different linters, test frameworks, or scanners:

1. Update relevant stage in `.gitlab-ci.yml`
2. Change `image:` to appropriate tool container
3. Adjust `script:` commands
4. Update artifact paths if needed

## License

This template is provided as-is for use in your projects. Customize freely to meet your needs.

## Support

For issues or questions:
1. Check GitLab CI/CD logs for detailed error messages
2. Review this README for troubleshooting steps
3. Consult the source project's documentation
4. Review `.gitlab-ci.yml.backup` for advanced examples

## References

- [GitLab CI/CD Documentation](https://docs.gitlab.com/ee/ci/)
- [Docker Multi-Stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [Trivy Documentation](https://aquasecurity.github.io/trivy/)
- [Docker Security Best Practices](https://docs.docker.com/develop/security-best-practices/)
