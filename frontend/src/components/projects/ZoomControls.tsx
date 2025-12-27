'use client';

import { useCallback } from 'react';
import { ZoomLevel, ZOOM_LEVELS } from '@/hooks/useCodeZoom';

interface ZoomControlsProps {
  zoomLevel: ZoomLevel;
  onZoomChange: (level: ZoomLevel) => void;
  className?: string;
}

/**
 * Compact slider control for zoom level.
 * Snaps to preset zoom levels (25%, 50%, 75%, 100%, 125%, 150%).
 * Matches existing toolbar styling with dark aesthetic.
 */
export function ZoomControls({ zoomLevel, onZoomChange, className = '' }: ZoomControlsProps) {
  // Find the current index in zoom levels array
  const currentIndex = ZOOM_LEVELS.indexOf(zoomLevel);
  const minIndex = 0;
  const maxIndex = ZOOM_LEVELS.length - 1;

  const handleSliderChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      e.stopPropagation();
      const index = parseInt(e.target.value, 10);
      const newLevel = ZOOM_LEVELS[index];
      if (newLevel !== undefined) {
        onZoomChange(newLevel);
      }
    },
    [onZoomChange]
  );

  return (
    <div
      className={`inline-flex items-center gap-2 ${className}`}
      role="group"
      aria-label="Code zoom level"
      onClick={(e) => e.stopPropagation()}
    >
      <input
        type="range"
        min={minIndex}
        max={maxIndex}
        value={currentIndex}
        onChange={handleSliderChange}
        className="w-20 h-1.5 bg-gray-200 rounded-full appearance-none cursor-pointer
          [&::-webkit-slider-thumb]:appearance-none
          [&::-webkit-slider-thumb]:w-3
          [&::-webkit-slider-thumb]:h-3
          [&::-webkit-slider-thumb]:rounded-full
          [&::-webkit-slider-thumb]:bg-gray-600
          [&::-webkit-slider-thumb]:hover:bg-gray-700
          [&::-webkit-slider-thumb]:transition-colors
          [&::-moz-range-thumb]:w-3
          [&::-moz-range-thumb]:h-3
          [&::-moz-range-thumb]:rounded-full
          [&::-moz-range-thumb]:bg-gray-600
          [&::-moz-range-thumb]:hover:bg-gray-700
          [&::-moz-range-thumb]:border-0
          [&::-moz-range-thumb]:transition-colors"
        aria-label="Zoom level slider"
        aria-valuemin={ZOOM_LEVELS[minIndex]}
        aria-valuemax={ZOOM_LEVELS[maxIndex]}
        aria-valuenow={zoomLevel}
        aria-valuetext={`${zoomLevel}%`}
      />
      <span className="text-xs font-medium text-gray-600 min-w-[36px] text-right">
        {zoomLevel}%
      </span>
    </div>
  );
}

export default ZoomControls;
