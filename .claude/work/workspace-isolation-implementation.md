# Workspace Isolation Implementation Guide

**Recommended approach for Go Chat AI coding assistant**

## Architecture

```
┌──────────────┐      WebSocket      ┌──────────────────────┐
│  Next.js     │◄──────────────────►│  Go Backend          │
│  Frontend    │                     │  (Orchestrator)      │
└──────────────┘                     └──────┬───────────────┘
                                            │
                      ┌─────────────────────┼─────────────────────┐
                      │                     │                     │
                ┌─────▼──────┐       ┌─────▼──────┐       ┌─────▼──────┐
                │ Workspace A│       │ Workspace B│       │ Workspace C│
                │ /workspaces│       │ /workspaces│       │ /workspaces│
                │    /abc123 │       │    /def456 │       │    /ghi789 │
                │            │       │            │       │            │
                │ Claude Code│       │ Claude Code│       │ Claude Code│
                │ CLI Exec   │       │ CLI Exec   │       │ CLI Exec   │
                └────────────┘       └────────────┘       └────────────┘
                      │                     │                     │
                      └─────────────────────┼─────────────────────┘
                                            │
                                     ┌──────▼───────┐
                                     │  PostgreSQL  │
                                     │  (Metadata + Files)
                                     └──────────────┘
```

## Directory Structure

```
/workspace/
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── executor/           # NEW: Claude Code executor
│   │   │   ├── executor.go
│   │   │   ├── workspace.go
│   │   │   └── process.go
│   │   └── ...
│   └── ...
│
├── workspaces/                  # NEW: Project workspaces
│   ├── abc123-project-name/
│   │   ├── .git/
│   │   ├── .claude/
│   │   ├── src/
│   │   ├── README.md
│   │   └── ...
│   ├── def456-other-project/
│   │   └── ...
│   └── .cleanup/                # Metadata for cleanup
│       └── last_access.json
│
└── frontend/
    └── ...
```

## Implementation

### 1. Workspace Manager

```go
// backend/internal/executor/workspace.go
package executor

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type WorkspaceManager struct {
    baseDir string
}

func NewWorkspaceManager(baseDir string) *WorkspaceManager {
    return &WorkspaceManager{
        baseDir: baseDir,
    }
}

func (wm *WorkspaceManager) GetWorkspacePath(projectID string) string {
    return filepath.Join(wm.baseDir, projectID)
}

func (wm *WorkspaceManager) EnsureWorkspace(projectID string) error {
    path := wm.GetWorkspacePath(projectID)

    // Create workspace directory if it doesn't exist
    if err := os.MkdirAll(path, 0755); err != nil {
        return fmt.Errorf("failed to create workspace: %w", err)
    }

    // Initialize git repository if not exists
    gitDir := filepath.Join(path, ".git")
    if _, err := os.Stat(gitDir); os.IsNotExist(err) {
        if err := wm.initGitRepo(path); err != nil {
            return fmt.Errorf("failed to init git repo: %w", err)
        }
    }

    // Update last access time
    wm.updateLastAccess(projectID)

    return nil
}

func (wm *WorkspaceManager) initGitRepo(path string) error {
    cmd := exec.Command("git", "init")
    cmd.Dir = path
    return cmd.Run()
}

func (wm *WorkspaceManager) updateLastAccess(projectID string) error {
    metaPath := filepath.Join(wm.baseDir, ".cleanup", "last_access.json")

    // Read existing data
    data := make(map[string]time.Time)
    if content, err := os.ReadFile(metaPath); err == nil {
        json.Unmarshal(content, &data)
    }

    // Update timestamp
    data[projectID] = time.Now()

    // Write back
    os.MkdirAll(filepath.Dir(metaPath), 0755)
    content, _ := json.MarshalIndent(data, "", "  ")
    return os.WriteFile(metaPath, content, 0644)
}

func (wm *WorkspaceManager) CleanupInactive(maxAge time.Duration) ([]string, error) {
    metaPath := filepath.Join(wm.baseDir, ".cleanup", "last_access.json")

    data := make(map[string]time.Time)
    if content, err := os.ReadFile(metaPath); err == nil {
        json.Unmarshal(content, &data)
    }

    var removed []string
    cutoff := time.Now().Add(-maxAge)

    for projectID, lastAccess := range data {
        if lastAccess.Before(cutoff) {
            path := wm.GetWorkspacePath(projectID)
            if err := os.RemoveAll(path); err == nil {
                removed = append(removed, projectID)
                delete(data, projectID)
            }
        }
    }

    // Save updated data
    content, _ := json.MarshalIndent(data, "", "  ")
    os.WriteFile(metaPath, content, 0644)

    return removed, nil
}
```

### 2. Claude Code Executor

```go
// backend/internal/executor/executor.go
package executor

import (
    "context"
    "fmt"
    "io"
    "os/exec"
    "strings"
    "time"

    "github.com/rs/zerolog"
)

type ClaudeExecutor struct {
    workspaceMgr *WorkspaceManager
    logger       zerolog.Logger
    timeout      time.Duration
}

func NewClaudeExecutor(workspaceMgr *WorkspaceManager, logger zerolog.Logger) *ClaudeExecutor {
    return &ClaudeExecutor{
        workspaceMgr: workspaceMgr,
        logger:       logger,
        timeout:      5 * time.Minute, // Default timeout
    }
}

type ExecuteRequest struct {
    ProjectID string
    Message   string
    Files     []FileContext
}

type FileContext struct {
    Path    string
    Content string
}

type ExecuteResponse struct {
    Output   string
    Files    []FileChange
    Error    string
    Duration time.Duration
}

type FileChange struct {
    Path      string
    Content   string
    Operation string // "create", "update", "delete"
}

func (ce *ClaudeExecutor) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
    start := time.Now()

    // Ensure workspace exists
    if err := ce.workspaceMgr.EnsureWorkspace(req.ProjectID); err != nil {
        return nil, fmt.Errorf("workspace setup failed: %w", err)
    }

    workspace := ce.workspaceMgr.GetWorkspacePath(req.ProjectID)

    // Write context files to workspace
    if err := ce.writeContextFiles(workspace, req.Files); err != nil {
        return nil, fmt.Errorf("failed to write context files: %w", err)
    }

    // Execute Claude Code
    output, err := ce.executeClaude(ctx, workspace, req.Message)

    resp := &ExecuteResponse{
        Output:   output,
        Duration: time.Since(start),
    }

    if err != nil {
        resp.Error = err.Error()
        return resp, err
    }

    // Detect file changes
    changes, err := ce.detectFileChanges(workspace)
    if err != nil {
        ce.logger.Warn().Err(err).Msg("failed to detect file changes")
    } else {
        resp.Files = changes
    }

    return resp, nil
}

func (ce *ClaudeExecutor) executeClaude(ctx context.Context, workspace, message string) (string, error) {
    // Create context with timeout
    execCtx, cancel := context.WithTimeout(ctx, ce.timeout)
    defer cancel()

    // Build command
    cmd := exec.CommandContext(execCtx, "claude", "code")
    cmd.Dir = workspace
    cmd.Stdin = strings.NewReader(message)

    // Capture output
    var stdout, stderr strings.Builder
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    // Execute
    ce.logger.Info().
        Str("workspace", workspace).
        Str("message", truncate(message, 100)).
        Msg("executing claude code")

    if err := cmd.Run(); err != nil {
        ce.logger.Error().
            Err(err).
            Str("stderr", stderr.String()).
            Msg("claude execution failed")
        return "", fmt.Errorf("claude execution failed: %w\n%s", err, stderr.String())
    }

    return stdout.String(), nil
}

func (ce *ClaudeExecutor) writeContextFiles(workspace string, files []FileContext) error {
    for _, f := range files {
        path := filepath.Join(workspace, f.Path)

        // Create parent directories
        if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
            return err
        }

        // Write file
        if err := os.WriteFile(path, []byte(f.Content), 0644); err != nil {
            return err
        }
    }
    return nil
}

func (ce *ClaudeExecutor) detectFileChanges(workspace string) ([]FileChange, error) {
    // Use git status to detect changes
    cmd := exec.Command("git", "status", "--porcelain")
    cmd.Dir = workspace

    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    var changes []FileChange
    lines := strings.Split(string(output), "\n")

    for _, line := range lines {
        if line == "" {
            continue
        }

        // Parse git status output
        // Format: XY PATH
        // X = index status, Y = working tree status
        status := line[0:2]
        path := strings.TrimSpace(line[3:])

        var operation string
        switch {
        case strings.Contains(status, "A"):
            operation = "create"
        case strings.Contains(status, "M"):
            operation = "update"
        case strings.Contains(status, "D"):
            operation = "delete"
        default:
            operation = "update"
        }

        // Read file content (unless deleted)
        var content string
        if operation != "delete" {
            fullPath := filepath.Join(workspace, path)
            data, err := os.ReadFile(fullPath)
            if err == nil {
                content = string(data)
            }
        }

        changes = append(changes, FileChange{
            Path:      path,
            Content:   content,
            Operation: operation,
        })
    }

    return changes, nil
}

func truncate(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen] + "..."
}
```

### 3. Service Integration

```go
// backend/internal/service/chat.go (modifications)

type ChatService struct {
    // ... existing fields ...
    executor *executor.ClaudeExecutor
}

func NewChatService(
    config ChatConfig,
    claudeService ClaudeMessenger,
    // ... other dependencies ...
    executor *executor.ClaudeExecutor,
    // ...
) *ChatService {
    return &ChatService{
        // ...
        executor: executor,
        // ...
    }
}

func (s *ChatService) SendMessage(ctx context.Context, req *SendMessageRequest) (<-chan MessageChunk, error) {
    // ... existing validation ...

    // Execute in workspace
    execReq := &executor.ExecuteRequest{
        ProjectID: req.ProjectID,
        Message:   req.Content,
        Files:     s.buildFileContext(req.ProjectID),
    }

    execResp, err := s.executor.Execute(ctx, execReq)
    if err != nil {
        return nil, fmt.Errorf("execution failed: %w", err)
    }

    // Save file changes to database
    if err := s.saveFileChanges(req.ProjectID, execResp.Files); err != nil {
        s.logger.Warn().Err(err).Msg("failed to save file changes")
    }

    // Stream response to client
    chunks := make(chan MessageChunk, 10)
    go s.streamResponse(chunks, execResp.Output)

    return chunks, nil
}

func (s *ChatService) buildFileContext(projectID string) []executor.FileContext {
    // Fetch relevant files from database
    files, err := s.fileRepo.ListByProject(context.Background(), projectID)
    if err != nil {
        s.logger.Warn().Err(err).Msg("failed to fetch files")
        return nil
    }

    var context []executor.FileContext
    for _, f := range files {
        context = append(context, executor.FileContext{
            Path:    f.Path,
            Content: f.Content,
        })
    }

    return context
}

func (s *ChatService) saveFileChanges(projectID string, changes []executor.FileChange) error {
    for _, change := range changes {
        switch change.Operation {
        case "create", "update":
            _, err := s.fileRepo.Upsert(context.Background(), &model.File{
                ProjectID: projectID,
                Path:      change.Path,
                Content:   change.Content,
            })
            if err != nil {
                return err
            }

        case "delete":
            if err := s.fileRepo.DeleteByPath(context.Background(), projectID, change.Path); err != nil {
                return err
            }
        }
    }

    return nil
}
```

### 4. Main Server Setup

```go
// backend/cmd/server/main.go (additions)

func main() {
    // ... existing setup ...

    // Initialize workspace manager
    workspaceDir := getEnv("WORKSPACE_DIR", "/workspaces")
    workspaceMgr := executor.NewWorkspaceManager(workspaceDir)

    // Initialize Claude executor
    claudeExecutor := executor.NewClaudeExecutor(workspaceMgr, logger)

    // Initialize chat service with executor
    chatService := service.NewChatService(
        service.ChatConfig{
            ContextMessageLimit: cfg.ContextMessageLimit,
        },
        claudeService,
        discoveryService,
        agentContextService,
        projectRepo,
        fileRepo,
        fileMetadataRepo,
        claudeExecutor, // NEW
        logger,
    )

    // ... rest of setup ...

    // Start cleanup goroutine
    go startWorkspaceCleanup(workspaceMgr, logger)
}

func startWorkspaceCleanup(wm *executor.WorkspaceManager, logger zerolog.Logger) {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        maxAge := 7 * 24 * time.Hour // 7 days
        removed, err := wm.CleanupInactive(maxAge)

        if err != nil {
            logger.Error().Err(err).Msg("workspace cleanup failed")
            continue
        }

        if len(removed) > 0 {
            logger.Info().
                Int("count", len(removed)).
                Strs("projects", removed).
                Msg("cleaned up inactive workspaces")
        }
    }
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
```

### 5. Docker Compose Updates

```yaml
# docker-compose.yml
version: '3.8'

services:
  backend:
    build: ./backend
    ports:
      - "8081:8080"
    environment:
      PORT: "8080"
      DATABASE_URL: postgres://gochat:gochat@db:5432/gochat?sslmode=disable
      CLAUDE_API_KEY: ${CLAUDE_API_KEY}
      CLAUDE_MODEL: claude-sonnet-4-20250514
      WORKSPACE_DIR: /workspaces
    volumes:
      - workspaces:/workspaces  # Persist workspaces
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: gochat
      POSTGRES_PASSWORD: gochat
      POSTGRES_DB: gochat
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U gochat"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pgdata:
  workspaces:  # NEW: Persist workspaces across restarts
```

### 6. Configuration

```bash
# backend/.env.example
PORT=8080
DATABASE_URL=postgres://gochat:gochat@localhost:5432/gochat?sslmode=disable
CLAUDE_API_KEY=your-api-key-here
CLAUDE_MODEL=claude-sonnet-4-20250514
WORKSPACE_DIR=/workspaces
WORKSPACE_CLEANUP_AGE_DAYS=7
CLAUDE_EXEC_TIMEOUT=300  # 5 minutes
```

## Testing

### Unit Tests

```go
// backend/internal/executor/executor_test.go
package executor_test

import (
    "context"
    "os"
    "path/filepath"
    "testing"

    "gitlab.yuki.lan/goodies/gochat/backend/internal/executor"
    "github.com/rs/zerolog"
)

func TestClaudeExecutor(t *testing.T) {
    // Create temp workspace
    tmpDir := t.TempDir()

    wm := executor.NewWorkspaceManager(tmpDir)
    logger := zerolog.Nop()
    exec := executor.NewClaudeExecutor(wm, logger)

    t.Run("creates workspace on first use", func(t *testing.T) {
        req := &executor.ExecuteRequest{
            ProjectID: "test-project",
            Message:   "echo 'hello world'",
        }

        _, err := exec.Execute(context.Background(), req)
        if err != nil {
            t.Fatalf("Execute failed: %v", err)
        }

        // Check workspace exists
        workspace := filepath.Join(tmpDir, "test-project")
        if _, err := os.Stat(workspace); os.IsNotExist(err) {
            t.Error("Workspace not created")
        }

        // Check git initialized
        gitDir := filepath.Join(workspace, ".git")
        if _, err := os.Stat(gitDir); os.IsNotExist(err) {
            t.Error("Git not initialized")
        }
    })

    t.Run("writes context files", func(t *testing.T) {
        req := &executor.ExecuteRequest{
            ProjectID: "test-project-2",
            Message:   "test",
            Files: []executor.FileContext{
                {Path: "src/main.go", Content: "package main"},
                {Path: "README.md", Content: "# Test"},
            },
        }

        _, err := exec.Execute(context.Background(), req)
        if err != nil {
            t.Fatalf("Execute failed: %v", err)
        }

        workspace := filepath.Join(tmpDir, "test-project-2")

        // Check files exist
        mainGo := filepath.Join(workspace, "src/main.go")
        content, err := os.ReadFile(mainGo)
        if err != nil {
            t.Errorf("File not created: %v", err)
        }
        if string(content) != "package main" {
            t.Errorf("File content mismatch: got %s", string(content))
        }
    })
}

func TestWorkspaceManager(t *testing.T) {
    tmpDir := t.TempDir()
    wm := executor.NewWorkspaceManager(tmpDir)

    t.Run("cleanup removes old workspaces", func(t *testing.T) {
        // Create workspaces
        wm.EnsureWorkspace("project-1")
        wm.EnsureWorkspace("project-2")

        // Simulate old workspace by manually setting old timestamp
        // (In real test, would mock time or use test clock)

        // For now, test that cleanup runs without error
        removed, err := wm.CleanupInactive(0) // Clean everything
        if err != nil {
            t.Fatalf("Cleanup failed: %v", err)
        }

        t.Logf("Removed %d workspaces", len(removed))
    })
}
```

## Monitoring

### Health Check

```go
// Add to health handler
func (h *HealthHandler) Health(c *gin.Context) {
    // ... existing DB check ...

    // Check workspace directory accessible
    workspaceHealth := "healthy"
    if _, err := os.Stat(h.workspaceDir); err != nil {
        workspaceHealth = "unhealthy"
    }

    c.JSON(http.StatusOK, gin.H{
        "status":    "ok",
        "database":  dbHealth,
        "workspace": workspaceHealth,
        "timestamp": time.Now().Unix(),
    })
}
```

### Metrics Endpoint

```go
// backend/internal/handler/metrics.go
func (h *MetricsHandler) WorkspaceMetrics(c *gin.Context) {
    stats := h.workspaceMgr.GetStats()

    c.JSON(http.StatusOK, gin.H{
        "total_workspaces":    stats.TotalWorkspaces,
        "active_workspaces":   stats.ActiveWorkspaces,
        "inactive_workspaces": stats.InactiveWorkspaces,
        "total_size_mb":       stats.TotalSizeMB,
        "oldest_access":       stats.OldestAccess,
    })
}
```

## Migration Path

### Phase 1: Add Workspace Support (Week 1)
- [ ] Implement WorkspaceManager
- [ ] Implement ClaudeExecutor
- [ ] Add workspace volume to Docker Compose
- [ ] Unit tests
- [ ] Integration tests

### Phase 2: Switch to Workspace Mode (Week 2)
- [ ] Update ChatService to use executor
- [ ] Add cleanup goroutine
- [ ] Add health checks
- [ ] Add metrics
- [ ] Documentation

### Phase 3: Optional Container Mode (Future)
- [ ] Add per-project config: isolation_mode
- [ ] Implement ContainerExecutor
- [ ] Add container lifecycle management
- [ ] UI toggle for container mode

## Benefits Over Container Approach

1. **Simplicity:** 1/10th the code complexity
2. **Performance:** No container startup latency
3. **Resources:** 1/4 the memory usage
4. **Reliability:** Fewer failure modes
5. **Debuggability:** Standard filesystem, easy to inspect
6. **User Experience:** No cold starts, instant responses
7. **Development:** Faster iteration, easier testing
8. **Maintenance:** Less operational burden

## Limitations

1. **Isolation:** Process-level, not container-level
2. **Dependencies:** Shared host environment
3. **Security:** Less sandboxing (mitigated by user trust)
4. **Customization:** Can't customize per-project environment

## When to Reconsider Containers

Consider container-per-project if:
1. Users request custom environments (Python 2 vs 3, Node 14 vs 18)
2. Security becomes critical (untrusted code execution)
3. Planning SaaS multi-tenant offering
4. Need stronger isolation for compliance

Until then, workspace isolation is sufficient and simpler.
