import {
  detectPhaseFromMessage,
  detectPhaseWithContext,
  groupMessagesByPhase,
  detectPhaseMarker,
  assignPhasesToMessages,
  shouldShowPhaseToggle,
  MIN_MESSAGES_FOR_PHASE_VIEW,
  MIN_MESSAGES_PER_SECTION,
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

describe('detectPhaseMarker', () => {
  it('returns null for messages without phase markers', () => {
    expect(detectPhaseMarker('Hello, this is a regular message')).toBeNull();
    expect(detectPhaseMarker('Let me help you build something')).toBeNull();
    expect(detectPhaseMarker('Here is some code: ```js\ncode\n```')).toBeNull();
  });

  it('detects [Beginning X phase] markers', () => {
    expect(detectPhaseMarker('[Beginning planning phase]')).toBe('planning');
    expect(detectPhaseMarker('[Beginning building phase]')).toBe('building');
    expect(detectPhaseMarker('[Beginning testing phase]')).toBe('testing');
    expect(detectPhaseMarker('[Beginning launch phase]')).toBe('launch');
  });

  it('detects [Entering X phase] markers', () => {
    expect(detectPhaseMarker('[Entering planning phase]')).toBe('planning');
    expect(detectPhaseMarker('[Entering building phase]')).toBe('building');
    expect(detectPhaseMarker('[Entering testing phase]')).toBe('testing');
    expect(detectPhaseMarker('[Entering launch phase]')).toBe('launch');
  });

  it('detects [Moving to X phase] markers', () => {
    expect(detectPhaseMarker('[Moving to planning phase]')).toBe('planning');
    expect(detectPhaseMarker('[Moving to building phase]')).toBe('building');
    expect(detectPhaseMarker('[Moving to testing phase]')).toBe('testing');
    expect(detectPhaseMarker('[Moving to launch phase]')).toBe('launch');
  });

  it('detects [Starting X phase] markers', () => {
    expect(detectPhaseMarker('[Starting planning phase]')).toBe('planning');
    expect(detectPhaseMarker('[Starting building phase]')).toBe('building');
    expect(detectPhaseMarker('[Starting testing phase]')).toBe('testing');
    expect(detectPhaseMarker('[Starting launch phase]')).toBe('launch');
  });

  it('is case-insensitive', () => {
    expect(detectPhaseMarker('[BEGINNING PLANNING PHASE]')).toBe('planning');
    expect(detectPhaseMarker('[Beginning Planning Phase]')).toBe('planning');
    expect(detectPhaseMarker('[beginning BUILDING phase]')).toBe('building');
  });

  it('detects markers embedded in larger content', () => {
    const content = 'Great! Now let me start coding.\n\n[Beginning building phase]\n\nHere is the implementation...';
    expect(detectPhaseMarker(content)).toBe('building');
  });

  it('returns null for discovery phase markers (not valid post-discovery)', () => {
    expect(detectPhaseMarker('[Beginning discovery phase]')).toBeNull();
  });
});

describe('assignPhasesToMessages - Sticky Phase Logic', () => {
  it('assigns discovery to all messages when discovery not complete', () => {
    const messages: Message[] = [
      createMessage('user', 'Hello'),
      createMessage('assistant', 'Hi there!'),
      createMessage('user', 'What can you build?'),
    ];

    const phases = assignPhasesToMessages(messages, -1);
    expect(phases).toEqual(['discovery', 'discovery', 'discovery']);
  });

  it('assigns discovery to all messages up to and including discoveryCompleteIndex', () => {
    const messages: Message[] = [
      createMessage('user', 'Hello'),
      createMessage('assistant', 'Hi there!'),
      createMessage('user', 'Ready to build'),
      createMessage('assistant', 'Let me plan this out'),
    ];

    // Discovery completed at index 1 (second message)
    const phases = assignPhasesToMessages(messages, 1);
    expect(phases[0]).toBe('discovery');
    expect(phases[1]).toBe('discovery');
    // Messages after discovery complete default to planning
    expect(phases[2]).toBe('planning');
    expect(phases[3]).toBe('planning');
  });

  it('defaults to planning for first post-discovery message without marker', () => {
    const messages: Message[] = [
      createMessage('user', 'Define my app'),
      createMessage('assistant', 'I understand'),
      createMessage('user', 'Lets go'),
      createMessage('assistant', 'Starting now'),
    ];

    const phases = assignPhasesToMessages(messages, 1);
    // Post-discovery messages default to planning
    expect(phases[2]).toBe('planning');
    expect(phases[3]).toBe('planning');
  });

  it('uses phase marker when present', () => {
    const messages: Message[] = [
      createMessage('user', 'Define my app'),
      createMessage('assistant', 'I understand'),
      createMessage('assistant', '[Beginning building phase]\n\nHere is code'),
      createMessage('user', 'Looks good'),
    ];

    const phases = assignPhasesToMessages(messages, 1);
    expect(phases[0]).toBe('discovery');
    expect(phases[1]).toBe('discovery');
    expect(phases[2]).toBe('building');
    expect(phases[3]).toBe('building'); // Sticky - inherits from previous
  });

  it('implements sticky phase inheritance correctly', () => {
    const messages: Message[] = [
      createMessage('user', 'Requirements'),
      createMessage('assistant', 'Discovery complete'),
      createMessage('assistant', '[Beginning planning phase]'),
      createMessage('user', 'Good plan'),
      createMessage('assistant', 'Thanks'),
      createMessage('assistant', '[Beginning building phase]'),
      createMessage('user', 'Nice code'),
      createMessage('assistant', 'More code'),
      createMessage('user', 'Test it'),
      createMessage('assistant', '[Beginning testing phase]'),
      createMessage('user', 'All tests pass'),
    ];

    const phases = assignPhasesToMessages(messages, 1);

    // Discovery phase
    expect(phases[0]).toBe('discovery');
    expect(phases[1]).toBe('discovery');
    // Planning phase (marked + sticky)
    expect(phases[2]).toBe('planning');
    expect(phases[3]).toBe('planning');
    expect(phases[4]).toBe('planning');
    // Building phase (marked + sticky)
    expect(phases[5]).toBe('building');
    expect(phases[6]).toBe('building');
    expect(phases[7]).toBe('building');
    expect(phases[8]).toBe('building');
    // Testing phase (marked + sticky)
    expect(phases[9]).toBe('testing');
    expect(phases[10]).toBe('testing');
  });

  it('handles transition between all phases', () => {
    const messages: Message[] = [
      createMessage('assistant', 'Welcome'),
      createMessage('assistant', '[Beginning planning phase]'),
      createMessage('assistant', '[Beginning building phase]'),
      createMessage('assistant', '[Beginning testing phase]'),
      createMessage('assistant', '[Beginning launch phase]'),
    ];

    const phases = assignPhasesToMessages(messages, 0);

    expect(phases[0]).toBe('discovery');
    expect(phases[1]).toBe('planning');
    expect(phases[2]).toBe('building');
    expect(phases[3]).toBe('testing');
    expect(phases[4]).toBe('launch');
  });
});

describe('groupMessagesByPhase with sticky phases', () => {
  it('groups all messages as discovery when discovery not complete', () => {
    const messages: Message[] = [
      createMessage('user', 'Hello'),
      createMessage('assistant', 'Hi'),
      createMessage('user', 'Help me'),
    ];

    const groups = groupMessagesByPhase(messages, -1);
    const discoveryMessages = groups.get('discovery') || [];

    expect(discoveryMessages.length).toBe(3);
    expect(groups.get('planning')?.length || 0).toBe(0);
  });

  it('separates discovery from post-discovery phases', () => {
    const messages: Message[] = [
      createMessage('user', 'Requirements'),
      createMessage('assistant', 'Understood'),
      createMessage('user', 'Start building'),
      createMessage('assistant', '[Beginning building phase]\n```code```'),
    ];

    const groups = groupMessagesByPhase(messages, 1);

    // Discovery has 2 messages
    // Post-discovery: 1 user message (planning by default) + 1 building marker message
    // Planning has only 1 message, so it gets merged into discovery (min section size = 2)
    // Building has only 1 message, so it gets merged into the previous non-empty phase
    // Result: discovery gets planning merged, then building gets merged
    // This results in 4 messages in discovery
    expect(groups.get('discovery')?.length).toBe(4);
    expect(groups.get('planning')?.length).toBe(0);
    expect(groups.get('building')?.length).toBe(0);
  });

  it('uses sticky inheritance for grouping', () => {
    const messages: Message[] = [
      createMessage('assistant', 'Discovery done'),
      createMessage('assistant', '[Beginning building phase]'),
      createMessage('user', 'Continue'),
      createMessage('assistant', 'More code'),
      createMessage('user', 'And more'),
    ];

    const groups = groupMessagesByPhase(messages, 0);

    const buildingMessages = groups.get('building') || [];
    // 4 messages after discovery: all should be in building (marker + sticky)
    expect(buildingMessages.length).toBe(4);
  });
});

describe('shouldShowPhaseToggle', () => {
  it('returns false when discovery not complete', () => {
    expect(shouldShowPhaseToggle(false, 15)).toBe(false);
    expect(shouldShowPhaseToggle(false, 100)).toBe(false);
  });

  it('returns false when fewer than MIN_MESSAGES_FOR_PHASE_VIEW messages', () => {
    expect(shouldShowPhaseToggle(true, 0)).toBe(false);
    expect(shouldShowPhaseToggle(true, 5)).toBe(false);
    expect(shouldShowPhaseToggle(true, MIN_MESSAGES_FOR_PHASE_VIEW - 1)).toBe(false);
  });

  it('returns true when discovery complete AND enough messages', () => {
    expect(shouldShowPhaseToggle(true, MIN_MESSAGES_FOR_PHASE_VIEW)).toBe(true);
    expect(shouldShowPhaseToggle(true, MIN_MESSAGES_FOR_PHASE_VIEW + 1)).toBe(true);
    expect(shouldShowPhaseToggle(true, 50)).toBe(true);
  });

  it('uses correct constant for minimum messages (10)', () => {
    expect(MIN_MESSAGES_FOR_PHASE_VIEW).toBe(10);
    expect(shouldShowPhaseToggle(true, 9)).toBe(false);
    expect(shouldShowPhaseToggle(true, 10)).toBe(true);
  });
});

describe('minimum section size enforcement', () => {
  it('exports MIN_MESSAGES_PER_SECTION constant as 2', () => {
    expect(MIN_MESSAGES_PER_SECTION).toBe(2);
  });

  it('merges small sections into previous phase', () => {
    const messages: Message[] = [
      // Discovery: 3 messages
      createMessage('user', 'Req 1'),
      createMessage('assistant', 'Got it'),
      createMessage('assistant', 'Summary'),
      // Planning: 1 message (should merge into discovery)
      createMessage('assistant', '[Beginning planning phase]'),
      // Building: 3 messages
      createMessage('assistant', '[Beginning building phase]'),
      createMessage('user', 'Looks good'),
      createMessage('assistant', 'More code'),
    ];

    const groups = groupMessagesByPhase(messages, 2);

    // Planning has only 1 message, should be merged into discovery
    const discoveryMessages = groups.get('discovery') || [];
    const planningMessages = groups.get('planning') || [];
    const buildingMessages = groups.get('building') || [];

    // Discovery (3) + Planning merged (1) = 4
    expect(discoveryMessages.length).toBe(4);
    expect(planningMessages.length).toBe(0);
    expect(buildingMessages.length).toBe(3);
  });

  it('keeps sections with 2 or more messages separate', () => {
    const messages: Message[] = [
      createMessage('user', 'Req'),
      createMessage('assistant', 'Got it'),
      createMessage('assistant', '[Beginning planning phase]'),
      createMessage('user', 'Plan question'),
      createMessage('assistant', '[Beginning building phase]'),
      createMessage('user', 'Build'),
    ];

    const groups = groupMessagesByPhase(messages, 1);

    const planningMessages = groups.get('planning') || [];
    const buildingMessages = groups.get('building') || [];

    // Planning has 2 messages, building has 2 messages - both should stay
    expect(planningMessages.length).toBe(2);
    expect(buildingMessages.length).toBe(2);
  });
});

// Legacy tests - kept for backwards compatibility with deprecated functions
describe('detectPhaseFromMessage (deprecated)', () => {
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

describe('detectPhaseWithContext (deprecated)', () => {
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
    expect(detectPhaseWithContext(messages, 0)).toBe('building');
  });

  it('returns discovery if no following assistant message', () => {
    const messages: Message[] = [
      createMessage('assistant', 'Here is the code:\n```js\ncode\n```'),
      createMessage('user', 'What else can you do?'),
    ];
    expect(detectPhaseWithContext(messages, 1)).toBe('discovery');
  });
});
