'use client';

import { useState, useCallback } from 'react';
import { FileNode, FileWithContent } from '@/types';
import { FileTree } from './FileTree';
import { FileRevealList } from './FileRevealList';
import { FilePreviewModal } from '../shared/FilePreviewModal';

export type ViewMode = 'tree' | 'reveal' | 'grouped';

interface FileExplorerProps {
  files: FileNode[];
  onLoadContent?: (fileId: string) => Promise<FileWithContent | null>;
  defaultViewMode?: ViewMode;
  showViewToggle?: boolean;
  isLoading?: boolean;
}

function TreeViewIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M4 6h16M4 10h16M4 14h16M4 18h16"
      />
    </svg>
  );
}

function CardViewIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M4 5a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM14 5a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1V5zM4 15a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1v-4zM14 15a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z"
      />
    </svg>
  );
}

function GroupViewIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
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
 * FileExplorer - Unified file browsing component
 *
 * Supports three view modes:
 * - tree: Traditional file tree (compact, for quick navigation)
 * - reveal: 2-tier reveal cards (descriptions first, code on expand)
 * - grouped: Files organized by functional group (App Map style)
 */
export function FileExplorer({
  files,
  onLoadContent,
  defaultViewMode = 'reveal',
  showViewToggle = true,
  isLoading = false,
}: FileExplorerProps) {
  const [viewMode, setViewMode] = useState<ViewMode>(defaultViewMode);
  const [selectedFile, setSelectedFile] = useState<string | undefined>();
  const [previewFile, setPreviewFile] = useState<FileWithContent | null>(null);
  const [isPreviewOpen, setIsPreviewOpen] = useState(false);

  // Handle file selection in tree view (opens preview modal)
  const handleFileSelect = useCallback(async (file: FileNode) => {
    setSelectedFile(file.id);

    if (onLoadContent) {
      const fileData = await onLoadContent(file.id);
      if (fileData) {
        setPreviewFile(fileData);
        setIsPreviewOpen(true);
      }
    }
  }, [onLoadContent]);

  const handleClosePreview = useCallback(() => {
    setIsPreviewOpen(false);
    setPreviewFile(null);
  }, []);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">
        <LoadingSpinner className="w-6 h-6 text-gray-400" />
      </div>
    );
  }

  if (files.length === 0) {
    return (
      <div className="text-center py-8 px-4">
        <div className="w-12 h-12 mx-auto mb-3 bg-gray-100 rounded-full flex items-center justify-center">
          <svg
            className="w-6 h-6 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
            />
          </svg>
        </div>
        <p className="text-gray-500 text-sm">No files yet</p>
        <p className="text-gray-400 text-xs mt-1">
          Files will appear here as you build
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      {/* View mode toggle */}
      {showViewToggle && (
        <div className="flex items-center gap-1 px-2 py-2 border-b border-gray-100">
          <button
            onClick={() => setViewMode('grouped')}
            className={`p-1.5 rounded transition-colors ${
              viewMode === 'grouped'
                ? 'bg-teal-100 text-teal-700'
                : 'text-gray-400 hover:bg-gray-100 hover:text-gray-600'
            }`}
            title="Group by function"
            aria-label="Group by function"
          >
            <GroupViewIcon className="w-4 h-4" />
          </button>
          <button
            onClick={() => setViewMode('reveal')}
            className={`p-1.5 rounded transition-colors ${
              viewMode === 'reveal'
                ? 'bg-teal-100 text-teal-700'
                : 'text-gray-400 hover:bg-gray-100 hover:text-gray-600'
            }`}
            title="Card view with descriptions"
            aria-label="Card view with descriptions"
          >
            <CardViewIcon className="w-4 h-4" />
          </button>
          <button
            onClick={() => setViewMode('tree')}
            className={`p-1.5 rounded transition-colors ${
              viewMode === 'tree'
                ? 'bg-teal-100 text-teal-700'
                : 'text-gray-400 hover:bg-gray-100 hover:text-gray-600'
            }`}
            title="Tree view"
            aria-label="Tree view"
          >
            <TreeViewIcon className="w-4 h-4" />
          </button>

          <span className="ml-auto text-xs text-gray-400">
            {viewMode === 'grouped' && 'By Purpose'}
            {viewMode === 'reveal' && 'Descriptions'}
            {viewMode === 'tree' && 'Files'}
          </span>
        </div>
      )}

      {/* File view content */}
      <div className="flex-1 overflow-y-auto">
        {viewMode === 'tree' && (
          <FileTree
            files={files}
            selectedFile={selectedFile}
            onFileSelect={handleFileSelect}
          />
        )}

        {viewMode === 'reveal' && (
          <div className="p-2">
            <FileRevealList
              files={files}
              onLoadContent={onLoadContent}
              groupByFunction={false}
              showEmptyState={false}
            />
          </div>
        )}

        {viewMode === 'grouped' && (
          <div className="p-2">
            <FileRevealList
              files={files}
              onLoadContent={onLoadContent}
              groupByFunction={true}
              showEmptyState={false}
            />
          </div>
        )}
      </div>

      {/* Preview modal for tree view */}
      <FilePreviewModal
        file={previewFile}
        isOpen={isPreviewOpen}
        onClose={handleClosePreview}
      />
    </div>
  );
}

export default FileExplorer;
