'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { FileItem, FileWithContent } from '@/types';
import { api } from '@/lib/api';

interface PreviewFile {
  path: string;
  content: string;
}

interface UsePreviewFilesReturn {
  previewFiles: PreviewFile[];
  isLoading: boolean;
  loadAllFiles: (files: FileItem[]) => Promise<void>;
}

/**
 * Hook for loading all file contents for preview
 * Only loads previewable files (HTML, CSS, JS)
 */
export function usePreviewFiles(): UsePreviewFilesReturn {
  const [previewFiles, setPreviewFiles] = useState<PreviewFile[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const loadingRef = useRef(false);

  const loadAllFiles = useCallback(async (files: FileItem[]) => {
    // Filter to only previewable file types
    const previewableFiles = files.filter(f =>
      f.path.endsWith('.html') ||
      f.path.endsWith('.css') ||
      f.path.endsWith('.js')
    );

    if (previewableFiles.length === 0) {
      setPreviewFiles([]);
      return;
    }

    // Prevent concurrent loads
    if (loadingRef.current) return;
    loadingRef.current = true;
    setIsLoading(true);

    try {
      // Load all file contents in parallel
      const loadedFiles = await Promise.all(
        previewableFiles.map(async (file) => {
          try {
            const fileWithContent = await api.getFile(file.id);
            return {
              path: file.path,
              content: fileWithContent.content,
            };
          } catch (error) {
            console.error(`[usePreviewFiles] Failed to load file ${file.path}:`, error);
            return null;
          }
        })
      );

      // Filter out failed loads
      const successfulFiles = loadedFiles.filter((f): f is PreviewFile => f !== null);
      setPreviewFiles(successfulFiles);
    } catch (error) {
      console.error('[usePreviewFiles] Failed to load files:', error);
    } finally {
      setIsLoading(false);
      loadingRef.current = false;
    }
  }, []);

  return {
    previewFiles,
    isLoading,
    loadAllFiles,
  };
}

export default usePreviewFiles;
