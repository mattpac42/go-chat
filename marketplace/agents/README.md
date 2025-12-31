# Agent Library

This directory contains the complete collection of specialized AI agents for various domains. These agents are designed to be copied into other Claude Code projects as needed.

## Purpose

The `/agents/` library serves as a repository of reusable domain-expert agents. Only essential agent-management agents (garden-guide, project-navigator) remain in `.claude/agents/` to minimize context usage. Domain-specific agents are stored here and copied into projects when needed.

## Agent Organization

### Software Engineering

**Strategic Level:**
- `software-strategic.md` - Software architecture design, system design, technical strategy, and architectural decision-making

**Tactical Level:**
- `software-tactical.md` - Code implementation, TDD, debugging, code reviews, refactoring, and performance optimization

### Platform Engineering

**Strategic Level:**
- `platform-strategic.md` - Infrastructure architecture, cloud strategy, platform design, and scalability planning

**Tactical Level:**
- `platform-tactical.md` - Infrastructure deployment, container orchestration, cloud services, and platform operations

### CI/CD & DevOps

**Strategic Level:**
- `cicd-strategic.md` - DevOps strategy, pipeline architecture, deployment patterns, and automation strategy

**Tactical Level:**
- `cicd-tactical.md` - Pipeline implementation, build automation, deployment scripts, and CI/CD tooling

### Cybersecurity

**Strategic Level:**
- `cybersecurity-strategic.md` - Security architecture, threat modeling, compliance strategy, and security governance

**Tactical Level:**
- `cybersecurity-tactical.md` - Security implementation, vulnerability scanning, penetration testing, and security tooling

### Site Reliability Engineering (SRE)

**Strategic Level:**
- `sre-strategic.md` - Reliability architecture, SLO/SLI design, incident management strategy, and observability planning

**Tactical Level:**
- `sre-tactical.md` - Monitoring implementation, incident response, capacity planning, and performance tuning

### Product Management

**Strategic Level:**
- `product-manager-strategic.md` - Product strategy, market analysis, competitive positioning, and product roadmap planning
- `product-visionary.md` - Product vision creation, strategic innovation, market opportunity identification, and long-term product direction
- `product-feature-architect.md` - Feature roadmap planning, technical feasibility analysis, and feature prioritization

**Tactical Level:**
- `product-manager-tactical.md` - Product requirements documentation (PRDs), feature planning, user story creation, and backlog management

### UX/UI Design

**Strategic Level:**
- `ux-strategic.md` - Design systems, UX strategy, information architecture, and design governance

**Tactical Level:**
- `ux-tactical.md` - UI implementation, component design, accessibility, and responsive design

### Finance

**Strategic Level:**
- `finance-strategic-officer.md` - Financial strategy, capital allocation, M&A analysis, and corporate finance

**Leadership Level:**
- `finance-fpa-director.md` - Financial planning and analysis leadership, budgeting strategy, and business partnership

**Tactical Level:**
- `finance-tactical-analyst.md` - Financial modeling, forecasting, variance analysis, and reporting
- `finance-tactical-controller.md` - Accounting operations, financial controls, compliance, and financial reporting

### Business Intelligence

**Tactical Level:**
- `business-intelligence-tactical.md` - Data analytics, BI dashboards, reporting, and business metrics

### Business Development

**Strategic Level:**
- `business-strategic.md` - Business strategy, market positioning, competitive analysis, and business development planning

**Tactical Level:**
- `business-opportunity-qualifier-tactical.md` - Bid/no-bid decisions, opportunity evaluation, Pwin calculations, and competitive positioning analysis
- `business-capture-manager-tactical.md` - Capture planning, customer engagement, solution positioning, teaming strategy, and pre-RFP activities
- `business-proposal-manager-tactical.md` - Proposal management, compliance verification, RFP analysis, and proposal production coordination
- `business-contract-strategist-tactical.md` - Contract strategy, pricing models, risk allocation, and contract negotiation support

### Project Management

**Strategic Level:**
- `project-management-strategic.md` - Program portfolio strategy, governance frameworks, and organizational project management

**Tactical Level:**
- `project-management-program-evaluator-tactical.md` - Program feasibility assessment, complexity evaluation, delivery methodology selection, and resource estimation
- `project-management-project-manager-tactical.md` - Project execution, schedule management, resource coordination, and stakeholder communication
- `project-management-portfolio-analyzer-tactical.md` - Portfolio analysis, project prioritization, capacity planning, and value optimization
- `project-management-pmo-analyst-tactical.md` - PMO operations, process improvement, metrics tracking, and governance compliance

### Task Management

Specialized agents for task analysis and documentation:
- `task-analysis.md` - Task breakdown and complexity analysis
- `task-document.md` - Task documentation and specification
- `task-enhance.md` - Task improvement and refinement
- `task-generate.md` - Task generation from requirements
- `task-validate.md` - Task validation and completeness checking

### Prompt Engineering

- `prompt-optimizer.md` - Prompt optimization and refinement for AI agent effectiveness

## How to Use

### Copy Agent to Project

```bash
# Copy a single agent to your project's .claude/agents/ directory
cp agents/software-tactical.md /path/to/project/.claude/agents/

# Copy multiple agents
cp agents/software-*.md /path/to/project/.claude/agents/
cp agents/platform-*.md /path/to/project/.claude/agents/
```

### Agent Selection Guide

**For code implementation tasks:**
- Use `software-tactical.md`

**For architecture decisions:**
- Use `software-strategic.md`

**For infrastructure work:**
- Use `platform-tactical.md`

**For platform strategy:**
- Use `platform-strategic.md`

**For CI/CD pipelines:**
- Use `cicd-tactical.md`

**For security implementation:**
- Use `cybersecurity-tactical.md`

**For product planning:**
- Use `product-manager-tactical.md` (PRDs and features)
- Use `product-visionary.md` (vision and innovation)
- Use `product-feature-architect.md` (feature roadmaps)

**For UI/UX work:**
- Use `ux-tactical.md`

**For financial analysis:**
- Use `finance-tactical-analyst.md`

**For data analysis:**
- Use `business-intelligence-tactical.md`

**For capture management:**
- Use `business-capture-manager-tactical.md` (pre-RFP activities and capture planning)

**For proposal management:**
- Use `business-proposal-manager-tactical.md` (RFP response and compliance)

**For program evaluation:**
- Use `project-management-program-evaluator-tactical.md` (feasibility assessment and delivery planning)

**For project execution:**
- Use `project-management-project-manager-tactical.md` (project delivery and coordination)

## Agent Naming Convention

Agents follow a consistent naming pattern:
- `[domain]-strategic.md` - Strategic/architectural level
- `[domain]-tactical.md` - Implementation/execution level
- `[domain]-[specialty].md` - Specialized roles (e.g., finance-fpa-director)

## Context Efficiency

By keeping domain agents in this library folder instead of `.claude/agents/`, projects maintain minimal context usage. Only copy agents you actively need for your specific project.

**Token Savings:** Moving 27 agents from `.claude/agents/` to `/agents/` saves approximately 3.5k-4k tokens per Claude Code session.

## Agent Metadata

All agents include frontmatter metadata:
- `name` - Agent identifier
- `description` - Brief agent purpose
- `color` - VSCode workspace color for visual identification
- `type` - Strategic or Tactical
- `domain` - Domain specialization

See `.claude/docs/agent-color-scheme.md` for complete color coding guide.

## Maintenance

When adding new agents to this library:
1. Follow the naming convention
2. Include complete frontmatter metadata
3. Add entry to this README under appropriate domain
4. Test agent in isolation before adding to library

## Related Documentation

- `.claude/docs/agent-invocation-examples.md` - Agent invocation patterns
- `.claude/docs/agent-color-scheme.md` - Agent color coding guide
- `CLAUDE.md` - Main agent orchestration rules
