import { render, screen } from '@testing-library/react';
import { CostSavingsCard, CostSavingsData } from '@/components/savings/CostSavingsCard';

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

describe('CostSavingsCard', () => {
  describe('time and value consistency', () => {
    it('shows 0 min and $0 value when dev hours is exactly 0', () => {
      const data: CostSavingsData = {
        pmMinutes: 0,
        devHours: 0,
        designerHours: 0,
        messageCount: 0,
        filesGenerated: 0,
        tokensUsed: 0,
      };

      render(<CostSavingsCard data={data} />);

      // Dev time should show 0 min with $0 value - consistent
      // The label is "Dev Time" (title case)
      const devSection = screen.getByText('Dev Time').closest('div')?.parentElement;
      expect(devSection).toHaveTextContent('0 min');
      expect(devSection).toHaveTextContent('$0.00 value');
    });

    it('shows non-zero time when dev hours is very small but non-zero', () => {
      const data: CostSavingsData = {
        pmMinutes: 0,
        devHours: 0.0025, // Very small: 0.15 minutes = 9 seconds
        designerHours: 0,
        messageCount: 1,
        filesGenerated: 0,
        tokensUsed: 500,
      };

      render(<CostSavingsCard data={data} />);

      // When there's a non-zero time but rounds to 0 minutes,
      // should show "<1 min" not "0 min" to be consistent with non-zero value
      // Currently shows "0 min" which is the bug we're fixing
      const devSection = screen.getByText('Dev Time').closest('div')?.parentElement;
      // This is the failing assertion - we want "<1 min" not "0 min"
      expect(devSection).toHaveTextContent('<1 min');
      // 0.0025 * 112.5 = 0.28125 -> $0.28
      expect(devSection).toHaveTextContent('$0.28 value');
    });

    it('shows <1 min for very small PM time values', () => {
      const data: CostSavingsData = {
        pmMinutes: 0.3, // Less than half a minute - would round to 0
        devHours: 0,
        designerHours: 0,
        messageCount: 1,
        filesGenerated: 0,
        tokensUsed: 500,
      };

      render(<CostSavingsCard data={data} />);

      // When there's non-zero time that rounds to 0, show "<1 min"
      const pmSection = screen.getByText('PM Time').closest('div')?.parentElement;
      expect(pmSection).toHaveTextContent('<1 min');
      // 0.3 / 60 * 80 = 0.40 -> $0.40
      expect(pmSection).toHaveTextContent('$0.40 value');
    });
  });

  describe('total calculation', () => {
    it('displays total that equals sum of PM + Dev + Designer values', () => {
      const data: CostSavingsData = {
        pmMinutes: 15, // 0.25 hours * $80 = $20
        devHours: 1, // 1 hour * $112.50 = $112.50
        designerHours: 2, // 2 hours * $95 = $190
        messageCount: 25,
        filesGenerated: 2,
        tokensUsed: 12500,
      };

      render(<CostSavingsCard data={data} />);

      // Total: $20 + $112.50 + $190 = $322.50
      // formatCurrency rounds to whole dollar for values >= $1
      expect(screen.getByText('$323')).toBeInTheDocument();
    });

    it('correctly sums small values in total', () => {
      const data: CostSavingsData = {
        pmMinutes: 1.5, // 0.025 hours * $80 = $2
        devHours: 0.0025, // 0.0025 hours * $112.50 = $0.28
        designerHours: 0,
        messageCount: 1,
        filesGenerated: 0,
        tokensUsed: 500,
      };

      render(<CostSavingsCard data={data} />);

      // Total: $2 + $0.28 = $2.28, should round to $2
      // The important thing is the total is displayed
      expect(screen.getByText('$2')).toBeInTheDocument();
    });
  });

  describe('formatting', () => {
    it('formats hours correctly', () => {
      const data: CostSavingsData = {
        pmMinutes: 0,
        devHours: 2.5,
        designerHours: 0,
        messageCount: 10,
        filesGenerated: 5,
        tokensUsed: 5000,
      };

      render(<CostSavingsCard data={data} />);

      const devSection = screen.getByText('Dev Time').closest('div')?.parentElement;
      expect(devSection).toHaveTextContent('2.5 hours');
    });

    it('converts hours to minutes for values less than 1 hour', () => {
      const data: CostSavingsData = {
        pmMinutes: 0,
        devHours: 0.5, // 30 minutes
        designerHours: 0,
        messageCount: 10,
        filesGenerated: 1,
        tokensUsed: 5000,
      };

      render(<CostSavingsCard data={data} />);

      const devSection = screen.getByText('Dev Time').closest('div')?.parentElement;
      expect(devSection).toHaveTextContent('30 min');
    });
  });
});
