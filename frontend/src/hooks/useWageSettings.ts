import { useState, useCallback, useEffect } from 'react';

/**
 * Wage settings for cost calculations
 */
export interface WageSettings {
  pmHourlyRate: number;       // default: 80
  devHourlyRate: number;      // default: 112.50
  designerHourlyRate: number; // default: 95
}

export const WAGE_SETTINGS_STORAGE_KEY = 'wage-settings';

export const DEFAULT_WAGE_SETTINGS: WageSettings = {
  pmHourlyRate: 80,
  devHourlyRate: 112.50,
  designerHourlyRate: 95,
};

/**
 * Parse and validate wage settings from stored JSON
 */
function parseStoredSettings(json: string | null): WageSettings | null {
  if (!json) return null;

  try {
    const parsed = JSON.parse(json);
    // Ensure all values are numbers
    return {
      pmHourlyRate: Number(parsed.pmHourlyRate ?? DEFAULT_WAGE_SETTINGS.pmHourlyRate),
      devHourlyRate: Number(parsed.devHourlyRate ?? DEFAULT_WAGE_SETTINGS.devHourlyRate),
      designerHourlyRate: Number(parsed.designerHourlyRate ?? DEFAULT_WAGE_SETTINGS.designerHourlyRate),
    };
  } catch {
    return null;
  }
}

/**
 * Load settings from localStorage
 */
function loadSettings(): WageSettings {
  if (typeof window === 'undefined') {
    return DEFAULT_WAGE_SETTINGS;
  }

  const stored = localStorage.getItem(WAGE_SETTINGS_STORAGE_KEY);
  return parseStoredSettings(stored) ?? DEFAULT_WAGE_SETTINGS;
}

/**
 * Save settings to localStorage
 */
function saveSettings(settings: WageSettings): void {
  if (typeof window === 'undefined') return;
  localStorage.setItem(WAGE_SETTINGS_STORAGE_KEY, JSON.stringify(settings));
}

/**
 * Clear settings from localStorage
 */
function clearSettings(): void {
  if (typeof window === 'undefined') return;
  localStorage.removeItem(WAGE_SETTINGS_STORAGE_KEY);
}

export interface UseWageSettingsResult {
  settings: WageSettings;
  updateSettings: (updates: Partial<WageSettings>) => void;
  resetSettings: () => void;
}

/**
 * Hook to manage wage settings with localStorage persistence
 *
 * @example
 * ```tsx
 * const { settings, updateSettings, resetSettings } = useWageSettings();
 *
 * // Access rates
 * console.log(settings.pmHourlyRate); // 80 by default
 *
 * // Update a single rate
 * updateSettings({ pmHourlyRate: 100 });
 *
 * // Reset to defaults
 * resetSettings();
 * ```
 */
export function useWageSettings(): UseWageSettingsResult {
  const [settings, setSettings] = useState<WageSettings>(loadSettings);

  // Sync with localStorage on mount (for SSR hydration)
  useEffect(() => {
    setSettings(loadSettings());
  }, []);

  const updateSettings = useCallback((updates: Partial<WageSettings>) => {
    setSettings((current) => {
      const newSettings = { ...current, ...updates };
      saveSettings(newSettings);
      return newSettings;
    });
  }, []);

  const resetSettings = useCallback(() => {
    clearSettings();
    setSettings(DEFAULT_WAGE_SETTINGS);
  }, []);

  return {
    settings,
    updateSettings,
    resetSettings,
  };
}
