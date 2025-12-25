# Architect Session: App Map Architecture Design

**Date**: 2025-12-25T12:24:26Z
**Agent**: architect
**Task**: Design the App Map system architecture for organizing files by PURPOSE rather than file paths

## Work Completed

Designed and documented the complete App Map system architecture including:

1. **Database Schema Design**
   - New `functional_groups` table for purpose-based groupings
   - Extended `files` table with metadata columns (short_description, long_description, purpose, visibility_tier, ai_confidence)
   - New `file_relationships` table for graph view connections
   - Default functional groups seed data (homepage, forms, backend, database, etc.)

2. **Backend API Design**
   - New models in `/workspace/backend/internal/model/appmap.go`
   - App Map endpoints supporting three view levels
   - File metadata update endpoints
   - Relationship management endpoints

3. **Claude Integration Design**
   - Enhanced system prompt for metadata generation via YAML front matter
   - Metadata parser for extracting descriptions and group assignments
   - Integration with existing chat service flow

4. **Data Flow Architecture**
   - Sequence diagram from user message to App Map display
   - WebSocket event types for real-time updates

## Decisions Made

- **ADR-001**: Store file metadata as columns in existing `files` table (not separate table) - simpler queries, always accessed together
- **ADR-002**: Functional groups are per-project with template seeding - allows customization per project type
- **ADR-003**: YAML front matter for Claude metadata generation - familiar pattern, easy parsing
- **ADR-004**: Fixed three-level view progression matching product vision learning journey

## Files Created

- `/workspace/.claude/work/design/app-map-architecture.md`: Complete technical design document with:
  - Database schema (SQL migrations)
  - Entity relationship diagrams
  - API endpoint specifications with examples
  - Go model definitions
  - Handler implementation patterns
  - Integration points with existing system
  - Migration strategy
  - Testing approach
  - Effort estimates (~26 hours total)

## Recommendations

1. **Immediate Next Steps**:
   - Create database migration file `003_app_map_schema.sql`
   - Implement backend repository and handler for App Map
   - Update Claude system prompt for metadata generation
   - Build frontend AppMap component with view level toggle

2. **Key Implementation Order**:
   - Phase 1: Database schema (2h)
   - Phase 2: Backend APIs (10h)
   - Phase 3: Claude integration (5h)
   - Phase 4: Frontend components (8h)

3. **Open Questions for Product**:
   - AI fallback strategy when metadata generation fails
   - User permissions for editing AI-generated metadata
   - Relationship validation approach
