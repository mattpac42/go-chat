'use client';

import { useState, useCallback, useRef, useEffect } from 'react';

interface CodeBlockProps {
  code: string;
  language?: string;
  defaultCollapsed?: boolean;
}

function ChevronIcon({ isOpen, className }: { isOpen: boolean; className?: string }) {
  return (
    <svg
      className={`w-4 h-4 transition-transform duration-200 ${isOpen ? 'rotate-90' : ''} ${className || ''}`}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9 5l7 7-7 7"
      />
    </svg>
  );
}

export function CodeBlock({ code, language = 'text', defaultCollapsed = true }: CodeBlockProps) {
  const [copied, setCopied] = useState(false);
  const [isCollapsed, setIsCollapsed] = useState(defaultCollapsed);
  const [contentHeight, setContentHeight] = useState<number | null>(null);
  const contentRef = useRef<HTMLDivElement>(null);

  // Calculate line count for the hint
  const lineCount = code.split('\n').length;

  // Measure content height for smooth animation
  useEffect(() => {
    if (contentRef.current) {
      setContentHeight(contentRef.current.scrollHeight);
    }
  }, [code]);

  const handleCopy = useCallback(async (e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent toggle when clicking copy
    try {
      await navigator.clipboard.writeText(code);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  }, [code]);

  const handleToggle = useCallback(() => {
    setIsCollapsed(prev => !prev);
  }, []);

  return (
    <div className="relative my-2 rounded-lg overflow-hidden bg-gray-800">
      {/* Clickable header with language label, line count, and copy button */}
      <div
        className="flex items-center justify-between px-4 py-2 bg-gray-900/50 border-b border-gray-700 cursor-pointer select-none hover:bg-gray-900/70 transition-colors"
        onClick={handleToggle}
        role="button"
        aria-expanded={!isCollapsed}
        aria-label={isCollapsed ? `Expand ${language} code block` : `Collapse ${language} code block`}
      >
        <div className="flex items-center gap-2">
          <ChevronIcon isOpen={!isCollapsed} className="text-gray-400" />
          <span className="text-xs text-gray-400 font-mono">{language}</span>
          {isCollapsed && (
            <span className="text-xs text-gray-500">
              {lineCount} {lineCount === 1 ? 'line' : 'lines'}
            </span>
          )}
        </div>
        <button
          onClick={handleCopy}
          className="flex items-center gap-1 px-2 py-1 text-xs text-gray-400 hover:text-white transition-colors rounded hover:bg-gray-700"
          aria-label={copied ? 'Copied!' : 'Copy code'}
        >
          {copied ? (
            <>
              <CheckIcon className="w-4 h-4" />
              <span>Copied!</span>
            </>
          ) : (
            <>
              <CopyIcon className="w-4 h-4" />
              <span>Copy</span>
            </>
          )}
        </button>
      </div>

      {/* Collapsible code content with smooth animation */}
      <div
        ref={contentRef}
        className="overflow-hidden transition-all duration-200 ease-in-out"
        style={{
          maxHeight: isCollapsed ? 0 : (contentHeight ?? 'none'),
          opacity: isCollapsed ? 0 : 1,
        }}
      >
        <div className="overflow-x-auto">
          <pre className="p-4 text-sm">
            <code className="text-gray-100 font-mono whitespace-pre">{code}</code>
          </pre>
        </div>
      </div>
    </div>
  );
}

function CopyIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
      />
    </svg>
  );
}

function CheckIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M5 13l4 4L19 7"
      />
    </svg>
  );
}
