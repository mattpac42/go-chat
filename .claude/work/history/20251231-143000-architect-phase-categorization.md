# Architect Session: Message Phase Categorization Design

**Date**: 2025-12-31T14:30:00Z
**Agent**: architect
**Task**: Design optimal approach for message phase categorization in chat interface

## Work Completed

Analyzed current phase categorization system and designed a retrospective, cost-efficient architecture for message phase grouping.

## Current State Analysis

### Existing Implementation
- **BuildPhaseProgress.tsx**: Real-time keyword detection via `detectPhaseFromMessage()`
- **useBuildPhase.ts**: Hook that tracks current phase and groups messages
- **MessageList.tsx**: Renders either timeline or phased view
- **Database**: Messages table has `agent_type` but no `phase` column
- **Backend**: Agent selection uses keyword matching, no phase emission

### Identified Problems
1. Keyword detection runs per-message on each render (eager, stateless)
2. Single "planning" message during Discovery causes premature grouping
3. No persistence of phase assignments
4. No distinction between "current phase" and "message's phase"

## Architecture Decision Record

### ADR-001: Deferred Phase Categorization Strategy

**Context**: Users see incorrect phase groupings because keyword detection happens in real-time without context of the overall conversation flow.

**Decision**: Implement a **deferred, marker-based categorization** with optional LLM-assisted analysis.

**Consequences**:
- (+) Accurate phase groupings that reflect actual project state
- (+) Minimal cost through explicit markers + incremental analysis
- (-) Phase grouping unavailable until discovery completes
- (-) Requires backend coordination for phase markers

---

## Recommended Architecture

### 1. When to Categorize

**Trigger Points** (in priority order):
1. **Discovery Completion**: Mark all prior messages as "discovery" phase
2. **Phase Transition Markers**: Root emits explicit `[PHASE:planning]` markers
3. **Toggle to Phase View**: Re-analyze only uncategorized messages
4. **Threshold Check**: Minimum 10 messages before showing phase toggle

**Rationale**: Discovery is a well-defined boundary (exists in DB as `discovery.stage='complete'`). This provides a clean cutpoint for "everything before is discovery."

```
Timeline:
[Discovery Phase - All messages automatically tagged 'discovery']
    |
    v  <- discoveryComplete event
[Building Phase Begins]
    |
    +-- Root emits [PHASE:planning] -> messages tagged 'planning'
    |
    +-- Root emits [PHASE:building] -> messages tagged 'building'
    ...
```

### 2. How to Persist Phases

**Recommended: Hybrid Storage**

| Layer | Storage | Purpose |
|-------|---------|---------|
| Source of Truth | Database `messages.build_phase` column | Permanent phase assignment |
| Cache | React state / localStorage | Fast access, survives page refresh |
| Computed | On-demand for uncached | Fallback for legacy messages |

**Database Migration**:
```sql
ALTER TABLE messages ADD COLUMN build_phase VARCHAR(20);
CREATE INDEX idx_messages_build_phase ON messages(build_phase)
  WHERE build_phase IS NOT NULL;
```

**API Response Enhancement**:
```typescript
interface Message {
  // existing fields...
  buildPhase?: 'discovery' | 'planning' | 'building' | 'testing' | 'launch';
}
```

### 3. Incremental vs Full Re-analysis

**Strategy: Anchor-Based Incremental Analysis**

```
Phase Assignment Rules:
1. Discovery messages: All messages before discovery_complete timestamp
2. Post-discovery messages:
   a. If message has explicit phase marker -> use it
   b. If previous message has phase -> inherit it (stickiness)
   c. If none of above -> use keyword detection (fallback)
```

**Key Insight**: Phases are "sticky" - once you enter "building", you stay there unless explicitly transitioned. This avoids the current problem where a single "planning" keyword causes incorrect grouping.

### 4. Cost Analysis

| Approach | Cost | Accuracy | Latency |
|----------|------|----------|---------|
| Current keyword matching | $0 | Low | Instant |
| Root phase markers | $0 | High | N/A (integrated) |
| LLM batch analysis | $0.002-0.01/call | High | 1-2s |
| Hybrid (recommended) | $0 + optional LLM | High | Instant + optional |

**Recommendation**: Use explicit markers from Root, with LLM analysis only when user requests "re-analyze phases" action.

### 5. Threshold Recommendations

| Setting | Recommended Value | Rationale |
|---------|-------------------|-----------|
| Min messages for phase toggle | 10 | Below this, grouping adds no value |
| Min messages per phase to show | 2 | Avoid orphan sections |
| Discovery auto-categorize trigger | On `discoveryComplete` | Natural boundary |
| Phase view default state | Collapsed for non-current | Reduce cognitive load |

**Configuration**:
```typescript
const PHASE_CONFIG = {
  minMessagesForToggle: 10,
  minMessagesPerPhase: 2,
  autoCategorizeBatchSize: 50, // for LLM analysis if used
  enablePhaseViewOnDiscoveryComplete: true,
};
```

---

## Implementation Guidance

### Phase 1: Backend Phase Markers (Recommended First)

Modify Root's prompts to emit explicit phase transitions:

```go
// In agent_context.go, add to productManagerPrompt:
const phaseTransitionGuidance = `
## Phase Transitions
When the conversation naturally moves to a new phase, emit a transition marker:
- Ready to plan architecture? Say: "[Beginning planning phase]"
- Starting implementation? Say: "[Beginning building phase]"
- Moving to testing? Say: "[Beginning testing phase]"
- Ready to deploy? Say: "[Beginning launch phase]"
`
```

Frontend parses these markers:
```typescript
function detectPhaseMarker(content: string): BuildPhase | null {
  const markers: Record<string, BuildPhase> = {
    '[Beginning planning phase]': 'planning',
    '[Beginning building phase]': 'building',
    '[Beginning testing phase]': 'testing',
    '[Beginning launch phase]': 'launch',
  };
  for (const [marker, phase] of Object.entries(markers)) {
    if (content.includes(marker)) return phase;
  }
  return null;
}
```

### Phase 2: Database Persistence

1. Add `build_phase` column to messages table
2. Backfill: All messages before project's `discovery.confirmed_at` get `phase='discovery'`
3. New messages inherit phase from previous or marker

### Phase 3: Smart View Toggle

```typescript
// In BuildPhaseProgress.tsx
const canShowPhaseView = useMemo(() => {
  return (
    discoveryComplete &&
    messages.length >= PHASE_CONFIG.minMessagesForToggle
  );
}, [discoveryComplete, messages.length]);

// Hide toggle entirely until conditions met
{canShowPhaseView && (
  <SegmentedControl
    options={['Timeline', 'By Phase']}
    value={showPhasedView ? 'By Phase' : 'Timeline'}
    onChange={...}
  />
)}
```

---

## Component Diagram

```
+------------------+     +-------------------+     +------------------+
|   ChatContainer  |---->|  useBuildPhase    |---->| BuildPhaseProgress|
+------------------+     +-------------------+     +------------------+
        |                        |                         |
        v                        v                         v
+------------------+     +-------------------+     +------------------+
|   MessageList    |     | Phase Assignment  |     | PhaseSection     |
|                  |     | Service           |     | (collapsible)    |
+------------------+     +-------------------+     +------------------+
        |                        |
        v                        v
+------------------+     +-------------------+
| MessageBubble    |     | messages.build_   |
| (shows phase     |     | phase column      |
|  badge on hover) |     | (DB persistence)  |
+------------------+     +-------------------+
```

---

## Decisions Made

1. **Deferred categorization**: Wait for discovery complete before any phase grouping
2. **Explicit markers over inference**: Root emits phase transitions, avoiding keyword ambiguity
3. **Database persistence**: Phase assignments survive sessions, enable analytics
4. **Sticky phases**: Once in a phase, stay there until explicit transition
5. **Progressive disclosure**: Hide phase toggle until meaningful grouping possible

## Files to Modify

| File | Change |
|------|--------|
| `/workspace/backend/migrations/008_message_phase.sql` | Add build_phase column |
| `/workspace/backend/internal/service/agent_context.go` | Add phase guidance to prompts |
| `/workspace/frontend/src/components/chat/BuildPhaseProgress.tsx` | Add phase marker detection |
| `/workspace/frontend/src/hooks/useBuildPhase.ts` | Implement sticky phase logic |
| `/workspace/frontend/src/types/index.ts` | Add buildPhase to Message type |

## Recommendations

1. **Implement Phase 1 first** - Backend markers provide immediate accuracy improvement with minimal effort
2. **Phase toggle UX** - Consider showing a "Phase view available!" notification when discovery completes
3. **Analytics opportunity** - Stored phases enable "time spent in each phase" metrics
4. **LLM analysis** - Keep as optional "re-analyze" button for edge cases, not default flow
