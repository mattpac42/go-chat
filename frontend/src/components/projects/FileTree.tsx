'use client';

import { useState } from 'react';
import { FileNode } from '@/types';

interface FileTreeProps {
  files: FileNode[];
  selectedFile?: string;
  onFileSelect?: (file: FileNode) => void;
}

interface FileTreeNodeProps {
  node: FileNode;
  depth: number;
  selectedFile?: string;
  onFileSelect?: (file: FileNode) => void;
}

// File extension to icon mapping
function getFileIcon(filename: string, language?: string): string {
  const ext = filename.split('.').pop()?.toLowerCase() || '';

  // Check language first
  if (language) {
    switch (language) {
      case 'typescript':
      case 'tsx':
        return 'ts';
      case 'javascript':
      case 'jsx':
        return 'js';
      case 'go':
        return 'go';
      case 'python':
        return 'py';
      case 'rust':
        return 'rs';
    }
  }

  // Fallback to extension
  switch (ext) {
    case 'ts':
    case 'tsx':
      return 'ts';
    case 'js':
    case 'jsx':
      return 'js';
    case 'go':
      return 'go';
    case 'py':
      return 'py';
    case 'rs':
      return 'rs';
    case 'json':
      return 'json';
    case 'md':
      return 'md';
    case 'css':
    case 'scss':
    case 'sass':
      return 'css';
    case 'html':
      return 'html';
    case 'sql':
      return 'sql';
    case 'yaml':
    case 'yml':
      return 'yaml';
    default:
      return 'file';
  }
}

function FileIcon({ type }: { type: string }) {
  const colors: Record<string, string> = {
    ts: 'text-blue-600',
    js: 'text-yellow-500',
    go: 'text-cyan-500',
    py: 'text-green-500',
    rs: 'text-orange-600',
    json: 'text-yellow-600',
    md: 'text-gray-600',
    css: 'text-pink-500',
    html: 'text-orange-500',
    sql: 'text-purple-500',
    yaml: 'text-red-400',
    file: 'text-gray-400',
  };

  return (
    <svg
      className={`w-4 h-4 flex-shrink-0 ${colors[type] || colors.file}`}
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

function FolderIcon({ isOpen }: { isOpen: boolean }) {
  return (
    <svg
      className="w-4 h-4 flex-shrink-0 text-amber-500"
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
      className={`w-3 h-3 flex-shrink-0 text-gray-400 transition-transform duration-150 ${
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

function FileTreeNode({ node, depth, selectedFile, onFileSelect }: FileTreeNodeProps) {
  const [isOpen, setIsOpen] = useState(depth === 0); // Auto-expand root level
  const isFolder = node.type === 'folder';
  const isSelected = node.id === selectedFile;
  const paddingLeft = depth * 16;

  const handleClick = () => {
    if (isFolder) {
      setIsOpen(!isOpen);
    } else {
      onFileSelect?.(node);
    }
  };

  return (
    <div>
      <button
        onClick={handleClick}
        className={`w-full flex items-center gap-1.5 py-1 px-2 text-sm text-left hover:bg-gray-100 rounded transition-colors ${
          isSelected ? 'bg-teal-50 text-teal-700' : 'text-gray-700'
        }`}
        style={{ paddingLeft: `${paddingLeft + 8}px` }}
      >
        {isFolder && <ChevronIcon isOpen={isOpen} />}
        {isFolder ? (
          <FolderIcon isOpen={isOpen} />
        ) : (
          <>
            <span className="w-3" /> {/* Spacer to align with folders */}
            <FileIcon type={getFileIcon(node.name, node.language)} />
          </>
        )}
        <span className="truncate">{node.name}</span>
      </button>
      {isFolder && isOpen && node.children && (
        <div>
          {node.children.map((child) => (
            <FileTreeNode
              key={child.id}
              node={child}
              depth={depth + 1}
              selectedFile={selectedFile}
              onFileSelect={onFileSelect}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export function FileTree({ files, selectedFile, onFileSelect }: FileTreeProps) {
  if (files.length === 0) {
    return (
      <div className="text-sm text-gray-500 p-4 text-center">
        No files generated yet
      </div>
    );
  }

  return (
    <div className="py-2">
      {files.map((node) => (
        <FileTreeNode
          key={node.id}
          node={node}
          depth={0}
          selectedFile={selectedFile}
          onFileSelect={onFileSelect}
        />
      ))}
    </div>
  );
}
