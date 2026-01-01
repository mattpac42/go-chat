'use client';

import { useState, useCallback, useRef, useEffect } from 'react';
import { Message, ChatState, ServerMessage, ConnectionStatus, CompletenessReport } from '@/types';
import { useWebSocket } from './useWebSocket';

interface UseChatOptions {
  projectId: string;
  initialMessages?: Message[];
  onFilesUpdated?: () => void;
}

interface UseChatReturn {
  messages: Message[];
  isLoading: boolean;
  error: string | null;
  connectionStatus: ConnectionStatus;
  reconnectAttempts: number;
  completenessReport: CompletenessReport | null;
  sendMessage: (content: string) => void;
  clearError: () => void;
  reconnect: () => void;
}

// Delay before showing connection error (to avoid flashing during project switch)
const CONNECTION_ERROR_DELAY = 800;

// Regex to strip complete discovery metadata tags
const DISCOVERY_DATA_PATTERN = /<!--DISCOVERY_DATA:[\s\S]*?-->/g;

/**
 * Get display-safe content by stripping discovery metadata
 * Handles both complete tags and incomplete tags still being streamed
 */
function getDisplayContent(content: string): string {
  // First, strip any complete discovery metadata tags
  let result = content.replace(DISCOVERY_DATA_PATTERN, '');

  // Then check for incomplete HTML comment that might be discovery data
  // This handles the case where we've received "<!--DISCOVERY_DATA:..." but not the closing "-->"
  const incompleteTagStart = result.lastIndexOf('<!--');
  if (incompleteTagStart !== -1) {
    // Check if there's a closing "-->" after this opening
    const closingTag = result.indexOf('-->', incompleteTagStart);
    if (closingTag === -1) {
      // No closing tag yet - this is an incomplete comment, hide it
      result = result.substring(0, incompleteTagStart);
    }
  }

  return result;
}

/**
 * Chat hook that integrates with WebSocket for real-time messaging
 *
 * Features:
 * - Manages chat message state
 * - Handles streaming message assembly (chunks -> full message)
 * - Manages loading states during AI response
 * - Provides error handling and reconnection
 */
export function useChat({ projectId, initialMessages = [], onFilesUpdated }: UseChatOptions): UseChatReturn {
  const [state, setState] = useState<ChatState>({
    messages: initialMessages,
    isLoading: false,
    error: null,
  });
  const [completenessReport, setCompletenessReport] = useState<CompletenessReport | null>(null);

  // Sync initialMessages when they change (e.g., welcome message loaded after discovery)
  // Only update if we have no messages and initialMessages has content
  useEffect(() => {
    if (state.messages.length === 0 && initialMessages.length > 0) {
      setState(prev => ({ ...prev, messages: initialMessages }));
    }
  }, [initialMessages, state.messages.length]);

  // Track streaming message by ID
  const streamingMessageRef = useRef<Map<string, Message>>(new Map());

  // Track delayed error timeout
  const errorTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  /**
   * Handle incoming WebSocket messages
   */
  const handleWebSocketMessage = useCallback((serverMessage: ServerMessage) => {
    switch (serverMessage.type) {
      case 'message_start': {
        // Create new streaming message
        const newMessage: Message = {
          id: serverMessage.messageId,
          projectId: serverMessage.projectId || projectId,
          role: 'assistant',
          content: '',
          timestamp: new Date().toISOString(),
          isStreaming: true,
        };
        streamingMessageRef.current.set(serverMessage.messageId, newMessage);
        setState(prev => ({
          ...prev,
          messages: [...prev.messages, newMessage],
          isLoading: true,
          error: null,
        }));
        break;
      }

      case 'message_chunk': {
        const streamingMessage = streamingMessageRef.current.get(serverMessage.messageId);
        if (streamingMessage && serverMessage.content) {
          // Accumulate content in the streaming message
          // Store raw content in ref for complete message assembly
          const rawContent = streamingMessage.content + serverMessage.content;
          const updatedMessage = {
            ...streamingMessage,
            content: rawContent,
          };
          streamingMessageRef.current.set(serverMessage.messageId, updatedMessage);

          // Get display-safe content (strips complete and incomplete metadata)
          const displayContent = getDisplayContent(rawContent);

          setState(prev => ({
            ...prev,
            messages: prev.messages.map(msg =>
              msg.id === serverMessage.messageId
                ? { ...msg, content: displayContent }
                : msg
            ),
          }));
        }
        break;
      }

      case 'message_complete': {
        // Finalize the streaming message
        streamingMessageRef.current.delete(serverMessage.messageId);
        // Strip any metadata from final content
        const finalContent = getDisplayContent(serverMessage.fullContent || '');
        setState(prev => ({
          ...prev,
          messages: prev.messages.map(msg =>
            msg.id === serverMessage.messageId
              ? {
                  ...msg,
                  content: finalContent || msg.content,
                  isStreaming: false,
                  agentType: serverMessage.agentType || msg.agentType,
                }
              : msg
          ),
          isLoading: false,
        }));
        // Update completeness report if provided
        if (serverMessage.completenessReport) {
          setCompletenessReport(serverMessage.completenessReport);
        }
        break;
      }

      case 'error': {
        // Handle error and clean up any streaming message
        const errorMessageId = serverMessage.messageId;
        if (errorMessageId) {
          streamingMessageRef.current.delete(errorMessageId);
          // Remove failed streaming message from UI
          setState(prev => ({
            ...prev,
            messages: prev.messages.filter(msg => msg.id !== errorMessageId),
            isLoading: false,
            error: serverMessage.error || 'An error occurred',
          }));
        } else {
          setState(prev => ({
            ...prev,
            isLoading: false,
            error: serverMessage.error || 'An error occurred',
          }));
        }
        break;
      }

      case 'files_updated': {
        // Trigger file refresh callback when files are created/updated via tool use
        onFilesUpdated?.();
        break;
      }
    }
  }, [projectId, onFilesUpdated]);

  /**
   * Clear any pending error timeout
   */
  const clearErrorTimeout = useCallback(() => {
    if (errorTimeoutRef.current) {
      clearTimeout(errorTimeoutRef.current);
      errorTimeoutRef.current = null;
    }
  }, []);

  /**
   * Handle WebSocket connection events
   */
  const handleConnect = useCallback(() => {
    // Clear any pending error timeout - connection succeeded
    clearErrorTimeout();
    setState(prev => ({
      ...prev,
      error: null,
    }));
  }, [clearErrorTimeout]);

  const handleDisconnect = useCallback(() => {
    // Clean up any in-progress streaming
    streamingMessageRef.current.clear();
    setState(prev => ({
      ...prev,
      isLoading: false,
    }));
  }, []);

  const handleError = useCallback(() => {
    // Delay showing error to avoid flash during project switch
    clearErrorTimeout();
    errorTimeoutRef.current = setTimeout(() => {
      setState(prev => ({
        ...prev,
        error: 'Connection error. Attempting to reconnect...',
      }));
    }, CONNECTION_ERROR_DELAY);
  }, [clearErrorTimeout]);

  // Cleanup error timeout on unmount
  useEffect(() => {
    return () => {
      clearErrorTimeout();
    };
  }, [clearErrorTimeout]);

  // Initialize WebSocket connection
  const {
    status: connectionStatus,
    sendMessage: wsSendMessage,
    connect: wsConnect,
    reconnectAttempts,
  } = useWebSocket({
    projectId,
    onMessage: handleWebSocketMessage,
    onConnect: handleConnect,
    onDisconnect: handleDisconnect,
    onError: handleError,
    autoConnect: true,
  });

  /**
   * Send a chat message
   */
  const sendMessage = useCallback((content: string) => {
    const trimmedContent = content.trim();
    if (!trimmedContent) return;

    // Add user message to state immediately
    const userMessage: Message = {
      id: `user-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      projectId,
      role: 'user',
      content: trimmedContent,
      timestamp: new Date().toISOString(),
    };

    setState(prev => ({
      ...prev,
      messages: [...prev.messages, userMessage],
      error: null,
      isLoading: true,
    }));

    // Send via WebSocket
    wsSendMessage({
      type: 'chat_message',
      projectId,
      content: trimmedContent,
      timestamp: new Date().toISOString(),
    });
  }, [projectId, wsSendMessage]);

  /**
   * Clear error state
   */
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }));
  }, []);

  /**
   * Manually trigger reconnection
   */
  const reconnect = useCallback(() => {
    wsConnect();
  }, [wsConnect]);

  return {
    messages: state.messages,
    isLoading: state.isLoading,
    error: state.error,
    connectionStatus,
    reconnectAttempts,
    completenessReport,
    sendMessage,
    clearError,
    reconnect,
  };
}
