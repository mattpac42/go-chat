---
name: garden-request
description: Generate improvement requests for the Garden team. Use when discovering bugs, missing features, or workflow improvements that should be upstreamed to The Garden repository. Triggers on "/garden-request", "report to garden", "garden improvement", "upstream request".
---

# Garden Request Generator

Create structured improvement reports for The Garden team.

## Purpose

When working in a project rooted from The Garden, you may discover:
- Missing features that should be in the base template
- Workflow improvements that would benefit all Garden projects
- Bugs in skills, agents, or scripts
- Documentation gaps

This skill creates a standardized report that can be sent to The Garden maintainers.

## Usage

```
/garden-request
```

## Workflow

### 1. Gather Information

Ask the user:
1. **Title**: Brief description of the improvement
2. **Category**: Bug | Feature | Enhancement | Documentation
3. **Priority**: Low | Medium | High | Critical
4. **Description**: What's the issue or improvement?
5. **Current behavior**: How does it work now?
6. **Proposed behavior**: How should it work?

### 2. Generate Report

Create report at `.claude/work/garden-requests/GR-XXX-slug.md`

Use the next available GR number:
```bash
ls .claude/work/garden-requests/ | grep -o 'GR-[0-9]*' | sort -t'-' -k2 -n | tail -1
```

### 3. Report Template

```markdown
# Garden Request: GR-XXX

## Title
[Brief title]

## Priority
[Low | Medium | High | Critical]

## Category
[Bug | Feature | Enhancement | Documentation]

## Source Project
[Project name from lineage.json or current directory]

## Date
[Current date]

---

## Summary

[2-3 sentence summary of the request]

## Current Behavior

[How it works now / what's missing]

## Proposed Behavior

[How it should work]

## Implementation Notes

[Technical details, files to modify, code snippets if relevant]

## Files to Modify

| File | Change |
|------|--------|
| [file path] | [description of change] |

## Testing

[How to verify the fix/feature works]

## Notes

[Additional context, backwards compatibility, etc.]
```

### 4. Output

After generating, display:
```
Created: .claude/work/garden-requests/GR-XXX-title-slug.md

To submit to Garden team:
1. Copy the report contents
2. Create issue in Garden repo, or
3. Include in next sync-baseline PR
```

## Report Categories

| Category | Use When |
|----------|----------|
| **Bug** | Something is broken or not working as documented |
| **Feature** | New capability that doesn't exist |
| **Enhancement** | Improvement to existing feature |
| **Documentation** | Missing or incorrect docs |

## Priority Guidelines

| Priority | Criteria |
|----------|----------|
| **Critical** | Blocks work, security issue, data loss |
| **High** | Significant workflow friction, affects many projects |
| **Medium** | Nice to have, improves experience |
| **Low** | Minor polish, edge cases |

## File Location

All reports are stored in `.claude/work/garden-requests/` with naming:
```
GR-001-short-description.md
GR-002-another-request.md
```

## Sending to Garden

Options for submitting requests:

1. **Copy/paste** - Copy report content to Garden repo issue
2. **Sync PR** - Include `.claude/work/garden-requests/` in a sync-baseline PR
3. **Direct commit** - If you have Garden repo access, commit directly

## Example

User: "The plant command doesn't initialize beads automatically"

Generated: `GR-001-beads-auto-init.md` with full technical details on how to add beads initialization to the plant workflow.
