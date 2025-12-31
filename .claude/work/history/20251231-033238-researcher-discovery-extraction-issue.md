# Researcher Session: Discovery Data Extraction Issue

**Date**: 2025-12-31 03:32:38
**Agent**: researcher
**Task**: Investigate why project title, summary, and users aren't being extracted/displayed for new projects

## Problem Summary

User reports that in a new project (http://localhost:3001/projects/4b45c8b2-6753-4af7-840b-bfa5a35f2e3a), the following data is not being extracted or displayed:
- Project title/name
- Summary/solves statement
- Users/user types

## Investigation Findings

### 1. Extraction Flow Architecture

The discovery data extraction happens in this flow:

1. **User sends message** → `/workspace/backend/internal/service/chat.go:ProcessMessage()`
2. **Claude responds** → Response streamed back
3. **Extract metadata** → `/workspace/backend/internal/service/discovery.go:ExtractAndSaveData()` (line 447)
4. **Parse DISCOVERY_DATA** → `parseResponseMetadata()` looks for `<!--DISCOVERY_DATA:{...}-->` comment
5. **Process extracted data** → `processExtractedData()` saves to database (line 665)
6. **Strip metadata from display** → `StripMetadata()` removes the HTML comment before showing to user

### 2. Recent Change That May Have Caused Issue

**CRITICAL FINDING**: In `/workspace/frontend/src/components/chat/MessageBubble.tsx`, the `stripDiscoveryMetadata()` function was recently added (lines 26-31):

```typescript
function stripDiscoveryMetadata(content: string): string {
  // Remove the discovery data HTML comment completely
  // Pattern: <!--DISCOVERY_DATA:{...}--> where {...} is JSON (may contain nested braces)
  // Use a non-greedy match that ends at the --> delimiter
  return content.replace(/<!--DISCOVERY_DATA:.*?-->/g, '');
}
```

This is called in `processAssistantContent()` at line 131:
```typescript
let cleanContent = stripDiscoveryMetadata(content);
```

**However**, this is on the FRONTEND only. The extraction happens on the BACKEND in `chat.go` at lines 169-182:

```go
// If in discovery mode, extract and save discovery data from response
if discovery != nil && !discovery.Stage.IsComplete() {
    if err := s.discoveryService.ExtractAndSaveData(ctx, discovery.ID, responseContent); err != nil {
        s.logger.Warn().
            Err(err).
            Str("projectId", projectID.String()).
            Str("discoveryId", discovery.ID.String()).
            Msg("failed to extract discovery data from response")
        // Continue - don't fail the message for discovery extraction errors
    }

    // Strip discovery metadata from response for display
    responseContent = StripMetadata(responseContent)
}
```

**The extraction happens BEFORE the metadata is stripped**, so the frontend change should NOT affect extraction.

### 3. Extraction Logic Analysis

#### Summary Stage Requirements (from `/workspace/backend/internal/service/prompts/discovery.go`)

The summary stage prompt (line 202-273) tells Claude:

```
CRITICAL METADATA REQUIREMENT:
You MUST include this metadata comment at the VERY END of your response:
<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"project_name":"Your Generated Name","solves_statement":"Your one sentence problem statement"}}-->
```

#### Extraction Code (from `/workspace/backend/internal/service/discovery.go`)

**Line 706-717**: The code looks for `project_name` in two places:
```go
// Extract project name (check top-level and nested under "summary")
if pn, ok := extracted["project_name"].(string); ok && pn != "" {
    s.logger.Info().Str("project_name", pn).Msg("found project_name at top level")
    update.ProjectName = &pn
    hasUpdates = true
} else if summary, ok := extracted["summary"].(map[string]interface{}); ok {
    if pn, ok := summary["project_name"].(string); ok && pn != "" {
        s.logger.Info().Str("project_name", pn).Msg("found project_name in summary")
        update.ProjectName = &pn
        hasUpdates = true
    }
}
```

**Line 720-724**: Warning if project_name not found:
```go
// Log if project_name not found in summary stage
if update.ProjectName == nil && discovery.Stage == model.StageSummary {
    s.logger.Warn().
        Interface("extracted", extracted).
        Msg("project_name not found in summary stage metadata")
}
```

#### Users Extraction (Line 744-778)

Users are extracted from the `users` array in the metadata:
```go
// Extract and save users (clear existing first to avoid duplicates)
if usersRaw, ok := extracted["users"].([]interface{}); ok && len(usersRaw) > 0 {
    // Clear existing users before adding new ones
    if err := s.repo.ClearUsers(ctx, discovery.ID); err != nil {
        s.logger.Warn().Err(err).Msg("failed to clear existing users")
    }
    for _, u := range usersRaw {
        if userMap, ok := u.(map[string]interface{}); ok {
            user := &model.DiscoveryUser{
                DiscoveryID: discovery.ID,
            }
            if desc, ok := userMap["description"].(string); ok {
                user.Description = desc
            }
            // ... more fields
        }
    }
}
```

### 4. Fallback Extraction

**Line 476-510**: If metadata parsing fails, the code has fallback regex extraction:

```go
// For summary stage, try fallback extraction if metadata doesn't have project_name
if discovery.Stage == model.StageSummary {
    if metadata.Extracted == nil {
        metadata.Extracted = make(map[string]interface{})
    }

    // Try fallback extraction for project_name
    if _, hasProjectName := metadata.Extracted["project_name"]; !hasProjectName {
        if fallbackName := extractProjectNameFromText(response); fallbackName != "" {
            s.logger.Info().
                Str("fallbackName", fallbackName).
                Msg("extracted project_name from response text (fallback)")
            metadata.Extracted["project_name"] = fallbackName
        }
    }
    // Similar for solves_statement
}
```

## Root Cause Analysis

### Most Likely Causes:

1. **Claude not including metadata comment**: The AI might not be following the prompt instruction to include the `<!--DISCOVERY_DATA:...-->` comment

2. **Metadata parsing failure**: The JSON in the comment might be malformed or not matching expected format

3. **Premature stage advancement**: The discovery might be advancing to summary stage before all required data is collected

4. **Users data not in metadata**: Earlier stages (personas stage) should have extracted users, but they might not be in the metadata for summary stage

### Expected Metadata Format

**Summary Stage** should have:
```json
{
  "stage_complete": true,
  "extracted": {
    "project_name": "Short Name",
    "solves_statement": "One sentence description"
  }
}
```

**Personas Stage** should have included:
```json
{
  "stage_complete": true,
  "extracted": {
    "users": [
      {
        "description": "user type",
        "count": 1,
        "has_permissions": true,
        "permission_notes": "what they can access"
      }
    ]
  }
}
```

## Files Analyzed

### Backend Files:
- `/workspace/backend/internal/service/chat.go` - Main chat processing and extraction trigger
- `/workspace/backend/internal/service/discovery.go` - Extraction logic (ExtractAndSaveData, processExtractedData)
- `/workspace/backend/internal/service/prompts/discovery.go` - Stage prompts with metadata requirements
- `/workspace/backend/internal/handler/discovery.go` - API endpoints for discovery
- `/workspace/backend/internal/repository/discovery.go` - Database operations (GetSummary)

### Frontend Files:
- `/workspace/frontend/src/components/chat/MessageBubble.tsx` - Strips metadata from display
- `/workspace/frontend/src/components/discovery/DiscoverySummaryCard.tsx` - Displays the summary
- `/workspace/frontend/src/hooks/useDiscovery.ts` - Fetches discovery state from API

## Recommended Next Steps

1. **Check backend logs** for the specific project:
   - Look for "extracted discovery metadata" messages
   - Look for "project_name not found in summary stage metadata" warnings
   - Check if extraction is even being triggered

2. **Inspect actual Claude responses** for the project:
   - Query the database to see what messages were saved
   - Check if any messages contain `<!--DISCOVERY_DATA:` comments
   - Verify the JSON structure if present

3. **Check discovery state in database**:
   ```sql
   SELECT * FROM project_discovery WHERE project_id = '4b45c8b2-6753-4af7-840b-bfa5a35f2e3a';
   SELECT * FROM discovery_users WHERE discovery_id = (SELECT id FROM project_discovery WHERE project_id = '4b45c8b2-6753-4af7-840b-bfa5a35f2e3a');
   SELECT * FROM discovery_features WHERE discovery_id = (SELECT id FROM project_discovery WHERE project_id = '4b45c8b2-6753-4af7-840b-bfa5a35f2e3a');
   ```

4. **Test metadata extraction** with a controlled message containing proper metadata format

5. **Review prompt compliance**: The prompt tells Claude to include metadata, but we should verify:
   - Is the system prompt being used correctly?
   - Is Claude consistently following the instruction?
   - Are there any conflicts with other instructions?

## Key Code Locations

- **Extraction trigger**: `/workspace/backend/internal/service/chat.go:171`
- **Metadata parsing**: `/workspace/backend/internal/service/discovery.go:587-601`
- **Data processing**: `/workspace/backend/internal/service/discovery.go:665-837`
- **Fallback extraction**: `/workspace/backend/internal/service/discovery.go:476-510`
- **Project name regex**: `/workspace/backend/internal/service/discovery.go:443-444`
- **Users extraction**: `/workspace/backend/internal/service/discovery.go:744-778`
- **Summary prompt**: `/workspace/backend/internal/service/prompts/discovery.go:202-273`
- **Frontend metadata stripping**: `/workspace/frontend/src/components/chat/MessageBubble.tsx:26-31`

## Conclusion

The extraction system has:
- ✅ Proper metadata parsing
- ✅ Fallback text extraction
- ✅ Comprehensive logging
- ✅ Correct flow (backend extracts before frontend strips)

The issue is likely:
- ❌ Claude not including metadata in responses
- ❌ Metadata format not matching expected structure
- ❌ Data from earlier stages not persisting

**Action Required**: Need to examine actual messages and logs for the specific project to determine which failure point is occurring.
