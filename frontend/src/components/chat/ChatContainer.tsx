'use client';

import { useEffect, useRef } from 'react';
import { useChat } from '@/hooks/useChat';
import { MessageList } from './MessageList';
import { ChatInput } from './ChatInput';
import { ConnectionStatus } from '@/components/shared/ConnectionStatus';
import { Message } from '@/types';

interface ChatContainerProps {
  projectId: string;
  projectTitle?: string;
  initialMessages?: Message[];
  onMenuClick?: () => void;
  onStreamingComplete?: () => void;
}

export function ChatContainer({
  projectId,
  projectTitle = 'New Project',
  initialMessages = [],
  onMenuClick,
  onStreamingComplete,
}: ChatContainerProps) {
  const {
    messages,
    isLoading,
    error,
    connectionStatus,
    reconnectAttempts,
    sendMessage,
    clearError,
    reconnect,
  } = useChat({
    projectId,
    initialMessages,
  });

  // Track previous loading state to detect when streaming completes
  const wasLoadingRef = useRef(false);

  useEffect(() => {
    // Detect transition from loading to not loading (streaming complete)
    if (wasLoadingRef.current && !isLoading) {
      onStreamingComplete?.();
    }
    wasLoadingRef.current = isLoading;
  }, [isLoading, onStreamingComplete]);

  return (
    <div className="flex flex-col h-full bg-white">
      {/* Header */}
      <header className="flex items-center justify-between px-4 py-3 border-b border-gray-200 bg-white sticky top-0 z-10 safe-area-pt">
        <div className="flex items-center gap-3">
          {onMenuClick && (
            <button
              onClick={onMenuClick}
              className="p-2 -ml-2 rounded-lg hover:bg-gray-100 md:hidden"
              aria-label="Open menu"
            >
              <MenuIcon className="w-6 h-6 text-gray-600" />
            </button>
          )}
          <h1 className="text-lg font-semibold text-gray-900 truncate">
            {projectTitle}
          </h1>
        </div>
        <ConnectionStatus
          status={connectionStatus}
          reconnectAttempts={reconnectAttempts}
          onReconnect={reconnect}
        />
      </header>

      {/* Error banner */}
      {error && (
        <div className="flex items-center justify-between px-4 py-2 bg-red-50 border-b border-red-100">
          <p className="text-sm text-red-600">{error}</p>
          <button
            onClick={clearError}
            className="text-sm text-red-600 hover:text-red-800 font-medium"
          >
            Dismiss
          </button>
        </div>
      )}

      {/* Messages */}
      <MessageList messages={messages} isLoading={isLoading} />

      {/* Input */}
      <ChatInput
        onSend={sendMessage}
        disabled={connectionStatus !== 'connected' || isLoading}
        placeholder={
          connectionStatus !== 'connected'
            ? 'Connecting...'
            : isLoading
            ? 'Waiting for response...'
            : 'Describe what you want to build...'
        }
      />

      {/* Bottom spacer for visual breathing room */}
      <div className="h-6 bg-gray-50 flex-shrink-0" />
    </div>
  );
}

function MenuIcon({ className }: { className?: string }) {
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
        d="M4 6h16M4 12h16M4 18h16"
      />
    </svg>
  );
}
