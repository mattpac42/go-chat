# Developer Session: Discovery Service Implementation

**Date**: 2025-12-26T02:40:30Z
**Agent**: developer
**Task**: Implement Discovery Service and stage-specific prompts for Phase 1 of Guided Discovery Flow

## Work Completed

Implemented the complete Discovery Service infrastructure for the guided discovery flow:

### 1. Created Discovery Prompts (`/workspace/backend/internal/service/prompts/discovery.go`)

- `DiscoveryPromptBuilder` struct with `Build(stage, context)` method
- `DiscoveryContext` struct for passing captured data between stages
- Stage-specific prompts for all 5 stages:
  - `welcomePrompt()` - Warm greeting, ask about business (Stage 1 of 5)
  - `problemPrompt(context)` - Ask about pain points and challenges (Stage 2 of 5)
  - `personasPrompt(context)` - Ask about who will use the app (Stage 3 of 5)
  - `mvpPrompt(context)` - Ask for THREE essential features (Stage 4 of 5)
  - `summaryPrompt(context)` - Generate summary, ask for confirmation (Stage 5 of 5)
- Each prompt includes:
  - Claude's role as "Product Guide"
  - Stage number indicator
  - Style guidelines (warm, no jargon, concise)
  - DO NOT rules (no code, no technical terms)
  - Metadata output format for data extraction
- `GetStageDisplayInfo()` helper for frontend display information

### 2. Created Discovery Service (`/workspace/backend/internal/service/discovery.go`)

Core methods implemented:
- `GetOrCreateDiscovery(ctx, projectID)` - Creates new or returns existing discovery
- `GetDiscovery(ctx, projectID)` - Retrieves discovery by project ID
- `GetDiscoveryByID(ctx, discoveryID)` - Retrieves discovery by ID
- `AdvanceStage(ctx, discoveryID)` - Progresses to next stage with validation
- `UpdateDiscoveryData(ctx, discoveryID, data)` - Updates captured data fields
- `AddUser(ctx, discoveryID, user)` - Adds user persona
- `AddFeature(ctx, discoveryID, feature)` - Adds feature (MVP or future)
- `GetSummary(ctx, discoveryID)` - Returns complete discovery summary
- `ConfirmDiscovery(ctx, discoveryID)` - Marks discovery complete
- `ResetDiscovery(ctx, discoveryID)` - Deletes and recreates discovery
- `GetSystemPrompt(ctx, projectID)` - Returns stage-appropriate prompt for chat integration
- `ExtractAndSaveData(ctx, discoveryID, response)` - Parses Claude response metadata
- `IsDiscoveryMode(ctx, projectID)` - Returns whether project is in discovery mode
- `GetDiscoveryStage(ctx, projectID)` - Returns current stage

Custom error types:
- `ErrInvalidStageTransition`
- `ErrDiscoveryNotFound`
- `ErrDiscoveryAlreadyComplete`

Utility function:
- `StripMetadata(response)` - Removes `<!--DISCOVERY_DATA:...-->` from response for display

### 3. Created Service Tests (`/workspace/backend/internal/service/discovery_test.go`)

Comprehensive test coverage including:
- `TestGetOrCreateDiscovery_CreatesNew` / `_ReturnsExisting`
- `TestAdvanceStage_ProgressesThroughStages`
- `TestAdvanceStage_ErrorsOnComplete` / `_ErrorsOnNotFound`
- `TestGetSystemPrompt_ReturnsStageAppropriatePrompt` / `_ReturnsEmptyForComplete`
- `TestConfirmDiscovery_MarksComplete` / `_ErrorsOnNonSummaryStage`
- `TestUpdateDiscoveryData` / `_ErrorsOnComplete`
- `TestAddUser` / `TestAddFeature`
- `TestGetSummary`
- `TestResetDiscovery`
- `TestIsDiscoveryMode`
- `TestExtractAndSaveData` (basic, with users, with features, no metadata)
- `TestStripMetadata`

## Decisions Made

- **Metadata extraction format**: Used HTML comment format `<!--DISCOVERY_DATA:{...}-->` which can be stripped for user display but parsed by backend
- **Stage progression**: Automatic advancement when `stage_complete: true` in metadata, with validation that prevents skipping stages
- **Error handling**: Custom error types for clear error handling upstream
- **Context building**: `buildPromptContext()` aggregates all captured data for prompt generation

## Files Modified

- `/workspace/backend/internal/service/prompts/discovery.go` (NEW): Discovery prompt builder and stage prompts
- `/workspace/backend/internal/service/discovery.go` (NEW): Discovery service with all methods
- `/workspace/backend/internal/service/discovery_test.go` (NEW): Comprehensive test suite

## Recommendations

1. **Integration with ChatService**: The `GetSystemPrompt()` method should be called in `ChatService.ProcessMessage()` to use discovery prompts when in discovery mode
2. **Response processing**: After receiving Claude responses, call `ExtractAndSaveData()` to parse metadata and advance stages automatically
3. **Display cleanup**: Use `StripMetadata()` before displaying responses to users
4. **Testing**: Run tests once Go runtime is available in the development environment
