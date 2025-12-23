# Secure Pipeline Template Structure

## Directory Layout

```
.claude/templates/secure-pipeline/
├── README.md                      # Comprehensive documentation
├── QUICKSTART.md                  # 5-minute setup guide
├── TEMPLATE_STRUCTURE.md          # This file - architecture overview
├── gitlab-ci.yml.template         # Main CI/CD pipeline configuration
├── Dockerfile.template            # Multi-stage application container
├── Dockerfile.dind.template       # Custom Docker-in-Docker with CA certs
├── docker-compose.yml.template    # Deployment orchestration
└── variables.env.example          # Complete variable reference
```

## Template File Relationships

```
┌─────────────────────────────────────────────────────────────────┐
│                     GitLab CI/CD Pipeline                       │
│                  (gitlab-ci.yml.template)                       │
└────────┬────────────────────────────────────┬───────────────────┘
         │                                    │
         │                                    │
    ┌────▼────────┐                     ┌────▼─────────────┐
    │   PREPARE   │                     │      BUILD       │
    │   STAGE     │                     │      STAGE       │
    │             │                     │                  │
    │ Build DinD  │────uses────┐        │  Build App       │
    │   Image     │            │        │  Container       │
    └─────────────┘            │        └────┬─────────────┘
         │                     │             │
         │                ┌────▼─────────────▼────┐
         │                │ Dockerfile.dind       │
         │                │    .template          │
         │                │                       │
         │                │ - CA certificates     │
         │                │ - Registry trust      │
         │                │ - DinD daemon config  │
         │                └───────────────────────┘
         │
         │                ┌───────────────────────┐
         └────creates────▶│  Custom DinD Image    │
                          │  (in registry)        │
                          └───────┬───────────────┘
                                  │
                          ┌───────▼───────────────┐
                          │ Dockerfile.template   │
                          │                       │
                          │ - Multi-stage build   │
                          │ - Non-root user       │
                          │ - Virtual environment │
                          │ - Health checks       │
                          └───────┬───────────────┘
                                  │
                          ┌───────▼───────────────┐
                          │  Application Image    │
                          │  (in registry)        │
                          └───────┬───────────────┘
                                  │
                          ┌───────▼───────────────┐
                          │   SCAN STAGE          │
                          │                       │
                          │ - Trivy security scan │
                          │ - Vulnerability check │
                          └───────┬───────────────┘
                                  │
                          ┌───────▼───────────────┐
                          │   RELEASE GATE        │
                          │                       │
                          │ - Validate scan       │
                          │ - Block on HIGH/CRIT  │
                          └───────┬───────────────┘
                                  │
                          ┌───────▼───────────────┐
                          │   DEPLOY STAGE        │
                          │                       │
                          │ - SSH to target       │
                          │ - Deploy with compose │
                          └───────┬───────────────┘
                                  │
                          ┌───────▼───────────────┐
                          │ docker-compose.yml    │
                          │     .template         │
                          │                       │
                          │ - Service definition  │
                          │ - Health checks       │
                          │ - Volume mounts       │
                          │ - Network config      │
                          └───────────────────────┘
```

## Variable Flow

```
┌──────────────────────────┐
│  variables.env.example   │  ◄─── Reference for all variables
│                          │
│  Documents:              │
│  - Application config    │
│  - Registry settings     │
│  - Deployment targets    │
│  - Security options      │
└────────┬─────────────────┘
         │
         │ Variables used in:
         │
         ├────────────────────────────────────────┐
         │                                        │
    ┌────▼──────────────┐              ┌─────────▼────────────┐
    │ gitlab-ci.yml     │              │ Dockerfile           │
    │                   │              │                      │
    │ APP_NAME          │              │ APP_USER             │
    │ APP_DOMAIN        │              │ DEPLOY_PORT          │
    │ DEPLOY_PORT       │              │ FLASK_APP            │
    │ CI_REGISTRY       │              │ WSGI_MODULE          │
    │ CI_REGISTRY_IMAGE │              │ WORKERS              │
    │ VERSION_FILE      │              └──────────────────────┘
    │ DEPLOY_HOST       │
    │ DEPLOY_USER       │              ┌──────────────────────┐
    │ DEPLOY_PATH       │              │ docker-compose.yml   │
    └───────────────────┘              │                      │
                                       │ APP_NAME             │
    ┌────────────────────┐             │ DEPLOY_PORT          │
    │ Dockerfile.dind    │             │ DATA_PATH            │
    │                    │             │ CI_REGISTRY_IMAGE    │
    │ CI_REGISTRY        │             └──────────────────────┘
    └────────────────────┘
```

## Security Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Security Layers                          │
└─────────────────────────────────────────────────────────────────┘

Layer 1: Build Security
┌─────────────────────────────────────────┐
│ Multi-Stage Build (Dockerfile.template) │
│ ┌─────────────────────────────────────┐ │
│ │ Builder Stage                       │ │
│ │ - Install dependencies              │ │
│ │ - Create virtual environment        │ │
│ └─────────────────────────────────────┘ │
│ ┌─────────────────────────────────────┐ │
│ │ Production Stage                    │ │
│ │ - Minimal base image                │ │
│ │ - Copy only venv                    │ │
│ │ - No build tools                    │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘

Layer 2: Runtime Security
┌─────────────────────────────────────────┐
│ Non-Root User Execution                 │
│ ┌─────────────────────────────────────┐ │
│ │ USER {{APP_USER}}                   │ │
│ │ - No root privileges                │ │
│ │ - Limited filesystem access         │ │
│ │ - Defense in depth                  │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘

Layer 3: Vulnerability Management
┌─────────────────────────────────────────┐
│ Trivy Security Scanning                 │
│ ┌─────────────────────────────────────┐ │
│ │ - Scan for CVEs                     │ │
│ │ - Check HIGH/CRITICAL severity      │ │
│ │ - Generate detailed reports         │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘

Layer 4: Release Controls
┌─────────────────────────────────────────┐
│ Release Gate Validation                 │
│ ┌─────────────────────────────────────┐ │
│ │ - Block vulnerable images           │ │
│ │ - Require scan completion           │ │
│ │ - Audit trail via artifacts         │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘

Layer 5: Deployment Security
┌─────────────────────────────────────────┐
│ SSH + Isolated Network                  │
│ ┌─────────────────────────────────────┐ │
│ │ - SSH key authentication            │ │
│ │ - Localhost port binding            │ │
│ │ - Bridge network isolation          │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

## Pipeline Stages Detail

```
┌──────────────────────────────────────────────────────────────────┐
│                    PREPARE STAGE (Optional)                      │
├──────────────────────────────────────────────────────────────────┤
│ Job: build_dind_image                                            │
│                                                                  │
│ Triggers:                                                        │
│  - Changes to Dockerfile.dind                                    │
│  - Manual trigger                                                │
│                                                                  │
│ Purpose:                                                         │
│  - Build custom Docker-in-Docker image                           │
│  - Install CA certificates for registry trust                    │
│  - Cache for subsequent pipeline runs                            │
│                                                                  │
│ Output: $CI_REGISTRY_IMAGE/docker-dind:27-with-cert              │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                     LINT STAGE (Conditional)                     │
├──────────────────────────────────────────────────────────────────┤
│ Job: lint                                                        │
│                                                                  │
│ Skip when: QUICK_DEPLOY=true                                     │
│                                                                  │
│ Tools:                                                           │
│  - Ruff (Python code quality)                                    │
│                                                                  │
│ Actions:                                                         │
│  - Auto-fix code style issues                                    │
│  - Validate final code quality                                   │
│                                                                  │
│ Exit: Fail on remaining violations                               │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                     TEST STAGE (Conditional)                     │
├──────────────────────────────────────────────────────────────────┤
│ Job: test                                                        │
│                                                                  │
│ Skip when: QUICK_DEPLOY=true                                     │
│                                                                  │
│ Tools:                                                           │
│  - pytest                                                        │
│                                                                  │
│ Actions:                                                         │
│  - Install dependencies                                          │
│  - Run all tests                                                 │
│                                                                  │
│ Exit: Fail on test failures                                      │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                       BUILD STAGE (Required)                     │
├──────────────────────────────────────────────────────────────────┤
│ Job: build_image                                                 │
│                                                                  │
│ Actions:                                                         │
│  1. Read VERSION file                                            │
│  2. Generate build metadata (.buildinfo)                         │
│  3. Build Docker image (multi-stage)                             │
│  4. Tag with multiple strategies:                                │
│     - VERSION-SHA (1.0.0-a1b2c3d)                                │
│     - sha-FULL_SHA (immutable reference)                         │
│     - VERSION (semver)                                           │
│     - latest (convenience)                                       │
│  5. Push all tags to registry                                    │
│  6. Export metadata to build.env                                 │
│                                                                  │
│ Artifacts:                                                       │
│  - build.env (dotenv report)                                     │
│  - .buildinfo (build metadata)                                   │
│                                                                  │
│ Outputs:                                                         │
│  - IMAGE_DIGEST                                                  │
│  - DOCKER_IMAGE                                                  │
│  - VERSION                                                       │
│  - VERSION_TAG                                                   │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                     SCAN STAGE (Conditional)                     │
├──────────────────────────────────────────────────────────────────┤
│ Job: scan_image                                                  │
│                                                                  │
│ Skip when: QUICK_DEPLOY=true                                     │
│                                                                  │
│ Tools:                                                           │
│  - Trivy vulnerability scanner                                   │
│                                                                  │
│ Actions:                                                         │
│  1. Pull built image from registry                              │
│  2. Scan for vulnerabilities (HIGH/CRITICAL)                     │
│  3. Generate JSON report                                         │
│  4. Display human-readable summary                               │
│                                                                  │
│ Artifacts:                                                       │
│  - trivy.json (30 day retention)                                 │
│                                                                  │
│ Exit: Allow failure (report only)                                │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                    RELEASE STAGE (Required)                      │
├──────────────────────────────────────────────────────────────────┤
│ Job: release_gate                                                │
│                                                                  │
│ Dependencies:                                                    │
│  - build_image (required)                                        │
│  - scan_image (optional)                                         │
│                                                                  │
│ Validation:                                                      │
│  1. Check QUICK_DEPLOY mode (skip checks if true)                │
│  2. Analyze trivy.json for HIGH/CRITICAL CVEs                    │
│  3. Block deployment if vulnerabilities found                    │
│  4. Create release artifact                                      │
│                                                                  │
│ Artifacts:                                                       │
│  - release/image.txt (approved image reference)                  │
│                                                                  │
│ Exit: Fail on HIGH/CRITICAL vulnerabilities                      │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                     DEPLOY STAGE (Manual/Auto)                   │
├──────────────────────────────────────────────────────────────────┤
│ Job: deploy_to_platform                                          │
│                                                                  │
│ Triggers:                                                        │
│  - Automatic on main branch                                      │
│  - Manual trigger from any branch                                │
│                                                                  │
│ Dependencies:                                                    │
│  - build_image (for DOCKER_IMAGE variable)                       │
│  - release_gate (ensures security validation)                    │
│                                                                  │
│ Actions:                                                         │
│  1. Setup SSH connection                                         │
│  2. Copy docker-compose.yml to target                            │
│  3. Login to registry on target                                  │
│  4. Pull latest image                                            │
│  5. Deploy with docker-compose up -d                             │
│  6. Wait for container startup                                   │
│  7. Run health check (curl)                                      │
│                                                                  │
│ Environment:                                                     │
│  - Name: production                                              │
│  - URL: https://$APP_DOMAIN                                      │
│                                                                  │
│ Exit: Fail on health check failure                               │
└──────────────────────────────────────────────────────────────────┘
```

## Deployment Flow

```
GitLab CI/CD Runner          Deployment Target Server
─────────────────────        ─────────────────────────

┌─────────────────┐
│ Deploy Job      │
│                 │
│ 1. Setup SSH    │
└────────┬────────┘
         │
         │ SSH Connection
         │ (DEPLOY_SSH_KEY)
         │
         ├──────────────────────────┐
         │                          │
    ┌────▼─────────┐                │
    │ SCP          │                │
    │ docker-      │                │
    │ compose.yml  │────────────┐   │
    └──────────────┘            │   │
                                │   │
                           ┌────▼───▼──────────────┐
                           │ /srv/{{APP_NAME}}/    │
                           │                       │
                           │ docker-compose.yml    │
                           └───────┬───────────────┘
                                   │
                           ┌───────▼───────────────┐
                           │ docker login          │
                           │ (CI_JOB_TOKEN)        │
                           └───────┬───────────────┘
                                   │
                           ┌───────▼───────────────┐
                           │ docker-compose pull   │
                           │                       │
                           │ Pull image from       │
                           │ GitLab registry       │
                           └───────┬───────────────┘
                                   │
                           ┌───────▼───────────────┐
                           │ docker-compose up -d  │
                           │                       │
                           │ - Stop old container  │
                           │ - Start new container │
                           │ - Mount volumes       │
                           └───────┬───────────────┘
                                   │
                           ┌───────▼───────────────┐
                           │ Container Running     │
                           │                       │
                           │ - Port: 127.0.0.1:    │
                           │   {{DEPLOY_PORT}}     │
                           │ - User: {{APP_USER}}  │
                           │ - Network: bridge     │
                           └───────┬───────────────┘
                                   │
         ┌──────────────────────┐  │
         │ Health Check (SSH)   │  │
         │ curl localhost:PORT  │◄─┘
         └──────────────────────┘
```

## Advanced Features (Optional)

These features are available in `.gitlab-ci.yml.backup` but not included in the basic template:

```
┌──────────────────────────────────────────────────────────────────┐
│                    OPTIONAL: SBOM Generation                     │
├──────────────────────────────────────────────────────────────────┤
│ Tool: Syft                                                       │
│                                                                  │
│ Actions:                                                         │
│  - Extract software components from image                        │
│  - Generate SPDX JSON format                                     │
│  - Include in release artifacts                                  │
│                                                                  │
│ Benefits:                                                        │
│  - Supply chain transparency                                     │
│  - License compliance                                            │
│  - Dependency tracking                                           │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                    OPTIONAL: Image Signing                       │
├──────────────────────────────────────────────────────────────────┤
│ Tool: Cosign                                                     │
│                                                                  │
│ Actions:                                                         │
│  - Sign image digest                                             │
│  - Sign SBOM                                                     │
│  - Sign vulnerability report                                     │
│  - Verify signatures in release gate                             │
│                                                                  │
│ Benefits:                                                        │
│  - Cryptographic provenance                                      │
│  - Tamper detection                                              │
│  - Supply chain security                                         │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                  OPTIONAL: Multiple Scanners                     │
├──────────────────────────────────────────────────────────────────┤
│ Tools: Trivy + Grype                                             │
│                                                                  │
│ Actions:                                                         │
│  - Run parallel scans                                            │
│  - Compare results                                               │
│  - Increase detection confidence                                 │
│                                                                  │
│ Benefits:                                                        │
│  - Reduced false negatives                                       │
│  - Comprehensive coverage                                        │
│  - Cross-validation                                              │
└──────────────────────────────────────────────────────────────────┘
```

## Customization Points

```
┌────────────────────────────────────────────────────────────────┐
│                    Template Customization                      │
└────────────────────────────────────────────────────────────────┘

1. Application Stack
   File: Dockerfile.template
   Change: Base image, build commands, runtime
   Example: python:3.12-slim → node:20-alpine

2. Test Framework
   File: gitlab-ci.yml.template (test stage)
   Change: Test commands and tools
   Example: pytest → jest

3. Linting Tools
   File: gitlab-ci.yml.template (lint stage)
   Change: Linter tool and configuration
   Example: ruff → eslint

4. Deployment Method
   File: gitlab-ci.yml.template (deploy stage)
   Change: Deployment mechanism
   Example: docker-compose → kubectl

5. Security Scanning
   File: gitlab-ci.yml.template (scan stage)
   Change: Scanner tool or add additional scanners
   Example: trivy → grype + trivy

6. Runner Tags
   File: gitlab-ci.yml.template (tags)
   Change: Runner selection
   Example: docker → kubernetes

7. Build Stages
   File: gitlab-ci.yml.template
   Change: Add or remove stages
   Example: Add integration_test stage

8. Environment Variables
   File: All templates
   Change: Add application-specific configuration
   Example: Add DATABASE_URL
```

## File Modification Priority

When customizing for your application:

```
MUST MODIFY (Required for basic functionality):
├── gitlab-ci.yml.template      ⭐⭐⭐⭐⭐
│   Replace ALL {{VARIABLES}}
├── Dockerfile.template         ⭐⭐⭐⭐⭐
│   Replace ALL {{VARIABLES}}
├── docker-compose.yml.template ⭐⭐⭐⭐⭐
│   Replace ALL {{VARIABLES}}
└── Dockerfile.dind.template    ⭐⭐⭐⭐
    Replace CI_REGISTRY variable

SHOULD MODIFY (Customize for your needs):
├── Dockerfile.template         ⭐⭐⭐
│   Adjust build steps for your app
├── gitlab-ci.yml.template      ⭐⭐⭐
│   Customize test/lint commands
└── docker-compose.yml.template ⭐⭐
    Add application-specific env vars

OPTIONAL MODIFY (Advanced features):
├── gitlab-ci.yml.template      ⭐⭐
│   Add SBOM/signing stages
└── Dockerfile.template         ⭐
    Add additional build optimizations
```

## Integration with Existing Projects

```
Scenario 1: New Project
┌──────────────────────────────────────┐
│ 1. Copy all template files          │
│ 2. Replace variables                │
│ 3. Create VERSION file              │
│ 4. Configure GitLab variables       │
│ 5. Commit and push                  │
└──────────────────────────────────────┘

Scenario 2: Existing .gitlab-ci.yml
┌──────────────────────────────────────┐
│ 1. Backup existing pipeline         │
│ 2. Copy template as .gitlab-ci.new  │
│ 3. Merge custom jobs/stages         │
│ 4. Test in feature branch           │
│ 5. Replace when validated           │
└──────────────────────────────────────┘

Scenario 3: Existing Dockerfile
┌──────────────────────────────────────┐
│ 1. Review Dockerfile.template       │
│ 2. Extract security patterns:       │
│    - Multi-stage build              │
│    - Non-root user                  │
│    - Health check                   │
│ 3. Merge into existing Dockerfile  │
│ 4. Test build locally               │
│ 5. Update pipeline to use patterns │
└──────────────────────────────────────┘
```

## Success Criteria

After implementing this template, you should have:

```
✅ Automated CI/CD pipeline
✅ Security vulnerability scanning
✅ Multi-stage optimized Docker images
✅ Non-root container execution
✅ Health check monitoring
✅ Release gate validation
✅ Automated deployment to target server
✅ Build metadata tracking
✅ Container image tagging strategy
✅ SSH-based secure deployment
```

## Next Steps

1. Review `README.md` for comprehensive documentation
2. Follow `QUICKSTART.md` for rapid setup
3. Customize templates for your application
4. Test in development environment
5. Deploy to production
6. Monitor and optimize

## References

- GitLab CI/CD: https://docs.gitlab.com/ee/ci/
- Docker Multi-Stage Builds: https://docs.docker.com/build/building/multi-stage/
- Trivy: https://aquasecurity.github.io/trivy/
- Docker Security: https://docs.docker.com/develop/security-best-practices/
