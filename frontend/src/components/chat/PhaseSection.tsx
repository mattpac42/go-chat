'use client';

import { useState, useCallback } from 'react';
import { Message } from '@/types';
import { BuildPhase } from './BuildPhaseProgress';
import { MessageBubble } from './MessageBubble';

interface PhaseSectionProps {
  phase: BuildPhase;
  label: string;
  messages: Message[];
  isCurrentPhase: boolean;
  defaultExpanded?: boolean;
  showBadgeForMessage?: (messageId: string) => boolean;
}

const PHASE_ICONS: Record<BuildPhase, React.ReactNode> = {
  discovery: <CompassIcon className="w-4 h-4" />,
  planning: <MapIcon className="w-4 h-4" />,
  building: <HammerIcon className="w-4 h-4" />,
  testing: <CheckBadgeIcon className="w-4 h-4" />,
  launch: <RocketIcon className="w-4 h-4" />,
};

const PHASE_COLORS: Record<BuildPhase, { bg: string; text: string; border: string }> = {
  discovery: { bg: 'bg-teal-50', text: 'text-teal-700', border: 'border-teal-200' },
  planning: { bg: 'bg-blue-50', text: 'text-blue-700', border: 'border-blue-200' },
  building: { bg: 'bg-amber-50', text: 'text-amber-700', border: 'border-amber-200' },
  testing: { bg: 'bg-purple-50', text: 'text-purple-700', border: 'border-purple-200' },
  launch: { bg: 'bg-green-50', text: 'text-green-700', border: 'border-green-200' },
};

export function PhaseSection({
  phase,
  label,
  messages,
  isCurrentPhase,
  defaultExpanded = true,
  showBadgeForMessage,
}: PhaseSectionProps) {
  const [isExpanded, setIsExpanded] = useState(defaultExpanded);
  const colors = PHASE_COLORS[phase];

  const toggleExpanded = useCallback(() => {
    setIsExpanded((prev) => !prev);
  }, []);

  if (messages.length === 0) {
    return null;
  }

  return (
    <div className="mb-4">
      {/* Section Header */}
      <button
        onClick={toggleExpanded}
        className={`w-full flex items-center justify-between px-4 py-2.5 rounded-lg transition-all ${
          colors.bg
        } ${colors.border} border ${
          isExpanded ? 'rounded-b-none' : ''
        } hover:opacity-90`}
        aria-expanded={isExpanded}
        aria-controls={`phase-${phase}-messages`}
      >
        <div className="flex items-center gap-3">
          <div className={`${colors.text}`}>{PHASE_ICONS[phase]}</div>
          <span className={`font-medium ${colors.text}`}>{label}</span>
          {isCurrentPhase && (
            <span className="px-2 py-0.5 text-xs font-medium bg-white/60 rounded-full text-gray-600">
              Current
            </span>
          )}
        </div>
        <div className="flex items-center gap-3">
          <span className={`text-sm ${colors.text} opacity-70`}>
            {messages.length} {messages.length === 1 ? 'message' : 'messages'}
          </span>
          <ChevronIcon
            className={`w-5 h-5 ${colors.text} transition-transform ${
              isExpanded ? 'rotate-180' : ''
            }`}
          />
        </div>
      </button>

      {/* Collapsible Content */}
      <div
        id={`phase-${phase}-messages`}
        className={`transition-all duration-300 overflow-hidden ${
          isExpanded
            ? 'max-h-[10000px] opacity-100'
            : 'max-h-0 opacity-0'
        }`}
      >
        <div className={`${colors.border} border border-t-0 rounded-b-lg p-4 bg-white`}>
          {messages.map((message) => (
            <MessageBubble
              key={message.id}
              message={message}
              showBadge={showBadgeForMessage?.(message.id) ?? false}
            />
          ))}
        </div>
      </div>
    </div>
  );
}

/**
 * Container that renders messages grouped by phase
 */
interface PhasedMessageListProps {
  messagesByPhase: Map<BuildPhase, Message[]>;
  currentPhase: BuildPhase;
  showBadgeForMessage?: (messageId: string) => boolean;
}

export function PhasedMessageList({
  messagesByPhase,
  currentPhase,
  showBadgeForMessage,
}: PhasedMessageListProps) {
  const phaseOrder: BuildPhase[] = ['discovery', 'planning', 'building', 'testing', 'launch'];
  const phaseLabels: Record<BuildPhase, string> = {
    discovery: 'Discovery',
    planning: 'Planning',
    building: 'Building',
    testing: 'Testing',
    launch: 'Launch',
  };

  return (
    <div className="space-y-2">
      {phaseOrder.map((phase) => {
        const messages = messagesByPhase.get(phase) || [];
        if (messages.length === 0) return null;

        return (
          <PhaseSection
            key={phase}
            phase={phase}
            label={phaseLabels[phase]}
            messages={messages}
            isCurrentPhase={phase === currentPhase}
            defaultExpanded={phase === currentPhase}
            showBadgeForMessage={showBadgeForMessage}
          />
        );
      })}
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

function ChevronIcon({ className }: { className?: string }) {
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
        d="M19 9l-7 7-7-7"
      />
    </svg>
  );
}
