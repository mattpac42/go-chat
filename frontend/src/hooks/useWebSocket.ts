'use client';

import { useState, useCallback, useRef, useEffect } from 'react';
import { ConnectionStatus, ServerMessage, ClientMessage } from '@/types';
import { getWebSocketUrl } from '@/lib/api';

interface UseWebSocketOptions {
  projectId: string;
  onMessage?: (message: ServerMessage) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
  onError?: (error: Event) => void;
  autoConnect?: boolean;
  maxReconnectAttempts?: number;
}

interface UseWebSocketReturn {
  status: ConnectionStatus;
  sendMessage: (message: ClientMessage) => void;
  connect: () => void;
  disconnect: () => void;
  reconnectAttempts: number;
}

// Exponential backoff configuration
const BASE_RECONNECT_DELAY = 1000; // 1 second
const MAX_RECONNECT_DELAY = 30000; // 30 seconds
const DEFAULT_MAX_RECONNECT_ATTEMPTS = 5;

/**
 * Calculate exponential backoff delay
 */
function getReconnectDelay(attempt: number): number {
  const delay = Math.min(
    BASE_RECONNECT_DELAY * Math.pow(2, attempt),
    MAX_RECONNECT_DELAY
  );
  // Add jitter (0-25% of delay) to prevent thundering herd
  const jitter = delay * Math.random() * 0.25;
  return delay + jitter;
}

/**
 * WebSocket hook for real-time chat communication
 *
 * Features:
 * - Connects to real backend WebSocket
 * - Handles connection states (connecting, connected, disconnected)
 * - Implements reconnection with exponential backoff
 * - Handles message types: message_start, message_chunk, message_complete, error
 */
export function useWebSocket(options: UseWebSocketOptions): UseWebSocketReturn {
  const {
    projectId,
    onMessage,
    onConnect,
    onDisconnect,
    onError,
    autoConnect = true,
    maxReconnectAttempts = DEFAULT_MAX_RECONNECT_ATTEMPTS,
  } = options;

  const [status, setStatus] = useState<ConnectionStatus>('disconnected');
  const [reconnectAttemptCount, setReconnectAttemptCount] = useState(0);

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const isManualDisconnectRef = useRef(false);
  const messageQueueRef = useRef<ClientMessage[]>([]);
  const reconnectAttemptRef = useRef(0);

  // Use refs for callbacks to avoid re-creating connect/disconnect
  const onMessageRef = useRef(onMessage);
  const onConnectRef = useRef(onConnect);
  const onDisconnectRef = useRef(onDisconnect);
  const onErrorRef = useRef(onError);

  // Update refs when callbacks change
  useEffect(() => {
    onMessageRef.current = onMessage;
    onConnectRef.current = onConnect;
    onDisconnectRef.current = onDisconnect;
    onErrorRef.current = onError;
  }, [onMessage, onConnect, onDisconnect, onError]);

  // Clear any pending reconnect timeout
  const clearReconnectTimeout = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
  }, []);

  // Connect to WebSocket
  const connect = useCallback(() => {
    // Clean up existing connection
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }

    clearReconnectTimeout();
    isManualDisconnectRef.current = false;
    setStatus('connecting');

    const url = getWebSocketUrl(projectId);
    console.log(`Connecting to WebSocket: ${url}`);

    try {
      const ws = new WebSocket(url);
      wsRef.current = ws;

      ws.onopen = () => {
        console.log('WebSocket connected');
        setStatus('connected');
        setReconnectAttemptCount(0);
        reconnectAttemptRef.current = 0;
        onConnectRef.current?.();

        // Send any queued messages
        while (messageQueueRef.current.length > 0) {
          const queuedMessage = messageQueueRef.current.shift();
          if (queuedMessage) {
            ws.send(JSON.stringify(queuedMessage));
          }
        }
      };

      ws.onmessage = (event: MessageEvent) => {
        try {
          const message: ServerMessage = JSON.parse(event.data);
          console.log('WebSocket message received:', message.type);
          onMessageRef.current?.(message);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      ws.onerror = (event: Event) => {
        console.error('WebSocket error:', event);
        onErrorRef.current?.(event);
      };

      ws.onclose = (event: CloseEvent) => {
        console.log(`WebSocket closed: code=${event.code}, reason=${event.reason}`);
        wsRef.current = null;
        setStatus('disconnected');
        onDisconnectRef.current?.();

        // Only attempt reconnection if not manually disconnected
        if (!isManualDisconnectRef.current && event.code !== 1000) {
          const attempt = reconnectAttemptRef.current;
          if (attempt < maxReconnectAttempts) {
            const delay = getReconnectDelay(attempt);
            console.log(`Scheduling reconnection attempt ${attempt + 1} in ${Math.round(delay)}ms`);
            reconnectAttemptRef.current = attempt + 1;
            setReconnectAttemptCount(attempt + 1);
            reconnectTimeoutRef.current = setTimeout(() => {
              connect();
            }, delay);
          } else {
            console.log(`Max reconnection attempts (${maxReconnectAttempts}) reached`);
          }
        }
      };
    } catch (error) {
      console.error('Failed to create WebSocket:', error);
      setStatus('disconnected');

      const attempt = reconnectAttemptRef.current;
      if (attempt < maxReconnectAttempts) {
        const delay = getReconnectDelay(attempt);
        reconnectAttemptRef.current = attempt + 1;
        setReconnectAttemptCount(attempt + 1);
        reconnectTimeoutRef.current = setTimeout(() => {
          connect();
        }, delay);
      }
    }
  }, [projectId, maxReconnectAttempts, clearReconnectTimeout]);

  // Disconnect from WebSocket
  const disconnect = useCallback(() => {
    isManualDisconnectRef.current = true;
    clearReconnectTimeout();

    if (wsRef.current) {
      wsRef.current.close(1000, 'Client disconnect');
      wsRef.current = null;
    }

    setStatus('disconnected');
    setReconnectAttemptCount(0);
    reconnectAttemptRef.current = 0;
    messageQueueRef.current = [];
  }, [clearReconnectTimeout]);

  // Send message through WebSocket
  const sendMessage = useCallback((message: ClientMessage) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    } else if (status === 'connecting') {
      // Queue message to send after connection
      console.log('Queueing message (WebSocket connecting)');
      messageQueueRef.current.push(message);
    } else {
      console.warn('WebSocket not connected, cannot send message');
    }
  }, [status]);

  // Auto-connect on mount and cleanup on unmount
  useEffect(() => {
    if (autoConnect) {
      connect();
    }

    return () => {
      disconnect();
    };
  }, [projectId, autoConnect, connect, disconnect]);

  return {
    status,
    sendMessage,
    connect,
    disconnect,
    reconnectAttempts: reconnectAttemptCount,
  };
}
