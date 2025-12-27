# Session 011 - Hybrid Summary Modal UX Fixes

**Date**: 2025-12-27
**Branch**: `feature/discovery-app-map-seeding`
**Context**: 90% (handoff triggered)

## Summary

Fixed multiple bugs in the hybrid summary modal UX and mock discovery service flow.

## Issues Fixed

### 1. Duplicate Welcome Messages
- Mock service was returning `welcome_response` multiple times
- Added `hasAssistantMessages()` check to skip welcome if already sent

### 2. Duplicate Stage Messages (2nd/3rd responses)
- Changed fixture selection from global `messageCount` to `countAssistantMessagesInStage()`
- Now counts responses per stage by matching content patterns

### 3. Markdown Formatting Lost (Bullets not rendering)
- `createMockStream` was using `strings.Fields()` which splits on whitespace including newlines
- Fixed to stream by character chunks to preserve formatting

### 4. Premature Discovery Complete Notification
- `summary_response.json` had `stage_complete: true` - changed to `false`
- Added `hasMeaningfulSummary` check (requires MVP features) before showing notification bar

### 5. Summary Modal Not Populated
- Fixed all fixture files: `users` was an object, changed to array format
- Added `ClearUsers()` method to prevent duplicate users from accumulating

### 6. Progress Circles Reset on "Edit in Chat"
- `handleEditDiscovery` was calling `resetDiscovery()` - removed that call
- Now just closes modal and lets user continue conversation

### 7. Project Title Not Updating Dynamically
- Added `onDiscoveryConfirmed` callback to `ChatContainer`
- Calls `fetchProjects()` after confirmation to refresh title
- Added `renameProjectFromDiscovery()` call in `ConfirmDiscovery` function

## Files Modified

### Backend
- `backend/internal/service/claude_mock.go` - Streaming fix, fixture selection logic
- `backend/internal/service/discovery.go` - ClearUsers call, rename on confirm
- `backend/internal/repository/discovery.go` - Added ClearUsers method
- `backend/testdata/discovery/*.json` - All fixtures: users format, stage_complete values

### Frontend
- `frontend/src/components/chat/ChatContainer.tsx` - Edit handler, summary visibility, confirm callback
- `frontend/src/components/ProjectPageClient.tsx` - Pass fetchProjects callback
- `frontend/src/components/HomeClient.tsx` - Pass fetchProjects callback
- `frontend/src/app/demo/discovery/page.tsx` - Extracted client component for build

## Pending

- User testing to verify all fixes work together
- Commit the changes (large changeset)
- Test Phase 3 Learning Journey system

## Key Decisions

- Used content pattern matching for stage detection in mock (more reliable than counters)
- Clear users before adding to prevent duplicates (simpler than upsert)
- Require MVP features for summary notification (prevents showing empty modal)
