# Session 017 - Demo Fixes Branch

**Date**: 2025-12-31
**Branch**: `demo-fixes`
**Commits**: 6 new commits

## Summary

Major session focused on demo polish and UX improvements. Created new branch `demo-fixes` with comprehensive fixes for the chat interface, cost savings display, and phase categorization system.

## Commits Made

1. **a1e5f1f** - Phase grouping, segmented control, preview modal, wage settings
2. **0cc9b06** - Preview localStorage, icon-only toggle, accurate costs
3. **f25103b** - Hide discovery metadata, badge overflow, chat input padding
4. **c35f4b2** - Phase categorization with Root markers, settings position, cost fix
5. **646a144** - Root layout padding for top/bottom spacing

## Features Implemented

### Phase Categorization System
- Hide phase toggle until discovery complete + 10 messages minimum
- Sticky phases: messages inherit phase from previous message
- Root can mark transitions with `[Beginning X phase]` markers
- Minimum 2 messages per section to avoid orphan sections

### Preview System
- Fullscreen modal with device frame selector (Desktop/Tablet/Mobile)
- Fixed localStorage access by adding `allow-same-origin` to sandbox
- Apps using localStorage now render correctly

### Cost Savings
- Wage settings modal with configurable PM/Dev/Designer rates ($80, $112.50, $95/hr defaults)
- Per-agent time tracking (Root=PM, Bloom=Designer, Harvest=Developer)
- Fixed "<1 min" display for small time values
- Settings gear moved to right of "Connected" status

### UI Polish
- Icon-only toggle (Clock/Layers) with tooltips
- Hidden DISCOVERY_DATA during streaming
- Root layout padding for badge/input spacing (pt-2 pb-2)

## Tests

- 210 tests passing
- New test files: BuildPhaseProgress, CostSavingsCard, useCostSavings, useWageSettings, WageSettingsModal

## Known Issues

1. **Project title/summary extraction** - Some new projects not getting metadata extracted. Likely Claude not including DISCOVERY_DATA consistently.

## Key Files Modified

- `frontend/src/components/chat/BuildPhaseProgress.tsx` - Phase categorization logic
- `frontend/src/components/chat/ChatContainer.tsx` - Header layout, settings position
- `frontend/src/components/preview/PreviewModal.tsx` - Fullscreen preview
- `frontend/src/hooks/useWageSettings.ts` - Wage configuration
- `frontend/src/hooks/useCostSavings.ts` - Per-agent cost calculations
- `frontend/src/app/layout.tsx` - Root padding
