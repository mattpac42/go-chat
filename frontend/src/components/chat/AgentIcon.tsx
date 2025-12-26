'use client';

import { AgentType, AGENT_CONFIG } from '@/types';

interface AgentIconProps {
  type: AgentType;
  className?: string;
}

/**
 * ProductGuideIcon - Target with checkmark (violet)
 * Represents goal-focused product management
 */
function ProductGuideIcon({ className = 'w-5 h-5' }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  );
}

/**
 * UXExpertIcon - Layout grid (orange)
 * Represents interface and user experience design
 */
function UXExpertIcon({ className = 'w-5 h-5' }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z"
      />
    </svg>
  );
}

/**
 * DeveloperIcon - Code brackets (emerald)
 * Represents code implementation and development
 */
function DeveloperIcon({ className = 'w-5 h-5' }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"
      />
    </svg>
  );
}

/**
 * AgentIcon - Renders the appropriate icon for each agent type
 * Uses the color from AGENT_CONFIG for consistent theming
 */
export function AgentIcon({ type, className = 'w-5 h-5' }: AgentIconProps) {
  const config = AGENT_CONFIG[type];
  const style = { color: config.color };

  switch (type) {
    case 'product_manager':
      return (
        <span style={style}>
          <ProductGuideIcon className={className} />
        </span>
      );
    case 'designer':
      return (
        <span style={style}>
          <UXExpertIcon className={className} />
        </span>
      );
    case 'developer':
      return (
        <span style={style}>
          <DeveloperIcon className={className} />
        </span>
      );
    default:
      return (
        <span style={style}>
          <DeveloperIcon className={className} />
        </span>
      );
  }
}
