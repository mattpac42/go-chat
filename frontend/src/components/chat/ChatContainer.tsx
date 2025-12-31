'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { useChat } from '@/hooks/useChat';
import { useDiscovery } from '@/hooks/useDiscovery';
import { useDiscoveryExperience } from '@/hooks/useDiscoveryExperience';
import { useBuildPhase } from '@/hooks/useBuildPhase';
import { MessageList, MessageListHandle } from './MessageList';
import { ChatInput } from './ChatInput';
import { BuildPhaseProgress } from './BuildPhaseProgress';
import { MilestoneToast } from './MilestoneToast';
import { ConnectionStatus } from '@/components/shared/ConnectionStatus';
import { DiscoveryProgress, DiscoveryStageDrawer } from '@/components/discovery';
import { DiscoverySummaryModal } from '@/components/discovery/DiscoverySummaryModal';
import { CostSavingsIcon } from '@/components/savings';
import { WageSettingsModal } from '@/components/settings';
import { Message } from '@/types';

interface ChatContainerProps {
  projectId: string;
  projectTitle?: string;
  initialMessages?: Message[];
  onMenuClick?: () => void;
  onStreamingComplete?: () => void;
  onDiscoveryConfirmed?: () => void;
  onRefetchMessages?: () => Promise<void>;
  onTitleUpdate?: (newTitle: string) => Promise<void>;
}

export function ChatContainer({
  projectId,
  projectTitle = 'New Project',
  initialMessages = [],
  onMenuClick,
  onStreamingComplete,
  onDiscoveryConfirmed,
  onRefetchMessages,
  onTitleUpdate,
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

  // Build phase tracking (post-discovery)
  const discoveryComplete = currentStage === 'complete';
  const {
    currentPhase,
    messagesByPhase,
    newlyCompletedPhase,
    acknowledgePhase,
    scrollToPhase,
  } = useBuildPhase({
    messages,
    isDiscoveryMode,
    discoveryComplete,
  });

  const [showStageDrawer, setShowStageDrawer] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const [isConfirming, setIsConfirming] = useState(false);
  const [showSummaryModal, setShowSummaryModal] = useState(false);
  const [isSkipping, setIsSkipping] = useState(false);
  const [isStartingDiscovery, setIsStartingDiscovery] = useState(false);
  const [isEditingTitle, setIsEditingTitle] = useState(false);
  const [editedTitle, setEditedTitle] = useState(projectTitle);
  const [isSavingTitle, setIsSavingTitle] = useState(false);
  const [showPhasedView, setShowPhasedView] = useState(false);
  const [showWageSettings, setShowWageSettings] = useState(false);
  const messageListRef = useRef<MessageListHandle>(null);
  const titleInputRef = useRef<HTMLInputElement>(null);

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

  // Sync editedTitle when projectTitle changes externally
  useEffect(() => {
    setEditedTitle(projectTitle);
  }, [projectTitle]);

  // Focus input when editing starts
  useEffect(() => {
    if (isEditingTitle && titleInputRef.current) {
      titleInputRef.current.focus();
      titleInputRef.current.select();
    }
  }, [isEditingTitle]);

  // Title editing handlers
  const handleTitleClick = useCallback(() => {
    if (onTitleUpdate) {
      setIsEditingTitle(true);
    }
  }, [onTitleUpdate]);

  const handleTitleSave = useCallback(async () => {
    const trimmedTitle = editedTitle.trim();
    if (!trimmedTitle || trimmedTitle === projectTitle) {
      setEditedTitle(projectTitle);
      setIsEditingTitle(false);
      return;
    }

    if (onTitleUpdate) {
      setIsSavingTitle(true);
      try {
        await onTitleUpdate(trimmedTitle);
        setIsEditingTitle(false);
      } catch (error) {
        console.error('Failed to update title:', error);
        setEditedTitle(projectTitle);
      } finally {
        setIsSavingTitle(false);
      }
    }
  }, [editedTitle, projectTitle, onTitleUpdate]);

  const handleTitleCancel = useCallback(() => {
    setEditedTitle(projectTitle);
    setIsEditingTitle(false);
  }, [projectTitle]);

  const handleTitleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleTitleSave();
    } else if (e.key === 'Escape') {
      handleTitleCancel();
    }
  }, [handleTitleSave, handleTitleCancel]);

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
    <div className="flex flex-col h-full bg-white overflow-visible">
      {/* Header */}
      <header className="flex items-center justify-between px-4 py-3 pt-8 border-b border-gray-200 bg-white sticky top-0 z-10 safe-area-pt overflow-visible">
        <div className="flex items-center gap-3 min-w-0 flex-1">
          {onMenuClick && (
            <button
              onClick={onMenuClick}
              className="p-2 -ml-2 rounded-lg hover:bg-gray-100 md:hidden flex-shrink-0"
              aria-label="Open menu"
            >
              <MenuIcon className="w-6 h-6 text-gray-600" />
            </button>
          )}
          {isEditingTitle ? (
            <input
              ref={titleInputRef}
              type="text"
              value={editedTitle}
              onChange={(e) => setEditedTitle(e.target.value)}
              onBlur={handleTitleSave}
              onKeyDown={handleTitleKeyDown}
              disabled={isSavingTitle}
              className="text-lg font-semibold text-gray-900 bg-transparent border-b-2 border-teal-400 outline-none px-1 py-0 min-w-0 flex-1 disabled:opacity-50"
              aria-label="Edit project title"
            />
          ) : (
            <h1
              onClick={handleTitleClick}
              className={`text-lg font-semibold text-gray-900 truncate ${
                onTitleUpdate ? 'cursor-pointer hover:text-teal-600 transition-colors' : ''
              }`}
              title={onTitleUpdate ? 'Click to edit title' : undefined}
            >
              {projectTitle}
            </h1>
          )}
        </div>
        <div className="flex items-center gap-2 sm:gap-4">
          {/* Cost Savings Icon - shown when there are messages */}
          {messages.length > 0 && (
            <CostSavingsIcon
              metrics={{
                messageCount: messages.length,
                filesGenerated: 0,
                pmMessageCount: messages.filter(m => m.role === 'assistant' && (m.agentType === 'product_manager' || m.agentType === 'product')).length,
                designerMessageCount: messages.filter(m => m.role === 'assistant' && m.agentType === 'designer').length,
                developerMessageCount: messages.filter(m => m.role === 'assistant' && m.agentType === 'developer').length,
              }}
              storageKey={`cost-savings-${projectId}`}
            />
          )}
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
          {/* Settings - shown when there are messages */}
          {messages.length > 0 && (
            <button
              onClick={() => setShowWageSettings(true)}
              className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
              aria-label="Wage Settings"
              title="Configure hourly rates"
            >
              <SettingsIcon className="w-5 h-5" />
            </button>
          )}
        </div>
      </header>

      {/* Build Phase Progress - shown after discovery is complete */}
      {discoveryComplete && messages.length > 0 && (
        <BuildPhaseProgress
          messages={messages}
          isDiscoveryMode={isDiscoveryMode}
          discoveryComplete={discoveryComplete}
          currentPhase={currentPhase}
          onPhaseClick={scrollToPhase}
          showPhasedView={showPhasedView}
          onTogglePhasedView={() => setShowPhasedView((prev) => !prev)}
        />
      )}

      {/* Milestone Toast - shows when a phase completes */}
      {newlyCompletedPhase && (
        <MilestoneToast
          phase={newlyCompletedPhase}
          onDismiss={() => acknowledgePhase(newlyCompletedPhase)}
        />
      )}

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
        currentStage={currentStage}
        showSkipDiscovery={experienceLoaded && hasCompletedBefore && !isSkipping}
        onSkipDiscovery={handleSkipDiscovery}
        onStartDiscovery={handleStartDiscovery}
        isStartingDiscovery={isStartingDiscovery}
        showPhasedView={showPhasedView && discoveryComplete}
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

      {/* Wage Settings Modal */}
      <WageSettingsModal
        isOpen={showWageSettings}
        onClose={() => setShowWageSettings(false)}
      />
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

function SettingsIcon({ className }: { className?: string }) {
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
        d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
      />
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
      />
    </svg>
  );
}
