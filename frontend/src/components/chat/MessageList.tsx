'use client';

import { useEffect, useRef, useLayoutEffect, useMemo, forwardRef, useImperativeHandle, useState } from 'react';
import { Message, AgentType } from '@/types';
import { DiscoveryStage } from '@/types/discovery';
import { MessageBubble } from './MessageBubble';
import { PersonaIntroduction, isIntroductionMessage, isTeamIntroMessage } from './PersonaIntroduction';
import { PhaseSection } from './PhaseSection';
import { BuildPhase, groupMessagesByPhase, detectCurrentPhase } from './BuildPhaseProgress';
import { useAgentIntroductions } from '@/hooks/useAgentIntroductions';
import { usePersonaIntroductions } from '@/hooks/usePersonaIntroductions';

interface MessageListProps {
  messages: Message[];
  projectId: string;
  isLoading?: boolean;
  isDiscoveryMode?: boolean;
  currentStage?: DiscoveryStage;
  showSkipDiscovery?: boolean;
  onSkipDiscovery?: () => void;
  onStartDiscovery?: () => void;
  isStartingDiscovery?: boolean;
  /** Show messages grouped by build phase with collapsible sections */
  showPhasedView?: boolean;
}

export interface MessageListHandle {
  scrollToBottom: () => void;
}

export const MessageList = forwardRef<MessageListHandle, MessageListProps>(function MessageList({
  messages,
  projectId,
  isLoading = false,
  isDiscoveryMode = false,
  currentStage = 'welcome',
  showSkipDiscovery = false,
  onSkipDiscovery,
  onStartDiscovery,
  isStartingDiscovery = false,
  showPhasedView = false,
}, ref) {
  const containerRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);
  const isInitialLoad = useRef(true);
  const prevMessageCount = useRef(0);

  // Expose scrollToBottom method to parent
  useImperativeHandle(ref, () => ({
    scrollToBottom: () => {
      if (bottomRef.current) {
        bottomRef.current.scrollIntoView({ behavior: 'smooth' });
      }
    },
  }));

  // Track which agents have been introduced to the user
  const { hasMetAgent, markAgentMet } = useAgentIntroductions(projectId);

  // Handle persona introductions when transitioning from discovery to building
  const { processMessagesWithIntroductions } = usePersonaIntroductions(
    projectId,
    currentStage,
    messages
  );

  // Process messages to inject persona introductions if needed
  const processedMessages = useMemo(() => {
    return processMessagesWithIntroductions(messages);
  }, [messages, processMessagesWithIntroductions]);

  // Process messages to determine which should show the "NEW" badge
  // An agent gets a badge only on their FIRST appearance in this session
  // AND only if the user hasn't met them before (across sessions)
  // Skip badge logic for introduction messages (they have their own styling)
  const messagesWithBadges = useMemo(() => {
    const seenInThisRender = new Set<AgentType>();

    return processedMessages.map((message) => {
      // Introduction messages have their own styling, don't show badges
      if (isIntroductionMessage(message)) {
        return { message, showBadge: false, isIntro: true, isTeamIntro: isTeamIntroMessage(message) };
      }

      // Only assistant messages with agentType can have badges
      if (message.role !== 'assistant' || !message.agentType) {
        return { message, showBadge: false, isIntro: false, isTeamIntro: false };
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

      return { message, showBadge, isIntro: false, isTeamIntro: false };
    });
  }, [processedMessages, hasMetAgent, markAgentMet]);

  // Group messages by build phase for phased view
  const messagesByPhase = useMemo(() => {
    return groupMessagesByPhase(messages);
  }, [messages]);

  // Get current phase for phased view
  const currentPhase = useMemo(() => {
    return detectCurrentPhase(messages);
  }, [messages]);

  // Create a lookup for badges by message id
  const badgeLookup = useMemo(() => {
    const lookup = new Map<string, boolean>();
    messagesWithBadges.forEach(({ message, showBadge }) => {
      lookup.set(message.id, showBadge);
    });
    return lookup;
  }, [messagesWithBadges]);

  // Instant scroll to bottom on initial load (no animation)
  useLayoutEffect(() => {
    if (messages.length > 0 && isInitialLoad.current && containerRef.current) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
      isInitialLoad.current = false;
      prevMessageCount.current = messages.length;
    }
  }, [messages]);

  // Smooth scroll for new messages and streaming updates
  useEffect(() => {
    // Skip if this is initial load
    if (isInitialLoad.current) {
      prevMessageCount.current = messages.length;
      return;
    }

    // Check if last message is streaming
    const lastMessage = messages[messages.length - 1];
    const isStreaming = lastMessage?.isStreaming;

    // Scroll on new message OR while streaming
    if (messages.length > prevMessageCount.current || isStreaming) {
      if (bottomRef.current) {
        bottomRef.current.scrollIntoView({ behavior: 'smooth' });
      }
    }
    prevMessageCount.current = messages.length;
  }, [messages]);

  if (messages.length === 0) {
    if (isDiscoveryMode) {
      return (
        <div className="flex-1 flex flex-col items-center justify-center p-4">
          {/* Intro Card */}
          <div className="bg-gradient-to-br from-teal-50 to-cyan-50 border border-teal-200 rounded-xl p-6 max-w-md w-full shadow-sm">
            {/* Header */}
            <div className="flex items-start gap-4">
              <div className="flex-shrink-0">
                <div className="w-12 h-12 rounded-full bg-teal-500 flex items-center justify-center">
                  <CompassIcon className="w-6 h-6 text-white" />
                </div>
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-teal-800 mb-1">
                  Let&apos;s figure out what you need
                </h3>
                <div className="flex items-center gap-2 text-sm text-teal-600">
                  <ClockIcon className="w-4 h-4" />
                  <span>About 5 minutes</span>
                </div>
              </div>
            </div>

            {/* Content */}
            <div className="mt-4 ml-16">
              <p className="text-sm text-gray-600 mb-3">
                Root will help you clarify:
              </p>
              <ul className="space-y-2">
                <li className="flex items-start gap-2 text-sm text-gray-700">
                  <CheckCircleIcon className="w-4 h-4 text-teal-500 mt-0.5 flex-shrink-0" />
                  <span>The problem you&apos;re solving</span>
                </li>
                <li className="flex items-start gap-2 text-sm text-gray-700">
                  <CheckCircleIcon className="w-4 h-4 text-teal-500 mt-0.5 flex-shrink-0" />
                  <span>Who will use your product</span>
                </li>
                <li className="flex items-start gap-2 text-sm text-gray-700">
                  <CheckCircleIcon className="w-4 h-4 text-teal-500 mt-0.5 flex-shrink-0" />
                  <span>Essential features to start with</span>
                </li>
              </ul>
            </div>

            {/* Primary CTA */}
            <div className="mt-6">
              <button
                onClick={onStartDiscovery}
                disabled={isStartingDiscovery}
                className="w-full px-4 py-3 bg-teal-500 text-white font-medium rounded-lg hover:bg-teal-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
              >
                {isStartingDiscovery ? (
                  <>
                    <WaitingIndicator />
                    <span>Starting...</span>
                  </>
                ) : (
                  <>
                    <SeedlingIcon className="w-5 h-5" />
                    <span>Let&apos;s solve my problem</span>
                  </>
                )}
              </button>
            </div>

            {/* Skip option for returning users */}
            {showSkipDiscovery && onSkipDiscovery && (
              <div className="mt-4 pt-4 border-t border-teal-200 flex justify-center">
                <button
                  onClick={onSkipDiscovery}
                  disabled={isStartingDiscovery}
                  className="text-sm text-teal-600 hover:text-teal-800 font-medium flex items-center gap-1.5 transition-colors disabled:opacity-50"
                >
                  <span>Done this before? Skip to building</span>
                  <ArrowRightIcon className="w-4 h-4" />
                </button>
              </div>
            )}
          </div>
        </div>
      );
    }

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

  // Phase labels for phased view
  const phaseLabels: Record<BuildPhase, string> = {
    discovery: 'Discovery',
    planning: 'Planning',
    building: 'Building',
    testing: 'Testing',
    launch: 'Launch',
  };

  return (
    <div
      ref={containerRef}
      className="flex-1 overflow-y-auto p-4"
      data-testid="message-list"
    >
      {showPhasedView ? (
        // Phased view: messages grouped by build phase with collapsible sections
        <>
          {(['discovery', 'planning', 'building', 'testing', 'launch'] as BuildPhase[]).map((phase) => {
            const phaseMessages = messagesByPhase.get(phase) || [];
            if (phaseMessages.length === 0) return null;

            return (
              <PhaseSection
                key={phase}
                phase={phase}
                label={phaseLabels[phase]}
                messages={phaseMessages}
                isCurrentPhase={phase === currentPhase}
                defaultExpanded={phase === currentPhase}
                showBadgeForMessage={(messageId) => badgeLookup.get(messageId) ?? false}
              />
            );
          })}
        </>
      ) : (
        // Standard view: messages in chronological order
        messagesWithBadges.map(({ message, showBadge, isIntro, isTeamIntro }) => (
          isIntro ? (
            <PersonaIntroduction
              key={message.id}
              message={message}
              isTeamIntro={isTeamIntro}
            />
          ) : (
            <MessageBubble
              key={message.id}
              message={message}
              showBadge={showBadge}
            />
          )
        ))
      )}

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
});

function TypingIndicator() {
  return (
    <div className="flex items-center gap-1" aria-label="AI is typing">
      <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
      <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
      <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
    </div>
  );
}

function WaitingIndicator() {
  return (
    <div className="flex items-center gap-1.5" aria-label="Waiting for Guide">
      <span className="w-2 h-2 bg-teal-400 rounded-full animate-pulse" style={{ animationDelay: '0ms' }} />
      <span className="w-2 h-2 bg-teal-400 rounded-full animate-pulse" style={{ animationDelay: '200ms' }} />
      <span className="w-2 h-2 bg-teal-400 rounded-full animate-pulse" style={{ animationDelay: '400ms' }} />
    </div>
  );
}

function CompassIcon({ className }: { className?: string }) {
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
        d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
      />
    </svg>
  );
}

function ClockIcon({ className }: { className?: string }) {
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
        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  );
}

function CheckCircleIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="currentColor"
      viewBox="0 0 20 20"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fillRule="evenodd"
        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
        clipRule="evenodd"
      />
    </svg>
  );
}

function ArrowRightIcon({ className }: { className?: string }) {
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
        d="M13 7l5 5m0 0l-5 5m5-5H6"
      />
    </svg>
  );
}

function SeedlingIcon({ className }: { className?: string }) {
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
        d="M12 19V12M12 12C12 9 14 7 17 7C17 10 15 12 12 12ZM12 12C12 9 10 7 7 7C7 10 9 12 12 12Z"
      />
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M12 22C12 22 12 19 12 19"
      />
    </svg>
  );
}
