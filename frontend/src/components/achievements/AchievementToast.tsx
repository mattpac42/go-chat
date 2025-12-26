'use client';

import { useEffect, useState } from 'react';
import { UserAchievement } from '@/types/achievements';

interface AchievementToastProps {
  achievement: UserAchievement;
  onDismiss: () => void;
  autoDismissMs?: number;
}

// Inline SVG Icons
function TrophyIcon({ className }: { className?: string }) {
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
        d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"
      />
    </svg>
  );
}

function StarIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="currentColor"
      viewBox="0 0 24 24"
    >
      <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
    </svg>
  );
}

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
 * AchievementToast - Animated toast notification when achievement is unlocked
 *
 * Features:
 * - Slide-up animation on mount
 * - Auto-dismiss after 5 seconds (configurable)
 * - Manual dismiss with close button
 * - Shows achievement name, description, and points earned
 * - Celebratory gradient background
 */
export function AchievementToast({
  achievement,
  onDismiss,
  autoDismissMs = 5000,
}: AchievementToastProps) {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    // Trigger animation on mount
    const animateIn = requestAnimationFrame(() => {
      setIsVisible(true);
    });

    // Auto-dismiss timer
    const timer = setTimeout(() => {
      setIsVisible(false);
      // Wait for exit animation before calling onDismiss
      setTimeout(onDismiss, 300);
    }, autoDismissMs);

    return () => {
      cancelAnimationFrame(animateIn);
      clearTimeout(timer);
    };
  }, [autoDismissMs, onDismiss]);

  const handleDismiss = () => {
    setIsVisible(false);
    setTimeout(onDismiss, 300);
  };

  const { achievement: ach } = achievement;
  if (!ach) return null;

  return (
    <div
      role="alert"
      aria-live="polite"
      className={`fixed bottom-4 right-4 z-50 transform transition-all duration-300 ease-out ${
        isVisible
          ? 'translate-y-0 opacity-100'
          : 'translate-y-8 opacity-0'
      }`}
    >
      <div className="bg-gradient-to-r from-amber-500 to-orange-500 text-white rounded-lg shadow-lg p-4 max-w-sm">
        <div className="flex items-start gap-3">
          {/* Trophy Icon */}
          <div className="flex-shrink-0 bg-white/20 rounded-full p-2">
            <TrophyIcon className="w-6 h-6" />
          </div>

          {/* Content */}
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <span className="font-semibold">Achievement Unlocked!</span>
              <span className="text-amber-200 text-sm flex items-center gap-1">
                <StarIcon className="w-3 h-3" />
                +{ach.points}
              </span>
            </div>
            <p className="font-bold text-lg">{ach.name}</p>
            <p className="text-white/80 text-sm">{ach.description}</p>
          </div>

          {/* Dismiss Button */}
          <button
            onClick={handleDismiss}
            className="flex-shrink-0 text-white/60 hover:text-white transition-colors"
            aria-label="Dismiss notification"
          >
            <CloseIcon className="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  );
}

export type { AchievementToastProps };
