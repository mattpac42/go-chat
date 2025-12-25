# Go Chat MVP Roadmap

**Created**: 2025-12-24
**Status**: Active
**MVP Timeframe**: 8 weeks (4 phases x 2 weeks)

---

## MVP Vision

**Prove one thing**: A non-technical user can describe an application in plain English, and Go Chat will build, deploy, and let them monitor it - all without touching code, command lines, or infrastructure.

**MVP Success Statement**: By the end of 8 weeks, 3 internal users will have successfully described, deployed, and monitored at least one application each, without any developer assistance.

---

## MVP Scope Definition

### What We Are Building

The smallest possible system that demonstrates end-to-end value:
1. **Input**: User describes what they want in a chat interface
2. **Process**: AI generates code, creates project in GitLab, deploys automatically
3. **Output**: Running application that user can access and monitor from their phone

### What We Are NOT Building (Deferred)

- Multi-turn refinement and iteration (v1 is single request -> deployed app)
- Multiple environment management (production only)
- Advanced monitoring/alerting (basic status only)
- User authentication/multi-tenancy (single user/operator)
- Database provisioning (static apps or bring-your-own-data initially)
- Native mobile apps (responsive web only)

### MVP User

**Single Persona Focus**: Sam the Small Business Owner
- Needs to build simple tools (inventory tracker, customer list, basic dashboard)
- Checks on things from phone throughout day
- Zero tolerance for technical complexity

---

## Phased Feature Breakdown

### Phase 1: Foundation (Weeks 1-2)
**Theme**: Core Chat + Basic Code Generation
**Goal**: User can have a conversation that produces working code

| Priority | Feature | Description | Done When |
|----------|---------|-------------|-----------|
| MUST | Chat UI | Simple web chat interface, mobile-responsive | User can type messages and see responses |
| MUST | AI Integration | Connect to AI model for conversation | AI responds to messages with context |
| MUST | Code Generation | AI generates application code from description | Describing "inventory tracker" produces runnable code |
| SHOULD | Project Context | AI remembers conversation context | Multi-message conversations work |
| DEFER | Clarifying Questions | AI asks follow-up questions | - |
| DEFER | Real-time Preview | Live preview during generation | - |

**Phase 1 Exit Criteria**:
- User can describe an app in chat
- AI generates appropriate code
- Code is viewable in the interface
- Works on mobile browser

---

### Phase 2: Infrastructure (Weeks 3-4)
**Theme**: Devcontainer + GitLab Integration
**Goal**: Generated code automatically becomes a real project

| Priority | Feature | Description | Done When |
|----------|---------|-------------|-----------|
| MUST | GitLab Project Creation | Auto-create project from generated code | New GitLab project exists after generation |
| MUST | Devcontainer Template | Standard devcontainer for generated apps | Project has working devcontainer.json |
| MUST | Auto-Push | Generated code pushed to GitLab automatically | Code appears in GitLab repo |
| SHOULD | Basic .gitlab-ci.yml | Skeleton CI/CD pipeline in generated projects | Pipeline file exists in repo |
| DEFER | Environment Variables | Secure config management | - |
| DEFER | Multiple Languages | Support beyond primary stack | - |

**Phase 2 Exit Criteria**:
- Generating an app creates a GitLab project
- Project has devcontainer and can be opened in VS Code
- Code is committed and pushed automatically
- No user interaction with GitLab required

---

### Phase 3: Deployment (Weeks 5-6)
**Theme**: Automated Pipelines
**Goal**: Projects deploy to production automatically

| Priority | Feature | Description | Done When |
|----------|---------|-------------|-----------|
| MUST | CI/CD Pipeline | Working pipeline that builds and deploys | Pushing code triggers deployment |
| MUST | Auto-Deploy | Apps deploy without user action | App is accessible at a URL after generation |
| MUST | Deployment URL | User gets link to their running app | Chat shows "Your app is live at [URL]" |
| SHOULD | Build Status | Show if deployment succeeded/failed | User sees success/failure in chat |
| DEFER | Rollback | Revert to previous version | - |
| DEFER | Staging Environment | Preview before production | - |

**Phase 3 Exit Criteria**:
- Generated apps automatically deploy
- User receives URL to access their app
- Deployment happens in under 10 minutes
- 90%+ deployment success rate

---

### Phase 4: Observability (Weeks 7-8)
**Theme**: Mobile-Friendly Monitoring
**Goal**: Users can check on their apps from anywhere

| Priority | Feature | Description | Done When |
|----------|---------|-------------|-----------|
| MUST | App Status Dashboard | Simple view of deployed apps | User sees list of their apps with status |
| MUST | Health Indicator | Is the app running or not | Green/red indicator per app |
| MUST | Mobile-Optimized View | Dashboard works great on phone | Fully usable on mobile browser |
| SHOULD | Deployment History | When was app last deployed | Timestamp visible |
| SHOULD | Access Logs | Basic request counts | User sees "X visitors today" |
| DEFER | Error Alerts | Push notifications for problems | - |
| DEFER | Performance Metrics | Response times, resources | - |

**Phase 4 Exit Criteria**:
- Dashboard shows all user's apps
- Health status visible at a glance
- Works seamlessly on mobile phone
- Non-technical user understands all information shown

---

## Sprint Structure

```
Week 1-2: Phase 1 (Foundation)
├── Sprint Goal: "Describe it, generate it"
├── Demo: Generate app code from chat conversation
└── Risk: AI code quality, prompt engineering

Week 3-4: Phase 2 (Infrastructure)
├── Sprint Goal: "Generate it, project it"
├── Demo: Chat generates complete GitLab project
└── Risk: GitLab API integration, devcontainer setup

Week 5-6: Phase 3 (Deployment)
├── Sprint Goal: "Project it, deploy it"
├── Demo: Chat to live URL in under 10 minutes
└── Risk: CI/CD reliability, deployment infrastructure

Week 7-8: Phase 4 (Observability)
├── Sprint Goal: "Deploy it, monitor it"
├── Demo: Complete flow on mobile device
└── Risk: Mobile UX, data aggregation
```

---

## Success Metrics for MVP

### Primary Success Metric
**3 internal users deploy 1+ applications without developer help**

### Supporting Metrics

| Metric | Target | Measurement Method |
|--------|--------|-------------------|
| End-to-end time | < 30 minutes | Timestamp: first message to live URL |
| Deployment success rate | > 90% | Successful deploys / total attempts |
| Mobile usability | 100% functional | All critical paths work on phone |
| User satisfaction | 4+ / 5 | Post-session survey |
| Zero technical support | 0 escalations | No developer intervention required |

### Validation Questions (Answer by MVP End)

1. Can non-technical users describe apps clearly enough?
2. Does AI generate usable code from descriptions?
3. Is the invisible infrastructure actually invisible?
4. Do users check status on mobile?
5. Is the value proposition proven?

---

## Technical Decisions for MVP

### Simplifying Constraints

| Decision | Rationale |
|----------|-----------|
| Single app type (web dashboard) | Prove concept before expanding |
| One programming language/stack | Reduce AI complexity |
| Production only (no staging) | Minimum viable deployment |
| Single user (no auth) | Defer multi-tenancy |
| Pre-configured templates | Reliability over flexibility |

### Technology Stack (Recommended)

- **Chat UI**: React + responsive CSS
- **AI**: Claude API for code generation
- **Backend**: Go (project namesake, simple deployment)
- **GitLab**: Self-hosted, API integration
- **CI/CD**: GitLab CI with Docker
- **Deployment**: Container-based (Kubernetes or Docker Compose)
- **Monitoring**: Basic health checks + log aggregation

---

## Risk Mitigation

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| AI generates broken code | High | High | Constrained templates, validation before deploy |
| GitLab API complexity | Medium | Medium | Spike early in Phase 2 |
| Deployment failures | Medium | High | Robust error handling, clear failure messages |
| Mobile UX issues | Low | Medium | Mobile-first design from Phase 1 |
| Scope creep | High | High | Strict DEFER discipline, weekly scope review |

---

## Definition of Done: MVP Complete

The MVP is complete when:

1. **Functional Completeness**
   - [ ] User can describe app in chat on mobile device
   - [ ] AI generates appropriate code
   - [ ] GitLab project created automatically
   - [ ] App deploys without user action
   - [ ] Live URL provided to user
   - [ ] User can check app status on phone

2. **Quality Gates**
   - [ ] 90%+ deployment success rate
   - [ ] < 30 minute end-to-end time
   - [ ] Works on iOS Safari and Android Chrome
   - [ ] No technical jargon visible to users

3. **Validation Complete**
   - [ ] 3 internal users have deployed apps
   - [ ] Zero developer assistance required
   - [ ] User satisfaction 4+ out of 5
   - [ ] Key assumptions documented as validated/invalidated

---

## Post-MVP Priorities

Once MVP is validated, immediate next priorities:

1. **Iteration capability** - Change apps through conversation
2. **Error handling** - Plain language explanations when things fail
3. **Multiple app types** - Expand beyond dashboards
4. **User accounts** - Multi-user support
5. **Alerting** - Notify users of problems

---

## Appendix: Feature Mapping to Vision Themes

| MVP Phase | Vision Theme | Features Included | Features Deferred |
|-----------|--------------|-------------------|-------------------|
| Phase 1 | Theme 1: Conversational Development | Chat UI, AI integration, code generation | Clarifying questions, real-time preview, multi-turn refinement |
| Phase 2 | Theme 2: Infrastructure Automation | GitLab integration, devcontainer, auto-push | Environment management, database provisioning |
| Phase 3 | Theme 2: Infrastructure Automation | CI/CD pipeline, auto-deploy | Rollback, staging |
| Phase 4 | Theme 3: Mobile Observability | Status dashboard, health indicator | Alerts, performance metrics |

---

**Document History**

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2025-12-24 | 1.0 | Initial MVP roadmap | Product Strategy |
