# Vision-First Workflow Guide

Complete guide to using the strategic product planning workflow with vision documents, feature roadmaps, and PRD generation.

---

## Overview

The Claude Agent System now supports a **vision-first workflow** that helps you plan major product initiatives strategically before diving into implementation. This workflow is optional but recommended for complex, multi-feature efforts.

### Workflow Stages

```
┌─────────────────────────────────────────────────────────────────┐
│                     VISION-FIRST WORKFLOW                        │
└─────────────────────────────────────────────────────────────────┘

Stage 0: Vision & Strategy (Optional - For Major Initiatives)
├─ 0.1: Create Product Vision → strategic-product-visionary
├─ 0.2: Create Feature Roadmap → strategic-feature-architect
└─ 0.3: Prepare PRD Priorities → Ready for Stage 1

Stage 1-4: Standard PRD Workflow (Existing)
├─ 1: Create PRDs → tactical-product-manager
├─ 2: Generate Tasks → tactical-product-manager
├─ 3: Implementation → Specialized agents
└─ 4: Completion or OBE → Archive
```

---

## When to Use Vision Workflow

### ✅ USE Vision Workflow For:

1. **Major Product Initiatives**
   - New product launches
   - Platform expansions
   - Multi-quarter projects
   - Cross-team efforts

2. **Strategic Planning Needs**
   - Need stakeholder alignment
   - Unclear product direction
   - Multiple competing priorities
   - Budget/resource planning required

3. **Complex Feature Sets**
   - 5+ related features
   - Multiple user personas
   - Cross-cutting capabilities
   - Technical architecture implications

4. **Market-Driven Development**
   - Competitive response needed
   - Market opportunity validation
   - Product-market fit assessment
   - Go-to-market planning required

### ❌ SKIP Vision Workflow For:

1. **Single Features**
   - One well-defined feature
   - Clear requirements already
   - Tactical improvements
   - Bug fixes or patches

2. **Urgent Work**
   - Critical bugs
   - Security patches
   - Emergency fixes
   - Time-sensitive work

3. **Maintenance Tasks**
   - Code refactoring
   - Technical debt
   - Performance optimization
   - Infrastructure updates

---

## Stage 0.1: Create Product Vision

### Agent: strategic-product-visionary

This agent conducts structured discovery interviews and creates comprehensive vision documents.

### Trigger Phrases

Tell the main Claude agent:
- "I want to create a product vision"
- "Help me define what we should build"
- "Create a strategic vision for [initiative]"
- "Conduct product discovery for [project]"

### Discovery Interview Process

The agent will ask structured questions in 5 phases:

**Phase 1: Business Context**
- What business problem are we solving?
- What are the key business objectives?
- Who are the key stakeholders?
- What constraints exist?

**Phase 2: User Understanding**
- Who are the primary users?
- What are their biggest pain points?
- What jobs are they trying to get done?
- What alternatives do they use?

**Phase 3: Market & Positioning**
- What is the market opportunity?
- What trends are we riding?
- What makes us different?
- What segments should we prioritize?

**Phase 4: Technical Context**
- What existing systems do we have?
- What technical capabilities do we need?
- What are our constraints?
- What integrations are required?

**Phase 5: Success Definition**
- How will we measure success?
- What are the KPIs?
- What timelines are we working with?
- What validation criteria prove we're on track?

### Output Documents

**1. product-vision.md** (Created in `/.claude/tasks/0_vision/`)

Contains:
- Vision statement (6-12 month horizon)
- Strategic context (business goals, market opportunity)
- Target users and personas (with jobs-to-be-done)
- Core value propositions
- **Strategic themes** (3-5 major feature domains)
- Success metrics (North Star, KPIs, OKRs)
- Product principles
- In-scope / out-of-scope boundaries

**2. strategic-themes.md** (Created in `/.claude/tasks/0_vision/`)

Contains:
- Theme overview matrix
- Detailed theme definitions
- Theme dependencies
- Theme-to-persona mapping
- Prioritization rationale
- Cross-cutting capabilities
- Implementation guidance

### Example Vision Output

```markdown
# Product Vision Statement

We envision a platform that empowers small business owners to manage
their entire customer lifecycle without technical expertise. By combining
intelligent automation with simple interfaces, we will eliminate the
complexity barrier and become the trusted operating system for businesses
under 50 employees.

## Strategic Themes

1. **Customer Acquisition Engine** - Help businesses find and convert customers
2. **Relationship Management** - Enable deeper customer relationships at scale
3. **Business Intelligence** - Provide actionable insights without complexity
4. **Workflow Automation** - Automate repetitive tasks intelligently
```

### Best Practices

1. **Be Honest During Discovery**
   - Share what you know AND what you don't know
   - Admit assumptions vs validated facts
   - Flag areas needing research

2. **Focus on Outcomes**
   - Emphasize user and business outcomes
   - Avoid jumping to solution features
   - Think 6-12 months out, not next sprint

3. **Limit Strategic Themes**
   - Keep to 3-5 themes maximum
   - Each theme should be distinct
   - Themes should cover 70%+ of planned work

4. **Define Clear Boundaries**
   - Be explicit about what's OUT of scope
   - Make strategic trade-offs visible
   - Set realistic timeframes

---

## Stage 0.2: Create Feature Roadmap

### Agent: strategic-feature-architect

This agent takes your vision and decomposes it into a prioritized feature roadmap.

### Trigger Phrases

After vision is approved:
- "Create a feature roadmap from this vision"
- "Break these themes into features"
- "Build a roadmap for the next 6 months"
- "Decompose the vision into epics"

### Roadmap Creation Process

The agent will:

1. **Analyze Strategic Themes**
   - Review each theme from vision
   - Identify required capabilities
   - Map to user personas and jobs

2. **Decompose into Feature Epics**
   - Create 10-15 feature epics
   - Define scope and objectives for each
   - Establish epic-level success metrics

3. **Map Dependencies**
   - Technical dependencies (APIs, platform, data)
   - Feature dependencies (prerequisites)
   - Business dependencies (partnerships, compliance)

4. **Prioritize Features**
   - Use RICE scoring (Reach × Impact × Confidence / Effort)
   - Or WSJF (Value + Criticality + Risk / Size)
   - Or Value vs Effort matrix

5. **Sequence into Phases**
   - Phase 1: Foundation (core infrastructure, platform)
   - Phase 2: Core Experience (primary user features)
   - Phase 3: Enhancement (optimization, advanced features)

6. **Prepare PRD Briefs**
   - Identify top 3-5 epics ready for PRD
   - Provide strategic context for each
   - Define success criteria

### Output Document

**feature-roadmap.md** (Created in `/.claude/tasks/0_vision/`)

Contains:
- Vision and strategy summary
- Roadmap timeline (phases and epics)
- Epic priority matrix
- Detailed epic definitions with:
  - Strategic theme alignment
  - Target users and value delivered
  - Key capabilities
  - Dependencies
  - Success metrics
  - Effort estimates
  - Priority scores
  - PRD readiness status
- Dependency map (visual)
- Prioritization framework details
- Theme coverage analysis
- MVP and phased rollout strategy
- Cross-cutting capabilities
- Success metrics by phase
- PRD generation plan

### Example Roadmap Structure

```markdown
## Phase 1: Foundation (Month 1-2)

### Epic A: Customer Database & Profiles
- Theme: Relationship Management
- Value: Enable businesses to track customer information centrally
- Capabilities: Contact management, custom fields, data import
- Dependencies: None (foundation)
- Effort: Medium | Priority: 9.0 | PRD Status: Ready

### Epic B: Email Marketing Engine
- Theme: Customer Acquisition Engine
- Value: Enable targeted email campaigns
- Capabilities: Email builder, list segmentation, analytics
- Dependencies: Epic A (needs customer database)
- Effort: Large | Priority: 8.5 | PRD Status: Ready
```

### Best Practices

1. **Keep Epics Focused**
   - Each epic should be completable in 1-2 months
   - Break large epics into smaller ones
   - Clear scope boundaries

2. **Explicit Dependencies**
   - Map all technical prerequisites
   - Identify feature sequencing needs
   - Plan for parallel workstreams

3. **Realistic Prioritization**
   - Don't make everything high priority
   - Balance quick wins with strategic bets
   - Consider resource constraints

4. **PRD Readiness**
   - Mark which epics are ready for PRD now
   - Identify what needs refinement
   - Note any blockers

---

## Stage 0.3: Transition to PRD Workflow

### Agent: tactical-product-manager

Once roadmap is approved, use it to guide PRD creation.

### Process

1. **Select Epic from Roadmap**
   - Start with highest priority epic marked "PRD Ready"
   - Use epic definition as strategic context
   - Reference theme objectives

2. **Create PRD**
   - Follow standard PRD workflow (Stage 1)
   - Include vision context in PRD header
   - Map epic capabilities to user stories
   - Use epic success metrics in PRD

3. **Continue Through PRD Workflow**
   - Generate tasks (Stage 2)
   - Implement features (Stage 3)
   - Complete or archive (Stage 4)

### Example Flow

```
User: "Create a PRD for Epic A: Customer Database"

Main Agent:
→ Delegates to tactical-product-manager
→ Provides context from:
   - Vision document (strategic themes)
   - Roadmap (epic definition)
   - Theme details (success metrics)

tactical-product-manager:
→ Creates: /.claude/tasks/1_backlog/001-customer-database/
→ Creates: prd-001-customer-database.md
→ Includes strategic context from vision and roadmap
```

---

## Complete Workflow Example

### Scenario: Launching New SaaS Platform

**Week 1-2: Vision Creation**

User: "I want to create a vision for our new SaaS platform for small businesses"

1. Main agent delegates to **strategic-product-visionary**
2. Agent conducts discovery interview (5 phases)
3. Agent creates:
   - `product-vision.md` with 4 strategic themes
   - `strategic-themes.md` with theme details
4. User reviews and approves vision

---

**Week 3: Feature Roadmap**

User: "Create a feature roadmap from this vision"

1. Main agent delegates to **strategic-feature-architect**
2. Agent analyzes 4 themes, creates 12 feature epics
3. Agent maps dependencies and priorities using RICE
4. Agent sequences into 3 phases over 6 months
5. Agent creates `feature-roadmap.md`
6. User reviews and approves roadmap

---

**Week 4: First PRD**

User: "Create a PRD for the Customer Database epic"

1. Main agent delegates to **tactical-product-manager**
2. Agent creates PRD using epic context from roadmap
3. Folder: `1_backlog/001-customer-database/`
4. File: `prd-001-customer-database.md`
5. User approves PRD

---

**Week 5: Task Generation**

User: "Generate tasks for this PRD"

1. Agent uses `2_generate-tasks.md` template
2. Creates `tasks-prd-001-customer-database.md`
3. Breaks epic into 8 parent tasks with subtasks
4. Assigns agents to each task

---

**Week 6+: Implementation**

User: "Let's start implementing"

1. Folder moves to `2_active/001-customer-database/`
2. Tasks executed one at a time
3. Progress tracked in task file
4. When complete, moves to `3_completed/`

---

**Repeat for Next Epic**

User: "Create a PRD for the Email Marketing epic"

(Repeat PRD workflow for each epic from roadmap)

---

## File Structure Reference

```
.claude/
├── tasks/
│   ├── 0_vision/
│   │   ├── README.md                    ← Guide to vision workflow
│   │   ├── product-vision.md            ← Created by strategic-product-visionary
│   │   ├── strategic-themes.md          ← Created by strategic-product-visionary
│   │   └── feature-roadmap.md           ← Created by strategic-feature-architect
│   │
│   ├── 0_obe/                          ← Cancelled/superseded features
│   │
│   ├── 1_backlog/                      ← PRDs awaiting implementation
│   │   └── 001-customer-database/
│   │       ├── prd-001-customer-database.md
│   │       └── tasks-prd-001-customer-database.md
│   │
│   ├── 2_active/                       ← Currently implementing
│   │   └── 002-email-marketing/
│   │       ├── prd-002-email-marketing.md
│   │       ├── tasks-prd-002-email-marketing.md
│   │       └── implementation-notes.md
│   │
│   └── 3_completed/                    ← Delivered features
│       └── 001-customer-database/
│           └── [all files from implementation]
│
└── templates/
    ├── product-vision-template.md       ← Vision document template
    ├── strategic-themes-template.md     ← Themes template
    ├── feature-roadmap-template.md      ← Roadmap template
    └── prd-template.md                  ← PRD template (existing)
```

---

## Agent Collaboration

### strategic-product-visionary

**Collaborates With**:
- **strategic-feature-architect** - Provides vision for roadmap decomposition
- **strategic-product-manager** - Provides vision for go-to-market strategy
- **strategic-ux-ui-designer** - Provides personas for UX strategy

**Provides**:
- Strategic themes as foundation for features
- User personas for all downstream work
- Success metrics for measuring outcomes
- Product principles for decision-making

---

### strategic-feature-architect

**Collaborates With**:
- **strategic-product-visionary** - Consumes vision to create roadmap
- **tactical-product-manager** - Provides epic briefs for PRD creation
- **strategic-software-engineer** - Validates technical feasibility
- **strategic-platform-engineer** - Identifies infrastructure dependencies

**Provides**:
- Feature epics ready for PRD generation
- Dependency maps for implementation sequencing
- Priority guidance for resource allocation
- Success metrics for each epic

---

### tactical-product-manager

**Collaborates With**:
- **strategic-feature-architect** - Consumes epic briefs to create PRDs
- **strategic-product-visionary** - References vision for strategic context
- **tactical-software-engineer** - Works on feature implementation
- **tactical-ux-ui-designer** - Collaborates on user experience

**Provides**:
- Detailed PRDs with user stories and acceptance criteria
- Task breakdowns for implementation
- Sprint planning and backlog grooming

---

## Vision Lifecycle Management

### Quarterly Reviews (Every 3 Months)

**Review Checklist**:
- [ ] Are we achieving vision-level OKRs?
- [ ] Are strategic themes still valid?
- [ ] Do we need to adjust priorities based on learnings?
- [ ] Are there new opportunities or threats?
- [ ] Should we add/remove/modify themes?

**Update Process**:
1. Review actual outcomes vs expected
2. Update `strategic-themes.md` with learnings
3. Adjust `feature-roadmap.md` priorities if needed
4. Communicate changes to stakeholders

---

### Vision Refresh (Every 6-12 Months)

**When to Refresh**:
- Current vision timeframe expires
- Major business strategy changes
- Market conditions shift significantly
- Product-market fit needs reassessment

**Refresh Process**:
1. Archive current vision documents:
   ```bash
   mkdir .claude/tasks/0_vision/archive/
   mv product-vision.md archive/2025-Q1-product-vision.md
   mv strategic-themes.md archive/2025-Q1-strategic-themes.md
   mv feature-roadmap.md archive/2025-Q1-feature-roadmap.md
   ```

2. Create new vision with **strategic-product-visionary**
3. Build new roadmap with **strategic-feature-architect**
4. Review existing PRDs - move outdated ones to OBE
5. Create new PRDs aligned with new vision

---

## Tips for Success

### 1. Vision Documents

✅ **Do**:
- Keep vision inspiring but achievable
- Limit to 3-5 strategic themes
- Base on real user research and data
- Define clear, measurable success metrics
- Be explicit about what's out of scope

❌ **Don't**:
- Create 20+ themes (too many)
- Jump straight to specific features
- Set unrealistic timelines (6-12 months is right)
- Leave success criteria vague
- Ignore technical constraints

---

### 2. Feature Roadmaps

✅ **Do**:
- Create 10-15 epics max (more is overwhelming)
- Map all dependencies explicitly
- Use objective prioritization frameworks
- Plan for MVP and iterative delivery
- Mark PRD readiness honestly

❌ **Don't**:
- Create 50+ micro-epics (too granular)
- Ignore technical dependencies
- Make everything "high priority"
- Plan for big bang releases
- Assume all epics are ready for PRD

---

### 3. Vision-to-Execution Alignment

✅ **Do**:
- Reference vision in every PRD
- Trace features back to strategic themes
- Use consistent success metrics
- Validate assumptions regularly
- Update documents based on learnings

❌ **Don't**:
- Create vision then ignore it
- Build features not in vision without updating vision
- Let vision and reality drift apart
- Set metrics then never measure
- Stick to plan when you should pivot

---

## Common Questions

### Q: Do I need to create a vision for every feature?
**A**: No. Vision workflow is for major initiatives. Single features can go straight to PRD.

### Q: Can I have multiple active visions?
**A**: Yes, for different product lines or business units. Create separate vision folders or use prefixes.

### Q: What if my vision changes mid-execution?
**A**: Update the vision documents and adjust the roadmap. Move affected PRDs to OBE if needed.

### Q: How detailed should strategic themes be?
**A**: High-level. Each theme should encompass 3-5 epics. Don't go deeper than epic-level at vision stage.

### Q: Can I skip the roadmap and go vision → PRD directly?
**A**: Not recommended. The roadmap ensures dependencies are mapped and priorities are clear.

### Q: How long does vision workflow take?
**A**: Typically:
- Vision creation: 2-4 hours (discovery + drafting)
- Roadmap creation: 2-3 hours (decomposition + prioritization)
- Total: Half day to full day for complete vision package

### Q: Should I get stakeholder approval at each stage?
**A**: Yes. Get approval after:
1. Vision document (strategic alignment)
2. Feature roadmap (priority alignment)
3. First PRD (execution alignment)

---

## Templates Available

All templates are in `/.claude/templates/`:

1. **product-vision-template.md** - Complete vision document structure
2. **strategic-themes-template.md** - Theme details and priorities
3. **feature-roadmap-template.md** - Epic definitions and roadmap
4. **prd-template.md** - Standard PRD (existing)

Each template includes:
- Complete section structure
- Example content
- Guidance on what to include
- Best practices

---

## Next Steps

### Getting Started with Vision Workflow

1. **Read** this guide thoroughly
2. **Review** example templates in `/.claude/templates/`
3. **Decide** if your initiative needs vision workflow
4. **Start** with discovery interview by requesting vision creation
5. **Follow** the stages: Vision → Roadmap → PRDs → Implementation

### Questions or Issues?

- Review `/.claude/tasks/0_vision/README.md` for directory-specific guidance
- Check `CLAUDE.md` section 9 for workflow rules
- Examine agent definitions in `/.claude/agents/` for capabilities
- Test with a small initiative before full product vision

---

**Remember**: Vision workflow is a tool to help you plan strategically. Use it when it adds value, skip it when it doesn't. The goal is better outcomes, not process for process sake.
