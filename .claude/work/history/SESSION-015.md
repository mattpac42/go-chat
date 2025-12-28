# Session 015 - Image Upload & Chat UX Improvements

## Date
2025-12-28

## Summary
Fixed clipboard paste bug for image uploads and added comprehensive file attachment support to chat input.

## Completed Tasks

### 1. Fixed Clipboard Paste Bug
- **Root cause**: Upload was going to Next.js (port 3001) instead of backend (port 8081), causing 404
- **Additional issue**: Blob URL was revoked prematurely due to React 18 Strict Mode double-mounting
- **Solution**:
  - Created Next.js API proxy route at `/api/projects/[id]/upload/route.ts`
  - Changed blob URL management from ref to state
  - Added proper error handling for upload failures

### 2. Added File Attachment Support
- Paperclip button to open file picker
- Drag-and-drop support with visual feedback (teal highlight)
- Supports: PNG, JPG, JPEG, GIF, WebP
- Reuses same upload flow as clipboard paste

### 3. Added Inline Title Editing
- Click project title in chat header to edit
- Enter to save, Escape to cancel, blur also saves
- Auto-focuses and selects text
- Updates sidebar project list automatically

## Files Modified

### Frontend
- `src/components/chat/ChatInput.tsx` - File attachment, drag-drop, fixed blob URL handling
- `src/components/chat/ChatContainer.tsx` - Inline title editing
- `src/components/ProjectPageClient.tsx` - Wire up title update handler
- `src/app/api/projects/[id]/upload/route.ts` - NEW: API proxy for uploads

## Technical Notes

### Upload Flow (Fixed)
1. Browser -> `/api/projects/:id/upload` (same origin, no CORS)
2. Next.js API route -> `localhost:8081/api/projects/:id/upload`
3. Backend processes with Claude Vision (~12 seconds)
4. Response proxied back through Next.js
5. Image markdown content appended to chat message

### Blob URL Fix
- Changed from `useRef` to `useState` for `pendingImageUrl`
- `ImagePreview` component now receives URL as prop instead of creating its own
- Prevents React Strict Mode from revoking URL during double-mount

## Git Status
```
 M src/components/ProjectPageClient.tsx
 M src/components/chat/ChatContainer.tsx
 M src/components/chat/ChatInput.tsx
?? src/app/api/
```

## Next Session Work
**Smart file naming for uploads** - Use Claude Vision output to generate intuitive filenames:
- Short description (1-3 words) becomes the filename
- Long description stays as metadata
- Instead of `image-2025-12-27.md`, get `bakery-menu-mockup.md`
