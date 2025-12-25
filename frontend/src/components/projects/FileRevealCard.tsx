'use client';

import { useState, useCallback, useRef, useEffect } from 'react';
import { FileNode, FileWithContent } from '@/types';

interface FileRevealCardProps {
  file: FileNode;
  onLoadContent?: (fileId: string) => Promise<FileWithContent | null>;
  isExpanded?: boolean;
  onToggleExpand?: () => void;
}

/**
 * File extension to language display name mapping
 */
function getLanguageDisplay(filename: string, language?: string): string {
  if (language) {
    const displayNames: Record<string, string> = {
      typescript: 'TypeScript',
      tsx: 'TypeScript React',
      javascript: 'JavaScript',
      jsx: 'JavaScript React',
      go: 'Go',
      python: 'Python',
      rust: 'Rust',
      json: 'JSON',
      markdown: 'Markdown',
      css: 'CSS',
      html: 'HTML',
      sql: 'SQL',
      yaml: 'YAML',
    };
    return displayNames[language] || language;
  }

  const ext = filename.split('.').pop()?.toLowerCase() || '';
  const extMap: Record<string, string> = {
    ts: 'TypeScript',
    tsx: 'TypeScript React',
    js: 'JavaScript',
    jsx: 'JavaScript React',
    go: 'Go',
    py: 'Python',
    rs: 'Rust',
    json: 'JSON',
    md: 'Markdown',
    css: 'CSS',
    scss: 'SCSS',
    html: 'HTML',
    sql: 'SQL',
    yaml: 'YAML',
    yml: 'YAML',
  };
  return extMap[ext] || 'File';
}

/**
 * Get accent color based on file type
 */
function getAccentColor(filename: string, language?: string): string {
  const ext = filename.split('.').pop()?.toLowerCase() || '';
  const lang = language?.toLowerCase() || ext;

  const colorMap: Record<string, string> = {
    typescript: 'border-blue-400',
    tsx: 'border-blue-400',
    ts: 'border-blue-400',
    javascript: 'border-yellow-400',
    jsx: 'border-yellow-400',
    js: 'border-yellow-400',
    go: 'border-cyan-400',
    python: 'border-green-400',
    py: 'border-green-400',
    rust: 'border-orange-400',
    rs: 'border-orange-400',
    json: 'border-yellow-500',
    css: 'border-pink-400',
    scss: 'border-pink-400',
    html: 'border-orange-500',
    sql: 'border-purple-400',
    yaml: 'border-red-400',
    yml: 'border-red-400',
    md: 'border-gray-400',
  };

  return colorMap[lang] || colorMap[ext] || 'border-gray-300';
}

function FileTypeIcon({ filename, language }: { filename: string; language?: string }) {
  const ext = filename.split('.').pop()?.toLowerCase() || '';
  const lang = language?.toLowerCase() || ext;

  // Color classes for each file type
  const colorClasses: Record<string, string> = {
    typescript: 'text-blue-600',
    tsx: 'text-blue-600',
    ts: 'text-blue-600',
    javascript: 'text-yellow-500',
    jsx: 'text-yellow-500',
    js: 'text-yellow-500',
    go: 'text-cyan-500',
    python: 'text-green-500',
    py: 'text-green-500',
    rust: 'text-orange-600',
    rs: 'text-orange-600',
    json: 'text-yellow-600',
    css: 'text-pink-500',
    scss: 'text-pink-500',
    html: 'text-orange-500',
    sql: 'text-purple-500',
    yaml: 'text-red-400',
    yml: 'text-red-400',
    md: 'text-gray-600',
  };

  const colorClass = colorClasses[lang] || colorClasses[ext] || 'text-gray-400';

  return (
    <svg
      className={`w-5 h-5 flex-shrink-0 ${colorClass}`}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
      />
    </svg>
  );
}

function ChevronIcon({ isOpen }: { isOpen: boolean }) {
  return (
    <svg
      className={`w-4 h-4 text-gray-400 transition-transform duration-200 ${
        isOpen ? 'rotate-90' : ''
      }`}
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

function CodeIcon({ className }: { className?: string }) {
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
        d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"
      />
    </svg>
  );
}

function CopyIcon({ className }: { className?: string }) {
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

function LoadingSpinner({ className }: { className?: string }) {
  return (
    <svg
      className={`animate-spin ${className}`}
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        className="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        strokeWidth="4"
      />
      <path
        className="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      />
    </svg>
  );
}

/**
 * FileRevealCard - 2-tier file reveal system
 *
 * Tier 1 (default): Shows human-readable description of what the file does
 * Tier 2 (expanded): Shows actual code with syntax highlighting
 */
export function FileRevealCard({
  file,
  onLoadContent,
  isExpanded: controlledExpanded,
  onToggleExpand,
}: FileRevealCardProps) {
  const [internalExpanded, setInternalExpanded] = useState(false);
  const [content, setContent] = useState<string | null>(file.content || null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const codeRef = useRef<HTMLPreElement>(null);

  // Support both controlled and uncontrolled modes
  const isExpanded = controlledExpanded !== undefined ? controlledExpanded : internalExpanded;

  const accentColor = getAccentColor(file.name, file.language);
  const languageDisplay = getLanguageDisplay(file.name, file.language);

  // Generate default description if not provided
  const shortDescription = file.shortDescription || generateDefaultDescription(file);
  const longDescription = file.longDescription;

  const handleToggle = useCallback(async () => {
    const newExpanded = !isExpanded;

    if (onToggleExpand) {
      onToggleExpand();
    } else {
      setInternalExpanded(newExpanded);
    }

    // Load content when expanding if not already loaded
    if (newExpanded && !content && onLoadContent) {
      setIsLoading(true);
      setError(null);
      try {
        const fileData = await onLoadContent(file.id);
        if (fileData) {
          setContent(fileData.content);
        } else {
          setError('Failed to load file content');
        }
      } catch {
        setError('Failed to load file content');
      } finally {
        setIsLoading(false);
      }
    }
  }, [isExpanded, content, file.id, onLoadContent, onToggleExpand]);

  const handleCopy = useCallback(async () => {
    if (!content) return;
    try {
      await navigator.clipboard.writeText(content);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  }, [content]);

  // Measure content height for smooth animation
  const [contentHeight, setContentHeight] = useState(0);
  const contentRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (contentRef.current && isExpanded) {
      setContentHeight(contentRef.current.scrollHeight);
    }
  }, [isExpanded, content]);

  return (
    <div
      className={`bg-white rounded-lg border-l-4 ${accentColor} shadow-sm hover:shadow-md transition-shadow duration-200`}
    >
      {/* Tier 1: Description Header (always visible) */}
      <button
        onClick={handleToggle}
        className="w-full text-left p-4 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:ring-inset rounded-lg"
        aria-expanded={isExpanded}
      >
        <div className="flex items-start gap-3">
          <FileTypeIcon filename={file.name} language={file.language} />

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <h3 className="font-medium text-gray-900 truncate">
                {file.name}
              </h3>
              <span className="text-xs px-2 py-0.5 bg-gray-100 rounded text-gray-500 flex-shrink-0">
                {languageDisplay}
              </span>
            </div>

            {/* Short description - Tier 1 primary content */}
            <p className="text-sm text-gray-600 leading-relaxed">
              {shortDescription}
            </p>

            {/* Long description if available */}
            {longDescription && (
              <p className="text-xs text-gray-500 mt-1 leading-relaxed">
                {longDescription}
              </p>
            )}

            {/* Functional group badge if available */}
            {file.functionalGroup && (
              <span className="inline-block mt-2 text-xs px-2 py-0.5 bg-teal-50 text-teal-700 rounded-full">
                {file.functionalGroup}
              </span>
            )}
          </div>

          <div className="flex items-center gap-2 flex-shrink-0">
            <span className="text-xs text-gray-400 hidden sm:inline">
              {isExpanded ? 'Hide code' : 'View code'}
            </span>
            <ChevronIcon isOpen={isExpanded} />
          </div>
        </div>
      </button>

      {/* Tier 2: Code Content (expandable) */}
      <div
        className="overflow-hidden transition-all duration-300 ease-in-out"
        style={{
          maxHeight: isExpanded ? `${contentHeight + 100}px` : '0px',
          opacity: isExpanded ? 1 : 0,
        }}
      >
        <div ref={contentRef} className="border-t border-gray-100">
          {/* Code toolbar */}
          <div className="flex items-center justify-between px-4 py-2 bg-gray-50 border-b border-gray-100">
            <div className="flex items-center gap-2">
              <CodeIcon className="w-4 h-4 text-gray-400" />
              <span className="text-xs text-gray-500 font-mono">
                {file.path}
              </span>
            </div>

            <button
              onClick={(e) => {
                e.stopPropagation();
                handleCopy();
              }}
              disabled={!content || isLoading}
              className="flex items-center gap-1.5 px-2 py-1 text-xs rounded hover:bg-gray-200 transition-colors disabled:opacity-50"
              aria-label={copied ? 'Copied!' : 'Copy code'}
            >
              {copied ? (
                <>
                  <CheckIcon className="w-3.5 h-3.5 text-green-500" />
                  <span className="text-green-600 hidden sm:inline">Copied!</span>
                </>
              ) : (
                <>
                  <CopyIcon className="w-3.5 h-3.5 text-gray-500" />
                  <span className="text-gray-600 hidden sm:inline">Copy</span>
                </>
              )}
            </button>
          </div>

          {/* Code content area */}
          <div className="max-h-80 overflow-auto bg-gray-900 p-4">
            {isLoading ? (
              <div className="flex items-center justify-center py-8">
                <LoadingSpinner className="w-6 h-6 text-gray-400" />
                <span className="ml-2 text-sm text-gray-400">Loading code...</span>
              </div>
            ) : error ? (
              <div className="flex items-center justify-center py-8 text-red-400 text-sm">
                {error}
              </div>
            ) : content ? (
              <pre
                ref={codeRef}
                className="text-sm font-mono text-gray-100 whitespace-pre-wrap break-words leading-relaxed"
              >
                {content}
              </pre>
            ) : (
              <div className="flex items-center justify-center py-8 text-gray-500 text-sm">
                No content available
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

/**
 * Generate a default description based on file name and path
 */
function generateDefaultDescription(file: FileNode): string {
  const name = file.name.toLowerCase();
  const path = file.path.toLowerCase();

  // Common patterns for React/Next.js projects
  if (name === 'page.tsx' || name === 'page.ts') {
    const folder = path.split('/').slice(-2, -1)[0] || 'root';
    return `Page component for the ${folder} route`;
  }
  if (name === 'layout.tsx' || name === 'layout.ts') {
    return 'Shared layout wrapper for this section';
  }
  if (name.endsWith('.test.tsx') || name.endsWith('.test.ts') || name.endsWith('.spec.ts')) {
    const baseName = name.replace(/\.(test|spec)\.(tsx?|jsx?)$/, '');
    return `Test suite for ${baseName}`;
  }
  if (name === 'index.ts' || name === 'index.tsx') {
    return 'Entry point and exports for this module';
  }
  if (path.includes('components/')) {
    const componentName = name.replace(/\.(tsx?|jsx?)$/, '');
    return `${componentName} user interface component`;
  }
  if (path.includes('hooks/')) {
    return 'Custom React hook for shared logic';
  }
  if (path.includes('lib/') || path.includes('utils/')) {
    return 'Utility functions and helpers';
  }
  if (path.includes('api/')) {
    return 'API route handler for backend functionality';
  }
  if (path.includes('types/')) {
    return 'TypeScript type definitions';
  }
  if (name.endsWith('.json')) {
    if (name === 'package.json') return 'Project dependencies and scripts';
    if (name === 'tsconfig.json') return 'TypeScript compiler configuration';
    return 'Configuration data';
  }
  if (name.endsWith('.md')) {
    if (name === 'readme.md') return 'Project documentation and setup instructions';
    return 'Documentation';
  }
  if (name.endsWith('.css') || name.endsWith('.scss')) {
    return 'Styling rules for the user interface';
  }

  // Default fallback
  const ext = file.language || name.split('.').pop() || 'file';
  return `${ext.charAt(0).toUpperCase() + ext.slice(1)} source file`;
}

export default FileRevealCard;
