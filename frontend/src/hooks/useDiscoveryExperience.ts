'use client';

import { useState, useCallback, useEffect } from 'react';

const DISCOVERY_COMPLETED_KEY = 'hasCompletedDiscovery';

/**
 * Hook to track if user has completed discovery before (across all projects)
 * Used to show "skip to building" option for returning users
 */
export function useDiscoveryExperience() {
  const [hasCompletedBefore, setHasCompletedBefore] = useState(false);
  const [isLoaded, setIsLoaded] = useState(false);

  // Load from localStorage on mount
  useEffect(() => {
    if (typeof window !== 'undefined') {
      const stored = localStorage.getItem(DISCOVERY_COMPLETED_KEY);
      setHasCompletedBefore(stored === 'true');
      setIsLoaded(true);
    }
  }, []);

  /**
   * Mark discovery as completed (called when any discovery finishes successfully)
   */
  const markDiscoveryCompleted = useCallback(() => {
    if (typeof window !== 'undefined') {
      localStorage.setItem(DISCOVERY_COMPLETED_KEY, 'true');
      setHasCompletedBefore(true);
    }
  }, []);

  /**
   * Clear the "has completed" flag (for testing/reset purposes)
   */
  const clearDiscoveryHistory = useCallback(() => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem(DISCOVERY_COMPLETED_KEY);
      setHasCompletedBefore(false);
    }
  }, []);

  return {
    hasCompletedBefore,
    isLoaded,
    markDiscoveryCompleted,
    clearDiscoveryHistory,
  };
}
