---
name: product-vision
description: Create product vision through structured discovery. Use when user wants to define product vision, strategic direction, or start a major product initiative. Triggers on phrases like "create a vision", "define product strategy", "what should we build", or starting new product lines.
---

# Product Vision

Guide discovery and creation of product vision documents.

## Workflow

1. **Discovery Interview** (ask sequentially, one at a time):
   - "What problem are you solving and for whom?"
   - "What's the business objective or opportunity?"
   - "Who are your target users? Describe 2-3 personas."
   - "What makes your solution unique?"
   - "What does success look like in 6-12 months?"

2. **Synthesize Vision**:
   - Create vision statement (6-12 month horizon)
   - Define 3-5 strategic themes
   - Establish success metrics/OKRs
   - Set scope boundaries (in/out)

3. **Create Documents**:
   - `.claude/work/0_vision/product-vision.md`
   - `.claude/work/0_vision/strategic-themes.md`

4. **Get Approval** before proceeding to roadmap

## Output: product-vision.md

```markdown
# Product Vision: [Name]

## Vision Statement
[One paragraph, 6-12 month horizon]

## Business Objectives
- [Objective 1]
- [Objective 2]

## Target Users
### Persona 1: [Name]
[Description, needs, pain points]

## Strategic Themes
1. [Theme]: [Description]
2. [Theme]: [Description]
3. [Theme]: [Description]

## Success Metrics
| Metric | Target | Measurement |
|--------|--------|-------------|

## Scope
**In scope**: [list]
**Out of scope**: [list]
```

## Next Step

After vision approval: "Create a feature roadmap from this vision" → triggers feature-roadmap skill
