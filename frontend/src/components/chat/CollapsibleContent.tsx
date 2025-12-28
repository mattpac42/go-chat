'use client';

import { useState, useRef, useEffect, ReactNode, Children, isValidElement } from 'react';

interface CollapsibleContentProps {
  children: ReactNode;
  paragraphThreshold?: number;
  visibleParagraphs?: number;
  isUserMessage?: boolean;
  /** When true, disables collapsing entirely (useful during streaming) */
  disabled?: boolean;
}

/**
 * Wrapper component that auto-collapses content when paragraphs exceed threshold.
 * Shows first N paragraphs with a "Show more..." link to expand.
 */
export function CollapsibleContent({
  children,
  paragraphThreshold = 3,
  visibleParagraphs = 2,
  isUserMessage = false,
  disabled = false,
}: CollapsibleContentProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const [collapsedHeight, setCollapsedHeight] = useState<number | null>(null);
  const [fullHeight, setFullHeight] = useState<number | null>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const measureRef = useRef<HTMLDivElement>(null);

  // Count paragraph-level elements (p, h1-h6, pre, ul, ol, blockquote)
  const childArray = Children.toArray(children);
  const paragraphCount = countParagraphElements(childArray);
  const shouldCollapse = !disabled && paragraphCount > paragraphThreshold;

  // Measure heights for smooth animation
  useEffect(() => {
    if (!shouldCollapse || !measureRef.current || !contentRef.current) return;

    // Measure full height
    setFullHeight(measureRef.current.scrollHeight);

    // Measure collapsed height based on visible paragraphs
    const contentDiv = contentRef.current;
    const blockElements = contentDiv.querySelectorAll(':scope > p, :scope > h1, :scope > h2, :scope > h3, :scope > h4, :scope > h5, :scope > h6, :scope > pre, :scope > ul, :scope > ol, :scope > blockquote, :scope > div');

    let height = 0;
    let count = 0;
    blockElements.forEach((el) => {
      if (count < visibleParagraphs) {
        const style = window.getComputedStyle(el);
        const marginTop = parseFloat(style.marginTop) || 0;
        const marginBottom = parseFloat(style.marginBottom) || 0;
        height += (el as HTMLElement).offsetHeight + marginTop + marginBottom;
        count++;
      }
    });

    // Add a small buffer
    setCollapsedHeight(Math.max(height, 60));
  }, [children, shouldCollapse, visibleParagraphs]);

  // If not enough content, render normally
  if (!shouldCollapse) {
    return <>{children}</>;
  }

  const hiddenCount = paragraphCount - visibleParagraphs;

  return (
    <div className="relative">
      {/* Hidden measuring element */}
      <div
        ref={measureRef}
        aria-hidden="true"
        className="absolute invisible pointer-events-none"
        style={{ width: '100%' }}
      >
        {children}
      </div>

      {/* Visible content with animation */}
      <div
        ref={contentRef}
        className="overflow-hidden transition-all duration-300 ease-in-out"
        style={{
          maxHeight: isExpanded
            ? (fullHeight ?? 'none')
            : (collapsedHeight ?? 'auto'),
        }}
      >
        {children}
      </div>

      {/* Gradient fade when collapsed */}
      {!isExpanded && (
        <div
          className={`absolute bottom-6 left-0 right-0 h-8 pointer-events-none ${
            isUserMessage
              ? 'bg-gradient-to-t from-teal-400 to-transparent'
              : 'bg-gradient-to-t from-gray-100 to-transparent'
          }`}
        />
      )}

      {/* Show more/less link */}
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className={`mt-1 text-sm transition-colors cursor-pointer ${
          isUserMessage
            ? 'text-teal-100 hover:text-white'
            : 'text-gray-500 hover:text-gray-700'
        }`}
        aria-expanded={isExpanded}
        aria-label={isExpanded ? 'Show less content' : 'Show more content'}
      >
        {isExpanded ? (
          <span className="flex items-center gap-1">
            <ChevronUpIcon className="w-3 h-3" />
            Show less
          </span>
        ) : (
          <span className="flex items-center gap-1">
            <ChevronDownIcon className="w-3 h-3" />
            Show more...
          </span>
        )}
      </button>
    </div>
  );
}

/**
 * Count paragraph-level elements in React children
 */
function countParagraphElements(children: ReturnType<typeof Children.toArray>): number {
  let count = 0;

  children.forEach((child) => {
    if (isValidElement(child)) {
      const type = child.type;
      // Check for paragraph-level element types
      if (
        type === 'p' ||
        type === 'h1' ||
        type === 'h2' ||
        type === 'h3' ||
        type === 'h4' ||
        type === 'h5' ||
        type === 'h6' ||
        type === 'pre' ||
        type === 'ul' ||
        type === 'ol' ||
        type === 'blockquote' ||
        type === 'div'
      ) {
        count++;
      }
      // Also check for custom components that render these elements
      // by looking at displayName or checking if it's a function component
      if (typeof type === 'function') {
        // Assume function components that wrap block elements count as 1
        count++;
      }
    }
  });

  return count;
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
