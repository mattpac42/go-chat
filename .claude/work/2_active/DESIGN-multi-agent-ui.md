# Design: Multi-Agent Frontend Integration (Phase 2C)

**Date**: 2025-12-26
**Status**: Ready for Implementation
**Prerequisite**: AgentContextService (complete), Agent System Prompts (complete)

---

## Overview

This design enables visual agent distinction in the chat UI. The backend already determines which agent (product_manager, designer, developer) should respond; this design covers how to surface that information to users.

---

## 1. Data Model Changes

### 1.1 Backend Message Model

**File**: `/workspace/backend/internal/model/message.go`

```go
// AgentType is already defined in prd.go, reuse it
// const (
//   AgentProductManager AgentType = "product_manager"
//   AgentDesigner       AgentType = "designer"
//   AgentDeveloper      AgentType = "developer"
// )

type Message struct {
    ID         uuid.UUID   `db:"id" json:"id"`
    ProjectID  uuid.UUID   `db:"project_id" json:"projectId,omitempty"`
    Role       Role        `db:"role" json:"role"`
    Content    string      `db:"content" json:"content"`
    CreatedAt  time.Time   `db:"created_at" json:"createdAt"`
    CodeBlocks []CodeBlock `db:"-" json:"codeBlocks,omitempty"`

    // NEW: Agent identification for assistant messages
    AgentType  *AgentType  `db:"agent_type" json:"agentType,omitempty"`
}
```

**Database Migration**:
```sql
ALTER TABLE messages ADD COLUMN agent_type VARCHAR(50);
```

### 1.2 WebSocket Response Changes

**File**: `/workspace/backend/internal/handler/websocket.go`

Add `agentType` to message responses:

```go
type WebSocketMessage struct {
    Type      string     `json:"type"`
    Content   string     `json:"content,omitempty"`
    MessageID string     `json:"messageId,omitempty"`
    Timestamp time.Time  `json:"timestamp"`
    AgentType string     `json:"agentType,omitempty"` // NEW
}

type MessageCompleteResponse struct {
    Type        string            `json:"type"`
    MessageID   string            `json:"messageId"`
    FullContent string            `json:"fullContent"`
    CodeBlocks  []model.CodeBlock `json:"codeBlocks"`
    Timestamp   time.Time         `json:"timestamp"`
    AgentType   string            `json:"agentType,omitempty"` // NEW
}
```

The `agentType` is determined by `AgentContextService.SelectAgent()` and passed through `ChatService.ProcessMessage()`.

### 1.3 Frontend Types

**File**: `/workspace/frontend/src/types/index.ts`

```typescript
// Agent types matching backend
export type AgentType = 'product_manager' | 'designer' | 'developer';

export interface Message {
  id: string;
  projectId: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: string;
  isStreaming?: boolean;

  // NEW: Agent identification
  agentType?: AgentType;
}

// WebSocket messages updated
export interface ServerMessage {
  type: 'message_start' | 'message_chunk' | 'message_complete' | 'error';
  projectId: string;
  messageId: string;
  content?: string;
  fullContent?: string;
  error?: string;
  agentType?: AgentType; // NEW
}
```

### 1.4 First Introduction Tracking

Track which agents the user has "met" in localStorage per project:

```typescript
// Key: `agents_introduced_${projectId}`
// Value: JSON array of AgentType strings
// Example: ["product_manager", "developer"]
```

This is managed in the frontend only - no backend persistence needed for MVP.

---

## 2. Agent Visual Specifications

### 2.1 Agent Configuration

Create a new configuration file:

**File**: `/workspace/frontend/src/config/agents.ts`

```typescript
import { AgentType } from '@/types';

export interface AgentConfig {
  type: AgentType;
  displayName: string;
  shortName: string;      // For mobile/compact views
  accentColor: string;    // Tailwind color class
  accentHex: string;      // For dynamic styling
  icon: 'target' | 'layout' | 'code';
  description: string;
}

export const AGENT_CONFIG: Record<AgentType, AgentConfig> = {
  product_manager: {
    type: 'product_manager',
    displayName: 'Product Guide',
    shortName: 'Guide',
    accentColor: 'violet-600',
    accentHex: '#7C3AED',
    icon: 'target',
    description: 'Guides vision, goals, and scope',
  },
  designer: {
    type: 'designer',
    displayName: 'UX Expert',
    shortName: 'UX',
    accentColor: 'orange-500',
    accentHex: '#F97316',
    icon: 'layout',
    description: 'Interface and user experience design',
  },
  developer: {
    type: 'developer',
    displayName: 'Developer',
    shortName: 'Dev',
    accentColor: 'emerald-500',
    accentHex: '#10B981',
    icon: 'code',
    description: 'Code implementation and features',
  },
};

export function getAgentConfig(type: AgentType): AgentConfig {
  return AGENT_CONFIG[type] || AGENT_CONFIG.developer;
}
```

### 2.2 Visual Treatment Summary

| Element | Specification |
|---------|---------------|
| Left border | 3px solid accent color |
| Icon size | 20x20px (desktop), 16x16px (mobile) |
| Label | `displayName` on desktop, `shortName` on mobile |
| Badge | "NEW" pill on first appearance, violet-100 background |
| Bubble background | Same as current (gray-100) |

---

## 3. Component Structure

### 3.1 New Components

```
frontend/src/components/chat/
  AgentIcon.tsx        # SVG icons for each agent type
  AgentHeader.tsx      # Agent label + icon + optional NEW badge
  MessageBubble.tsx    # Updated to use AgentHeader
```

### 3.2 AgentIcon Component

**File**: `/workspace/frontend/src/components/chat/AgentIcon.tsx`

```typescript
interface AgentIconProps {
  icon: 'target' | 'layout' | 'code';
  className?: string;
}

export function AgentIcon({ icon, className = 'w-5 h-5' }: AgentIconProps) {
  switch (icon) {
    case 'target':
      return (
        <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      );
    case 'layout':
      return (
        <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
            d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z" />
        </svg>
      );
    case 'code':
      return (
        <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
            d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
        </svg>
      );
  }
}
```

### 3.3 AgentHeader Component

**File**: `/workspace/frontend/src/components/chat/AgentHeader.tsx`

```typescript
import { AgentType } from '@/types';
import { getAgentConfig } from '@/config/agents';
import { AgentIcon } from './AgentIcon';

interface AgentHeaderProps {
  agentType: AgentType;
  isFirstAppearance?: boolean;
  timestamp?: string;
  isMobile?: boolean;
}

export function AgentHeader({
  agentType,
  isFirstAppearance,
  timestamp,
  isMobile
}: AgentHeaderProps) {
  const config = getAgentConfig(agentType);

  return (
    <div className="flex items-center gap-2 mb-2">
      <span style={{ color: config.accentHex }}>
        <AgentIcon
          icon={config.icon}
          className={isMobile ? 'w-4 h-4' : 'w-5 h-5'}
        />
      </span>
      <span
        className="text-sm font-medium"
        style={{ color: config.accentHex }}
      >
        {isMobile ? config.shortName : config.displayName}
      </span>
      {isFirstAppearance && (
        <span className="px-1.5 py-0.5 text-xs font-medium bg-violet-100 text-violet-700 rounded">
          NEW
        </span>
      )}
      {timestamp && (
        <span className="ml-auto text-xs text-gray-500">
          {timestamp}
        </span>
      )}
    </div>
  );
}
```

### 3.4 Updated MessageBubble

Key changes to `/workspace/frontend/src/components/chat/MessageBubble.tsx`:

```typescript
interface MessageBubbleProps {
  message: Message;
  showCodeBlocks?: boolean;
  isFirstAgentAppearance?: boolean; // NEW
  isMobile?: boolean;               // NEW
}

export function MessageBubble({
  message,
  showCodeBlocks = false,
  isFirstAgentAppearance = false,
  isMobile = false
}: MessageBubbleProps) {
  const isUser = message.role === 'user';
  const agentConfig = message.agentType ? getAgentConfig(message.agentType) : null;

  return (
    <div className={`flex ${isUser ? 'justify-end' : 'justify-start'} mb-4`}>
      <div
        className={`max-w-[85%] md:max-w-[70%] rounded-2xl px-4 py-3 ${
          isUser
            ? 'bg-teal-400 text-white rounded-br-md'
            : 'bg-gray-100 text-gray-900 rounded-bl-md'
        }`}
        style={
          // Add left border accent for agent messages
          !isUser && agentConfig
            ? { borderLeft: `3px solid ${agentConfig.accentHex}` }
            : undefined
        }
      >
        {/* Agent header for assistant messages */}
        {!isUser && message.agentType && (
          <AgentHeader
            agentType={message.agentType}
            isFirstAppearance={isFirstAgentAppearance}
            isMobile={isMobile}
          />
        )}

        {/* Existing message content rendering... */}
      </div>
    </div>
  );
}
```

---

## 4. First Introduction Flow

### 4.1 Hook for Agent Introduction Tracking

**File**: `/workspace/frontend/src/hooks/useAgentIntroductions.ts`

```typescript
import { useState, useEffect, useCallback } from 'react';
import { AgentType } from '@/types';

const STORAGE_KEY_PREFIX = 'agents_introduced_';

export function useAgentIntroductions(projectId: string) {
  const [introducedAgents, setIntroducedAgents] = useState<Set<AgentType>>(
    new Set()
  );

  // Load from localStorage on mount
  useEffect(() => {
    const stored = localStorage.getItem(`${STORAGE_KEY_PREFIX}${projectId}`);
    if (stored) {
      try {
        const agents = JSON.parse(stored) as AgentType[];
        setIntroducedAgents(new Set(agents));
      } catch {
        // Invalid data, start fresh
      }
    }
  }, [projectId]);

  // Mark agent as introduced
  const markIntroduced = useCallback((agentType: AgentType) => {
    setIntroducedAgents((prev) => {
      if (prev.has(agentType)) return prev;

      const next = new Set(prev);
      next.add(agentType);

      // Persist to localStorage
      localStorage.setItem(
        `${STORAGE_KEY_PREFIX}${projectId}`,
        JSON.stringify(Array.from(next))
      );

      return next;
    });
  }, [projectId]);

  // Check if this is first appearance
  const isFirstAppearance = useCallback(
    (agentType: AgentType) => !introducedAgents.has(agentType),
    [introducedAgents]
  );

  // Reset introductions (for testing or new projects)
  const resetIntroductions = useCallback(() => {
    localStorage.removeItem(`${STORAGE_KEY_PREFIX}${projectId}`);
    setIntroducedAgents(new Set());
  }, [projectId]);

  return {
    introducedAgents,
    markIntroduced,
    isFirstAppearance,
    resetIntroductions,
  };
}
```

### 4.2 Integration in MessageList

Update `MessageList.tsx` to track introductions:

```typescript
import { useAgentIntroductions } from '@/hooks/useAgentIntroductions';

export function MessageList({ messages, projectId, isLoading }: MessageListProps) {
  const { isFirstAppearance, markIntroduced } = useAgentIntroductions(projectId);
  const [isMobile, setIsMobile] = useState(false);

  // Track mobile state
  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 768);
    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  // Process messages to determine first appearances
  const processedMessages = useMemo(() => {
    const seenAgents = new Set<AgentType>();

    return messages.map((message) => {
      if (message.role === 'assistant' && message.agentType) {
        const isFirst = !seenAgents.has(message.agentType) &&
                        isFirstAppearance(message.agentType);
        seenAgents.add(message.agentType);

        // Mark as introduced when displayed
        if (isFirst) {
          markIntroduced(message.agentType);
        }

        return { ...message, isFirstAgentAppearance: isFirst };
      }
      return { ...message, isFirstAgentAppearance: false };
    });
  }, [messages, isFirstAppearance, markIntroduced]);

  return (
    <div ref={containerRef} className="flex-1 overflow-y-auto p-4">
      {processedMessages.map((message) => (
        <MessageBubble
          key={message.id}
          message={message}
          isFirstAgentAppearance={message.isFirstAgentAppearance}
          isMobile={isMobile}
        />
      ))}
      {/* ... loading indicator ... */}
    </div>
  );
}
```

---

## 5. Backend Integration Points

### 5.1 ChatService Changes

**File**: `/workspace/backend/internal/service/chat.go`

The `ProcessMessage` method needs to:

1. Call `AgentContextService.GetContextForMessage()` to determine agent
2. Include agent type in the response

```go
type ProcessMessageResult struct {
    Content    string
    CodeBlocks []model.CodeBlock
    AgentType  model.AgentType  // NEW
}

func (s *ChatService) ProcessMessage(
    ctx context.Context,
    projectID uuid.UUID,
    message string,
    onChunk func(string),
) (*ProcessMessageResult, error) {
    // Get agent context
    agentCtx, err := s.agentContextService.GetContextForMessage(ctx, projectID, message)
    if err != nil {
        // Fall back to developer if context service fails
        agentCtx = &model.AgentContext{Agent: model.AgentDeveloper}
    }

    // Get system prompt for selected agent
    systemPrompt, _ := s.agentContextService.GetSystemPrompt(ctx, agentCtx)

    // ... existing Claude API call with systemPrompt ...

    return &ProcessMessageResult{
        Content:    fullContent,
        CodeBlocks: codeBlocks,
        AgentType:  agentCtx.Agent,
    }, nil
}
```

### 5.2 WebSocket Handler Changes

**File**: `/workspace/backend/internal/handler/websocket.go`

Pass agent type to `message_start` and `message_complete`:

```go
func (h *WebSocketHandler) handleChatMessage(...) {
    // ... existing code ...

    result, err := h.chatService.ProcessMessage(chatCtx, projectID, msg.Content, onChunk)
    if err != nil {
        // ...
    }

    // Include agent type in message_start (for streaming indicator)
    h.sendMessageStart(conn, mu, messageID, string(result.AgentType))

    // Include agent type in message_complete
    h.sendMessageComplete(conn, mu, messageID, result.Content, result.CodeBlocks, string(result.AgentType))
}
```

---

## 6. One-Voice Rule Enforcement

The "one agent speaks at a time" rule is enforced entirely in the backend:

1. `AgentContextService.SelectAgent()` returns exactly one agent
2. No multi-agent responses in a single message
3. Hand-offs between agents are expressed in message content, not separate messages

The frontend simply displays whatever agent type is specified in each message.

---

## 7. Implementation Checklist

### Backend Tasks
- [ ] Add `agent_type` column to messages table (migration)
- [ ] Update Message model with AgentType field
- [ ] Update WebSocket message types with agentType
- [ ] Wire AgentContextService into ChatService
- [ ] Update ProcessMessage to return agent type
- [ ] Update websocket handler to include agent type in responses

### Frontend Tasks
- [ ] Add AgentType to frontend types
- [ ] Create `/config/agents.ts` configuration
- [ ] Create `AgentIcon.tsx` component
- [ ] Create `AgentHeader.tsx` component
- [ ] Update `MessageBubble.tsx` with agent styling
- [ ] Create `useAgentIntroductions` hook
- [ ] Update `MessageList.tsx` to pass introduction props
- [ ] Update `useChat` hook to handle agentType in server messages
- [ ] Add mobile responsiveness for agent headers

### Testing
- [ ] Verify agent header displays correctly for each agent type
- [ ] Verify NEW badge appears only on first appearance per project
- [ ] Verify mobile compact view shows short names
- [ ] Verify left border accent colors match config
- [ ] Test localStorage persistence of introduced agents

---

## 8. File Changes Summary

| File | Change Type | Description |
|------|-------------|-------------|
| `backend/internal/model/message.go` | Modify | Add AgentType field |
| `backend/internal/handler/websocket.go` | Modify | Add agentType to responses |
| `backend/internal/service/chat.go` | Modify | Integrate AgentContextService |
| `frontend/src/types/index.ts` | Modify | Add AgentType, update Message |
| `frontend/src/config/agents.ts` | New | Agent configuration |
| `frontend/src/components/chat/AgentIcon.tsx` | New | Agent icon component |
| `frontend/src/components/chat/AgentHeader.tsx` | New | Agent header component |
| `frontend/src/components/chat/MessageBubble.tsx` | Modify | Add agent styling |
| `frontend/src/components/chat/MessageList.tsx` | Modify | Track introductions |
| `frontend/src/hooks/useAgentIntroductions.ts` | New | Introduction tracking hook |
| `frontend/src/hooks/useChat.ts` | Modify | Handle agentType in messages |

---

## 9. Future Considerations (Not MVP)

- **@mentions**: Allow users to directly address specific agents
- **Team drawer**: Mobile drawer showing all available agents
- **Agent hand-off messages**: Explicit "Let me bring in our UX expert..." transitions
- **Agent availability states**: Show which agents are "active" vs "available"

These are documented in the strategic UX design but deferred for MVP.
