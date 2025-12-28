'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { useChat } from '@/hooks/useChat';
import { useDiscovery } from '@/hooks/useDiscovery';
import { useDiscoveryExperience } from '@/hooks/useDiscoveryExperience';
import { MessageList, MessageListHandle } from './MessageList';
import { ChatInput } from './ChatInput';
import { ConnectionStatus } from '@/components/shared/ConnectionStatus';
import { DiscoveryProgress, DiscoveryStageDrawer } from '@/components/discovery';
import { DiscoverySummaryModal } from '@/components/discovery/DiscoverySummaryModal';
import { Message } from '@/types';

interface ChatContainerProps {
  projectId: string;
  projectTitle?: string;
  initialMessages?: Message[];
  onMenuClick?: () => void;
  onStreamingComplete?: () => void;
  onDiscoveryConfirmed?: () => void;
  onRefetchMessages?: () => Promise<void>;
}

export function ChatContainer({
  projectId,
  projectTitle = 'New Project',
  initialMessages = [],
  onMenuClick,
  onStreamingComplete,
  onDiscoveryConfirmed,
  onRefetchMessages,
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

  // Discovery integration
  const {
    isDiscoveryMode,
    currentStage,
    summary,
    confirmDiscovery,
    resetDiscovery,
    skipDiscovery,
    refetch: refetchDiscovery,
  } = useDiscovery(projectId);

  // Track if user has completed discovery before (for skip option)
  const { hasCompletedBefore, isLoaded: experienceLoaded, markDiscoveryCompleted } = useDiscoveryExperience();

  const [showStageDrawer, setShowStageDrawer] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const [isConfirming, setIsConfirming] = useState(false);
  const [showSummaryModal, setShowSummaryModal] = useState(false);
  const [isSkipping, setIsSkipping] = useState(false);
  const [isStartingDiscovery, setIsStartingDiscovery] = useState(false);
  const messageListRef = useRef<MessageListHandle>(null);

  // Track previous loading state to detect when streaming completes
  const wasLoadingRef = useRef(false);

  useEffect(() => {
    // Detect transition from loading to not loading (streaming complete)
    if (wasLoadingRef.current && !isLoading) {
      onStreamingComplete?.();
      // Refetch discovery state after each response to update progress
      refetchDiscovery();
    }
    wasLoadingRef.current = isLoading;
  }, [isLoading, onStreamingComplete, refetchDiscovery]);

  // Mobile detection
  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 768);
    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  // Discovery summary handlers
  const handleConfirmDiscovery = async () => {
    setIsConfirming(true);
    try {
      await confirmDiscovery();
      // Mark that user has completed discovery (for skip option on future projects)
      markDiscoveryCompleted();
      // Notify parent to refetch project data (title may have changed)
      onDiscoveryConfirmed?.();
      // Close the modal after successful confirmation
      setShowSummaryModal(false);
    } catch {
      // Error is handled in the hook
    } finally {
      setIsConfirming(false);
    }
  };

  const handleEditDiscovery = () => {
    // Close modal and return to chat for editing
    // Don't reset discovery - just let user continue the conversation
    setShowSummaryModal(false);
  };

  const handleStartOver = () => {
    // Close modal and reset discovery
    setShowSummaryModal(false);
    resetDiscovery();
  };

  // Handler for skipping discovery (for returning users)
  const handleSkipDiscovery = useCallback(async () => {
    setIsSkipping(true);
    try {
      const success = await skipDiscovery();
      if (success) {
        markDiscoveryCompleted();
      }
    } finally {
      setIsSkipping(false);
    }
  }, [skipDiscovery, markDiscoveryCompleted]);

  // Handler for starting discovery (user clicks CTA)
  const handleStartDiscovery = useCallback(async () => {
    if (!onRefetchMessages) return;

    setIsStartingDiscovery(true);
    try {
      // Give the backend a moment to generate the welcome message, then refetch
      await new Promise(resolve => setTimeout(resolve, 500));
      await onRefetchMessages();
      // If still no messages after first attempt, try again
      if (messages.length === 0) {
        await new Promise(resolve => setTimeout(resolve, 2000));
        await onRefetchMessages();
      }
    } finally {
      setIsStartingDiscovery(false);
    }
  }, [onRefetchMessages, messages.length]);

  // Check if we should show the summary card
  // Only show when we have meaningful summary data (at least MVP features)
  const hasMeaningfulSummary = summary !== null &&
    summary.mvpFeatures && summary.mvpFeatures.length > 0;
  const showSummaryCard = currentStage === 'summary' && hasMeaningfulSummary;

  // Check if discovery is complete and we have a summary to show
  const showViewSummaryButton = currentStage === 'complete' && summary !== null;


  // Check if we are waiting for the first message in discovery mode
  const isWaitingForDiscoveryStart = isDiscoveryMode && messages.length === 0;

  // Stage-aware input placeholder
  const getPlaceholder = () => {
    if (connectionStatus !== 'connected') return 'Connecting...';
    if (isStartingDiscovery) return 'Root is joining...';
    if (isWaitingForDiscoveryStart) return 'Click the button above to start...';
    if (isLoading) return 'Waiting for response...';
    if (!isDiscoveryMode) return 'Describe what you want to build...';

    switch (currentStage) {
      case 'welcome':
        return 'Tell me about yourself...';
      case 'problem':
        return 'What challenges do you face?';
      case 'personas':
        return 'Who will use this?';
      case 'mvp':
        return 'What features are essential?';
      case 'summary':
        return 'Ready to start building?';
      default:
        return 'Describe what you want to build...';
    }
  };

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
        <div className="flex items-center gap-4">
          {showViewSummaryButton && (
            <button
              onClick={() => setShowSummaryModal(true)}
              className="px-3 py-1.5 text-sm font-medium text-teal-600 bg-teal-50 rounded-lg hover:bg-teal-100 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:ring-offset-1 transition-colors flex items-center gap-1.5"
            >
              <DocumentIcon className="w-4 h-4" />
              <span className="hidden sm:inline">Project Summary</span>
            </button>
          )}
          {isDiscoveryMode && (
            <DiscoveryProgress
              currentStage={currentStage}
              onStageClick={() => setShowStageDrawer(true)}
              isMobile={isMobile}
            />
          )}
          <ConnectionStatus
            status={connectionStatus}
            reconnectAttempts={reconnectAttempts}
            onReconnect={reconnect}
          />
        </div>
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
      <MessageList
        ref={messageListRef}
        messages={messages}
        projectId={projectId}
        isLoading={isLoading}
        isDiscoveryMode={isDiscoveryMode}
        showSkipDiscovery={experienceLoaded && hasCompletedBefore && !isSkipping}
        onSkipDiscovery={handleSkipDiscovery}
        onStartDiscovery={handleStartDiscovery}
        isStartingDiscovery={isStartingDiscovery}
      />

      {/* Discovery Complete Notification Bar - slim bar that opens modal */}
      {showSummaryCard && (
        <button
          onClick={() => setShowSummaryModal(true)}
          className="w-full px-4 py-3 bg-gradient-to-r from-teal-500 to-teal-600 text-white flex items-center justify-center gap-2 hover:from-teal-600 hover:to-teal-700 transition-all group"
        >
          <CheckCircleIcon className="w-5 h-5" />
          <span className="font-medium">Discovery complete!</span>
          <span className="text-teal-100">Review your project summary</span>
          <ChevronRightIcon className="w-4 h-4 group-hover:translate-x-0.5 transition-transform" />
        </button>
      )}

      {/* Input */}
      <ChatInput
        projectId={projectId}
        onSend={sendMessage}
        disabled={connectionStatus !== 'connected' || isLoading || isWaitingForDiscoveryStart || isStartingDiscovery}
        placeholder={getPlaceholder()}
      />

      {/* Discovery stage drawer for mobile */}
      <DiscoveryStageDrawer
        isOpen={showStageDrawer}
        onClose={() => setShowStageDrawer(false)}
        currentStage={currentStage}
      />

      {/* Summary modal - action mode when in summary stage, view-only when complete */}
      {summary && (
        <DiscoverySummaryModal
          isOpen={showSummaryModal}
          onClose={() => setShowSummaryModal(false)}
          summary={summary}
          messageCount={messages.length}
          // Pass action handlers only when in summary stage (pre-confirmation)
          {...(showSummaryCard ? {
            onConfirm: handleConfirmDiscovery,
            onEdit: handleEditDiscovery,
            onStartOver: handleStartOver,
            isConfirming: isConfirming,
          } : {})}
        />
      )}
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

function CheckCircleIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  );
}

function ChevronRightIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9 5l7 7-7 7"
      />
    </svg>
  );
}

function DocumentIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
      />
    </svg>
  );
}
