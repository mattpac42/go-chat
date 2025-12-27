'use client';

import { useState, useCallback, useRef, useMemo, useEffect, KeyboardEvent, ClipboardEvent } from 'react';
import { useMessageHistory } from '@/hooks/useMessageHistory';

// Allowed MIME types for image paste (matches backend)
const ALLOWED_IMAGE_TYPES = ['image/png', 'image/jpeg', 'image/gif', 'image/webp'];

interface ChatInputProps {
  projectId: string;
  onSend: (message: string) => void;
  disabled?: boolean;
  placeholder?: string;
}

interface ImagePreviewProps {
  file: File;
  onRemove: () => void;
}

function ImagePreview({ file, onRemove }: ImagePreviewProps) {
  const previewUrl = useMemo(() => URL.createObjectURL(file), [file]);

  useEffect(() => {
    return () => {
      URL.revokeObjectURL(previewUrl);
    };
  }, [previewUrl]);

  return (
    <div className="relative inline-block" data-testid="image-preview">
      <img
        src={previewUrl}
        alt="Preview of pasted image"
        className="max-h-32 rounded border border-gray-200"
        data-testid="image-preview-img"
      />
      <button
        type="button"
        onClick={onRemove}
        className="absolute -top-2 -right-2 w-6 h-6 flex items-center justify-center rounded-full bg-gray-700 text-white hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500"
        aria-label="Remove image"
        data-testid="image-preview-remove"
      >
        <span aria-hidden="true">x</span>
      </button>
    </div>
  );
}

export function ChatInput({
  projectId,
  onSend,
  disabled = false,
  placeholder = 'Type a message...',
}: ChatInputProps) {
  const [message, setMessage] = useState('');
  const [pendingImage, setPendingImage] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const pendingImageUrlRef = useRef<string | null>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const { addToHistory, navigateHistory, resetNavigation } = useMessageHistory(projectId);

  const handleSend = useCallback(async () => {
    const trimmedMessage = message.trim();
    const hasContent = trimmedMessage || pendingImage;

    if (!hasContent || disabled || isUploading) {
      return;
    }

    setIsUploading(true);

    try {
      // If there's a pending image, upload it first
      if (pendingImage) {
        const formData = new FormData();
        formData.append('file', pendingImage);

        const response = await fetch(`/api/projects/${projectId}/upload`, {
          method: 'POST',
          body: formData,
        });

        if (!response.ok) {
          console.error('Failed to upload image');
          // Still continue with the message if upload fails
        }
      }

      // Send the message (even if empty when image was uploaded)
      if (trimmedMessage) {
        addToHistory(trimmedMessage);
        onSend(trimmedMessage);
      } else if (pendingImage) {
        // When only image, just call onSend with empty string or a placeholder
        // This triggers the chat flow without text
        onSend('');
      }

      // Clear state
      setMessage('');
      setPendingImage(null);
      if (pendingImageUrlRef.current) {
        URL.revokeObjectURL(pendingImageUrlRef.current);
        pendingImageUrlRef.current = null;
      }

      // Reset textarea height
      if (textareaRef.current) {
        textareaRef.current.style.height = 'auto';
      }
    } finally {
      setIsUploading(false);
    }
  }, [message, pendingImage, disabled, isUploading, projectId, onSend, addToHistory]);

  const handleKeyDown = useCallback(
    (e: KeyboardEvent<HTMLTextAreaElement>) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        handleSend();
        return;
      }

      // Handle up arrow for history navigation (only at start of input)
      if (e.key === 'ArrowUp') {
        const textarea = e.currentTarget;
        // Only trigger when cursor is at position 0
        if (textarea.selectionStart === 0 && textarea.selectionEnd === 0) {
          const historicMessage = navigateHistory('up', message);
          if (historicMessage !== null) {
            e.preventDefault();
            setMessage(historicMessage);
            // Move cursor to end of text
            requestAnimationFrame(() => {
              if (textareaRef.current) {
                textareaRef.current.selectionStart = historicMessage.length;
                textareaRef.current.selectionEnd = historicMessage.length;
              }
            });
          }
        }
        return;
      }

      // Handle down arrow for history navigation
      if (e.key === 'ArrowDown') {
        const textarea = e.currentTarget;
        // Only trigger when cursor is at end of input
        if (textarea.selectionStart === message.length && textarea.selectionEnd === message.length) {
          const historicMessage = navigateHistory('down', message);
          if (historicMessage !== null) {
            e.preventDefault();
            setMessage(historicMessage);
            // Move cursor to end of text
            requestAnimationFrame(() => {
              if (textareaRef.current) {
                textareaRef.current.selectionStart = historicMessage.length;
                textareaRef.current.selectionEnd = historicMessage.length;
              }
            });
          }
        }
      }
    },
    [handleSend, message, navigateHistory]
  );

  const handleChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const value = e.target.value;
    // Limit to 4000 characters as per PRD
    if (value.length <= 4000) {
      setMessage(value);
      // Reset history navigation when user types
      resetNavigation();
      // Auto-resize textarea
      if (textareaRef.current) {
        textareaRef.current.style.height = 'auto';
        textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 120)}px`;
      }
    }
  }, [resetNavigation]);

  const handlePaste = useCallback((e: ClipboardEvent<HTMLTextAreaElement>) => {
    const items = e.clipboardData.items;
    for (let i = 0; i < items.length; i++) {
      const item = items[i];
      if (item.type.indexOf('image') !== -1 && ALLOWED_IMAGE_TYPES.includes(item.type)) {
        e.preventDefault();
        const file = item.getAsFile();
        if (file) {
          // Revoke previous URL if exists
          if (pendingImageUrlRef.current) {
            URL.revokeObjectURL(pendingImageUrlRef.current);
          }
          setPendingImage(file);
          // Store the URL for later cleanup
          pendingImageUrlRef.current = URL.createObjectURL(file);
        }
        return;
      }
    }
    // If no image found, let the default paste behavior happen
  }, []);

  const handleRemoveImage = useCallback(() => {
    if (pendingImageUrlRef.current) {
      URL.revokeObjectURL(pendingImageUrlRef.current);
      pendingImageUrlRef.current = null;
    }
    setPendingImage(null);
  }, []);

  // Determine if send button should be enabled
  const canSend = !disabled && !isUploading && (message.trim() || pendingImage);

  return (
    <div className="flex flex-col gap-2 p-4 pb-5 border-t border-gray-200 bg-white safe-area-pb">
      {/* Image preview area */}
      {pendingImage && (
        <div className="px-1">
          <ImagePreview file={pendingImage} onRemove={handleRemoveImage} />
        </div>
      )}

      {/* Input area */}
      <div className="flex items-end gap-2">
        <textarea
          ref={textareaRef}
          value={message}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
          onPaste={handlePaste}
          disabled={disabled || isUploading}
          placeholder={placeholder}
          rows={1}
          className="flex-1 resize-none rounded-xl border-2 border-gray-300 px-4 py-3 text-base focus:outline-none focus:border-teal-400 disabled:bg-gray-100 disabled:cursor-not-allowed min-h-[48px] max-h-[120px]"
          aria-label="Message input"
          data-testid="chat-input"
        />
        <button
          onClick={handleSend}
          disabled={!canSend}
          className="flex items-center justify-center w-12 h-12 rounded-full bg-teal-400 text-white hover:bg-teal-500 focus:outline-none focus:ring-2 focus:ring-teal-400 focus:ring-offset-2 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
          aria-label="Send message"
          data-testid="send-button"
        >
          <SendIcon className="w-5 h-5" />
        </button>
      </div>
    </div>
  );
}

function SendIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
      />
    </svg>
  );
}
