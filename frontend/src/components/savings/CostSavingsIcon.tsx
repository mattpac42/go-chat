'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { CostSavingsCard } from './CostSavingsCard';
import { useCostSavings, type SessionMetrics } from '@/hooks/useCostSavings';

interface CostSavingsIconProps {
  metrics: SessionMetrics;
  /** Storage key for tracking last viewed savings (unique per project) */
  storageKey?: string;
}

/**
 * Dollar sign icon component
 */
function DollarIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      viewBox="0 0 24 24"
      aria-hidden="true"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M12 6v12m-3-2.818.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  );
}

/**
 * Close icon for the popover
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

/**
 * Format currency for display in the badge
 */
function formatCompactCurrency(value: number): string {
  if (value < 1) {
    return '$0';
  }
  if (value >= 1000) {
    return `$${(value / 1000).toFixed(1)}k`;
  }
  return `$${Math.round(value)}`;
}

/**
 * CostSavingsIcon - A compact icon button that shows cost savings in a popover
 *
 * Features:
 * - Dollar icon with subtle badge showing total savings
 * - Click to reveal CostSavingsCard in a popover
 * - Pulse animation when savings have increased since last view
 * - Tracks last viewed timestamp in localStorage
 */
export function CostSavingsIcon({
  metrics,
  storageKey = 'cost-savings-last-viewed',
}: CostSavingsIconProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [hasNewSavings, setHasNewSavings] = useState(false);
  const popoverRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);

  // Calculate current savings
  const { data, totalValue } = useCostSavings(metrics);

  // Track last viewed savings value
  useEffect(() => {
    if (typeof window === 'undefined') return;

    const stored = localStorage.getItem(storageKey);
    if (stored) {
      const lastViewedValue = parseFloat(stored);
      // Mark as "new" if current value is significantly higher (at least $1 more)
      if (totalValue > lastViewedValue + 1) {
        setHasNewSavings(true);
      }
    } else if (totalValue > 1) {
      // First time seeing any savings
      setHasNewSavings(true);
    }
  }, [totalValue, storageKey]);

  // Mark as viewed when popover opens
  const handleOpen = useCallback(() => {
    setIsOpen(true);
    setHasNewSavings(false);
    if (typeof window !== 'undefined') {
      localStorage.setItem(storageKey, totalValue.toString());
    }
  }, [totalValue, storageKey]);

  const handleClose = useCallback(() => {
    setIsOpen(false);
  }, []);

  // Handle clicking outside to close
  useEffect(() => {
    if (!isOpen) return;

    const handleClickOutside = (event: MouseEvent) => {
      if (
        popoverRef.current &&
        !popoverRef.current.contains(event.target as Node) &&
        buttonRef.current &&
        !buttonRef.current.contains(event.target as Node)
      ) {
        handleClose();
      }
    };

    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        handleClose();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    document.addEventListener('keydown', handleEscape);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      document.removeEventListener('keydown', handleEscape);
    };
  }, [isOpen, handleClose]);

  // Don't show if no meaningful savings yet
  if (totalValue < 1) {
    return null;
  }

  return (
    <div className="relative">
      {/* Icon Button */}
      <button
        ref={buttonRef}
        onClick={isOpen ? handleClose : handleOpen}
        className={`
          relative p-2 rounded-lg transition-all duration-200
          ${isOpen
            ? 'bg-teal-100 text-teal-700'
            : 'text-teal-600 hover:bg-teal-50 hover:text-teal-700'
          }
          focus:outline-none focus:ring-2 focus:ring-teal-500 focus:ring-offset-1
          ${hasNewSavings ? 'animate-pulse-subtle' : ''}
        `}
        aria-label={`View cost savings: ${formatCompactCurrency(totalValue)} saved`}
        aria-expanded={isOpen}
        aria-haspopup="dialog"
      >
        <DollarIcon className="w-5 h-5" />

        {/* Savings badge */}
        <span className={`
          absolute -top-1 -right-1 min-w-[1.25rem] h-5 px-1
          text-xs font-semibold text-white bg-teal-500 rounded-full
          flex items-center justify-center
          ${hasNewSavings ? 'animate-bounce-subtle' : ''}
        `}>
          {formatCompactCurrency(totalValue)}
        </span>

        {/* Pulse ring for new savings */}
        {hasNewSavings && (
          <span className="absolute inset-0 rounded-lg animate-ping-slow bg-teal-400 opacity-30" />
        )}
      </button>

      {/* Popover */}
      {isOpen && (
        <div
          ref={popoverRef}
          role="dialog"
          aria-label="Cost savings details"
          className="absolute right-0 top-full mt-2 z-50 w-80 sm:w-96 animate-fade-in"
        >
          {/* Arrow */}
          <div className="absolute -top-2 right-4 w-4 h-4 bg-white rotate-45 border-l border-t border-gray-200" />

          {/* Content */}
          <div className="relative bg-white rounded-lg shadow-lg border border-gray-200 overflow-hidden">
            {/* Header with close button */}
            <div className="flex items-center justify-between px-4 py-2 bg-gray-50 border-b border-gray-200">
              <span className="text-sm font-medium text-gray-700">Your Savings</span>
              <button
                onClick={handleClose}
                className="p-1 text-gray-400 hover:text-gray-600 rounded hover:bg-gray-200 transition-colors"
                aria-label="Close savings"
              >
                <CloseIcon className="w-4 h-4" />
              </button>
            </div>

            {/* Cost Savings Card */}
            <div className="p-0">
              <CostSavingsCard data={data} showDetailed />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
