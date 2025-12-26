'use client';

import { useState, useCallback } from 'react';
import { FileNode, FileWithContent } from '@/types';
import { FileRevealCard, RevealTier } from './FileRevealCard';

// Track file expansion state
interface FileState {
  tier: RevealTier;
  pinned: boolean; // true if user explicitly expanded via chevron/code button
}

interface FileRevealListProps {
  files: FileNode[];
  onLoadContent?: (fileId: string) => Promise<FileWithContent | null>;
  groupByFunction?: boolean;
  showEmptyState?: boolean;
}

function FolderIcon({ isOpen }: { isOpen: boolean }) {
  return (
    <svg
      className="w-5 h-5 flex-shrink-0 text-amber-500"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      {isOpen ? (
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M5 19a2 2 0 01-2-2V7a2 2 0 012-2h4l2 2h4a2 2 0 012 2v1M5 19h14a2 2 0 002-2v-5a2 2 0 00-2-2H9a2 2 0 00-2 2v5a2 2 0 01-2 2z"
        />
      ) : (
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
        />
      )}
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

interface FolderSectionProps {
  folder: FileNode;
  depth: number;
  onLoadContent?: (fileId: string) => Promise<FileWithContent | null>;
  expandedFiles: Set<string>;
  onToggleFile: (fileId: string) => void;
  getFileTier: (fileId: string) => RevealTier;
  onCardClick: (fileId: string, hasLongDesc: boolean) => void;
  onIntentionalExpand: (fileId: string, tier: RevealTier) => void;
}

function FolderSection({
  folder,
  depth,
  onLoadContent,
  expandedFiles,
  onToggleFile,
  getFileTier,
  onCardClick,
  onIntentionalExpand,
}: FolderSectionProps) {
  const [isOpen, setIsOpen] = useState(depth < 2); // Auto-expand first two levels

  const files = folder.children?.filter((c) => c.type === 'file') || [];
  const folders = folder.children?.filter((c) => c.type === 'folder') || [];

  return (
    <div className="mb-2">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full flex items-center gap-2 px-3 py-2 text-left hover:bg-gray-50 rounded-lg transition-colors"
        style={{ paddingLeft: `${depth * 16 + 12}px` }}
      >
        <ChevronIcon isOpen={isOpen} />
        <FolderIcon isOpen={isOpen} />
        <span className="font-medium text-gray-700">{folder.name}</span>
        {folder.functionalGroup && (
          <span className="text-xs px-2 py-0.5 bg-teal-50 text-teal-600 rounded-full ml-auto">
            {folder.functionalGroup}
          </span>
        )}
      </button>

      {isOpen && (
        <div
          className="mt-1 space-y-2"
          style={{ paddingLeft: `${depth * 16 + 24}px` }}
        >
          {/* Render subfolders */}
          {folders.map((subfolder) => (
            <FolderSection
              key={subfolder.id}
              folder={subfolder}
              depth={depth + 1}
              onLoadContent={onLoadContent}
              expandedFiles={expandedFiles}
              onToggleFile={onToggleFile}
              getFileTier={getFileTier}
              onCardClick={onCardClick}
              onIntentionalExpand={onIntentionalExpand}
            />
          ))}

          {/* Render files as reveal cards */}
          {files.map((file) => (
            <FileRevealCard
              key={file.id}
              file={file}
              onLoadContent={onLoadContent}
              tier={getFileTier(file.id)}
              onCardClick={() => onCardClick(file.id, !!file.longDescription)}
              onIntentionalExpand={(tier) => onIntentionalExpand(file.id, tier)}
            />
          ))}
        </div>
      )}
    </div>
  );
}

interface FunctionalGroupSectionProps {
  groupName: string;
  files: FileNode[];
  onLoadContent?: (fileId: string) => Promise<FileWithContent | null>;
  getFileTier: (fileId: string) => RevealTier;
  onCardClick: (fileId: string, hasLongDesc: boolean) => void;
  onIntentionalExpand: (fileId: string, tier: RevealTier) => void;
}

function FunctionalGroupSection({
  groupName,
  files,
  onLoadContent,
  getFileTier,
  onCardClick,
  onIntentionalExpand,
}: FunctionalGroupSectionProps) {
  const [isOpen, setIsOpen] = useState(true);

  // Color mapping for functional groups
  const groupColors: Record<string, string> = {
    'Homepage': 'bg-blue-50 text-blue-700 border-blue-200',
    'Navigation': 'bg-indigo-50 text-indigo-700 border-indigo-200',
    'Contact Form': 'bg-green-50 text-green-700 border-green-200',
    'About Page': 'bg-purple-50 text-purple-700 border-purple-200',
    'Backend Services': 'bg-orange-50 text-orange-700 border-orange-200',
    'Database': 'bg-red-50 text-red-700 border-red-200',
    'Configuration': 'bg-gray-100 text-gray-700 border-gray-200',
    'Documentation': 'bg-yellow-50 text-yellow-700 border-yellow-200',
    'UI Components': 'bg-pink-50 text-pink-700 border-pink-200',
    'User Interface': 'bg-cyan-50 text-cyan-700 border-cyan-200',
    'Pages': 'bg-teal-50 text-teal-700 border-teal-200',
  };

  const colorClasses = groupColors[groupName] || 'bg-gray-50 text-gray-700 border-gray-200';

  return (
    <div className="mb-4">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className={`w-full flex items-center gap-3 px-4 py-3 text-left border rounded-lg transition-colors ${colorClasses}`}
      >
        <ChevronIcon isOpen={isOpen} />
        <span className="font-medium">{groupName}</span>
        <span className="ml-auto text-sm opacity-75">
          {files.length} {files.length === 1 ? 'file' : 'files'}
        </span>
      </button>

      {isOpen && (
        <div className="mt-2 ml-4 space-y-2">
          {files.map((file) => (
            <FileRevealCard
              key={file.id}
              file={file}
              onLoadContent={onLoadContent}
              tier={getFileTier(file.id)}
              onCardClick={() => onCardClick(file.id, !!file.longDescription)}
              onIntentionalExpand={(tier) => onIntentionalExpand(file.id, tier)}
            />
          ))}
        </div>
      )}
    </div>
  );
}

/**
 * Flatten file tree to get all files
 */
function flattenFileTree(nodes: FileNode[]): FileNode[] {
  const result: FileNode[] = [];

  function traverse(node: FileNode) {
    if (node.type === 'file') {
      result.push(node);
    }
    if (node.children) {
      node.children.forEach(traverse);
    }
  }

  nodes.forEach(traverse);
  return result;
}

/**
 * Group files by functional group
 */
function groupFilesByFunction(nodes: FileNode[]): Map<string, FileNode[]> {
  const files = flattenFileTree(nodes);
  const groups = new Map<string, FileNode[]>();

  // Define group order for consistent display
  const groupOrder = [
    'Homepage',
    'Navigation',
    'Contact Form',
    'About Page',
    'UI Components',
    'Pages',
    'Backend Services',
    'Database',
    'Configuration',
    'Documentation',
    'Other',
  ];

  // Initialize groups in order
  for (const group of groupOrder) {
    groups.set(group, []);
  }

  // Populate groups
  for (const file of files) {
    const group = file.functionalGroup || 'Other';
    if (!groups.has(group)) {
      groups.set(group, []);
    }
    groups.get(group)!.push(file);
  }

  // Remove empty groups
  const keysToDelete: string[] = [];
  groups.forEach((value, key) => {
    if (value.length === 0) {
      keysToDelete.push(key);
    }
  });
  keysToDelete.forEach((key) => groups.delete(key));

  return groups;
}

/**
 * FileRevealList - Displays files with 2-tier reveal system
 *
 * Can display files in two modes:
 * 1. Folder structure (traditional tree, but with reveal cards for files)
 * 2. Functional groups (organized by purpose, not file path)
 */
export function FileRevealList({
  files,
  onLoadContent,
  groupByFunction = false,
  showEmptyState = true,
}: FileRevealListProps) {
  const [fileStates, setFileStates] = useState<Map<string, FileState>>(new Map());

  // Get tier for a file
  const getFileTier = useCallback((fileId: string): RevealTier => {
    return fileStates.get(fileId)?.tier ?? 'collapsed';
  }, [fileStates]);

  // Handle casual card click - collapse non-pinned files, toggle this file (collapsed <-> details only)
  const handleCardClick = useCallback((fileId: string, hasLongDesc: boolean) => {
    setFileStates((prev) => {
      const next = new Map(prev);

      // Collapse all non-pinned files
      next.forEach((state, id) => {
        if (id !== fileId && !state.pinned) {
          next.set(id, { tier: 'collapsed', pinned: false });
        }
      });

      // Toggle this file's tier between collapsed and details only (not code)
      const current = prev.get(fileId)?.tier ?? 'collapsed';
      let nextTier: RevealTier;
      if (current === 'collapsed') nextTier = hasLongDesc ? 'details' : 'collapsed';
      else nextTier = 'collapsed'; // from details or code -> collapsed

      next.set(fileId, { tier: nextTier, pinned: false });
      return next;
    });
  }, []);

  // Handle intentional expand (chevron/code button) - pins the file open
  const handleIntentionalExpand = useCallback((fileId: string, tier: RevealTier) => {
    setFileStates((prev) => {
      const next = new Map(prev);
      next.set(fileId, { tier, pinned: tier !== 'collapsed' });
      return next;
    });
  }, []);

  // Legacy handler for folder sections (simplified)
  const handleToggleFile = useCallback((fileId: string) => {
    setFileStates((prev) => {
      const next = new Map(prev);
      const current = prev.get(fileId);
      if (current && current.tier !== 'collapsed') {
        next.set(fileId, { tier: 'collapsed', pinned: false });
      } else {
        next.set(fileId, { tier: 'details', pinned: true });
      }
      return next;
    });
  }, []);

  if (files.length === 0 && showEmptyState) {
    return (
      <div className="text-center py-8">
        <div className="w-16 h-16 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
          <svg
            className="w-8 h-8 text-gray-400"
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
        <p className="text-gray-500 text-sm">No files generated yet</p>
        <p className="text-gray-400 text-xs mt-1">
          Start a conversation to create your first file
        </p>
      </div>
    );
  }

  // Functional group view
  if (groupByFunction) {
    const groups = groupFilesByFunction(files);

    return (
      <div className="space-y-2">
        {Array.from(groups.entries()).map(([groupName, groupFiles]) => (
          <FunctionalGroupSection
            key={groupName}
            groupName={groupName}
            files={groupFiles}
            onLoadContent={onLoadContent}
            getFileTier={getFileTier}
            onCardClick={handleCardClick}
            onIntentionalExpand={handleIntentionalExpand}
          />
        ))}
      </div>
    );
  }

  // Folder structure view
  const rootFolders = files.filter((f) => f.type === 'folder');
  const rootFiles = files.filter((f) => f.type === 'file');

  // Convert fileStates to expandedFiles for legacy FolderSection
  const expandedFiles = new Set(
    Array.from(fileStates.entries())
      .filter(([, state]) => state.tier !== 'collapsed')
      .map(([id]) => id)
  );

  return (
    <div className="space-y-2">
      {/* Root folders */}
      {rootFolders.map((folder) => (
        <FolderSection
          key={folder.id}
          folder={folder}
          depth={0}
          onLoadContent={onLoadContent}
          expandedFiles={expandedFiles}
          onToggleFile={handleToggleFile}
          getFileTier={getFileTier}
          onCardClick={handleCardClick}
          onIntentionalExpand={handleIntentionalExpand}
        />
      ))}

      {/* Root files */}
      {rootFiles.map((file) => (
        <FileRevealCard
          key={file.id}
          file={file}
          onLoadContent={onLoadContent}
          tier={getFileTier(file.id)}
          onCardClick={() => handleCardClick(file.id, !!file.longDescription)}
          onIntentionalExpand={(tier) => handleIntentionalExpand(file.id, tier)}
        />
      ))}
    </div>
  );
}

export default FileRevealList;
