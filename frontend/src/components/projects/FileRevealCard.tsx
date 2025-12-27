'use client';

import { useState, useCallback, useRef, useEffect } from 'react';
import { FileNode, FileWithContent } from '@/types';
import { API_BASE_URL } from '@/lib/api';
import { useCodeZoom } from '@/hooks/useCodeZoom';
import { ZoomControls } from './ZoomControls';

interface FileRevealCardProps {
  file: FileNode;
  onLoadContent?: (fileId: string) => Promise<FileWithContent | null>;
  /** Current tier state (controlled) */
  tier?: RevealTier;
  /** Called when user clicks card (cycles tiers) */
  onCardClick?: () => void;
  /** Called when user explicitly expands via chevron/code (should "pin" open) */
  onIntentionalExpand?: (tier: RevealTier) => void;
  /** Hide functional group tag (used in purpose view to avoid redundancy) */
  hideFunctionalGroup?: boolean;
}

type RevealTier = 'collapsed' | 'details' | 'code';

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
 * Convert full language name to 2-3 letter abbreviation for compact display
 */
function getAbbreviatedLanguage(language: string): string {
  const abbrevMap: Record<string, string> = {
    'TypeScript': 'TS',
    'TypeScript React': 'TSX',
    'JavaScript': 'JS',
    'JavaScript React': 'JSX',
    'Go': 'GO',
    'Python': 'PY',
    'Rust': 'RS',
    'JSON': 'JSON',
    'Markdown': 'MD',
    'CSS': 'CSS',
    'SCSS': 'SCSS',
    'HTML': 'HTML',
    'SQL': 'SQL',
    'YAML': 'YAML',
    'File': '···',
  };
  return abbrevMap[language] || language.substring(0, 3).toUpperCase();
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

function DownloadIcon({ className }: { className?: string }) {
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
        d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
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
 * FileRevealCard - 3-tier file reveal system
 *
 * Tier 1 (collapsed): File name + short description only
 * Tier 2 (details): Shows long description
 * Tier 3 (code): Shows code (long description stays visible for context)
 *
 * Linear progression: collapsed -> details -> code -> collapsed
 */
export { type RevealTier };

export function FileRevealCard({
  file,
  onLoadContent,
  tier: controlledTier,
  onCardClick,
  onIntentionalExpand,
  hideFunctionalGroup = false,
}: FileRevealCardProps) {
  const [internalTier, setInternalTier] = useState<RevealTier>('collapsed');
  const [content, setContent] = useState<string | null>(file.content || null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const codeRef = useRef<HTMLPreElement>(null);
  const { zoomLevel, setZoom, getZoomStyle } = useCodeZoom();

  // Support controlled and uncontrolled modes
  const tier = controlledTier ?? internalTier;
  const setTier = (newTier: RevealTier) => {
    if (controlledTier === undefined) {
      setInternalTier(newTier);
    }
  };

  const accentColor = getAccentColor(file.name, file.language);
  const languageDisplay = getLanguageDisplay(file.name, file.language);

  const shortDescription = file.shortDescription || generateDefaultDescription(file);
  const longDescription = file.longDescription;
  const hasLongDescription = !!longDescription;

  // Get next tier for chevron and card body click (toggle between collapsed <-> details only)
  // Chevron should NOT cycle to code view - code is only accessible via the code button
  const getNextTierToggle = useCallback((): RevealTier => {
    // If in code view, go to collapsed
    if (tier === 'code') return 'collapsed';
    // If collapsed, go to details (or stay collapsed if no long description)
    if (tier === 'collapsed') return hasLongDescription ? 'details' : 'collapsed';
    // If in details, go back to collapsed
    return 'collapsed';
  }, [tier, hasLongDescription]);

  // Handle main card click - toggles between collapsed and details only (not code)
  const handleCardClick = useCallback(() => {
    const nextTier = getNextTierToggle();
    if (onCardClick) {
      onCardClick();
    } else {
      setTier(nextTier);
    }
  }, [getNextTierToggle, onCardClick]);

  // Handle chevron click - toggles between collapsed and details (not code)
  const handleChevronClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    const nextTier = getNextTierToggle();
    if (onIntentionalExpand) {
      onIntentionalExpand(nextTier);
    } else {
      setTier(nextTier);
    }
  }, [getNextTierToggle, onIntentionalExpand]);

  // Handle code button click - jump to code or collapse
  const handleViewCode = useCallback(async (e: React.MouseEvent) => {
    e.stopPropagation();

    const newTier: RevealTier = tier === 'code' ? 'collapsed' : 'code';

    if (onIntentionalExpand) {
      onIntentionalExpand(newTier);
    } else {
      setTier(newTier);
    }

    // Load content if going to code and not already loaded
    if (newTier === 'code' && !content && onLoadContent) {
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
  }, [tier, content, file.id, onLoadContent, onIntentionalExpand]);

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
    if (contentRef.current && tier === 'code') {
      setContentHeight(contentRef.current.scrollHeight);
    }
  }, [tier, content]);

  const showDetails = tier === 'details' || tier === 'code';
  const showCode = tier === 'code';

  return (
    <div
      className={`bg-white rounded-lg border-l-4 ${accentColor} shadow-sm hover:shadow-md transition-all duration-200 focus-within:ring-2 focus-within:ring-teal-500 focus-within:ring-offset-1`}
    >
      {/* Tier 1: Header (always visible) - click to show/hide details */}
      <div
        onClick={handleCardClick}
        className="w-full text-left p-3 cursor-pointer rounded-t-lg"
        role="button"
        tabIndex={0}
        aria-expanded={showDetails}
        onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') handleCardClick(); }}
      >
        {/* Row 1: Filename + abbreviated language badge */}
        <div className="flex items-center justify-between gap-2 mb-1">
          <h3 className="font-medium text-gray-900 truncate flex-1 min-w-0">
            {file.name}
          </h3>
          <span className="text-xs px-1.5 py-0.5 bg-gray-100 rounded text-gray-500 flex-shrink-0 font-mono">
            {getAbbreviatedLanguage(languageDisplay)}
          </span>
        </div>

        {/* Row 2: Icon + description + actions */}
        <div className="flex items-center gap-2">
          <FileTypeIcon filename={file.name} language={file.language} />

          <p className={`flex-1 min-w-0 text-sm text-gray-600 ${showDetails ? '' : 'truncate'}`}>
            {shortDescription}
          </p>

          <div className="flex items-center gap-0.5 flex-shrink-0">
            {/* Download icon */}
            <button
              onClick={(e) => {
                e.stopPropagation();
                window.open(`${API_BASE_URL}/api/files/${file.id}/download`, '_blank');
              }}
              className="p-1 rounded hover:bg-gray-100 transition-colors"
              aria-label="Download file"
              title="Download file"
            >
              <DownloadIcon className="w-3.5 h-3.5 text-gray-400" />
            </button>
            {/* Code icon */}
            <button
              onClick={handleViewCode}
              className="p-1 rounded hover:bg-gray-100 transition-colors"
              aria-label={showCode ? 'Hide code' : 'View code'}
            >
              <CodeIcon className={`w-3.5 h-3.5 ${showCode ? 'text-teal-500' : 'text-gray-400'}`} />
            </button>
            {/* Chevron */}
            <button
              onClick={handleChevronClick}
              className="p-1 rounded hover:bg-gray-100 transition-colors"
              aria-label={tier === 'collapsed' ? 'Expand' : 'Collapse'}
            >
              <ChevronIcon isOpen={showDetails || showCode} />
            </button>
          </div>
        </div>

        {/* Functional group badge if available */}
        {file.functionalGroup && !hideFunctionalGroup && (
          <span className="inline-block mt-2 text-xs px-2 py-0.5 bg-teal-50 text-teal-700 rounded-full">
            {file.functionalGroup}
          </span>
        )}
      </div>

      {/* Tier 2: Long Description (stays visible with code for context) */}
      {hasLongDescription && (showDetails || showCode) && (
        <div className="px-4 pb-4 pt-0">
          <div className="pl-8 border-l-2 border-gray-100 ml-2">
            <p className="text-sm text-gray-500 leading-relaxed">
              {longDescription}
            </p>
          </div>
        </div>
      )}

      {/* Tier 3: Code Content (expandable) */}
      <div
        className="overflow-hidden transition-all duration-300 ease-in-out"
        style={{
          maxHeight: showCode ? `${contentHeight + 100}px` : '0px',
          opacity: showCode ? 1 : 0,
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

            <div className="flex items-center gap-2">
              <ZoomControls
                zoomLevel={zoomLevel}
                onZoomChange={setZoom}
              />
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
                className="text-sm font-mono text-gray-100 whitespace-pre-wrap break-words"
                style={getZoomStyle()}
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
