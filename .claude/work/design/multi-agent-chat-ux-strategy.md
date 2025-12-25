# Multi-Agent Chat Experience: Strategic UX Design

**Date**: 2025-12-25
**Author**: UX Strategic Designer
**Status**: Design Recommendation
**Version**: 1.0

---

## Executive Summary

The multi-agent chat experience transforms the current single-AI conversation into a collaborative team experience. This aligns with the "people first" philosophy by making users feel they are working with friendly human experts rather than a monolithic system. The design must balance the richness of multiple perspectives with the simplicity required for non-technical users.

**Recommendation**: Implement a **Coordinated Team Model** where a primary "Guide" agent leads the conversation and seamlessly introduces specialists when their expertise is needed, preventing chaos while maintaining the collaborative feel.

---

## Strategic Alignment

### Connection to Product Vision

| Vision Principle | Multi-Agent Support |
|------------------|---------------------|
| People first, product second | Multiple agents feel like a human team invested in the user's success |
| Teach, don't just build | Different experts can explain their domains in context |
| Guide before generating | Product Visionary leads discovery; other agents join for implementation |
| Conversation over configuration | Natural hand-offs feel like human experts collaborating |
| Progressive disclosure | Users meet agents gradually as complexity increases |

### Connection to Strategic Themes

- **Theme 1 (Guided Discovery)**: Product Visionary agent leads the discovery phase
- **Theme 2 (Conversational Development)**: Developer and Architect agents implement
- **Theme 5 (Learning Journey)**: Each agent teaches their domain, building user literacy

---

## Design Framework

### The Coordinated Team Model

Rather than multiple agents chiming in chaotically, the experience is orchestrated by a **Guide** agent who:
1. Leads the conversation and maintains context
2. Introduces specialists when relevant
3. Summarizes and synthesizes multi-agent input
4. Prevents information overload

```
+---------------------+
|       GUIDE         |  <-- Always visible, coordinates the team
| (Product Visionary) |
+----------+----------+
           |
    +------+------+------+------+
    |      |      |      |      |
    v      v      v      v      v
  UX    Arch   Dev   Research  ...

(Specialists appear contextually, one at a time)
```

---

## Visual Design Recommendations

### 1. Agent Visual Distinction

**Recommendation**: Use a consistent visual language combining avatar, color accent, and role label.

#### Agent Design Tokens

| Agent | Icon | Accent Color | Label |
|-------|------|--------------|-------|
| Product Visionary | Target/compass | Deep Purple (#7C3AED) | "Product Guide" |
| UX Designer | Layout/grid | Coral (#F97316) | "UX Expert" |
| Architect | Blueprint/structure | Steel Blue (#3B82F6) | "Technical Architect" |
| Developer | Code brackets | Emerald (#10B981) | "Developer" |
| Researcher | Magnifying glass | Amber (#F59E0B) | "Researcher" |

#### Visual Treatment

```
+------------------------------------------+
|  [Icon] Product Guide              3:42p |
|  +--------------------------------------+|
|  | Let's understand your bakery first. ||
|  | What's the biggest challenge you    ||
|  | face daily?                         ||
|  +--------------------------------------+|
|         [subtle purple left border]      |
+------------------------------------------+
```

**Key Design Decisions**:
- **Left border accent**: 3px colored border on the left of agent messages (not full bubble color) to be subtle yet distinctive
- **Avatar + Label**: Small icon (24px) with role label, not agent "name" to feel professional, not cartoonish
- **Consistent bubble style**: Same gray background as current AI messages; distinction comes from accent and label
- **Timestamp placement**: Same position as current design for consistency

#### ASCII Mockup: Message Thread

```
+----------------------------------------------------------+
|  You                                              10:30a |
|                          +-----------------------------+ |
|                          | I want to build an app for | |
|                          | my bakery                  | |
|                          +-----------------------------+ |
+----------------------------------------------------------+

+----------------------------------------------------------+
| [*] Product Guide                                 10:30a |
| +------------------------------------------------------+ |
| | Let's understand your bakery first.                  | |
| |                                                      | |
| | What's the biggest challenge you face daily?         | |
| +------------------------------------------------------+ |
| [purple accent bar on left]                              |
+----------------------------------------------------------+

+----------------------------------------------------------+
|  You                                              10:31a |
|                          +-----------------------------+ |
|                          | Tracking custom cake orders | |
|                          | is chaos                   | |
|                          +-----------------------------+ |
+----------------------------------------------------------+

+----------------------------------------------------------+
| [*] Product Guide                                 10:31a |
| +------------------------------------------------------+ |
| | So order management is the core problem.             | |
| | Who needs to use this system?                        | |
| +------------------------------------------------------+ |
+----------------------------------------------------------+

+----------------------------------------------------------+
|  You                                              10:32a |
|                          +-----------------------------+ |
|                          | Me and my two employees    | |
|                          +-----------------------------+ |
+----------------------------------------------------------+

+----------------------------------------------------------+
| [Grid] UX Expert                                  10:32a |
| +------------------------------------------------------+ |
| | For a small team, I'd suggest a simple dashboard     | |
| | view.                                                | |
| |                                                      | |
| | Should orders be visible to everyone or do you       | |
| | need permissions?                                    | |
| +------------------------------------------------------+ |
| [coral accent bar on left]                               |
+----------------------------------------------------------+
```

---

### 2. Agent Introduction System

**Problem**: Users meeting 5+ agents immediately would be overwhelming.
**Solution**: Progressive agent introduction aligned with conversation stage.

#### Introduction Flow

```
PHASE: Discovery (Guided)
+-----------------------------+
| Active: Product Guide       |
| Hidden: All others          |
+-----------------------------+
        |
        v (User defines problem/scope)

PHASE: Design Decisions
+-----------------------------+
| Active: Product Guide       |
| Introducing: UX Expert      |
| Hidden: Architect, Dev, Res |
+-----------------------------+
        |
        v (Design decisions made)

PHASE: Architecture
+-----------------------------+
| Active: Product Guide       |
| Available: UX Expert        |
| Introducing: Architect      |
| Hidden: Developer, Research |
+-----------------------------+
        |
        v (Ready to build)

PHASE: Implementation
+-----------------------------+
| Active: Product Guide       |
| Available: All introduced   |
| Active: Developer           |
+-----------------------------+
```

#### First Introduction Pattern

When a new agent joins the conversation, use a subtle introduction:

```
+----------------------------------------------------------+
| [*] Product Guide                                 10:32a |
| +------------------------------------------------------+ |
| | Now that we know who's using this, let me bring in   | |
| | our UX expert to help with the interface design.     | |
| +------------------------------------------------------+ |
+----------------------------------------------------------+

+----------------------------------------------------------+
| [Grid] UX Expert                          [NEW] - 10:32a |
| +------------------------------------------------------+ |
| | Hi! For a small team like yours, I'd suggest a       | |
| | simple dashboard view...                             | |
| +------------------------------------------------------+ |
| [coral accent bar on left]                               |
+----------------------------------------------------------+
```

**The [NEW] badge** appears only on an agent's first message and fades after the conversation continues.

---

### 3. @Mention System

**Recommendation**: Allow @mentions but make them optional and discoverable, not required.

#### Mention Behavior

```
Input Field:
+----------------------------------------------------------+
| Type a message... or @mention an expert                  |
|                                                      [>] |
+----------------------------------------------------------+

When user types "@":
+----------------------------------------------------------+
| @                                                        |
|                                                      [>] |
+----------------------------------------------------------+
| +------------------------------------------------------+ |
| | [*] Product Guide - Vision, goals, scope            | |
| | [Grid] UX Expert - Interface, user flows            | |
| | [Arch] Architect - System design, technical         | |
| | [<>] Developer - Code, implementation               | |
| | [?] Researcher - Options, analysis                  | |
| +------------------------------------------------------+ |
+----------------------------------------------------------+
```

#### When @Mentions Are Used

| Scenario | Behavior |
|----------|----------|
| No @mention | Guide routes to appropriate agent automatically |
| @Developer | Developer responds directly |
| @Developer @Architect | Both respond (rare, allow team discussion) |
| User asks technical question without @mention | Guide may hand off: "Let me ask our Developer..." |

**Mobile Consideration**: On mobile, show a quick-access bar above the keyboard:
```
+--[*]--[Grid]--[Arch]--[<>]--[?]--+
| Type a message...            [>] |
+----------------------------------+
```

---

### 4. Preventing Chaos: The "One Voice at a Time" Rule

**Critical Design Principle**: In any given response, **only one agent speaks at length**. Others may be quoted or referenced briefly.

#### Anti-Chaos Patterns

**Pattern 1: Single Agent Response (Most Common)**
```
| [Grid] UX Expert                                  10:35a |
| +------------------------------------------------------+ |
| | I recommend a simple list view with...               | |
| +------------------------------------------------------+ |
```

**Pattern 2: Hand-off Response**
```
| [*] Product Guide                                 10:35a |
| +------------------------------------------------------+ |
| | Great question about the database.                   | |
| | @Architect, can you weigh in?                        | |
| +------------------------------------------------------+ |

| [Arch] Architect                                  10:35a |
| +------------------------------------------------------+ |
| | For a small bakery app, I'd recommend...             | |
| +------------------------------------------------------+ |
```

**Pattern 3: Summary with Attribution (Rare)**
When multiple perspectives genuinely needed:
```
| [*] Product Guide                                 10:35a |
| +------------------------------------------------------+ |
| | Let me summarize what we discussed:                  | |
| |                                                      | |
| | UX: Simple dashboard for small teams                 | |
| | Architect: Lightweight database, no complexity       | |
| |                                                      | |
| | Ready to start building?                             | |
| +------------------------------------------------------+ |
```

**Never** have 3+ agents respond back-to-back in a single exchange.

---

### 5. The Lead Agent (Guide) Role

**Recommendation**: The Product Visionary serves as the "Guide" throughout the entire journey, with elevated visual presence.

#### Guide Responsibilities

1. **Initiate conversations** - Always the first to greet users
2. **Lead discovery** - Owns the guided discovery phase entirely
3. **Introduce specialists** - Brings in other agents contextually
4. **Synthesize input** - Summarizes when multiple agents contribute
5. **Keep focus** - Redirects tangents back to user goals
6. **Celebrate progress** - Acknowledges milestones

#### Visual Hierarchy

```
Header Bar (Persistent):
+----------------------------------------------------------+
| [*] Bakery Order Manager          [Team: 3 of 5 active]  |
+----------------------------------------------------------+

The Guide has slightly more prominent avatar (28px vs 24px)
```

---

### 6. Mobile Experience

#### Mobile Layout Adaptations

**Compact Agent Headers**:
```
Mobile Message (< 768px):
+------------------------------------+
| [*] Guide                   10:30a |
| +--------------------------------+ |
| | Let's understand your bakery  | |
| | first...                      | |
| +--------------------------------+ |
+------------------------------------+
```

**Quick Mention Bar**:
```
+------------------------------------+
| [*] [Grid] [Arch] [<>] [?]        |
+------------------------------------+
| Type...                       [>] |
+------------------------------------+
```

**Team Drawer** (swipe from right or tap header):
```
+------------------------------------+
|            Team                 X |
+------------------------------------+
| [*] Product Guide        ACTIVE   |
|     Guides vision and goals       |
+------------------------------------+
| [Grid] UX Expert         ACTIVE   |
|     Interface and user flows      |
+------------------------------------+
| [Arch] Architect         AVAILABLE|
|     System design                 |
+------------------------------------+
| [<>] Developer           LATER    |
|     Code implementation           |
+------------------------------------+
| [?] Researcher           LATER    |
|     Research and analysis         |
+------------------------------------+
```

**Swipe Gestures**:
- Swipe right on message: Quote/reference in reply
- Long press on agent header: View agent info

---

### 7. Agent State Indicators

#### States

| State | Meaning | Visual |
|-------|---------|--------|
| Active | Currently participating | Full color icon, no badge |
| Available | Introduced, can be @mentioned | Full color icon, subtle "ready" dot |
| Later | Not yet relevant to conversation | Grayed icon, only in team drawer |
| Thinking | Currently generating response | Pulsing icon + "typing" indicator |

```
Agent Thinking State:
+----------------------------------------------------------+
| [*] Product Guide                                        |
| +------------------------------------------------------+ |
| | ...                                    [typing dots] | |
| +------------------------------------------------------+ |
+----------------------------------------------------------+
```

---

## Risk Analysis and Mitigations

### Risk 1: Information Overload
**Risk Level**: High
**Description**: Multiple agents responding creates cognitive overload for non-technical users.
**Mitigation**:
- One-voice-at-a-time rule (strictly enforced)
- Guide summarizes multi-agent discussions
- Progressive agent introduction (never see all 5 at once initially)

### Risk 2: Confusion About Who to Address
**Risk Level**: Medium
**Description**: Users may not know which agent to talk to.
**Mitigation**:
- @mentions are optional, never required
- Guide auto-routes questions appropriately
- Clear agent role descriptions in team drawer

### Risk 3: Uncanny Valley Effect
**Risk Level**: Medium
**Description**: Agents feel fake or robotic, not like a "friendly team."
**Mitigation**:
- Role-based labels ("UX Expert") not human names ("Sarah")
- Professional but warm tone guidelines per agent
- Avoid excessive personality or quirks

### Risk 4: Lost Context
**Risk Level**: Medium
**Description**: Switching agents causes loss of conversation context.
**Mitigation**:
- All agents share full context (not siloed)
- Guide explicitly references prior discussions
- "You mentioned to our UX expert that..."

### Risk 5: Mobile Usability Degradation
**Risk Level**: Low-Medium
**Description**: Agent visuals take too much space on mobile.
**Mitigation**:
- Compact headers on mobile
- Agent icons only (no full labels in compact mode)
- Team drawer for full details

### Risk 6: Slower Response Times
**Risk Level**: Low
**Description**: Routing between agents adds latency.
**Mitigation**:
- Pre-load likely next agent during response
- Visual "handing off to..." indicator if delay occurs
- Target: <500ms for agent hand-off

---

## People First Alignment

This multi-agent design directly supports the "People First" principle:

| People First Aspect | How Multi-Agent Supports It |
|--------------------|-----------------------------|
| Feel like working with humans | Multiple experts feel more human than one omniscient AI |
| Not dependent on a black box | Users see who contributed what, building trust |
| Understanding their app | Each expert explains their domain in accessible terms |
| Empowerment | Users learn who to ask for what kind of help |
| Reduced intimidation | Specialists feel approachable, like colleagues |

---

## Implementation Considerations

### Message Data Model Extension

```typescript
interface Message {
  id: string;
  projectId: string;
  role: 'user' | 'agent';
  content: string;
  timestamp: string;
  // New fields for multi-agent
  agentType?: 'guide' | 'ux' | 'architect' | 'developer' | 'researcher';
  agentLabel?: string;  // "Product Guide", "UX Expert", etc.
  isFirstAppearance?: boolean;  // Show [NEW] badge
  handoffFrom?: string;  // Which agent handed off to this one
}
```

### Agent Routing Logic

```
1. User sends message
2. Backend analyzes message intent
3. If current agent can handle -> respond
4. If specialist needed:
   a. Current agent generates hand-off message
   b. Specialist generates response
   c. Both sent as sequential messages
5. Guide can inject summary at any point
```

### Phased Rollout Recommendation

| Phase | Scope | Timeline |
|-------|-------|----------|
| 1 | Single Guide agent with visual treatment | Week 1-2 |
| 2 | Add UX Expert for design decisions | Week 3-4 |
| 3 | Add Architect + Developer for implementation | Week 5-6 |
| 4 | Full team + @mentions + mobile optimization | Week 7-8 |

---

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| User comprehension of agent roles | 80%+ can identify role purpose | Post-session survey |
| Agent hand-off smoothness | <5% report confusion during hand-offs | User feedback |
| Time to first meaningful response | <2 seconds (same as single agent) | Analytics |
| Mobile engagement with team features | 60%+ view team drawer | Analytics |
| User preference vs. single agent | 70%+ prefer multi-agent | A/B test |

---

## Recommendations Summary

1. **Implement Coordinated Team Model** with Product Guide as lead agent
2. **Use subtle visual distinction** (left border accent, icon + label) not dramatic color differences
3. **Enforce one-voice-at-a-time** rule to prevent chaos
4. **Make @mentions optional** with auto-routing as default
5. **Introduce agents progressively** aligned with conversation phase
6. **Design mobile-first** with compact headers and team drawer
7. **Roll out incrementally** starting with Guide-only, adding specialists

---

## Next Steps

1. Create component specifications for multi-agent MessageBubble
2. Design agent introduction animation/transition
3. Define agent personality and tone guidelines
4. Prototype team drawer for mobile
5. User test with 3-5 non-technical participants

---

## Appendix: Extended ASCII Mockups

### Full Conversation Flow

```
+==============================================================+
|  [*] Bakery Order Manager                    [Team 2/5]  [=] |
+==============================================================+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:30a |
| +----------------------------------------------------------+ |
| | Welcome! I'm here to help you build something great.     | |
| |                                                          | |
| | Tell me about your business - what do you do?            | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|  You                                                  10:31a |
|                          +-------------------------------+   |
|                          | I run a bakery that does      |   |
|                          | custom cake orders            |   |
|                          +-------------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:31a |
| +----------------------------------------------------------+ |
| | Nice! Custom cakes - that sounds like a creative         | |
| | business.                                                | |
| |                                                          | |
| | What's the biggest challenge you face running this       | |
| | day-to-day?                                              | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|  You                                                  10:32a |
|                          +-------------------------------+   |
|                          | Tracking orders is chaos.     |   |
|                          | We use paper and things       |   |
|                          | get lost.                     |   |
|                          +-------------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:32a |
| +----------------------------------------------------------+ |
| | Paper-based order tracking that's getting out of hand.   | |
| | That's a problem we can definitely solve.                | |
| |                                                          | |
| | Who else besides you would use this system?              | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|  You                                                  10:33a |
|                          +-------------------------------+   |
|                          | Me and my two employees       |   |
|                          +-------------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:33a |
| +----------------------------------------------------------+ |
| | A small team of 3 - perfect for a simple solution.       | |
| |                                                          | |
| | Let me bring in our UX expert to help think through      | |
| | how this should look and work for your team.             | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [Grid] UX Expert                              [NEW] - 10:33a |
| +----------------------------------------------------------+ |
| | Hi! For a small team like yours, I'd suggest a simple    | |
| | shared dashboard - everyone sees the same orders.        | |
| |                                                          | |
| | Quick question: do all three of you need to add and      | |
| | edit orders, or mostly just view them?                   | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|  You                                                  10:34a |
|                          +-------------------------------+   |
|                          | We all need to add and edit   |   |
|                          +-------------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [Grid] UX Expert                                      10:34a |
| +----------------------------------------------------------+ |
| | Got it - everyone can add and edit. I'm picturing:       | |
| |                                                          | |
| | * An orders list (showing all current orders)            | |
| | * An order form (to add or edit an order)                | |
| | * Maybe a calendar view (to see orders by date)          | |
| |                                                          | |
| | Does that sound right for what you need?                 | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|  You                                                  10:35a |
|                          +-------------------------------+   |
|                          | Yes! The calendar would be    |   |
|                          | really helpful                |   |
|                          +-------------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:35a |
| +----------------------------------------------------------+ |
| | Great - we have a clear picture now:                     | |
| |                                                          | |
| | * Order tracking for your bakery                         | |
| | * 3-person team, everyone can add/edit                   | |
| | * List view + calendar view                              | |
| |                                                          | |
| | Ready to start building this?                            | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|  You                                                  10:35a |
|                          +-------------------------------+   |
|                          | Yes!                          |   |
|                          +-------------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:36a |
| +----------------------------------------------------------+ |
| | Excellent! Let me hand this off to our developer to      | |
| | start creating your first files.                         | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [<>] Developer                                [NEW] - 10:36a |
| +----------------------------------------------------------+ |
| | I'll create your Bakery Order Tracker now.               | |
| |                                                          | |
| | Starting with the order list and form, then adding       | |
| | the calendar view...                                     | |
| |                                                          | |
| | [Creating files...]                                      | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+==============================================================+
| [*] [Grid] [<>]                                              |
+--------------------------------------------------------------+
| Type a message... or @mention an expert                  [>] |
+--------------------------------------------------------------+
```

### Mobile Team Drawer

```
+--------------------------------+
|  Bakery Order Manager    [X]  |
+--------------------------------+
|         Your Team             |
+--------------------------------+
|                               |
| [*] Product Guide    LEADING  |
| Guides vision, goals, scope   |
|                               |
+--------------------------------+
|                               |
| [Grid] UX Expert     ACTIVE   |
| Interface and user flows      |
| "Simple dashboard for 3"      |
|                               |
+--------------------------------+
|                               |
| [<>] Developer       ACTIVE   |
| Code implementation           |
| "Creating files..."           |
|                               |
+--------------------------------+
|                               |
| [Arch] Architect     AVAILABLE|
| System design & architecture  |
| Tap to ask a question         |
|                               |
+--------------------------------+
|                               |
| [?] Researcher       LATER    |
| Research and analysis         |
| (Available when needed)       |
|                               |
+--------------------------------+
```

---

**Document History**

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2025-12-25 | 1.0 | Initial strategic design | UX Strategic Designer |
