import { renderHook, act } from '@testing-library/react';
import { useWageSettings, DEFAULT_WAGE_SETTINGS, WAGE_SETTINGS_STORAGE_KEY } from '@/hooks/useWageSettings';

describe('useWageSettings', () => {
  beforeEach(() => {
    // Clear localStorage before each test
    localStorage.clear();
  });

  afterEach(() => {
    localStorage.clear();
  });

  describe('default values', () => {
    it('returns default wage settings when no stored settings exist', () => {
      const { result } = renderHook(() => useWageSettings());

      expect(result.current.settings).toEqual({
        pmHourlyRate: 80,
        devHourlyRate: 112.50,
        designerHourlyRate: 95,
      });
    });

    it('exports correct default values', () => {
      expect(DEFAULT_WAGE_SETTINGS).toEqual({
        pmHourlyRate: 80,
        devHourlyRate: 112.50,
        designerHourlyRate: 95,
      });
    });
  });

  describe('localStorage persistence', () => {
    it('loads settings from localStorage on mount', () => {
      const customSettings = {
        pmHourlyRate: 100,
        devHourlyRate: 150,
        designerHourlyRate: 120,
      };
      localStorage.setItem(WAGE_SETTINGS_STORAGE_KEY, JSON.stringify(customSettings));

      const { result } = renderHook(() => useWageSettings());

      expect(result.current.settings).toEqual(customSettings);
    });

    it('saves settings to localStorage when updated', () => {
      const { result } = renderHook(() => useWageSettings());

      act(() => {
        result.current.updateSettings({
          pmHourlyRate: 90,
          devHourlyRate: 120,
          designerHourlyRate: 100,
        });
      });

      const stored = JSON.parse(localStorage.getItem(WAGE_SETTINGS_STORAGE_KEY) || '{}');
      expect(stored).toEqual({
        pmHourlyRate: 90,
        devHourlyRate: 120,
        designerHourlyRate: 100,
      });
    });

    it('handles partial updates correctly', () => {
      const { result } = renderHook(() => useWageSettings());

      act(() => {
        result.current.updateSettings({ pmHourlyRate: 90 });
      });

      expect(result.current.settings).toEqual({
        pmHourlyRate: 90,
        devHourlyRate: 112.50,
        designerHourlyRate: 95,
      });
    });

    it('handles invalid JSON in localStorage gracefully', () => {
      localStorage.setItem(WAGE_SETTINGS_STORAGE_KEY, 'invalid-json');

      const { result } = renderHook(() => useWageSettings());

      expect(result.current.settings).toEqual(DEFAULT_WAGE_SETTINGS);
    });
  });

  describe('updateSettings', () => {
    it('updates the settings state', () => {
      const { result } = renderHook(() => useWageSettings());

      act(() => {
        result.current.updateSettings({
          pmHourlyRate: 100,
          devHourlyRate: 150,
          designerHourlyRate: 120,
        });
      });

      expect(result.current.settings).toEqual({
        pmHourlyRate: 100,
        devHourlyRate: 150,
        designerHourlyRate: 120,
      });
    });

    it('preserves unchanged values during partial update', () => {
      const { result } = renderHook(() => useWageSettings());

      act(() => {
        result.current.updateSettings({ devHourlyRate: 200 });
      });

      expect(result.current.settings.pmHourlyRate).toBe(80);
      expect(result.current.settings.devHourlyRate).toBe(200);
      expect(result.current.settings.designerHourlyRate).toBe(95);
    });
  });

  describe('resetSettings', () => {
    it('resets settings to defaults', () => {
      const { result } = renderHook(() => useWageSettings());

      // First update to custom values
      act(() => {
        result.current.updateSettings({
          pmHourlyRate: 100,
          devHourlyRate: 200,
          designerHourlyRate: 150,
        });
      });

      // Then reset
      act(() => {
        result.current.resetSettings();
      });

      expect(result.current.settings).toEqual(DEFAULT_WAGE_SETTINGS);
    });

    it('clears localStorage on reset', () => {
      const { result } = renderHook(() => useWageSettings());

      act(() => {
        result.current.updateSettings({ pmHourlyRate: 100 });
      });

      expect(localStorage.getItem(WAGE_SETTINGS_STORAGE_KEY)).not.toBeNull();

      act(() => {
        result.current.resetSettings();
      });

      expect(localStorage.getItem(WAGE_SETTINGS_STORAGE_KEY)).toBeNull();
    });
  });

  describe('type safety', () => {
    it('handles string values converted to numbers', () => {
      // Simulate corrupted/type-unsafe localStorage data
      localStorage.setItem(
        WAGE_SETTINGS_STORAGE_KEY,
        JSON.stringify({
          pmHourlyRate: '100',
          devHourlyRate: '150',
          designerHourlyRate: '120',
        })
      );

      const { result } = renderHook(() => useWageSettings());

      // Should still work with number-like strings
      expect(typeof result.current.settings.pmHourlyRate).toBe('number');
    });
  });
});
