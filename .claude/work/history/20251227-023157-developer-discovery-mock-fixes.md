# Developer Session: Discovery Mock Service Fixes

**Date**: 2025-12-27 02:31:57
**Agent**: developer
**Task**: Fix several bugs in the discovery mock service and fixtures

## Work Completed

### Bug 1: Duplicate Welcome Messages (Fixed)
- Modified `determineFixtureKey` in `/workspace/backend/internal/service/claude_mock.go`
- Added `hasAssistantMessages()` helper function to check if conversation already has assistant messages
- Changed welcome logic to only return `welcome_response` if no assistant messages exist yet
- This prevents duplicate welcome messages when `GenerateWelcomeMessage` already sent one

### Bug 2: Duplicate Second/Third Messages (Fixed)
- Added `countUserMessages()` helper function to count user messages in conversation
- Changed stage fixture selection to use user message count instead of global `messageCount`
- Each stage now selects initial vs followup fixtures based on how many user messages are in the current conversation
- This prevents overlapping message count thresholds from returning the same fixture multiple times

### Bug 3: Premature Discovery Complete Notification (Fixed)
- Updated `/workspace/backend/testdata/discovery/summary_response.json`
- Changed `stage_complete: true` to `stage_complete: false`
- Only `summary_confirm_response.json` should have `stage_complete: true`
- This prevents the frontend from prematurely showing the discovery complete notification

### Bug 4: Users Field Format Mismatch (Fixed)
Fixed all fixtures that had `users` as an object instead of array:
- `/workspace/backend/testdata/discovery/summary_response.json`
- `/workspace/backend/testdata/discovery/summary_confirm_response.json`
- `/workspace/backend/testdata/discovery/mvp_followup_response.json`
- `/workspace/backend/testdata/discovery/mvp_response.json`
- `/workspace/backend/testdata/discovery/mvp_confirm_response.json`
- `/workspace/backend/testdata/discovery/personas_followup_response.json`
- `/workspace/backend/testdata/discovery/complete_response.json`

Changed from object format:
```json
"users": {"description": "...", "count": 3, ...}
```

To array format expected by `processExtractedData`:
```json
"users": [
  {"description": "Owner/baker", "count": 1, ...},
  {"description": "Employees who take orders", "count": 2, ...}
]
```

## Files Modified

- `/workspace/backend/internal/service/claude_mock.go`:
  - Added `hasAssistantMessages()` helper function
  - Added `countUserMessages()` helper function
  - Rewrote `determineFixtureKey()` logic to use conversation content instead of global message count

- `/workspace/backend/testdata/discovery/summary_response.json`:
  - Set `stage_complete: false`
  - Converted `users` from object to array

- `/workspace/backend/testdata/discovery/summary_confirm_response.json`:
  - Converted `users` from object to array

- `/workspace/backend/testdata/discovery/mvp_followup_response.json`:
  - Converted `users` from object to array

- `/workspace/backend/testdata/discovery/mvp_response.json`:
  - Converted `users` from object to array

- `/workspace/backend/testdata/discovery/mvp_confirm_response.json`:
  - Converted `users` from object to array

- `/workspace/backend/testdata/discovery/personas_followup_response.json`:
  - Converted `users` from object to array

- `/workspace/backend/testdata/discovery/complete_response.json`:
  - Converted `users` from object to array

## Decisions Made

1. **Conversation-based detection over global counter**: The global `messageCount` was unreliable since the mock service is a singleton. Using conversation content analysis is more reliable and reflects actual state.

2. **Separate user types in fixtures**: Split the single "Owner/baker and two employees" user object into two separate user entries - this matches what the `processExtractedData` code expects when iterating over `[]interface{}`.

3. **Summary stage progression**: The `summary_response.json` now sets `stage_complete: false` so the user can see the summary before being prompted to confirm, preventing premature completion notifications.

## Recommendations

1. Consider adding unit tests for the new `hasAssistantMessages()` and `countUserMessages()` helper functions
2. Manual testing of the full discovery flow is recommended to verify the fixes work end-to-end
3. The mock service improvements could be enhanced further by tracking which fixtures were already returned to completely prevent duplicates
