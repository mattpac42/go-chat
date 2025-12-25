import { api, ApiError, getWebSocketUrl } from '@/lib/api';

// Get the API base URL from environment, matching the lib/api.ts behavior
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Mock fetch globally
const mockFetch = jest.fn();
global.fetch = mockFetch;

describe('API Client', () => {
  beforeEach(() => {
    mockFetch.mockClear();
  });

  describe('listProjects', () => {
    it('returns projects on success', async () => {
      const mockProjects = [
        { id: '1', title: 'Project 1', createdAt: '2025-01-01', updatedAt: '2025-01-01' },
        { id: '2', title: 'Project 2', createdAt: '2025-01-02', updatedAt: '2025-01-02' },
      ];

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ projects: mockProjects }),
      });

      const result = await api.listProjects();

      expect(mockFetch).toHaveBeenCalledWith(
        `${API_BASE_URL}/api/projects`,
        expect.objectContaining({
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
        })
      );
      expect(result).toEqual(mockProjects);
    });

    it('throws ApiError on HTTP error', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => ({ error: 'Internal server error' }),
      });

      await expect(api.listProjects()).rejects.toThrow(ApiError);
    });

    it('includes status code in ApiError', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => ({ error: 'Internal server error' }),
      });

      try {
        await api.listProjects();
        fail('Expected ApiError to be thrown');
      } catch (error) {
        expect(error).toBeInstanceOf(ApiError);
        expect((error as ApiError).statusCode).toBe(500);
        expect((error as ApiError).message).toBe('Internal server error');
      }
    });
  });

  describe('createProject', () => {
    it('creates a project with title', async () => {
      const mockProject = {
        id: '123',
        title: 'New Project',
        createdAt: '2025-01-01',
        updatedAt: '2025-01-01',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockProject,
      });

      const result = await api.createProject({ title: 'New Project' });

      expect(mockFetch).toHaveBeenCalledWith(
        `${API_BASE_URL}/api/projects`,
        expect.objectContaining({
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ title: 'New Project' }),
        })
      );
      expect(result).toEqual(mockProject);
    });

    it('creates a project without title', async () => {
      const mockProject = {
        id: '123',
        title: 'Untitled',
        createdAt: '2025-01-01',
        updatedAt: '2025-01-01',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockProject,
      });

      const result = await api.createProject();

      expect(mockFetch).toHaveBeenCalledWith(
        `${API_BASE_URL}/api/projects`,
        expect.objectContaining({
          body: JSON.stringify({}),
        })
      );
      expect(result).toEqual(mockProject);
    });
  });

  describe('getProject', () => {
    it('returns project with messages', async () => {
      const mockProject = {
        id: '1',
        title: 'Project 1',
        createdAt: '2025-01-01',
        updatedAt: '2025-01-01',
        messages: [
          { id: 'msg-1', projectId: '1', role: 'user', content: 'Hello', timestamp: '2025-01-01' },
        ],
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockProject,
      });

      const result = await api.getProject('1');

      expect(mockFetch).toHaveBeenCalledWith(
        `${API_BASE_URL}/api/projects/1`,
        expect.objectContaining({
          method: 'GET',
        })
      );
      expect(result).toEqual(mockProject);
    });

    it('throws ApiError for not found', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => ({ error: 'Project not found' }),
      });

      await expect(api.getProject('nonexistent')).rejects.toThrow(ApiError);
    });
  });

  describe('deleteProject', () => {
    it('deletes a project successfully', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
      });

      await expect(api.deleteProject('1')).resolves.toBeUndefined();

      expect(mockFetch).toHaveBeenCalledWith(
        `${API_BASE_URL}/api/projects/1`,
        expect.objectContaining({
          method: 'DELETE',
        })
      );
    });

    it('throws ApiError on failure', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => ({ error: 'Project not found' }),
      });

      await expect(api.deleteProject('nonexistent')).rejects.toThrow(ApiError);
    });
  });

  describe('getWebSocketUrl', () => {
    it('returns correct WebSocket URL', () => {
      const url = getWebSocketUrl('project-123');
      const wsProtocol = API_BASE_URL.startsWith('https') ? 'wss' : 'ws';
      const baseUrl = API_BASE_URL.replace(/^https?/, wsProtocol);
      expect(url).toBe(`${baseUrl}/ws/chat?projectId=project-123`);
    });
  });
});

describe('ApiError', () => {
  it('creates error with status code and message', () => {
    const error = new ApiError(404, 'Not found', 'NOT_FOUND');

    expect(error).toBeInstanceOf(Error);
    expect(error.statusCode).toBe(404);
    expect(error.message).toBe('Not found');
    expect(error.code).toBe('NOT_FOUND');
    expect(error.name).toBe('ApiError');
  });
});
