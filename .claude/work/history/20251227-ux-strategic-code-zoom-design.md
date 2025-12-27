# UX Strategic Session: Code Zoom Feature Design

**Date**: 2025-12-27
**Agent**: ux-strategic
**Task**: Design recommendation for code zoom feature in file explorer

## Work Completed

Analyzed existing code viewing components (FilePreviewModal, FileRevealCard, FileRevealList, FileExplorer) and provided strategic UX design recommendations for implementing a code zoom feature.

## Decisions Made

- **Zoom method**: Preset levels (75/100/125/150%) with keyboard shortcuts (Cmd/Ctrl +/-) - Provides predictable, accessible options without over-complexity
- **Control placement**: Code toolbar integration - Contextual, always visible when relevant, grouped with related actions
- **Scope**: Global preference with localStorage persistence - Reduces cognitive overhead vs per-file memory
- **Avoid slider/scroll zoom**: Too imprecise or accident-prone for code viewing

## Files Analyzed

- `/workspace/frontend/src/components/projects/FileExplorer.tsx`: Main file explorer with tree/reveal/grouped views
- `/workspace/frontend/src/components/shared/FilePreviewModal.tsx`: Side panel (desktop) / bottom sheet (mobile) for code preview
- `/workspace/frontend/src/components/projects/FileRevealCard.tsx`: Inline expandable code cards with 3-tier reveal
- `/workspace/frontend/src/components/projects/FileRevealList.tsx`: File list with folder structure and functional grouping

## Deliverables

### Design Recommendations

1. **Primary interaction**: 4 preset zoom buttons (75%, 100%, 125%, 150%)
2. **Secondary interaction**: Keyboard shortcuts matching browser/IDE conventions
3. **Placement**: Existing code toolbar areas in both FileRevealCard and FilePreviewModal
4. **Persistence**: Global localStorage preference, not per-file

### Accessibility Requirements

- Full keyboard navigation for zoom controls
- ARIA labels with descriptive text
- Screen reader announcements for zoom changes
- Maintain WCAG 2.1 contrast ratios at all zoom levels
- Proportional line-height scaling

### Implementation Guidance for Tactical Designer

1. Create `useCodeZoom.ts` hook for state management and localStorage
2. Create `ZoomControls.tsx` compact segmented button component
3. Integrate into FileRevealCard toolbar (line ~503)
4. Integrate into FilePreviewModal header (line ~220)
5. Apply zoom via CSS variable or inline style on `<pre>` elements

### Font Size Mapping

| Level | Size |
|-------|------|
| 75% | 11px |
| 100% | 14px (current default) |
| 125% | 18px |
| 150% | 21px |

## Recommendations

- Implementation is straightforward and can proceed with tactical designer
- Consider adding a "Reset to default" action for accessibility
- Monitor for horizontal overflow issues at 150% zoom on narrow viewports
