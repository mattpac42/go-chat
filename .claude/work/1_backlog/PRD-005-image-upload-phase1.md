# PRD-005: Image Upload with Claude Vision (Phase 1)

**Version**: 1.0
**Created**: 2025-12-27
**Author**: Product Manager (Tactical)
**Status**: Backlog
**Phase**: Phase 1 of File Upload Feature
**Dependencies**: Existing Claude API integration

---

## Executive Summary

Enable users to upload images (screenshots, photos, diagrams) during AI-guided discovery. Images are converted to markdown descriptions using Claude Vision API and stored as markdown files in the project's `sources/` folder, providing context seeding for AI-assisted project planning.

---

## 1. Problem Statement

### Current Pain Points

Users participating in AI-guided discovery often have existing visual materials that would inform their project:

1. **Screenshots of existing systems** they want to replicate or improve
2. **Photos of paper forms** being digitized
3. **Diagrams and wireframes** showing desired functionality
4. **Reference images** from competitor applications
5. **Quick captures** from snipping tools during conversation

Currently, users cannot share these visual materials with the AI. They must:
- Manually describe images in text (error-prone, time-consuming)
- Reference external URLs the AI cannot access
- Skip providing visual context entirely

This creates an **information gap** between what users can see and what the AI understands.

### Business Impact

- **Slower discovery sessions**: Users spend time describing visuals instead of focusing on requirements
- **Lost context**: Important visual details are missed or poorly communicated
- **Reduced output quality**: AI generates requirements without understanding user's visual reference points
- **Poor user experience**: Users expect modern AI tools to "see" what they share

---

## 2. Proposed Solution

### Solution Overview

Allow users to upload images directly in the chat interface. Images are:

1. **Uploaded** via drag-drop, file picker, or clipboard paste
2. **Converted** to descriptive markdown using Claude Vision API
3. **Stored** as `.md` files in the project's `sources/` folder
4. **Displayed** in the file tree under "Source Materials" functional group
5. **Available** as context for subsequent AI interactions

### Why Phase 1 Focuses on Images

| Factor | Images | PDF/DOCX |
|--------|--------|----------|
| Conversion complexity | Low (single Claude API call) | Medium-High (multi-step extraction) |
| Existing integration | Claude Vision ready | Requires new libraries |
| Use case frequency | High (screenshots common) | Medium |
| Error rate | Low | Higher (complex layouts) |
| Implementation time | 2-3 days | 3-5 days |

Images provide the **highest value with lowest risk** for the initial implementation.

---

## 3. User Stories

### Primary User Stories

#### US-001: Upload Image via File Picker
**As a** user in discovery conversation
**I want to** upload an image file from my computer
**So that** the AI can understand my visual reference and incorporate it into project planning

**Acceptance Criteria:**
- [ ] Upload button/icon visible in chat input area
- [ ] Clicking opens native file picker
- [ ] File picker filters to supported image types (PNG, JPG, JPEG, GIF, WebP)
- [ ] Selected file uploads with progress indicator
- [ ] Conversion completes within 10 seconds for typical images
- [ ] Resulting markdown file appears in file tree under "Source Materials"
- [ ] Success message confirms upload and conversion

#### US-002: Upload Image via Drag and Drop
**As a** user with an image file on my desktop
**I want to** drag and drop the image directly into the chat area
**So that** I can quickly share visual context without navigating menus

**Acceptance Criteria:**
- [ ] Chat area shows visual drop zone when file is dragged over
- [ ] Drop zone indicates supported file types
- [ ] Dropping image initiates upload and conversion
- [ ] Same success flow as file picker upload
- [ ] Invalid file types show clear error message

#### US-003: Paste Screenshot from Clipboard
**As a** user who just captured a screenshot with snipping tool
**I want to** paste the image directly into the chat with Ctrl+V/Cmd+V
**So that** I can share quick captures without saving files first

**Acceptance Criteria:**
- [ ] Ctrl+V / Cmd+V in chat area detects clipboard image
- [ ] Pasted image shows preview before upload
- [ ] User can confirm or cancel upload
- [ ] Confirmed paste triggers same conversion flow
- [ ] Generated filename includes timestamp (e.g., `screenshot-20251227-143022.md`)

#### US-004: View Converted Markdown
**As a** user who uploaded an image
**I want to** see the markdown description the AI generated
**So that** I can verify the AI understood my image correctly

**Acceptance Criteria:**
- [ ] Converted file visible in file tree under "Source Materials"
- [ ] Clicking file shows markdown content in file viewer
- [ ] Markdown includes: image description, extracted text (if any), and structural notes
- [ ] Original filename preserved in markdown front matter

#### US-005: Handle Upload Errors
**As a** user uploading an image
**I want to** see clear error messages when something goes wrong
**So that** I know what to fix or try differently

**Acceptance Criteria:**
- [ ] File too large (>10MB): "File exceeds 10MB limit. Please resize or compress."
- [ ] Unsupported type: "File type not supported. Please use PNG, JPG, JPEG, GIF, or WebP."
- [ ] Conversion failure: "Could not process image. Please try a different image."
- [ ] Network error: "Upload failed. Please check your connection and try again."
- [ ] All errors dismissible and non-blocking

### Secondary User Stories

#### US-006: Reference Uploaded Material in Chat
**As a** user who uploaded images
**I want** the AI to automatically consider my uploads in subsequent responses
**So that** my visual context informs project requirements

**Acceptance Criteria:**
- [ ] Uploaded markdown files included in system context for AI
- [ ] AI responses can reference uploaded materials by name
- [ ] Multiple uploads handled (up to 10 files per project)

#### US-007: Delete Uploaded Material
**As a** user who uploaded an incorrect image
**I want to** delete the converted markdown file
**So that** irrelevant content does not pollute my project

**Acceptance Criteria:**
- [ ] Delete action available in file tree context menu
- [ ] Confirmation dialog before deletion
- [ ] File removed from `sources/` folder and database
- [ ] File tree updates immediately

---

## 4. Acceptance Criteria (Feature Level)

### Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-01 | Support PNG, JPG/JPEG, GIF, WebP image formats | Must |
| FR-02 | Convert images to markdown using Claude Vision API | Must |
| FR-03 | Store converted markdown in `sources/` folder | Must |
| FR-04 | Support file picker upload | Must |
| FR-05 | Support drag-and-drop upload | Must |
| FR-06 | Support clipboard paste upload | Must |
| FR-07 | Display files in file tree under "Source Materials" | Must |
| FR-08 | Enforce 10MB file size limit | Must |
| FR-09 | Enforce 10 files per project limit | Must |
| FR-10 | Show upload progress indicator | Should |
| FR-11 | Allow deletion of uploaded files | Should |
| FR-12 | Include original filename in markdown metadata | Should |

### Non-Functional Requirements

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-01 | Image conversion time | < 10 seconds for images under 2MB |
| NFR-02 | Upload progress feedback | Visual feedback within 500ms of action |
| NFR-03 | Error message display | < 1 second from error occurrence |
| NFR-04 | Memory usage during upload | < 50MB additional memory per upload |
| NFR-05 | Concurrent uploads | Support at least 3 simultaneous uploads |

---

## 5. Success Metrics

### Primary Metrics

| Metric | Target | Measurement Method |
|--------|--------|-------------------|
| Feature adoption | 40% of discovery sessions include at least 1 image upload | Analytics: uploads per session |
| Conversion success rate | > 95% of uploads convert successfully | Backend logs: conversion status |
| User completion rate | > 90% of started uploads complete | Frontend analytics: upload funnel |

### Secondary Metrics

| Metric | Target | Measurement Method |
|--------|--------|-------------------|
| Average conversion time | < 5 seconds | Backend logs: processing duration |
| Clipboard paste usage | > 30% of uploads via paste | Frontend analytics: upload method |
| Error rate by type | < 5% size errors, < 2% format errors | Backend logs: error categories |
| Files per project | Average 2-3 files | Database: files per project |

### Qualitative Indicators

- Users reference uploaded materials in chat naturally
- AI responses demonstrate understanding of uploaded content
- Reduced "describe what you're looking at" prompts in conversations

---

## 6. Out of Scope (Phase 2+)

The following items are explicitly excluded from Phase 1:

### Phase 2: Document Upload
- **PDF file upload and text extraction**
- **DOCX file upload and conversion**
- Complex document layout preservation
- Multi-page document handling

### Phase 3: Advanced Features
- Image editing/cropping before upload
- Batch upload (multiple files at once)
- Re-conversion of existing uploads
- Version history for uploaded files

### Not Planned
- Video file upload
- Audio file upload
- Cloud storage integration (Google Drive, Dropbox)
- URL-based image import
- Image search/embedding capabilities

---

## 7. Technical Notes

### Existing Integration Points

**Claude Vision API** (already integrated in backend):
- Current `ClaudeService` in `/workspace/backend/internal/service/claude.go`
- Uses `api.anthropic.com/v1/messages` endpoint
- Vision capability requires adding `image` content block to message format

**Extension Required**:
```go
// Current: text-only message
type ClaudeMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// Required: support image content blocks
type ClaudeMessage struct {
    Role    string      `json:"role"`
    Content interface{} `json:"content"` // string or []ContentBlock
}

type ContentBlock struct {
    Type   string       `json:"type"` // "text" or "image"
    Text   string       `json:"text,omitempty"`
    Source *ImageSource `json:"source,omitempty"`
}

type ImageSource struct {
    Type      string `json:"type"`       // "base64"
    MediaType string `json:"media_type"` // "image/png", etc.
    Data      string `json:"data"`       // base64 encoded
}
```

### Database Schema Extension

Reference architecture decision from ADR-003:
```sql
-- Add source tracking to files table
ALTER TABLE files ADD COLUMN IF NOT EXISTS source_type VARCHAR(20) DEFAULT 'generated';
-- Values: 'generated' (AI), 'uploaded' (user), 'converted' (upload + conversion)

-- Track original upload metadata
CREATE TABLE IF NOT EXISTS file_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    original_filename VARCHAR(255) NOT NULL,
    original_mime_type VARCHAR(100) NOT NULL,
    original_size_bytes BIGINT NOT NULL,
    conversion_status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Conversion Output Format

Converted markdown files should follow this structure:
```markdown
---
original_filename: screenshot-order-form.png
original_type: image/png
converted_at: 2025-12-27T14:30:22Z
short_description: Screenshot of paper order form with customer and product fields
long_description: User-uploaded screenshot showing a paper-based order form currently used by the business
functional_group: Source Materials
---

# Image: screenshot-order-form.png

## Description
This image shows a paper order form with the following elements...

## Text Extracted
- "Customer Name" (field label)
- "Product Selection" (section header)
- ...

## Structural Notes
- Two-column layout
- Header section at top
- Signature area at bottom
```

### API Endpoint Design

```
POST /api/projects/:id/upload
Content-Type: multipart/form-data

Request:
- file: binary (required, max 10MB)
- source: string (optional, "filepicker" | "dragdrop" | "clipboard")

Response (200):
{
  "file": {
    "id": "uuid",
    "path": "sources/screenshot-order-form.md",
    "filename": "screenshot-order-form.md",
    "language": "markdown"
  },
  "source": {
    "originalFilename": "screenshot-order-form.png",
    "originalMimeType": "image/png",
    "originalSizeBytes": 245678,
    "conversionStatus": "completed"
  }
}

Response (413): File too large
Response (415): Unsupported media type
Response (422): Conversion failed
```

### Frontend Components Needed

1. **UploadButton** - File picker trigger in chat input
2. **DropZone** - Drag-and-drop overlay for chat area
3. **ClipboardHandler** - Intercept Ctrl+V for image paste
4. **UploadProgress** - Progress indicator during upload/conversion
5. **PastePreview** - Confirmation dialog for clipboard paste

---

## 8. UX Flow

### Happy Path: Clipboard Paste

```
User captures screenshot with snipping tool
         |
         v
User presses Ctrl+V in chat input
         |
         v
+---------------------------+
|  Paste Preview Dialog     |
|  +---------------------+  |
|  | [Image Preview]     |  |
|  +---------------------+  |
|  This image will be     |
|  converted to markdown  |
|  and added to your      |
|  project sources.       |
|                         |
|  [Cancel]  [Upload]     |
+---------------------------+
         |
         v (User clicks Upload)
+---------------------------+
|  Converting...           |
|  [=====>    ] 60%        |
+---------------------------+
         |
         v
+---------------------------+
|  Upload Complete          |
|  screenshot-20251227.md  |
|  added to Source Materials|
|                          |
|  [View File] [Dismiss]   |
+---------------------------+
         |
         v
File appears in tree under "Source Materials"
```

### Error Path: File Too Large

```
User drags 15MB image to chat
         |
         v
+---------------------------+
|  Upload Error            |
|                          |
|  File exceeds 10MB limit |
|  (15.2MB uploaded)       |
|                          |
|  Please resize or        |
|  compress your image     |
|  and try again.          |
|                          |
|  [Dismiss]               |
+---------------------------+
```

---

## 9. Implementation Recommendations

### Suggested Sprint Breakdown

**Sprint 1: Backend Foundation (2-3 days)**
- Database migration for `file_sources` table
- Extend `ClaudeService` for Vision API calls
- Create `/api/projects/:id/upload` endpoint
- Image-to-markdown conversion logic

**Sprint 2: Frontend Integration (2-3 days)**
- Upload button in chat input
- Drag-and-drop zone implementation
- Clipboard paste handler
- Progress and error UI components

**Sprint 3: Polish and Testing (1-2 days)**
- File tree integration ("Source Materials" group)
- End-to-end testing
- Error handling edge cases
- Performance optimization

### Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Claude Vision API rate limits | Implement retry with exponential backoff |
| Large image conversion timeout | Set 30-second timeout, show user-friendly message |
| Browser clipboard API inconsistency | Feature-detect and gracefully degrade |
| Memory pressure from large images | Stream uploads, don't hold full image in memory |

---

## 10. Open Questions

1. **Filename generation for clipboard pastes**: Use timestamp-based naming or prompt user for a name?
   - **Recommendation**: Timestamp-based with optional rename later

2. **Conversion quality feedback**: Should users be able to edit the converted markdown?
   - **Recommendation**: Yes, enable editing in file viewer (existing functionality)

3. **Multiple rapid uploads**: Queue or parallel processing?
   - **Recommendation**: Queue with max 3 concurrent conversions

4. **AI context inclusion**: Include all source files or only recent ones?
   - **Recommendation**: Include all (up to 10 file limit enforced)

---

## Document History

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2025-12-27 | 1.0 | Initial PRD creation | Product Manager (Tactical) |

---

## Related Documents

- `/workspace/.claude/work/history/20251227-203101-architect-file-upload-conversion.md` - Architecture decisions
- `/workspace/.claude/work/history/20251227-203151-product-strategic-file-upload-strategy.md` - Strategic analysis
- `/workspace/.claude/work/1_backlog/PRD-004-discovery-enhancements.md` - Parent feature context
- `/workspace/backend/internal/service/claude.go` - Existing Claude integration
