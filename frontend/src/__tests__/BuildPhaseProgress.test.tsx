import {
  detectPhaseFromMessage,
  detectPhaseWithContext,
  groupMessagesByPhase,
  BuildPhase,
} from '@/components/chat/BuildPhaseProgress';
import { Message } from '@/types';

// Helper to create a test message
function createMessage(
  role: 'user' | 'assistant',
  content: string,
  id?: string
): Message {
  return {
    id: id || Math.random().toString(36).substring(7),
    projectId: 'test-project',
    role,
    content,
    timestamp: new Date().toISOString(),
  };
}

describe('detectPhaseFromMessage', () => {
  it('detects discovery phase for generic messages', () => {
    const message = createMessage('user', 'Hello, I need help with my app');
    expect(detectPhaseFromMessage(message)).toBe('discovery');
  });

  it('detects planning phase from architecture keywords', () => {
    const message = createMessage(
      'assistant',
      'Here is the architecture design for your app'
    );
    expect(detectPhaseFromMessage(message)).toBe('planning');
  });

  it('detects building phase from code blocks', () => {
    const message = createMessage(
      'assistant',
      'Here is the implementation:\n```typescript\nconst x = 1;\n```'
    );
    expect(detectPhaseFromMessage(message)).toBe('building');
  });

  it('detects testing phase from test keywords', () => {
    const message = createMessage(
      'assistant',
      'Running the test suite now to verify the changes'
    );
    expect(detectPhaseFromMessage(message)).toBe('testing');
  });

  it('detects launch phase from deploy keywords', () => {
    const message = createMessage(
      'assistant',
      'The application has been deployed successfully'
    );
    expect(detectPhaseFromMessage(message)).toBe('launch');
  });
});

describe('detectPhaseWithContext', () => {
  it('returns detected phase for assistant messages', () => {
    const messages: Message[] = [
      createMessage('assistant', 'Here is the architecture'),
    ];
    expect(detectPhaseWithContext(messages, 0)).toBe('planning');
  });

  it('returns detected phase for user messages with keywords', () => {
    const messages: Message[] = [
      createMessage('user', 'Can you test this?'),
    ];
    expect(detectPhaseWithContext(messages, 0)).toBe('testing');
  });

  it('inherits phase from following assistant message for keyword-less user messages', () => {
    const messages: Message[] = [
      createMessage('user', 'Sounds good, go ahead'),
      createMessage('assistant', 'Creating the component:\n```jsx\ncode\n```'),
    ];
    // User message has no keywords but following assistant is in building phase
    expect(detectPhaseWithContext(messages, 0)).toBe('building');
  });

  it('skips over other user messages to find assistant', () => {
    const messages: Message[] = [
      createMessage('user', 'First message'),
      createMessage('user', 'Second message'),
      createMessage('assistant', 'Running tests now'),
    ];
    // Both user messages should inherit testing from assistant
    expect(detectPhaseWithContext(messages, 0)).toBe('testing');
    expect(detectPhaseWithContext(messages, 1)).toBe('testing');
  });

  it('returns discovery if no following assistant message', () => {
    const messages: Message[] = [
      createMessage('assistant', 'Here is the code:\n```js\ncode\n```'),
      createMessage('user', 'What else can you do?'),
    ];
    // Last user message has no following assistant
    expect(detectPhaseWithContext(messages, 1)).toBe('discovery');
  });
});

describe('groupMessagesByPhase', () => {
  it('groups messages correctly by phase keywords', () => {
    const messages: Message[] = [
      createMessage('user', 'What can you help me with?'),
      createMessage('assistant', 'I can help you build an app'),
      createMessage(
        'assistant',
        'Here is the architecture for your solution'
      ),
    ];

    const groups = groupMessagesByPhase(messages);

    // All phases should be initialized
    expect(groups.has('discovery')).toBe(true);
    expect(groups.has('planning')).toBe(true);
    expect(groups.has('building')).toBe(true);
    expect(groups.has('testing')).toBe(true);
    expect(groups.has('launch')).toBe(true);
  });

  describe('context inheritance for user messages', () => {
    it('inherits phase from following assistant message when user message has no keywords', () => {
      const messages: Message[] = [
        createMessage('user', 'Can you write the login function?'),
        createMessage(
          'assistant',
          'Here is the implementation:\n```typescript\nfunction login() {}\n```'
        ),
      ];

      const groups = groupMessagesByPhase(messages);

      // User message should be in 'building' phase (inherited from assistant response)
      // not 'discovery' (which would happen if keyword matching failed)
      const buildingMessages = groups.get('building') || [];
      const discoveryMessages = groups.get('discovery') || [];

      expect(buildingMessages.length).toBe(2);
      expect(buildingMessages[0].role).toBe('user');
      expect(buildingMessages[1].role).toBe('assistant');
      expect(discoveryMessages.length).toBe(0);
    });

    it('inherits testing phase from following assistant message', () => {
      const messages: Message[] = [
        createMessage('user', 'Please make sure it works'),
        createMessage('assistant', 'Running tests now to verify everything'),
      ];

      const groups = groupMessagesByPhase(messages);

      const testingMessages = groups.get('testing') || [];
      expect(testingMessages.length).toBe(2);
      expect(testingMessages[0].role).toBe('user');
    });

    it('inherits planning phase from following assistant message', () => {
      const messages: Message[] = [
        createMessage('user', 'What should we build first?'),
        createMessage(
          'assistant',
          'Let me describe the architecture for this'
        ),
      ];

      const groups = groupMessagesByPhase(messages);

      const planningMessages = groups.get('planning') || [];
      expect(planningMessages.length).toBe(2);
      expect(planningMessages[0].role).toBe('user');
    });

    it('handles multiple conversation turns with correct phase inheritance', () => {
      const messages: Message[] = [
        // Discovery phase
        createMessage('user', 'I want to build a todo app'),
        createMessage('assistant', 'Great! Let me understand your requirements'),
        // Planning phase
        createMessage('user', 'What do you suggest?'),
        createMessage('assistant', 'Here is the architecture I recommend'),
        // Building phase
        createMessage('user', 'Looks good, go ahead'),
        createMessage('assistant', 'Creating the component:\n```jsx\nfunction TodoList() {}\n```'),
      ];

      const groups = groupMessagesByPhase(messages);

      const discoveryMessages = groups.get('discovery') || [];
      const planningMessages = groups.get('planning') || [];
      const buildingMessages = groups.get('building') || [];

      // First two messages in discovery
      expect(discoveryMessages.length).toBe(2);
      // Next two in planning
      expect(planningMessages.length).toBe(2);
      // Last two in building
      expect(buildingMessages.length).toBe(2);
    });

    it('keeps user message in discovery if no assistant message follows', () => {
      const messages: Message[] = [
        createMessage('assistant', 'Here is the implementation:\n```js\ncode\n```'),
        createMessage('user', 'What about adding more features?'),
      ];

      const groups = groupMessagesByPhase(messages);

      const buildingMessages = groups.get('building') || [];
      const discoveryMessages = groups.get('discovery') || [];

      // Assistant message in building
      expect(buildingMessages.length).toBe(1);
      // User message at end stays in discovery (no following assistant to inherit from)
      expect(discoveryMessages.length).toBe(1);
      expect(discoveryMessages[0].role).toBe('user');
    });

    it('preserves user messages that have actual phase keywords', () => {
      const messages: Message[] = [
        createMessage('user', 'Can you test the implementation?'),
        createMessage('assistant', 'Running the test suite now'),
      ];

      const groups = groupMessagesByPhase(messages);

      const testingMessages = groups.get('testing') || [];
      // Both should be in testing - user message has 'test' keyword
      expect(testingMessages.length).toBe(2);
    });
  });
});
