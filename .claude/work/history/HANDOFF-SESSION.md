# Handoff - Session 015

## Immediate Context

Working on **Chat UX & Image Upload improvements**. Fixed clipboard paste bug, added file attachment button and drag-drop, added inline title editing.

## Branch

`main` - Changes not yet committed

## Next Task: Smart File Naming for Uploads

**User Request**: When uploading images/screenshots, rename them based on Claude Vision's analysis instead of generic names like "image-2025-12-27.md".

**Proposed Approach**:
1. Modify Claude Vision prompt to extract:
   - Short name (1-3 words) - becomes filename
   - Long description - becomes metadata
2. Update `backend/internal/handler/upload.go`:
   - Parse Claude's response for short name
   - Sanitize and use as filename
   - Store long description in file_metadata
3. Update Vision prompt in `VisionPrompt` constant

**Example**:
- Before: `sources/image-2025-12-27.md`
- After: `sources/bakery-menu-mockup.md`

## Last Tasks Completed

1. Fixed clipboard paste bug (Next.js API proxy)
2. Added file attachment button (paperclip icon)
3. Added drag-and-drop support
4. Added inline title editing in chat header

## Critical Files to Read

1. `backend/internal/handler/upload.go` - Current upload logic, Vision prompt
2. `backend/internal/service/claude_vision.go` - Claude Vision API integration
3. `backend/internal/repository/file_metadata.go` - Metadata storage

## Uncommitted Changes

```
 M frontend/src/components/ProjectPageClient.tsx
 M frontend/src/components/chat/ChatContainer.tsx
 M frontend/src/components/chat/ChatInput.tsx
?? frontend/src/app/api/projects/[id]/upload/route.ts
```

## Suggested First Action

1. Commit current changes first (image upload fixes + inline title editing)
2. Then implement smart file naming in upload handler

## Session History

- SESSION-014: Image upload, clipboard paste (bug discovered)
- SESSION-015: Fixed paste bug, file attachment, drag-drop, inline title edit
