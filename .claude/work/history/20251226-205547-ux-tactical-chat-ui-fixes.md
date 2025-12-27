# UX Tactical Session: Chat UI Fixes

**Date**: 2025-12-26T20:55:47
**Agent**: ux-tactical
**Task**: Fix several UI issues in the chat interface

## Work Completed

Fixed 5 UI issues in the chat interface:

1. **Summary card overlapping messages**: Added `hasBottomCard` prop to MessageList that reduces bottom padding when summary card is visible, ensuring messages don't get hidden behind the card.

2. **User message text color**: Added explicit `text-white` class along with `prose-li:text-white`, `prose-ol:text-white`, `prose-ul:text-white` to ensure all text in user message bubbles renders white on the teal background.

3. **Numbered list formatting**: Changed from `list-inside` to `list-outside` with `pl-5` padding for both `ul` and `ol` elements. Added `li` component with proper styling to ensure list numbers/bullets align correctly with their content.

4. **Summary card accessibility**: Added collapsible summary card functionality with "Collapse" button when expanded and "Show Project Summary" button when collapsed. Uses chevron icons to indicate state.

5. **Button label clarity**: Changed "Edit Details" to "Start Over" on DiscoverySummaryCard to better communicate the destructive action of resetting discovery.

## Decisions Made

- **Collapsible summary card**: Chose to add collapse/expand functionality rather than dismiss-on-click to ensure users always have access to the summary when in summary stage
- **List styling**: Switched from `list-inside` to `list-outside` because `list-inside` causes rendering issues with markdown where the number appears on a separate line from content
- **Prose text colors**: Added multiple prose color utilities because prose-invert alone wasn't sufficient to override the default gray colors

## Files Modified

- `/workspace/frontend/src/components/chat/ChatContainer.tsx`: Added `hasBottomCard` prop to MessageList, implemented collapsible summary card with chevron icons
- `/workspace/frontend/src/components/chat/MessageList.tsx`: Added `hasBottomCard` prop and conditional padding
- `/workspace/frontend/src/components/chat/MessageBubble.tsx`: Fixed prose text colors for user messages, fixed list styling from list-inside to list-outside with proper padding
- `/workspace/frontend/src/components/discovery/DiscoverySummaryCard.tsx`: Changed "Edit Details" button to "Start Over"

## Recommendations

- Consider adding animation to the summary card collapse/expand transition
- The pre-existing TypeScript errors in test files should be addressed separately
- The demo/discovery page has a prerendering error that should be investigated
