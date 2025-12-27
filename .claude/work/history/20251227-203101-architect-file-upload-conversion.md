# Architect Session: File Upload with Document Conversion

**Date**: 2025-12-27T20:31:01Z
**Agent**: architect
**Task**: Design architecture for file upload feature with PDF/DOCX/Image to Markdown conversion

## Work Completed

Analyzed current codebase architecture and designed file upload system with document conversion capabilities.

## Current Architecture Analysis

### Backend (Go)
- **Framework**: Gin web framework
- **Database**: PostgreSQL with sqlx
- **Files Table Schema**: `id`, `project_id`, `path`, `filename`, `language`, `content` (TEXT), `created_at`
- **File Metadata Table**: `file_id`, `short_description`, `long_description`, `functional_group`
- **Repository Pattern**: Clean separation with `FileRepository` interface

### Frontend (Next.js/React)
- **File Explorer**: Supports tree, reveal, and grouped views
- **Hook**: `useFiles` for file management
- **API Client**: `api.getProjectFiles()`, `api.getFile()`

## Architecture Recommendations

---

## ADR-001: Document Conversion Strategy

### Context
Users need to upload PDF, DOCX, and image files that should be converted to markdown and stored as regular text files in the database.

### Decision
**Hybrid Approach with Go Libraries + Claude API Fallback**

| Format | Primary Solution | Fallback | Rationale |
|--------|-----------------|----------|-----------|
| PDF | `ledongthuc/pdf` | Claude Vision API | Go-native is fast for text PDFs; Claude handles scanned/complex |
| DOCX | `unidoc/unioffice` or custom parser | - | Go-native sufficient for standard documents |
| Images | Claude Vision API | Tesseract (optional) | Claude provides excellent OCR + understanding |

### Consequences
- **Positive**: Low latency for simple documents, graceful degradation for complex ones
- **Negative**: Additional Claude API costs for images and complex PDFs
- **Trade-off**: Balance between cost and quality

---

## ADR-002: Conversion Pipeline Architecture

### Context
Need to determine where and how document conversion should occur.

### Decision
**Synchronous Backend Processing with Timeout and Size Limits**

```
[Upload Request]
      |
      v
[Size Check] --> [REJECT if > 10MB]
      |
      v
[MIME Detection]
      |
      v
[Format Router]
  /    |    \
PDF  DOCX  Image
  \    |    /
      v
[Converter Service]
      |
      v
[Markdown Output]
      |
      v
[Save to files table]
```

### Key Design Points

1. **Synchronous for Small Files (<2MB)**
   - Convert inline during upload request
   - Return immediately with file ID

2. **Async with Polling for Large Files (2-10MB)**
   - Return upload ID immediately
   - Poll `/api/uploads/{id}/status`
   - WebSocket notification when complete

3. **Reject >10MB**
   - Return 413 with guidance to split document

### Consequences
- **Positive**: Simple implementation for most use cases
- **Negative**: Sync path may block for large PDFs
- **Mitigation**: Async path for files over threshold

---

## ADR-003: Database Schema Changes

### Context
Need to track upload source, conversion status, and preserve original file reference.

### Decision
**Add `file_sources` table and `source_type` column**

```sql
-- Migration: 008_file_uploads.sql

-- Track original uploaded files and conversion status
CREATE TABLE IF NOT EXISTS file_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    original_filename VARCHAR(255) NOT NULL,
    original_mime_type VARCHAR(100) NOT NULL,
    original_size_bytes BIGINT NOT NULL,
    conversion_status VARCHAR(20) DEFAULT 'pending', -- pending, processing, completed, failed
    conversion_error TEXT,
    conversion_started_at TIMESTAMP WITH TIME ZONE,
    conversion_completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_file_sources_file_id ON file_sources(file_id);
CREATE INDEX idx_file_sources_status ON file_sources(conversion_status);

-- Add source type to files table
ALTER TABLE files ADD COLUMN IF NOT EXISTS source_type VARCHAR(20) DEFAULT 'generated';
-- Values: 'generated' (from AI), 'uploaded' (user upload), 'converted' (upload + conversion)

COMMENT ON TABLE file_sources IS 'Tracks original uploaded files before markdown conversion';
COMMENT ON COLUMN files.source_type IS 'How the file was created: generated, uploaded, or converted';
```

### Consequences
- **Positive**: Full audit trail, can reconstruct original metadata
- **Negative**: Additional table to maintain
- **Trade-off**: Not storing original binary (per requirements) saves storage but loses ability to re-convert

---

## ADR-004: Conversion Library Selection

### Context
Need Go libraries for PDF and DOCX parsing.

### Decision

#### PDF Conversion
| Option | Pros | Cons | Recommendation |
|--------|------|------|----------------|
| `ledongthuc/pdf` | Pure Go, simple API | Text-only, no layout | **Primary** |
| `pdfcpu/pdfcpu` | Full PDF manipulation | Complex, overkill | Optional |
| External (pdftotext) | Best extraction | Container dependency | Fallback |

**Recommendation**: Use `ledongthuc/pdf` for text extraction, fall back to Claude Vision for scanned/image PDFs.

#### DOCX Conversion
| Option | Pros | Cons | Recommendation |
|--------|------|------|----------------|
| `unidoc/unioffice` | Full OOXML support | Commercial license | If budget allows |
| Custom `archive/zip` + XML | Free, adequate | More code | **Primary** |
| Pandoc (external) | Best quality | Container dependency | Optional |

**Recommendation**: Custom DOCX parser using Go's `archive/zip` to extract `word/document.xml` and convert to markdown. DOCX is well-structured XML.

#### Image OCR
| Option | Pros | Cons | Recommendation |
|--------|------|------|----------------|
| Claude Vision API | Excellent understanding | Cost per image | **Primary** |
| Tesseract (Go binding) | Free, local | Lower quality | Fallback |
| Google Vision API | Good quality | Additional API | Optional |

**Recommendation**: Claude Vision API for images - already integrated, excellent quality, understands context.

---

## ADR-005: Error Handling Strategy

### Context
Document conversion can fail for many reasons.

### Decision
**Graceful Degradation with User Feedback**

```
Conversion Flow:
1. Try primary converter
2. On failure, try fallback if available
3. On fallback failure, store with error metadata
4. Return partial success with error details

Error Response Format:
{
  "success": true,  // Upload succeeded
  "file": {...},    // Basic file info
  "conversion": {
    "status": "partial",  // "complete", "partial", "failed"
    "warnings": ["Table on page 3 may have formatting issues"],
    "errors": [],
    "originalPages": 12,
    "convertedPages": 12
  }
}
```

### Consequences
- **Positive**: Users always get feedback, partial content is better than none
- **Negative**: May deliver imperfect conversion
- **Mitigation**: Clear UI indicators for conversion quality

---

## ADR-006: Original File Preservation Policy

### Context
Requirements state "not stored as original binary" - need to clarify approach.

### Decision
**No Original Binary Storage (Memory-Only Processing)**

1. Upload binary to memory
2. Convert to markdown
3. Store only markdown
4. Discard binary after conversion

**Rationale**:
- Reduces storage costs significantly
- Simplifies backup/restore
- Markdown is version-control friendly
- User can re-upload if needed

**Alternative Considered**: Temporary blob storage with 24-hour expiry for re-conversion. Rejected due to complexity.

---

## Component Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (Next.js)                        │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │ FileUpload  │  │ FileExplorer│  │   ConversionStatus     │ │
│  │  Component  │  │ (existing)  │  │      Component         │ │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘ │
│         │                                      │                │
│         └──────────────┬───────────────────────┘                │
│                        │                                        │
│                  api.uploadFile()                               │
└────────────────────────┼────────────────────────────────────────┘
                         │
┌────────────────────────┼────────────────────────────────────────┐
│                        │           Backend (Go)                  │
├────────────────────────┼────────────────────────────────────────┤
│                        v                                        │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    UploadHandler                             ││
│  │  POST /api/projects/:id/upload                              ││
│  └─────────────────────────────────────────────────────────────┘│
│                        │                                        │
│                        v                                        │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                 ConversionService                           ││
│  │  ┌───────────────────────────────────────────────────────┐ ││
│  │  │ Converters:                                           │ ││
│  │  │  - PDFConverter (ledongthuc/pdf + Claude fallback)    │ ││
│  │  │  - DOCXConverter (custom XML parser)                  │ ││
│  │  │  - ImageConverter (Claude Vision API)                 │ ││
│  │  └───────────────────────────────────────────────────────┘ ││
│  └─────────────────────────────────────────────────────────────┘│
│                        │                                        │
│                        v                                        │
│  ┌──────────────┐  ┌──────────────────────┐                    │
│  │ FileRepository│  │ FileSourceRepository │                    │
│  │  (existing)   │  │     (new)            │                    │
│  └──────────────┘  └──────────────────────┘                    │
│         │                    │                                  │
└─────────┴────────────────────┴──────────────────────────────────┘
                         │
                         v
                   ┌──────────┐
                   │PostgreSQL│
                   │  files   │
                   │file_srcs │
                   └──────────┘
```

---

## API Design

### Upload Endpoint
```
POST /api/projects/:id/upload
Content-Type: multipart/form-data

Body:
- file: binary (required)
- path: string (optional, defaults to filename)

Response (200):
{
  "file": {
    "id": "uuid",
    "path": "uploads/document.md",
    "filename": "document.md",
    "language": "markdown",
    "createdAt": "2025-12-27T..."
  },
  "source": {
    "originalFilename": "quarterly-report.pdf",
    "originalMimeType": "application/pdf",
    "originalSizeBytes": 1245678,
    "conversionStatus": "completed",
    "conversionDuration": "1.2s"
  }
}

Response (413):
{
  "error": "file_too_large",
  "message": "Maximum file size is 10MB",
  "limit": 10485760,
  "received": 15728640
}

Response (415):
{
  "error": "unsupported_type",
  "message": "Unsupported file type: application/zip",
  "supportedTypes": ["application/pdf", "image/png", "image/jpeg", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"]
}
```

### Conversion Status (for async)
```
GET /api/uploads/:id/status

Response:
{
  "id": "uuid",
  "status": "processing", // pending, processing, completed, failed
  "progress": 45,         // percentage
  "fileId": null,         // populated when complete
  "error": null
}
```

---

## Recommendations Summary

| Question | Recommendation |
|----------|---------------|
| PDF library | `ledongthuc/pdf` + Claude Vision fallback |
| DOCX library | Custom parser with `archive/zip` |
| Image OCR | Claude Vision API (already integrated) |
| Conversion location | Backend, sync for <2MB, async for >2MB |
| Size limits | 10MB max, chunking not recommended |
| Failure handling | Graceful degradation with status feedback |
| Original preservation | No binary storage, markdown only |
| Schema changes | Add `file_sources` table + `source_type` column |

---

## Files Modified
None (architecture documentation only)

## Recommendations for Implementation

1. **Phase 1**: Image upload with Claude Vision (simplest, reuses existing integration)
2. **Phase 2**: PDF upload with text extraction
3. **Phase 3**: DOCX support
4. **Phase 4**: Async processing for large files

Estimated implementation effort: 3-4 days for Phase 1-2, additional 2 days for Phase 3-4.

## Next Steps for Developer Agent

1. Create migration `008_file_uploads.sql`
2. Implement `ConversionService` interface
3. Start with `ImageConverter` using existing Claude client
4. Add `UploadHandler` to file handler
5. Update frontend with upload component
