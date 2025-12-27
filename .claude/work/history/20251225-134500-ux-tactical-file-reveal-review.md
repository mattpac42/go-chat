# UX Tactical Session: File Reveal Interaction Review

**Date**: 2025-12-25T13:45:00Z
**Agent**: ux-tactical
**Task**: Review 3-tier file reveal interaction in FileRevealCard and propose improvements

## Work Completed

Analyzed the current FileRevealCard component implementation and identified usability issues based on user feedback. Documented the existing state machine, identified four key UX problems, and proposed a linear progressive disclosure model that aligns with natural "show more / show less" thinking.

## Decisions Made

- **Linear progression over independent toggles**: Users think in terms of expanding and collapsing, not managing two separate toggles. The proposed model uses a single axis: collapsed -> details -> code -> collapsed.

- **Keep long description visible with code**: When viewing code, the long description should remain visible to provide context about what the file does. Hiding it removes valuable context.

- **Chevron always visible**: Every file can expand (even if just to code), so the chevron should always be present. It rotates when content is expanded.

- **Code button as shortcut**: The code button provides a quick way to jump directly to code view or collapse from code view, without requiring multiple clicks through the details tier.

## Files Reviewed

- `/workspace/frontend/src/components/projects/FileRevealCard.tsx`: Main component with interaction logic

## Recommendations

### Immediate Implementation Changes

1. **Update handleCardClick** to progress through all three tiers linearly, including collapsing from code state

2. **Update handleViewCode** to always collapse to base state (not to details) for consistency

3. **Remove `!showCode` condition** from long description rendering so it remains visible alongside code

4. **Always render chevron** regardless of hasLongDescription, since all files can expand to code view

5. **Add meaningful aria-labels** that describe the next action based on current state

### New Interaction Pattern

| State | Card Click | Code Button |
|-------|------------|-------------|
| Collapsed | Show details (or code) | Show code |
| Details | Show code | Show code |
| Code | Collapse all | Collapse all |

## Next Steps

Developer should implement the recommended changes to FileRevealCard.tsx. The changes are isolated to the component and do not require backend modifications.
