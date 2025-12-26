'use client';

import { useState, useCallback, useEffect } from 'react';
import { API_BASE_URL } from '@/lib/api';
import { UserProgress, UserAchievement, Nudge, LearningEvent } from '@/types/achievements';

interface UseAchievementsReturn {
  progress: UserProgress | null;
  achievements: UserAchievement[];
  unseenAchievements: UserAchievement[];
  currentNudge: Nudge | null;
  isLoading: boolean;
  error: string | null;
  recordEvent: (event: LearningEvent) => Promise<UserAchievement[]>;
  markAchievementSeen: (id: string) => Promise<void>;
  dismissNudge: () => Promise<void>;
  acceptNudge: () => Promise<void>;
  refetch: () => Promise<void>;
}

/**
 * Hook for managing achievements and learning progress
 * @param projectId - The project ID to track progress for
 */
export function useAchievements(projectId: string | null): UseAchievementsReturn {
  const [progress, setProgress] = useState<UserProgress | null>(null);
  const [achievements, setAchievements] = useState<UserAchievement[]>([]);
  const [unseenAchievements, setUnseenAchievements] = useState<UserAchievement[]>([]);
  const [currentNudge, setCurrentNudge] = useState<Nudge | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  /**
   * Fetch user progress from API
   */
  const fetchProgress = useCallback(async () => {
    if (!projectId) return;

    try {
      const response = await fetch(`${API_BASE_URL}/api/projects/${projectId}/progress`);
      if (response.ok) {
        const data = await response.json();
        setProgress(data);
      }
    } catch (err) {
      console.error('Failed to fetch progress:', err);
    }
  }, [projectId]);

  /**
   * Fetch achievements from API
   */
  const fetchAchievements = useCallback(async () => {
    if (!projectId) return;

    try {
      const [allRes, unseenRes] = await Promise.all([
        fetch(`${API_BASE_URL}/api/projects/${projectId}/achievements`),
        fetch(`${API_BASE_URL}/api/projects/${projectId}/achievements/unseen`)
      ]);

      if (allRes.ok) {
        const data = await allRes.json();
        setAchievements(data.achievements || []);
      }

      if (unseenRes.ok) {
        const data = await unseenRes.json();
        setUnseenAchievements(data.achievements || []);
      }
    } catch (err) {
      console.error('Failed to fetch achievements:', err);
    }
  }, [projectId]);

  /**
   * Fetch current nudge from API
   */
  const fetchNudge = useCallback(async () => {
    if (!projectId) return;

    try {
      const response = await fetch(`${API_BASE_URL}/api/projects/${projectId}/nudge`);
      if (response.ok) {
        const data = await response.json();
        setCurrentNudge(data);
      }
    } catch (err) {
      console.error('Failed to fetch nudge:', err);
    }
  }, [projectId]);

  /**
   * Record a learning event and check for newly unlocked achievements
   */
  const recordEvent = useCallback(async (event: LearningEvent): Promise<UserAchievement[]> => {
    if (!projectId) return [];

    try {
      const response = await fetch(`${API_BASE_URL}/api/projects/${projectId}/events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(event)
      });

      if (response.ok) {
        const data = await response.json();

        // Refresh progress and achievements
        await Promise.all([fetchProgress(), fetchAchievements()]);

        return data.unlocked || [];
      }
    } catch (err) {
      console.error('Failed to record event:', err);
    }

    return [];
  }, [projectId, fetchProgress, fetchAchievements]);

  /**
   * Mark an achievement as seen
   */
  const markAchievementSeen = useCallback(async (id: string) => {
    if (!projectId) return;

    try {
      await fetch(`${API_BASE_URL}/api/projects/${projectId}/achievements/${id}/seen`, {
        method: 'POST'
      });

      setUnseenAchievements(prev => prev.filter(a => a.id !== id));
    } catch (err) {
      console.error('Failed to mark achievement seen:', err);
    }
  }, [projectId]);

  /**
   * Dismiss the current nudge
   */
  const dismissNudge = useCallback(async () => {
    if (!projectId) return;

    try {
      await fetch(`${API_BASE_URL}/api/projects/${projectId}/nudge/dismiss`, {
        method: 'POST'
      });

      setCurrentNudge(null);
    } catch (err) {
      console.error('Failed to dismiss nudge:', err);
    }
  }, [projectId]);

  /**
   * Accept the current nudge and perform its action
   */
  const acceptNudge = useCallback(async () => {
    if (!projectId) return;

    try {
      await fetch(`${API_BASE_URL}/api/projects/${projectId}/nudge/accept`, {
        method: 'POST'
      });

      setCurrentNudge(null);
      // Refetch to get next nudge if any
      await fetchNudge();
    } catch (err) {
      console.error('Failed to accept nudge:', err);
    }
  }, [projectId, fetchNudge]);

  /**
   * Refetch all data
   */
  const refetch = useCallback(async () => {
    await Promise.all([fetchProgress(), fetchAchievements(), fetchNudge()]);
  }, [fetchProgress, fetchAchievements, fetchNudge]);

  // Initial fetch
  useEffect(() => {
    if (!projectId) {
      setIsLoading(false);
      setProgress(null);
      setAchievements([]);
      setUnseenAchievements([]);
      setCurrentNudge(null);
      return;
    }

    setIsLoading(true);
    setError(null);

    Promise.all([fetchProgress(), fetchAchievements(), fetchNudge()])
      .catch(err => {
        const errorMessage = err instanceof Error
          ? err.message
          : 'Failed to load achievements';
        setError(errorMessage);
        console.error('Failed to fetch achievements data:', err);
      })
      .finally(() => setIsLoading(false));
  }, [projectId, fetchProgress, fetchAchievements, fetchNudge]);

  return {
    progress,
    achievements,
    unseenAchievements,
    currentNudge,
    isLoading,
    error,
    recordEvent,
    markAchievementSeen,
    dismissNudge,
    acceptNudge,
    refetch
  };
}
