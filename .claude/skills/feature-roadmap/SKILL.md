---
name: feature-roadmap
description: Decompose product vision into prioritized feature roadmap. Use when user has an approved product vision and wants to create a feature roadmap, break themes into epics, or plan implementation phases. Triggers on "create roadmap", "break down the vision", "plan features".
---

# Feature Roadmap

Transform strategic themes into prioritized feature epics.

## Prerequisites

- Product vision document exists
- Strategic themes defined
- User has approved vision

## Workflow

1. **Load Context**:
   - Read `.claude/work/0_vision/product-vision.md`
   - Read `.claude/work/0_vision/strategic-themes.md`

2. **Decompose Themes**:
   - Break each theme into 2-4 feature epics
   - Target 10-15 total epics
   - Define scope and success metrics per epic

3. **Prioritize**:
   - Score using RICE, WSJF, or Value vs Effort
   - Map dependencies between epics
   - Sequence into phases

4. **Create Roadmap**:
   - Phase 1: Foundation (months 1-2)
   - Phase 2: Core (months 3-4)
   - Phase 3: Enhancement (months 5-6)

5. **Prepare PRD Briefs** for top 3-5 epics

## Output: feature-roadmap.md

```markdown
# Feature Roadmap

## Phase 1: Foundation
| Epic | Theme | Priority | Dependencies | Status |
|------|-------|----------|--------------|--------|

### Epic: [Name]
**Theme**: [Parent theme]
**Scope**: [What's included]
**Success**: [Metrics]
**PRD Ready**: Yes/No

## Phase 2: Core
[Same structure]

## Phase 3: Enhancement
[Same structure]

## Dependency Map
[Epic relationships]

## PRD Generation Plan
1. [Epic] - Ready for PRD
2. [Epic] - Ready for PRD
```

## Next Step

After roadmap approval: "Create a PRD for [epic]" → triggers prd-create skill
