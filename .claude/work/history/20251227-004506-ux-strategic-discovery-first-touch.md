# UX Strategic Session: Discovery First-Touch Experience

**Date**: 2025-12-27
**Agent**: ux-strategic
**Task**: Design recommendations for discovery flow first-touch experience

## Work Completed

Analyzed the first-touch user experience problem in the guided discovery flow. The core issue is a mismatch between user-initiated empty state patterns and the bot-initiated nature of the discovery conversation.

### Problem Identified

1. Current empty state shows "Start a conversation" and "Describe what you want to build"
2. User types comprehensive description expecting to kick off the project
3. Bot's welcome message arrives asking basic questions ("Tell me about yourself")
4. User feels confused - their detailed input gets processed as a simple stage response

### Root Cause

The discovery flow is architecturally bot-initiated but the UI presents a user-initiated pattern. The welcome message from the Product Guide only appears AFTER the user sends their first message (reactive), not BEFORE (proactive).

## Recommendations Delivered

### 1. Bot-First Discovery Pattern (Priority: P0)
- Auto-send welcome message when discovery is created
- User sees bot message immediately upon entering new project
- Eliminates the misleading empty state entirely

### 2. Discovery-Aware Empty State (Priority: P0)
- If any loading delay, show "Your Product Guide is joining..."
- Replace "Describe what you want to build" with contextual loading state
- Sets correct expectation that bot will speak first

### 3. Input State Management (Priority: P1)
- Disable input until welcome message appears
- Placeholder: "Waiting for Guide..."
- Prevents race condition where user types before bot speaks

### 4. Discovery Mode Introduction Card (Priority: P2)
- Brief orientation before first bot message
- Sets time expectation ("about 5 minutes")
- Explains value of discovery process

### 5. Returning User Fast-Path (Priority: P3)
- Detect returning users and offer skip option
- "Start Fresh Discovery" vs "Quick Start"
- Design already exists in guided-discovery-ux.md

## Files Reviewed

- `/workspace/frontend/src/components/chat/MessageList.tsx` - Empty state component
- `/workspace/frontend/src/components/chat/ChatContainer.tsx` - Input placeholders and discovery integration
- `/workspace/frontend/src/hooks/useDiscovery.ts` - Discovery state management
- `/workspace/backend/internal/service/chat.go` - Message processing flow
- `/workspace/.claude/work/design/guided-discovery-ux.md` - Existing UX specification

## Decisions Made

- **Bot-first pattern over user-first**: Discovery is a guided experience, not a blank canvas. The bot should always speak first.
- **Backend initialization over frontend synthetic messages**: Auto-creating the welcome message on the backend ensures persistence and accurate conversation history.
- **Disable input during initialization**: Prevents user error and wasted effort from typing before context is set.

## Recommendations

For tactical implementation:
1. Backend: Modify discovery creation to auto-send welcome message as first assistant message
2. Frontend `MessageList.tsx`: Add discovery-aware loading state (lines 97-107)
3. Frontend `ChatContainer.tsx`: Disable input when discovery mode active and messages.length === 0
4. Update input placeholder to "Waiting for Guide..." during initialization

## Artifacts

Full design analysis included in the response with:
- Information architecture diagram
- Implementation priority matrix
- Design principles applied
- Summary of changes by component
