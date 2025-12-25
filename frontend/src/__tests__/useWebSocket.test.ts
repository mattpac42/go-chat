import { renderHook, act, waitFor } from '@testing-library/react';
import { useWebSocket } from '@/hooks/useWebSocket';

// Mock WebSocket
class MockWebSocket {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSING = 2;
  static CLOSED = 3;

  url: string;
  readyState: number = MockWebSocket.CONNECTING;
  onopen: (() => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  onclose: ((event: { code: number; reason: string }) => void) | null = null;
  sentMessages: string[] = [];

  constructor(url: string) {
    this.url = url;
    // Simulate async connection
    setTimeout(() => {
      if (this.readyState !== MockWebSocket.CLOSED) {
        this.readyState = MockWebSocket.OPEN;
        this.onopen?.();
      }
    }, 0);
  }

  send(data: string) {
    this.sentMessages.push(data);
  }

  close(code = 1000, reason = '') {
    this.readyState = MockWebSocket.CLOSED;
    this.onclose?.({ code, reason });
  }

  // Test helpers
  simulateMessage(data: object) {
    this.onmessage?.({ data: JSON.stringify(data) });
  }

  simulateError() {
    this.onerror?.(new Event('error'));
  }
}

// Store WebSocket instances for testing
let mockWsInstances: MockWebSocket[] = [];

beforeEach(() => {
  mockWsInstances = [];
  (global as unknown as { WebSocket: typeof MockWebSocket }).WebSocket = class extends MockWebSocket {
    constructor(url: string) {
      super(url);
      mockWsInstances.push(this);
    }
  } as unknown as typeof MockWebSocket;
});

afterEach(() => {
  mockWsInstances = [];
});

describe('useWebSocket', () => {
  it('connects to WebSocket on mount', async () => {
    const { result } = renderHook(() =>
      useWebSocket({ projectId: 'test-project' })
    );

    expect(result.current.status).toBe('connecting');

    await waitFor(() => {
      expect(result.current.status).toBe('connected');
    });

    expect(mockWsInstances.length).toBe(1);
    expect(mockWsInstances[0].url).toContain('projectId=test-project');
  });

  it('calls onConnect when connected', async () => {
    const onConnect = jest.fn();
    renderHook(() =>
      useWebSocket({ projectId: 'test-project', onConnect })
    );

    await waitFor(() => {
      expect(onConnect).toHaveBeenCalled();
    });
  });

  it('sends messages through WebSocket', async () => {
    const { result } = renderHook(() =>
      useWebSocket({ projectId: 'test-project' })
    );

    await waitFor(() => {
      expect(result.current.status).toBe('connected');
    });

    act(() => {
      result.current.sendMessage({
        type: 'chat_message',
        projectId: 'test-project',
        content: 'Hello',
        timestamp: '2025-01-01T00:00:00Z',
      });
    });

    expect(mockWsInstances[0].sentMessages.length).toBe(1);
    expect(JSON.parse(mockWsInstances[0].sentMessages[0])).toEqual({
      type: 'chat_message',
      projectId: 'test-project',
      content: 'Hello',
      timestamp: '2025-01-01T00:00:00Z',
    });
  });

  it('handles incoming messages', async () => {
    const onMessage = jest.fn();
    const { result } = renderHook(() =>
      useWebSocket({ projectId: 'test-project', onMessage })
    );

    await waitFor(() => {
      expect(result.current.status).toBe('connected');
    });

    act(() => {
      mockWsInstances[0].simulateMessage({
        type: 'message_start',
        projectId: 'test-project',
        messageId: 'msg-1',
      });
    });

    expect(onMessage).toHaveBeenCalledWith({
      type: 'message_start',
      projectId: 'test-project',
      messageId: 'msg-1',
    });
  });

  it('disconnects manually', async () => {
    const onDisconnect = jest.fn();
    const { result } = renderHook(() =>
      useWebSocket({ projectId: 'test-project', onDisconnect })
    );

    await waitFor(() => {
      expect(result.current.status).toBe('connected');
    });

    act(() => {
      result.current.disconnect();
    });

    expect(result.current.status).toBe('disconnected');
  });

  it('queues messages while connecting', async () => {
    const { result } = renderHook(() =>
      useWebSocket({ projectId: 'test-project' })
    );

    // Send message while still connecting
    act(() => {
      result.current.sendMessage({
        type: 'chat_message',
        projectId: 'test-project',
        content: 'Queued message',
        timestamp: '2025-01-01T00:00:00Z',
      });
    });

    // Wait for connection
    await waitFor(() => {
      expect(result.current.status).toBe('connected');
    });

    // Message should have been sent after connection
    expect(mockWsInstances[0].sentMessages.length).toBe(1);
    expect(JSON.parse(mockWsInstances[0].sentMessages[0]).content).toBe('Queued message');
  });

  it('reconnects when projectId changes', async () => {
    const { result, rerender } = renderHook(
      ({ projectId }) => useWebSocket({ projectId }),
      { initialProps: { projectId: 'project-1' } }
    );

    await waitFor(() => {
      expect(result.current.status).toBe('connected');
    });

    expect(mockWsInstances.length).toBe(1);

    // Change projectId
    rerender({ projectId: 'project-2' });

    await waitFor(() => {
      expect(mockWsInstances.length).toBe(2);
    });

    expect(mockWsInstances[1].url).toContain('projectId=project-2');
  });

  it('does not auto-connect when autoConnect is false', () => {
    const { result } = renderHook(() =>
      useWebSocket({ projectId: 'test-project', autoConnect: false })
    );

    expect(result.current.status).toBe('disconnected');
    expect(mockWsInstances.length).toBe(0);
  });

  it('connects manually when autoConnect is false', async () => {
    const { result } = renderHook(() =>
      useWebSocket({ projectId: 'test-project', autoConnect: false })
    );

    expect(mockWsInstances.length).toBe(0);

    act(() => {
      result.current.connect();
    });

    await waitFor(() => {
      expect(result.current.status).toBe('connected');
    });

    expect(mockWsInstances.length).toBe(1);
  });
});
