# UX Strategic Session: Multi-Agent Chat Experience

**Date**: 2025-12-25
**Agent**: ux-strategic
**Task**: Design multi-agent chat experience for "code without coding" platform

## Work Completed

Created comprehensive strategic UX design document for multi-agent chat experience at `/workspace/.claude/work/design/multi-agent-chat-ux-strategy.md`.

### Key Design Decisions

1. **Coordinated Team Model**: Product Guide leads conversations and introduces specialists contextually, preventing chaos while maintaining collaborative feel

2. **Visual Distinction Strategy**: Subtle left-border accent colors (3px) + icon + role label. Not dramatic color changes that would overwhelm.

3. **Agent Visual Tokens**:
   - Product Guide: Purple (#7C3AED) with target/compass icon
   - UX Expert: Coral (#F97316) with layout/grid icon
   - Architect: Steel Blue (#3B82F6) with blueprint icon
   - Developer: Emerald (#10B981) with code brackets icon
   - Researcher: Amber (#F59E0B) with magnifying glass icon

4. **One-Voice-at-a-Time Rule**: Only one agent speaks at length in any response to prevent information overload

5. **Progressive Agent Introduction**: Users meet agents gradually aligned with conversation phase (Discovery -> Design -> Architecture -> Implementation)

6. **@Mentions Optional**: Auto-routing by Guide is default; @mentions available for power users but never required

7. **Mobile-First Adaptations**: Compact headers, quick-mention bar above keyboard, swipe-accessible team drawer

## Answers to Key Questions

| Question | Recommendation |
|----------|----------------|
| How should agents be distinguished? | Left-border accent color + small icon + role label (not human names) |
| Should users @mention agents? | Optional, with discoverable UI; Guide auto-routes by default |
| How to prevent chaos? | One-voice-at-a-time rule + Guide as coordinator + progressive introduction |
| Should there be a lead agent? | Yes - Product Guide leads and introduces specialists |
| Mobile experience? | Compact headers, team drawer, quick-mention bar |
| Risks? | Information overload, confusion about who to address, uncanny valley, lost context |
| People first support? | Multiple experts feel human; users see who contributed what; builds trust |

## Decisions Made

- **Coordinated vs Autonomous**: Chose coordinated model where Guide orchestrates other agents
- **Visual subtlety**: Chose subtle accents over dramatic per-agent styling to reduce visual noise
- **Role labels vs Names**: Chose professional role labels ("UX Expert") over human names ("Sarah")
- **Progressive disclosure**: Agents introduced gradually, not all visible from start
- **Mobile parity**: Full feature parity on mobile with adapted UI, not reduced functionality

## Files Modified

- `/workspace/.claude/work/design/multi-agent-chat-ux-strategy.md`: Created comprehensive design document

## Risks Identified

1. **Information Overload** (High): Mitigated by one-voice rule and Guide summarization
2. **Who to Address Confusion** (Medium): Mitigated by optional @mentions with auto-routing
3. **Uncanny Valley** (Medium): Mitigated by role-based labels, professional tone
4. **Lost Context** (Medium): Mitigated by shared context across agents
5. **Mobile Usability** (Low-Medium): Mitigated by mobile-first compact design
6. **Response Latency** (Low): Mitigated by pre-loading and visual indicators

## Implementation Phases Recommended

1. Phase 1 (Week 1-2): Single Guide agent with multi-agent visual treatment
2. Phase 2 (Week 3-4): Add UX Expert for design decisions
3. Phase 3 (Week 5-6): Add Architect + Developer for implementation
4. Phase 4 (Week 7-8): Full team + @mentions + mobile optimization

## Recommendations

1. Begin with Guide-only prototype to validate visual treatment before adding agents
2. Create agent personality/tone guidelines before implementation
3. User test with non-technical participants early to validate team metaphor
4. Consider achievement for "meeting" each agent as part of learning progression
5. Ensure backend routing logic is robust before UI implementation

## Next Steps for Tactical UX

1. Create component specifications for AgentMessageBubble component
2. Design agent introduction animation/transition
3. Prototype team drawer for mobile in Figma or code
4. Define copy/tone guidelines per agent role
