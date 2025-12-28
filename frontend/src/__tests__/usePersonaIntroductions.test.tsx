import { renderHook, act } from '@testing-library/react';
import { usePersonaIntroductions } from '@/hooks/usePersonaIntroductions';
import { Message } from '@/types';

// Mock localStorage
const mockLocalStorage = (() => {
  let store: Record<string, string> = {};
  return {
    getItem: jest.fn((key: string) => store[key] || null),
    setItem: jest.fn((key: string, value: string) => {
      store[key] = value;
    }),
    removeItem: jest.fn((key: string) => {
      delete store[key];
    }),
    clear: () => {
      store = {};
    },
  };
})();

Object.defineProperty(window, 'localStorage', { value: mockLocalStorage });

describe('usePersonaIntroductions', () => {
  const projectId = 'test-project-1';

  const discoveryMessages: Message[] = [
    {
      id: 'msg-1',
      projectId,
      role: 'assistant',
      content: 'Welcome! Tell me about your project.',
      timestamp: '2025-12-24T10:00:00Z',
      agentType: 'product_manager',
    },
    {
      id: 'msg-2',
      projectId,
      role: 'user',
      content: 'I want to build an app.',
      timestamp: '2025-12-24T10:01:00Z',
    },
  ];

  const buildingMessages: Message[] = [
    ...discoveryMessages,
    {
      id: 'msg-3',
      projectId,
      role: 'assistant',
      content: "I'll design the interface for you.",
      timestamp: '2025-12-24T10:02:00Z',
      agentType: 'designer',
    },
  ];

  beforeEach(() => {
    mockLocalStorage.clear();
    jest.clearAllMocks();
  });

  it('should not show introductions during discovery', () => {
    const { result } = renderHook(() =>
      usePersonaIntroductions(projectId, 'welcome', discoveryMessages)
    );

    expect(result.current.shouldShowIntroductions).toBe(false);
    expect(result.current.introductionMessages).toHaveLength(0);
  });

  it('should not show introductions when discovery is complete but no building messages', () => {
    const { result } = renderHook(() =>
      usePersonaIntroductions(projectId, 'complete', discoveryMessages)
    );

    expect(result.current.shouldShowIntroductions).toBe(false);
  });

  it('should show introductions when discovery is complete and building starts', () => {
    const { result } = renderHook(() =>
      usePersonaIntroductions(projectId, 'complete', buildingMessages)
    );

    expect(result.current.shouldShowIntroductions).toBe(true);
    expect(result.current.introductionMessages.length).toBeGreaterThan(0);
  });

  it('should generate correct introduction messages', () => {
    const { result } = renderHook(() =>
      usePersonaIntroductions(projectId, 'complete', buildingMessages)
    );

    const introMessages = result.current.introductionMessages;

    // Should have Root's team intro + designer intro + developer intro
    expect(introMessages).toHaveLength(3);

    // First message is Root's team introduction
    expect(introMessages[0].agentType).toBe('product_manager');
    expect(introMessages[0].content).toContain('introduce the team');

    // Second is designer intro
    expect(introMessages[1].agentType).toBe('designer');
    expect(introMessages[1].content).toContain('design');

    // Third is developer intro
    expect(introMessages[2].agentType).toBe('developer');
    expect(introMessages[2].content).toContain('build');
  });

  it('should persist introduction state to localStorage', () => {
    const { result } = renderHook(() =>
      usePersonaIntroductions(projectId, 'complete', buildingMessages)
    );

    act(() => {
      result.current.markIntroductionsShown();
    });

    expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
      `team_introduced_${projectId}`,
      'true'
    );
  });

  it('should not show introductions after they have been shown', () => {
    // First, mark as shown
    mockLocalStorage.setItem(`team_introduced_${projectId}`, 'true');

    const { result } = renderHook(() =>
      usePersonaIntroductions(projectId, 'complete', buildingMessages)
    );

    expect(result.current.shouldShowIntroductions).toBe(false);
    expect(result.current.hasIntroducedTeam).toBe(true);
  });

  it('should process messages with introductions injected', () => {
    const { result } = renderHook(() =>
      usePersonaIntroductions(projectId, 'complete', buildingMessages)
    );

    const processed = result.current.processMessagesWithIntroductions(
      buildingMessages
    );

    // Original messages + 3 intro messages
    expect(processed.length).toBe(buildingMessages.length + 3);

    // Intro messages should be inserted before the first designer message
    const introStartIndex = processed.findIndex((m) =>
      m.id.startsWith('intro-')
    );
    expect(introStartIndex).toBe(2); // After the 2 discovery messages
  });

  it('should reset introductions', () => {
    const { result } = renderHook(() =>
      usePersonaIntroductions(projectId, 'complete', buildingMessages)
    );

    act(() => {
      result.current.markIntroductionsShown();
    });

    act(() => {
      result.current.resetIntroductions();
    });

    expect(mockLocalStorage.removeItem).toHaveBeenCalledWith(
      `team_introduced_${projectId}`
    );
    expect(result.current.hasIntroducedTeam).toBe(false);
  });
});
