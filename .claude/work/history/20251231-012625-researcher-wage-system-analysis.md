# Researcher Session: Wage/Cost System Analysis

**Date**: 2025-12-31 01:26:25
**Agent**: researcher
**Task**: Research current wage/cost calculation system to understand how to add configurable hourly wages and a designer/UX/UI role

## Work Completed

Conducted comprehensive exploration of the codebase to understand the current cost savings calculation system, including:
- Wage rate definitions and usage
- Cost calculation logic
- Role tracking (PM, Developer, Designer)
- Component architecture
- Data flow patterns

## Current Implementation

### 1. Wage Rate Constants (Hardcoded)

**Location**: `/workspace/frontend/src/hooks/useCostSavings.ts` (lines 5-6)

```typescript
const PM_HOURLY_RATE = 80;        // $80/hr for PM consulting
const DEV_HOURLY_RATE = 112.50;   // $112.50/hr for development
```

**Duplicated in**: `/workspace/frontend/src/components/savings/CostSavingsCard.tsx` (lines 18-19)

These rates are **hardcoded constants** - not currently configurable via settings.

### 2. Current Roles Tracked

The system currently recognizes **three agent types**:

**Location**: `/workspace/backend/internal/model/prd.go` (line 111)
```go
AgentProductManager AgentType = "product_manager"
AgentDesigner       AgentType = "designer"
AgentDeveloper      AgentType = "developer"
```

**Frontend Types**: `/workspace/frontend/src/types/index.ts` (line 10)
```typescript
export type AgentType = 'product_manager' | 'product' | 'designer' | 'developer';
```

**Note**: The `designer` agent type **already exists** in the system! It represents the "Bloom" persona.

### 3. Current Wage Calculations

**Only PM and Developer rates are used in cost calculations:**

#### PM Value Calculation
- **Input**: Minutes of PM-equivalent work
- **Formula**: `(minutes / 60) * PM_HOURLY_RATE`
- **Location**: `/workspace/frontend/src/hooks/useCostSavings.ts` (lines 68-70)

#### Developer Value Calculation
- **Input**: Hours of dev-equivalent work
- **Formula**: `hours * DEV_HOURLY_RATE`
- **Location**: `/workspace/frontend/src/hooks/useCostSavings.ts` (lines 75-77)

#### Total Value
```typescript
const totalValue = pmValue + devValue;
```

**Designer work is NOT currently included in cost savings calculations!**

### 4. Work Estimation Logic

**PM Minutes Estimation** (line 44-46):
```typescript
function estimatePmMinutes(messageCount: number): number {
  return messageCount * PM_MINUTES_PER_MESSAGE; // 1.5 min per message
}
```

**Dev Hours Estimation** (lines 52-56):
```typescript
function estimateDevHours(filesGenerated: number, messageCount: number): number {
  const fileBasedHours = filesGenerated * DEV_HOURS_PER_FILE; // 0.5 hrs per file
  const messageBasedHours = (messageCount / 100) * DEV_HOURS_PER_100_MESSAGES; // 0.25 hrs per 100 msgs
  return fileBasedHours + messageBasedHours;
}
```

**No designer work estimation currently exists.**

### 5. Data Flow Architecture

```
SessionMetrics {messageCount, filesGenerated, tokensUsed}
    ↓
useCostSavings() hook
    ↓ calculates
CostSavingsResult {data, pmValue, devValue, totalValue, aiCost, savingsMultiplier}
    ↓ passed to
CostSavingsCard / CostSavingsIcon components
    ↓ displays
Cost savings UI with breakdown
```

### 6. Components Using Cost Data

1. **CostSavingsCard** (`/workspace/frontend/src/components/savings/CostSavingsCard.tsx`)
   - Shows PM time + monetary value
   - Shows Dev time + monetary value
   - Shows total value delivered
   - Shows AI cost comparison
   - Displays disclaimer with rates

2. **CostSavingsIcon** (`/workspace/frontend/src/components/savings/CostSavingsIcon.tsx`)
   - Compact dollar icon with badge
   - Popover showing CostSavingsCard
   - Tracks "new savings" notifications via localStorage

3. **DiscoverySummaryModal** (`/workspace/frontend/src/components/discovery/DiscoverySummaryModal.tsx`)
   - Shows cost savings for discovery phase
   - Currently only uses messageCount (no files generated in discovery)

## What Needs to Change

### 1. Add Settings/Configuration System

**Current State**: No settings page or configuration system exists.

**Evidence**:
- Search for settings/config files found only node_modules
- No settings page in `/workspace/frontend/src/app`
- Only pages are: `page.tsx`, `projects/[id]/page.tsx`, `demo/discovery/page.tsx`

**Needed**:
- Create settings storage (localStorage or backend API)
- Create settings UI (modal or dedicated page)
- Create settings context/hook for global access

### 2. Make Wage Rates Configurable

**Files to Update**:

1. **Create settings storage**
   - Define settings interface with wage rates
   - Implement localStorage or API persistence
   - Create useSettings hook

2. **Update useCostSavings.ts**
   - Remove hardcoded constants
   - Accept rates as parameters or get from settings hook
   - Add DESIGNER_HOURLY_RATE (suggested default: $95/hr)

3. **Update CostSavingsCard.tsx**
   - Remove duplicate rate constants
   - Get rates from settings/hook
   - Update disclaimer to show all three rates

### 3. Add Designer Role to Cost Calculations

**Currently Missing**:
- No designer work estimation
- No designer hourly rate
- No designer value in totals

**What to Add**:

1. **Designer Rate Constant**
   ```typescript
   const DESIGNER_HOURLY_RATE = 95; // Suggested midpoint between PM ($80) and Dev ($112.50)
   ```

2. **Designer Work Estimation**
   - Need estimation logic (e.g., messages with "design" keywords, UI-related files)
   - Could track messages by agentType === 'designer'

3. **Update CostSavingsData interface** (in CostSavingsCard.tsx, lines 4-10)
   ```typescript
   export interface CostSavingsData {
     pmMinutes: number;
     devHours: number;
     designerHours: number; // NEW
     messageCount: number;
     filesGenerated: number;
     tokensUsed: number;
   }
   ```

4. **Update calculation flow**
   ```typescript
   const totalValue = pmValue + devValue + designerValue; // Add designer
   ```

5. **Update UI components**
   - Add third "ValueCard" for designer in CostSavingsCard
   - Update summary to show "~X hours of consulting" from all three roles

### 4. Track Designer Work Metrics

**Backend Already Tracks Agent Type**:
- Messages table has `agent_type` column (migration 006)
- Can be "product_manager", "designer", or "developer"

**Frontend Can Query**:
```typescript
// Count messages by agent type
const designerMessageCount = messages.filter(m => m.agentType === 'designer').length;
```

**Estimation Approach**:
```typescript
function estimateDesignerHours(messageCount: number, designerMessageCount: number): number {
  // Similar to dev hours, but for design work
  const DESIGNER_HOURS_PER_MESSAGE = 0.3; // ~18 min per design message
  return designerMessageCount * DESIGNER_HOURS_PER_MESSAGE;
}
```

## Implementation Recommendations

### Phase 1: Add Settings Infrastructure (Developer)
1. Create `/workspace/frontend/src/hooks/useWageSettings.ts`
   - Define WageSettings interface
   - Implement localStorage persistence
   - Provide default rates
   - Export hook for getting/setting rates

2. Create `/workspace/frontend/src/components/settings/WageSettingsModal.tsx`
   - Form to adjust PM, Designer, Developer hourly rates
   - Number inputs with validation
   - Save/Reset buttons
   - Show example calculations

3. Add settings icon to ChatContainer header (near CostSavingsIcon)

### Phase 2: Integrate Designer Role (Developer)
1. Update `useCostSavings` hook
   - Accept wage rates from settings
   - Add designer estimation logic
   - Track messages by agent type
   - Calculate designer value

2. Update `CostSavingsData` interface
   - Add designerHours field

3. Update `CostSavingsCard` component
   - Add third value card for designer
   - Use configurable rates
   - Update disclaimer

### Phase 3: Enhanced Tracking (Developer, Optional)
1. Pass actual message breakdown to useCostSavings
   - Instead of just messageCount, pass counts by agent type
   - More accurate work estimation

2. Add per-agent time tracking
   - Track which agent spent how much time
   - Could show breakdown: "2 hours PM planning, 1 hour design, 4 hours dev"

## Key Files Reference

### Frontend Files
- `/workspace/frontend/src/hooks/useCostSavings.ts` - Main calculation logic
- `/workspace/frontend/src/components/savings/CostSavingsCard.tsx` - Display component
- `/workspace/frontend/src/components/savings/CostSavingsIcon.tsx` - Icon with popover
- `/workspace/frontend/src/components/chat/ChatContainer.tsx` - Where icon appears
- `/workspace/frontend/src/types/index.ts` - AgentType definition

### Backend Files
- `/workspace/backend/internal/model/prd.go` - AgentType enum
- `/workspace/backend/internal/service/agent_context.go` - Agent selection logic
- `/workspace/backend/migrations/006_message_agent.sql` - Messages agent_type column

## Current Persona Names

- **Root** (product_manager) - Discovery guide
- **Bloom** (designer) - Design/UX persona (female pronouns)
- **Harvest** (developer) - Development persona (male pronouns)

## Technical Notes

1. **Duplicate Constants**: PM_HOURLY_RATE and DEV_HOURLY_RATE are defined in BOTH:
   - `/workspace/frontend/src/hooks/useCostSavings.ts`
   - `/workspace/frontend/src/components/savings/CostSavingsCard.tsx`

   After adding settings, these duplicates should be eliminated.

2. **No Backend Wage Storage**: Wages are purely frontend constants. Settings could be:
   - **Client-side only** (localStorage) - simpler, no backend changes
   - **Server-persisted** - would require backend settings table and API

3. **Designer Agent Exists But Not Counted**: The system already has full designer agent support in backend, frontend types, and message tracking. Only the cost calculation is missing!

4. **AI Cost Calculation**: Uses token-based estimation ($0.003 per 1K input, $0.015 per 1K output). This doesn't need to change for wage settings.

## Suggested Default Rates

Based on industry midpoints:
- **PM**: $80/hr (current, keep)
- **Designer**: $95/hr (suggested new, between PM and Dev)
- **Developer**: $112.50/hr (current, keep)

This creates a logical progression: PM < Designer < Developer

## Next Steps

The main agent should delegate to:
- **developer**: Implement settings infrastructure and integrate designer role into calculations
- **ux-tactical**: Design the wage settings UI/UX
