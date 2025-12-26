'use client';

import { useEffect, useRef, useState } from 'react';
import { useChat } from '@/hooks/useChat';
import { useDiscovery } from '@/hooks/useDiscovery';
import { MessageList } from './MessageList';
import { ChatInput } from './ChatInput';
import { ConnectionStatus } from '@/components/shared/ConnectionStatus';
import { DiscoveryProgress, DiscoveryStageDrawer, DiscoverySummaryCard } from '@/components/discovery';
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
      <MessageList messages={messages} isLoading={isLoading} />

      {/* Discovery Summary Card - shown when discovery reaches summary stage */}
      {showSummaryCard && (
        <div className="px-4 py-3 border-t border-gray-100 bg-gray-50">
          <DiscoverySummaryCard
            summary={summary}
            onEdit={handleEditDiscovery}
            onConfirm={handleConfirmDiscovery}
            isConfirming={isConfirming}
          />
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
