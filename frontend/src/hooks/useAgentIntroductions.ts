'use client';

import { useState, useEffect, useCallback } from 'react';
import { AgentType } from '@/types';

const STORAGE_KEY_PREFIX = 'agents_introduced_';

/**
 * Hook to track which agents have been "introduced" (met) per project.
 * Used to show a "NEW" badge on the first appearance of each agent type.
 * State is persisted in localStorage per project.
 */
export function useAgentIntroductions(projectId: string) {
  const [introducedAgents, setIntroducedAgents] = useState<Set<AgentType>>(
    new Set()
  );

  // Load from localStorage on mount or when projectId changes
  useEffect(() => {
    if (!projectId) return;

    const stored = localStorage.getItem(`${STORAGE_KEY_PREFIX}${projectId}`);
    if (stored) {
      try {
        const agents = JSON.parse(stored) as AgentType[];
        setIntroducedAgents(new Set(agents));
      } catch {
        // Invalid data, start fresh
        setIntroducedAgents(new Set());
      }
    } else {
      setIntroducedAgents(new Set());
    }
  }, [projectId]);

  // Check if the user has met this agent before
  const hasMetAgent = useCallback(
    (agentType: AgentType): boolean => {
      return introducedAgents.has(agentType);
    },
    [introducedAgents]
  );

  // Mark an agent as met/introduced
  const markAgentMet = useCallback(
    (agentType: AgentType) => {
      if (!projectId) return;

      setIntroducedAgents((prev) => {
        if (prev.has(agentType)) return prev;

        const next = new Set(prev);
        next.add(agentType);

        // Persist to localStorage
        localStorage.setItem(
          `${STORAGE_KEY_PREFIX}${projectId}`,
          JSON.stringify(Array.from(next))
        );

        return next;
      });
    },
    [projectId]
  );

  // Reset all introductions for this project (useful for testing)
  const resetIntroductions = useCallback(() => {
    if (!projectId) return;

    localStorage.removeItem(`${STORAGE_KEY_PREFIX}${projectId}`);
    setIntroducedAgents(new Set());
  }, [projectId]);

  return {
    hasMetAgent,
    markAgentMet,
    resetIntroductions,
  };
}
