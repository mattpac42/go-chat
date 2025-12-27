'use client';

import { ZoomLevel, ZOOM_LEVELS } from '@/hooks/useCodeZoom';

interface ZoomControlsProps {
  zoomLevel: ZoomLevel;
  onZoomChange: (level: ZoomLevel) => void;
  className?: string;
}

/**
 * Compact segmented button for zoom control.
 * Displays preset zoom levels (75%, 100%, 125%, 150%) with visual indication of active level.
 * Matches existing toolbar styling with gray background and small text.
 */
export function ZoomControls({ zoomLevel, onZoomChange, className = '' }: ZoomControlsProps) {
  return (
    <div
      className={`inline-flex items-center rounded-md bg-gray-100 p-0.5 ${className}`}
      role="group"
      aria-label="Code zoom level"
    >
      {ZOOM_LEVELS.map((level) => {
        const isActive = level === zoomLevel;
        return (
          <button
            key={level}
            onClick={(e) => {
              e.stopPropagation();
              onZoomChange(level);
            }}
            className={`
              px-2 py-0.5 text-xs font-medium rounded transition-colors
              ${isActive
                ? 'bg-white text-gray-900 shadow-sm'
                : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
              }
            `}
            aria-label={`Zoom ${level}%`}
            aria-pressed={isActive}
          >
            {level}%
          </button>
        );
      })}
    </div>
  );
}

export default ZoomControls;
