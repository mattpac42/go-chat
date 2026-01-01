# Completeness Checker Specification

## Overview

The Completeness Checker is a validation layer that runs after file creation to detect missing dependencies, broken references, and incomplete builds. It automatically identifies issues and either fixes them or prompts the user.

## Goals

1. **Detect** missing files referenced by other files
2. **Validate** basic syntax/structure of generated code
3. **Auto-fix** simple issues without user intervention
4. **Surface** problems clearly to non-technical users
5. **Prevent** broken previews from being shown

---

## When It Runs

| Trigger | Action |
|---------|--------|
| After each `write_file` tool completes | Quick check for that file's dependencies |
| After message completes (all files written) | Full project validation |
| Before preview is rendered | Block preview if critical issues |
| On session timeout/error | Log incomplete state for recovery |

---

## What It Checks

### 1. HTML File Validation

```html
<!-- Check these references exist -->
<script src="app.js"></script>          → Does app.js exist?
<link href="styles.css" rel="stylesheet"> → Does styles.css exist?
<img src="logo.png">                     → Does logo.png exist?
```

**Extraction regex:**
```go
scriptPattern := `<script[^>]+src=["']([^"']+)["']`
linkPattern := `<link[^>]+href=["']([^"']+\.css)["']`
imgPattern := `<img[^>]+src=["']([^"']+)["']`
```

### 2. JavaScript/TypeScript Validation

```javascript
// Check these imports resolve
import { App } from './App.js'           → Does App.js exist?
import styles from './styles.css'        → Does styles.css exist?
require('./utils')                       → Does utils.js exist?
```

**Extraction regex:**
```go
importPattern := `import\s+.*from\s+['"]([^'"]+)['"]`
requirePattern := `require\s*\(\s*['"]([^'"]+)['"]\s*\)`
```

### 3. CSS Validation

```css
/* Check these references exist */
@import url('reset.css');               → Does reset.css exist?
background: url('../images/bg.png');    → Does bg.png exist?
```

### 4. Package.json Validation (Future)

```json
{
  "main": "index.js",        → Does index.js exist?
  "scripts": {
    "start": "node server.js" → Does server.js exist?
  }
}
```

### 5. Basic Syntax Checks

| File Type | Check |
|-----------|-------|
| HTML | Valid opening/closing tags, DOCTYPE present |
| JSON | Valid JSON parse |
| JS/TS | Balanced braces/brackets (basic) |
| CSS | Balanced braces |

---

## Issue Severity Levels

| Level | Description | Action |
|-------|-------------|--------|
| **Critical** | App won't run (missing main JS file) | Block preview, auto-fix |
| **Warning** | Degraded experience (missing image) | Show preview with notice |
| **Info** | Best practice (missing favicon) | Log only |

### Severity Rules

```go
func getSeverity(missingFile string, referencedBy string) Severity {
    ext := filepath.Ext(missingFile)

    // Critical: Missing JS/TS referenced by HTML
    if (ext == ".js" || ext == ".ts") && strings.HasSuffix(referencedBy, ".html") {
        return Critical
    }

    // Critical: Missing CSS referenced by HTML
    if ext == ".css" && strings.HasSuffix(referencedBy, ".html") {
        return Critical
    }

    // Warning: Missing images
    if ext == ".png" || ext == ".jpg" || ext == ".svg" {
        return Warning
    }

    // Info: Everything else
    return Info
}
```

---

## Auto-Fix Behavior

### When to Auto-Fix

| Condition | Action |
|-----------|--------|
| Single critical file missing | Auto-generate it |
| Multiple files missing (>3) | Ask user before proceeding |
| Non-critical missing | Log and continue |
| Syntax error | Cannot auto-fix, report to user |

### Auto-Fix Flow

```
1. Detect: index.html references app.js, but app.js doesn't exist
2. Analyze: Look at index.html to understand what app.js should do
3. Generate: Create app.js with appropriate functionality
4. Validate: Re-run checker to confirm fix worked
5. Notify: "I noticed app.js was missing and created it for you"
```

### Auto-Fix Prompt Template

```
You are completing an incomplete build. The following file is missing:

Missing file: {{filename}}
Referenced by: {{referencingFile}}
Reference context: {{surroundingCode}}

Project files that exist:
{{existingFiles}}

Based on the reference context and existing files, generate the missing file.
Focus on making the app functional - the user can refine later.
```

---

## Data Model

### CompletenessReport

```go
type CompletenessReport struct {
    ProjectID     uuid.UUID           `json:"projectId"`
    CheckedAt     time.Time           `json:"checkedAt"`
    Status        ReportStatus        `json:"status"` // "pass", "warning", "critical"
    Issues        []CompletenessIssue `json:"issues"`
    FilesChecked  int                 `json:"filesChecked"`
    AutoFixable   int                 `json:"autoFixable"`
}

type CompletenessIssue struct {
    ID            string    `json:"id"`
    Severity      string    `json:"severity"` // "critical", "warning", "info"
    Type          string    `json:"type"`     // "missing_file", "syntax_error", "broken_reference"
    MissingFile   string    `json:"missingFile,omitempty"`
    ReferencedBy  string    `json:"referencedBy,omitempty"`
    ReferenceType string    `json:"referenceType,omitempty"` // "script", "stylesheet", "import", "image"
    LineNumber    int       `json:"lineNumber,omitempty"`
    Context       string    `json:"context,omitempty"` // Surrounding code
    AutoFixable   bool      `json:"autoFixable"`
    FixApplied    bool      `json:"fixApplied"`
}

type ReportStatus string

const (
    StatusPass     ReportStatus = "pass"
    StatusWarning  ReportStatus = "warning"
    StatusCritical ReportStatus = "critical"
)
```

---

## Service Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      ChatService                             │
│                                                              │
│  1. Process message                                          │
│  2. Write files                                              │
│  3. ──► Call CompletenessChecker.Check()                    │
│  4. If critical issues:                                      │
│     └──► Call CompletenessChecker.AutoFix()                 │
│  5. Return response (with any fix notifications)            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  CompletenessChecker                         │
│                                                              │
│  Check(projectID) → CompletenessReport                      │
│    1. Get all project files                                  │
│    2. For each file, extract references                      │
│    3. Verify each reference resolves                         │
│    4. Run syntax checks                                      │
│    5. Build report                                           │
│                                                              │
│  AutoFix(issue) → bool                                       │
│    1. Build context from existing files                      │
│    2. Call Claude to generate missing file                   │
│    3. Write file                                             │
│    4. Re-validate                                            │
└─────────────────────────────────────────────────────────────┘
```

---

## API Changes

### New Endpoint: Get Completeness Report

```
GET /api/projects/:id/completeness
```

**Response:**
```json
{
  "status": "critical",
  "checkedAt": "2024-01-01T12:00:00Z",
  "filesChecked": 2,
  "issues": [
    {
      "id": "issue-1",
      "severity": "critical",
      "type": "missing_file",
      "missingFile": "app.js",
      "referencedBy": "index.html",
      "referenceType": "script",
      "lineNumber": 15,
      "context": "<script src=\"app.js\"></script>",
      "autoFixable": true,
      "fixApplied": false
    }
  ],
  "autoFixable": 1
}
```

### New Endpoint: Trigger Auto-Fix

```
POST /api/projects/:id/completeness/fix
```

**Request:**
```json
{
  "issueIds": ["issue-1"]  // Optional, fix all if empty
}
```

**Response:**
```json
{
  "fixed": ["app.js"],
  "failed": [],
  "newReport": { ... }
}
```

---

## Frontend Integration

### Preview Panel Warning

When completeness check fails, show banner above preview:

```
┌────────────────────────────────────────────────────┐
│ ⚠️  Some files are missing                          │
│                                                     │
│ • app.js (referenced by index.html)                │
│                                                     │
│ [Fix Now]  [Show Preview Anyway]                   │
└────────────────────────────────────────────────────┘
```

### File Tree Badge

Show indicator on file tree when issues exist:

```
Files (2) ⚠️
├── index.html ✓
└── styles.css ✓

Missing:
└── app.js (click to create)
```

### Chat Notification

After auto-fix, show in chat:

```
┌────────────────────────────────────────────────────┐
│ 🔧 Auto-Fix Applied                                 │
│                                                     │
│ I noticed app.js was missing and created it.       │
│ Your app should now work correctly.                │
│                                                     │
│ [View app.js]                                      │
└────────────────────────────────────────────────────┘
```

---

## Implementation Phases

### Phase 1: Detection (MVP)
- [ ] Extract references from HTML files
- [ ] Check if referenced files exist
- [ ] Build and return CompletenessReport
- [ ] Log issues (no auto-fix yet)

### Phase 2: User Notification
- [ ] Show warning banner in preview panel
- [ ] Add "Missing files" section to file tree
- [ ] Endpoint for manual "create missing file" action

### Phase 3: Auto-Fix
- [ ] Auto-fix prompt template
- [ ] Trigger fix after build completes
- [ ] Re-validate after fix
- [ ] Chat notification of fixes applied

### Phase 4: Extended Validation
- [ ] JavaScript import resolution
- [ ] CSS @import resolution
- [ ] Basic syntax validation
- [ ] Package.json validation

---

## Edge Cases

| Scenario | Handling |
|----------|----------|
| External URL reference (`<script src="https://...">`) | Skip validation |
| CDN references | Skip validation |
| Dynamic imports (`import()`) | Skip (too complex) |
| Circular references during fix | Detect and break cycle |
| Fix generates another missing file | Cap at 3 fix iterations |
| Large project (100+ files) | Run async, cache results |

---

## Success Metrics

| Metric | Target |
|--------|--------|
| % of builds with broken previews | < 5% (down from ~20%) |
| Auto-fix success rate | > 80% |
| Time to detect issues | < 500ms |
| User-reported "app doesn't work" | Reduce by 50% |

---

## Open Questions

1. **Should fixes be automatic or require user confirmation?**
   - Recommendation: Auto-fix critical issues, confirm for multiple files

2. **How to handle conflicting fixes?**
   - E.g., both HTML and CSS reference missing image differently

3. **Should we validate against a "known good" template?**
   - E.g., "This looks like a React app but missing package.json"

4. **Cost concerns with auto-fix Claude calls?**
   - Cap at N fixes per build, use smaller model for simple cases
