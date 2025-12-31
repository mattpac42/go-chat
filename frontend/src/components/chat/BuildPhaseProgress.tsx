'use client';

import { useMemo } from 'react';
import { Message } from '@/types';

export type BuildPhase = 'discovery' | 'planning' | 'building' | 'testing' | 'launch';

export interface PhaseInfo {
  phase: BuildPhase;
  label: string;
  icon: React.ReactNode;
  description: string;
}

const PHASES: PhaseInfo[] = [
  {
    phase: 'discovery',
    label: 'Discovery',
    icon: <CompassIcon className="w-4 h-4" />,
    description: 'Understanding your needs',
  },
  {
    phase: 'planning',
    label: 'Planning',
    icon: <MapIcon className="w-4 h-4" />,
    description: 'Designing the architecture',
  },
  {
    phase: 'building',
    label: 'Building',
    icon: <HammerIcon className="w-4 h-4" />,
    description: 'Writing the code',
  },
  {
    phase: 'testing',
    label: 'Testing',
    icon: <CheckBadgeIcon className="w-4 h-4" />,
    description: 'Verifying everything works',
  },
  {
    phase: 'launch',
    label: 'Launch',
    icon: <RocketIcon className="w-4 h-4" />,
    description: 'Deploying your project',
  },
];

export function getPhaseIndex(phase: BuildPhase): number {
  return PHASES.findIndex((p) => p.phase === phase);
}

/**
 * Detect the current build phase based on message content.
 * Uses simple heuristics to determine which phase we're in.
 */
export function detectPhaseFromMessage(message: Message): BuildPhase {
  const content = message.content.toLowerCase();

  // Launch indicators
  if (
    content.includes('deploy') ||
    content.includes('launch') ||
    content.includes('production') ||
    content.includes('go live') ||
    content.includes('deployed successfully')
  ) {
    return 'launch';
  }

  // Testing indicators
  if (
    content.includes('test') ||
    content.includes('testing') ||
    content.includes('spec') ||
    content.includes('coverage') ||
    content.includes('verify') ||
    content.includes('assertion')
  ) {
    return 'testing';
  }

  // Building indicators (code blocks or implementation talk)
  if (
    content.includes('```') ||
    content.includes('implement') ||
    content.includes('function') ||
    content.includes('component') ||
    content.includes('created file') ||
    content.includes('updated file')
  ) {
    return 'building';
  }

  // Planning indicators
  if (
    content.includes('architecture') ||
    content.includes('structure') ||
    content.includes('file plan') ||
    content.includes('design') ||
    content.includes('schema') ||
    content.includes('api routes') ||
    content.includes('database')
  ) {
    return 'planning';
  }

  // Default to discovery
  return 'discovery';
}

/**
 * Detect the current overall phase from all messages
 */
export function detectCurrentPhase(messages: Message[]): BuildPhase {
  if (messages.length === 0) return 'discovery';

  // Check the last few messages to determine current phase
  const recentMessages = messages.slice(-5);
  const phases = recentMessages.map(detectPhaseFromMessage);

  // Return the highest (most advanced) phase detected
  const phaseIndices = phases.map(getPhaseIndex);
  const maxIndex = Math.max(...phaseIndices);
  return PHASES[maxIndex].phase;
}

/**
 * Detect phase with context - for user messages that default to 'discovery',
 * inherit the phase from the following assistant message if available.
 */
export function detectPhaseWithContext(
  messages: Message[],
  index: number
): BuildPhase {
  const message = messages[index];
  const detectedPhase = detectPhaseFromMessage(message);

  // If not a user message or if it matched a specific phase, use as-is
  if (message.role !== 'user' || detectedPhase !== 'discovery') {
    return detectedPhase;
  }

  // For user messages that defaulted to 'discovery', look ahead
  // to find the next assistant message and inherit its phase
  for (let i = index + 1; i < messages.length; i++) {
    if (messages[i].role === 'assistant') {
      const assistantPhase = detectPhaseFromMessage(messages[i]);
      return assistantPhase;
    }
  }

  // No following assistant message found, keep as discovery
  return 'discovery';
}

/**
 * Group messages by their detected phase
 */
export function groupMessagesByPhase(messages: Message[]): Map<BuildPhase, Message[]> {
  const groups = new Map<BuildPhase, Message[]>();

  // Initialize all phases with empty arrays
  PHASES.forEach((p) => groups.set(p.phase, []));

  messages.forEach((message, index) => {
    const phase = detectPhaseWithContext(messages, index);
    const existing = groups.get(phase) || [];
    existing.push(message);
    groups.set(phase, existing);
  });

  return groups;
}

interface BuildPhaseProgressProps {
  messages: Message[];
  isDiscoveryMode: boolean;
  discoveryComplete: boolean;
  currentPhase?: BuildPhase;
  onPhaseClick?: (phase: BuildPhase) => void;
  showPhasedView?: boolean;
  onTogglePhasedView?: () => void;
}

export function BuildPhaseProgress({
  messages,
  isDiscoveryMode,
  discoveryComplete,
  currentPhase: overridePhase,
  onPhaseClick,
  showPhasedView = false,
  onTogglePhasedView,
}: BuildPhaseProgressProps) {
  const currentPhase = useMemo(() => {
    if (overridePhase) return overridePhase;
    if (isDiscoveryMode && !discoveryComplete) return 'discovery';
    return detectCurrentPhase(messages);
  }, [messages, isDiscoveryMode, discoveryComplete, overridePhase]);

  const currentIndex = getPhaseIndex(currentPhase);

  return (
    <div className="w-full bg-white border-b border-gray-200 px-4 py-3">
      <div className="max-w-4xl mx-auto">
        {/* Desktop view */}
        <div className="hidden md:flex items-center gap-4">
          {/* View segmented control */}
          {onTogglePhasedView && (
            <div className="flex items-center rounded-lg border border-gray-200 bg-gray-50 p-0.5">
              <button
                onClick={() => showPhasedView && onTogglePhasedView()}
                className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm font-medium transition-all ${
                  !showPhasedView
                    ? 'bg-teal-500 text-white shadow-sm'
                    : 'text-gray-600 hover:text-gray-800'
                }`}
                aria-pressed={!showPhasedView}
              >
                <ListIcon className="w-4 h-4" />
                <span className="hidden xl:inline">Timeline</span>
              </button>
              <button
                onClick={() => !showPhasedView && onTogglePhasedView()}
                className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm font-medium transition-all ${
                  showPhasedView
                    ? 'bg-teal-500 text-white shadow-sm'
                    : 'text-gray-600 hover:text-gray-800'
                }`}
                aria-pressed={showPhasedView}
              >
                <GridIcon className="w-4 h-4" />
                <span className="hidden xl:inline">By Phase</span>
              </button>
            </div>
          )}

          {/* Phase indicators */}
          <div className="flex items-center justify-between flex-1">
          {PHASES.map((phaseInfo, index) => {
            const isCompleted = index < currentIndex;
            const isCurrent = index === currentIndex;
            const isUpcoming = index > currentIndex;

            return (
              <div key={phaseInfo.phase} className="flex items-center flex-1">
                <button
                  onClick={() => onPhaseClick?.(phaseInfo.phase)}
                  className={`flex items-center gap-2 px-3 py-2 rounded-lg transition-all ${
                    isCurrent
                      ? 'bg-teal-50 text-teal-700'
                      : isCompleted
                      ? 'text-teal-600 hover:bg-teal-50/50'
                      : 'text-gray-400'
                  } ${onPhaseClick ? 'cursor-pointer' : 'cursor-default'}`}
                  disabled={!onPhaseClick}
                >
                  <div
                    className={`flex items-center justify-center w-8 h-8 rounded-full transition-colors ${
                      isCurrent
                        ? 'bg-teal-500 text-white'
                        : isCompleted
                        ? 'bg-teal-100 text-teal-600'
                        : 'bg-gray-100 text-gray-400'
                    }`}
                  >
                    {isCompleted ? (
                      <CheckIcon className="w-4 h-4" />
                    ) : (
                      phaseInfo.icon
                    )}
                  </div>
                  <div className="hidden lg:block">
                    <div
                      className={`text-sm font-medium ${
                        isCurrent
                          ? 'text-teal-700'
                          : isCompleted
                          ? 'text-teal-600'
                          : 'text-gray-400'
                      }`}
                    >
                      {phaseInfo.label}
                    </div>
                    {isCurrent && (
                      <div className="text-xs text-teal-600/70">
                        {phaseInfo.description}
                      </div>
                    )}
                  </div>
                </button>

                {/* Connector line */}
                {index < PHASES.length - 1 && (
                  <div
                    className={`flex-1 h-0.5 mx-2 transition-colors ${
                      index < currentIndex ? 'bg-teal-300' : 'bg-gray-200'
                    }`}
                  />
                )}
              </div>
            );
          })}
          </div>
        </div>

        {/* Mobile view - compact progress bar */}
        <div className="md:hidden">
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center gap-2">
              <div className="flex items-center justify-center w-7 h-7 rounded-full bg-teal-500 text-white">
                {PHASES[currentIndex].icon}
              </div>
              <div>
                <div className="text-sm font-medium text-teal-700">
                  {PHASES[currentIndex].label}
                </div>
                <div className="text-xs text-gray-500">
                  Step {currentIndex + 1} of {PHASES.length}
                </div>
              </div>
            </div>
          </div>

          {/* Progress dots */}
          <div className="flex items-center gap-1.5">
            {PHASES.map((phaseInfo, index) => {
              const isCompleted = index < currentIndex;
              const isCurrent = index === currentIndex;

              return (
                <button
                  key={phaseInfo.phase}
                  onClick={() => onPhaseClick?.(phaseInfo.phase)}
                  className={`flex-1 h-1.5 rounded-full transition-colors ${
                    isCurrent
                      ? 'bg-teal-500'
                      : isCompleted
                      ? 'bg-teal-300'
                      : 'bg-gray-200'
                  }`}
                  aria-label={`${phaseInfo.label}: ${
                    isCompleted ? 'completed' : isCurrent ? 'current' : 'upcoming'
                  }`}
                />
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}

// Icons
function CompassIcon({ className }: { className?: string }) {
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
        d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
      />
    </svg>
  );
}

function MapIcon({ className }: { className?: string }) {
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
        d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"
      />
    </svg>
  );
}

function HammerIcon({ className }: { className?: string }) {
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
        d="M15.59 14.37a6 6 0 01-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 006.16-12.12A14.98 14.98 0 009.63 8.41m6 6a14.98 14.98 0 01-5.63 1.59m5.63-1.59L21 3m-6.41 11.37a6 6 0 11-5.84-7.38v4.8m5.84 2.58L3 21"
      />
    </svg>
  );
}

function CheckBadgeIcon({ className }: { className?: string }) {
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

function ListIcon({ className }: { className?: string }) {
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
        d="M4 6h16M4 12h16M4 18h16"
      />
    </svg>
  );
}

function GridIcon({ className }: { className?: string }) {
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
        d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
      />
    </svg>
  );
}
