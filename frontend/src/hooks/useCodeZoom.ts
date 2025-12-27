'use client';

import { useState, useCallback, useEffect } from 'react';

const STORAGE_KEY = 'codeZoomLevel';
const ZOOM_LEVELS = [75, 100, 125, 150] as const;
type ZoomLevel = (typeof ZOOM_LEVELS)[number];

interface UseCodeZoomReturn {
  zoomLevel: ZoomLevel;
  setZoom: (level: ZoomLevel) => void;
  zoomIn: () => void;
  zoomOut: () => void;
  zoomLevels: readonly number[];
  getZoomStyle: () => React.CSSProperties;
}

/**
 * Hook for managing code zoom level across the application.
 * - Persists zoom preference to localStorage
 * - Supports keyboard shortcuts (Cmd/Ctrl + Plus/Minus)
 * - Provides zoom style object for applying to code elements
 */
export function useCodeZoom(): UseCodeZoomReturn {
  const [zoomLevel, setZoomLevel] = useState<ZoomLevel>(100);

  // Load zoom level from localStorage on mount
  useEffect(() => {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) {
      const parsed = parseInt(stored, 10);
      if (ZOOM_LEVELS.includes(parsed as ZoomLevel)) {
        setZoomLevel(parsed as ZoomLevel);
      }
    }
  }, []);

  // Save zoom level to localStorage when it changes
  const setZoom = useCallback((level: ZoomLevel) => {
    setZoomLevel(level);
    localStorage.setItem(STORAGE_KEY, String(level));
  }, []);

  // Zoom in to the next level
  const zoomIn = useCallback(() => {
    setZoomLevel((current) => {
      const currentIndex = ZOOM_LEVELS.indexOf(current);
      const nextIndex = Math.min(currentIndex + 1, ZOOM_LEVELS.length - 1);
      const newLevel = ZOOM_LEVELS[nextIndex];
      localStorage.setItem(STORAGE_KEY, String(newLevel));
      return newLevel;
    });
  }, []);

  // Zoom out to the previous level
  const zoomOut = useCallback(() => {
    setZoomLevel((current) => {
      const currentIndex = ZOOM_LEVELS.indexOf(current);
      const prevIndex = Math.max(currentIndex - 1, 0);
      const newLevel = ZOOM_LEVELS[prevIndex];
      localStorage.setItem(STORAGE_KEY, String(newLevel));
      return newLevel;
    });
  }, []);

  // Keyboard shortcut handler
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const isMod = e.metaKey || e.ctrlKey;

      if (isMod && (e.key === '=' || e.key === '+')) {
        e.preventDefault();
        zoomIn();
      } else if (isMod && e.key === '-') {
        e.preventDefault();
        zoomOut();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [zoomIn, zoomOut]);

  // Generate style object for zoom
  const getZoomStyle = useCallback((): React.CSSProperties => {
    return {
      fontSize: `${zoomLevel}%`,
      lineHeight: 1.6,
    };
  }, [zoomLevel]);

  return {
    zoomLevel,
    setZoom,
    zoomIn,
    zoomOut,
    zoomLevels: ZOOM_LEVELS,
    getZoomStyle,
  };
}

export { ZOOM_LEVELS };
export type { ZoomLevel };
