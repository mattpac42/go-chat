'use client';

import { useEffect, useRef, useLayoutEffect } from 'react';
import { Message } from '@/types';
import { MessageBubble } from './MessageBubble';

interface MessageListProps {
  messages: Message[];
  isLoading?: boolean;
}

export function MessageList({ messages, isLoading = false }: MessageListProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);
  const isInitialLoad = useRef(true);
  const prevMessageCount = useRef(0);

  // Instant scroll to bottom on initial load (no animation)
  useLayoutEffect(() => {
    if (messages.length > 0 && isInitialLoad.current && containerRef.current) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
      isInitialLoad.current = false;
      prevMessageCount.current = messages.length;
    }
  }, [messages]);

  // Smooth scroll for new messages only (after initial load)
  useEffect(() => {
    // Skip if this is initial load or no new messages
    if (isInitialLoad.current || messages.length <= prevMessageCount.current) {
      prevMessageCount.current = messages.length;
      return;
    }

    // New message added - smooth scroll
    if (bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: 'smooth' });
    }
    prevMessageCount.current = messages.length;
  }, [messages]);

  if (messages.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center p-8">
        <div className="text-center text-gray-500">
          <p className="text-lg font-medium mb-2">Start a conversation</p>
          <p className="text-sm">
            Describe what you want to build and I will help you create it.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div
      ref={containerRef}
      className="flex-1 overflow-y-auto p-4"
      data-testid="message-list"
    >
      {messages.map((message) => (
        <MessageBubble key={message.id} message={message} />
      ))}

      {/* Loading indicator */}
      {isLoading && !messages.some((m) => m.isStreaming) && (
        <div className="flex justify-start mb-4">
          <div className="bg-gray-100 rounded-2xl rounded-bl-md px-4 py-3">
            <TypingIndicator />
          </div>
        </div>
      )}

      {/* Scroll anchor */}
      <div ref={bottomRef} />
    </div>
  );
}

function TypingIndicator() {
  return (
    <div className="flex items-center gap-1" aria-label="AI is typing">
      <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
      <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
      <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
    </div>
  );
}
