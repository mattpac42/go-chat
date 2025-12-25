# App Map Architecture Design

**Date**: 2025-12-25
**Author**: Architect Agent
**Status**: Draft
**Related Vision**: PRODUCT_VISION.md Theme 2 - Conversational Development Experience

---

## Executive Summary

This document defines the technical architecture for the App Map system, which organizes files by PURPOSE rather than file paths. The App Map provides non-technical users with an intuitive mental model of their application through functional groupings like "Homepage", "Contact Form", and "Backend Services".

The design extends the existing file storage system with metadata and grouping capabilities while maintaining backward compatibility with the current file tree display.

---

## 1. Database Schema Design

### 1.1 New Tables

#### `functional_groups` Table

Stores the functional groupings that organize files by purpose.

```sql
-- 003_app_map_schema.sql
-- App Map: Functional groups and file metadata for purpose-based organization

-- Functional groups table
CREATE TABLE IF NOT EXISTS functional_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(150) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    color VARCHAR(7),
    sort_order INTEGER DEFAULT 0,
    parent_group_id UUID REFERENCES functional_groups(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Ensure unique group names within a project
    CONSTRAINT unique_group_name_per_project UNIQUE (project_id, name)
);

-- Indexes for functional_groups
CREATE INDEX IF NOT EXISTS idx_functional_groups_project_id ON functional_groups(project_id);
CREATE INDEX IF NOT EXISTS idx_functional_groups_parent ON functional_groups(parent_group_id);

COMMENT ON TABLE functional_groups IS 'Purpose-based groupings for App Map visualization';
COMMENT ON COLUMN functional_groups.name IS 'Machine-readable identifier (e.g., homepage, contact_form)';
COMMENT ON COLUMN functional_groups.display_name IS 'Human-readable name (e.g., Homepage, Contact Form)';
COMMENT ON COLUMN functional_groups.icon IS 'Icon identifier for UI display';
COMMENT ON COLUMN functional_groups.color IS 'Hex color code for visual distinction';
COMMENT ON COLUMN functional_groups.sort_order IS 'Display order within parent group';
COMMENT ON COLUMN functional_groups.parent_group_id IS 'For nested group hierarchies';
```

#### Extended `files` Table

Add metadata columns to existing files table for App Map integration.

```sql
-- Add App Map metadata columns to files table
ALTER TABLE files
    ADD COLUMN IF NOT EXISTS functional_group_id UUID REFERENCES functional_groups(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS short_description VARCHAR(200),
    ADD COLUMN IF NOT EXISTS long_description TEXT,
    ADD COLUMN IF NOT EXISTS purpose VARCHAR(100),
    ADD COLUMN IF NOT EXISTS visibility_tier SMALLINT DEFAULT 1 CHECK (visibility_tier IN (1, 2)),
    ADD COLUMN IF NOT EXISTS ai_confidence DECIMAL(3,2),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Index for group lookups
CREATE INDEX IF NOT EXISTS idx_files_functional_group ON files(functional_group_id);
CREATE INDEX IF NOT EXISTS idx_files_purpose ON files(purpose);

COMMENT ON COLUMN files.functional_group_id IS 'Reference to functional group in App Map';
COMMENT ON COLUMN files.short_description IS 'One-line human-readable description of file purpose';
COMMENT ON COLUMN files.long_description IS 'Detailed explanation of what the file does';
COMMENT ON COLUMN files.purpose IS 'High-level purpose category (e.g., ui, logic, data, config)';
COMMENT ON COLUMN files.visibility_tier IS '1=show description only, 2=show code';
COMMENT ON COLUMN files.ai_confidence IS 'AI confidence score for auto-assigned metadata (0.00-1.00)';
```

#### `file_relationships` Table

Stores relationships between files for graph view visualization.

```sql
-- File relationships for graph view
CREATE TABLE IF NOT EXISTS file_relationships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    source_file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    target_file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    relationship_type VARCHAR(50) NOT NULL,
    description TEXT,
    ai_confidence DECIMAL(3,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Prevent duplicate relationships
    CONSTRAINT unique_file_relationship UNIQUE (source_file_id, target_file_id, relationship_type)
);

-- Indexes for relationship queries
CREATE INDEX IF NOT EXISTS idx_file_relationships_project ON file_relationships(project_id);
CREATE INDEX IF NOT EXISTS idx_file_relationships_source ON file_relationships(source_file_id);
CREATE INDEX IF NOT EXISTS idx_file_relationships_target ON file_relationships(target_file_id);

COMMENT ON TABLE file_relationships IS 'Connections between files for graph visualization';
COMMENT ON COLUMN file_relationships.relationship_type IS 'Type: imports, calls, renders, uses_data, configures';
COMMENT ON COLUMN file_relationships.ai_confidence IS 'AI confidence in relationship detection';
```

### 1.2 Default Functional Groups

Seed data for common functional groups:

```sql
-- Seed default functional groups (inserted per-project on creation)
-- These are template groups to be copied to new projects

CREATE TABLE IF NOT EXISTS default_functional_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(150) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    color VARCHAR(7),
    sort_order INTEGER DEFAULT 0,
    parent_name VARCHAR(100)
);

INSERT INTO default_functional_groups (name, display_name, description, icon, color, sort_order) VALUES
    ('frontend', 'Frontend', 'User interface and client-side code', 'monitor', '#3B82F6', 1),
    ('homepage', 'Homepage', 'Main landing page components', 'home', '#10B981', 2),
    ('forms', 'Forms', 'User input and data collection', 'edit', '#8B5CF6', 3),
    ('navigation', 'Navigation', 'Menus, headers, and routing', 'menu', '#F59E0B', 4),
    ('backend', 'Backend Services', 'Server-side logic and APIs', 'server', '#EF4444', 5),
    ('database', 'Database', 'Data models and storage', 'database', '#06B6D4', 6),
    ('authentication', 'Authentication', 'Login, signup, and user sessions', 'lock', '#EC4899', 7),
    ('configuration', 'Configuration', 'Settings and environment setup', 'settings', '#6B7280', 8),
    ('utilities', 'Utilities', 'Helper functions and shared code', 'tool', '#84CC16', 9);
```

### 1.3 Entity Relationship Diagram

```
+------------------+       +----------------------+       +------------------+
|    projects      |       | functional_groups    |       |      files       |
+------------------+       +----------------------+       +------------------+
| id (PK)          |<------| project_id (FK)      |       | id (PK)          |
| title            |       | id (PK)              |<------| functional_group_id (FK)
| created_at       |       | name                 |       | project_id (FK)  |
| updated_at       |       | display_name         |       | path             |
+------------------+       | description          |       | filename         |
                           | icon                 |       | language         |
                           | color                |       | content          |
                           | sort_order           |       | short_description|
                           | parent_group_id (FK) |----+  | long_description |
                           | created_at           |    |  | purpose          |
                           | updated_at           |    |  | visibility_tier  |
                           +----------------------+    |  | ai_confidence    |
                                     ^                 |  | created_at       |
                                     |                 |  | updated_at       |
                                     +-----------------+  +------------------+
                                                                   |
                           +----------------------+                |
                           | file_relationships   |<---------------+
                           +----------------------+
                           | id (PK)              |
                           | project_id (FK)      |
                           | source_file_id (FK)  |
                           | target_file_id (FK)  |
                           | relationship_type    |
                           | description          |
                           | ai_confidence        |
                           | created_at           |
                           +----------------------+
```

---

## 2. Backend API Design

### 2.1 New Models

**File**: `/workspace/backend/internal/model/appmap.go`

```go
package model

import (
    "time"
    "github.com/google/uuid"
)

// FunctionalGroup represents a purpose-based grouping in the App Map.
type FunctionalGroup struct {
    ID            uuid.UUID  `db:"id" json:"id"`
    ProjectID     uuid.UUID  `db:"project_id" json:"projectId"`
    Name          string     `db:"name" json:"name"`
    DisplayName   string     `db:"display_name" json:"displayName"`
    Description   *string    `db:"description" json:"description,omitempty"`
    Icon          *string    `db:"icon" json:"icon,omitempty"`
    Color         *string    `db:"color" json:"color,omitempty"`
    SortOrder     int        `db:"sort_order" json:"sortOrder"`
    ParentGroupID *uuid.UUID `db:"parent_group_id" json:"parentGroupId,omitempty"`
    CreatedAt     time.Time  `db:"created_at" json:"createdAt"`
    UpdatedAt     time.Time  `db:"updated_at" json:"updatedAt"`
}

// FunctionalGroupWithFiles includes files belonging to the group.
type FunctionalGroupWithFiles struct {
    FunctionalGroup
    Files []FileWithMetadata `json:"files"`
    Children []FunctionalGroupWithFiles `json:"children,omitempty"`
}

// FileWithMetadata extends File with App Map metadata.
type FileWithMetadata struct {
    ID                uuid.UUID  `db:"id" json:"id"`
    ProjectID         uuid.UUID  `db:"project_id" json:"projectId"`
    Path              string     `db:"path" json:"path"`
    Filename          string     `db:"filename" json:"filename"`
    Language          *string    `db:"language" json:"language,omitempty"`
    ShortDescription  *string    `db:"short_description" json:"shortDescription,omitempty"`
    LongDescription   *string    `db:"long_description" json:"longDescription,omitempty"`
    Purpose           *string    `db:"purpose" json:"purpose,omitempty"`
    VisibilityTier    int        `db:"visibility_tier" json:"visibilityTier"`
    FunctionalGroupID *uuid.UUID `db:"functional_group_id" json:"functionalGroupId,omitempty"`
    AIConfidence      *float64   `db:"ai_confidence" json:"aiConfidence,omitempty"`
    CreatedAt         time.Time  `db:"created_at" json:"createdAt"`
    UpdatedAt         time.Time  `db:"updated_at" json:"updatedAt"`
}

// FileRelationship represents a connection between two files.
type FileRelationship struct {
    ID               uuid.UUID `db:"id" json:"id"`
    ProjectID        uuid.UUID `db:"project_id" json:"projectId"`
    SourceFileID     uuid.UUID `db:"source_file_id" json:"sourceFileId"`
    TargetFileID     uuid.UUID `db:"target_file_id" json:"targetFileId"`
    RelationshipType string    `db:"relationship_type" json:"relationshipType"`
    Description      *string   `db:"description" json:"description,omitempty"`
    AIConfidence     *float64  `db:"ai_confidence" json:"aiConfidence,omitempty"`
    CreatedAt        time.Time `db:"created_at" json:"createdAt"`
}

// AppMapResponse is the full App Map structure for a project.
type AppMapResponse struct {
    ProjectID     uuid.UUID                  `json:"projectId"`
    ViewLevel     int                        `json:"viewLevel"`
    Groups        []FunctionalGroupWithFiles `json:"groups"`
    Relationships []FileRelationship         `json:"relationships,omitempty"`
    UngroupedFiles []FileWithMetadata        `json:"ungroupedFiles,omitempty"`
}

// FileMetadataUpdate is used when Claude assigns metadata to a file.
type FileMetadataUpdate struct {
    ShortDescription  string   `json:"shortDescription"`
    LongDescription   string   `json:"longDescription"`
    Purpose           string   `json:"purpose"`
    FunctionalGroup   string   `json:"functionalGroup"`
    RelatedFiles      []string `json:"relatedFiles,omitempty"`
    AIConfidence      float64  `json:"aiConfidence"`
}
```

### 2.2 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/projects/:id/appmap` | Get full App Map for project |
| GET | `/api/projects/:id/appmap?level=1\|2\|3` | Get App Map at specific view level |
| GET | `/api/projects/:id/groups` | List all functional groups |
| POST | `/api/projects/:id/groups` | Create a functional group |
| PUT | `/api/groups/:id` | Update a functional group |
| DELETE | `/api/groups/:id` | Delete a functional group |
| PUT | `/api/files/:id/metadata` | Update file App Map metadata |
| PUT | `/api/files/:id/group` | Assign file to a functional group |
| GET | `/api/projects/:id/relationships` | Get all file relationships |
| POST | `/api/projects/:id/relationships` | Create file relationship |

### 2.3 Response Examples

#### GET `/api/projects/:id/appmap?level=1`

Level 1 (Functional Only) - For beginners:

```json
{
  "projectId": "550e8400-e29b-41d4-a716-446655440000",
  "viewLevel": 1,
  "groups": [
    {
      "id": "group-1",
      "name": "homepage",
      "displayName": "Homepage",
      "description": "The main landing page that visitors see first",
      "icon": "home",
      "color": "#10B981",
      "sortOrder": 1,
      "files": []
    },
    {
      "id": "group-2",
      "name": "contact_form",
      "displayName": "Contact Form",
      "description": "Allows visitors to send you messages",
      "icon": "edit",
      "color": "#8B5CF6",
      "sortOrder": 2,
      "files": []
    }
  ]
}
```

#### GET `/api/projects/:id/appmap?level=2`

Level 2 (Functional + Files) - Intermediate:

```json
{
  "projectId": "550e8400-e29b-41d4-a716-446655440000",
  "viewLevel": 2,
  "groups": [
    {
      "id": "group-1",
      "name": "homepage",
      "displayName": "Homepage",
      "description": "The main landing page that visitors see first",
      "icon": "home",
      "color": "#10B981",
      "sortOrder": 1,
      "files": [
        {
          "id": "file-1",
          "filename": "index.html",
          "shortDescription": "The main structure of your homepage",
          "purpose": "ui",
          "visibilityTier": 1
        },
        {
          "id": "file-2",
          "filename": "styles.css",
          "shortDescription": "Visual styling for the homepage",
          "purpose": "ui",
          "visibilityTier": 1
        }
      ]
    }
  ]
}
```

#### GET `/api/projects/:id/appmap?level=3`

Level 3 (Full Technical) - Advanced, includes file paths:

```json
{
  "projectId": "550e8400-e29b-41d4-a716-446655440000",
  "viewLevel": 3,
  "groups": [
    {
      "id": "group-1",
      "name": "homepage",
      "displayName": "Homepage",
      "files": [
        {
          "id": "file-1",
          "path": "src/pages/index.html",
          "filename": "index.html",
          "language": "html",
          "shortDescription": "The main structure of your homepage",
          "longDescription": "This HTML file defines the structure of your landing page including the header, hero section, features list, and footer.",
          "purpose": "ui",
          "visibilityTier": 2
        }
      ]
    }
  ],
  "relationships": [
    {
      "sourceFileId": "file-1",
      "targetFileId": "file-2",
      "relationshipType": "imports",
      "description": "The homepage imports styles from the CSS file"
    }
  ]
}
```

### 2.4 Handler Implementation

**File**: `/workspace/backend/internal/handler/appmap.go`

```go
package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// AppMapHandler handles App Map endpoints.
type AppMapHandler struct {
    appMapRepo    repository.AppMapRepository
    projectRepo   repository.ProjectRepository
}

// NewAppMapHandler creates a new AppMapHandler.
func NewAppMapHandler(appMapRepo repository.AppMapRepository, projectRepo repository.ProjectRepository) *AppMapHandler {
    return &AppMapHandler{
        appMapRepo:  appMapRepo,
        projectRepo: projectRepo,
    }
}

// GetAppMap returns the App Map for a project.
// GET /api/projects/:id/appmap
func (h *AppMapHandler) GetAppMap(c *gin.Context) {
    projectID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
        return
    }

    // Parse view level (default to 2)
    level := 2
    if levelParam := c.Query("level"); levelParam != "" {
        if parsed, err := strconv.Atoi(levelParam); err == nil && parsed >= 1 && parsed <= 3 {
            level = parsed
        }
    }

    appMap, err := h.appMapRepo.GetAppMap(c.Request.Context(), projectID, level)
    if err != nil {
        if err == repository.ErrNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get app map"})
        return
    }

    c.JSON(http.StatusOK, appMap)
}

// UpdateFileMetadata updates App Map metadata for a file.
// PUT /api/files/:id/metadata
func (h *AppMapHandler) UpdateFileMetadata(c *gin.Context) {
    fileID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
        return
    }

    var update model.FileMetadataUpdate
    if err := c.ShouldBindJSON(&update); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    file, err := h.appMapRepo.UpdateFileMetadata(c.Request.Context(), fileID, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update metadata"})
        return
    }

    c.JSON(http.StatusOK, file)
}
```

---

## 3. Data Flow: Claude File Creation to App Map

### 3.1 Enhanced System Prompt

The Claude system prompt must be extended to generate App Map metadata:

```go
const appMapSystemPrompt = `You are Go Chat. You create files for users.

FILE FORMAT REQUIREMENT:
Every code block MUST include metadata in YAML front matter format:

` + "```" + `html:index.html
---
short_description: Main homepage structure with navigation and hero section
long_description: This HTML file defines the complete structure of the landing page including a responsive navigation bar, hero section with call-to-action, feature highlights grid, testimonials carousel, and footer with contact information.
purpose: ui
functional_group: homepage
related_files:
  - styles.css
  - app.js
---
<!DOCTYPE html>
<html>...
` + "```" + `

METADATA FIELDS (all required):
- short_description: One sentence describing what this file does (for Tier 1 display)
- long_description: Detailed explanation of the file's purpose and contents
- purpose: One of [ui, logic, data, config, style, api, test, docs]
- functional_group: One of [homepage, forms, navigation, backend, database, authentication, configuration, utilities] or suggest a new group name
- related_files: List of other files this one connects to

Be concise. Generate working code with complete metadata.`
```

### 3.2 Enhanced Markdown Parser

Extend the existing codeblock parser to extract App Map metadata:

**File**: `/workspace/backend/internal/pkg/markdown/metadata.go`

```go
package markdown

import (
    "regexp"
    "strings"

    "gopkg.in/yaml.v3"
)

// FileMetadata represents App Map metadata extracted from code blocks.
type FileMetadata struct {
    ShortDescription string   `yaml:"short_description"`
    LongDescription  string   `yaml:"long_description"`
    Purpose          string   `yaml:"purpose"`
    FunctionalGroup  string   `yaml:"functional_group"`
    RelatedFiles     []string `yaml:"related_files"`
}

// CodeBlockWithMetadata extends CodeBlock with App Map metadata.
type CodeBlockWithMetadata struct {
    CodeBlock
    Metadata *FileMetadata
}

// ParseMetadataFromContent extracts YAML front matter from file content.
func ParseMetadataFromContent(content string) (*FileMetadata, string, error) {
    // Check for YAML front matter pattern: ---\n...\n---
    pattern := regexp.MustCompile(`(?s)^---\n(.+?)\n---\n(.*)$`)
    matches := pattern.FindStringSubmatch(content)

    if len(matches) != 3 {
        // No metadata found, return original content
        return nil, content, nil
    }

    yamlContent := matches[1]
    fileContent := matches[2]

    var metadata FileMetadata
    if err := yaml.Unmarshal([]byte(yamlContent), &metadata); err != nil {
        // Invalid YAML, return original content
        return nil, content, nil
    }

    return &metadata, fileContent, nil
}

// ExtractCodeBlocksWithMetadata parses markdown and extracts code blocks with metadata.
func ExtractCodeBlocksWithMetadata(markdown string) []CodeBlockWithMetadata {
    blocks := ExtractCodeBlocks(markdown)
    result := make([]CodeBlockWithMetadata, 0, len(blocks))

    for _, block := range blocks {
        metadata, cleanContent, _ := ParseMetadataFromContent(block.Content)
        result = append(result, CodeBlockWithMetadata{
            CodeBlock: CodeBlock{
                Language: block.Language,
                Filename: block.Filename,
                Content:  cleanContent,
            },
            Metadata: metadata,
        })
    }

    return result
}
```

### 3.3 Enhanced File Save Flow

**Sequence Diagram**:

```
User                  WebSocket           ChatService          ClaudeService        FileRepo         AppMapRepo
  |                       |                    |                     |                  |                 |
  |-- Send Message ------>|                    |                     |                  |                 |
  |                       |-- Process -------->|                     |                  |                 |
  |                       |                    |-- Stream Request -->|                  |                 |
  |                       |                    |<-- Stream Chunks ---|                  |                 |
  |<-- Chunk (real-time) -|<-------------------|                     |                  |                 |
  |                       |                    |                     |                  |                 |
  |                       |                    |-- Parse Markdown ---|                  |                 |
  |                       |                    |-- Extract Blocks ---|                  |                 |
  |                       |                    |   with Metadata     |                  |                 |
  |                       |                    |                     |                  |                 |
  |                       |                    |-- Save File --------|----------------->|                 |
  |                       |                    |                     |                  |                 |
  |                       |                    |-- Save Metadata ----|------------------|---------------->|
  |                       |                    |   (group, desc)     |                  |                 |
  |                       |                    |                     |                  |                 |
  |                       |                    |-- Create Relations -|------------------|---------------->|
  |                       |                    |   (if related_files)|                  |                 |
  |                       |                    |                     |                  |                 |
  |<-- File Created Event-|<-------------------|                     |                  |                 |
```

### 3.4 Modified Chat Service

**File**: `/workspace/backend/internal/service/chat.go` (additions)

```go
// processFileWithMetadata saves a file and its App Map metadata.
func (s *ChatService) processFileWithMetadata(
    ctx context.Context,
    projectID uuid.UUID,
    block markdown.CodeBlockWithMetadata,
) error {
    // Save the file content (existing logic)
    file, err := s.fileRepo.SaveFile(ctx, projectID, block.Filename, block.Language, block.Content)
    if err != nil {
        return fmt.Errorf("failed to save file: %w", err)
    }

    // If metadata was extracted, save it
    if block.Metadata != nil {
        // Find or create functional group
        groupID, err := s.appMapRepo.FindOrCreateGroup(ctx, projectID, block.Metadata.FunctionalGroup)
        if err != nil {
            s.logger.Warn().Err(err).Msg("failed to assign functional group")
            // Non-fatal: continue without group assignment
        }

        // Update file metadata
        metadataUpdate := model.FileMetadataUpdate{
            ShortDescription: block.Metadata.ShortDescription,
            LongDescription:  block.Metadata.LongDescription,
            Purpose:          block.Metadata.Purpose,
            FunctionalGroup:  block.Metadata.FunctionalGroup,
            AIConfidence:     0.85, // Default confidence for Claude-generated metadata
        }

        if _, err := s.appMapRepo.UpdateFileMetadata(ctx, file.ID, metadataUpdate); err != nil {
            s.logger.Warn().Err(err).Msg("failed to update file metadata")
        }

        // Create relationships for related files
        for _, relatedPath := range block.Metadata.RelatedFiles {
            if err := s.appMapRepo.CreateRelationship(ctx, projectID, file.ID, relatedPath, "references"); err != nil {
                s.logger.Warn().Err(err).Str("related", relatedPath).Msg("failed to create relationship")
            }
        }
    }

    return nil
}
```

---

## 4. Integration Points

### 4.1 Existing System Integration

| Component | Current State | Required Changes |
|-----------|---------------|------------------|
| `backend/internal/repository/file.go` | Saves files without metadata | Add metadata columns to save/update |
| `backend/internal/handler/file.go` | Returns basic file info | Extend responses with metadata |
| `backend/internal/service/chat.go` | Parses markdown, saves files | Add metadata extraction and saving |
| `backend/internal/pkg/markdown/codeblock.go` | Extracts language:filename | Extend to parse YAML front matter |
| `frontend/src/hooks/useFiles.ts` | Builds file tree from paths | Add App Map view building |
| `frontend/src/components/projects/FileTree.tsx` | Path-based tree display | Add purpose-based grouping view |

### 4.2 New Components Required

**Backend**:
- `/workspace/backend/internal/model/appmap.go` - App Map models
- `/workspace/backend/internal/repository/appmap.go` - App Map data access
- `/workspace/backend/internal/handler/appmap.go` - App Map API handlers
- `/workspace/backend/internal/pkg/markdown/metadata.go` - Metadata parser

**Frontend**:
- `/workspace/frontend/src/components/projects/AppMap.tsx` - Main App Map component
- `/workspace/frontend/src/components/projects/FunctionalGroup.tsx` - Group display
- `/workspace/frontend/src/components/projects/ViewLevelToggle.tsx` - Level 1/2/3 switcher
- `/workspace/frontend/src/hooks/useAppMap.ts` - App Map data hook
- `/workspace/frontend/src/types/appmap.ts` - TypeScript types

### 4.3 WebSocket Events

Add new WebSocket event types for real-time App Map updates:

```typescript
// New server message types
interface AppMapUpdateMessage {
  type: 'appmap_update';
  projectId: string;
  update: {
    type: 'file_added' | 'file_updated' | 'group_created' | 'relationship_added';
    data: FileWithMetadata | FunctionalGroup | FileRelationship;
  };
}
```

---

## 5. Migration Strategy

### 5.1 Backward Compatibility

The design ensures backward compatibility:

1. **Existing files**: Files without metadata continue to work; they appear in "Ungrouped" section
2. **Existing API**: `/api/projects/:id/files` endpoint unchanged
3. **File tree**: Remains available as Level 3 view

### 5.2 Migration Steps

1. **Phase 1 - Database Schema** (Day 1)
   - Apply migration `003_app_map_schema.sql`
   - Add new columns to files table
   - Create functional_groups and file_relationships tables

2. **Phase 2 - Backend APIs** (Days 2-3)
   - Implement new repository methods
   - Add App Map handler endpoints
   - Extend markdown parser for metadata

3. **Phase 3 - Claude Integration** (Days 3-4)
   - Update system prompt for metadata generation
   - Modify chat service to process metadata
   - Test with sample conversations

4. **Phase 4 - Frontend** (Days 5-7)
   - Implement App Map component
   - Add view level toggle
   - Integrate with existing project page

---

## 6. Architecture Decision Records

### ADR-001: File Metadata Storage Location

**Context**: Need to store App Map metadata (descriptions, purpose, group assignment) for files.

**Decision**: Add columns to existing `files` table rather than a separate metadata table.

**Rationale**:
- Metadata is always accessed with file data
- Simpler queries (no JOINs for basic operations)
- Nullable columns handle files without metadata

**Consequences**:
- Files table grows wider
- All file queries return metadata (slightly larger payloads)
- Migration is straightforward

### ADR-002: Functional Groups per Project

**Context**: Functional groups could be global templates or project-specific.

**Decision**: Functional groups are stored per-project with optional seeding from templates.

**Rationale**:
- Different projects may have different grouping needs
- Users can rename/customize groups
- Allows project-specific groups (e.g., "Checkout Flow" for e-commerce)

**Consequences**:
- Default groups must be copied to each new project
- More storage overhead for group definitions
- Greater flexibility for users

### ADR-003: YAML Front Matter for Metadata

**Context**: Claude needs to generate metadata alongside code.

**Decision**: Use YAML front matter format within code blocks.

**Rationale**:
- Familiar pattern (common in static site generators)
- Easy to parse and extract
- Clean separation from actual code
- Claude handles YAML generation well

**Consequences**:
- Requires prompt engineering to ensure consistent format
- Parser must handle malformed YAML gracefully
- Slight increase in response token usage

### ADR-004: Three-Level View Progression

**Context**: Need to progressively disclose technical complexity.

**Decision**: Implement fixed three-level system (Functional Only, Functional + Files, Full Technical).

**Rationale**:
- Matches product vision's learning progression
- Clear mental model for users
- Simpler to implement than arbitrary granularity

**Consequences**:
- May need intermediate levels later
- All three views need distinct UI implementations
- Level transitions need careful UX design

---

## 7. Testing Strategy

### 7.1 Unit Tests

- Metadata parser: Valid YAML, malformed YAML, missing metadata
- Repository: CRUD operations for groups, metadata, relationships
- Handler: API response formats at each view level

### 7.2 Integration Tests

- End-to-end file creation with metadata extraction
- App Map retrieval with populated data
- Relationship creation from Claude output

### 7.3 Acceptance Criteria

1. Claude-generated files include metadata in response
2. Metadata is correctly parsed and stored
3. App Map API returns correct structure at each level
4. Files without metadata appear in "Ungrouped" section
5. View level toggle changes visible information

---

## 8. Open Questions

1. **AI Fallback**: What happens when Claude fails to generate metadata? Options:
   - Use AI to analyze code and generate metadata post-hoc
   - Leave metadata empty, prompt user to categorize
   - Use heuristics based on file path and language

2. **User Override**: Should users be able to:
   - Create custom functional groups?
   - Manually assign files to groups?
   - Edit AI-generated descriptions?

3. **Relationship Accuracy**: How to validate AI-detected relationships?
   - Static analysis for imports/dependencies
   - Trust AI confidence scores
   - Manual user verification

---

## 9. Implementation Recommendations

### Priority Order

1. **Critical Path** (must have for MVP):
   - Database schema with file metadata columns
   - Backend APIs for App Map retrieval at 3 levels
   - Claude prompt enhancement for metadata generation
   - Basic frontend App Map component with level toggle

2. **Important** (high value, implement soon after MVP):
   - Functional groups management
   - File relationship storage and display
   - Real-time App Map updates via WebSocket

3. **Nice to Have** (can defer):
   - Graph view visualization
   - User-created custom groups
   - Relationship editing

### Estimated Effort

| Component | Estimated Time |
|-----------|---------------|
| Database migration | 2 hours |
| Backend repository + models | 4 hours |
| Backend handlers | 3 hours |
| Markdown parser enhancement | 2 hours |
| Chat service integration | 3 hours |
| Frontend AppMap component | 6 hours |
| Frontend view toggle | 2 hours |
| Testing and refinement | 4 hours |
| **Total** | **~26 hours** |

---

## 10. Summary

The App Map architecture extends Go Chat's existing file system with:

1. **Purpose-based organization** through functional groups stored per-project
2. **Rich metadata** on files including descriptions, purpose, and confidence scores
3. **Progressive disclosure** via three view levels matching user skill progression
4. **AI integration** through YAML front matter in Claude's code generation
5. **Relationship tracking** for advanced graph visualization

The design prioritizes backward compatibility, leveraging existing infrastructure while adding the metadata and grouping capabilities needed for the App Map vision.

---

**Document History**

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2025-12-25 | 1.0 | Initial architecture design | Architect Agent |
