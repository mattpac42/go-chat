// Project types
export interface Project {
  id: string;
  title: string;
  createdAt: string;
  updatedAt: string;
}

// Message types
export interface Message {
  id: string;
  projectId: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: string;
  isStreaming?: boolean;
}

// File types
export interface FileItem {
  id: string;
  path: string;
  filename: string;
  language?: string;
  createdAt: string;
}

export interface FileWithContent extends FileItem {
  projectId: string;
  content: string;
}

export interface FileNode {
  id: string;
  name: string;
  path: string;
  type: 'file' | 'folder';
  language?: string;
  content?: string;
  children?: FileNode[];
}

// WebSocket message types
export interface ClientMessage {
  type: 'chat_message';
  projectId: string;
  content: string;
  timestamp: string;
}

export interface ServerMessage {
  type: 'message_start' | 'message_chunk' | 'message_complete' | 'error';
  projectId: string;
  messageId: string;
  content?: string;
  fullContent?: string;
  error?: string;
}

// Connection status
export type ConnectionStatus = 'connected' | 'connecting' | 'disconnected';

// Chat state
export interface ChatState {
  messages: Message[];
  isLoading: boolean;
  error: string | null;
}
