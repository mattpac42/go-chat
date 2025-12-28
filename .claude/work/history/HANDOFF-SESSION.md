# Handoff - Session 014

## Immediate Context

Working on **Chat UX improvements & new features**. Major work on progressive disclosure, image upload, cost savings, and sidebar simplification.

## Branch

`main` (v1.0.0 released this session)

## 🚨 BUG TO FIX FIRST

**Clipboard paste shows preview but image not sent to Claude**

When user pastes a screenshot:
1. ✅ Preview appears at bottom of chat input ("Preview of pasted image")
2. ✅ X button to remove works
3. ❌ When message is sent, Root says "I don't see a screenshot attached"

**The image is not being uploaded or included in the message.**

Files to check:
- `frontend/src/components/chat/ChatInput.tsx` - handleSend function, upload logic
- `backend/internal/handler/upload.go` - upload endpoint
- Check if upload is actually being called
- Check if image data is being sent to Claude in the message

Screenshot of bug attached to this session.

## Last Tasks Completed

1. Progressive disclosure for long chat content
2. Image upload backend with Claude Vision
3. Clipboard paste for screenshots
4. Up arrow message history in chat
5. Cost savings card (integrated into summary modal)
6. Simplified streaming (removed raw text toggle)
7. Fixed auto-scroll during streaming
8. Removed card view from file explorer (kept Purpose + Tree only)

## Critical Files to Read

1. `frontend/src/components/chat/ChatInput.tsx` - Clipboard paste + message history (BUG HERE)
2. `backend/internal/handler/upload.go` - Image upload endpoint
3. `frontend/src/components/chat/MessageBubble.tsx` - Simplified streaming, no CollapsibleList
4. `frontend/src/hooks/useMessageHistory.ts` - Up arrow navigation
5. `frontend/src/components/savings/CostSavingsCard.tsx` - Cost savings display

## New Features Added

### Image Upload (Phase 1)
- `POST /api/projects/:id/upload` - multipart form upload
- Supports PNG, JPG, GIF, WebP (max 10MB)
- Claude Vision converts to markdown
- Stored in `sources/` folder with "Source Materials" group
- Migration: `008_file_sources.sql`

### Chat Input Enhancements
- **Clipboard paste**: Ctrl/Cmd+V pastes screenshots directly
- **Up arrow history**: Cycles through previous messages (50 max, per-project)
- Preview before sending with remove button

### Cost Savings
- Shows PM + Dev time equivalents
- Integrated into DiscoverySummaryModal
- Uses message count to estimate value

## Changes Summary

### Chat UI
- Removed "Show raw output" toggle during streaming
- Content renders normally while streaming (no separate path)
- Streaming indicator dots at end of message
- Progressive disclosure disabled during streaming
- Removed CollapsibleList (lists render fully now)
- Auto-scroll works during streaming

### File Explorer
- Removed "Descriptions" card view
- Only "By Purpose" and "Files" (tree) remain
- Default is "By Purpose"

## Git Status

All changes committed and pushed to main.

## What's Next

1. **FIX BUG**: Clipboard paste not sending image (see above)

2. **Then test**:
   - Up arrow message history
   - Cost savings in summary modal

3. **Pending PRD work**:
   - Phase 2: PDF & DOCX upload
   - Collapsible sections by heading (user requested)
   - File picker button + drag-drop UI for chat

4. **Run migration** if not done:
   ```bash
   psql -d your_db -f backend/migrations/008_file_sources.sql
   ```

## Session History

- SESSION-013: Chat UI polish, collapsible code, zoom, file metadata
- SESSION-014: Image upload, clipboard paste, message history, cost savings, streaming fixes
