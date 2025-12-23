# Strategic Product Agents - Quick Start Guide

Quick reference for using the new strategic product planning agents in the Claude Agent System.

---

## New Agents Overview

Two new strategic agents have been added to help with high-level product planning:

### 1. strategic-product-visionary
**Purpose**: Capture big picture product vision through structured discovery

**Use When**:
- Starting a major product initiative
- Need to define "what" and "why" before "how"
- Stakeholder alignment required
- Multiple competing ideas need strategic coherence

**Capabilities**:
- Conducts structured discovery interviews
- Creates product vision documents (6-12 month horizon)
- Defines strategic themes (3-5 major feature domains)
- Establishes user personas and value propositions
- Sets vision-level success metrics (OKRs, KPIs)

**Output Files**:
- `/.claude/tasks/0_vision/product-vision.md`
- `/.claude/tasks/0_vision/strategic-themes.md`

**Example Usage**:
```
User: "I want to create a product vision for our new platform"

→ Main agent delegates to strategic-product-visionary
→ Agent conducts discovery interview
→ Agent creates vision and themes documents
→ Ready for feature roadmap creation
```

---

### 2. strategic-feature-architect
**Purpose**: Decompose vision into feature epics and create prioritized roadmaps

**Use When**:
- You have an approved product vision
- Need to break vision into implementable features
- Need to sequence features based on dependencies
- Want to prioritize using objective frameworks (RICE, WSJF)

**Capabilities**:
- Decomposes strategic themes into 10-15 feature epics
- Maps technical and business dependencies
- Prioritizes features using proven frameworks
- Sequences implementation into phases
- Prepares epic briefs for PRD generation

**Output Files**:
- `/.claude/tasks/0_vision/feature-roadmap.md`

**Example Usage**:
```
User: "Create a feature roadmap from this vision"

→ Main agent delegates to strategic-feature-architect
→ Agent analyzes strategic themes
→ Agent creates prioritized roadmap with epics
→ Ready for PRD generation
```

---

## Integration with Existing Workflow

### Before (PRD-First Workflow)

```
1. Create PRD → tactical-product-manager
2. Generate Tasks → tactical-product-manager
3. Implement → Various agents
4. Complete → Archive
```

### After (Vision-First Workflow - Optional)

```
0. Vision & Roadmap (NEW - Optional for major initiatives)
   0.1. Create Vision → strategic-product-visionary
   0.2. Create Roadmap → strategic-feature-architect
   0.3. Prepare PRD priorities

1. Create PRDs → tactical-product-manager
   (Uses roadmap epics as input)

2. Generate Tasks → tactical-product-manager
3. Implement → Various agents
4. Complete → Archive
```

**Key Point**: Vision workflow is OPTIONAL and sits BEFORE the existing PRD workflow.

---

## When to Use Vision Workflow

### ✅ Use Vision Workflow:
- Major product launches or new product lines
- Strategic pivots or direction changes
- Complex initiatives with 5+ related features
- When stakeholder alignment is critical
- Multi-quarter projects
- Need to validate market opportunity

### ❌ Skip Vision Workflow:
- Single, well-defined features
- Urgent bug fixes or patches
- Tactical improvements
- Requirements are already clear
- Time-sensitive work

---

## Quick Start: 3-Step Process

### Step 1: Create Product Vision

**Request**:
```
"I want to create a product vision for [your initiative]"
```

**What Happens**:
1. Main agent delegates to **strategic-product-visionary**
2. Agent asks discovery questions (5 phases):
   - Business context
   - User understanding
   - Market positioning
   - Technical context
   - Success definition
3. Agent creates vision documents

**Time Required**: 2-4 hours (discovery + drafting)

**Output**:
- Product vision statement
- 3-5 strategic themes
- User personas with jobs-to-be-done
- Core value propositions
- Success metrics (OKRs, KPIs)

---

### Step 2: Create Feature Roadmap

**Request**:
```
"Create a feature roadmap from this vision"
```

**What Happens**:
1. Main agent delegates to **strategic-feature-architect**
2. Agent analyzes strategic themes
3. Agent decomposes into 10-15 feature epics
4. Agent maps dependencies
5. Agent prioritizes using RICE or WSJF
6. Agent sequences into phases

**Time Required**: 2-3 hours (decomposition + prioritization)

**Output**:
- Feature epics with scope and objectives
- Dependency mapping
- Priority scores with rationale
- Phased roadmap (Foundation → Core → Enhancement)
- PRD generation priorities

---

### Step 3: Create PRDs from Roadmap

**Request**:
```
"Create a PRD for [Epic Name from roadmap]"
```

**What Happens**:
1. Main agent delegates to **tactical-product-manager**
2. Agent uses epic definition from roadmap
3. Agent includes vision context
4. Agent creates detailed PRD
5. Follow existing PRD workflow (tasks → implementation)

**Time Required**: Varies per feature (standard PRD workflow)

**Output**:
- Standard PRD in `1_backlog/NNN-feature-name/`
- Strategic context included from vision
- Success metrics aligned with roadmap

---

## File Structure

```
.claude/
├── agents/
│   ├── strategic-product-visionary.md    ← NEW: Vision creation agent
│   ├── strategic-feature-architect.md    ← NEW: Roadmap creation agent
│   ├── tactical-product-manager.md       ← EXISTING: PRD creation agent
│   └── ...
│
├── templates/
│   ├── product-vision-template.md        ← NEW: Vision document template
│   ├── strategic-themes-template.md      ← NEW: Themes template
│   ├── feature-roadmap-template.md       ← NEW: Roadmap template
│   ├── prd-template.md                   ← EXISTING: PRD template
│   └── ...
│
├── tasks/
│   ├── 0_vision/                         ← NEW: Vision documents
│   │   ├── README.md                     ← NEW: Vision workflow guide
│   │   ├── product-vision.md             ← Created by visionary agent
│   │   ├── strategic-themes.md           ← Created by visionary agent
│   │   └── feature-roadmap.md            ← Created by architect agent
│   │
│   ├── 0_obe/                            ← EXISTING: Cancelled features
│   ├── 1_backlog/                        ← EXISTING: PRDs to implement
│   ├── 2_active/                         ← EXISTING: In-progress
│   └── 3_completed/                      ← EXISTING: Delivered
│
└── docs/
    ├── vision-workflow-guide.md          ← NEW: Complete workflow guide
    └── strategic-agents-quick-start.md   ← NEW: This file
```

---

## CLAUDE.md Updates

The following sections in `CLAUDE.md` have been updated:

### 1. Task-to-Agent Mapping (Section 1)
Added:
- Product vision creation → strategic-product-visionary
- Feature roadmap planning → strategic-feature-architect

### 2. System Structure (Section 2)
Added:
- Task workflow structure explanation
- 0_vision directory documentation

### 3. Agent Selection Criteria (Section 3)
Added:
- Product vision creation agent selection
- Feature roadmap planning agent selection

### 4. Advanced Workflows (Section 9)
Added:
- Complete Product Vision Workflow documentation
- Integration with existing PRD workflow
- When to use vision vs direct PRD

---

## Example: Complete Flow

### Scenario: Launching Customer Platform

**Week 1: Vision**
```
User: "Create a vision for our customer engagement platform"

→ strategic-product-visionary conducts discovery
→ Outputs:
   - product-vision.md (4 strategic themes identified)
   - strategic-themes.md (themes detailed)
```

**Week 2: Roadmap**
```
User: "Create a feature roadmap from this vision"

→ strategic-feature-architect decomposes themes
→ Outputs:
   - feature-roadmap.md (12 epics across 3 phases)
   - Epic priorities using RICE scoring
   - Dependencies mapped
```

**Week 3+: PRD Generation**
```
User: "Create a PRD for Customer Database epic"

→ tactical-product-manager creates PRD
→ Uses epic definition from roadmap
→ Includes vision context
→ Outputs: 1_backlog/001-customer-database/prd-001-customer-database.md
```

**Week 4+: Implementation**
```
User: "Generate tasks for this PRD"
User: "Start implementation"

→ Follow existing PRD workflow
→ Tasks → Implementation → Completion
```

---

## Key Concepts

### Strategic Themes
- High-level feature domains (e.g., "Customer Acquisition", "Analytics")
- Typically 3-5 themes per vision
- Each theme encompasses multiple feature epics
- Themes last 6-12 months

### Feature Epics
- Major features that support a strategic theme
- Decomposed from themes by strategic-feature-architect
- Typically 10-15 epics per roadmap
- Each epic becomes a PRD

### Dependency Mapping
- Technical dependencies (infrastructure, APIs, data)
- Feature dependencies (prerequisites)
- Business dependencies (partnerships, compliance)
- Used to sequence implementation

### Prioritization Frameworks

**RICE Scoring**:
- Reach × Impact × Confidence / Effort
- Objective scoring based on data

**WSJF (Weighted Shortest Job First)**:
- (Value + Criticality + Risk) / Size
- Balances value with speed

**Value vs Effort Matrix**:
- Quadrants: Quick Wins, Strategic Bets, Fill-Ins, Avoid
- Visual prioritization tool

---

## Agent Handoffs

### Vision → Roadmap Handoff
```
strategic-product-visionary creates:
- product-vision.md
- strategic-themes.md

↓ (Strategic themes feed into)

strategic-feature-architect consumes:
- Strategic themes
- User personas
- Success metrics

strategic-feature-architect creates:
- feature-roadmap.md
- Epic definitions
- Dependency maps
```

### Roadmap → PRD Handoff
```
strategic-feature-architect creates:
- Epic definitions
- Priority guidance
- Success metrics

↓ (Epics feed into)

tactical-product-manager consumes:
- Epic definition
- Strategic context
- Theme objectives

tactical-product-manager creates:
- Detailed PRD
- User stories
- Acceptance criteria
```

---

## Best Practices

### 1. Vision Documents
✅ Keep to 3-5 strategic themes max
✅ Base on real user research and data
✅ Define measurable success metrics
✅ Be explicit about what's OUT of scope
✅ Review and update quarterly

❌ Don't create 10+ themes
❌ Don't jump straight to specific features
❌ Don't set unrealistic timelines (> 12 months)
❌ Don't ignore technical constraints

### 2. Feature Roadmaps
✅ Create 10-15 epics (not 50 micro-features)
✅ Map all dependencies explicitly
✅ Use objective prioritization (RICE, WSJF)
✅ Plan for MVP and iterative delivery
✅ Mark PRD readiness honestly

❌ Don't create too many epics (overwhelming)
❌ Don't make everything "high priority"
❌ Don't ignore dependencies
❌ Don't plan big bang releases

### 3. Vision Lifecycle
✅ Review quarterly (every 3 months)
✅ Refresh vision every 6-12 months
✅ Archive old visions
✅ Update based on learnings
✅ Adjust roadmap as needed

❌ Don't create vision then ignore it
❌ Don't stick to plan when you should pivot
❌ Don't let vision and reality drift

---

## Common Questions

**Q: Is vision workflow required?**
A: No, it's optional. Use for major initiatives, skip for single features.

**Q: Can I skip the roadmap and go vision → PRD?**
A: Not recommended. Roadmap ensures dependencies are clear.

**Q: How long does this take?**
A: Vision (2-4 hours) + Roadmap (2-3 hours) = ~1 day for complete package

**Q: Can I have multiple visions?**
A: Yes, for different product lines or initiatives.

**Q: What if vision changes mid-execution?**
A: Update vision docs, adjust roadmap, move affected PRDs to OBE if needed.

**Q: Do I need stakeholder approval?**
A: Yes, recommended after: Vision, Roadmap, First PRD.

---

## Resources

### Documentation
- **Complete Guide**: `/.claude/docs/vision-workflow-guide.md`
- **Directory Guide**: `/.claude/tasks/0_vision/README.md`
- **CLAUDE.md Section 9**: Advanced Workflows

### Templates
- **Vision**: `/.claude/templates/product-vision-template.md`
- **Themes**: `/.claude/templates/strategic-themes-template.md`
- **Roadmap**: `/.claude/templates/feature-roadmap-template.md`

### Agent Definitions
- **Visionary**: `/.claude/agents/strategic-product-visionary.md`
- **Architect**: `/.claude/agents/strategic-feature-architect.md`
- **Product Manager**: `/.claude/agents/tactical-product-manager.md`

---

## Getting Help

### If Vision Creation Feels Stuck
- Be honest about what you don't know
- Focus on outcomes, not features
- Share assumptions vs facts
- Ask for research time if needed

### If Roadmap Feels Overwhelming
- Start with top 5 epics, add more later
- Use simple prioritization first (Value vs Effort)
- Don't perfect the plan, learn as you go
- Break large epics into smaller ones

### If Integration Feels Unclear
- Review the example flow above
- Start with one epic → one PRD
- Follow existing PRD workflow
- Reference vision docs for context

---

## Summary

**New Capabilities**:
- Strategic product vision creation
- Feature roadmap planning with dependencies
- Epic-to-PRD workflow integration

**When to Use**:
- Major initiatives (yes)
- Single features (no)

**How to Start**:
1. Request vision creation
2. Review and approve vision
3. Request roadmap creation
4. Create PRDs from top epics
5. Follow existing workflow

**Key Benefit**:
Strategic clarity before tactical execution - build the right things, in the right order, for the right reasons.

---

**Ready to start? Request: "I want to create a product vision for [your initiative]"**
