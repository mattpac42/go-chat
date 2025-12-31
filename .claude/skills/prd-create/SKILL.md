---
name: prd-create
description: Create Product Requirements Document for a feature. Use when user wants to create a PRD, define feature requirements, or document a feature for implementation. Triggers on "create PRD", "write requirements", "document this feature", "PRD for [feature]".
---

# PRD Creation

Guide structured creation of Product Requirements Documents.

## Workflow

1. **Gather Context**:
   - Check for vision/roadmap context
   - If from roadmap, load epic definition
   - If standalone, conduct discovery

2. **Discovery Questions** (ask sequentially):
   - "What problem does this feature solve?"
   - "Who are the target users?"
   - "What's the expected outcome/success metric?"
   - "Are there constraints or dependencies?"
   - "What's the scope boundary?"

3. **Create PRD Folder**:
   - Determine next number (001, 002, etc.)
   - Create `.claude/work/1_backlog/[NNN]-[feature-name]/`

4. **Write PRD**:
   - Use structured template
   - Include all sections
   - Get user approval

## Output: prd-[feature].md

```markdown
# PRD: [Feature Name]

**Version**: 1.0
**Status**: Draft

## Problem Statement
[What problem, why it matters]

## Proposed Solution
[High-level approach]

## User Stories
1. As a [user], I want to [action] so that [benefit]

## Requirements

### Functional
| ID | Requirement | Priority |
|----|-------------|----------|
| FR-1 | [Requirement] | High |

### Non-Functional
| ID | Requirement | Criteria |
|----|-------------|----------|

## Success Metrics
| Metric | Target |
|--------|--------|

## Scope
**In**: [list]
**Out**: [list]

## Dependencies
- [Dependency]
```

## File Location

`.claude/work/1_backlog/[NNN]-[feature-name]/prd-[feature-name].md`

## Next Step

After PRD approval: "Generate tasks for this PRD" → triggers task-generate skill
