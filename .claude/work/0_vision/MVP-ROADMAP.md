# Go Chat MVP Roadmap

**Created**: 2025-12-24
**Last Updated**: 2025-12-26
**Status**: Active
**Version**: 2.0

---

## MVP Vision

**Prove one thing**: A non-technical user can be guided through product thinking, describe an application in plain English, understand what's being built through the App Map, and deploy it - all without touching code, command lines, or infrastructure.

**Core Philosophy**: "Code Without Coding" - Users learn systems thinking, not syntax. The AI handles the "how" while users focus on the "what" and "why."

**MVP Success Statement**: By the end of validation, 5 internal users will have successfully completed guided discovery, built applications they understand, and deployed them without developer assistance.

---

## Phase Summary

| Phase | Theme | Status | Focus |
|-------|-------|--------|-------|
| 1 | Foundation | ✅ Complete | Chat UI, AI Integration, Code Generation |
| 2A | Infrastructure | ✅ Complete | GitLab, Devcontainer, App Map, File Reveal |
| 2B | Guided Discovery | 🔄 Next | Product thinking before code generation |
| 2C | Multi-Agent Team | Planned | Coordinated specialist agents |
| 3 | Learning Journey | Planned | Achievements, progression, teaching |
| 4 | Deployment & Observability | Planned | Auto-deploy, status dashboard, monitoring |

---

## Phase 1: Foundation ✅ COMPLETE

**Theme**: Core Chat + Basic Code Generation
**Status**: Complete

### Delivered
- [x] Chat UI - Mobile-responsive web interface
- [x] AI Integration - Claude API for conversations
- [x] Code Generation - AI generates application code from descriptions
- [x] Project Context - Multi-message conversation memory
- [x] Real-time Streaming - Live response display

---

## Phase 2A: Infrastructure ✅ COMPLETE

**Theme**: Devcontainer + GitLab + App Map
**Status**: Complete

### Delivered
- [x] GitLab Project Creation - Auto-create projects from generated code
- [x] Devcontainer Template - Standard devcontainer for generated apps
- [x] Auto-Push - Generated code pushed to GitLab automatically
- [x] App Map Architecture - Files organized by PURPOSE, not paths
- [x] 2-Tier File Reveal - Descriptions first, code on request
- [x] 3-Level View Progression - Purpose → Tree → Code views
- [x] File Downloads - ZIP and single file downloads
- [x] Mobile Files Panel - Responsive file browsing

---

## Phase 2B: Guided Discovery 🔄 NEXT

**Theme**: Product Thinking Before Code (Vision Theme 1)
**Goal**: Guide users through structured discovery before any code is generated

### Why This Matters
The Product Vision emphasizes "Guide before generating" as a core principle. Users who complete discovery:
- Build better-defined applications
- Have 50% fewer major change requests
- Understand their product's purpose and structure
- Feel empowered, not dependent

### Features

| Priority | Feature | Description | Done When |
|----------|---------|-------------|-----------|
| MUST | Discovery Flow UI | Chat-based guided discovery experience | User can start discovery from new project |
| MUST | Product Vision Questions | AI guides articulation of what and why | User has clear vision statement captured |
| MUST | User Persona Definition | Guided exercise for who uses the app | At least one persona defined |
| MUST | MVP Scope Definition | Collaborative feature scoping | Core features identified and prioritized |
| SHOULD | Phased Roadmap | AI helps think about build order | Features organized into phases |
| SHOULD | Discovery Summary | Visual summary of discovery outputs | User sees their product definition |
| SHOULD | Discovery → App Map Seeding | Discovery outputs seed initial structure | App Map pre-populated from discovery |
| COULD | Skip Option | Allow experienced users to skip | Returning users can bypass discovery |

### Exit Criteria
- [ ] User can start new project with guided discovery
- [ ] 5-10 minute discovery flow completes before code generation
- [ ] Vision, persona, and MVP scope captured
- [ ] Discovery outputs visible in project summary
- [ ] Works seamlessly on mobile

### Success Metrics
- Discovery completion rate: 85%+
- Time to complete: 5-10 minutes
- User satisfaction: 4.5+ out of 5
- Reduction in post-build pivots: 50%

---

## Phase 2C: Multi-Agent Team Experience

**Theme**: Coordinated Specialist Agents (Vision Theme 6)
**Goal**: Transform single-AI chat into collaborative team experience

### Features

| Priority | Feature | Description |
|----------|---------|-------------|
| MUST | Product Guide Role | Lead agent that coordinates all conversations |
| MUST | Specialist Agents | UX Expert, Architect, Developer introduced contextually |
| MUST | Visual Agent Distinction | Icons, accent colors, role labels per agent |
| MUST | One-Voice Rule | Single agent speaks at a time, Guide summarizes |
| SHOULD | Progressive Introduction | Agents introduced as relevant to phase |
| SHOULD | [NEW] Badge | First-time agent introduction indicator |
| COULD | @Mention System | Users can directly address specific agents |

### Exit Criteria
- [ ] Product Guide leads all conversations
- [ ] Specialists appear contextually (UX in design, Dev in implementation)
- [ ] Visual distinction clear but not overwhelming
- [ ] Users understand agent roles (80%+)

---

## Phase 3: Learning Journey & Progression

**Theme**: Teaching, Not Just Building (Vision Theme 5)
**Goal**: Create structured learning path from functional understanding to technical literacy

### Features

| Priority | Feature | Description |
|----------|---------|-------------|
| MUST | Achievement System | Gamified milestones for learning |
| MUST | Progression Tracking | Track user advancement through levels |
| SHOULD | Nudge System | Contextual suggestions to explore deeper |
| SHOULD | Concept Explanations | On-demand plain-language explanations |
| COULD | VS Code Export | Graduation pathway to traditional development |

### Achievement Examples
- "First Look" - Viewed code for the first time
- "Connection Maker" - Understood component relationships
- "Level Up" - Advanced from Level 1 to Level 2 view
- "Explorer" - Viewed full technical tree
- "Graduate" - Exported project to VS Code

### Exit Criteria
- [ ] Achievement system implemented
- [ ] At least 5 achievements defined and tracking
- [ ] Users can see their progression
- [ ] Nudges appear contextually

---

## Phase 4: Deployment & Observability

**Theme**: Auto-Deploy + Mobile Monitoring (Vision Themes 3 & 4)
**Goal**: Applications deploy automatically; users monitor from anywhere

### Features

| Priority | Feature | Description |
|----------|---------|-------------|
| MUST | CI/CD Pipeline | Working pipeline that builds and deploys |
| MUST | Auto-Deploy | Apps deploy without user action |
| MUST | Deployment URL | User gets link to running app |
| MUST | App Status Dashboard | View all deployed apps with health |
| MUST | Health Indicator | Green/red status per app |
| SHOULD | Build Status | Show deployment success/failure in chat |
| SHOULD | Deployment History | When app was last deployed |
| COULD | Access Logs | Basic request counts |

### Exit Criteria
- [ ] Generated apps automatically deploy
- [ ] User receives URL to access their app
- [ ] Dashboard shows all apps with status
- [ ] Works on mobile (70%+ mobile sessions)

---

## Success Metrics Summary

### Primary Success Metric
**Applications deployed by non-technical users who completed guided discovery**

### Key Performance Indicators

| KPI | Target |
|-----|--------|
| Discovery completion rate | 85%+ |
| Time from discovery to prototype | < 1 hour |
| Deployment success rate | 95%+ |
| Mobile session percentage | 70%+ |
| User satisfaction | 4+ out of 5 |
| View level progression (30 days) | 50%+ advance |

---

## Technical Decisions

### Stack
- **Frontend**: Next.js + React + Tailwind CSS
- **Backend**: Go
- **AI**: Claude API
- **Version Control**: Self-hosted GitLab
- **Deployment**: Container-based (GitLab CI)

### Architecture Principles
- Mobile-first design
- Purpose-based file organization (App Map)
- Progressive disclosure of complexity
- Plain language always (no jargon)

---

## Alignment with Product Vision

| Vision Theme | Roadmap Phase | Status |
|--------------|---------------|--------|
| Theme 1: Guided Discovery | Phase 2B | Next |
| Theme 2: Conversational Development | Phase 1 + 2A | ✅ Complete |
| Theme 3: Infrastructure Automation | Phase 2A + 4 | Partial |
| Theme 4: Mobile-First Observability | Phase 4 | Planned |
| Theme 5: Learning Journey | Phase 3 | Planned |
| Theme 6: Multi-Agent Team | Phase 2C | Planned |
| Theme 7: Multi-Tenant Platform | Post-MVP | Deferred |

---

**Document History**

| Date | Version | Changes |
|------|---------|---------|
| 2025-12-24 | 1.0 | Initial MVP roadmap |
| 2025-12-26 | 2.0 | Realigned with Product Vision themes; added Phase 2B (Guided Discovery), 2C (Multi-Agent), Phase 3 (Learning); marked Phase 1 and 2A complete |
