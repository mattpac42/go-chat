# Developer Session: Per-Agent Cost Savings Calculations

**Date**: 2025-12-31T01:59:48
**Agent**: developer
**Task**: Fix cost savings calculations to count time per agent type instead of using total message count

## Work Completed

Implemented accurate per-agent cost savings calculations by:

1. **Updated SessionMetrics interface** (`/workspace/frontend/src/hooks/useCostSavings.ts`)
   - Added `pmMessageCount?: number` for product_manager agent (Root)
   - Added `developerMessageCount?: number` for developer agent (Harvest)
   - Kept existing `designerMessageCount?: number` for designer agent (Bloom)

2. **Updated estimation functions**
   - `estimatePmMinutes(pmMessageCount, messageCount)` - Uses pmMessageCount when available, falls back to messageCount for backward compatibility
   - `estimateDevHours(filesGenerated, developerMessageCount, messageCount)` - Uses developerMessageCount when available, falls back to messageCount

3. **Updated useCostSavings hook**
   - Extracts pmMessageCount and developerMessageCount from metrics
   - Passes appropriate counts to estimation functions with fallback support

4. **Updated ChatContainer.tsx**
   - Now passes per-agent counts by filtering messages by agentType:
     - `pmMessageCount`: Messages with agentType 'product_manager' or 'product' (legacy)
     - `designerMessageCount`: Messages with agentType 'designer'
     - `developerMessageCount`: Messages with agentType 'developer'

5. **Created comprehensive test suite** (`/workspace/frontend/src/__tests__/useCostSavings.test.ts`)
   - 16 tests covering estimation functions and hook behavior
   - Tests verify fallback to messageCount when agent-specific counts not provided
   - Tests verify accurate calculations when per-agent counts are provided

## Decisions Made

- **Backward compatibility**: Used nullish coalescing (`??`) to fall back to total messageCount when agent-specific counts are not provided. This ensures existing usage continues to work.
- **Agent type mapping**: Included 'product' as legacy alias for 'product_manager' to handle both old and new agent type values.
- **Test isolation**: Mocked useWageSettings to isolate cost calculation tests from wage settings implementation.

## Files Modified

- `/workspace/frontend/src/hooks/useCostSavings.ts`: Added pmMessageCount and developerMessageCount to SessionMetrics, updated estimation functions with fallback logic
- `/workspace/frontend/src/components/chat/ChatContainer.tsx`: Updated CostSavingsIcon props to filter messages by agentType
- `/workspace/frontend/src/__tests__/useCostSavings.test.ts`: New test file with 16 tests

## Test Results

All 186 tests pass, including 16 new tests for useCostSavings.

## Recommendations

- Consider tracking actual files generated per developer session for more accurate developer value calculation
- The filesGenerated metric is currently hardcoded to 0 in ChatContainer - this should be updated when file tracking is implemented
