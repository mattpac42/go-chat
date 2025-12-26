'use client';

import { AgentType, AGENT_CONFIG } from '@/types';
import { AgentIcon } from './AgentIcon';

interface AgentHeaderProps {
  agentType: AgentType;
  showBadge?: boolean;
  timestamp?: string;
}

/**
 * AgentHeader - Displays agent identification above assistant messages
 *
 * Features:
 * - Agent icon with accent color
 * - Display name (full on desktop, short on mobile via responsive classes)
 * - Optional "NEW" badge for first introduction
 * - Optional timestamp
 */
export function AgentHeader({
  agentType,
  showBadge = false,
  timestamp,
}: AgentHeaderProps) {
  const config = AGENT_CONFIG[agentType];

  return (
    <div className="flex items-center gap-2 mb-2">
      {/* Agent icon */}
      <AgentIcon
        type={agentType}
        className="w-4 h-4 md:w-5 md:h-5"
      />

      {/* Agent name - responsive: short name on mobile, full name on desktop */}
      <span
        className="text-sm font-medium"
        style={{ color: config.color }}
      >
        <span className="md:hidden">{config.shortName}</span>
        <span className="hidden md:inline">{config.displayName}</span>
      </span>

      {/* NEW badge for first introduction */}
      {showBadge && (
        <span className="px-1.5 py-0.5 text-xs font-medium bg-violet-100 text-violet-700 rounded">
          NEW
        </span>
      )}

      {/* Optional timestamp */}
      {timestamp && (
        <span className="ml-auto text-xs text-gray-500">
          {timestamp}
        </span>
      )}
    </div>
  );
}
