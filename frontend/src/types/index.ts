// Discovery types
export * from './discovery';

// Achievement and Learning Journey types
export * from './achievements';

// Agent types for multi-agent chat UI
// Note: 'product' is a legacy alias for 'product_manager' - kept for backward compatibility
// The 'product_manager' type is displayed as "Root" (the discovery/foundation agent)
export type AgentType = 'product_manager' | 'product' | 'designer' | 'developer';

export interface AgentConfig {
  displayName: string;
  shortName: string;
  color: string;
  bgColor: string;
  /** Root's description when introducing this agent */
  rootIntro: string;
  /** The agent's self-introduction message */
  selfIntro: string;
}

export const AGENT_CONFIG: Record<AgentType, AgentConfig> = {
  product_manager: {
    displayName: 'Root',
    shortName: 'Root',
    color: '#0D9488',
    bgColor: '#CCFBF1',
    rootIntro: '', // Root doesn't introduce themselves
    selfIntro: '', // Root doesn't have a self-intro
  }, // Teal - discovery/foundation
  product: {
    displayName: 'Root',
    shortName: 'Root',
    color: '#0D9488',
    bgColor: '#CCFBF1',
    rootIntro: '',
    selfIntro: '',
  }, // Legacy alias
  designer: {
    displayName: 'Bloom',
    shortName: 'Bloom',
    color: '#F97316',
    bgColor: '#FFF7ED',
    rootIntro: "Meet Bloom, our designer - she'll craft the visual experience and user interface.",
    selfIntro: "Hi! I'm excited to design your project. I'll focus on making it intuitive and beautiful.",
  }, // Orange - ideas flourish into designs
  developer: {
    displayName: 'Harvest',
    shortName: 'Harvest',
    color: '#10B981',
    bgColor: '#ECFDF5',
    rootIntro: "And Harvest, our developer - he'll write the code and bring everything to life.",
    selfIntro: "Hey! Ready to build this. I'll handle the technical implementation and make sure everything works smoothly.",
  }, // Green - bringing ideas to fruition
};

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
  agentType?: AgentType;
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
  agentType?: AgentType;
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
