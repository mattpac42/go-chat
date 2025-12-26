'use client';

import { useEffect, useRef, useState } from 'react';
import { useChat } from '@/hooks/useChat';
import { useDiscovery } from '@/hooks/useDiscovery';
import { MessageList } from './MessageList';
import { ChatInput } from './ChatInput';
import { ConnectionStatus } from '@/components/shared/ConnectionStatus';
import { DiscoveryProgress, DiscoveryStageDrawer, DiscoverySummaryCard } from '@/components/discovery';
import { DiscoverySummaryModal } from '@/components/discovery/DiscoverySummaryModal';
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

  // Discovery integration
  const {
    isDiscoveryMode,
    currentStage,
    summary,
    confirmDiscovery,
    resetDiscovery,
    refetch: refetchDiscovery,
  } = useDiscovery(projectId);
  const [showStageDrawer, setShowStageDrawer] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const [isConfirming, setIsConfirming] = useState(false);
  const [isSummaryCollapsed, setIsSummaryCollapsed] = useState(false);
  const [showSummaryModal, setShowSummaryModal] = useState(false);

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
    } catch {
      // Error is handled in the hook
    } finally {
      setIsConfirming(false);
    }
  };

  const handleEditDiscovery = () => {
    resetDiscovery();
  };

  // Check if we should show the summary card
  const showSummaryCard = currentStage === 'summary' && summary !== null;

  // Check if discovery is complete and we have a summary to show
  const showViewSummaryButton = currentStage === 'complete' && summary !== null;

  // Stage-aware input placeholder
  const getPlaceholder = () => {
    if (connectionStatus !== 'connected') return 'Connecting...';
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
      <MessageList messages={messages} projectId={projectId} isLoading={isLoading} hasBottomCard={showSummaryCard} />

      {/* Discovery Summary Card - shown when discovery reaches summary stage */}
      {showSummaryCard && (
        <div className="border-t border-gray-100 bg-gray-50">
          {isSummaryCollapsed ? (
            <button
              onClick={() => setIsSummaryCollapsed(false)}
              className="w-full px-4 py-3 flex items-center justify-between text-sm text-gray-600 hover:bg-gray-100 transition-colors"
            >
              <span className="font-medium">Show Project Summary</span>
              <ChevronUpIcon className="w-4 h-4" />
            </button>
          ) : (
            <div className="px-4 py-3">
              <div className="flex justify-end mb-2">
                <button
                  onClick={() => setIsSummaryCollapsed(true)}
                  className="text-xs text-gray-500 hover:text-gray-700 flex items-center gap-1"
                >
                  <span>Collapse</span>
                  <ChevronDownIcon className="w-3 h-3" />
                </button>
              </div>
              <DiscoverySummaryCard
                summary={summary}
                onEdit={handleEditDiscovery}
                onConfirm={handleConfirmDiscovery}
                isConfirming={isConfirming}
              />
            </div>
          )}
        </div>
      )}

      {/* Input */}
      <ChatInput
        onSend={sendMessage}
        disabled={connectionStatus !== 'connected' || isLoading}
        placeholder={getPlaceholder()}
      />

      {/* Discovery stage drawer for mobile */}
      <DiscoveryStageDrawer
        isOpen={showStageDrawer}
        onClose={() => setShowStageDrawer(false)}
        currentStage={currentStage}
      />

      {/* Summary modal for completed discovery */}
      {summary && (
        <DiscoverySummaryModal
          isOpen={showSummaryModal}
          onClose={() => setShowSummaryModal(false)}
          summary={summary}
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

function ChevronUpIcon({ className }: { className?: string }) {
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
        d="M5 15l7-7 7 7"
      />
    </svg>
  );
}

function ChevronDownIcon({ className }: { className?: string }) {
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
        d="M19 9l-7 7-7-7"
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
