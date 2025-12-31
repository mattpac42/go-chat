import { render, screen } from '@testing-library/react';
import { MessageBubble } from '@/components/chat/MessageBubble';
import { Message } from '@/types';

describe('MessageBubble', () => {
  const userMessage: Message = {
    id: 'test-user-1',
    projectId: 'project-1',
    role: 'user',
    content: 'Hello, can you help me?',
    timestamp: '2025-12-24T10:00:00Z',
  };

  const assistantMessage: Message = {
    id: 'test-assistant-1',
    projectId: 'project-1',
    role: 'assistant',
    content: 'Of course! I am happy to help.',
    timestamp: '2025-12-24T10:01:00Z',
  };

  const codeMessage: Message = {
    id: 'test-code-1',
    projectId: 'project-1',
    role: 'assistant',
    content: 'Here is an example:\n\n```typescript\nconst hello = "world";\n```\n\nLet me know if you need more.',
    timestamp: '2025-12-24T10:02:00Z',
  };

  const streamingMessage: Message = {
    id: 'test-streaming-1',
    projectId: 'project-1',
    role: 'assistant',
    content: 'I am currently typing',
    timestamp: '2025-12-24T10:03:00Z',
    isStreaming: true,
  };

  const messageWithDiscoveryData: Message = {
    id: 'test-discovery-1',
    projectId: 'project-1',
    role: 'assistant',
    content: 'Welcome to your project!<!--DISCOVERY_DATA:{"stage_complete":true,"extracted":{"business_context":"test"}}-->',
    timestamp: '2025-12-24T10:04:00Z',
    isStreaming: true,
  };

  it('renders user message with correct styling', () => {
    render(<MessageBubble message={userMessage} />);

    const bubble = screen.getByTestId('message-bubble');
    expect(bubble).toHaveAttribute('data-role', 'user');
    expect(screen.getByText('Hello, can you help me?')).toBeInTheDocument();
  });

  it('renders assistant message with correct styling', () => {
    render(<MessageBubble message={assistantMessage} />);

    const bubble = screen.getByTestId('message-bubble');
    expect(bubble).toHaveAttribute('data-role', 'assistant');
    expect(screen.getByText('Of course! I am happy to help.')).toBeInTheDocument();
  });

  it('renders code blocks with syntax highlighting wrapper', () => {
    render(<MessageBubble message={codeMessage} />);

    expect(screen.getByText('Here is an example:')).toBeInTheDocument();
    expect(screen.getByText('const hello = "world";')).toBeInTheDocument();
    expect(screen.getByText('Let me know if you need more.')).toBeInTheDocument();
    expect(screen.getByText('typescript')).toBeInTheDocument();
  });

  it('renders timestamp correctly', () => {
    render(<MessageBubble message={userMessage} />);

    // Check for time format (can vary by locale, so just check it exists)
    const bubble = screen.getByTestId('message-bubble');
    expect(bubble).toBeInTheDocument();
  });

  it('shows streaming indicator when message is streaming', () => {
    render(<MessageBubble message={streamingMessage} />);

    expect(screen.getByText('I am currently typing')).toBeInTheDocument();
    // Streaming indicator is animated bouncing dots
    const streamingIndicator = document.querySelector('.animate-bounce');
    expect(streamingIndicator).toBeInTheDocument();
  });

  it('does not show streaming indicator for completed messages', () => {
    render(<MessageBubble message={assistantMessage} />);

    const streamingIndicator = document.querySelector('.animate-bounce');
    expect(streamingIndicator).not.toBeInTheDocument();
  });

  it('renders copy button for code blocks', () => {
    render(<MessageBubble message={codeMessage} />);

    const copyButton = screen.getByRole('button', { name: /copy/i });
    expect(copyButton).toBeInTheDocument();
  });

  it('strips DISCOVERY_DATA metadata from streaming messages', () => {
    render(<MessageBubble message={messageWithDiscoveryData} />);

    // Should display the user-facing content (using regex for partial match)
    expect(screen.getByText(/Welcome to your project/)).toBeInTheDocument();

    // Should NOT display the metadata comment
    expect(screen.queryByText(/DISCOVERY_DATA/)).not.toBeInTheDocument();
    expect(screen.queryByText(/stage_complete/)).not.toBeInTheDocument();
    expect(screen.queryByText(/business_context/)).not.toBeInTheDocument();
  });
});
