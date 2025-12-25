# Developer Session: Claude Prompt Metadata Generation

**Date**: 2025-12-25T12:35:38
**Agent**: developer
**Task**: Update Claude prompt to generate file metadata (descriptions, functional groups) when creating files

## Work Completed

Implemented the complete metadata generation pipeline for the App Map feature:

1. **Updated Claude System Prompt** (`/workspace/backend/internal/service/claude.go`)
   - Added YAML front matter format requirements to the system prompt
   - Specified required metadata fields: short_description, long_description, functional_group
   - Provided clear examples of correct formatting for Claude to follow

2. **Created YAML Metadata Parser** (`/workspace/backend/internal/pkg/markdown/metadata.go`)
   - `FileMetadata` struct for storing parsed metadata
   - `ParseMetadataFromContent()` function extracts YAML front matter from file content
   - `ExtractCodeBlocksWithMetadata()` combines code block extraction with metadata parsing
   - Handles edge cases: no metadata, invalid YAML, empty front matter, Windows line endings

3. **Integrated Metadata into Chat Service** (`/workspace/backend/internal/service/chat.go`)
   - Added `fileMetadataRepo` dependency to ChatService
   - Updated `ProcessMessage()` to use `ExtractCodeBlocksWithMetadata()`
   - Files are saved with metadata via `fileMetadataRepo.Upsert()` when metadata is present
   - Created `inferFilenamesFromUserMessageWithMetadata()` for metadata-aware filename inference

4. **Updated Server Wiring** (`/workspace/backend/cmd/server/main.go`)
   - Added `fileMetadataRepo` initialization and injection

5. **Added Mock Repository** (`/workspace/backend/internal/repository/mock.go`)
   - `MockFileRepository` for testing file operations

6. **Comprehensive Unit Tests**
   - `/workspace/backend/internal/pkg/markdown/metadata_test.go`: Parser tests
   - `/workspace/backend/internal/service/chat_test.go`: Integration tests for metadata extraction

## Decisions Made

- **YAML front matter format**: Chosen because it is a familiar pattern (common in static site generators), easy to parse, and Claude handles YAML generation well
- **Graceful degradation**: If metadata parsing fails, the file is still saved without metadata (no errors thrown)
- **Upsert behavior**: Metadata is upserted (create or update) to handle file regeneration scenarios

## Files Modified

- `/workspace/backend/internal/service/claude.go`: Updated system prompt with metadata format
- `/workspace/backend/internal/pkg/markdown/metadata.go`: NEW - YAML metadata parser
- `/workspace/backend/internal/pkg/markdown/metadata_test.go`: NEW - Parser unit tests
- `/workspace/backend/internal/service/chat.go`: Integrated metadata extraction and storage
- `/workspace/backend/internal/service/chat_test.go`: Added metadata extraction tests, fixed test signatures
- `/workspace/backend/internal/repository/mock.go`: Added MockFileRepository
- `/workspace/backend/cmd/server/main.go`: Wired FileMetadataRepository

## Test Coverage

- Parser tests: Valid YAML, no front matter, invalid YAML, empty front matter, partial metadata, leading whitespace, Windows line endings
- Integration tests: File with metadata extraction, file without metadata handling

## Recommendations

1. Run full test suite once Go is available in the environment to verify compilation
2. Consider adding logging for debugging metadata parsing in production
3. Future enhancement: Add validation for functional_group values against known groups
