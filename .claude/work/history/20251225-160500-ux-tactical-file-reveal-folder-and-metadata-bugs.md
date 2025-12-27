# UX Tactical Session: File Reveal Folder and Metadata Bugs

**Date**: 2025-12-25 16:05:00
**Agent**: ux-tactical
**Task**: Fix two UI bugs - folder file expansion and AI metadata display

## Work Completed

### Bug 1: Files in folders don't respond to chevron/code clicks

**Problem**: In `FolderSection`, `FileRevealCard` was rendered with a controlled `tier` prop but without `onIntentionalExpand`. In controlled mode, the internal `setTier` does nothing, so chevron and code button clicks had no effect for files inside folders.

**Solution**:
1. Added three new props to `FolderSectionProps` interface:
   - `getFileTier: (fileId: string) => RevealTier`
   - `onCardClick: (fileId: string, hasLongDesc: boolean) => void`
   - `onIntentionalExpand: (fileId: string, tier: RevealTier) => void`
2. Updated `FolderSection` component to destructure and use these props
3. Passed handlers to `FileRevealCard` inside `FolderSection` (matching `FunctionalGroupSection` pattern)
4. Updated recursive `FolderSection` calls to propagate these handlers
5. Updated parent calls in `FileRevealList` to provide the handlers

### Bug 2: AI metadata showing as raw text in messages

**Problem**: Claude's responses included inline code with metadata format like:
```
`javascript:server.js --- short_description: ... long_description: ... functional_group: ... --- const express = require...`
```
This rendered as raw bold text because the metadata wasn't being stripped.

**Solution**:
Added regex to `preprocessMarkdown()` to handle inline code with metadata format:
```javascript
processed = processed.replace(
  /`(\w+):([^\s`]+)\s+---\s+[\s\S]*?---\s+([\s\S]*?)`/g,
  '`$3`'
);
```
This strips the language:filename prefix and metadata section, keeping only the actual code.

## Decisions Made
- Mirrored `FunctionalGroupSection` pattern for consistency since it already correctly passes all handlers
- Used regex capture groups to preserve only the actual code content, discarding language hints and metadata

## Files Modified
- `/workspace/frontend/src/components/projects/FileRevealList.tsx`: Added props to FolderSectionProps interface, updated FolderSection to use handlers, passed handlers to FileRevealCard
- `/workspace/frontend/src/components/chat/MessageBubble.tsx`: Added inline code metadata stripping regex to preprocessMarkdown()

## Recommendations
- Test both folder view and functional group view to verify expand/collapse works consistently
- Test message rendering with various AI responses containing inline code metadata
- Consider adding unit tests for the preprocessMarkdown function edge cases
