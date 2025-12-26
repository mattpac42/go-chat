'use client';

import { Nudge } from '@/types/achievements';

interface NudgePopoverProps {
  nudge: Nudge;
  onAccept: () => void;
  onDismiss: () => void;
}

// Inline SVG Icons - Dynamic icon selection based on nudge.icon
function EyeIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
      />
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
      />
    </svg>
  );
}

function GitBranchIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M6 3v12M18 9a3 3 0 100-6 3 3 0 000 6zM6 21a3 3 0 100-6 3 3 0 000 6zM18 9a9 9 0 01-9 9"
      />
    </svg>
  );
}

function TrendingUpIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
      />
    </svg>
  );
}

function LightbulbIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
      />
    </svg>
  );
}

function NetworkIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M8.288 15.038a5.25 5.25 0 017.424 0M5.106 11.856c3.807-3.808 9.98-3.808 13.788 0M1.924 8.674c5.565-5.565 14.587-5.565 20.152 0M12 20h.01"
      />
    </svg>
  );
}

function DownloadIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
      />
    </svg>
  );
}

function ChevronRightIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9 5l7 7-7 7"
      />
    </svg>
  );
}

function CloseIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M6 18L18 6M6 6l12 12"
      />
    </svg>
  );
}

// Map icon names to components
const ICON_MAP: Record<string, React.FC<{ className?: string }>> = {
  eye: EyeIcon,
  'git-branch': GitBranchIcon,
  'trending-up': TrendingUpIcon,
  lightbulb: LightbulbIcon,
  network: NetworkIcon,
  download: DownloadIcon,
};

/**
 * NudgePopover - Contextual suggestion popover with accept/dismiss buttons
 *
 * Features:
 * - Fixed position at bottom of screen
 * - Slide-up animation
 * - Dynamic icon based on nudge type
 * - Accept action button and dismiss option
 * - Responsive: full-width on mobile, fixed-width on desktop
 */
export function NudgePopover({
  nudge,
  onAccept,
  onDismiss,
}: NudgePopoverProps) {
  // Get icon component based on nudge.icon string
  const IconComponent = ICON_MAP[nudge.icon] || LightbulbIcon;

  return (
    <div
      role="dialog"
      aria-label={nudge.title}
      className="fixed bottom-20 left-4 right-4 md:left-auto md:right-4 md:w-80 z-40 animate-[slideUp_0.3s_ease-out]"
    >
      <div className="bg-teal-600 text-white rounded-lg shadow-lg p-4">
        <div className="flex items-start gap-3">
          {/* Icon */}
          <div
            className="flex-shrink-0 bg-white/20 rounded-full p-2"
            aria-hidden="true"
          >
            <IconComponent className="w-5 h-5" />
          </div>

          {/* Content */}
          <div className="flex-1 min-w-0">
            <p className="font-semibold">{nudge.title}</p>
            <p className="text-white/80 text-sm mt-1">{nudge.message}</p>

            {/* Actions */}
            <div className="flex items-center gap-2 mt-3">
              <button
                onClick={onAccept}
                className="flex items-center gap-1 bg-white text-teal-600 px-3 py-1.5 rounded-md text-sm font-medium hover:bg-white/90 transition-colors focus:outline-none focus:ring-2 focus:ring-white/50"
              >
                {nudge.action}
                <ChevronRightIcon className="w-4 h-4" />
              </button>
              <button
                onClick={onDismiss}
                className="text-white/60 hover:text-white text-sm transition-colors"
              >
                Maybe later
              </button>
            </div>
          </div>

          {/* Close Button */}
          <button
            onClick={onDismiss}
            className="flex-shrink-0 text-white/60 hover:text-white transition-colors focus:outline-none focus:ring-2 focus:ring-white/50 rounded"
            aria-label="Dismiss suggestion"
          >
            <CloseIcon className="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  );
}

export type { NudgePopoverProps };
