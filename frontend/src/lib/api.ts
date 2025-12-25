/**
 * REST API client for project CRUD operations
 */

import { Project, Message, FileItem, FileWithContent } from '@/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

/**
 * API error class with status code and message
 */
export class ApiError extends Error {
  constructor(
    public statusCode: number,
    message: string,
    public code?: string
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Project with messages for GET /api/projects/:id
 */
export interface ProjectWithMessages extends Project {
  messages: Message[];
}

/**
 * Create project request body
 */
export interface CreateProjectRequest {
  title?: string;
}

/**
 * Update project request body
 */
export interface UpdateProjectRequest {
  title: string;
}

/**
 * Generic API response handler
 */
async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let errorMessage = `HTTP error ${response.status}`;
    let errorCode: string | undefined;

    try {
      const errorData = await response.json();
      errorMessage = errorData.error || errorData.message || errorMessage;
      errorCode = errorData.code;
    } catch {
      // Use default error message if parsing fails
    }

    throw new ApiError(response.status, errorMessage, errorCode);
  }

  return response.json();
}

/**
 * API client singleton
 */
export const api = {
  /**
   * List all projects
   * GET /api/projects
   */
  async listProjects(): Promise<Project[]> {
    const response = await fetch(`${API_BASE_URL}/api/projects`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    const data = await handleResponse<{ projects: Project[] }>(response);
    return data.projects || [];
  },

  /**
   * Create a new project
   * POST /api/projects
   */
  async createProject(data?: CreateProjectRequest): Promise<Project> {
    const response = await fetch(`${API_BASE_URL}/api/projects`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data || {}),
    });

    return handleResponse<Project>(response);
  },

  /**
   * Get a project with its messages
   * GET /api/projects/:id
   */
  async getProject(id: string): Promise<ProjectWithMessages> {
    const response = await fetch(`${API_BASE_URL}/api/projects/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    return handleResponse<ProjectWithMessages>(response);
  },

  /**
   * Delete a project
   * DELETE /api/projects/:id
   */
  async deleteProject(id: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/api/projects/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      let errorMessage = `HTTP error ${response.status}`;
      let errorCode: string | undefined;

      try {
        const errorData = await response.json();
        errorMessage = errorData.error || errorData.message || errorMessage;
        errorCode = errorData.code;
      } catch {
        // Use default error message if parsing fails
      }

      throw new ApiError(response.status, errorMessage, errorCode);
    }
  },

  /**
   * Update a project (rename)
   * PATCH /api/projects/:id
   */
  async updateProject(id: string, data: UpdateProjectRequest): Promise<Project> {
    const response = await fetch(`${API_BASE_URL}/api/projects/${id}`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    return handleResponse<Project>(response);
  },

  /**
   * Get files for a project
   * GET /api/projects/:id/files
   */
  async getProjectFiles(projectId: string): Promise<FileItem[]> {
    const response = await fetch(`${API_BASE_URL}/api/projects/${projectId}/files`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    const data = await handleResponse<{ files: FileItem[] }>(response);
    return data.files || [];
  },

  /**
   * Get a file by ID
   * GET /api/files/:id
   */
  async getFile(id: string): Promise<FileWithContent> {
    const response = await fetch(`${API_BASE_URL}/api/files/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    return handleResponse<FileWithContent>(response);
  },
};

/**
 * Get the WebSocket URL for chat
 */
export function getWebSocketUrl(projectId: string): string {
  const wsProtocol = API_BASE_URL.startsWith('https') ? 'wss' : 'ws';
  const baseUrl = API_BASE_URL.replace(/^https?/, wsProtocol);
  return `${baseUrl}/ws/chat?projectId=${projectId}`;
}

export default api;
