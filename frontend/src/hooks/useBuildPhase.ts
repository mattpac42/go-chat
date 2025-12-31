'use client';

import { useState, useEffect, useMemo, useCallback, useRef } from 'react';
import { Message } from '@/types';
import {
  BuildPhase,
  detectCurrentPhase,
  groupMessagesByPhase,
  getPhaseIndex,
} from '@/components/chat/BuildPhaseProgress';

interface UseBuildPhaseOptions {
  messages: Message[];
  isDiscoveryMode: boolean;
  discoveryComplete: boolean;
}

interface UseBuildPhaseReturn {
  /** The current detected phase */
  currentPhase: BuildPhase;
  /** Messages grouped by phase */
  messagesByPhase: Map<BuildPhase, Message[]>;
  /** Array of phases that have been completed */
  completedPhases: BuildPhase[];
  /** Phase that just completed (for toast) */
  newlyCompletedPhase: BuildPhase | null;
  /** Acknowledge a completed phase (dismiss toast) */
  acknowledgePhase: (phase: BuildPhase) => void;
  /** Whether to show phased view (grouped messages) */
  showPhasedView: boolean;
  /** Toggle phased view */
  togglePhasedView: () => void;
  /** Scroll to a specific phase */
  scrollToPhase: (phase: BuildPhase) => void;
}

/**
 * Hook to track build phases based on message content and discovery state
 */
export function useBuildPhase({
  messages,
  isDiscoveryMode,
  discoveryComplete,
}: UseBuildPhaseOptions): UseBuildPhaseReturn {
  const [acknowledgedPhases, setAcknowledgedPhases] = useState<Set<BuildPhase>>(
    new Set()
  );
  const [showPhasedView, setShowPhasedView] = useState(false);
  const previousPhaseRef = useRef<BuildPhase>('discovery');

  // Track the message count when discovery was completed
  const discoveryCompleteAtRef = useRef<number>(-1);

  // Update discovery complete index when discovery transitions to complete
  useEffect(() => {
    if (discoveryComplete && discoveryCompleteAtRef.current === -1) {
      // Discovery just completed - remember the current message index
      discoveryCompleteAtRef.current = messages.length - 1;
    } else if (!discoveryComplete) {
      // Discovery was reset
      discoveryCompleteAtRef.current = -1;
    }
  }, [discoveryComplete, messages.length]);

  // Calculate discovery complete index
  const discoveryCompleteIndex = useMemo(() => {
    if (!discoveryComplete) return -1;
    // If we've tracked when discovery completed, use that
    if (discoveryCompleteAtRef.current >= 0) {
      return discoveryCompleteAtRef.current;
    }
    // Fallback: if discoveryComplete is true but we haven't tracked the index,
    // all current messages are discovery (edge case on initial load)
    return messages.length - 1;
  }, [discoveryComplete, messages.length]);

  // Determine current phase
  const currentPhase = useMemo(() => {
    if (isDiscoveryMode && !discoveryComplete) {
      return 'discovery';
    }
    return detectCurrentPhase(messages, discoveryCompleteIndex);
  }, [messages, isDiscoveryMode, discoveryComplete, discoveryCompleteIndex]);

  // Group messages by phase
  const messagesByPhase = useMemo(() => {
    return groupMessagesByPhase(messages, discoveryCompleteIndex);
  }, [messages, discoveryCompleteIndex]);

  // Determine completed phases
  const completedPhases = useMemo(() => {
    const currentIndex = getPhaseIndex(currentPhase);
    const phases: BuildPhase[] = ['discovery', 'planning', 'building', 'testing', 'launch'];
    return phases.slice(0, currentIndex);
  }, [currentPhase]);

  // Detect newly completed phase
  const newlyCompletedPhase = useMemo(() => {
    const previousIndex = getPhaseIndex(previousPhaseRef.current);
    const currentIndex = getPhaseIndex(currentPhase);

    // If we moved forward, the previous phase just completed
    if (currentIndex > previousIndex && !acknowledgedPhases.has(previousPhaseRef.current)) {
      return previousPhaseRef.current;
    }

    // Special case: discovery complete
    if (discoveryComplete && !acknowledgedPhases.has('discovery')) {
      return 'discovery';
    }

    return null;
  }, [currentPhase, discoveryComplete, acknowledgedPhases]);

  // Update previous phase ref
  useEffect(() => {
    previousPhaseRef.current = currentPhase;
  }, [currentPhase]);

  // Acknowledge a completed phase (dismiss the toast)
  const acknowledgePhase = useCallback((phase: BuildPhase) => {
    setAcknowledgedPhases((prev) => {
      const newSet = new Set(prev);
      newSet.add(phase);
      return newSet;
    });
  }, []);

  // Toggle phased view
  const togglePhasedView = useCallback(() => {
    setShowPhasedView((prev) => !prev);
  }, []);

  // Scroll to a specific phase section
  const scrollToPhase = useCallback((phase: BuildPhase) => {
    const element = document.getElementById(`phase-${phase}-messages`);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  }, []);

  return {
    currentPhase,
    messagesByPhase,
    completedPhases,
    newlyCompletedPhase,
    acknowledgePhase,
    showPhasedView,
    togglePhasedView,
    scrollToPhase,
  };
}
