# Session 012 - Agent Personas & Discovery UX

## Date
2025-12-27

## Summary

Focused on improving agent personas and discovery UX flow:

1. **Agent Naming (Go and Grow theme)**
   - Renamed Product Guide → **Root** (discovery/foundation)
   - Renamed Designer → **Bloom** (ideas flourish into designs)
   - Renamed Developer → **Harvest** (bringing ideas to fruition)

2. **Agent Consolidation**
   - Deleted `software-tactical.md` (duplicate of `developer.md`)
   - Deleted `software-strategic.md` (duplicate of `architect.md`)
   - Updated `lineage.json` to reflect changes

3. **Test Coverage Analysis**
   - Backend: 209 test cases across 17 files
   - Frontend: 163 test cases across 5 files
   - Identified gap: ux-tactical agent was implementing code without test requirements
   - Updated `developer.md` to explicitly include frontend TDD
   - Updated `ux-tactical.md` to hand off to developer for implementation

4. **Discovery UX Improvement**
   - Replaced "waiting/polling" pattern with user-initiated CTA
   - Added "Let's solve my problem" button with seedling icon
   - Removed automatic polling for welcome message
   - User clicks button → fetch welcome message → start discovery

## Files Modified

### Agent Files
- `.claude/agents/developer.md` - Added frontend TDD, ux-tactical collaborator
- `.claude/agents/ux-tactical.md` - Clarified design-only scope, developer handoff
- Deleted: `.claude/agents/software-tactical.md`
- Deleted: `.claude/agents/software-strategic.md`
- `.claude/lineage.json` - Removed deleted agents

### Backend
- `backend/internal/service/agent_context.go` - Root/Bloom/Harvest prompts
- `backend/internal/service/prompts/discovery.go` - Root persona in discovery
- `backend/internal/service/discovery.go` - Comment update
- `backend/internal/service/discovery_test.go` - Test for Root
- `backend/internal/service/chat_test.go` - Tests for Root
- `backend/testdata/discovery/README.md` - Doc update

### Frontend
- `frontend/src/types/index.ts` - Root/Bloom/Harvest display names + colors
- `frontend/src/components/chat/MessageList.tsx` - CTA button, seedling icon
- `frontend/src/components/chat/ChatContainer.tsx` - Start discovery handler
- `frontend/src/components/discovery/DiscoveryIntroCard.tsx` - Root text
- `frontend/src/hooks/useProjects.ts` - Removed polling logic
- `frontend/src/components/HomeClient.tsx` - Pass onRefetchMessages
- `frontend/src/components/ProjectPageClient.tsx` - Pass onRefetchMessages

## Key Decisions

1. **User-initiated discovery start** - Better UX than waiting/polling
2. **Go and Grow naming** - Root → Bloom → Harvest progression
3. **Agent consolidation** - developer + architect cover all needs
4. **TDD handoff** - ux-tactical designs, developer implements with tests

## Pending

- [ ] Test the new discovery CTA flow end-to-end
- [ ] Commit all changes (~40 modified files)
- [ ] Run backend tests to verify Root assertions pass
