'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { Message, AgentType, AGENT_CONFIG } from '@/types';
import { DiscoveryStage } from '@/types/discovery';

const STORAGE_KEY_PREFIX = 'team_introduced_';

/**
 * Personas that should be introduced when transitioning to building phase
 */
const TEAM_PERSONAS: AgentType[] = ['designer', 'developer'];

/**
 * Hook to manage team persona introductions when transitioning from discovery to building.
 *
 * When discovery completes, Root introduces the team (Bloom and Harvest),
 * then each persona introduces themselves briefly.
 *
 * @param projectId - The project ID to track introductions for
 * @param currentStage - Current discovery stage
 * @param messages - Current messages in the chat
 */
export function usePersonaIntroductions(
  projectId: string,
  currentStage: DiscoveryStage,
  messages: Message[]
) {
  const [hasIntroducedTeam, setHasIntroducedTeam] = useState(false);
  const [introMessagesShown, setIntroMessagesShown] = useState(false);

  // Load persisted state on mount
  useEffect(() => {
    if (!projectId) return;

    const stored = localStorage.getItem(`${STORAGE_KEY_PREFIX}${projectId}`);
    if (stored === 'true') {
      setHasIntroducedTeam(true);
      setIntroMessagesShown(true);
    }
  }, [projectId]);

  // Detect when we should show introductions:
  // - Discovery just completed (stage is 'complete')
  // - We haven't shown introductions yet
  // - There's at least one message from a non-Root agent OR we're transitioning
  const shouldShowIntroductions = useMemo(() => {
    if (hasIntroducedTeam || introMessagesShown) return false;
    if (currentStage !== 'complete') return false;

    // Check if there's a message from designer or developer (building phase started)
    const hasNonRootMessage = messages.some(
      (m) =>
        m.role === 'assistant' &&
        m.agentType &&
        TEAM_PERSONAS.includes(m.agentType)
    );

    return hasNonRootMessage;
  }, [hasIntroducedTeam, introMessagesShown, currentStage, messages]);

  // Mark introductions as shown and persist
  const markIntroductionsShown = useCallback(() => {
    if (!projectId) return;

    setHasIntroducedTeam(true);
    setIntroMessagesShown(true);
    localStorage.setItem(`${STORAGE_KEY_PREFIX}${projectId}`, 'true');
  }, [projectId]);

  // Generate the introduction messages that should be injected
  const introductionMessages = useMemo((): Message[] => {
    if (!shouldShowIntroductions) return [];

    const now = new Date().toISOString();
    const introMessages: Message[] = [];

    // Root's team introduction message
    introMessages.push({
      id: `intro-root-team-${projectId}`,
      projectId,
      role: 'assistant',
      content: `Great! Now that we understand your project, let me introduce the team who'll help build it.\n\n${AGENT_CONFIG.designer.rootIntro}\n\n${AGENT_CONFIG.developer.rootIntro}`,
      timestamp: now,
      agentType: 'product_manager',
    });

    // Each persona's self-introduction
    for (const persona of TEAM_PERSONAS) {
      const config = AGENT_CONFIG[persona];
      introMessages.push({
        id: `intro-${persona}-${projectId}`,
        projectId,
        role: 'assistant',
        content: config.selfIntro,
        timestamp: now,
        agentType: persona,
      });
    }

    return introMessages;
  }, [shouldShowIntroductions, projectId]);

  // Get the index where introductions should be inserted
  // (before the first non-Root agent message)
  const getInsertionIndex = useCallback(
    (msgs: Message[]): number => {
      if (!shouldShowIntroductions) return -1;

      const firstNonRootIndex = msgs.findIndex(
        (m) =>
          m.role === 'assistant' &&
          m.agentType &&
          TEAM_PERSONAS.includes(m.agentType)
      );

      return firstNonRootIndex;
    },
    [shouldShowIntroductions]
  );

  // Process messages to inject introductions at the right place
  const processMessagesWithIntroductions = useCallback(
    (msgs: Message[]): Message[] => {
      if (!shouldShowIntroductions || introductionMessages.length === 0) {
        return msgs;
      }

      const insertionIndex = getInsertionIndex(msgs);
      if (insertionIndex === -1) return msgs;

      // Insert introduction messages before the first non-Root agent message
      const result = [
        ...msgs.slice(0, insertionIndex),
        ...introductionMessages,
        ...msgs.slice(insertionIndex),
      ];

      // Mark as shown after processing
      // Use setTimeout to avoid state update during render
      setTimeout(() => markIntroductionsShown(), 0);

      return result;
    },
    [
      shouldShowIntroductions,
      introductionMessages,
      getInsertionIndex,
      markIntroductionsShown,
    ]
  );

  // Reset introductions for this project (useful for testing)
  const resetIntroductions = useCallback(() => {
    if (!projectId) return;

    localStorage.removeItem(`${STORAGE_KEY_PREFIX}${projectId}`);
    setHasIntroducedTeam(false);
    setIntroMessagesShown(false);
  }, [projectId]);

  return {
    shouldShowIntroductions,
    introductionMessages,
    processMessagesWithIntroductions,
    markIntroductionsShown,
    resetIntroductions,
    hasIntroducedTeam,
  };
}
