'use client';

import { useState, useEffect, useCallback } from 'react';
import { BuildPhase } from './BuildPhaseProgress';

interface MilestoneToastProps {
  phase: BuildPhase;
  onDismiss: () => void;
  autoHideDelay?: number;
}

const PHASE_CONFIG: Record<
  BuildPhase,
  { label: string; icon: React.ReactNode; bgColor: string; textColor: string }
> = {
  discovery: {
    label: 'Discovery Complete',
    icon: <CheckIcon className="w-4 h-4" />,
    bgColor: 'bg-teal-500',
    textColor: 'text-white',
  },
  planning: {
    label: 'Planning Done',
    icon: <MapCheckIcon className="w-4 h-4" />,
    bgColor: 'bg-blue-500',
    textColor: 'text-white',
  },
  building: {
    label: 'Build Complete',
    icon: <HammerCheckIcon className="w-4 h-4" />,
    bgColor: 'bg-amber-500',
    textColor: 'text-white',
  },
  testing: {
    label: 'Tests Passing',
    icon: <TestCheckIcon className="w-4 h-4" />,
    bgColor: 'bg-purple-500',
    textColor: 'text-white',
  },
  launch: {
    label: 'Launched!',
    icon: <RocketIcon className="w-4 h-4" />,
    bgColor: 'bg-green-500',
    textColor: 'text-white',
  },
};

export function MilestoneToast({
  phase,
  onDismiss,
  autoHideDelay = 4000,
}: MilestoneToastProps) {
  const [isVisible, setIsVisible] = useState(false);
  const [isExiting, setIsExiting] = useState(false);
  const config = PHASE_CONFIG[phase];

  const handleDismiss = useCallback(() => {
    setIsExiting(true);
    setTimeout(() => {
      setIsVisible(false);
      onDismiss();
    }, 300);
  }, [onDismiss]);

  // Enter animation
  useEffect(() => {
    const timer = setTimeout(() => setIsVisible(true), 50);
    return () => clearTimeout(timer);
  }, []);

  // Auto-hide after delay
  useEffect(() => {
    if (autoHideDelay > 0) {
      const timer = setTimeout(handleDismiss, autoHideDelay);
      return () => clearTimeout(timer);
    }
  }, [autoHideDelay, handleDismiss]);

  return (
    <div
      className={`fixed top-20 left-1/2 transform -translate-x-1/2 z-50 transition-all duration-300 ${
        isVisible && !isExiting
          ? 'translate-y-0 opacity-100'
          : '-translate-y-4 opacity-0'
      }`}
      role="alert"
      aria-live="polite"
    >
      <button
        onClick={handleDismiss}
        className={`flex items-center gap-2 px-4 py-2.5 ${config.bgColor} ${config.textColor} rounded-full shadow-lg hover:opacity-90 transition-opacity`}
      >
        {config.icon}
        <span className="font-medium text-sm">{config.label}</span>
      </button>
    </div>
  );
}

/**
 * Container component that manages multiple milestone toasts
 */
interface MilestoneToastContainerProps {
  completedPhases: BuildPhase[];
  onPhaseAcknowledged: (phase: BuildPhase) => void;
}

export function MilestoneToastContainer({
  completedPhases,
  onPhaseAcknowledged,
}: MilestoneToastContainerProps) {
  // Show only the most recent completed phase
  const latestPhase = completedPhases[completedPhases.length - 1];

  if (!latestPhase) {
    return null;
  }

  return (
    <MilestoneToast
      phase={latestPhase}
      onDismiss={() => onPhaseAcknowledged(latestPhase)}
    />
  );
}

/**
 * Hook to track phase transitions and trigger toasts
 */
export function useMilestoneToasts() {
  const [acknowledgedPhases, setAcknowledgedPhases] = useState<Set<BuildPhase>>(
    new Set()
  );
  const [pendingToast, setPendingToast] = useState<BuildPhase | null>(null);

  const showToast = useCallback((phase: BuildPhase) => {
    if (!acknowledgedPhases.has(phase)) {
      setPendingToast(phase);
    }
  }, [acknowledgedPhases]);

  const acknowledgePhase = useCallback((phase: BuildPhase) => {
    setAcknowledgedPhases((prev) => {
      const newSet = new Set(prev);
      newSet.add(phase);
      return newSet;
    });
    setPendingToast(null);
  }, []);

  const resetAcknowledgements = useCallback(() => {
    setAcknowledgedPhases(new Set());
    setPendingToast(null);
  }, []);

  return {
    pendingToast,
    acknowledgedPhases,
    showToast,
    acknowledgePhase,
    resetAcknowledgements,
  };
}

// Icons
function CheckIcon({ className }: { className?: string }) {
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
        d="M5 13l4 4L19 7"
      />
    </svg>
  );
}

function MapCheckIcon({ className }: { className?: string }) {
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
        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  );
}

function HammerCheckIcon({ className }: { className?: string }) {
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
        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  );
}

function TestCheckIcon({ className }: { className?: string }) {
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
        d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"
      />
    </svg>
  );
}

function RocketIcon({ className }: { className?: string }) {
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
        d="M15.59 14.37a6 6 0 01-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 006.16-12.12A14.98 14.98 0 009.63 8.41m6 6a14.98 14.98 0 01-5.63 1.59m5.63-1.59L21 3m-11.97 14.09a6 6 0 01-4.78-7.84"
      />
    </svg>
  );
}
