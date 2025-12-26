'use client';

import { UserProgress, LearningLevel } from '@/types/achievements';

interface ProgressBadgeProps {
  progress: UserProgress;
  compact?: boolean;
  onClick?: () => void;
}

// Inline SVG Icons
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

// Level configuration
const LEVEL_NAMES: Record<LearningLevel, string> = {
  1: 'Explorer',
  2: 'Navigator',
  3: 'Technologist',
  4: 'Developer',
};

const LEVEL_COLORS: Record<LearningLevel, string> = {
  1: 'bg-emerald-500',
  2: 'bg-blue-500',
  3: 'bg-purple-500',
  4: 'bg-amber-500',
};

const LEVEL_TEXT_COLORS: Record<LearningLevel, string> = {
  1: 'text-emerald-500',
  2: 'text-blue-500',
  3: 'text-purple-500',
  4: 'text-amber-500',
};

// Points required to advance to each level
const LEVEL_THRESHOLDS: Record<LearningLevel, number> = {
  1: 0,    // Start here
  2: 50,   // Need 50 points for Level 2
  3: 150,  // Need 150 points for Level 3
  4: 300,  // Need 300 points for Level 4
};

/**
 * ProgressBadge - Shows current learning level and points progress
 *
 * Two modes:
 * - compact: Small inline badge showing level number and points
 * - full: Card with level name, points, and progress bar to next level
 */
export function ProgressBadge({
  progress,
  compact = false,
  onClick,
}: ProgressBadgeProps) {
  const currentLevel = progress.currentLevel;
  const levelName = LEVEL_NAMES[currentLevel];
  const levelColor = LEVEL_COLORS[currentLevel];

  // Calculate progress to next level
  const currentThreshold = LEVEL_THRESHOLDS[currentLevel];
  const nextLevel = Math.min(currentLevel + 1, 4) as LearningLevel;
  const nextThreshold = LEVEL_THRESHOLDS[nextLevel];
  const pointsInCurrentLevel = progress.totalPoints - currentThreshold;
  const pointsNeededForNext = nextThreshold - currentThreshold;
  const progressPercent =
    currentLevel >= 4
      ? 100
      : Math.min((pointsInCurrentLevel / pointsNeededForNext) * 100, 100);

  if (compact) {
    return (
      <button
        onClick={onClick}
        className="flex items-center gap-1.5 px-2 py-1 rounded-full bg-gray-100 hover:bg-gray-200 transition-colors"
        aria-label={`Level ${currentLevel}: ${levelName}, ${progress.totalPoints} points`}
      >
        <div
          className={`w-2 h-2 rounded-full ${levelColor}`}
          aria-hidden="true"
        />
        <span className="text-xs font-medium text-gray-600">
          Lvl {currentLevel}
        </span>
        <span className="text-xs text-gray-400 flex items-center gap-0.5">
          <StarIcon className="w-3 h-3" />
          {progress.totalPoints}
        </span>
      </button>
    );
  }

  return (
    <div
      className={`bg-white rounded-lg shadow-sm border border-gray-200 p-4 ${
        onClick ? 'cursor-pointer hover:shadow-md' : ''
      } transition-shadow`}
      onClick={onClick}
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
      onKeyDown={
        onClick
          ? (e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                onClick();
              }
            }
          : undefined
      }
      aria-label={
        onClick
          ? `Level ${currentLevel}: ${levelName}, ${progress.totalPoints} points. Click for details.`
          : undefined
      }
    >
      {/* Header: Level and Points */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <div
            className={`w-3 h-3 rounded-full ${levelColor}`}
            aria-hidden="true"
          />
          <span className="font-medium text-gray-800">
            Level {currentLevel}: {levelName}
          </span>
        </div>
        <div className="flex items-center gap-1 text-amber-500">
          <StarIcon className="w-4 h-4" />
          <span className="font-semibold">{progress.totalPoints}</span>
        </div>
      </div>

      {/* Progress to next level */}
      {currentLevel < 4 ? (
        <div className="space-y-1">
          <div className="flex justify-between text-xs text-gray-500">
            <span>Progress to Level {nextLevel}</span>
            <span>
              {pointsInCurrentLevel}/{pointsNeededForNext} points
            </span>
          </div>
          <div
            className="h-2 bg-gray-100 rounded-full overflow-hidden"
            role="progressbar"
            aria-valuenow={pointsInCurrentLevel}
            aria-valuemin={0}
            aria-valuemax={pointsNeededForNext}
            aria-label={`${pointsInCurrentLevel} of ${pointsNeededForNext} points to Level ${nextLevel}`}
          >
            <div
              className={`h-full ${levelColor} transition-all duration-300 ease-out`}
              style={{ width: `${progressPercent}%` }}
            />
          </div>
        </div>
      ) : (
        <div className="text-sm text-gray-500">
          Max level reached! Keep exploring to unlock more achievements.
        </div>
      )}

      {/* Stats summary */}
      <div className="mt-3 pt-3 border-t border-gray-100 grid grid-cols-3 gap-2 text-center">
        <div>
          <div className="text-lg font-semibold text-gray-700">
            {progress.filesViewedCount}
          </div>
          <div className="text-xs text-gray-400">Files Viewed</div>
        </div>
        <div>
          <div className="text-lg font-semibold text-gray-700">
            {progress.codeViewsCount}
          </div>
          <div className="text-xs text-gray-400">Code Views</div>
        </div>
        <div>
          <div className="text-lg font-semibold text-gray-700">
            {progress.treeExpansionsCount}
          </div>
          <div className="text-xs text-gray-400">Tree Expands</div>
        </div>
      </div>
    </div>
  );
}

export type { ProgressBadgeProps };
