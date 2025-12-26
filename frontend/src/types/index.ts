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
  // Metadata from App Map
  shortDescription?: string;
  longDescription?: string;
  functionalGroup?: string;
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
  /** Tier 1: Human-readable short description of what the file does */
  shortDescription?: string;
  /** Tier 1: Longer explanation of the file's purpose and why it exists */
  longDescription?: string;
  /** Functional group for App Map organization (e.g., "Homepage", "Backend Services") */
  functionalGroup?: string;
}

/**
 * File metadata for the 2-tier reveal system
 * Tier 1: Human-readable descriptions (shown by default)
 * Tier 2: Actual code (shown on user request)
 */
export interface FileMetadata {
  shortDescription: string;
  longDescription: string;
  language: string;
  functionalGroup: string;
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
