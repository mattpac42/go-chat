'use client';

import { useState, useCallback, useEffect, useMemo } from 'react';
import {
  ProjectDiscovery,
  DiscoveryStage,
  DiscoverySummary,
} from '@/types/discovery';
import { API_BASE_URL } from '@/lib/api';

/**
 * Maps discovery stages to progress numbers (1-5)
 */
const STAGE_PROGRESS: Record<DiscoveryStage, number> = {
  welcome: 1,
  problem: 2,
  personas: 3,
  mvp: 4,
  summary: 5,
  complete: 5,
};

interface UseDiscoveryReturn {
  /** Current discovery state from API */
  discovery: ProjectDiscovery | null;
  /** Whether the project is in discovery mode (not yet complete) */
  isDiscoveryMode: boolean;
  /** Current discovery stage */
  currentStage: DiscoveryStage;
  /** Stage progress (1-5) */
  stageProgress: number;
  /** Discovery summary (available after summary stage) */
  summary: DiscoverySummary | null;
  /** Loading state */
  isLoading: boolean;
  /** Error message if any */
  error: string | null;
  /** Confirm and complete the discovery process */
  confirmDiscovery: () => Promise<void>;
  /** Reset discovery to start over */
  resetDiscovery: () => Promise<void>;
  /** Refetch discovery state from API */
  refetch: () => Promise<void>;
}

/**
 * Hook for managing project discovery state
 * @param projectId - The project ID to fetch discovery for
 */
export function useDiscovery(projectId: string | null): UseDiscoveryReturn {
  const [discovery, setDiscovery] = useState<ProjectDiscovery | null>(null);
  const [summary, setSummary] = useState<DiscoverySummary | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  /**
   * Fetch discovery state from API
   */
  const fetchDiscovery = useCallback(async () => {
    if (!projectId) {
      setIsLoading(false);
      setDiscovery(null);
      setSummary(null);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(`${API_BASE_URL}/api/projects/${projectId}/discovery`);

      if (response.status === 404) {
        // No discovery yet - this is normal for new projects
        setDiscovery(null);
        setSummary(null);
        setIsLoading(false);
        return;
      }

      if (!response.ok) {
        throw new Error(`Failed to fetch discovery: ${response.statusText}`);
      }

      const data = await response.json();
      setDiscovery(data.discovery || data);

      // Set summary if available
      if (data.summary) {
        setSummary(data.summary);
      }
    } catch (err) {
      const errorMessage = err instanceof Error
        ? err.message
        : 'Failed to load discovery state';
      setError(errorMessage);
      console.error('Failed to fetch discovery:', err);
    } finally {
      setIsLoading(false);
    }
  }, [projectId]);

  /**
   * Confirm and complete the discovery process
   */
  const confirmDiscovery = useCallback(async () => {
    if (!projectId) {
      setError('No project ID provided');
      return;
    }

    setError(null);

    try {
      const response = await fetch(`${API_BASE_URL}/api/projects/${projectId}/discovery/confirm`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to confirm discovery: ${response.statusText}`);
      }

      const data = await response.json();
      setDiscovery(data.discovery || data);

      if (data.summary) {
        setSummary(data.summary);
      }
    } catch (err) {
      const errorMessage = err instanceof Error
        ? err.message
        : 'Failed to confirm discovery';
      setError(errorMessage);
      console.error('Failed to confirm discovery:', err);
      throw err;
    }
  }, [projectId]);

  /**
   * Reset discovery to start over
   */
  const resetDiscovery = useCallback(async () => {
    if (!projectId) {
      setError('No project ID provided');
      return;
    }

    setError(null);

    try {
      const response = await fetch(`${API_BASE_URL}/api/projects/${projectId}/discovery`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error(`Failed to reset discovery: ${response.statusText}`);
      }

      // Clear local state
      setDiscovery(null);
      setSummary(null);
    } catch (err) {
      const errorMessage = err instanceof Error
        ? err.message
        : 'Failed to reset discovery';
      setError(errorMessage);
      console.error('Failed to reset discovery:', err);
      throw err;
    }
  }, [projectId]);

  // Derived state: is in discovery mode (not complete)
  const isDiscoveryMode = useMemo(() => {
    if (!discovery) return true; // New projects start in discovery
    return discovery.stage !== 'complete';
  }, [discovery]);

  // Derived state: current stage
  const currentStage = useMemo((): DiscoveryStage => {
    if (!discovery) return 'welcome';
    return discovery.stage;
  }, [discovery]);

  // Derived state: stage progress (1-5)
  const stageProgress = useMemo(() => {
    return STAGE_PROGRESS[currentStage];
  }, [currentStage]);

  // Fetch discovery on mount or when projectId changes
  useEffect(() => {
    fetchDiscovery();
  }, [fetchDiscovery]);

  return {
    discovery,
    isDiscoveryMode,
    currentStage,
    stageProgress,
    summary,
    isLoading,
    error,
    confirmDiscovery,
    resetDiscovery,
    refetch: fetchDiscovery,
  };
}
