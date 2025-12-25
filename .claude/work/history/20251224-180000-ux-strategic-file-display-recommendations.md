# UX Strategic Session: File Display Flow Recommendations

**Date**: 2025-12-24T18:00:00Z
**Agent**: ux-strategic
**Task**: Provide design direction for file generation UX in code chat application

## Work Completed

### Analysis Conducted
- Reviewed current frontend implementation (HomeClient, FileTree, FilePreviewModal, FilePill, MessageBubble, ChatContainer)
- Referenced previous ux-tactical research session
- Analyzed patterns from Cursor, Bolt.new, v0.dev, Replit

### Deliverables Provided
Comprehensive UX recommendations addressing:
1. File click behavior (slide-over panel / bottom sheet)
2. Code display in chat (file pills with optional expansion)
3. Additional UX enhancements (progressive disclosure, mobile access, status indicators)

## Decisions Made
- **File Preview Pattern**: Keep existing slide-over (desktop) / bottom sheet (mobile) - non-destructive, maintains context
- **Code Display**: Replace inline code blocks with file pills by default - reduces cognitive load for non-technical users
- **Progressive Disclosure**: Three-tier hierarchy (summary > pills > full code) to serve both technical and non-technical users
- **Mobile Strategy**: Add persistent "Files" button in chat header for quick access

## Design Principles Established
1. Hide complexity by default
2. Progressive disclosure for power users
3. Contextual actions where needed
4. Mobile-first patterns
5. Conversational continuity

## Files Analyzed
- `/workspace/frontend/src/components/HomeClient.tsx`
- `/workspace/frontend/src/components/projects/FileTree.tsx`
- `/workspace/frontend/src/components/shared/FilePreviewModal.tsx`
- `/workspace/frontend/src/components/chat/FilePill.tsx`
- `/workspace/frontend/src/components/chat/MessageBubble.tsx`
- `/workspace/frontend/src/components/chat/ChatContainer.tsx`

## Recommendations for Implementation
1. Modify MessageBubble to parse code blocks and render FilePill components
2. Add expand/collapse toggle to code blocks
3. Connect file pills to FilePreviewModal via shared state
4. Add mobile file access button to ChatContainer header
5. Implement syntax highlighting in file preview
6. Add file status indicators (new/modified)

## Handoff Notes
Tactical implementation should focus on P0 items first:
- File pill display in messages (replace inline code)
- File tree sidebar access on mobile
