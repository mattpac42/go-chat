'use client';

import { useState, useCallback, useRef, KeyboardEvent, ClipboardEvent, DragEvent, ChangeEvent } from 'react';
import { useMessageHistory } from '@/hooks/useMessageHistory';

// Allowed MIME types for image upload (matches backend)
const ALLOWED_IMAGE_TYPES = ['image/png', 'image/jpeg', 'image/gif', 'image/webp'];
const ALLOWED_EXTENSIONS = '.png,.jpg,.jpeg,.gif,.webp';

interface ChatInputProps {
  projectId: string;
  onSend: (message: string) => void;
  disabled?: boolean;
  placeholder?: string;
}

interface ImagePreviewProps {
  previewUrl: string;
  onRemove: () => void;
}

function ImagePreview({ previewUrl, onRemove }: ImagePreviewProps) {
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
  const [pendingImageUrl, setPendingImageUrl] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [isDragOver, setIsDragOver] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { addToHistory, navigateHistory, resetNavigation } = useMessageHistory(projectId);

  // Helper to set a pending image file
  const setImageFile = useCallback((file: File) => {
    if (!ALLOWED_IMAGE_TYPES.includes(file.type)) {
      console.warn('Unsupported file type:', file.type);
      return;
    }
    // Revoke previous URL if exists
    if (pendingImageUrl) {
      URL.revokeObjectURL(pendingImageUrl);
    }
    setPendingImage(file);
    setPendingImageUrl(URL.createObjectURL(file));
  }, [pendingImageUrl]);

  const handleSend = useCallback(async () => {
    const trimmedMessage = message.trim();
    const hasContent = trimmedMessage || pendingImage;

    if (!hasContent || disabled || isUploading) {
      return;
    }

    setIsUploading(true);

    try {
      let imageContent = '';

      // If there's a pending image, upload it first
      if (pendingImage) {
        try {
          const formData = new FormData();
          formData.append('file', pendingImage);

          const response = await fetch(`/api/projects/${projectId}/upload`, {
            method: 'POST',
            body: formData,
          });

          if (response.ok) {
            // Parse the upload response to get the image description
            const data = await response.json();
            if (data.file?.content) {
              imageContent = data.file.content;
            }
          } else {
            console.error('Failed to upload image:', response.status, response.statusText);
          }
        } catch (uploadError) {
          console.error('Image upload failed:', uploadError);
          // Continue without image content - user can retry
        }
      }

      // Build the final message with image content if available
      let finalMessage = trimmedMessage;
      if (imageContent) {
        if (trimmedMessage) {
          // User provided a message with the image
          finalMessage = `${trimmedMessage}\n\n---\n**Uploaded Image:**\n${imageContent}`;
        } else {
          // Image only, no text message
          finalMessage = `[Uploaded Image]\n\n${imageContent}`;
        }
      }

      // Send the message if we have content
      if (finalMessage) {
        addToHistory(trimmedMessage || '[Image]');
        onSend(finalMessage);
      }

      // Clear state
      setMessage('');
      setPendingImage(null);
      if (pendingImageUrl) {
        URL.revokeObjectURL(pendingImageUrl);
        setPendingImageUrl(null);
      }

      // Reset textarea height
      if (textareaRef.current) {
        textareaRef.current.style.height = 'auto';
      }
    } finally {
      setIsUploading(false);
    }
  }, [message, pendingImage, pendingImageUrl, disabled, isUploading, projectId, onSend, addToHistory]);

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
          setImageFile(file);
        }
        return;
      }
    }
    // If no image found, let the default paste behavior happen
  }, [setImageFile]);

  const handleRemoveImage = useCallback(() => {
    if (pendingImageUrl) {
      URL.revokeObjectURL(pendingImageUrl);
      setPendingImageUrl(null);
    }
    setPendingImage(null);
  }, [pendingImageUrl]);

  // Handle file input change (from file picker)
  const handleFileChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setImageFile(file);
    }
    // Reset input so same file can be selected again
    e.target.value = '';
  }, [setImageFile]);

  // Handle attachment button click
  const handleAttachClick = useCallback(() => {
    fileInputRef.current?.click();
  }, []);

  // Drag and drop handlers
  const handleDragOver = useCallback((e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.dataTransfer.types.includes('Files')) {
      setIsDragOver(true);
    }
  }, []);

  const handleDragLeave = useCallback((e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragOver(false);
  }, []);

  const handleDrop = useCallback((e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragOver(false);

    const files = e.dataTransfer.files;
    if (files.length > 0) {
      const file = files[0];
      if (ALLOWED_IMAGE_TYPES.includes(file.type)) {
        setImageFile(file);
      } else {
        console.warn('Unsupported file type:', file.type);
      }
    }
  }, [setImageFile]);

  // Determine if send button should be enabled
  const canSend = !disabled && !isUploading && (message.trim() || pendingImage);

  return (
    <div
      className={`flex flex-col gap-2 p-4 pb-5 border-t border-gray-200 bg-white safe-area-pb transition-colors ${
        isDragOver ? 'bg-teal-50 border-t-teal-400' : ''
      }`}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      {/* Hidden file input */}
      <input
        ref={fileInputRef}
        type="file"
        accept={ALLOWED_EXTENSIONS}
        onChange={handleFileChange}
        className="hidden"
        aria-hidden="true"
        data-testid="file-input"
      />

      {/* Drag overlay indicator */}
      {isDragOver && (
        <div className="flex items-center justify-center py-4 text-teal-600 font-medium">
          <PaperclipIcon className="w-5 h-5 mr-2" />
          Drop image here
        </div>
      )}

      {/* Image preview area */}
      {pendingImageUrl && !isDragOver && (
        <div className="px-1">
          <ImagePreview previewUrl={pendingImageUrl} onRemove={handleRemoveImage} />
        </div>
      )}

      {/* Input area */}
      <div className="flex items-end gap-2">
        {/* Attachment button */}
        <button
          type="button"
          onClick={handleAttachClick}
          disabled={disabled || isUploading}
          className="flex items-center justify-center w-10 h-10 rounded-full text-gray-500 hover:text-gray-700 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-teal-400 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          aria-label="Attach image"
          data-testid="attach-button"
        >
          <PaperclipIcon className="w-5 h-5" />
        </button>

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

function PaperclipIcon({ className }: { className?: string }) {
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
        d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"
      />
    </svg>
  );
}
