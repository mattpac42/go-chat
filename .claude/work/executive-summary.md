# Executive Summary: Devcontainer-Per-Project Evaluation

**Date:** 2026-01-01
**Recommendation:** DO NOT implement devcontainer-per-project for local tool

## Quick Decision Matrix

| Factor | Current | Containers | Workspaces |
|--------|---------|-----------|-----------|
| Dev Time | - | 8 weeks | 1 week |
| Complexity | Low | High | Low |
| User Resources | 500MB | 2-10GB | 1GB |
| Cold Start | 0s | 10-30s | 0s |
| Isolation | None | Strong | Medium |
| Maintenance | 2h/mo | 8h/mo | 2h/mo |
| Best For | Current state | SaaS | Local tool |

## The Answer to Each Question

### 1. Container Orchestration
**Recommended:** Docker Compose + Docker SDK (if you must use containers)
- Kubernetes is overkill for local tool
- Podman adds user friction
- Simple dynamic container creation with Docker API

**Better Alternative:** None needed for workspace approach

### 2. Resource Requirements
**Per Container:** 1-2GB RAM, 1-2 CPUs
**Realistic Limits:** 5-8 concurrent projects on typical dev machine
**Cleanup Required:** Idle timeout, LRU eviction, regular cleanup

**Workspace Approach:** 100-200MB per project, 20+ concurrent projects

### 3. Lifecycle Management
**When to create:** Lazy (on first message)
**When to stop:** 60min idle timeout
**When to remove:** 7 days after stop
**Persistence:** Host filesystem volumes

**Workspace Approach:** Create on first use, clean up after 7 days inactive

### 4. Communication
**Phase 1 (MVP):** Docker exec (works with existing Claude Code CLI)
**Phase 2 (SaaS):** HTTP REST API (requires Claude Code modification)

**Workspace Approach:** Direct CLI execution in project directory

### 5. Security
**Container Isolation:** seccomp, capabilities, read-only root, network policies
**Credentials:** Proxy pattern (orchestrator holds API key)
**Filesystem:** Bind mounts to project directories only

**Workspace Approach:** Process isolation, OS-level permissions

### 6. Cost/Complexity
**Development:** 8 weeks vs 1 week
**Code:** +3600 LOC (220% increase)
**Maintenance:** 4x ongoing effort
**User Resources:** 4-20x higher

**ROI:** Not worth it for local tool

### 7. Alternatives (THIS IS THE ANSWER)

**Recommended: Workspace Isolation**
```
/workspaces/
├── project-1/  → Execute: cd project-1 && claude code
├── project-2/  → Execute: cd project-2 && claude code
└── project-3/  → Execute: cd project-3 && claude code
```

**Why:**
- 80% of benefits, 20% of complexity
- No container overhead
- Instant response (no cold start)
- Easy debugging
- Low resources
- Simple implementation

**Limitations:**
- Shared host environment
- Less isolation
- Can't customize per-project environment

**When these limitations matter:** Upgrade to containers

## Recommendation

### For Current State (Local Development Tool)

**Implement Workspace Isolation:**
1. Week 1: Core implementation (WorkspaceManager + ClaudeExecutor)
2. Week 2: Polish (cleanup, monitoring, docs)
3. Total: 2 weeks, ~800 LOC, low complexity

**Skip container-per-project entirely.**

### For Future SaaS Pivot

**Then reconsider containers:**
1. Phase 1: Prove workspace approach works
2. Phase 2: Add optional container mode (user opt-in)
3. Phase 3: Default to containers when multi-tenant

### Hybrid Approach (Best of Both)

```yaml
# Per-project configuration
{
  "project_id": "abc123",
  "isolation_mode": "workspace",  // or "container"
}
```

- **Default:** workspace (fast, simple)
- **Optional:** container (when user needs it)
- **Future-proof:** Can scale to SaaS later

## Implementation Artifacts

Created three documents:

1. **devcontainer-per-project-evaluation.md**
   - Full 7-section analysis
   - Container orchestration options
   - Resource requirements
   - Lifecycle management
   - Communication patterns
   - Security considerations
   - Complete alternatives analysis

2. **workspace-isolation-implementation.md**
   - Full Go implementation
   - WorkspaceManager code
   - ClaudeExecutor code
   - Service integration
   - Docker Compose updates
   - Testing strategy
   - Migration roadmap

3. **executive-summary.md** (this document)
   - Quick decision matrix
   - Clear recommendations
   - Implementation plan

## Next Steps

### If You Agree (Recommended)

1. Review workspace-isolation-implementation.md
2. Start Phase 1 implementation (1 week)
3. Test with real projects
4. Deploy and validate

### If You Disagree

1. Review specific concerns
2. Identify which benefits of containers are critical
3. Consider hybrid approach (workspace + optional containers)
4. Evaluate phased migration

## Bottom Line

**Containers are over-engineering for a local development tool.**

Use workspace isolation:
- Simpler
- Faster
- Cheaper
- Good enough
- Upgradeable later if needed

The 80/20 rule: Get 80% of the value for 20% of the cost.
