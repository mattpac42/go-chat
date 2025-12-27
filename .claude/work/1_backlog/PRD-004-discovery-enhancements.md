# PRD-004: Discovery Enhancements

**Version**: 1.0
**Created**: 2025-12-26
**Author**: Product Manager
**Status**: Backlog
**Phase**: Future
**Dependencies**: PRD-003 (Guided Discovery)

---

## Feature 1: Cost Estimation

### Problem
Users start building without understanding the potential development cost. This leads to:
- Sticker shock when they see final estimates
- Abandoned projects mid-development
- Unrealistic expectations about MVP scope

### Proposed Solution
Show an estimated cost range during the Summary stage based on:
- Number of features (MVP + future)
- User count and permission complexity
- Integration requirements
- Hosting/infrastructure needs

### UX Concept
```
+----------------------------------------------+
|  ESTIMATED COST                              |
|  $500 - $1,500 for MVP                       |
|  Based on: 3 features, 3 users, basic auth   |
|                                              |
|  [?] How is this calculated?                 |
+----------------------------------------------+
```

### Implementation Notes
- Backend: Create pricing model based on feature complexity
- Store cost factors per feature type
- Display range (not exact) to set expectations
- Allow user to see breakdown

### Effort Estimate
- Backend: 2-3 days (pricing model, API)
- Frontend: 1-2 days (UI components)
- Calibration: Ongoing

---

## Feature 2: Reference File Upload (Markdown Extraction)

### Problem
Users have existing materials that could inform the project:
- Screenshots of competitor apps
- PDF wireframes or mockups
- Training documents for workflows
- Photos of paper forms being replaced

Currently no way to share these with the AI.

### Proposed Solution
**Simplified approach: Convert to markdown on upload, no file storage.**

Allow file upload during discovery that:
1. Accepts file upload (image, PDF, DOCX)
2. Extracts content immediately using Claude Vision / text extraction
3. Converts to markdown and stores in database
4. Discards original file (no S3/storage needed)
5. Shows extracted content to user for review/editing

### Why Markdown-Only?
- Files are reference material, not deliverables
- Agents need the *content*, not the *file*
- No storage infrastructure = simpler + cheaper
- Fits naturally into chat context window
- User can edit extraction if AI missed something

### Data Model
```typescript
interface DiscoveryAttachment {
  id: string;
  discoveryId: string;
  filename: string;           // "order-form.pdf" (for display)
  fileType: string;           // "application/pdf"
  extractedMarkdown: string;  // The usable content
  extractionMethod: string;   // "claude-vision" | "pdftotext" | "manual"
  uploadedAt: string;
  editedAt?: string;          // If user manually edited
  // NO original file storage
}
```

### UX Concept

**Upload Flow:**
```
+----------------------------------------------+
| [+] Add reference files                      |
|                                              |
| Drag files here or click to upload           |
| Supports: Images, PDFs, Word docs            |
+----------------------------------------------+
```

**After Processing (shows extracted content):**
```
+--------------------------------------------+
| order-form.pdf                          [x]|
+--------------------------------------------+
| ## Extracted Content                       |
|                                            |
| **Order Form Fields:**                     |
| - Customer Name (text input)               |
| - Cake Type (dropdown: chocolate, vanilla) |
| - Delivery Date (date picker)              |
| - Special Instructions (textarea)          |
|                                            |
| **Layout Notes:**                          |
| Two-column form with submit button at      |
| bottom right.                              |
|                                            |
| [Edit] if extraction needs fixing          |
+--------------------------------------------+
```

### Implementation Notes

**Backend:**
- POST `/api/projects/:id/discovery/attachments` - upload + extract
- GET `/api/projects/:id/discovery/attachments` - list attachments
- PUT `/api/projects/:id/discovery/attachments/:id` - edit markdown
- DELETE `/api/projects/:id/discovery/attachments/:id` - remove

**Extraction Pipeline:**
```go
func ExtractToMarkdown(file multipart.File, contentType string) (string, error) {
    switch {
    case strings.HasPrefix(contentType, "image/"):
        return extractWithClaudeVision(file)
    case contentType == "application/pdf":
        return extractPDFText(file)
    case isWordDoc(contentType):
        return extractWordText(file)
    default:
        return "", ErrUnsupportedFileType
    }
}
```

**Frontend:**
- Drag-and-drop upload zone component
- Processing spinner during extraction
- Markdown preview with edit button
- Delete attachment button

**Integration with Agents:**
- Attachments included in system prompt context
- Format: `## Reference: {filename}\n{extractedMarkdown}`

### Effort Estimate
- Backend: 1.5 days (upload endpoint, extraction, CRUD)
- Frontend: 1 day (upload UI, preview, edit)
- Testing: 0.5 day

**Total: ~3 days** (down from 5-6 days with file storage)

---

## Priority

| Feature | Business Value | Technical Effort | Priority |
|---------|---------------|------------------|----------|
| Cost Estimation | High | Medium (~3-4 days) | P1 |
| File Upload (Markdown) | High | Low (~3 days) | P1 |

Both features are now similar effort - could implement in parallel.

---

## Open Questions

1. **Cost Model**: How do we calibrate pricing? Per feature? Per complexity score?
2. **File Storage**: S3 vs local? Cost implications?
3. **File Size Limits**: What's reasonable? 10MB? 50MB?
4. **File Retention**: How long do we keep uploaded files?
5. **Privacy**: How do we handle sensitive documents?

---

**Document History**

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2025-12-26 | 1.0 | Initial feature proposals | Product Manager |
