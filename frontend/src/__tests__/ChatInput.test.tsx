import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ChatInput } from '@/components/chat/ChatInput';

describe('ChatInput', () => {
  const mockOnSend = jest.fn();

  beforeEach(() => {
    mockOnSend.mockClear();
  });

  it('renders input and send button', () => {
    render(<ChatInput onSend={mockOnSend} />);

    expect(screen.getByTestId('chat-input')).toBeInTheDocument();
    expect(screen.getByTestId('send-button')).toBeInTheDocument();
  });

  it('renders with custom placeholder', () => {
    render(<ChatInput onSend={mockOnSend} placeholder="Custom placeholder" />);

    expect(screen.getByPlaceholderText('Custom placeholder')).toBeInTheDocument();
  });

  it('calls onSend when send button is clicked', async () => {
    const user = userEvent.setup();
    render(<ChatInput onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');
    const sendButton = screen.getByTestId('send-button');

    await user.type(input, 'Hello, world!');
    await user.click(sendButton);

    expect(mockOnSend).toHaveBeenCalledWith('Hello, world!');
  });

  it('calls onSend when Enter is pressed', async () => {
    const user = userEvent.setup();
    render(<ChatInput onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');

    await user.type(input, 'Hello, world!{enter}');

    expect(mockOnSend).toHaveBeenCalledWith('Hello, world!');
  });

  it('does not call onSend when Shift+Enter is pressed', async () => {
    const user = userEvent.setup();
    render(<ChatInput onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');

    await user.type(input, 'Hello, world!');
    await user.keyboard('{Shift>}{Enter}{/Shift}');

    expect(mockOnSend).not.toHaveBeenCalled();
  });

  it('clears input after sending', async () => {
    const user = userEvent.setup();
    render(<ChatInput onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');

    await user.type(input, 'Hello, world!');
    await user.click(screen.getByTestId('send-button'));

    expect(input).toHaveValue('');
  });

  it('does not send empty messages', async () => {
    const user = userEvent.setup();
    render(<ChatInput onSend={mockOnSend} />);

    const sendButton = screen.getByTestId('send-button');
    await user.click(sendButton);

    expect(mockOnSend).not.toHaveBeenCalled();
  });

  it('does not send whitespace-only messages', async () => {
    const user = userEvent.setup();
    render(<ChatInput onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');

    await user.type(input, '   ');
    await user.click(screen.getByTestId('send-button'));

    expect(mockOnSend).not.toHaveBeenCalled();
  });

  it('disables input and button when disabled prop is true', () => {
    render(<ChatInput onSend={mockOnSend} disabled />);

    expect(screen.getByTestId('chat-input')).toBeDisabled();
    expect(screen.getByTestId('send-button')).toBeDisabled();
  });

  it('trims message before sending', async () => {
    const user = userEvent.setup();
    render(<ChatInput onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');

    await user.type(input, '  Hello, world!  ');
    await user.click(screen.getByTestId('send-button'));

    expect(mockOnSend).toHaveBeenCalledWith('Hello, world!');
  });

  it('enforces character limit of 4000 characters', async () => {
    const user = userEvent.setup();
    render(<ChatInput onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');
    const longText = 'a'.repeat(4100);

    // Using fireEvent for this test as userEvent.type would be too slow
    fireEvent.change(input, { target: { value: longText } });

    // Should not exceed 4000 characters
    expect((input as HTMLTextAreaElement).value.length).toBeLessThanOrEqual(4000);
  });

  it('has proper accessibility attributes', () => {
    render(<ChatInput onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');
    const sendButton = screen.getByTestId('send-button');

    expect(input).toHaveAttribute('aria-label', 'Message input');
    expect(sendButton).toHaveAttribute('aria-label', 'Send message');
  });
});
