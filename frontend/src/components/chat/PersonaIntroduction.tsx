'use client';

import { Message, AgentType, AGENT_CONFIG } from '@/types';
import { AgentHeader } from './AgentHeader';

interface PersonaIntroductionProps {
  message: Message;
  isTeamIntro?: boolean; // True for Root's team introduction message
}

/**
 * A special message bubble for persona introductions.
 * These appear when transitioning from discovery to building phase.
 *
 * Team introductions (from Root) get a distinctive gradient background.
 * Individual persona intros get their agent color styling.
 */
export function PersonaIntroduction({
  message,
  isTeamIntro = false,
}: PersonaIntroductionProps) {
  const agentType = message.agentType || 'product_manager';
  const agentConfig = AGENT_CONFIG[agentType];

  // Team introduction from Root gets special styling
  if (isTeamIntro) {
    return (
      <div className="flex justify-start mb-4" data-testid="team-introduction">
        <div className="max-w-[85%] md:max-w-[70%] rounded-2xl px-4 py-4 rounded-bl-md bg-gradient-to-br from-teal-50 to-cyan-50 border border-teal-200">
          {/* Header with special "Team" label */}
          <div className="flex items-center gap-2 mb-3">
            <div
              className="w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-medium"
              style={{ backgroundColor: agentConfig.color }}
            >
              <TeamIcon className="w-4 h-4" />
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-gray-900">
                {agentConfig.displayName}
              </span>
              <span className="text-xs text-teal-600 bg-teal-100 px-2 py-0.5 rounded-full">
                Introducing the team
              </span>
            </div>
          </div>

          {/* Content */}
          <div className="text-gray-700 whitespace-pre-line leading-relaxed">
            {message.content}
          </div>
        </div>
      </div>
    );
  }

  // Individual persona introduction
  return (
    <div
      className="flex justify-start mb-4"
      data-testid={`persona-intro-${agentType}`}
    >
      <div
        className="max-w-[85%] md:max-w-[70%] rounded-2xl px-4 py-3 rounded-bl-md"
        style={{
          backgroundColor: agentConfig.bgColor,
          borderLeft: `3px solid ${agentConfig.color}`,
        }}
      >
        {/* Agent header with wave emoji */}
        <div className="flex items-center gap-2 mb-2">
          <AgentHeader agentType={agentType} showBadge={true} />
        </div>

        {/* Content with friendly styling */}
        <div className="text-gray-700 leading-relaxed">{message.content}</div>
      </div>
    </div>
  );
}

/**
 * Helper to check if a message is an introduction message
 */
export function isIntroductionMessage(message: Message): boolean {
  return message.id.startsWith('intro-');
}

/**
 * Helper to check if a message is Root's team introduction
 */
export function isTeamIntroMessage(message: Message): boolean {
  return message.id.startsWith('intro-root-team-');
}

function TeamIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
      />
    </svg>
  );
}
