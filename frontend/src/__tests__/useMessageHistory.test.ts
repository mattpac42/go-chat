import { renderHook, act } from '@testing-library/react';
import { useMessageHistory } from '@/hooks/useMessageHistory';

// Mock sessionStorage
const mockSessionStorage: { [key: string]: string } = {};

beforeEach(() => {
  // Clear mock storage
  Object.keys(mockSessionStorage).forEach((key) => {
    delete mockSessionStorage[key];
  });

  // Mock sessionStorage
  Object.defineProperty(window, 'sessionStorage', {
    value: {
      getItem: jest.fn((key: string) => mockSessionStorage[key] || null),
      setItem: jest.fn((key: string, value: string) => {
        mockSessionStorage[key] = value;
      }),
      removeItem: jest.fn((key: string) => {
        delete mockSessionStorage[key];
      }),
      clear: jest.fn(() => {
        Object.keys(mockSessionStorage).forEach((key) => {
          delete mockSessionStorage[key];
        });
      }),
    },
    writable: true,
  });
});

describe('useMessageHistory', () => {
  const projectId = 'test-project-123';

  describe('addToHistory', () => {
    it('adds a message to history', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      act(() => {
        result.current.addToHistory('Hello, world!');
      });

      // Navigate back to verify message was added
      const message = result.current.navigateHistory('up', '');
      expect(message).toBe('Hello, world!');
    });

    it('does not add duplicate consecutive messages', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      act(() => {
        result.current.addToHistory('Hello');
        result.current.addToHistory('Hello');
        result.current.addToHistory('Hello');
      });

      // Navigate back - should only find one "Hello"
      const first = result.current.navigateHistory('up', '');
      expect(first).toBe('Hello');

      // Navigating up again should return null (no more history)
      const second = result.current.navigateHistory('up', '');
      expect(second).toBeNull();
    });

    it('allows different consecutive messages', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      act(() => {
        result.current.addToHistory('First');
        result.current.addToHistory('Second');
        result.current.addToHistory('Third');
      });

      // Navigate back through history
      expect(result.current.navigateHistory('up', '')).toBe('Third');
      expect(result.current.navigateHistory('up', '')).toBe('Second');
      expect(result.current.navigateHistory('up', '')).toBe('First');
    });

    it('stores messages in sessionStorage with project-specific key', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      act(() => {
        result.current.addToHistory('Test message');
      });

      expect(window.sessionStorage.setItem).toHaveBeenCalledWith(
        `message-history-${projectId}`,
        expect.any(String)
      );

      const stored = JSON.parse(mockSessionStorage[`message-history-${projectId}`]);
      expect(stored).toContain('Test message');
    });

    it('limits history to 50 messages', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      // Add 60 messages
      act(() => {
        for (let i = 1; i <= 60; i++) {
          result.current.addToHistory(`Message ${i}`);
        }
      });

      // Navigate through all history - should only get 50 messages
      let count = 0;
      let message = result.current.navigateHistory('up', '');
      while (message !== null) {
        count++;
        message = result.current.navigateHistory('up', '');
      }

      expect(count).toBe(50);

      // Reset and verify the oldest messages were removed
      act(() => {
        result.current.resetNavigation();
      });

      // Most recent should be "Message 60"
      expect(result.current.navigateHistory('up', '')).toBe('Message 60');
    });

    it('does not add empty messages', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      act(() => {
        result.current.addToHistory('');
        result.current.addToHistory('   ');
      });

      // Should have no history
      const message = result.current.navigateHistory('up', '');
      expect(message).toBeNull();
    });
  });

  describe('navigateHistory', () => {
    it('returns null when history is empty', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      const message = result.current.navigateHistory('up', '');
      expect(message).toBeNull();
    });

    it('preserves draft message when navigating', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      act(() => {
        result.current.addToHistory('Sent message');
      });

      // Navigate up with a draft
      const navigatedTo = result.current.navigateHistory('up', 'My draft');
      expect(navigatedTo).toBe('Sent message');

      // Navigate down should restore the draft
      const restoredDraft = result.current.navigateHistory('down', 'Sent message');
      expect(restoredDraft).toBe('My draft');
    });

    it('returns null when navigating down past the draft', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      act(() => {
        result.current.addToHistory('Sent message');
      });

      // Navigate up, then down twice
      result.current.navigateHistory('up', 'draft');
      result.current.navigateHistory('down', 'Sent message');
      const pastDraft = result.current.navigateHistory('down', 'draft');

      expect(pastDraft).toBeNull();
    });

    it('returns null when navigating up past oldest message', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      act(() => {
        result.current.addToHistory('Only message');
      });

      // Navigate up once (gets the message)
      result.current.navigateHistory('up', '');
      // Navigate up again (should return null)
      const pastOldest = result.current.navigateHistory('up', 'Only message');

      expect(pastOldest).toBeNull();
    });

    it('navigates correctly through multiple messages', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      act(() => {
        result.current.addToHistory('First');
        result.current.addToHistory('Second');
        result.current.addToHistory('Third');
      });

      // Navigate up through history
      expect(result.current.navigateHistory('up', 'draft')).toBe('Third');
      expect(result.current.navigateHistory('up', 'Third')).toBe('Second');
      expect(result.current.navigateHistory('up', 'Second')).toBe('First');
      expect(result.current.navigateHistory('up', 'First')).toBeNull(); // past oldest

      // Navigate back down
      expect(result.current.navigateHistory('down', 'First')).toBe('Second');
      expect(result.current.navigateHistory('down', 'Second')).toBe('Third');
      expect(result.current.navigateHistory('down', 'Third')).toBe('draft');
      expect(result.current.navigateHistory('down', 'draft')).toBeNull(); // past draft
    });
  });

  describe('resetNavigation', () => {
    it('resets navigation index to default position', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      act(() => {
        result.current.addToHistory('First');
        result.current.addToHistory('Second');
      });

      // Navigate up a couple times
      result.current.navigateHistory('up', '');
      result.current.navigateHistory('up', 'Second');

      // Reset navigation
      act(() => {
        result.current.resetNavigation();
      });

      // Should start fresh from most recent
      expect(result.current.navigateHistory('up', '')).toBe('Second');
    });

    it('clears the stored draft', () => {
      const { result } = renderHook(() => useMessageHistory(projectId));

      act(() => {
        result.current.addToHistory('Message');
      });

      // Navigate with a draft
      result.current.navigateHistory('up', 'My important draft');

      // Reset navigation
      act(() => {
        result.current.resetNavigation();
      });

      // Navigate up and back down - draft should be empty now
      result.current.navigateHistory('up', 'new draft');
      const afterReset = result.current.navigateHistory('down', 'Message');
      expect(afterReset).toBe('new draft');
    });
  });

  describe('project isolation', () => {
    it('maintains separate history for different projects', () => {
      const { result: result1 } = renderHook(() => useMessageHistory('project-1'));
      const { result: result2 } = renderHook(() => useMessageHistory('project-2'));

      act(() => {
        result1.current.addToHistory('Project 1 message');
        result2.current.addToHistory('Project 2 message');
      });

      expect(result1.current.navigateHistory('up', '')).toBe('Project 1 message');
      expect(result2.current.navigateHistory('up', '')).toBe('Project 2 message');
    });
  });

  describe('sessionStorage persistence', () => {
    it('loads existing history from sessionStorage on mount', () => {
      // Pre-populate sessionStorage
      mockSessionStorage[`message-history-${projectId}`] = JSON.stringify([
        'Existing message 1',
        'Existing message 2',
      ]);

      const { result } = renderHook(() => useMessageHistory(projectId));

      // Should be able to navigate through existing history
      expect(result.current.navigateHistory('up', '')).toBe('Existing message 2');
      expect(result.current.navigateHistory('up', 'Existing message 2')).toBe('Existing message 1');
    });

    it('handles corrupted sessionStorage gracefully', () => {
      // Pre-populate with invalid JSON
      mockSessionStorage[`message-history-${projectId}`] = 'not valid json';

      const { result } = renderHook(() => useMessageHistory(projectId));

      // Should not throw and should start with empty history
      expect(result.current.navigateHistory('up', '')).toBeNull();

      // Should still be able to add new messages
      act(() => {
        result.current.addToHistory('New message');
      });

      expect(result.current.navigateHistory('up', '')).toBe('New message');
    });
  });
});
