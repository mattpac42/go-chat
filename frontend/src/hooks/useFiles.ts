'use client';

import { useState, useCallback, useEffect } from 'react';
import { FileItem, FileNode, FileWithContent } from '@/types';
import { api, ApiError } from '@/lib/api';

interface UseFilesReturn {
  files: FileItem[];
  fileTree: FileNode[];
  isLoading: boolean;
  error: string | null;
  fetchFiles: () => Promise<void>;
  getFile: (id: string) => Promise<FileWithContent | null>;
  clearError: () => void;
}

/**
 * Build a tree structure from flat file list
 */
function buildFileTree(files: FileItem[]): FileNode[] {
  const root: FileNode[] = [];
  const nodeMap = new Map<string, FileNode>();

  // Sort files by path to ensure parent folders are processed first
  const sortedFiles = [...files].sort((a, b) => a.path.localeCompare(b.path));

  for (const file of sortedFiles) {
    const parts = file.path.split('/');
    let currentPath = '';
    let currentLevel = root;

    // Create folder nodes for each path segment except the last (which is the file)
    for (let i = 0; i < parts.length - 1; i++) {
      const part = parts[i];
      currentPath = currentPath ? `${currentPath}/${part}` : part;

      let folderNode = nodeMap.get(currentPath);
      if (!folderNode) {
        folderNode = {
          id: `folder-${currentPath}`,
          name: part,
          path: currentPath,
          type: 'folder',
          children: [],
        };
        nodeMap.set(currentPath, folderNode);
        currentLevel.push(folderNode);
      }
      currentLevel = folderNode.children!;
    }

    // Add the file node
    const fileNode: FileNode = {
      id: file.id,
      name: file.filename,
      path: file.path,
      type: 'file',
      language: file.language,
    };
    currentLevel.push(fileNode);
  }

  // Sort each level: folders first, then files, alphabetically
  const sortLevel = (nodes: FileNode[]): FileNode[] => {
    return nodes.sort((a, b) => {
      if (a.type !== b.type) {
        return a.type === 'folder' ? -1 : 1;
      }
      return a.name.localeCompare(b.name);
    }).map(node => {
      if (node.children) {
        node.children = sortLevel(node.children);
      }
      return node;
    });
  };

  return sortLevel(root);
}

/**
 * Hook for managing files for a project
 */
export function useFiles(projectId: string): UseFilesReturn {
  const [files, setFiles] = useState<FileItem[]>([]);
  const [fileTree, setFileTree] = useState<FileNode[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  /**
   * Fetch all files for the project
   */
  const fetchFiles = useCallback(async () => {
    console.log('[useFiles] fetchFiles called, projectId:', projectId);
    if (!projectId) {
      setFiles([]);
      setFileTree([]);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      console.log('[useFiles] Fetching files for project:', projectId);
      const data = await api.getProjectFiles(projectId);
      console.log('[useFiles] Got files:', data);
      setFiles(data);
      setFileTree(buildFileTree(data));
    } catch (err) {
      const errorMessage = err instanceof ApiError
        ? err.message
        : 'Failed to load files';
      setError(errorMessage);
      console.error('[useFiles] Failed to fetch files:', err);
    } finally {
      setIsLoading(false);
    }
  }, [projectId]);

  /**
   * Get a single file with content
   */
  const getFile = useCallback(async (id: string): Promise<FileWithContent | null> => {
    try {
      return await api.getFile(id);
    } catch (err) {
      const errorMessage = err instanceof ApiError
        ? err.message
        : 'Failed to load file';
      setError(errorMessage);
      console.error('Failed to fetch file:', err);
      return null;
    }
  }, []);

  /**
   * Clear error state
   */
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Fetch files on mount or when projectId changes
  useEffect(() => {
    fetchFiles();
  }, [fetchFiles]);

  return {
    files,
    fileTree,
    isLoading,
    error,
    fetchFiles,
    getFile,
    clearError,
  };
}
