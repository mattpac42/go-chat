# Devcontainer-Per-Project Infrastructure Evaluation

**Date:** 2026-01-01
**Context:** Evaluating infrastructure implications of running a devcontainer per project for Go Chat AI coding assistant

## Current Architecture

```
┌──────────────┐      WebSocket      ┌──────────────┐
│  Next.js     │◄──────────────────►│  Go Backend  │
│  Frontend    │                     │  (Port 8080) │
└──────────────┘                     └──────┬───────┘
                                            │
                                     ┌──────▼───────┐
                                     │  PostgreSQL  │
                                     │  (Port 5432) │
                                     └──────────────┘
                                            │
                                     Files stored in DB
                                     Claude API calls from backend
```

**Characteristics:**
- Single monolithic backend
- Files stored in PostgreSQL
- Backend makes Claude API calls directly
- No isolation between projects
- Simple Docker Compose setup

## Proposed Architecture

```
┌──────────────┐      WebSocket      ┌──────────────┐
│  Next.js     │◄──────────────────►│  Go Backend  │
│  Frontend    │                     │  (Orchestrator)│
└──────────────┘                     └──────┬───────┘
                                            │
                      ┌─────────────────────┼─────────────────────┐
                      │                     │                     │
                ┌─────▼──────┐       ┌─────▼──────┐       ┌─────▼──────┐
                │ Devcontainer│       │ Devcontainer│       │ Devcontainer│
                │  Project A  │       │  Project B  │       │  Project C  │
                │             │       │             │       │             │
                │ Claude Code │       │ Claude Code │       │ Claude Code │
                │  Running    │       │  Running    │       │  Running    │
                └─────────────┘       └─────────────┘       └─────────────┘
                      │                     │                     │
                      └─────────────────────┼─────────────────────┘
                                            │
                                     ┌──────▼───────┐
                                     │  PostgreSQL  │
                                     │  (Metadata)  │
                                     └──────────────┘
```

**Characteristics:**
- Backend becomes orchestrator/API gateway
- One devcontainer per project
- Claude Code CLI runs inside each container
- File operations happen in container filesystem
- Need container lifecycle management

---

## Evaluation Criteria

### 1. Container Orchestration

#### Option A: Docker Compose (Recommended for Local Tool)

**Pros:**
- Already in use (docker-compose.yml exists)
- Simple local development workflow
- Easy to manage on single machine
- No cluster complexity
- Native Docker socket access
- Familiar to developers

**Cons:**
- Not designed for dynamic container creation
- Limited scaling capabilities
- Manual container lifecycle management
- No built-in service discovery beyond DNS

**Implementation:**
```yaml
# Dynamic compose generation approach
services:
  backend:
    # orchestrator

  db:
    # metadata storage

  # Generated dynamically per project:
  project-abc123:
    build:
      context: .devcontainer
    volumes:
      - ./projects/abc123:/workspace
    environment:
      - CLAUDE_API_KEY=${CLAUDE_API_KEY}
    networks:
      - gochat
```

**Verdict:** Good fit for local tool. Use Docker SDK to dynamically create/remove containers.

#### Option B: Kubernetes

**Pros:**
- Designed for container orchestration
- Built-in service discovery
- Auto-scaling capabilities
- Resource limits and quotas
- Health checks and restarts
- Production-ready

**Cons:**
- Massive overkill for local development
- Complex setup (minikube, k3s, or kind)
- Higher resource overhead
- Steeper learning curve
- Not typical for "local tool"

**Verdict:** NOT recommended unless planning SaaS pivot.

#### Option C: Podman + systemd

**Pros:**
- Rootless containers (better security)
- systemd integration for lifecycle
- Compatible with Docker compose format
- No daemon requirement
- Native cgroup support

**Cons:**
- Less common than Docker
- Network complexity without daemon
- Limited tooling ecosystem
- User needs Podman installed

**Verdict:** Interesting for security, but adds user friction.

**RECOMMENDATION:** Docker Compose + Docker SDK for dynamic management

---

### 2. Resource Requirements

#### Per-Container Resources

**Base Devcontainer:**
- Base image: ~200-500MB (Node.js devcontainer)
- Claude Code CLI: ~50MB
- Project files: 10-500MB (varies)
- Runtime memory: 256-512MB idle, 1-2GB active
- CPU: 0.5-2 cores during active work

**Estimated Requirements:**

| Concurrent Projects | Memory | CPU | Disk |
|-------------------|--------|-----|------|
| 1 project | 2GB | 2 cores | 1GB |
| 5 projects | 6GB | 4 cores | 5GB |
| 10 projects | 10GB | 6 cores | 10GB |
| 20 projects | 18GB | 8 cores | 20GB |

**Typical Developer Machine:**
- MacBook Pro M1/M2: 16GB RAM, 8 cores → ~5-8 concurrent projects
- Linux workstation: 32GB RAM, 16 cores → ~15-20 concurrent projects
- Budget laptop: 8GB RAM, 4 cores → ~2-3 concurrent projects

#### Resource Limits

```yaml
# Per-container limits
services:
  project-abc123:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G
        reservations:
          cpus: '0.5'
          memory: 512M
```

**Cleanup Strategy:**

1. **Idle Timeout:** Stop containers after 30-60 minutes of inactivity
2. **Max Containers:** Limit to 10 concurrent containers per user
3. **LRU Eviction:** Stop least recently used when hitting limit
4. **Disk Cleanup:** Remove stopped containers after 7 days
5. **Image Pruning:** Weekly cleanup of unused images

**RECOMMENDATION:**
- Default limit: 5 concurrent containers
- Idle timeout: 1 hour
- Stop (don't remove) on idle
- Persist workspace volumes

---

### 3. Lifecycle Management

#### Container States

```
Created → Starting → Running → Active → Idle → Stopped → Removed
    ↑                                                        │
    └────────────────────────────────────────────────────────┘
```

#### State Transitions

**When to CREATE:**
- User opens project for first time
- User explicitly requests new environment
- Container was removed and needs recreation

**When to START:**
- User sends message to project
- User opens project (if stopped)
- Health check detects container down

**When to STOP:**
- Idle timeout reached (60 minutes)
- User closes project
- Resource limit reached (LRU)
- User explicitly stops

**When to REMOVE:**
- User deletes project
- Container stopped for >7 days
- Disk space critical
- User explicitly removes

#### Persistence Strategy

```
/workspace/projects/
├── abc123/                    # Project ID
│   ├── workspace/             # Mounted into container
│   │   ├── .git/
│   │   ├── src/
│   │   └── ...
│   ├── .container/            # Container metadata
│   │   ├── state.json         # Current state
│   │   ├── last_active        # Timestamp
│   │   └── container_id       # Docker container ID
│   └── logs/                  # Container logs
└── def456/
    └── ...
```

**Database Schema:**

```sql
CREATE TABLE project_containers (
    project_id UUID PRIMARY KEY,
    container_id VARCHAR(64),
    state VARCHAR(20), -- created, running, stopped, removed
    last_active TIMESTAMP,
    created_at TIMESTAMP,
    stopped_at TIMESTAMP,
    resource_limits JSONB,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);
```

**RECOMMENDATION:**
- Lazy creation (on first message)
- Persist workspace on host filesystem
- Track state in database
- Implement health checks

---

### 4. Communication Architecture

#### Option A: HTTP REST (Recommended)

**Architecture:**
```
Frontend → Backend (Orchestrator) → HTTP → Container (Port 8080)
                                           Claude Code REST API
```

**Implementation:**
```go
// Backend orchestrator
func (h *ProjectHandler) SendMessage(c *gin.Context) {
    projectID := c.Param("id")

    // Ensure container is running
    container := h.containerMgr.EnsureRunning(projectID)

    // Forward request to container
    resp, err := http.Post(
        fmt.Sprintf("http://%s:8080/api/chat", container.IP),
        "application/json",
        c.Request.Body,
    )

    // Stream response back
    io.Copy(c.Writer, resp.Body)
}
```

**Pros:**
- Simple, well-understood protocol
- Easy to debug (curl, browser)
- No special container networking
- Can use existing HTTP libraries
- Stateless (scales well)

**Cons:**
- Need REST API in Claude Code (doesn't exist)
- Each request is separate (no session)
- Overhead for many small requests

#### Option B: WebSocket (Bidirectional)

**Architecture:**
```
Frontend → WS → Backend (Orchestrator) → WS → Container
                                              Claude Code WS Server
```

**Pros:**
- Bidirectional communication
- Lower latency
- Streaming responses natural
- Connection-oriented (session state)

**Cons:**
- More complex connection management
- Need WS server in Claude Code
- Connection drops on container restart
- More complex debugging

#### Option C: Unix Socket (Container-only)

**Architecture:**
```
Backend (Orchestrator) → Docker Exec → Container
                         Claude Code stdin/stdout
```

**Implementation:**
```go
func (c *ContainerManager) ExecuteCommand(projectID, command string) (string, error) {
    execConfig := types.ExecConfig{
        Cmd:          []string{"claude", "code", "--json"},
        AttachStdout: true,
        AttachStdin:  true,
    }

    exec, err := c.client.ContainerExecCreate(ctx, containerID, execConfig)
    // ...
}
```

**Pros:**
- No network exposure
- Direct Docker API
- No modifications to Claude Code
- Works with existing CLI

**Cons:**
- Exec overhead per command
- No streaming (without complexity)
- Harder to debug
- Not as clean as REST/WS

#### Option D: gRPC (High Performance)

**Pros:**
- High performance
- Built-in streaming
- Type-safe protocol
- Good for microservices

**Cons:**
- Most complex implementation
- Overkill for local tool
- Requires protobuf definitions
- Needs gRPC server in Claude Code

**RECOMMENDATION:**

Start with **Docker Exec (Option C)** because:
1. No modifications to Claude Code needed
2. Works with existing CLI
3. Simplest implementation
4. Good for MVP/local tool

Migrate to **HTTP REST** if/when:
1. Planning SaaS offering
2. Need better performance
3. Want cleaner architecture
4. Claude Code adds REST API

---

### 5. Security

#### Isolation Layers

**1. Container Isolation**
```yaml
services:
  project-abc123:
    security_opt:
      - no-new-privileges:true
      - seccomp=unconfined  # Claude Code needs broad access
    cap_drop:
      - ALL
    cap_add:
      - CHOWN
      - DAC_OVERRIDE
      - SETGID
      - SETUID
```

**2. Network Isolation**
```yaml
networks:
  gochat-internal:
    driver: bridge
    internal: true  # No internet access

  gochat-external:
    driver: bridge  # Internet for Claude API

services:
  project-abc123:
    networks:
      - gochat-internal  # Talk to orchestrator
      - gochat-external  # Talk to Claude API
```

**3. Filesystem Isolation**
```yaml
services:
  project-abc123:
    volumes:
      - ./projects/abc123:/workspace:rw  # Project files
      - /tmp                              # Ephemeral
    read_only: true  # Root FS read-only
    tmpfs:
      - /tmp:rw,noexec,nosuid,size=1g
```

#### Credential Management

**Problem:** Each container needs Claude API key

**Option A: Environment Variable (Simple)**
```yaml
environment:
  - CLAUDE_API_KEY=${CLAUDE_API_KEY}
```
- Simple, works immediately
- Key visible in `docker inspect`
- Not secure for multi-user

**Option B: Docker Secrets (Better)**
```yaml
secrets:
  claude_api_key:
    external: true

services:
  project-abc123:
    secrets:
      - claude_api_key
```
- More secure
- Requires Docker Swarm mode (overkill for local)
- Or external secret file

**Option C: Secrets Manager (SaaS)**
- HashiCorp Vault
- AWS Secrets Manager
- Too complex for local tool

**Option D: Proxy Pattern (Recommended for Local)**
```
Container → Backend (Orchestrator) → Claude API
            (Orchestrator injects key)
```
- Container never sees API key
- Orchestrator manages credentials
- Can add rate limiting, quotas
- Better security model

**RECOMMENDATION:**
- **Local tool:** Proxy pattern (orchestrator calls Claude API)
- **SaaS:** Docker secrets + secrets manager

#### Additional Security Measures

1. **Resource Limits:** Prevent DoS via CPU/memory exhaustion
2. **Network Policies:** Restrict container-to-container communication
3. **Image Scanning:** Scan devcontainer images for vulnerabilities
4. **User Namespaces:** Map container root to non-root host user
5. **Audit Logging:** Log all container creation/deletion events

---

### 6. Cost/Complexity Analysis

#### Development Complexity

| Component | Current | Proposed | Delta |
|-----------|---------|----------|-------|
| Backend code | Simple HTTP server | Container orchestration logic | +1500 LOC |
| Container management | None | Create/start/stop/remove | +800 LOC |
| Networking | Single backend | Per-project routing | +400 LOC |
| Monitoring | Basic health check | Per-container health | +300 LOC |
| Testing | Unit + integration | + Container tests | +600 LOC |
| **Total** | **~3000 LOC** | **~6600 LOC** | **+3600 LOC (2.2x)** |

#### Operational Complexity

**Current:**
```bash
docker-compose up
# Done!
```

**Proposed:**
```bash
docker-compose up                    # Start orchestrator
# Wait for project activity...
# Auto-create container per project
# Manage lifecycle
# Monitor resources
# Clean up idle containers
# Handle failures
```

**New failure modes:**
- Container creation fails
- Container crashes during operation
- Resource exhaustion (disk/memory/CPU)
- Network communication failures
- Docker daemon issues
- Volume mount permission problems

**Debugging complexity:**
- Which container has the issue?
- Container logs scattered
- Network routing problems
- Resource contention
- State synchronization

#### User Experience

**Current:**
1. Open browser
2. Create project
3. Chat
4. Done

**Proposed:**
1. Open browser
2. Create project
3. Send first message → **wait for container creation (10-30s)**
4. Chat
5. Project idle → container stops → **next message slower**
6. Many projects → **resource warnings**
7. Disk full → **cleanup needed**

**UX Concerns:**
- First message latency (cold start)
- Unexpected slowdowns (container stopped)
- Resource limit warnings
- "Too many projects" errors
- Confusing error messages

#### Cost Comparison

**Local Tool (Current):**
- Development: 2 weeks initial
- Maintenance: ~2 hours/month
- User resources: 500MB RAM, 1 CPU
- No infrastructure costs

**Local Tool (Proposed):**
- Development: 6-8 weeks initial
- Maintenance: ~8 hours/month
- User resources: 2-10GB RAM, 2-8 CPUs
- No infrastructure costs
- **Trade-off:** 3x dev time, 4x maintenance, 4-20x resources

**SaaS (Current):**
- Infrastructure: ~$50/month (basic VPS)
- Scaling: Vertical (bigger server)

**SaaS (Proposed):**
- Infrastructure: ~$200-500/month (cluster)
- Scaling: Horizontal (more nodes)
- **Trade-off:** 4-10x infrastructure cost

---

### 7. Alternatives Analysis

#### Alternative 1: Workspace Isolation (No Containers)

```
Backend
├── /workspaces/abc123/
├── /workspaces/def456/
└── /workspaces/ghi789/
```

**Approach:**
- Each project gets isolated directory
- Backend executes Claude Code CLI in project directory
- Use process isolation instead of containers

**Pros:**
- Much simpler (no Docker complexity)
- Lower resource overhead
- Faster execution (no container startup)
- Easier debugging

**Cons:**
- Less isolation (shared filesystem)
- Security concerns (process escaping)
- No environment customization
- Dependencies conflict

**Code Example:**
```go
func (s *ChatService) ExecuteClaude(projectID, message string) (string, error) {
    workspace := filepath.Join("/workspaces", projectID)

    cmd := exec.Command("claude", "code", "--workspace", workspace)
    cmd.Stdin = strings.NewReader(message)

    output, err := cmd.CombinedOutput()
    return string(output), err
}
```

#### Alternative 2: Shared Devcontainer with Git Branches

**Approach:**
- One devcontainer
- Each project is a git branch
- Switch branches for different projects

**Pros:**
- Single container (simple)
- Git handles isolation
- Low resources

**Cons:**
- Can't work on multiple projects simultaneously
- Branch switching overhead
- State contamination
- Terrible UX

#### Alternative 3: VM per Project (Firecracker)

**Approach:**
- Use Firecracker microVMs
- Full OS isolation
- Similar to AWS Lambda

**Pros:**
- Better isolation than containers
- Fast startup (~125ms)
- Lower memory overhead than VMs

**Cons:**
- Linux only (KVM required)
- Complex setup
- Overkill for local tool
- Not macOS/Windows compatible

#### Alternative 4: WebAssembly Sandbox

**Approach:**
- Compile Claude Code to WASM
- Run in sandboxed WASM runtime
- Ultimate isolation

**Pros:**
- Perfect sandboxing
- Cross-platform
- Very low overhead
- Future-proof

**Cons:**
- Claude Code not available as WASM
- Limited filesystem access
- Immature tooling
- Major development effort

#### Alternative 5: Hybrid Approach (Recommended)

**Approach:**
- Default: Workspace isolation (Alternative 1)
- Optional: Devcontainer mode (power users)
- Configurable per-project

**Pros:**
- Simple default path
- Advanced option available
- User choice
- Gradual migration

**Cons:**
- Two codepaths to maintain
- More testing needed

**Implementation:**
```go
type ProjectConfig struct {
    ID              string
    IsolationMode   string  // "workspace" or "container"
}

func (s *ChatService) ExecuteClaude(project *ProjectConfig, msg string) (string, error) {
    switch project.IsolationMode {
    case "container":
        return s.containerExec(project, msg)
    case "workspace":
        return s.workspaceExec(project, msg)
    default:
        return s.workspaceExec(project, msg)
    }
}
```

---

## Recommendations

### For Local Development Tool (Current State)

**DO NOT implement devcontainer-per-project.**

**Instead:**

1. **Use Workspace Isolation (Alternative 1)**
   - Each project gets isolated directory
   - Execute Claude Code CLI in project context
   - Simple, fast, low overhead
   - Good enough for local tool

2. **Add Optional Container Mode**
   - Per-project setting: "Use container"
   - For projects that need specific environments
   - User opts in (not default)

3. **Focus on Core Value**
   - Make chat experience better
   - Improve Claude Code integration
   - Better file management
   - Not infrastructure complexity

### If Planning SaaS Pivot

**Then consider devcontainer-per-project:**

1. **Phase 1: MVP with Docker Exec**
   - Docker Compose for orchestration
   - Docker SDK for dynamic containers
   - Docker exec for communication
   - Simple, validates concept

2. **Phase 2: REST API**
   - Add REST server to Claude Code
   - HTTP communication
   - Better architecture
   - Prepare for scale

3. **Phase 3: Kubernetes**
   - When scaling beyond single node
   - Multi-tenant infrastructure
   - Auto-scaling, self-healing
   - Production-ready

### Hybrid Recommendation (Best of Both)

```yaml
# docker-compose.yml
services:
  backend:
    # Orchestrator
    environment:
      - DEFAULT_ISOLATION=workspace
      - ALLOW_CONTAINER_MODE=true
      - MAX_CONTAINERS=5

  # Containers created dynamically when requested
```

**Benefits:**
- Low complexity default (workspace)
- High isolation option (container)
- User choice
- Path to SaaS if needed

---

## Implementation Roadmap

### If Proceeding with Containers (Against Recommendation)

**Week 1-2: Foundation**
- [ ] Docker SDK integration
- [ ] Container lifecycle manager
- [ ] Basic create/start/stop/remove
- [ ] Workspace volume mounting

**Week 3-4: Communication**
- [ ] Docker exec integration
- [ ] Message routing
- [ ] Error handling
- [ ] Response streaming

**Week 5-6: Lifecycle**
- [ ] Idle timeout detection
- [ ] Health checks
- [ ] Auto-restart on failure
- [ ] Graceful shutdown

**Week 7-8: Polish**
- [ ] Resource limits
- [ ] Cleanup strategies
- [ ] Monitoring/logging
- [ ] User documentation

### Recommended Alternative (Workspace Isolation)

**Week 1: Implementation**
- [ ] Workspace directory structure
- [ ] Claude Code CLI integration
- [ ] Process execution
- [ ] File management

**Week 2: Polish**
- [ ] Error handling
- [ ] Testing
- [ ] Documentation
- [ ] Optional container mode flag

---

## Conclusion

**For a local development tool, devcontainer-per-project is over-engineering.**

**Recommendation: Start simple with workspace isolation.**

**Key Decision Factors:**

| Factor | Workspace | Container |
|--------|-----------|-----------|
| Development time | 1 week | 8 weeks |
| Maintenance burden | Low | High |
| User resources | Low | High |
| Isolation strength | Medium | High |
| Complexity | Low | High |
| UX friction | None | High |
| Future-proof | Sufficient | Over-prepared |

**The 80/20 rule applies:**
- 80% of value from workspace isolation
- 20% more value from containers
- 400% more complexity

**Start with workspace isolation. Add container mode later if proven necessary.**
