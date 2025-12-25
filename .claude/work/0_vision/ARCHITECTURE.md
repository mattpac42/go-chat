# Go Chat MVP - Technical Architecture

**Created**: 2025-12-24
**Author**: Architect Agent
**Status**: Draft
**Document Type**: Architecture Decision Record (ADR) + System Design

---

## Executive Summary

This document defines the technical architecture for Go Chat MVP - a conversational development platform enabling non-technical users to build and deploy applications through natural language chat. The architecture prioritizes simplicity, rapid iteration, and maintainability over scalability, recognizing this is an MVP targeting 5-20 users in the first 6 months.

**Key Architectural Decisions**:
- Monolithic backend with clear module boundaries (not microservices)
- Single-tenant first, multi-tenant by month 6
- Self-hosted GitLab CE for version control and CI/CD
- Devcontainers for isolated development environments
- Claude API for AI-powered code generation

---

## 1. System Overview

### 1.1 High-Level Architecture Diagram

```
+------------------------------------------------------------------+
|                         USER DEVICES                               |
|   +------------+    +------------+    +------------+              |
|   |  Desktop   |    |   Mobile   |    |   Tablet   |              |
|   |  Browser   |    |  Browser   |    |  Browser   |              |
|   +-----+------+    +-----+------+    +-----+------+              |
+---------|-----------------|-----------------|-----------------------+
          |                 |                 |
          v                 v                 v
+------------------------------------------------------------------+
|                    REVERSE PROXY (Caddy/Traefik)                  |
|                 - TLS Termination                                  |
|                 - Request Routing                                  |
|                 - Rate Limiting                                    |
+---------------------------+--------------------------------------+
                            |
          +-----------------+------------------+
          |                                    |
          v                                    v
+-------------------+              +------------------------+
|  FRONTEND (SPA)   |              |   BACKEND (Monolith)   |
|                   |              |                        |
| - Chat Interface  |   REST/WS    | +--------------------+ |
| - Project List    |<------------>| |   API Gateway      | |
| - Build Status    |              | +--------------------+ |
| - Error Display   |              |          |            |
|                   |              |          v            |
| Tech: SvelteKit   |              | +--------------------+ |
+-------------------+              | | Chat Orchestrator  | |
                                   | +--------------------+ |
                                   |          |            |
                                   |    +-----+-----+      |
                                   |    |     |     |      |
                                   |    v     v     v      |
                                   | +---+ +---+ +---+     |
                                   | |AI | |DC | |GL |     |
                                   | |Svc| |Svc| |Svc|     |
                                   | +---+ +---+ +---+     |
                                   |                       |
                                   | Tech: Go              |
                                   +----------+------------+
                                              |
          +-----------------------------------+-----------------------------------+
          |                    |                    |                             |
          v                    v                    v                             v
+------------------+  +------------------+  +------------------+  +------------------+
|   PostgreSQL     |  |   Claude API     |  |   GitLab CE      |  | Container Runtime|
|                  |  |   (External)     |  |   (Self-hosted)  |  |   (Docker/Podman)|
| - Users          |  |                  |  |                  |  |                  |
| - Projects       |  | - Code Gen       |  | - Git Repos      |  | - Devcontainers  |
| - Conversations  |  | - NL Processing  |  | - CI/CD          |  | - App Previews   |
| - Build History  |  |                  |  | - Registry       |  |                  |
+------------------+  +------------------+  +------------------+  +------------------+
```

### 1.2 Component Interaction Flow

```
User Message Flow:
==================

[User] ---(chat)---> [Frontend] ---(WebSocket)---> [Backend]
                                                       |
                                                       v
                                              [Chat Orchestrator]
                                                       |
                     +----------------+----------------+----------------+
                     |                |                |                |
                     v                v                v                v
              [AI Service]    [Devcontainer    [GitLab         [Database]
               Claude API      Service]         Service]
                     |                |                |
                     v                v                v
              Code Generation   Container Mgmt   Git Operations
                     |                |                |
                     +-------+--------+----------------+
                             |
                             v
                     [Response Stream]
                             |
                             v
                        [Frontend]
                             |
                             v
                         [User]
```

### 1.3 Key Technology Choices

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Frontend | Next.js 14+ (App Router) | React ecosystem, excellent mobile PWA support, existing devcontainer |
| Backend | Go 1.22+ with Gin + Gorilla WebSocket | Fast, mature HTTP framework, industry-standard WebSocket |
| Database | PostgreSQL 16 | Battle-tested, JSON support, full-text search for conversations |
| AI/LLM | Claude API | Superior code generation, large context window, conversation memory |
| Version Control | GitLab CE | Self-hosted, integrated CI/CD, container registry, API-first |
| Containers | Podman (local) / Docker (prod) | OCI-compatible, rootless local dev, Docker for CI/CD |
| Reverse Proxy | Caddy | Automatic HTTPS, simple config, Go-native |

### 1.4 Development Environment Strategy

| Environment | Container Runtime | Compose Tool | Notes |
|-------------|-------------------|--------------|-------|
| Local/Dev (MacOS) | Podman (rootless) | podman-compose | Via .devcontainer/, native MacOS testing supported |
| Production | Docker | docker-compose | GitLab CI/CD → yuki-platform server (192.168.0.200) |

**Key insight**: Podman and Docker are OCI-compatible. Same Dockerfile, same images work in both environments.

**Gold Bricks Pipeline**: Production deployments use the secure pipeline templates at `.claude/templates/secure-pipeline/` which build Docker images, push to `gitlab.yuki.lan:5050`, and SSH deploy to production.

---

## 2. Core Components

### 2.1 Chat Interface (Frontend)

**Purpose**: Provide a conversational interface for users to describe applications and view build progress.

**Technology**: Next.js 14+ (App Router) with TypeScript

**Why Next.js**:
- React ecosystem with massive community and library support
- App Router provides streaming and React Server Components
- Already configured in project devcontainer
- Excellent PWA and mobile support
- Tailwind CSS integration for rapid styling

**Key Features**:
- Real-time chat with streaming responses (WebSocket)
- Project list and switching
- Build status visualization (progress bars, logs)
- Plain-language error display
- Mobile-responsive design (progressive enhancement)
- Offline message queue (sync when reconnected)

**Component Structure**:
```
src/
  app/
    layout.tsx           # Root layout with providers
    page.tsx             # Home/project list
    projects/[id]/
      page.tsx           # Project chat interface
  components/
    chat/
      ChatContainer.tsx
      MessageList.tsx
      MessageBubble.tsx
      ChatInput.tsx
      CodeBlock.tsx
    projects/
      ProjectList.tsx
      ProjectCard.tsx
    shared/
      ConnectionStatus.tsx
  hooks/
    useWebSocket.ts
    useChat.ts
    useProjects.ts
  lib/
    api.ts               # REST API client
    websocket.ts         # WebSocket handler
```

**API Contract (Frontend <-> Backend)**:
```typescript
// WebSocket Messages
interface ChatMessage {
  type: 'user_message' | 'assistant_message' | 'status_update' | 'error';
  projectId: string;
  content: string;
  timestamp: string;
  metadata?: {
    buildId?: string;
    progress?: number;
    stage?: string;
  };
}

// REST Endpoints
GET    /api/projects                 # List user projects
POST   /api/projects                 # Create new project
GET    /api/projects/:id             # Get project details
GET    /api/projects/:id/messages    # Get conversation history
GET    /api/projects/:id/builds      # Get build history
GET    /api/projects/:id/status      # Get current status
```

### 2.2 AI/LLM Integration (Claude API)

**Purpose**: Translate natural language requests into code, provide contextual responses, and guide users through application development.

**Technology**: Claude API (Anthropic)

**Why Claude**:
- Superior code generation quality compared to alternatives
- 200K token context window supports full conversation history
- Strong instruction-following for consistent output formats
- Tool use capability for structured code generation
- Reasonable pricing for MVP scale

**Integration Architecture**:
```
+-------------------+     +-------------------+     +-------------------+
|  Chat Orchestrator|---->|    AI Service     |---->|    Claude API     |
|                   |     |                   |     |                   |
|  - Conversation   |     |  - Prompt Mgmt    |     |  - Code Gen       |
|    Context        |     |  - Token Tracking |     |  - NL Response    |
|  - Project State  |     |  - Response Parse |     |  - Error Explain  |
|  - User Intent    |     |  - Error Handling |     |                   |
+-------------------+     +-------------------+     +-------------------+
```

**Prompt Strategy**:
```
System Prompt Layers:
1. Base System Prompt
   - Role definition (helpful coding assistant)
   - Output format requirements (structured JSON for code blocks)
   - Safety guidelines

2. Project Context
   - Current project type and tech stack
   - Existing files and structure
   - Previous decisions made

3. Conversation History
   - Recent messages (sliding window)
   - Key decisions and clarifications

4. Current Request
   - User's latest message
   - Any detected intent hints
```

**Response Parsing**:
```go
type AIResponse struct {
    Explanation  string        `json:"explanation"`
    CodeBlocks   []CodeBlock   `json:"code_blocks"`
    Questions    []string      `json:"questions,omitempty"`
    NextSteps    []string      `json:"next_steps,omitempty"`
    BuildTrigger bool          `json:"build_trigger"`
}

type CodeBlock struct {
    Filename    string `json:"filename"`
    Language    string `json:"language"`
    Content     string `json:"content"`
    Action      string `json:"action"` // create, update, delete
    Explanation string `json:"explanation"`
}
```

**Token Budget Management**:
- Track token usage per conversation
- Implement sliding window for history (keep last N messages)
- Summarize older context when approaching limits
- User-facing indicator of conversation length

### 2.3 Devcontainer Orchestration

**Purpose**: Provide isolated, reproducible development environments for each project.

**Technology**: Docker/Podman with Devcontainer specification

**Why Devcontainers**:
- Industry-standard specification
- IDE integration (VSCode, JetBrains)
- Reproducible environments
- Pre-built images available
- Feature extensibility

**Architecture**:
```
+---------------------+     +---------------------+     +---------------------+
|  Devcontainer       |     |  Container          |     |  Docker/Podman      |
|  Service            |---->|  Manager            |---->|  Daemon             |
|                     |     |                     |     |                     |
|  - Create env       |     |  - Image pull/build |     |  - Container run    |
|  - Start/stop       |     |  - Volume mgmt      |     |  - Network config   |
|  - Port forwarding  |     |  - Cleanup          |     |  - Resource limits  |
+---------------------+     +---------------------+     +---------------------+
           |
           v
+---------------------+
|  Devcontainer       |
|  Templates          |
|                     |
|  - Node.js/React    |
|  - Python/Flask     |
|  - Go               |
|  - Static Site      |
+---------------------+
```

**Container Lifecycle**:
```
1. Project Creation
   - Select template based on detected app type
   - Generate devcontainer.json
   - Create Docker network for project

2. Development Session
   - Start container on demand
   - Mount project volume
   - Forward ports for preview
   - Stream logs to frontend

3. Build Process
   - Execute build commands in container
   - Capture output for display
   - Extract artifacts for deployment

4. Cleanup
   - Stop containers after inactivity (30 min default)
   - Preserve volumes for state
   - Clean orphaned containers daily
```

**Resource Limits (MVP)**:
```yaml
# Per-container limits
resources:
  memory: 2GB
  cpu: 1.0
  storage: 5GB

# Per-user limits (MVP)
max_containers: 3
max_projects: 10
```

**Devcontainer Template Example**:
```json
{
  "name": "Node.js Project",
  "image": "mcr.microsoft.com/devcontainers/javascript-node:18",
  "features": {
    "ghcr.io/devcontainers/features/git:1": {}
  },
  "forwardPorts": [3000],
  "postCreateCommand": "npm install",
  "customizations": {
    "vscode": {
      "extensions": ["dbaeumer.vscode-eslint"]
    }
  }
}
```

### 2.4 GitLab Integration

**Purpose**: Provide version control, CI/CD pipelines, and container registry for all projects.

**Technology**: GitLab CE (Community Edition), self-hosted

**Why Self-Hosted GitLab CE**:
- Full API access without rate limits
- Data sovereignty (user code stays on our infrastructure)
- Integrated CI/CD (no external service needed)
- Built-in container registry
- Free for unlimited users/repos
- Familiar to developers for debugging

**Integration Points**:
```
+---------------------+     +---------------------+     +---------------------+
|  GitLab Service     |     |  GitLab CE          |     |  GitLab Runner      |
|  (go-chat)          |---->|  (API)              |---->|  (Docker executor)  |
|                     |     |                     |     |                     |
|  - Create repos     |     |  - REST API         |     |  - CI jobs          |
|  - Push commits     |     |  - Webhooks         |     |  - Build images     |
|  - Trigger CI       |     |  - Git over HTTPS   |     |  - Run tests        |
|  - Get status       |     |  - Container Reg    |     |  - Deploy           |
+---------------------+     +---------------------+     +---------------------+
```

**Repository Structure** (Generated per project):
```
project-name/
  .devcontainer/
    devcontainer.json
    Dockerfile          # If custom base needed
  .gitlab-ci.yml        # Auto-generated CI config
  src/                  # Application source
  tests/                # Generated tests
  README.md             # Auto-generated docs
  package.json          # Or requirements.txt, go.mod, etc.
```

**CI Pipeline Template**:
```yaml
stages:
  - validate
  - build
  - test
  - deploy

variables:
  DOCKER_TLS_CERTDIR: ""

validate:
  stage: validate
  script:
    - echo "Validating project structure..."
    - test -f package.json || test -f requirements.txt || test -f go.mod

build:
  stage: build
  script:
    - echo "Building application..."
    # Template-specific build commands

test:
  stage: test
  script:
    - echo "Running tests..."
    # Template-specific test commands

deploy:
  stage: deploy
  only:
    - main
  script:
    - echo "Deploying to preview environment..."
    # Deploy to preview URL
```

**Webhook Flow**:
```
GitLab -> Webhook -> go-chat Backend -> Update Project Status -> Notify Frontend
```

### 2.5 Deployment Pipeline

**Purpose**: Automatically build and deploy user applications to accessible preview URLs.

**MVP Deployment Model**:
For MVP, we use a simplified single-server deployment model where all components run on one host.

```
+------------------------------------------------------------------+
|                      SINGLE HOST (MVP)                            |
|                                                                    |
|  +------------+  +------------+  +------------+  +------------+   |
|  | go-chat    |  | GitLab CE  |  | PostgreSQL |  | Caddy      |   |
|  | Backend    |  |            |  |            |  | (Proxy)    |   |
|  +------------+  +------------+  +------------+  +------------+   |
|                                                                    |
|  +---------------------------------------------------------+      |
|  |                User Application Containers               |      |
|  |  +----------+  +----------+  +----------+  +----------+ |      |
|  |  | proj-1   |  | proj-2   |  | proj-3   |  | proj-4   | |      |
|  |  | :3001    |  | :3002    |  | :3003    |  | :3004    | |      |
|  |  +----------+  +----------+  +----------+  +----------+ |      |
|  +---------------------------------------------------------+      |
|                                                                    |
+------------------------------------------------------------------+
```

**Preview URL Scheme**:
```
https://<project-slug>.preview.gochat.local
https://my-app.preview.gochat.local
```

**Deployment Flow**:
```
1. User says "deploy this" or build succeeds
2. Backend triggers GitLab CI pipeline
3. CI builds container image, pushes to registry
4. Backend pulls image, starts container
5. Caddy route created for preview URL
6. User notified with preview link
```

---

## 3. Tech Stack Details

### 3.1 Frontend: Next.js

**Version**: Next.js 14+ with React 18+

**Dependencies**:
```json
{
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "typescript": "^5.0.0",
    "tailwindcss": "^3.4.0",
    "@types/react": "^18.2.0",
    "@types/node": "^20.0.0"
  }
}
```

**Styling**: Tailwind CSS for rapid, mobile-first development

**State Management**: React hooks (useState, useReducer, useContext)

**Build Output**: Node.js server for WebSocket support and SSR

### 3.2 Backend: Go

**Version**: Go 1.22+

**Key Libraries**:
```go
// go.mod dependencies
require (
    github.com/gin-gonic/gin v1.9.1       // HTTP router
    github.com/gorilla/websocket v1.5.1    // WebSocket support
    github.com/jmoiron/sqlx v1.3.5         // SQL extensions
    github.com/lib/pq v1.10.9              // PostgreSQL driver
    github.com/docker/docker v24.0.7       // Docker client
    github.com/xanzy/go-gitlab v0.95.0     // GitLab API client
    github.com/anthropics/anthropic-sdk-go // Claude API (or HTTP client)
    github.com/rs/zerolog v1.31.0          // Structured logging
    github.com/golang-migrate/migrate v4.16.2 // DB migrations
)
```

**Project Structure**:
```
cmd/
  server/
    main.go              # Entry point
internal/
  api/
    handlers/            # HTTP handlers
    middleware/          # Auth, logging, etc.
    routes.go            # Route definitions
  chat/
    orchestrator.go      # Conversation management
    intent.go            # Intent detection
  ai/
    client.go            # Claude API client
    prompts.go           # Prompt templates
    parser.go            # Response parsing
  devcontainer/
    manager.go           # Container lifecycle
    templates/           # Devcontainer templates
  gitlab/
    client.go            # GitLab API wrapper
    webhooks.go          # Webhook handlers
    pipelines.go         # CI/CD management
  storage/
    postgres.go          # Database operations
    models.go            # Data models
  config/
    config.go            # Configuration loading
migrations/
  *.sql                  # Database migrations
```

### 3.3 Database: PostgreSQL

**Version**: PostgreSQL 16

**Schema Design**:
```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    template VARCHAR(50) NOT NULL,
    gitlab_project_id INTEGER,
    gitlab_repo_url TEXT,
    preview_url TEXT,
    status VARCHAR(50) DEFAULT 'created',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Conversations table (stores full chat history)
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id),
    role VARCHAR(20) NOT NULL, -- 'user', 'assistant', 'system'
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    tokens_used INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Builds table
CREATE TABLE builds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id),
    gitlab_pipeline_id INTEGER,
    status VARCHAR(50) NOT NULL,
    stage VARCHAR(50),
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    logs TEXT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_projects_user ON projects(user_id);
CREATE INDEX idx_conversations_project ON conversations(project_id);
CREATE INDEX idx_conversations_created ON conversations(created_at);
CREATE INDEX idx_builds_project ON builds(project_id);
CREATE INDEX idx_builds_status ON builds(status);
```

### 3.4 Container Orchestration

**Technology**: Docker Engine with Docker Compose for MVP

**Why Docker (not Kubernetes for MVP)**:
- Simpler to operate and debug
- Single-node is sufficient for MVP scale
- Lower resource overhead
- Faster iteration on container management code
- Kubernetes can be added later when scaling requires it

**Docker Compose for Platform Services**:
```yaml
version: '3.8'

services:
  go-chat-backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@postgres:5432/gochat
      - GITLAB_URL=http://gitlab:80
      - CLAUDE_API_KEY=${CLAUDE_API_KEY}
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    depends_on:
      - postgres
      - gitlab

  go-chat-frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - API_URL=http://go-chat-backend:8080

  postgres:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=gochat
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=gochat
    volumes:
      - postgres_data:/var/lib/postgresql/data

  gitlab:
    image: gitlab/gitlab-ce:latest
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - gitlab_config:/etc/gitlab
      - gitlab_logs:/var/log/gitlab
      - gitlab_data:/var/opt/gitlab

  caddy:
    image: caddy:2-alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy_data:/data

volumes:
  postgres_data:
  gitlab_config:
  gitlab_logs:
  gitlab_data:
  caddy_data:
```

### 3.5 GitLab CE Integration Approach

**Setup Requirements**:
- GitLab CE instance (Docker or native install)
- Personal Access Token with api, read_repository, write_repository scopes
- GitLab Runner configured with Docker executor
- Webhook secret for pipeline notifications

**API Integration Pattern**:
```go
type GitLabService struct {
    client *gitlab.Client
    config GitLabConfig
}

// Key operations
func (g *GitLabService) CreateProject(name, description string) (*Project, error)
func (g *GitLabService) PushFiles(projectID int, files []File, message string) error
func (g *GitLabService) TriggerPipeline(projectID int, ref string) (*Pipeline, error)
func (g *GitLabService) GetPipelineStatus(projectID, pipelineID int) (*PipelineStatus, error)
func (g *GitLabService) SetupWebhook(projectID int, url string) error
```

---

## 4. MVP Architecture Decisions

### 4.1 Build vs Buy/Integrate Matrix

| Capability | Decision | Rationale |
|------------|----------|-----------|
| Chat UI | Build | Core differentiator, must match our UX vision |
| AI Code Gen | Integrate (Claude) | Best-in-class, no value in building our own LLM |
| Version Control | Integrate (GitLab) | Industry standard, excellent API |
| CI/CD | Integrate (GitLab CI) | Bundled with GitLab, reduces complexity |
| Container Registry | Integrate (GitLab) | Bundled with GitLab |
| Database | Use (PostgreSQL) | Industry standard, no custom needs |
| User Auth | Build (simple) | MVP needs basic auth only, simple to implement |
| Devcontainer Mgmt | Build | Custom needs for our workflow, thin layer over Docker |

### 4.2 Single-Tenant First

**Decision**: Build single-tenant architecture for MVP

**Rationale**:
- Faster to build and iterate
- Simpler security model
- Easier debugging and support
- Multi-tenancy complexity deferred until validated

**Single-Tenant Implications**:
- All users share the same database (but data is user-scoped)
- All projects run on same Docker host
- No resource isolation between users
- Simple flat configuration

**Multi-Tenancy Migration Path** (Month 6):
1. Add organization/tenant table
2. Scope all queries by tenant_id
3. Add resource quotas per tenant
4. Separate container networks per tenant
5. Consider separate GitLab groups per tenant

### 4.3 Key Interfaces

**Frontend <-> Backend**:
```
Protocol: HTTPS + WSS (WebSocket Secure)
Format: JSON
Auth: JWT tokens in Authorization header
Endpoints: RESTful for CRUD, WebSocket for chat streaming
```

**Backend <-> Claude API**:
```
Protocol: HTTPS
Format: JSON (Anthropic API format)
Auth: API key in header
Pattern: Request-response with streaming for long generations
```

**Backend <-> GitLab**:
```
Protocol: HTTPS
Format: JSON (GitLab API v4)
Auth: Personal Access Token
Pattern:
  - Outbound: REST API calls
  - Inbound: Webhooks for pipeline status
```

**Backend <-> Docker**:
```
Protocol: Unix socket (/var/run/docker.sock)
Format: Docker API
Auth: Socket access (host-level permission)
Pattern: API calls for container lifecycle
```

---

## 5. Risk Assessment

### 5.1 Technical Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Claude API rate limits/downtime | Medium | High | Implement retry logic, queue requests, consider fallback LLM |
| Container resource exhaustion | Medium | Medium | Strict resource limits, cleanup jobs, monitoring |
| GitLab integration complexity | Low | Medium | Use official client library, extensive testing |
| WebSocket connection stability | Low | Medium | Implement reconnection logic, message queue for offline |
| Database scaling | Low | Low | PostgreSQL handles MVP scale easily, optimize later |

### 5.2 Integration Complexity Points

**High Complexity**:
1. **AI Response Parsing**
   - Challenge: Claude outputs can be unpredictable
   - Mitigation: Strict prompt engineering, validation layer, fallback parsing

2. **Container Lifecycle Management**
   - Challenge: Many edge cases (crashes, hangs, resource leaks)
   - Mitigation: Health checks, timeouts, automatic cleanup

**Medium Complexity**:
1. **GitLab Pipeline Orchestration**
   - Challenge: Async webhook handling, status synchronization
   - Mitigation: Idempotent handlers, status polling fallback

2. **Real-time Chat Streaming**
   - Challenge: Coordinating AI streaming with UI updates
   - Mitigation: Clear message protocol, buffering strategy

**Low Complexity**:
1. **User Authentication**
   - Using proven patterns (JWT, bcrypt)
   - Can integrate OAuth later

2. **Database Operations**
   - Standard CRUD, well-understood patterns

### 5.3 Dependency Risks

| Dependency | Risk | Mitigation |
|------------|------|------------|
| Claude API | Pricing changes, API changes | Abstract AI interface, monitor alternatives |
| GitLab CE | Security vulnerabilities | Regular updates, pin versions |
| Docker | License changes (already happened) | Podman as drop-in replacement |

---

## 6. Deployment Architecture

### 6.1 MVP Deployment (Single Server)

**Recommended Specs**:
- CPU: 8 cores
- RAM: 32 GB
- Storage: 500 GB SSD
- Network: 1 Gbps

**Cost Estimate**: ~$200-400/month on major cloud providers

**Deployment Method**: Docker Compose on single VM

### 6.2 Configuration Management

**Environment Variables**:
```bash
# Core
GO_CHAT_ENV=production
GO_CHAT_PORT=8080
GO_CHAT_SECRET_KEY=<generated>

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/gochat

# Claude API
CLAUDE_API_KEY=<api-key>
CLAUDE_MODEL=claude-3-opus-20240229

# GitLab
GITLAB_URL=https://gitlab.gochat.local
GITLAB_TOKEN=<personal-access-token>
GITLAB_WEBHOOK_SECRET=<generated>

# Docker
DOCKER_HOST=unix:///var/run/docker.sock

# Features
ENABLE_SIGNUP=true
MAX_PROJECTS_PER_USER=10
CONTAINER_MEMORY_LIMIT=2g
```

### 6.3 Monitoring (MVP)

**Essential Metrics**:
- API response times (P50, P95, P99)
- WebSocket connection count
- Container count and resource usage
- GitLab pipeline success rate
- Claude API token usage and costs

**Tools**:
- Application logs: Structured JSON to stdout, aggregate with Loki or CloudWatch
- Metrics: Prometheus + Grafana (or cloud-native equivalents)
- Error tracking: Sentry (free tier sufficient for MVP)

---

## 7. Security Considerations

### 7.1 MVP Security Requirements

| Area | Requirement | Implementation |
|------|-------------|----------------|
| Transport | TLS everywhere | Caddy automatic HTTPS |
| Auth | Secure password storage | bcrypt with appropriate cost factor |
| Sessions | Secure token handling | JWT with short expiry, HTTP-only cookies |
| API Keys | No client exposure | All API calls server-side |
| Containers | Isolation | Separate Docker networks per project |
| Secrets | No hardcoding | Environment variables, secret manager later |

### 7.2 Container Security

```yaml
# Security settings for user containers
security_opt:
  - no-new-privileges:true
cap_drop:
  - ALL
cap_add:
  - CHOWN
  - SETUID
  - SETGID
read_only: false  # Need writable for dev
tmpfs:
  - /tmp
```

---

## 8. Development Workflow

### 8.1 Local Development Setup

```bash
# Prerequisites
- Docker Desktop or Podman
- Go 1.22+
- Node.js 20+
- pnpm (for frontend)

# Setup
git clone <repo>
cd go-chat
cp .env.example .env
docker compose -f docker-compose.dev.yml up -d postgres gitlab
make migrate
make run  # Starts backend and frontend in dev mode
```

### 8.2 Testing Strategy

| Layer | Tool | Coverage Target |
|-------|------|-----------------|
| Unit (Go) | go test | 70%+ for core logic |
| Unit (Frontend) | Vitest | 60%+ for components |
| Integration | go test + testcontainers | Key workflows |
| E2E | Playwright | Critical user paths |

---

## 9. Future Considerations

### 9.1 Scaling Path

When MVP validation succeeds, consider:

1. **Horizontal Scaling** (Month 6-9)
   - Stateless backend behind load balancer
   - Redis for session/cache
   - Separate container hosts

2. **Multi-Region** (Month 12+)
   - Database replication
   - CDN for frontend
   - Regional container hosts

3. **Kubernetes Migration** (When justified)
   - Only when managing 100+ containers becomes painful
   - Or when multi-tenant isolation requires it

### 9.2 Feature Architecture Hooks

The MVP architecture accommodates future features:

- **Templates Marketplace**: Add templates table, template registry
- **Team Collaboration**: Add teams table, project permissions
- **Custom Domains**: Caddy dynamic configuration
- **Billing**: Add usage metering to existing tables

---

## 10. Architecture Decision Records

### ADR-001: Monolithic Backend Over Microservices

**Status**: Accepted

**Context**: Need to decide on backend architecture pattern for MVP.

**Decision**: Build a modular monolith in Go rather than microservices.

**Consequences**:
- Positive: Faster development, simpler deployment, easier debugging
- Positive: Can extract services later if needed
- Negative: Must maintain discipline with module boundaries
- Negative: Single point of failure (acceptable for MVP)

### ADR-002: Next.js for Frontend

**Status**: Accepted

**Context**: Need to choose frontend framework optimized for mobile-first chat interface.

**Decision**: Use Next.js 14+ with App Router instead of SvelteKit.

**Consequences**:
- Positive: React ecosystem with massive library support
- Positive: Already configured in project devcontainer
- Positive: Excellent streaming support via App Router
- Positive: Large developer community
- Negative: Larger bundle size than Svelte (mitigated by code splitting)

### ADR-003: Self-Hosted GitLab Over GitHub/Cloud GitLab

**Status**: Accepted

**Context**: Need version control and CI/CD for user projects.

**Decision**: Self-host GitLab CE rather than use GitHub or GitLab SaaS.

**Consequences**:
- Positive: Full API access, no rate limits
- Positive: Data sovereignty, user code on our infrastructure
- Positive: Free for unlimited users (CE edition)
- Negative: Operational burden of running GitLab
- Negative: Must handle security updates ourselves

### ADR-004: Docker Over Kubernetes for MVP

**Status**: Accepted

**Context**: Need container orchestration for devcontainers and user applications.

**Decision**: Use Docker Engine with Docker Compose, not Kubernetes.

**Consequences**:
- Positive: Simpler to operate and debug
- Positive: Lower resource overhead
- Positive: Faster development iteration
- Negative: Manual scaling
- Negative: Will need migration if we outgrow single node

---

## Appendix A: Technology Evaluation Matrix

### Frontend Framework Comparison (Updated)

| Criteria | Next.js | SvelteKit | Nuxt | Weight |
|----------|---------|-----------|------|--------|
| Bundle Size | 3 | 5 | 3 | Medium |
| Mobile Performance | 4 | 5 | 4 | High |
| Developer Experience | 5 | 5 | 4 | Medium |
| Ecosystem | 5 | 3 | 4 | High |
| Existing Setup | 5 | 1 | 1 | High |
| **Weighted Score** | **4.5** | 3.6 | 3.0 | |

**Note**: Next.js selected due to existing devcontainer setup and React ecosystem benefits.

### Backend Language Comparison

| Criteria | Go | Node.js | Python | Rust | Weight |
|----------|-----|---------|--------|------|--------|
| Performance | 5 | 3 | 2 | 5 | Medium |
| Docker Integration | 5 | 4 | 4 | 4 | High |
| Development Speed | 4 | 5 | 5 | 2 | High |
| Concurrency Model | 5 | 3 | 2 | 5 | Medium |
| Deployment Simplicity | 5 | 3 | 3 | 4 | High |
| **Weighted Score** | **4.7** | 3.5 | 3.2 | 3.8 | |

---

## Document History

| Date | Version | Author | Changes |
|------|---------|--------|---------|
| 2025-12-24 | 1.0 | Architect Agent | Initial MVP architecture |
