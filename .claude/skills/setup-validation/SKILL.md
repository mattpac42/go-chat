---
name: setup-validation
description: Verify Garden environment is correctly configured. Use when initializing a new project, when issues are detected with Garden setup, or when user asks about setup status. Validates folder structure, required files, and configuration.
---

# Setup Validation

Verify Garden setup is complete and valid.

## Checks

### Required Structure
```
.garden/
├── PROTOCOLS.md
├── PROJECT.md
├── agents/ (5 core agents)
├── skills/
├── commands/
├── work/{backlog,active,done}
└── config/plugins.json
```

### Required Agents
developer, architect, product, platform, researcher

### Required Config
- plugins.json must be valid JSON
- settings.json if present must be valid

## Output Format

```
Garden Setup Validation

Structure: ✅
Files: ✅
Config: ✅

Status: READY
```

Or with issues:

```
Garden Setup Validation

Structure: ❌
  • Missing: work/active/

Status: ISSUES FOUND
Fix: mkdir -p .garden/work/active
```

## Remediation

For missing items, offer to:
1. Create missing folders
2. Generate files from templates
3. Initialize default configurations
