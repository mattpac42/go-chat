# PRD-001: Tool Use for File Operations

## Overview

Replace unreliable markdown code block parsing with Claude's native Tool Use API for file operations, ensuring 100% reliable file creation.

## Problem Statement

Currently, the go-chat backend parses markdown code blocks from Claude's responses to extract files. This requires Claude to follow a specific format (`\`\`\`language:filename`), but compliance is inconsistent (~70-80%), resulting in files not being created even when Claude claims to have created them.

## Solution

Implement Claude's Tool Use API with explicit `write_file`, `read_file`, and `list_files` tools. Tool calls are structured JSON, eliminating parsing ambiguity and providing a feedback loop where Claude receives confirmation of file operations.

## Phases

### Phase 1: Tool Use Implementation (This PRD)
- Add tool definitions to Claude API requests
- Handle `tool_use` response blocks
- Implement tool result feedback
- Keep markdown parsing as fallback

### Phase 2: Container Execution (Future PRD)
- Add `run_command` tool
- Container pool for code execution
- Build/test capabilities

### Phase 3: Full Devcontainer (Future PRD)
- Persistent project environments
- GitLab integration
- Full development workflow

## Requirements

### Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-1 | Define `write_file` tool with path and content parameters | Must |
| FR-2 | Define `read_file` tool with path parameter | Must |
| FR-3 | Define `list_files` tool with optional path parameter | Should |
| FR-4 | Parse `tool_use` content blocks from Claude responses | Must |
| FR-5 | Execute file operations and return `tool_result` | Must |
| FR-6 | Support multi-turn conversations with tool results | Must |
| FR-7 | Maintain backward compatibility with markdown parsing | Should |

### Non-Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| NFR-1 | No increase in response latency beyond tool execution time | Must |
| NFR-2 | File operations must be atomic (write succeeds or fails cleanly) | Must |
| NFR-3 | Tool errors must be reported back to Claude for recovery | Must |

## Technical Design

### Tool Definitions

```json
{
  "name": "write_file",
  "description": "Create or overwrite a file with the given content",
  "input_schema": {
    "type": "object",
    "properties": {
      "path": {
        "type": "string",
        "description": "File path relative to project root"
      },
      "content": {
        "type": "string",
        "description": "File content to write"
      }
    },
    "required": ["path", "content"]
  }
}
```

### Response Flow

```
1. User sends message
2. Backend sends to Claude with tools defined
3. Claude responds with tool_use block
4. Backend executes tool, saves file
5. Backend sends tool_result back to Claude
6. Claude continues or completes response
7. WebSocket streams final response to frontend
```

### Files to Modify

| File | Changes |
|------|---------|
| `backend/internal/service/claude.go` | Add tool definitions, parse tool_use |
| `backend/internal/service/chat.go` | Execute tools, send tool_result |
| `backend/internal/handler/websocket.go` | Handle tool execution events |

## Success Criteria

1. Files created via tool use appear in file explorer 100% of the time
2. Claude receives confirmation of file creation
3. Existing markdown-based projects continue to work
4. No regression in chat response quality

## Out of Scope

- Code execution (Phase 2)
- Container management (Phase 2/3)
- File editing tools (can be added later)
- Directory operations

## Timeline

- Implementation: 1-2 weeks
- Testing: 2-3 days
- Rollout: Immediate (with fallback)
