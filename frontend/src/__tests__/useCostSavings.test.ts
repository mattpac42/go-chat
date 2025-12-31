import { renderHook } from '@testing-library/react';
import {
  useCostSavings,
  estimatePmMinutes,
  estimateDevHours,
  estimateDesignerHours,
  SessionMetrics,
} from '@/hooks/useCostSavings';

// Mock the useWageSettings hook
jest.mock('@/hooks/useWageSettings', () => ({
  useWageSettings: () => ({
    settings: {
      pmHourlyRate: 80,
      devHourlyRate: 112.5,
      designerHourlyRate: 95,
    },
  }),
  DEFAULT_WAGE_SETTINGS: {
    pmHourlyRate: 80,
    devHourlyRate: 112.5,
    designerHourlyRate: 95,
  },
}));

describe('useCostSavings', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  describe('estimatePmMinutes', () => {
    it('uses pmMessageCount when provided', () => {
      // 5 PM messages * 1.5 min = 7.5 minutes
      const result = estimatePmMinutes(5, 100);
      expect(result).toBe(7.5);
    });

    it('falls back to messageCount when pmMessageCount is undefined', () => {
      // 100 messages * 1.5 min = 150 minutes
      const result = estimatePmMinutes(undefined, 100);
      expect(result).toBe(150);
    });

    it('uses zero when pmMessageCount is 0', () => {
      const result = estimatePmMinutes(0, 100);
      expect(result).toBe(0);
    });
  });

  describe('estimateDevHours', () => {
    it('uses developerMessageCount when provided', () => {
      // 5 files * 0.5 hours = 2.5 hours from files
      // 20 dev messages / 100 * 0.25 = 0.05 hours from messages
      // Total = 2.55 hours
      const result = estimateDevHours(5, 20, 100);
      expect(result).toBe(2.55);
    });

    it('falls back to messageCount when developerMessageCount is undefined', () => {
      // 5 files * 0.5 hours = 2.5 hours from files
      // 100 messages / 100 * 0.25 = 0.25 hours from messages
      // Total = 2.75 hours
      const result = estimateDevHours(5, undefined, 100);
      expect(result).toBe(2.75);
    });

    it('uses zero messages when developerMessageCount is 0', () => {
      // 5 files * 0.5 hours = 2.5 hours from files
      // 0 dev messages / 100 * 0.25 = 0 hours from messages
      // Total = 2.5 hours
      const result = estimateDevHours(5, 0, 100);
      expect(result).toBe(2.5);
    });

    it('calculates correctly with no files', () => {
      // 0 files * 0.5 hours = 0 hours from files
      // 50 dev messages / 100 * 0.25 = 0.125 hours from messages
      const result = estimateDevHours(0, 50, 100);
      expect(result).toBe(0.125);
    });
  });

  describe('estimateDesignerHours', () => {
    it('calculates designer hours from message count', () => {
      // 4 designer messages * 0.5 hours = 2 hours
      const result = estimateDesignerHours(4);
      expect(result).toBe(2);
    });

    it('returns 0 for no designer messages', () => {
      const result = estimateDesignerHours(0);
      expect(result).toBe(0);
    });
  });

  describe('useCostSavings hook', () => {
    it('calculates PM value using pmMessageCount when provided', () => {
      const metrics: SessionMetrics = {
        messageCount: 100,
        filesGenerated: 0,
        pmMessageCount: 10, // Only 10 PM messages, not 100
        developerMessageCount: 5,
        designerMessageCount: 3,
      };

      const { result } = renderHook(() => useCostSavings(metrics));

      // PM: 10 messages * 1.5 min = 15 min = 0.25 hours * $80/hr = $20
      expect(result.current.data.pmMinutes).toBe(15);
      expect(result.current.pmValue).toBe(20);
    });

    it('calculates developer value using developerMessageCount when provided', () => {
      const metrics: SessionMetrics = {
        messageCount: 100,
        filesGenerated: 2,
        pmMessageCount: 10,
        developerMessageCount: 20, // Only 20 developer messages
        designerMessageCount: 0,
      };

      const { result } = renderHook(() => useCostSavings(metrics));

      // Dev: 2 files * 0.5 hrs = 1 hr + 20/100 * 0.25 = 0.05 hrs = 1.05 hrs
      // 1.05 hrs * $112.50/hr = $118.125
      expect(result.current.data.devHours).toBe(1.05);
      expect(result.current.devValue).toBe(118.125);
    });

    it('calculates designer value using designerMessageCount', () => {
      const metrics: SessionMetrics = {
        messageCount: 100,
        filesGenerated: 0,
        pmMessageCount: 10,
        developerMessageCount: 5,
        designerMessageCount: 4, // 4 designer messages
      };

      const { result } = renderHook(() => useCostSavings(metrics));

      // Designer: 4 messages * 0.5 hrs = 2 hrs * $95/hr = $190
      expect(result.current.data.designerHours).toBe(2);
      expect(result.current.designerValue).toBe(190);
    });

    it('falls back to messageCount when agent-specific counts are not provided', () => {
      const metrics: SessionMetrics = {
        messageCount: 50,
        filesGenerated: 2,
        // No agent-specific counts provided
      };

      const { result } = renderHook(() => useCostSavings(metrics));

      // PM: 50 messages * 1.5 min = 75 min
      expect(result.current.data.pmMinutes).toBe(75);

      // Dev: 2 files * 0.5 hrs = 1 hr + 50/100 * 0.25 = 0.125 hrs = 1.125 hrs
      expect(result.current.data.devHours).toBe(1.125);

      // Designer: 0 (default)
      expect(result.current.data.designerHours).toBe(0);
    });

    it('calculates total value correctly with per-agent counts', () => {
      const metrics: SessionMetrics = {
        messageCount: 100,
        filesGenerated: 2,
        pmMessageCount: 10,
        developerMessageCount: 20,
        designerMessageCount: 4,
      };

      const { result } = renderHook(() => useCostSavings(metrics));

      // PM: 10 * 1.5 = 15 min = 0.25 hrs * $80 = $20
      // Dev: 2 * 0.5 + 20/100 * 0.25 = 1.05 hrs * $112.50 = $118.125
      // Designer: 4 * 0.5 = 2 hrs * $95 = $190
      // Total: $20 + $118.125 + $190 = $328.125
      expect(result.current.totalValue).toBe(328.125);
    });

    it('handles zero counts gracefully', () => {
      const metrics: SessionMetrics = {
        messageCount: 0,
        filesGenerated: 0,
        pmMessageCount: 0,
        developerMessageCount: 0,
        designerMessageCount: 0,
      };

      const { result } = renderHook(() => useCostSavings(metrics));

      expect(result.current.data.pmMinutes).toBe(0);
      expect(result.current.data.devHours).toBe(0);
      expect(result.current.data.designerHours).toBe(0);
      expect(result.current.totalValue).toBe(0);
    });

    it('more accurately reflects value when using per-agent counts vs fallback', () => {
      // Without per-agent counts (old behavior)
      const metricsOld: SessionMetrics = {
        messageCount: 100,
        filesGenerated: 2,
      };

      // With per-agent counts (new behavior)
      const metricsNew: SessionMetrics = {
        messageCount: 100,
        filesGenerated: 2,
        pmMessageCount: 10, // Only 10 of 100 messages are PM
        developerMessageCount: 20, // Only 20 of 100 messages are developer
        designerMessageCount: 5, // Only 5 of 100 messages are designer
      };

      const { result: resultOld } = renderHook(() => useCostSavings(metricsOld));
      const { result: resultNew } = renderHook(() => useCostSavings(metricsNew));

      // Old behavior inflates PM value because it uses total message count
      // PM: 100 * 1.5 = 150 min = 2.5 hrs * $80 = $200
      expect(resultOld.current.data.pmMinutes).toBe(150);

      // New behavior uses only PM messages
      // PM: 10 * 1.5 = 15 min = 0.25 hrs * $80 = $20
      expect(resultNew.current.data.pmMinutes).toBe(15);

      // New calculation should be lower (more accurate) than old
      expect(resultNew.current.pmValue).toBeLessThan(resultOld.current.pmValue);
      expect(resultNew.current.devValue).toBeLessThan(resultOld.current.devValue);
    });
  });
});
