# UX Tactical Session: Frontend UX Fixes

**Date**: 2025-12-25T21:23:33
**Agent**: ux-tactical
**Task**: Fix 4 frontend UX issues in Go Chat application

## Work Completed

### 1. Project Title Sync (Issue #1)
Modified `useProject` hook to accept an optional `projects` array parameter and sync title changes in real-time when the projects list is updated after a rename operation.

### 2. Streaming Code UX (Issue #2)
Replaced the raw streaming content display with a "Generating response..." indicator that has bouncing dots. Added an expandable "Show raw output" toggle so users can optionally view streaming content if desired. The final rendered markdown only appears after streaming completes.

### 3. File Tile Click Behavior (Issue #3)
Changed the card body click behavior to toggle only between `collapsed` and `details` tiers. Code view is now only accessible via the dedicated code button (`</>`). The chevron still cycles through all three tiers for power users.

### 4. Markdown Rendering (Issue #4)
Added a `preprocessMarkdown()` function that fixes common markdown parsing issues before rendering:
- Ensures code blocks have proper newlines after opening \`\`\`language
- Ensures code blocks have proper newlines before closing \`\`\`
- Ensures proper spacing after closing code blocks

## Decisions Made

- **Title sync via props**: Chose to pass the `projects` array to `useProject` rather than implementing a global state or event system, keeping changes minimal and localized.
- **Streaming UX**: Opted for a clean "Generating..." message with optional raw output toggle rather than partial markdown rendering, which eliminates flickering and broken markdown during streaming.
- **Two-tier toggle**: Decided card body clicks should only toggle collapsed/details (not code) to reduce accidental code view activation. Code button remains the intentional way to see code.

## Files Modified

- `/workspace/frontend/src/hooks/useProjects.ts`: Added `projects` parameter to `useProject` hook with title sync effect
- `/workspace/frontend/src/components/ProjectPageClient.tsx`: Pass `projects` to `useProject` call
- `/workspace/frontend/src/components/chat/MessageBubble.tsx`: Added streaming UI with expandable raw content, added `preprocessMarkdown()` function
- `/workspace/frontend/src/components/projects/FileRevealCard.tsx`: Split `getNextTier` into `getNextTierFull` (for chevron) and `getNextTierToggle` (for card body)
- `/workspace/frontend/src/components/projects/FileRevealList.tsx`: Updated `handleCardClick` to use two-tier toggle logic, fixed FolderSection props

## Recommendations

- Test the streaming UX with actual AI responses to verify the "Generating..." message displays correctly
- Consider adding syntax highlighting to the raw output view during streaming for better readability
- Monitor user feedback on the two-tier file toggle; may need adjustment based on usage patterns
