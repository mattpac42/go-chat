# Guided Discovery Flow: UX Design Document

**Date**: 2025-12-25
**Author**: UX Tactical Designer
**Status**: Design Specification
**Version**: 1.0

---

## Executive Summary

The Guided Discovery Flow is a 5-10 minute AI-led conversation that occurs BEFORE any code generation. Led by the Product Guide agent, this flow helps non-technical users articulate:

1. **Product Vision** - What they want to build and why
2. **User Personas** - Who will use the application
3. **MVP Scope** - Minimum viable feature set
4. **Phased Roadmap** - Build order and priority

This document provides UX patterns, wireframes, and specifications for implementing the discovery flow within the existing Go Chat interface.

---

## Design Principles

| Principle | Application to Discovery Flow |
|-----------|-------------------------------|
| Conversation over configuration | No forms, no wizards - just natural dialogue |
| Guide before generating | Discovery completes before any code is created |
| Progressive disclosure | Reveal complexity only as users demonstrate readiness |
| Mobile-first | All designs optimized for mobile viewports first |
| Plain language always | No jargon, no technical terms during discovery |

---

## Flow Diagram

```
+==============================================================+
|                    GUIDED DISCOVERY FLOW                      |
+==============================================================+

     START
       |
       v
+-------------------------------+
|   STAGE 1: WELCOME            |  Duration: ~1 min
|   - Greet user                |
|   - Set expectations          |
|   - Ask about business/context|
+-------------------------------+
       |
       v
+-------------------------------+
|   STAGE 2: PROBLEM DISCOVERY  |  Duration: ~2 min
|   - Identify pain points      |
|   - Understand current state  |
|   - Clarify goals             |
+-------------------------------+
       |
       v
+-------------------------------+
|   STAGE 3: USER PERSONAS      |  Duration: ~2 min
|   - Who will use this?        |
|   - Different user types?     |
|   - Permission levels?        |
+-------------------------------+
       |
       v
+-------------------------------+
|   STAGE 4: MVP SCOPE          |  Duration: ~2 min
|   - Essential features        |
|   - Nice-to-haves (later)     |
|   - Priority ranking          |
+-------------------------------+
       |
       v
+-------------------------------+
|   STAGE 5: SUMMARY            |  Duration: ~1 min
|   - Recap all decisions       |
|   - Confirm understanding     |
|   - Transition to development |
+-------------------------------+
       |
       v
     DEVELOPMENT PHASE
     (Code generation begins)
```

### Alternative Paths

```
                    +------------------------+
                    |  RETURNING USER?       |
                    +------------------------+
                              |
              +---------------+---------------+
              |                               |
              v                               v
    +------------------+           +------------------+
    |  NEW PROJECT     |           |  FAST TRACK      |
    |  Full discovery  |           |  - Show template |
    |  5-10 minutes    |           |  - Quick confirm |
    |                  |           |  - 1-2 minutes   |
    +------------------+           +------------------+
```

---

## Progress Indicator Design

### Design Philosophy

The progress indicator must:
- Show users where they are in the journey
- Feel lightweight, not like a rigid form
- Be unobtrusive on mobile
- Encourage completion without creating pressure

### Progress Bar Component

#### Desktop/Tablet (768px+)

```
+==============================================================+
|  [*] Product Guide             [1/5 Welcome] ----o----o----o |
+==============================================================+
```

**Specifications**:
- Position: Right side of header, inline with agent indicator
- Style: Horizontal dots with labels
- Current stage: Filled dot + stage name
- Future stages: Empty dots only
- Completed stages: Filled dots, no checkmarks (too clinical)

#### Mobile (<768px)

```
+----------------------------------------+
| [*] Guide               [o o o o o] 1/5|
+----------------------------------------+
```

**Specifications**:
- Position: Right side of header
- Style: Compact dot progression (no labels)
- Stage number visible as "1/5"
- Tapping dots shows stage names in tooltip

### ASCII Wireframes: Progress States

**Stage 1 - Welcome**
```
Desktop:
+----------------------------------------------------------+
| [*] Product Guide        Welcome  o   o   o   o    1 of 5 |
+----------------------------------------------------------+

Mobile:
+--------------------------------------+
| [*] Guide             [o o o o o] 1/5|
+--------------------------------------+
```

**Stage 3 - Personas (midway)**
```
Desktop:
+----------------------------------------------------------+
| [*] Product Guide          o   o  Personas  o   o  3 of 5 |
+----------------------------------------------------------+

Mobile:
+--------------------------------------+
| [*] Guide             [o o o o o] 3/5|
+--------------------------------------+
```

**Stage 5 - Summary (final)**
```
Desktop:
+----------------------------------------------------------+
| [*] Product Guide          o   o   o   o  Summary  5 of 5 |
+----------------------------------------------------------+

Mobile:
+--------------------------------------+
| [*] Guide             [o o o o o] 5/5|
+--------------------------------------+
```

---

## Stage Wireframes

### Stage 1: Welcome

**Purpose**: Greet user, set expectations, learn about their context

```
+==============================================================+
|  [*] New Project                              [o o o o o] 1/5|
+==============================================================+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:30a |
| +----------------------------------------------------------+ |
| | Welcome! I'm here to help you turn your idea into a      | |
| | working application.                                     | |
| |                                                          | |
| | Before we start building, let's take a few minutes to    | |
| | understand exactly what you need. This will help us      | |
| | create something that truly fits your business.          | |
| |                                                          | |
| | First, tell me a bit about yourself - what do you do?    | |
| +----------------------------------------------------------+ |
| [purple accent bar on left]                                  |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| Type a message...                                        [>] |
+--------------------------------------------------------------+
```

**Key UX Elements**:
- Warm, welcoming tone
- Sets expectation for "a few minutes"
- Open-ended first question (not yes/no)
- No technical language

### Stage 2: Problem Discovery

**Purpose**: Understand pain points and goals

```
+==============================================================+
|  [*] New Project                              [o o o o o] 2/5|
+==============================================================+

+--------------------------------------------------------------+
|  You                                                  10:31a |
|                            +-----------------------------+   |
|                            | I run a small bakery that   |   |
|                            | does custom cake orders     |   |
|                            +-----------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:31a |
| +----------------------------------------------------------+ |
| | A custom cake bakery - that sounds creative!             | |
| |                                                          | |
| | What's the biggest challenge you face running this       | |
| | day-to-day? What takes up the most time or causes        | |
| | the most headaches?                                      | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|  You                                                  10:32a |
|                            +-----------------------------+   |
|                            | Tracking orders is chaos.   |   |
|                            | We use paper and WhatsApp   |   |
|                            | and things get lost.        |   |
|                            +-----------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:32a |
| +----------------------------------------------------------+ |
| | So order tracking is the core problem. Paper and         | |
| | WhatsApp aren't working at scale.                        | |
| |                                                          | |
| | What would a perfect day look like if this problem       | |
| | was solved? What would you be able to do that you        | |
| | can't do now?                                            | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| Type a message...                                        [>] |
+--------------------------------------------------------------+
```

**Key UX Elements**:
- Reflects back user's words (validation)
- Asks about outcomes, not features
- Future-focused ("perfect day")
- Empathetic language ("headaches", "chaos")

### Stage 3: User Personas

**Purpose**: Identify who will use the application

```
+==============================================================+
|  [*] Cake Order Manager                       [o o o o o] 3/5|
+==============================================================+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:34a |
| +----------------------------------------------------------+ |
| | Great - so we're building an order management system.    | |
| |                                                          | |
| | Now let's think about who will actually use this.        | |
| | Besides yourself, who else needs access?                 | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|  You                                                  10:35a |
|                            +-----------------------------+   |
|                            | Me and my two employees.    |   |
|                            | They take orders and I      |   |
|                            | do the baking.              |   |
|                            +-----------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:35a |
| +----------------------------------------------------------+ |
| | So three people total:                                   | |
| |                                                          | |
| | * You (the owner/baker)                                  | |
| | * Two employees who take orders                          | |
| |                                                          | |
| | Should your employees be able to do everything, or       | |
| | are there some things only you should handle - like      | |
| | seeing sales numbers or changing prices?                 | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| Type a message...                                        [>] |
+--------------------------------------------------------------+
```

**Key UX Elements**:
- Summarizes understanding with bullet points
- Asks about permissions in plain terms
- Uses concrete examples ("sales numbers", "prices")
- Project title updates as context develops

### Stage 4: MVP Scope

**Purpose**: Define essential features and priorities

```
+==============================================================+
|  [*] Cake Order Manager                       [o o o o o] 4/5|
+==============================================================+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:37a |
| +----------------------------------------------------------+ |
| | Perfect. Now let's figure out what this needs to do.     | |
| |                                                          | |
| | If you could only have THREE things in version one,      | |
| | what would be absolutely essential?                      | |
| |                                                          | |
| | (Don't worry - we can add more later. This is just       | |
| | to get started quickly.)                                 | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|  You                                                  10:38a |
|                            +-----------------------------+   |
|                            | 1. See all orders in one    |   |
|                            |    place                    |   |
|                            | 2. Add new orders easily    |   |
|                            | 3. Know when orders are due |   |
|                            +-----------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:38a |
| +----------------------------------------------------------+ |
| | Those are great priorities:                              | |
| |                                                          | |
| | 1. Order list (see everything)                           | |
| | 2. Order form (add easily)                               | |
| | 3. Due dates (know what's coming)                        | |
| |                                                          | |
| | Anything else you mentioned that we should plan for      | |
| | in a future version? Like a calendar view, or customer   | |
| | tracking?                                                | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|  You                                                  10:39a |
|                            +-----------------------------+   |
|                            | Yes! Calendar would be      |   |
|                            | amazing for planning.       |   |
|                            +-----------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:39a |
| +----------------------------------------------------------+ |
| | Perfect - we'll add calendar view to the roadmap for     | |
| | version two.                                             | |
| |                                                          | |
| | So our plan is:                                          | |
| |                                                          | |
| | Version 1 (MVP):                                         | |
| | * Order list                                             | |
| | * Order form                                             | |
| | * Due date visibility                                    | |
| |                                                          | |
| | Version 2:                                               | |
| | * Calendar view                                          | |
| |                                                          | |
| | Sound right?                                             | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| Type a message...                                        [>] |
+--------------------------------------------------------------+
```

**Key UX Elements**:
- Limits scope with "THREE things" constraint
- Explicitly defers future items ("Version 2")
- Translates user language into feature names
- Asks for confirmation before proceeding

### Stage 5: Summary and Confirmation

**Purpose**: Confirm everything before development begins

```
+==============================================================+
|  [*] Cake Order Manager                       [o o o o o] 5/5|
+==============================================================+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:40a |
| +----------------------------------------------------------+ |
| | Here's what we're going to build:                        | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|         +----------------------------------------------+     |
|         |  DISCOVERY SUMMARY                           |     |
|         +----------------------------------------------+     |
|         |                                              |     |
|         |  Project: Cake Order Manager                 |     |
|         |                                              |     |
|         |  What It Solves:                             |     |
|         |  Replaces paper and WhatsApp order           |     |
|         |  tracking with a simple digital system       |     |
|         |                                              |     |
|         |  Who Uses It:                                |     |
|         |  * You (owner/baker) - full access           |     |
|         |  * 2 employees - order management only       |     |
|         |                                              |     |
|         |  Version 1 (MVP):                            |     |
|         |  * Order list view                           |     |
|         |  * Order creation form                       |     |
|         |  * Due date tracking                         |     |
|         |                                              |     |
|         |  Coming Later:                               |     |
|         |  * Calendar view (V2)                        |     |
|         |                                              |     |
|         +----------------------------------------------+     |
|                                                              |
|         [  Edit Details  ]      [  Start Building -->  ]     |
|                                                              |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:40a |
| +----------------------------------------------------------+ |
| | Does this capture what you need? You can edit any        | |
| | details now, or we can start building!                   | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| Type a message...                                        [>] |
+--------------------------------------------------------------+
```

**Key UX Elements**:
- Visual summary card (not just text)
- Clear section headers
- Two distinct actions: Edit or Proceed
- Builds anticipation ("Start Building")
- Card appears inline in chat (not modal)

---

## Summary Card Component

### Desktop Design

```
+----------------------------------------------------------------------+
|                        DISCOVERY SUMMARY                              |
+----------------------------------------------------------------------+
|                                                                       |
|  +------------------------+  +-----------------------------------+    |
|  |  PROJECT               |  |  USERS                           |    |
|  |  Cake Order Manager    |  |  [*] You - full access           |    |
|  |                        |  |  [*] 2 employees - orders only   |    |
|  +------------------------+  +-----------------------------------+    |
|                                                                       |
|  +------------------------+  +-----------------------------------+    |
|  |  SOLVES                |  |  MVP FEATURES                    |    |
|  |  Replaces paper and    |  |  * Order list view               |    |
|  |  WhatsApp with digital |  |  * Order creation form           |    |
|  |  order tracking        |  |  * Due date tracking             |    |
|  +------------------------+  +-----------------------------------+    |
|                                                                       |
|  +-------------------------------------------------------------------+|
|  |  ROADMAP                                                          ||
|  |  V1: Order list, form, due dates                                  ||
|  |  V2: Calendar view                                                ||
|  +-------------------------------------------------------------------+|
|                                                                       |
|  +--------------------+                      +---------------------+  |
|  |   Edit Details     |                      |  Start Building --> |  |
|  +--------------------+                      +---------------------+  |
|                                                                       |
+----------------------------------------------------------------------+
```

### Mobile Design

```
+--------------------------------------+
|        DISCOVERY SUMMARY             |
+--------------------------------------+
|                                      |
|  PROJECT                             |
|  Cake Order Manager                  |
|                                      |
+--------------------------------------+
|                                      |
|  SOLVES                              |
|  Replaces paper and WhatsApp         |
|  with digital order tracking         |
|                                      |
+--------------------------------------+
|                                      |
|  USERS                               |
|  [*] You - full access               |
|  [*] 2 employees - orders only       |
|                                      |
+--------------------------------------+
|                                      |
|  MVP FEATURES                        |
|  * Order list view                   |
|  * Order creation form               |
|  * Due date tracking                 |
|                                      |
+--------------------------------------+
|                                      |
|  COMING LATER                        |
|  * Calendar view (V2)                |
|                                      |
+--------------------------------------+
|                                      |
|  +--------------------------------+  |
|  |        Edit Details            |  |
|  +--------------------------------+  |
|                                      |
|  +--------------------------------+  |
|  |    Start Building -->          |  |
|  +--------------------------------+  |
|                                      |
+--------------------------------------+
```

**Specifications**:
- Mobile: Stack all sections vertically
- Desktop: 2-column grid where appropriate
- Primary action (Start Building): Teal fill, right-aligned
- Secondary action (Edit): Ghost button, left-aligned
- Card has subtle shadow and rounded corners

---

## Edge Cases and Navigation

### 1. Go Back / Edit Previous Answer

**Scenario**: User wants to change something they said earlier

```
+--------------------------------------------------------------+
|  You                                                  10:42a |
|                            +-----------------------------+   |
|                            | Actually, I want to change  |   |
|                            | what I said about users     |   |
|                            +-----------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:42a |
| +----------------------------------------------------------+ |
| | No problem! Let's revisit who will use this.             | |
| |                                                          | |
| | You mentioned yourself and two employees. What would     | |
| | you like to change?                                      | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+
```

**Behavior**:
- No explicit "back button" - handled conversationally
- Guide recalls previous context
- Progress indicator does NOT go backward (conversation continues forward)
- If user fundamentally changes scope, Guide acknowledges and updates summary

### 2. Skip Discovery (Returning User)

**Scenario**: User has created projects before and wants to move fast

```
+==============================================================+
|  [*] New Project                                   [Fast Track]|
+==============================================================+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:30a |
| +----------------------------------------------------------+ |
| | Welcome back! I see you've built projects before.        | |
| |                                                          | |
| | Would you like to:                                       | |
| |                                                          | |
| | [Use Previous Template]  Quick start from your last      | |
| |                          project's structure             | |
| |                                                          | |
| | [Fresh Discovery]        New 5-minute conversation to    | |
| |                          define this project             | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| Type or tap an option...                                 [>] |
+--------------------------------------------------------------+
```

**Template Selection Flow (if chosen)**:

```
+--------------------------------------------------------------+
| [*] Product Guide                                     10:30a |
| +----------------------------------------------------------+ |
| | Your previous projects:                                  | |
| |                                                          | |
| | +------------------------------------------------------+ | |
| | | [1] Cake Order Manager                               | | |
| | |     Order tracking with 3 users                      | | |
| | +------------------------------------------------------+ | |
| | | [2] Inventory Tracker                                | | |
| | |     Stock management for small shop                  | | |
| | +------------------------------------------------------+ | |
| |                                                          | |
| | Which would you like to use as a starting point?         | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+
```

### 3. User Gives Too Little Information

**Scenario**: User provides one-word or vague answers

```
+--------------------------------------------------------------+
|  You                                                  10:32a |
|                            +-----------------------------+   |
|                            | Tracking                    |   |
|                            +-----------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:32a |
| +----------------------------------------------------------+ |
| | Tracking - that's a start! Help me understand a bit      | |
| | more:                                                    | |
| |                                                          | |
| | What exactly do you need to track? For example:          | |
| | - Orders and deliveries?                                 | |
| | - Inventory or stock?                                    | |
| | - Customer information?                                  | |
| | - Something else?                                        | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+
```

**Behavior**:
- Guide provides examples to prompt elaboration
- Uses bullet points to make options scannable
- Never judges or criticizes brief answers
- Gently guides toward more detail

### 4. User Gives Too Much Information

**Scenario**: User writes a paragraph describing everything

```
+--------------------------------------------------------------+
|  You                                                  10:32a |
|                            +-----------------------------+   |
|                            | I run a bakery and we do    |   |
|                            | custom cakes and we need    |   |
|                            | to track orders and also    |   |
|                            | inventory and customer info |   |
|                            | and payments and we have 3  |   |
|                            | employees and we're open    |   |
|                            | 6 days a week and...        |   |
|                            +-----------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:32a |
| +----------------------------------------------------------+ |
| | Lots of great context! Let me break this down.           | |
| |                                                          | |
| | I'm hearing a few different needs:                       | |
| | * Order tracking                                         | |
| | * Inventory management                                   | |
| | * Customer information                                   | |
| | * Payment tracking                                       | |
| |                                                          | |
| | For version one, which of these is the MOST urgent?      | |
| | We can add the others later.                             | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+
```

**Behavior**:
- Guide summarizes and organizes the information
- Helps user prioritize
- Explicitly defers non-essential items
- Uses "version one" framing to reduce scope anxiety

### 5. User Wants to Skip a Stage

**Scenario**: User says "I don't know" or tries to skip

```
+--------------------------------------------------------------+
|  You                                                  10:35a |
|                            +-----------------------------+   |
|                            | I don't really know who     |   |
|                            | will use it yet             |   |
|                            +-----------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:35a |
| +----------------------------------------------------------+ |
| | That's okay! Let's keep it simple for now.               | |
| |                                                          | |
| | I'll assume it's just you to start. We can easily add    | |
| | more users later once you figure that out.               | |
| |                                                          | |
| | Let's move on to what features you need...               | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+
```

**Behavior**:
- Accept uncertainty gracefully
- Provide a sensible default
- Explicitly note it can be changed
- Move forward without friction

### 6. User Abandons Mid-Discovery

**Scenario**: User closes browser or navigates away

**On Return (same session)**:

```
+--------------------------------------------------------------+
| [*] Product Guide                                     10:45a |
| +----------------------------------------------------------+ |
| | Welcome back! We were in the middle of defining your     | |
| | project.                                                 | |
| |                                                          | |
| | Here's what we have so far:                              | |
| | * You're building an order tracking system               | |
| | * For you and 2 employees                                | |
| |                                                          | |
| | Ready to continue where we left off?                     | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+
```

**On Return (new session, incomplete discovery)**:

```
+--------------------------------------------------------------+
| [*] Product Guide                                     10:45a |
| +----------------------------------------------------------+ |
| | Hi again! You started defining a project last time but   | |
| | didn't finish.                                           | |
| |                                                          | |
| | [Continue Previous]     Resume your order tracking       | |
| |                         project                          | |
| |                                                          | |
| | [Start Fresh]           Begin a new project              | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+
```

### 7. Edit After Summary

**Scenario**: User taps "Edit Details" on the summary card

```
+--------------------------------------------------------------+
| [*] Product Guide                                     10:42a |
| +----------------------------------------------------------+ |
| | What would you like to change?                           | |
| |                                                          | |
| | [Project Name]    "Cake Order Manager"                   | |
| | [Users]           3 users with different access          | |
| | [MVP Features]    Order list, form, due dates            | |
| | [Roadmap]         Calendar view in V2                    | |
| |                                                          | |
| | Tap a section to edit, or just tell me what to change.   | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+
```

**Behavior**:
- Show clickable sections for quick access
- Also accept natural language edits
- After edit, regenerate summary
- Maintain conversation flow (no modal dialogs)

---

## Conversation Starters and Prompts

### Opening Messages by Stage

| Stage | Guide's Opening |
|-------|-----------------|
| 1. Welcome | "Welcome! I'm here to help you turn your idea into a working application. Before we start building, let's take a few minutes to understand exactly what you need. First, tell me a bit about yourself - what do you do?" |
| 2. Problem | "Thanks for sharing that! What's the biggest challenge you face running this day-to-day? What takes up the most time or causes the most headaches?" |
| 3. Personas | "Now let's think about who will actually use this. Besides yourself, who else needs access?" |
| 4. MVP | "If you could only have THREE things in version one, what would be absolutely essential? (Don't worry - we can add more later.)" |
| 5. Summary | "Here's what we're going to build: [Summary Card]" |

### Recovery Prompts

| Situation | Guide's Response |
|-----------|------------------|
| No response for 30s | "Take your time - I'm here when you're ready." |
| User says "I don't know" | "That's okay! Let me give you some options..." |
| User goes off-topic | "That's interesting! Let me note that for later. For now, let's focus on [current stage]..." |
| User asks for help | "No problem! Right now I'm trying to understand [current stage]. Can you tell me about [specific prompt]?" |

---

## Mobile-Specific Considerations

### Touch Targets

- All tappable elements: minimum 44x44px
- Summary card buttons: full-width on mobile
- Progress dots: 32x32px tap area (even if visual is smaller)

### Keyboard Behavior

- Keyboard should not overlap active input
- "Next" button on iOS keyboard submits message
- Auto-scroll to bottom when keyboard appears

### Summary Card on Mobile

- Card spans full width (edge-to-edge with padding)
- Sections stack vertically
- Buttons stack vertically (Start Building on top)
- Card is scrollable if content exceeds viewport

### Progress Indicator Drawer

When user taps the progress dots on mobile:

```
+--------------------------------------+
|           Discovery Progress      [X]|
+--------------------------------------+
|                                      |
|  [o] Welcome           COMPLETED     |
|      Set the stage                   |
|                                      |
+--------------------------------------+
|                                      |
|  [o] Problem Discovery COMPLETED     |
|      Identified pain points          |
|                                      |
+--------------------------------------+
|                                      |
|  [@] User Personas     CURRENT       |
|      Defining who uses this          |
|                                      |
+--------------------------------------+
|                                      |
|  [ ] MVP Scope         NEXT          |
|      Essential features              |
|                                      |
+--------------------------------------+
|                                      |
|  [ ] Summary           UPCOMING      |
|      Confirm and begin               |
|                                      |
+--------------------------------------+
```

---

## Transition to Development Phase

After user confirms the summary ("Start Building"), a clear visual transition occurs:

```
+--------------------------------------------------------------+
|  You                                                  10:41a |
|                            +-----------------------------+   |
|                            | Yes, let's do it!           |   |
|                            +-----------------------------+   |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
|                     DISCOVERY COMPLETE                        |
|                                                              |
|          [checkmark animation]                               |
|                                                              |
|          Now let's bring this to life...                     |
|                                                              |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [*] Product Guide                                     10:41a |
| +----------------------------------------------------------+ |
| | Excellent! Discovery is complete. Now let me hand this   | |
| | off to our developer to start creating your first files. | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+

+--------------------------------------------------------------+
| [<>] Developer                                [NEW] - 10:42a |
| +----------------------------------------------------------+ |
| | I'm starting on your Cake Order Manager now.             | |
| |                                                          | |
| | First, I'll create the order list so you can see all     | |
| | your orders in one place...                              | |
| |                                                          | |
| | [Creating files...]                                      | |
| +----------------------------------------------------------+ |
+--------------------------------------------------------------+
```

**Transition Elements**:
- "Discovery Complete" banner appears inline (not modal)
- Brief animation or visual celebration
- Guide explicitly hands off to Developer
- Developer introduces themselves with [NEW] badge
- Progress indicator disappears or transforms to development phase indicator

---

## Data Model for Discovery State

```typescript
interface DiscoveryState {
  projectId: string;
  stage: 'welcome' | 'problem' | 'personas' | 'mvp' | 'summary' | 'complete';
  stageStartedAt: string;

  // Captured data
  businessContext?: string;    // What they do
  problemStatement?: string;   // Pain points
  goals?: string[];           // What success looks like

  users?: {
    description: string;
    count: number;
    hasPermissions: boolean;
    permissionNotes?: string;
  };

  mvpFeatures?: {
    name: string;
    priority: number;
  }[];

  futureFeatures?: {
    name: string;
    version: string;
  }[];

  summary?: {
    projectName: string;
    solvesStatement: string;
    confirmedAt?: string;
  };

  // Metadata
  isReturningUser: boolean;
  usedTemplate?: string;
  editHistory?: {
    stage: string;
    originalValue: string;
    newValue: string;
    editedAt: string;
  }[];
}
```

---

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Discovery completion rate | 85%+ | Users who reach "Start Building" |
| Time to complete discovery | 5-10 minutes | Average time from welcome to confirmation |
| User satisfaction | 4.5/5 | Post-discovery survey |
| Stage drop-off points | <10% per stage | Analytics on stage abandonment |
| Edit frequency on summary | <30% | Users who tap "Edit Details" |
| Returning user fast-track adoption | 60%+ | Returning users who choose template |

---

## Implementation Recommendations

### Phase 1: Core Discovery Flow
- Implement 5-stage conversation with Product Guide
- Basic progress indicator (dots only)
- Simple text summary (no card initially)

### Phase 2: Enhanced Summary
- Design and implement summary card component
- Add edit functionality
- Implement mobile-optimized layouts

### Phase 3: Returning User Experience
- Add template detection and selection
- Implement session recovery
- Add "fast track" option

### Phase 4: Polish
- Add transition animations
- Implement progress drawer on mobile
- Add analytics tracking for all stages

---

## Appendix: Complete Mobile Flow

```
SCREEN 1: Welcome
+--------------------------------------+
| [=] New Project       [o o o o o] 1/5|
+--------------------------------------+
|                                      |
| [*] Product Guide             10:30a |
| +----------------------------------+ |
| | Welcome! I'm here to help you   | |
| | turn your idea into a working   | |
| | application.                    | |
| |                                 | |
| | Before we start building, let's | |
| | take a few minutes to understand| |
| | exactly what you need.          | |
| |                                 | |
| | First, tell me a bit about      | |
| | yourself - what do you do?      | |
| +----------------------------------+ |
|                                      |
+--------------------------------------+
| Type a message...                [>] |
+--------------------------------------+

SCREEN 2: Problem Discovery
+--------------------------------------+
| [=] New Project       [o o o o o] 2/5|
+--------------------------------------+
|                                      |
|                    +---------------+ |
|                    | I run a       | |
|                    | bakery        | |
|                    +---------------+ |
|                                      |
| [*] Product Guide             10:31a |
| +----------------------------------+ |
| | A bakery - nice! What's the     | |
| | biggest challenge you face      | |
| | running this day-to-day?        | |
| +----------------------------------+ |
|                                      |
+--------------------------------------+
| Type a message...                [>] |
+--------------------------------------+

SCREEN 3: User Personas
+--------------------------------------+
| [=] Cake Orders       [o o o o o] 3/5|
+--------------------------------------+
|                                      |
| [*] Product Guide             10:34a |
| +----------------------------------+ |
| | Who else besides you would use  | |
| | this system?                    | |
| +----------------------------------+ |
|                                      |
|                    +---------------+ |
|                    | Me and 2      | |
|                    | employees     | |
|                    +---------------+ |
|                                      |
| [*] Product Guide             10:35a |
| +----------------------------------+ |
| | So three people total.          | |
| |                                 | |
| | Should they all have the same   | |
| | access, or are there things     | |
| | only you should handle?         | |
| +----------------------------------+ |
|                                      |
+--------------------------------------+
| Type a message...                [>] |
+--------------------------------------+

SCREEN 4: MVP Scope
+--------------------------------------+
| [=] Cake Orders       [o o o o o] 4/5|
+--------------------------------------+
|                                      |
| [*] Product Guide             10:37a |
| +----------------------------------+ |
| | If you could only have THREE    | |
| | things in version one, what     | |
| | would be essential?             | |
| +----------------------------------+ |
|                                      |
|                    +---------------+ |
|                    | Order list,   | |
|                    | add orders,   | |
|                    | due dates     | |
|                    +---------------+ |
|                                      |
| [*] Product Guide             10:38a |
| +----------------------------------+ |
| | Great priorities! Anything for  | |
| | a future version?               | |
| +----------------------------------+ |
|                                      |
+--------------------------------------+
| Type a message...                [>] |
+--------------------------------------+

SCREEN 5: Summary
+--------------------------------------+
| [=] Cake Orders       [o o o o o] 5/5|
+--------------------------------------+
|                                      |
| [*] Product Guide             10:40a |
| +----------------------------------+ |
| | Here's what we're building:     | |
| +----------------------------------+ |
|                                      |
| +----------------------------------+ |
| |      DISCOVERY SUMMARY          | |
| +----------------------------------+ |
| |                                 | |
| | PROJECT                         | |
| | Cake Order Manager              | |
| |                                 | |
| | SOLVES                          | |
| | Paper and WhatsApp chaos        | |
| |                                 | |
| | USERS                           | |
| | You + 2 employees               | |
| |                                 | |
| | MVP FEATURES                    | |
| | * Order list                    | |
| | * Order form                    | |
| | * Due dates                     | |
| |                                 | |
| | LATER                           | |
| | * Calendar (V2)                 | |
| |                                 | |
| | +------------------------------+| |
| | |       Edit Details           || |
| | +------------------------------+| |
| |                                 | |
| | +==============================+| |
| | ||    Start Building -->      ||| |
| | +==============================+| |
| |                                 | |
| +----------------------------------+ |
|                                      |
+--------------------------------------+
| Type a message...                [>] |
+--------------------------------------+
```

---

**Document History**

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2025-12-25 | 1.0 | Initial UX design specification | UX Tactical Designer |
