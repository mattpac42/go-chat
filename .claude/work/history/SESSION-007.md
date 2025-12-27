# Session 007 - Mock Claude Service & Discovery Flow Fixes

**Date**: 2025-12-26
**Branch**: `feature/guided-discovery`

## Summary

Implemented mock Claude service toggle for testing discovery flow without API costs, fixed multiple issues with discovery state progression, and added automatic project renaming after discovery completes.

## Completed Work

### 1. Mock Claude Service Toggle
- Added `ClaudeMessenger` interface in `backend/internal/service/claude.go`
- Updated `ChatService` to accept interface instead of concrete type
- Added `USE_MOCK_CLAUDE=true` environment variable support in `main.go`
- Mock service loads fixtures from `testdata/discovery/*.json`

### 2. Fixed Discovery Fixtures
- Updated all followup fixtures to have `stage_complete: true`:
  - `welcome_response.json`
  - `problem_followup_response.json`
  - `personas_followup_response.json`
  - `mvp_followup_response.json`
- Mock service now appends `<!--DISCOVERY_DATA:...-->` metadata to responses

### 3. Fixed Frontend Discovery Hook
- Added `API_BASE_URL` to `useDiscovery.ts` (was calling wrong URL)
- Added `refetchDiscovery()` call in `ChatContainer.tsx` after each message
- Progress indicator now updates as stages advance

### 4. UI Fix - Chat Input Border
- Changed from `ring-2` to `border-2` in `ChatInput.tsx`
- Focus border was getting clipped by viewport edge
- Removed bottom spacer div that was covering the input

### 5. Auto-Rename Project After Discovery
- Added `ProjectRepository` to `DiscoveryService`
- New `renameProjectFromDiscovery()` method called when stage advances to "complete"
- Updated summary prompt to request 1-3 word project names
- Extraction now checks both top-level and nested `summary` object for project_name

## Modified Files

### Backend
- `backend/cmd/server/main.go` - Mock service toggle, pass projectRepo to discovery
- `backend/internal/service/claude.go` - Added ClaudeMessenger interface
- `backend/internal/service/claude_mock.go` - Appends metadata to responses
- `backend/internal/service/chat.go` - Uses ClaudeMessenger interface
- `backend/internal/service/discovery.go` - Added projectRepo, auto-rename logic
- `backend/internal/service/prompts/discovery.go` - 1-3 word project name requirement
- `backend/testdata/discovery/*.json` - Fixed stage_complete flags

### Frontend
- `frontend/src/hooks/useDiscovery.ts` - Fixed API URLs
- `frontend/src/components/chat/ChatContainer.tsx` - Refetch discovery after messages
- `frontend/src/components/chat/ChatInput.tsx` - Fixed border clipping

## Key Decisions

1. **Interface over concrete type** - `ClaudeMessenger` interface allows swapping real/mock service
2. **Metadata in response** - Mock service appends `<!--DISCOVERY_DATA:...-->` to match real API behavior
3. **Auto-rename at completion** - Project renamed using discovered name when discovery completes
4. **Nested extraction** - Supports project_name at top-level or inside `summary` object

## Testing Notes

### Mock Mode
```bash
cd backend && USE_MOCK_CLAUDE=true go run ./cmd/server
```

### Real API Mode
```bash
cd backend && unset USE_MOCK_CLAUDE && source .env && go run ./cmd/server
```

### Database Migration Required
```bash
psql postgres://gochat:gochat@localhost:5432/gochat -f backend/migrations/004_discovery.sql
```

## Pending Work

1. **DiscoverySummaryCard** - Display summary inline in chat when discovery completes
2. **Edit Details flow** - Allow editing discovery data before confirming
3. **End-to-end testing** - Test full flow with real Claude API
4. **Commit changes** - All changes are uncommitted

## Git Status

Multiple files modified, nothing committed this session.
