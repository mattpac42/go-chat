# Architect Session: PRD Auto-Generation Design

**Date**: 2025-12-26T19:18:17
**Agent**: Architect
**Task**: Design PRD auto-generation and agent context system for Go Chat

## Work Completed

Designed comprehensive architecture for automatic PRD generation from completed discovery and agent context system for specialized agents (Product Manager, Designer, Developer).

### Key Architectural Decisions

1. **One PRD per feature** (not per phase): Enables granular tracking, focused agent context, and flexible prioritization
2. **Async generation on discovery confirm**: User doesn't wait; PRDs appear progressively
3. **Claude generates PRD content**: Uses structured discovery data as input, outputs user stories and acceptance criteria
4. **Active PRD concept**: Tracks which feature is currently being worked on for focused agent context
5. **Intent-based agent routing**: User messages matched to relevant PRDs, then routed to appropriate agent

### Deliverables Created

1. **Full Design Document**: `/workspace/.claude/work/2_active/DESIGN-prd-generation.md`
   - Architecture overview with flow diagrams
   - Data model design (PRD, UserStory, AcceptanceCriterion, AgentContext)
   - SQL migration for PRD tables
   - Service interface definitions (PRDService, AgentContextService)
   - Agent persona system prompts
   - Integration points with existing Discovery and Chat services
   - Sequence diagrams for PRD generation and agent context flows
   - 7-day implementation roadmap

## Decisions Made

| Decision | Rationale |
|----------|-----------|
| One PRD per feature | Granular tracking, focused context, reduced token usage |
| Async generation | Non-blocking UX, progressive PRD availability |
| Claude for PRD content | Best-in-class generation, structured output |
| Active PRD tracking | Clear focus for agent responses |
| Keyword-based PRD matching | Simple initial implementation, can enhance with embeddings later |
| Agent selection by intent | Product (scope), Designer (UI/UX), Developer (code) |

## Files Modified

- Created: `/workspace/.claude/work/2_active/DESIGN-prd-generation.md` - Complete design specification

## Recommendations

### Immediate Next Steps (Developer)

1. **Create migration 005_prds.sql** using the provided SQL schema
2. **Add PRD model** to `backend/internal/model/prd.go`
3. **Implement PRDRepository** in `backend/internal/repository/prd.go`
4. **Implement PRDService** with Claude integration for generation
5. **Hook PRD generation** into `DiscoveryService.ConfirmDiscovery()`

### Frontend Work (UX Tactical)

1. **PRD list view** in project sidebar/drawer
2. **Active PRD badge** in chat header
3. **PRD status indicators** (generating, draft, ready, in_progress, complete)

### Future Enhancements

- WebSocket notifications for PRD generation progress
- PRD versioning for feature change tracking
- Industry-specific PRD templates
- Embedding-based PRD matching (vs keyword)

## Summary

Designed a complete PRD auto-generation system that transforms discovery data into actionable PRDs, one per feature. The agent context system routes user messages to the appropriate agent (Product Manager, Designer, Developer) with relevant PRD context. The design integrates cleanly with existing Discovery and Chat services, with a clear 7-day implementation path.
