# Developer Session: Smart File Naming for Image Uploads

**Date**: 2025-12-28 03:38:21
**Agent**: developer
**Task**: Implement smart file naming for image uploads using Claude Vision analysis

## Work Completed

Implemented smart file naming feature that renames uploaded images based on Claude Vision's content analysis instead of using the original filename.

### Changes Made

1. **Updated VisionPrompt** (`/workspace/backend/internal/handler/upload.go` lines 23-33):
   - Modified prompt to request a short descriptive filename (1-3 words, kebab-case) on the first line prefixed with "FILENAME:"
   - Added example response format for Claude to follow

2. **Created `parseVisionResponse` function** (`/workspace/backend/internal/handler/upload.go` lines 254-300):
   - Extracts the filename from the first line after "FILENAME:" prefix
   - Returns remaining content without the FILENAME line
   - Falls back to "image-upload" if parsing fails

3. **Created `sanitizeSmartFilename` function** (`/workspace/backend/internal/handler/upload.go` lines 302-337):
   - Converts to lowercase
   - Replaces spaces and underscores with dashes
   - Removes invalid characters (keeps only alphanumeric and dashes)
   - Truncates to 40 characters
   - Falls back to "image-upload" if result is empty

4. **Updated Upload handler** (`/workspace/backend/internal/handler/upload.go` lines 159-172):
   - Parses vision response to extract smart filename and content
   - Uses smart filename instead of sanitized original filename
   - Added debug logging for filename generation

5. **Added comprehensive tests** (`/workspace/backend/internal/handler/upload_test.go`):
   - 12 table-driven test cases for `parseVisionResponse` function
   - Updated integration tests to validate smart filename behavior
   - Added fallback test when vision response lacks FILENAME prefix

## Decisions Made

- **Default fallback filename**: Used "image-upload" as the fallback when parsing fails, as it's more descriptive than generic alternatives
- **Underscore handling**: Convert underscores to dashes for consistency with kebab-case
- **40-char limit**: Maintained same limit as existing `sanitizeUploadFilename` for consistency

## Files Modified

- `/workspace/backend/internal/handler/upload.go`: Added smart filename parsing and updated handler logic
- `/workspace/backend/internal/handler/upload_test.go`: Added tests for new functionality and updated existing tests

## Recommendations

1. **Run tests in Go environment**: Tests were written but could not be executed in the current devcontainer (Node.js only). Run `go test -v ./internal/handler/...` in a Go environment to verify
2. **Monitor Claude Vision responses**: Ensure Claude consistently follows the FILENAME format; fallback behavior handles edge cases
3. **Consider caching**: If vision API costs are a concern, the filename extraction could be done client-side with a simpler heuristic

## Example Flow

Before:
- User uploads "Screenshot 2025-12-27 at 3.24.23 PM.png"
- File saved as: `sources/screenshot-2025-12-27-at-32423-pm-2025-12-27.md`

After:
- Claude Vision sees a bakery menu mockup
- Response: "FILENAME: bakery-menu-mockup\n\n## Menu Design\n..."
- File saved as: `sources/bakery-menu-mockup-2025-12-27.md`
