import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { WageSettingsModal } from '@/components/settings/WageSettingsModal';
import { DEFAULT_WAGE_SETTINGS, WAGE_SETTINGS_STORAGE_KEY } from '@/hooks/useWageSettings';

describe('WageSettingsModal', () => {
  const mockOnClose = jest.fn();

  beforeEach(() => {
    localStorage.clear();
    mockOnClose.mockClear();
  });

  afterEach(() => {
    localStorage.clear();
  });

  describe('rendering', () => {
    it('does not render when isOpen is false', () => {
      render(<WageSettingsModal isOpen={false} onClose={mockOnClose} />);
      expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });

    it('renders modal when isOpen is true', () => {
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });

    it('displays all three input fields with correct labels', () => {
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      expect(screen.getByLabelText('Product Manager')).toBeInTheDocument();
      expect(screen.getByLabelText('Developer')).toBeInTheDocument();
      expect(screen.getByLabelText('Designer/UX')).toBeInTheDocument();
    });

    it('displays default values in inputs', () => {
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      expect(screen.getByLabelText('Product Manager')).toHaveValue(80);
      expect(screen.getByLabelText('Developer')).toHaveValue(112.5);
      expect(screen.getByLabelText('Designer/UX')).toHaveValue(95);
    });

    it('displays stored values from localStorage', () => {
      const customSettings = {
        pmHourlyRate: 100,
        devHourlyRate: 150,
        designerHourlyRate: 120,
      };
      localStorage.setItem(WAGE_SETTINGS_STORAGE_KEY, JSON.stringify(customSettings));

      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      expect(screen.getByLabelText('Product Manager')).toHaveValue(100);
      expect(screen.getByLabelText('Developer')).toHaveValue(150);
      expect(screen.getByLabelText('Designer/UX')).toHaveValue(120);
    });

    it('displays Save and Cancel buttons', () => {
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument();
    });

    it('displays dollar sign prefix for inputs', () => {
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      // Each input group should have a $ prefix
      const dollarSigns = screen.getAllByText('$');
      expect(dollarSigns.length).toBe(3);
    });
  });

  describe('interactions', () => {
    it('updates input value when user types', async () => {
      const user = userEvent.setup();
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      const pmInput = screen.getByLabelText('Product Manager');
      await user.clear(pmInput);
      await user.type(pmInput, '100');

      expect(pmInput).toHaveValue(100);
    });

    it('calls onClose when Cancel button is clicked', async () => {
      const user = userEvent.setup();
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      await user.click(screen.getByRole('button', { name: /cancel/i }));

      expect(mockOnClose).toHaveBeenCalledTimes(1);
    });

    it('saves settings to localStorage when Save is clicked', async () => {
      const user = userEvent.setup();
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      const pmInput = screen.getByLabelText('Product Manager');
      await user.clear(pmInput);
      await user.type(pmInput, '100');

      await user.click(screen.getByRole('button', { name: /save/i }));

      const stored = JSON.parse(localStorage.getItem(WAGE_SETTINGS_STORAGE_KEY) || '{}');
      expect(stored.pmHourlyRate).toBe(100);
    });

    it('calls onClose after saving', async () => {
      const user = userEvent.setup();
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      await user.click(screen.getByRole('button', { name: /save/i }));

      expect(mockOnClose).toHaveBeenCalledTimes(1);
    });

    it('closes modal when clicking overlay/backdrop', async () => {
      const user = userEvent.setup();
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      const backdrop = screen.getByTestId('modal-backdrop');
      await user.click(backdrop);

      expect(mockOnClose).toHaveBeenCalledTimes(1);
    });

    it('closes modal when pressing Escape key', async () => {
      const user = userEvent.setup();
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      await user.keyboard('{Escape}');

      expect(mockOnClose).toHaveBeenCalledTimes(1);
    });

    it('does not close when clicking inside modal content', async () => {
      const user = userEvent.setup();
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      const modalContent = screen.getByRole('dialog');
      await user.click(modalContent);

      expect(mockOnClose).not.toHaveBeenCalled();
    });
  });

  describe('validation', () => {
    it('prevents negative values', async () => {
      const user = userEvent.setup();
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      const pmInput = screen.getByLabelText('Product Manager');
      await user.clear(pmInput);
      await user.type(pmInput, '-50');

      // Input should have min="0" attribute
      expect(pmInput).toHaveAttribute('min', '0');
    });

    it('saves valid values only when Save is clicked', async () => {
      const user = userEvent.setup();
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      const pmInput = screen.getByLabelText('Product Manager');
      const devInput = screen.getByLabelText('Developer');
      const designerInput = screen.getByLabelText('Designer/UX');

      await user.clear(pmInput);
      await user.type(pmInput, '90');
      await user.clear(devInput);
      await user.type(devInput, '140');
      await user.clear(designerInput);
      await user.type(designerInput, '110');

      await user.click(screen.getByRole('button', { name: /save/i }));

      const stored = JSON.parse(localStorage.getItem(WAGE_SETTINGS_STORAGE_KEY) || '{}');
      expect(stored).toEqual({
        pmHourlyRate: 90,
        devHourlyRate: 140,
        designerHourlyRate: 110,
      });
    });
  });

  describe('reset functionality', () => {
    it('displays Reset to Defaults button', () => {
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);
      expect(screen.getByRole('button', { name: /reset to defaults/i })).toBeInTheDocument();
    });

    it('resets all values to defaults when Reset is clicked', async () => {
      const user = userEvent.setup();

      // Start with custom values
      const customSettings = {
        pmHourlyRate: 100,
        devHourlyRate: 150,
        designerHourlyRate: 120,
      };
      localStorage.setItem(WAGE_SETTINGS_STORAGE_KEY, JSON.stringify(customSettings));

      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      await user.click(screen.getByRole('button', { name: /reset to defaults/i }));

      expect(screen.getByLabelText('Product Manager')).toHaveValue(DEFAULT_WAGE_SETTINGS.pmHourlyRate);
      expect(screen.getByLabelText('Developer')).toHaveValue(DEFAULT_WAGE_SETTINGS.devHourlyRate);
      expect(screen.getByLabelText('Designer/UX')).toHaveValue(DEFAULT_WAGE_SETTINGS.designerHourlyRate);
    });
  });

  describe('accessibility', () => {
    it('has proper aria-label for the modal', () => {
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);
      expect(screen.getByRole('dialog')).toHaveAttribute('aria-label', 'Wage Settings');
    });

    it('inputs have number type', () => {
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      expect(screen.getByLabelText('Product Manager')).toHaveAttribute('type', 'number');
      expect(screen.getByLabelText('Developer')).toHaveAttribute('type', 'number');
      expect(screen.getByLabelText('Designer/UX')).toHaveAttribute('type', 'number');
    });

    it('inputs have step attribute for decimal support', () => {
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      expect(screen.getByLabelText('Product Manager')).toHaveAttribute('step', '0.01');
      expect(screen.getByLabelText('Developer')).toHaveAttribute('step', '0.01');
      expect(screen.getByLabelText('Designer/UX')).toHaveAttribute('step', '0.01');
    });
  });

  describe('version display', () => {
    it('displays version information in the modal footer', () => {
      render(<WageSettingsModal isOpen={true} onClose={mockOnClose} />);

      expect(screen.getByTestId('version-display')).toBeInTheDocument();
    });
  });
});
