# Architect Session: Learning Journey System Design

**Date**: 2025-12-26T21:56:52
**Agent**: architect
**Task**: Design architecture for Phase 3: Learning Journey & Progression system

## Work Completed

Designed comprehensive architecture for the gamified learning system including:

1. **Database Schema** (migration 007)
   - `achievements` - Definition table with trigger configs
   - `user_progress` - Per-project learning stats and level
   - `user_achievements` - Junction table for unlocked achievements
   - `nudge_history` - Track nudge interactions

2. **Backend Services**
   - `AchievementService` - Event processing, trigger evaluation, point tracking
   - `NudgeService` - Contextual suggestion logic based on progress

3. **API Endpoints**
   - GET/POST for progress, achievements, events, nudges
   - Event-driven design where UI fires events, backend evaluates triggers

4. **Frontend Components**
   - `AchievementToast` - Celebration notification on unlock
   - `ProgressBadge` - Level and points display
   - `NudgePopover` - Contextual learning suggestions
   - `useAchievements` hook for state management

5. **Event Integration**
   - Mapped trigger points: code_viewed, tree_expanded, level_changed, etc.
   - Integration examples for FilePreviewModal, FileTree

## Decisions Made

- **Event-driven triggers**: UI fires generic events, backend evaluates all achievement conditions
- **Project-scoped progress**: Each project has independent learning state (no global user auth needed)
- **Nudge frequency**: One nudge per type per project, with cooldown periods
- **Trigger types**: event (single action), count (accumulation), milestone (complex conditions)

## Files Created

- `/workspace/.claude/work/2_active/DESIGN-learning-journey.md` - Full architecture design

## Key Artifacts

### 10 Seed Achievements
1. First Look - Viewed code for the first time
2. Curious Mind - Viewed 5 different files
3. Deep Diver - Viewed 10 different files
4. Connection Maker - Expanded tree relationships
5. Big Picture - Viewed full application tree
6. Level Up - Advanced to Level 2
7. Explorer - Reached Level 3 (technical)
8. Technologist - Reached Level 4 (developer)
9. Graduate - Exported to VS Code
10. Self-Sufficient - Made direct code edit

### Implementation Phases
1. Database & Models (1 day)
2. Repository Layer (1 day)
3. Service Layer (2 days)
4. REST API (1 day)
5. Frontend Types & Hook (1 day)
6. Frontend Components (2 days)
7. Event Integration (1 day)

**Total: 9 days estimated**

## Recommendations

1. **Start with database migration** - Create 007_learning_journey.sql first
2. **Implement core flow first** - Focus on code_viewed -> first_look achievement end-to-end
3. **Add achievements incrementally** - Start with 3-4, expand after validation
4. **Test trigger evaluation** - Unit test the trigger matching logic thoroughly
5. **Consider WebSocket** - For real-time achievement notifications (future enhancement)

## Architecture Diagram

```
Frontend                    Backend                     Database
--------                    -------                     --------
FilePreviewModal
      |
      v
recordEvent({              POST /events
  type: 'code_viewed'  ->  AchievementService  ->  user_progress
})                            |                    user_achievements
      |                       v
      |                   evaluateTrigger()
      v                       |
AchievementToast  <----  unlocked achievements
```
