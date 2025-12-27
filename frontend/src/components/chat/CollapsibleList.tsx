'use client';

import { useState, useRef, useEffect, ReactNode, Children, isValidElement } from 'react';

interface CollapsibleListProps {
  children: ReactNode;
  listType: 'ul' | 'ol';
  className?: string;
  visibleCount?: number;
  threshold?: number;
}

/**
 * A list component that auto-collapses when items exceed threshold.
 * Shows first N items with a "Show X more items..." link to expand.
 */
export function CollapsibleList({
  children,
  listType,
  className = '',
  visibleCount = 4,
  threshold = 5,
}: CollapsibleListProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const [collapsedHeight, setCollapsedHeight] = useState<number | null>(null);
  const [fullHeight, setFullHeight] = useState<number | null>(null);
  const listRef = useRef<HTMLUListElement | HTMLOListElement>(null);
  const measureRef = useRef<HTMLDivElement>(null);

  // Convert children to array for counting and slicing
  const childArray = Children.toArray(children);
  const totalItems = childArray.length;
  const shouldCollapse = totalItems > threshold;
  const hiddenCount = totalItems - visibleCount;

  // Measure heights for smooth animation
  useEffect(() => {
    if (!shouldCollapse || !measureRef.current) return;

    // Measure full height
    const measureDiv = measureRef.current;
    measureDiv.style.position = 'absolute';
    measureDiv.style.visibility = 'hidden';
    measureDiv.style.height = 'auto';
    setFullHeight(measureDiv.scrollHeight);

    // Measure collapsed height (first N items)
    const visibleItems = measureDiv.querySelectorAll(':scope > li');
    let height = 0;
    visibleItems.forEach((item, index) => {
      if (index < visibleCount) {
        height += (item as HTMLElement).offsetHeight;
      }
    });
    // Add some padding for the list
    setCollapsedHeight(height + 8);
  }, [children, shouldCollapse, visibleCount]);

  // If not enough items, render normally
  if (!shouldCollapse) {
    const ListTag = listType;
    return <ListTag className={className}>{children}</ListTag>;
  }

  const ListTag = listType;

  return (
    <div className="relative">
      {/* Hidden measuring element */}
      <div
        ref={measureRef}
        aria-hidden="true"
        className="absolute invisible pointer-events-none"
        style={{ width: '100%' }}
      >
        <ListTag className={className}>{children}</ListTag>
      </div>

      {/* Visible list with animation */}
      <div
        className="overflow-hidden transition-all duration-300 ease-in-out"
        style={{
          maxHeight: isExpanded
            ? (fullHeight ?? 'none')
            : (collapsedHeight ?? 'auto'),
        }}
      >
        <ListTag ref={listRef as any} className={className}>
          {children}
        </ListTag>
      </div>

      {/* Show more/less link */}
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="mt-1 text-sm text-gray-500 hover:text-gray-700 transition-colors cursor-pointer"
        aria-expanded={isExpanded}
        aria-label={isExpanded ? 'Show fewer items' : `Show ${hiddenCount} more items`}
      >
        {isExpanded ? (
          <span className="flex items-center gap-1">
            <ChevronUpIcon className="w-3 h-3" />
            Show less
          </span>
        ) : (
          <span className="flex items-center gap-1">
            <ChevronDownIcon className="w-3 h-3" />
            Show {hiddenCount} more {hiddenCount === 1 ? 'item' : 'items'}...
          </span>
        )}
      </button>
    </div>
  );
}

function ChevronDownIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M19 9l-7 7-7-7"
      />
    </svg>
  );
}

function ChevronUpIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M5 15l7-7 7 7"
      />
    </svg>
  );
}
