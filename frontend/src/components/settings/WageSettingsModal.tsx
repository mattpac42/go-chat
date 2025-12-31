'use client';

import { useState, useEffect, useCallback } from 'react';
import { useWageSettings, DEFAULT_WAGE_SETTINGS, type WageSettings } from '@/hooks/useWageSettings';

interface WageSettingsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

/**
 * Close icon for the modal
 */
function CloseIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M6 18L18 6M6 6l12 12"
      />
    </svg>
  );
}

interface RateInputProps {
  id: string;
  label: string;
  value: number;
  onChange: (value: number) => void;
}

function RateInput({ id, label, value, onChange }: RateInputProps) {
  return (
    <div className="flex flex-col gap-1.5">
      <label htmlFor={id} className="text-sm font-medium text-gray-700">
        {label}
      </label>
      <div className="relative">
        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500 text-sm">
          $
        </span>
        <input
          id={id}
          type="number"
          min="0"
          step="0.01"
          value={value}
          onChange={(e) => onChange(parseFloat(e.target.value) || 0)}
          className="w-full pl-7 pr-12 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-teal-500 focus:border-teal-500 text-gray-900"
        />
        <span className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 text-sm">
          /hr
        </span>
      </div>
    </div>
  );
}

/**
 * WageSettingsModal - Modal for configuring hourly rates for cost calculations
 *
 * Features:
 * - Three input fields for PM, Developer, and Designer rates
 * - Persists to localStorage via useWageSettings hook
 * - Save, Cancel, and Reset to Defaults options
 * - Closes on Escape key or backdrop click
 */
export function WageSettingsModal({ isOpen, onClose }: WageSettingsModalProps) {
  const { settings, updateSettings, resetSettings } = useWageSettings();

  // Local state for form values (allows Cancel to discard changes)
  const [formValues, setFormValues] = useState<WageSettings>(settings);

  // Sync form values when settings change or modal opens
  useEffect(() => {
    if (isOpen) {
      setFormValues(settings);
    }
  }, [isOpen, settings]);

  // Handle escape key
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  const handleSave = useCallback(() => {
    updateSettings(formValues);
    onClose();
  }, [formValues, updateSettings, onClose]);

  const handleReset = useCallback(() => {
    setFormValues(DEFAULT_WAGE_SETTINGS);
    resetSettings();
  }, [resetSettings]);

  const handleBackdropClick = useCallback((e: React.MouseEvent) => {
    // Only close if clicking the backdrop itself
    if (e.target === e.currentTarget) {
      onClose();
    }
  }, [onClose]);

  if (!isOpen) return null;

  return (
    <div
      data-testid="modal-backdrop"
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      onClick={handleBackdropClick}
    >
      <div
        role="dialog"
        aria-label="Wage Settings"
        className="bg-white rounded-xl shadow-xl w-full max-w-md mx-4 overflow-hidden"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200">
          <h2 className="text-lg font-semibold text-gray-900">Wage Settings</h2>
          <button
            onClick={onClose}
            className="p-1 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors"
            aria-label="Close"
          >
            <CloseIcon className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="px-6 py-5 space-y-4">
          <p className="text-sm text-gray-600 mb-4">
            Configure hourly rates used to calculate the value of AI-assisted work.
          </p>

          <RateInput
            id="pm-rate"
            label="Product Manager"
            value={formValues.pmHourlyRate}
            onChange={(value) => setFormValues((prev) => ({ ...prev, pmHourlyRate: value }))}
          />

          <RateInput
            id="dev-rate"
            label="Developer"
            value={formValues.devHourlyRate}
            onChange={(value) => setFormValues((prev) => ({ ...prev, devHourlyRate: value }))}
          />

          <RateInput
            id="designer-rate"
            label="Designer/UX"
            value={formValues.designerHourlyRate}
            onChange={(value) => setFormValues((prev) => ({ ...prev, designerHourlyRate: value }))}
          />
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between px-6 py-4 bg-gray-50 border-t border-gray-200">
          <button
            onClick={handleReset}
            className="text-sm text-gray-500 hover:text-gray-700 transition-colors"
          >
            Reset to Defaults
          </button>
          <div className="flex gap-3">
            <button
              onClick={onClose}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleSave}
              className="px-4 py-2 text-sm font-medium text-white bg-teal-600 rounded-lg hover:bg-teal-700 transition-colors"
            >
              Save
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
