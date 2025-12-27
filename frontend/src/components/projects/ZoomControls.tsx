'use client';

import { useCallback } from 'react';
import { ZoomLevel, ZOOM_LEVELS } from '@/hooks/useCodeZoom';

interface ZoomControlsProps {
  zoomLevel: ZoomLevel;
  onZoomChange: (level: ZoomLevel) => void;
  className?: string;
}

function MinusIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
    </svg>
  );
}

function PlusIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
    </svg>
  );
}

/**
 * Compact +/- button controls for zoom level.
 * Snaps to preset zoom levels (25%, 50%, 75%, 100%, 125%, 150%).
 * Each click adjusts zoom by one step (approximately 25%).
 */
export function ZoomControls({ zoomLevel, onZoomChange, className = '' }: ZoomControlsProps) {
  const currentIndex = ZOOM_LEVELS.indexOf(zoomLevel);
  const canZoomOut = currentIndex > 0;
  const canZoomIn = currentIndex < ZOOM_LEVELS.length - 1;

  const handleZoomOut = useCallback(
    (e: React.MouseEvent) => {
      e.stopPropagation();
      if (canZoomOut) {
        const newLevel = ZOOM_LEVELS[currentIndex - 1];
        if (newLevel !== undefined) {
          onZoomChange(newLevel);
        }
      }
    },
    [canZoomOut, currentIndex, onZoomChange]
  );

  const handleZoomIn = useCallback(
    (e: React.MouseEvent) => {
      e.stopPropagation();
      if (canZoomIn) {
        const newLevel = ZOOM_LEVELS[currentIndex + 1];
        if (newLevel !== undefined) {
          onZoomChange(newLevel);
        }
      }
    },
    [canZoomIn, currentIndex, onZoomChange]
  );

  return (
    <div
      className={`inline-flex items-center gap-1 ${className}`}
      role="group"
      aria-label="Code zoom level"
      onClick={(e) => e.stopPropagation()}
    >
      <button
        onClick={handleZoomOut}
        disabled={!canZoomOut}
        className="flex items-center justify-center w-6 h-6 text-gray-600 hover:text-gray-800 hover:bg-gray-200 rounded transition-colors disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:bg-transparent"
        aria-label="Zoom out"
        title="Zoom out"
      >
        <MinusIcon className="w-3.5 h-3.5" />
      </button>
      <span className="text-xs font-medium text-gray-600 min-w-[36px] text-center tabular-nums">
        {zoomLevel}%
      </span>
      <button
        onClick={handleZoomIn}
        disabled={!canZoomIn}
        className="flex items-center justify-center w-6 h-6 text-gray-600 hover:text-gray-800 hover:bg-gray-200 rounded transition-colors disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:bg-transparent"
        aria-label="Zoom in"
        title="Zoom in"
      >
        <PlusIcon className="w-3.5 h-3.5" />
      </button>
    </div>
  );
}

export default ZoomControls;
