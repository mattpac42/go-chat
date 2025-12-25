# PRD-002: Backend API + AI Integration

**Version**: 1.0
**Created**: 2025-12-24
**Author**: Product Manager (Tactical)
**Status**: Draft
**Phase**: 1 - Foundation (Weeks 1-2)

---

## Overview

### Problem Statement

The Go Chat frontend needs a robust backend service that can manage project/conversation data, handle real-time communication, and integrate with Claude AI for code generation. The backend must handle streaming AI responses to provide immediate feedback, maintain conversation context for coherent multi-turn interactions, and extract code blocks from AI responses for proper rendering.

### Proposed Solution

A Go-based backend service using the Gin framework for REST APIs and Gorilla WebSocket for real-time communication. The service integrates with Claude API for AI code generation, stores conversation data in PostgreSQL, and streams AI responses to the frontend in real-time.

### Goals

1. Provide reliable REST APIs for project and conversation management
2. Enable real-time, streaming AI responses via WebSocket
3. Maintain conversation context for multi-turn interactions
4. Extract and identify code blocks from AI responses
5. Support production deployment via GitLab CI/CD

### Non-Goals

- User authentication and authorization
- Rate limiting (single user MVP)
- Multi-tenancy and user isolation
- Database migrations tooling (manual for MVP)
- Background job processing
- Caching layer (Redis)
- Horizontal scaling considerations

---

## User Stories

### Primary User

**Go Chat Frontend Application**
- Requires reliable API endpoints for data operations
- Needs real-time streaming for AI responses
- Expects structured responses with code block metadata

### Derived User Stories

1. **US-001**: As the frontend, I want to create, list, and delete projects via REST API, so that users can manage multiple conversations.

2. **US-002**: As the frontend, I want to retrieve a project with its full message history, so that users can continue previous conversations.

3. **US-003**: As the frontend, I want to send a chat message via WebSocket and receive the AI response as a stream, so that users see responses in real-time.

4. **US-004**: As the frontend, I want the AI to remember previous messages in the conversation, so that users can have coherent multi-turn discussions.

5. **US-005**: As the frontend, I want code blocks in AI responses to be identified with language metadata, so that I can render them with proper syntax highlighting.

6. **US-006**: As the operator, I want structured logging and health endpoints, so that I can monitor the service in production.

---

## Requirements

### Functional Requirements

| ID | Requirement | Priority | Acceptance Criteria |
|----|-------------|----------|---------------------|
| FR-001 | Create project | High | POST /api/projects creates new project, returns project ID |
| FR-002 | List projects | High | GET /api/projects returns all projects with titles and timestamps |
| FR-003 | Get project with messages | High | GET /api/projects/:id returns project metadata and full message history |
| FR-004 | Delete project | Medium | DELETE /api/projects/:id removes project and all messages |
| FR-005 | WebSocket connection | High | /ws/chat accepts WebSocket connections and maintains state |
| FR-006 | Send message via WebSocket | High | Client can send chat_message and receive acknowledgment |
| FR-007 | AI response streaming | High | AI responses stream to client as message_chunk events |
| FR-008 | Conversation context | High | AI receives last N messages as context (configurable, default 20) |
| FR-009 | Code block extraction | High | AI responses parsed to identify code blocks with language |
| FR-010 | Message persistence | High | All messages (user and AI) stored in database |
| FR-011 | Health endpoint | Medium | GET /health returns service status and database connectivity |
| FR-012 | Graceful shutdown | Medium | Service handles SIGTERM, completes in-flight requests |
| FR-013 | Error handling | High | All endpoints return structured error responses |
| FR-014 | Connection heartbeat | Medium | WebSocket ping/pong keeps connections alive |
| FR-015 | Auto-generate project title | Low | First message used to generate project title via AI |

### Non-Functional Requirements

| ID | Requirement | Criteria |
|----|-------------|----------|
| NFR-001 | Response Time | REST API responses < 100ms (excluding AI calls) |
| NFR-002 | WebSocket Latency | Message delivery < 50ms once received from AI |
| NFR-003 | Concurrent Connections | Support 10 concurrent WebSocket connections (single user + development) |
| NFR-004 | Availability | 99% uptime during business hours |
| NFR-005 | Logging | Structured JSON logs with request ID, timestamps, levels |
| NFR-006 | Memory | < 256MB memory usage under normal load |
| NFR-007 | Startup Time | Service ready to accept requests < 5 seconds |
| NFR-008 | Security | No credentials in logs; HTTPS/WSS in production |

---

## Technical Requirements

### Technology Stack

| Component | Technology | Version | Rationale |
|-----------|------------|---------|-----------|
| Language | Go | 1.22+ | Performance, simplicity, project namesake |
| HTTP Framework | Gin | 1.9+ | Fast, mature, good middleware ecosystem |
| WebSocket | Gorilla WebSocket | 1.5+ | Industry standard, well-maintained |
| Database | PostgreSQL | 16 | Reliability, JSON support, production-ready |
| DB Driver | sqlx | 1.3+ | Enhanced database/sql with struct scanning |
| Logging | zerolog | 1.31+ | Structured JSON logging, zero-allocation |
| Configuration | envconfig | 1.2+ | Environment-based configuration |
| AI Client | HTTP (Claude API) | v1 | Direct HTTP to Claude messages API |

### Architecture

```
cmd/
├── server/
│   └── main.go              # Application entry point
internal/
├── config/
│   └── config.go            # Configuration loading
├── handler/
│   ├── health.go            # Health check handlers
│   ├── project.go           # Project CRUD handlers
│   └── websocket.go         # WebSocket handler
├── middleware/
│   ├── logging.go           # Request logging
│   ├── cors.go              # CORS configuration
│   └── recovery.go          # Panic recovery
├── model/
│   ├── project.go           # Project entity
│   └── message.go           # Message entity
├── repository/
│   ├── project.go           # Project database operations
│   └── message.go           # Message database operations
├── service/
│   ├── chat.go              # Chat orchestration service
│   └── claude.go            # Claude API client
├── websocket/
│   ├── hub.go               # WebSocket connection hub
│   ├── client.go            # Client connection handler
│   └── message.go           # WebSocket message types
└── pkg/
    └── markdown/
        └── codeblock.go     # Code block extraction
```

### Database Schema

```sql
-- projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL DEFAULT 'New Project',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- messages table
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('user', 'assistant')),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_messages_project_id ON messages(project_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
```

---

## API Contracts

### REST API Endpoints

#### Health Check

```
GET /health

Response 200:
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-12-24T10:30:00Z"
}
```

#### List Projects

```
GET /api/projects

Response 200:
{
  "projects": [
    {
      "id": "uuid",
      "title": "Inventory Tracker",
      "messageCount": 5,
      "createdAt": "2025-12-24T10:00:00Z",
      "updatedAt": "2025-12-24T10:30:00Z"
    }
  ]
}
```

#### Create Project

```
POST /api/projects
Content-Type: application/json

Request:
{
  "title": "My New Project"  // optional, defaults to "New Project"
}

Response 201:
{
  "id": "uuid",
  "title": "My New Project",
  "createdAt": "2025-12-24T10:00:00Z",
  "updatedAt": "2025-12-24T10:00:00Z"
}
```

#### Get Project with Messages

```
GET /api/projects/:id

Response 200:
{
  "id": "uuid",
  "title": "Inventory Tracker",
  "createdAt": "2025-12-24T10:00:00Z",
  "updatedAt": "2025-12-24T10:30:00Z",
  "messages": [
    {
      "id": "uuid",
      "role": "user",
      "content": "Build me an inventory tracker",
      "createdAt": "2025-12-24T10:00:00Z"
    },
    {
      "id": "uuid",
      "role": "assistant",
      "content": "I'll create an inventory tracker for you...",
      "codeBlocks": [
        {
          "language": "python",
          "code": "...",
          "startLine": 5,
          "endLine": 25
        }
      ],
      "createdAt": "2025-12-24T10:00:05Z"
    }
  ]
}
```

#### Delete Project

```
DELETE /api/projects/:id

Response 204: (no content)

Response 404:
{
  "error": "project not found"
}
```

### WebSocket Protocol

#### Connection

```
WS /ws/chat?projectId={uuid}
```

#### Client to Server Messages

```typescript
// Send chat message
{
  "type": "chat_message",
  "content": "Build me an inventory tracker",
  "timestamp": "2025-12-24T10:00:00Z"
}

// Ping (keep-alive)
{
  "type": "ping"
}
```

#### Server to Client Messages

```typescript
// Message acknowledged, streaming starts
{
  "type": "message_start",
  "messageId": "uuid",
  "timestamp": "2025-12-24T10:00:00Z"
}

// Streaming chunk
{
  "type": "message_chunk",
  "messageId": "uuid",
  "content": "I'll create",  // incremental content
  "timestamp": "2025-12-24T10:00:01Z"
}

// Message complete
{
  "type": "message_complete",
  "messageId": "uuid",
  "fullContent": "I'll create an inventory tracker...",
  "codeBlocks": [
    {
      "language": "python",
      "code": "class InventoryItem:\n    ...",
      "startIndex": 150,
      "endIndex": 450
    }
  ],
  "timestamp": "2025-12-24T10:00:30Z"
}

// Error
{
  "type": "error",
  "error": "Failed to generate response",
  "code": "AI_ERROR",
  "timestamp": "2025-12-24T10:00:05Z"
}

// Pong response
{
  "type": "pong"
}
```

---

## Claude API Integration

### System Prompt

```
You are Go Chat, an AI assistant that helps users build applications by writing code.

When the user describes what they want to build:
1. Understand their requirements
2. Ask clarifying questions if needed (but try to infer reasonable defaults)
3. Generate complete, working code
4. Explain what the code does in simple terms

Always format code in markdown code blocks with the language specified:
```python
# your code here
```

Keep explanations brief and non-technical. The user may not be a developer.
Focus on creating working, practical solutions.
```

### Request Format

```json
{
  "model": "claude-sonnet-4-20250514",
  "max_tokens": 4096,
  "system": "[system prompt above]",
  "messages": [
    {"role": "user", "content": "Build me an inventory tracker"},
    {"role": "assistant", "content": "I'll create..."},
    {"role": "user", "content": "Add a search feature"}
  ],
  "stream": true
}
```

### Context Window Management

- Include last 20 messages (configurable via `CONTEXT_MESSAGE_LIMIT`)
- If total tokens exceed 80% of model limit, truncate oldest messages
- Always include system prompt
- Estimate tokens: ~4 chars per token for English text

---

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | 8080 | HTTP server port |
| `DATABASE_URL` | Yes | - | PostgreSQL connection string |
| `CLAUDE_API_KEY` | Yes | - | Anthropic API key |
| `CLAUDE_MODEL` | No | claude-sonnet-4-20250514 | Claude model to use |
| `CLAUDE_MAX_TOKENS` | No | 4096 | Max tokens in response |
| `CONTEXT_MESSAGE_LIMIT` | No | 20 | Messages to include as context |
| `LOG_LEVEL` | No | info | Log level (debug, info, warn, error) |
| `CORS_ORIGINS` | No | * | Allowed CORS origins |
| `WS_PING_INTERVAL` | No | 30 | WebSocket ping interval (seconds) |
| `WS_PONG_TIMEOUT` | No | 60 | WebSocket pong timeout (seconds) |

---

## Acceptance Criteria

### Phase 1 Exit Criteria (from Roadmap)

- [ ] AI responds to messages with context
- [ ] Describing "inventory tracker" produces runnable code
- [ ] Multi-message conversations work

### Detailed Acceptance Tests

1. **AC-001**: POST /api/projects creates a new project; GET /api/projects lists it
2. **AC-002**: WebSocket connection to /ws/chat?projectId=X succeeds and stays open
3. **AC-003**: Sending "Build me an inventory tracker" via WebSocket returns streaming response with Python code
4. **AC-004**: Response includes code block metadata with language="python"
5. **AC-005**: Second message "Add a search feature" references the previous inventory tracker context
6. **AC-006**: GET /api/projects/:id returns all messages from the conversation
7. **AC-007**: DELETE /api/projects/:id removes project and all messages
8. **AC-008**: GET /health returns healthy status when database is connected
9. **AC-009**: WebSocket reconnection after disconnect resumes with full message history
10. **AC-010**: Service starts and accepts requests within 5 seconds

---

## Dependencies

### Internal Dependencies

- PostgreSQL database provisioned and accessible
- Network access to Anthropic API (api.anthropic.com)

### External Dependencies

| Dependency | Type | Fallback |
|------------|------|----------|
| Claude API | Critical | Service degraded; returns error |
| PostgreSQL | Critical | Service unavailable |

---

## Out of Scope

The following are explicitly NOT included in this PRD:

1. User authentication and authorization
2. Rate limiting
3. Message editing or deletion
4. File/image attachments
5. Multiple AI model support
6. Conversation branching
7. Response regeneration
8. Token usage tracking/billing
9. Caching layer
10. API versioning (v1 prefix not required for MVP)
11. OpenAPI/Swagger documentation
12. Database migrations tooling

---

## Deployment

### Local Development

```yaml
# docker-compose.yml (simplified)
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/gochat
      - CLAUDE_API_KEY=${CLAUDE_API_KEY}
    depends_on:
      - db

  db:
    image: postgres:16
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=gochat
```

### Production Deployment

Production uses GitLab CI/CD with Gold Bricks pipeline:

1. GitLab CI builds Docker image on push to main
2. Image pushed to registry at `gitlab.yuki.lan:5050`
3. SSH deploy to yuki-platform server (192.168.0.200)
4. Service runs via docker-compose on production host

See `/workspace/.claude/templates/secure-pipeline/` for pipeline templates.

---

## Risks and Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Claude API latency/timeouts | High | Medium | Implement 60-second timeout; stream partial responses; clear error messages |
| WebSocket connection drops | Medium | Medium | Client-side reconnection; heartbeat ping/pong |
| Database connection pool exhaustion | High | Low | Configure appropriate pool size; connection timeouts |
| Large AI responses exceed memory | Medium | Low | Stream responses; don't buffer full response |
| Code block parsing fails | Low | Medium | Graceful fallback; treat as plain text |

---

## Success Metrics

| Metric | Target | Measurement Method |
|--------|--------|-------------------|
| API response time (non-AI) | < 100ms p95 | Structured logging with timing |
| WebSocket message delivery | < 50ms | Timestamp comparison |
| AI response start time | < 2 seconds | Time from request to first chunk |
| Error rate | < 1% | Error logs / total requests |
| Successful code generation | > 80% | Responses containing code blocks |

---

## Timeline

| Milestone | Description | Target |
|-----------|-------------|--------|
| Database schema ready | Tables created, local DB running | Day 1 |
| REST API complete | Project CRUD working | Day 4 |
| WebSocket functional | Connection and basic messaging | Day 6 |
| Claude integration | AI responses streaming | Day 9 |
| Code block extraction | Metadata in responses | Day 11 |
| Integration tested | End-to-end with frontend | Day 13 |
| Phase 1 complete | All acceptance criteria met | Day 14 |

---

## Approval

| Role | Name | Date | Status |
|------|------|------|--------|
| Product | | | [ ] |
| Engineering | | | [ ] |
| Platform | | | [ ] |

---

**Document History**

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2025-12-24 | 1.0 | Initial PRD | Product Manager |
