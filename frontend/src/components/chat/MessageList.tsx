'use client';

import { useEffect, useRef, useLayoutEffect, useMemo } from 'react';
import { Message, AgentType } from '@/types';
import { MessageBubble } from './MessageBubble';
import { useAgentIntroductions } from '@/hooks/useAgentIntroductions';

interface MessageListProps {
  messages: Message[];
  projectId: string;
  isLoading?: boolean;
  hasBottomCard?: boolean;
}

export function MessageList({
  messages,
  projectId,
  isLoading = false,
  hasBottomCard = false,
}: MessageListProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);
  const isInitialLoad = useRef(true);
  const prevMessageCount = useRef(0);

  // Track which agents have been introduced to the user
  const { hasMetAgent, markAgentMet } = useAgentIntroductions(projectId);

  // Process messages to determine which should show the "NEW" badge
  // An agent gets a badge only on their FIRST appearance in this session
  // AND only if the user hasn't met them before (across sessions)
  const messagesWithBadges = useMemo(() => {
    const seenInThisRender = new Set<AgentType>();

    return messages.map((message) => {
      // Only assistant messages with agentType can have badges
      if (message.role !== 'assistant' || !message.agentType) {
        return { message, showBadge: false };
      }

      const agentType = message.agentType;

      // First time seeing this agent in this render AND user hasn't met them
      const isFirstInRender = !seenInThisRender.has(agentType);
      const isNewToUser = !hasMetAgent(agentType);
      const showBadge = isFirstInRender && isNewToUser;

      // Mark as seen in this render
      seenInThisRender.add(agentType);

      // Mark the agent as met (persists to localStorage)
      if (showBadge) {
        markAgentMet(agentType);
      }

      return { message, showBadge };
    });
  }, [messages, hasMetAgent, markAgentMet]);

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
      className={`flex-1 overflow-y-auto p-4 ${hasBottomCard ? 'pb-[500px]' : ''}`}
      data-testid="message-list"
    >
      {messagesWithBadges.map(({ message, showBadge }) => (
        <MessageBubble
          key={message.id}
          message={message}
          showBadge={showBadge}
        />
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
