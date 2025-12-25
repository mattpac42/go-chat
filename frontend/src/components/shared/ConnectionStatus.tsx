'use client';

import { ConnectionStatus as ConnectionStatusType } from '@/types';

interface ConnectionStatusProps {
  status: ConnectionStatusType;
  reconnectAttempts?: number;
  onReconnect?: () => void;
}

export function ConnectionStatus({
  status,
  reconnectAttempts = 0,
  onReconnect,
}: ConnectionStatusProps) {
  const config = {
    connected: {
      color: 'bg-green-500',
      label: 'Connected',
      textColor: 'text-gray-500',
    },
    connecting: {
      color: 'bg-yellow-500',
      label: reconnectAttempts > 0 ? `Reconnecting (${reconnectAttempts}/5)...` : 'Connecting...',
      textColor: 'text-yellow-600',
    },
    disconnected: {
      color: 'bg-red-500',
      label: 'Disconnected',
      textColor: 'text-red-600',
    },
  };

  const { color, label, textColor } = config[status];

  return (
    <div className="flex items-center gap-2">
      <span
        className={`w-2 h-2 rounded-full ${color} ${
          status === 'connecting' ? 'animate-pulse' : ''
        }`}
        aria-hidden="true"
      />
      <span className={`text-sm ${textColor}`}>{label}</span>
      {status === 'disconnected' && onReconnect && (
        <button
          onClick={onReconnect}
          className="text-sm text-teal-500 hover:text-teal-700 font-medium ml-2"
        >
          Retry
        </button>
      )}
    </div>
  );
}
