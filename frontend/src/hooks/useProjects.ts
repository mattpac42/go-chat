'use client';

import { useState, useCallback, useEffect } from 'react';
import { Project, Message } from '@/types';
import { api, ProjectWithMessages, ApiError } from '@/lib/api';

interface UseProjectsReturn {
  projects: Project[];
  isLoading: boolean;
  error: string | null;
  fetchProjects: () => Promise<void>;
  createProject: (title?: string) => Promise<Project | null>;
  deleteProject: (id: string) => Promise<boolean>;
  renameProject: (id: string, title: string) => Promise<Project | null>;
  clearError: () => void;
}

interface UseProjectReturn {
  project: Project | null;
  messages: Message[];
  isLoading: boolean;
  error: string | null;
  fetchProject: () => Promise<void>;
  clearError: () => void;
}

/**
 * Hook for managing the list of projects
 */
export function useProjects(): UseProjectsReturn {
  const [projects, setProjects] = useState<Project[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  /**
   * Fetch all projects from the API
   */
  const fetchProjects = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await api.listProjects();
      // Sort by updatedAt descending (most recent first)
      const sorted = [...data].sort(
        (a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
      );
      setProjects(sorted);
    } catch (err) {
      const errorMessage = err instanceof ApiError
        ? err.message
        : 'Failed to load projects';
      setError(errorMessage);
      console.error('Failed to fetch projects:', err);
    } finally {
      setIsLoading(false);
    }
  }, []);

  /**
   * Create a new project
   */
  const createProject = useCallback(async (title?: string): Promise<Project | null> => {
    setError(null);

    try {
      const newProject = await api.createProject({ title });
      // Add to the beginning of the list (most recent)
      setProjects(prev => [newProject, ...prev]);
      return newProject;
    } catch (err) {
      const errorMessage = err instanceof ApiError
        ? err.message
        : 'Failed to create project';
      setError(errorMessage);
      console.error('Failed to create project:', err);
      return null;
    }
  }, []);

  /**
   * Delete a project
   */
  const deleteProject = useCallback(async (id: string): Promise<boolean> => {
    setError(null);

    try {
      await api.deleteProject(id);
      setProjects(prev => prev.filter(p => p.id !== id));
      return true;
    } catch (err) {
      const errorMessage = err instanceof ApiError
        ? err.message
        : 'Failed to delete project';
      setError(errorMessage);
      console.error('Failed to delete project:', err);
      return false;
    }
  }, []);

  /**
   * Rename a project
   */
  const renameProject = useCallback(async (id: string, title: string): Promise<Project | null> => {
    setError(null);

    try {
      const updatedProject = await api.updateProject(id, { title });
      // Update the project in the list
      setProjects(prev =>
        prev.map(p => p.id === id ? { ...p, title: updatedProject.title, updatedAt: updatedProject.updatedAt } : p)
      );
      return updatedProject;
    } catch (err) {
      const errorMessage = err instanceof ApiError
        ? err.message
        : 'Failed to rename project';
      setError(errorMessage);
      console.error('Failed to rename project:', err);
      return null;
    }
  }, []);

  /**
   * Clear error state
   */
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Fetch projects on mount
  useEffect(() => {
    fetchProjects();
  }, [fetchProjects]);

  return {
    projects,
    isLoading,
    error,
    fetchProjects,
    createProject,
    deleteProject,
    renameProject,
    clearError,
  };
}

/**
 * Hook for fetching a single project with messages
 */
export function useProject(projectId: string): UseProjectReturn {
  const [project, setProject] = useState<Project | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  /**
   * Fetch project with messages from the API
   */
  const fetchProject = useCallback(async () => {
    if (!projectId) {
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const data: ProjectWithMessages = await api.getProject(projectId);
      setProject({
        id: data.id,
        title: data.title,
        createdAt: data.createdAt,
        updatedAt: data.updatedAt,
      });
      // Transform messages from API format (createdAt) to frontend format (timestamp)
      const transformedMessages: Message[] = (data.messages || []).map((msg) => ({
        ...msg,
        timestamp: (msg as { createdAt?: string }).createdAt || msg.timestamp || new Date().toISOString(),
      }));
      setMessages(transformedMessages);
    } catch (err) {
      const errorMessage = err instanceof ApiError
        ? err.message
        : 'Failed to load project';
      setError(errorMessage);
      console.error('Failed to fetch project:', err);
    } finally {
      setIsLoading(false);
    }
  }, [projectId]);

  /**
   * Clear error state
   */
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Fetch project on mount or when projectId changes
  useEffect(() => {
    fetchProject();
  }, [fetchProject]);

  return {
    project,
    messages,
    isLoading,
    error,
    fetchProject,
    clearError,
  };
}
