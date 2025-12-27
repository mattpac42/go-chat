'use client';

import { useState, useCallback, useEffect, useRef } from 'react';

const MAX_HISTORY_SIZE = 50;
const STORAGE_KEY_PREFIX = 'message-history-';

/**
 * Hook for managing chat message history with up/down arrow navigation
 * Similar to terminal command history behavior
 */
export function useMessageHistory(projectId: string) {
  const [history, setHistory] = useState<string[]>([]);
  const navigationIndexRef = useRef<number>(-1); // -1 means at draft position
  const draftRef = useRef<string>('');
  const storageKey = `${STORAGE_KEY_PREFIX}${projectId}`;

  // Load history from sessionStorage on mount
  useEffect(() => {
    if (typeof window === 'undefined') return;

    try {
      const stored = sessionStorage.getItem(storageKey);
      if (stored) {
        const parsed = JSON.parse(stored);
        if (Array.isArray(parsed)) {
          setHistory(parsed);
        }
      }
    } catch {
      // Handle corrupted storage gracefully - start with empty history
      setHistory([]);
    }
  }, [storageKey]);

  // Save history to sessionStorage whenever it changes
  const saveToStorage = useCallback(
    (newHistory: string[]) => {
      if (typeof window === 'undefined') return;
      try {
        sessionStorage.setItem(storageKey, JSON.stringify(newHistory));
      } catch {
        // Ignore storage errors (e.g., quota exceeded)
      }
    },
    [storageKey]
  );

  /**
   * Add a message to history
   * Does not add empty messages or duplicate consecutive messages
   */
  const addToHistory = useCallback(
    (message: string) => {
      const trimmed = message.trim();
      if (!trimmed) return;

      setHistory((prev) => {
        // Don't add duplicate consecutive messages
        if (prev.length > 0 && prev[prev.length - 1] === trimmed) {
          return prev;
        }

        // Add new message and limit size
        const newHistory = [...prev, trimmed];
        if (newHistory.length > MAX_HISTORY_SIZE) {
          // Remove oldest messages
          const trimmedHistory = newHistory.slice(-MAX_HISTORY_SIZE);
          saveToStorage(trimmedHistory);
          return trimmedHistory;
        }

        saveToStorage(newHistory);
        return newHistory;
      });

      // Reset navigation state when a new message is sent
      navigationIndexRef.current = -1;
      draftRef.current = '';
    },
    [saveToStorage]
  );

  /**
   * Navigate through message history
   * @param direction 'up' to go back in history, 'down' to go forward
   * @param currentValue The current input value (to preserve as draft)
   * @returns The message to display, or null if navigation is not possible
   */
  const navigateHistory = useCallback(
    (direction: 'up' | 'down', currentValue: string): string | null => {
      const currentIndex = navigationIndexRef.current;

      if (direction === 'up') {
        // Going back in history
        if (history.length === 0) return null;

        // If we're at draft position (-1), save the draft and go to most recent
        if (currentIndex === -1) {
          draftRef.current = currentValue;
          navigationIndexRef.current = history.length - 1;
          return history[history.length - 1];
        }

        // If we're at the oldest message, can't go further back
        if (currentIndex === 0) {
          return null;
        }

        // Go to older message
        navigationIndexRef.current = currentIndex - 1;
        return history[currentIndex - 1];
      } else {
        // Going forward in history (toward draft)
        // If we're at draft position, can't go further
        if (currentIndex === -1) {
          return null;
        }

        // If we're at the most recent message, go back to draft
        if (currentIndex === history.length - 1) {
          navigationIndexRef.current = -1;
          return draftRef.current;
        }

        // Go to newer message
        navigationIndexRef.current = currentIndex + 1;
        return history[currentIndex + 1];
      }
    },
    [history]
  );

  /**
   * Reset navigation to default state (at draft position)
   * Call this when user types to cancel history navigation
   */
  const resetNavigation = useCallback(() => {
    navigationIndexRef.current = -1;
    draftRef.current = '';
  }, []);

  return {
    addToHistory,
    navigateHistory,
    resetNavigation,
  };
}
