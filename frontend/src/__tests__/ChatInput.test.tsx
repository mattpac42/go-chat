import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ChatInput } from '@/components/chat/ChatInput';

// Mock fetch for upload tests
const mockFetch = jest.fn();
global.fetch = mockFetch;

// Mock sessionStorage for message history
const mockSessionStorage: { [key: string]: string } = {};

// Mock URL.createObjectURL and URL.revokeObjectURL
const mockCreateObjectURL = jest.fn((file: File) => `blob:http://localhost/${file.name}`);
const mockRevokeObjectURL = jest.fn();
global.URL.createObjectURL = mockCreateObjectURL;
global.URL.revokeObjectURL = mockRevokeObjectURL;

// Helper to create a mock image file
function createMockImageFile(name: string, type: string = 'image/png', size: number = 1024): File {
  const content = new Array(size).fill('x').join('');
  return new File([content], name, { type });
}

// Helper to create a paste event with an image
function createPasteEventWithImage(file: File) {
  const clipboardData = {
    items: [
      {
        type: file.type,
        kind: 'file',
        getAsFile: () => file,
      },
    ] as unknown as DataTransferItemList,
    types: ['Files'],
    getData: () => '',
    setData: () => {},
    clearData: () => {},
    files: [file] as unknown as FileList,
  } as unknown as DataTransfer;

  return {
    clipboardData,
    preventDefault: jest.fn(),
    stopPropagation: jest.fn(),
  };
}

// Helper to create a paste event with only text
function createPasteEventWithText() {
  const clipboardData = {
    items: [
      {
        type: 'text/plain',
        kind: 'string',
        getAsFile: () => null,
      },
    ] as unknown as DataTransferItemList,
    types: ['text/plain'],
    getData: () => 'pasted text',
    setData: () => {},
    clearData: () => {},
    files: [] as unknown as FileList,
  } as unknown as DataTransfer;

  return {
    clipboardData,
    preventDefault: jest.fn(),
    stopPropagation: jest.fn(),
  };
}

beforeEach(() => {
  mockFetch.mockClear();
  mockCreateObjectURL.mockClear();
  mockRevokeObjectURL.mockClear();

  Object.keys(mockSessionStorage).forEach((key) => {
    delete mockSessionStorage[key];
  });

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

describe('ChatInput', () => {
  const mockOnSend = jest.fn();
  const testProjectId = 'test-project-123';

  beforeEach(() => {
    mockOnSend.mockClear();
  });

  it('renders input and send button', () => {
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

    expect(screen.getByTestId('chat-input')).toBeInTheDocument();
    expect(screen.getByTestId('send-button')).toBeInTheDocument();
  });

  it('renders with custom placeholder', () => {
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} placeholder="Custom placeholder" />);

    expect(screen.getByPlaceholderText('Custom placeholder')).toBeInTheDocument();
  });

  it('calls onSend when send button is clicked', async () => {
    const user = userEvent.setup();
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');
    const sendButton = screen.getByTestId('send-button');

    await user.type(input, 'Hello, world!');
    await user.click(sendButton);

    expect(mockOnSend).toHaveBeenCalledWith('Hello, world!');
  });

  it('calls onSend when Enter is pressed', async () => {
    const user = userEvent.setup();
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');

    await user.type(input, 'Hello, world!{enter}');

    expect(mockOnSend).toHaveBeenCalledWith('Hello, world!');
  });

  it('does not call onSend when Shift+Enter is pressed', async () => {
    const user = userEvent.setup();
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');

    await user.type(input, 'Hello, world!');
    await user.keyboard('{Shift>}{Enter}{/Shift}');

    expect(mockOnSend).not.toHaveBeenCalled();
  });

  it('clears input after sending', async () => {
    const user = userEvent.setup();
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');

    await user.type(input, 'Hello, world!');
    await user.click(screen.getByTestId('send-button'));

    expect(input).toHaveValue('');
  });

  it('does not send empty messages', async () => {
    const user = userEvent.setup();
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

    const sendButton = screen.getByTestId('send-button');
    await user.click(sendButton);

    expect(mockOnSend).not.toHaveBeenCalled();
  });

  it('does not send whitespace-only messages', async () => {
    const user = userEvent.setup();
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');

    await user.type(input, '   ');
    await user.click(screen.getByTestId('send-button'));

    expect(mockOnSend).not.toHaveBeenCalled();
  });

  it('disables input and button when disabled prop is true', () => {
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} disabled />);

    expect(screen.getByTestId('chat-input')).toBeDisabled();
    expect(screen.getByTestId('send-button')).toBeDisabled();
  });

  it('trims message before sending', async () => {
    const user = userEvent.setup();
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');

    await user.type(input, '  Hello, world!  ');
    await user.click(screen.getByTestId('send-button'));

    expect(mockOnSend).toHaveBeenCalledWith('Hello, world!');
  });

  it('enforces character limit of 4000 characters', async () => {
    const user = userEvent.setup();
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');
    const longText = 'a'.repeat(4100);

    // Using fireEvent for this test as userEvent.type would be too slow
    fireEvent.change(input, { target: { value: longText } });

    // Should not exceed 4000 characters
    expect((input as HTMLTextAreaElement).value.length).toBeLessThanOrEqual(4000);
  });

  it('has proper accessibility attributes', () => {
    render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

    const input = screen.getByTestId('chat-input');
    const sendButton = screen.getByTestId('send-button');

    expect(input).toHaveAttribute('aria-label', 'Message input');
    expect(sendButton).toHaveAttribute('aria-label', 'Send message');
  });

  describe('message history navigation', () => {
    // Helper to simulate keydown with specific cursor position
    const simulateArrowKey = (
      input: HTMLTextAreaElement,
      key: 'ArrowUp' | 'ArrowDown',
      selectionStart: number,
      selectionEnd: number
    ) => {
      // Set cursor position
      input.setSelectionRange(selectionStart, selectionEnd);
      // Fire keydown
      fireEvent.keyDown(input, { key });
    };

    it('navigates through message history with up arrow when cursor is at position 0', async () => {
      const user = userEvent.setup();
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input') as HTMLTextAreaElement;

      // Send some messages first
      await user.type(input, 'First message{enter}');
      await user.type(input, 'Second message{enter}');

      // Press up arrow (cursor is at position 0 in empty input)
      simulateArrowKey(input, 'ArrowUp', 0, 0);

      // Should show the most recent message
      expect(input).toHaveValue('Second message');
    });

    it('navigates back through multiple messages with up arrow', async () => {
      const user = userEvent.setup();
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input') as HTMLTextAreaElement;

      // Send some messages
      await user.type(input, 'First{enter}');
      await user.type(input, 'Second{enter}');
      await user.type(input, 'Third{enter}');

      // Navigate back through history (each time cursor at 0)
      simulateArrowKey(input, 'ArrowUp', 0, 0);
      expect(input).toHaveValue('Third');

      simulateArrowKey(input, 'ArrowUp', 0, 0);
      expect(input).toHaveValue('Second');

      simulateArrowKey(input, 'ArrowUp', 0, 0);
      expect(input).toHaveValue('First');

      // Can't go further back
      simulateArrowKey(input, 'ArrowUp', 0, 0);
      expect(input).toHaveValue('First');
    });

    it('navigates forward through history with down arrow when cursor is at end', async () => {
      const user = userEvent.setup();
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input') as HTMLTextAreaElement;

      // Send some messages
      await user.type(input, 'First{enter}');
      await user.type(input, 'Second{enter}');

      // Navigate back
      simulateArrowKey(input, 'ArrowUp', 0, 0);
      expect(input).toHaveValue('Second');

      simulateArrowKey(input, 'ArrowUp', 0, 0);
      expect(input).toHaveValue('First');

      // Now navigate forward (cursor at end of current message)
      simulateArrowKey(input, 'ArrowDown', 5, 5); // "First" is 5 chars
      expect(input).toHaveValue('Second');

      simulateArrowKey(input, 'ArrowDown', 6, 6); // "Second" is 6 chars
      expect(input).toHaveValue(''); // back to empty draft
    });

    it('preserves draft message when navigating history', async () => {
      const user = userEvent.setup();
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input') as HTMLTextAreaElement;

      // Send a message first
      await user.type(input, 'Sent message{enter}');

      // Start typing a draft
      await user.type(input, 'My draft');

      // Navigate up (with cursor at position 0)
      simulateArrowKey(input, 'ArrowUp', 0, 0);
      expect(input).toHaveValue('Sent message');

      // Navigate back down (cursor at end)
      simulateArrowKey(input, 'ArrowDown', 12, 12); // "Sent message" is 12 chars
      expect(input).toHaveValue('My draft');
    });

    it('does not trigger history when cursor is not at position 0', async () => {
      const user = userEvent.setup();
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input') as HTMLTextAreaElement;

      // Send a message
      await user.type(input, 'Sent message{enter}');

      // Type something
      await user.type(input, 'Current text');

      // Up arrow with cursor in middle - should not navigate
      simulateArrowKey(input, 'ArrowUp', 5, 5);
      expect(input).toHaveValue('Current text');
    });

    it('does not trigger down navigation when cursor is not at end', async () => {
      const user = userEvent.setup();
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input') as HTMLTextAreaElement;

      // Send messages
      await user.type(input, 'First{enter}');
      await user.type(input, 'Second{enter}');

      // Navigate back
      simulateArrowKey(input, 'ArrowUp', 0, 0);
      simulateArrowKey(input, 'ArrowUp', 0, 0);
      expect(input).toHaveValue('First');

      // Down arrow with cursor in middle - should not navigate
      simulateArrowKey(input, 'ArrowDown', 2, 2);
      expect(input).toHaveValue('First');
    });

    it('stores messages in sessionStorage', async () => {
      const user = userEvent.setup();
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');

      await user.type(input, 'Test message{enter}');

      expect(window.sessionStorage.setItem).toHaveBeenCalledWith(
        `message-history-${testProjectId}`,
        expect.any(String)
      );
    });

    it('does not add duplicate consecutive messages to history', async () => {
      const user = userEvent.setup();
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input') as HTMLTextAreaElement;

      // Send the same message multiple times
      await user.type(input, 'Same message{enter}');
      await user.type(input, 'Same message{enter}');
      await user.type(input, 'Same message{enter}');

      // Navigate back - should only find one message
      simulateArrowKey(input, 'ArrowUp', 0, 0);
      expect(input).toHaveValue('Same message');

      // Can't go further back (no more history)
      simulateArrowKey(input, 'ArrowUp', 0, 0);
      expect(input).toHaveValue('Same message');
    });
  });

  describe('clipboard image paste support', () => {
    it('captures pasted images and shows preview', async () => {
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('screenshot.png', 'image/png');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      // Preview should be shown
      await waitFor(() => {
        expect(screen.getByTestId('image-preview')).toBeInTheDocument();
      });

      // Verify the preview image uses blob URL
      const previewImg = screen.getByTestId('image-preview-img') as HTMLImageElement;
      expect(previewImg.src).toContain('blob:');
    });

    it('prevents default for paste events containing images', () => {
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('screenshot.png', 'image/png');

      // Create the actual event and spy on preventDefault
      const pasteEvent = new Event('paste', { bubbles: true, cancelable: true }) as unknown as ClipboardEventInit & Event;
      const preventDefaultSpy = jest.spyOn(pasteEvent, 'preventDefault');

      // Add clipboardData
      Object.defineProperty(pasteEvent, 'clipboardData', {
        value: {
          items: [
            {
              type: imageFile.type,
              kind: 'file',
              getAsFile: () => imageFile,
            },
          ],
        },
      });

      act(() => {
        input.dispatchEvent(pasteEvent);
      });

      expect(preventDefaultSpy).toHaveBeenCalled();
    });

    it('allows normal text paste to proceed', () => {
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const pasteEvent = createPasteEventWithText();

      fireEvent.paste(input, pasteEvent);

      // Should not show image preview for text paste
      expect(screen.queryByTestId('image-preview')).not.toBeInTheDocument();
      // Should not have prevented default
      expect(pasteEvent.preventDefault).not.toHaveBeenCalled();
    });

    it('shows remove button to clear pasted image', async () => {
      const user = userEvent.setup();
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('screenshot.png', 'image/png');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      await waitFor(() => {
        expect(screen.getByTestId('image-preview-remove')).toBeInTheDocument();
      });

      // Click remove button
      await user.click(screen.getByTestId('image-preview-remove'));

      // Preview should be removed
      expect(screen.queryByTestId('image-preview')).not.toBeInTheDocument();
    });

    it('revokes object URL when image is removed', async () => {
      const user = userEvent.setup();
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('screenshot.png', 'image/png');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      await waitFor(() => {
        expect(screen.getByTestId('image-preview-remove')).toBeInTheDocument();
      });

      await user.click(screen.getByTestId('image-preview-remove'));

      // Should have revoked the blob URL
      expect(mockRevokeObjectURL).toHaveBeenCalled();
    });

    it('accepts PNG images', async () => {
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('test.png', 'image/png');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      await waitFor(() => {
        expect(screen.getByTestId('image-preview')).toBeInTheDocument();
      });
    });

    it('accepts JPEG images', async () => {
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('test.jpg', 'image/jpeg');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      await waitFor(() => {
        expect(screen.getByTestId('image-preview')).toBeInTheDocument();
      });
    });

    it('accepts GIF images', async () => {
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('test.gif', 'image/gif');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      await waitFor(() => {
        expect(screen.getByTestId('image-preview')).toBeInTheDocument();
      });
    });

    it('accepts WebP images', async () => {
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('test.webp', 'image/webp');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      await waitFor(() => {
        expect(screen.getByTestId('image-preview')).toBeInTheDocument();
      });
    });

    it('uploads image before sending message', async () => {
      const user = userEvent.setup();
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          file: { id: 'file-123', path: 'sources/screenshot.md' },
        }),
      });

      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('screenshot.png', 'image/png');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      await waitFor(() => {
        expect(screen.getByTestId('image-preview')).toBeInTheDocument();
      });

      // Type a message and send
      await user.type(input, 'Check out this screenshot');
      await user.click(screen.getByTestId('send-button'));

      // Should have called fetch with the upload endpoint
      await waitFor(() => {
        expect(mockFetch).toHaveBeenCalledWith(
          `/api/projects/${testProjectId}/upload`,
          expect.objectContaining({
            method: 'POST',
            body: expect.any(FormData),
          })
        );
      });

      // Should have sent the message
      expect(mockOnSend).toHaveBeenCalledWith('Check out this screenshot');
    });

    it('clears image preview after successful send', async () => {
      const user = userEvent.setup();
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          file: { id: 'file-123', path: 'sources/screenshot.md' },
        }),
      });

      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('screenshot.png', 'image/png');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      await waitFor(() => {
        expect(screen.getByTestId('image-preview')).toBeInTheDocument();
      });

      await user.type(input, 'Message with image');
      await user.click(screen.getByTestId('send-button'));

      // Wait for the upload and send to complete
      await waitFor(() => {
        expect(mockOnSend).toHaveBeenCalled();
      });

      // Preview should be cleared
      await waitFor(() => {
        expect(screen.queryByTestId('image-preview')).not.toBeInTheDocument();
      });
    });

    it('can send message with only an image (no text)', async () => {
      const user = userEvent.setup();
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          file: { id: 'file-123', path: 'sources/screenshot.md' },
        }),
      });

      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('screenshot.png', 'image/png');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      await waitFor(() => {
        expect(screen.getByTestId('image-preview')).toBeInTheDocument();
      });

      // Send without any text - button should be enabled because we have an image
      const sendButton = screen.getByTestId('send-button');
      expect(sendButton).not.toBeDisabled();

      await user.click(sendButton);

      // Should have uploaded the image
      await waitFor(() => {
        expect(mockFetch).toHaveBeenCalled();
      });
    });

    it('shows preview with max height of 128px', async () => {
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('screenshot.png', 'image/png');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      await waitFor(() => {
        const previewImg = screen.getByTestId('image-preview-img');
        expect(previewImg).toHaveClass('max-h-32');
      });
    });

    it('replaces existing image when new one is pasted', async () => {
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');

      // Paste first image
      const imageFile1 = createMockImageFile('first.png', 'image/png');
      fireEvent.paste(input, createPasteEventWithImage(imageFile1));

      await waitFor(() => {
        expect(screen.getByTestId('image-preview')).toBeInTheDocument();
      });

      // Paste second image
      const imageFile2 = createMockImageFile('second.png', 'image/png');
      fireEvent.paste(input, createPasteEventWithImage(imageFile2));

      // Should still only have one preview
      await waitFor(() => {
        const previews = screen.getAllByTestId('image-preview');
        expect(previews.length).toBe(1);
      });

      // Old URL should have been revoked
      expect(mockRevokeObjectURL).toHaveBeenCalled();
    });

    it('disables send button during upload', async () => {
      // Create a delayed fetch response
      let resolveUpload: () => void;
      mockFetch.mockReturnValueOnce(
        new Promise((resolve) => {
          resolveUpload = () =>
            resolve({
              ok: true,
              json: async () => ({
                file: { id: 'file-123', path: 'sources/screenshot.md' },
              }),
            });
        })
      );

      const user = userEvent.setup();
      render(<ChatInput projectId={testProjectId} onSend={mockOnSend} />);

      const input = screen.getByTestId('chat-input');
      const imageFile = createMockImageFile('screenshot.png', 'image/png');
      const pasteEvent = createPasteEventWithImage(imageFile);

      fireEvent.paste(input, pasteEvent);

      await waitFor(() => {
        expect(screen.getByTestId('image-preview')).toBeInTheDocument();
      });

      await user.type(input, 'Test message');

      // Click send
      await user.click(screen.getByTestId('send-button'));

      // Button should be disabled during upload
      await waitFor(() => {
        expect(screen.getByTestId('send-button')).toBeDisabled();
      });

      // Resolve the upload
      act(() => {
        resolveUpload!();
      });

      // Wait for completion
      await waitFor(() => {
        expect(mockOnSend).toHaveBeenCalled();
      });
    });
  });
});
