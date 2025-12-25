# PRD-001: Chat UI (Frontend)

**Version**: 1.0
**Created**: 2025-12-24
**Author**: Product Manager (Tactical)
**Status**: Draft
**Phase**: 1 - Foundation (Weeks 1-2)

---

## Overview

### Problem Statement

Non-technical users need a simple, intuitive way to describe applications they want to build. Traditional development tools require coding knowledge, complex interfaces, and desktop computers. Our target user (Sam the Small Business Owner) checks things from their phone throughout the day and has zero tolerance for technical complexity.

### Proposed Solution

A mobile-first, responsive chat interface that allows users to describe applications in plain English and receive AI-generated code responses. The interface uses real-time WebSocket connections for streaming responses and provides clear visual feedback during code generation.

### Goals

1. Provide a frictionless chat experience that works seamlessly on mobile browsers
2. Display AI-generated code with proper syntax highlighting and formatting
3. Enable real-time streaming of AI responses for immediate feedback
4. Support multiple projects with easy switching between conversations

### Non-Goals

- Native mobile applications (iOS/Android apps)
- User authentication and multi-tenancy
- Offline functionality
- Advanced code editing capabilities
- File upload/attachment support
- Voice input for messages

---

## User Stories

### Primary User

**Sam the Small Business Owner**
- Needs to build simple tools (inventory tracker, customer list, basic dashboard)
- Checks on things from phone throughout day
- Zero tolerance for technical complexity
- Uses iOS Safari or Android Chrome

### User Stories

1. **US-001**: As a small business owner, I want to open Go Chat on my phone and immediately start describing what I need, so that I don't waste time with setup or learning a complex interface.

2. **US-002**: As a user, I want to type a description of an app in plain English and see the AI's response stream in real-time, so that I know the system is working and can read the response as it generates.

3. **US-003**: As a user, I want to see generated code with proper formatting and syntax highlighting, so that I can recognize it as code even if I don't understand it.

4. **US-004**: As a user, I want clear visual feedback when the AI is "thinking" or generating a response, so that I know to wait and the app hasn't frozen.

5. **US-005**: As a user, I want to start multiple project conversations and switch between them, so that I can work on different ideas without losing previous context.

6. **US-006**: As a user, I want the chat interface to reconnect automatically if my connection drops, so that I don't lose my conversation when moving between WiFi and cellular.

---

## Requirements

### Functional Requirements

| ID | Requirement | Priority | Acceptance Criteria |
|----|-------------|----------|---------------------|
| FR-001 | Chat message input with send button | High | User can type messages up to 4000 characters and send via button or Enter key |
| FR-002 | Message display (user and AI) | High | Messages display in chronological order with clear visual distinction between user and AI |
| FR-003 | Real-time response streaming | High | AI responses stream character-by-character or chunk-by-chunk with < 100ms perceived latency |
| FR-004 | Code block rendering | High | Code blocks render with syntax highlighting, language label, and copy button |
| FR-005 | Loading/generating state | High | Visual indicator (animation, text) shows when AI is processing |
| FR-006 | Project list view | High | User can see list of their projects/conversations with titles and timestamps |
| FR-007 | Project switching | High | User can switch between projects without page refresh |
| FR-008 | New project creation | High | User can start a new project conversation with one tap |
| FR-009 | Mobile-responsive layout | High | UI adapts to viewport sizes from 320px to 2560px width |
| FR-010 | Connection status indicator | Medium | User sees current connection state (connected, connecting, disconnected) |
| FR-011 | Auto-reconnection | Medium | WebSocket reconnects automatically with exponential backoff (max 5 attempts) |
| FR-012 | Message history loading | Medium | Previous messages load when opening an existing project |
| FR-013 | Scroll to bottom on new message | Medium | Chat auto-scrolls to show latest message unless user has scrolled up |
| FR-014 | Empty state for new project | Low | New project shows helpful prompt suggesting what to type |
| FR-015 | Error message display | Medium | User-friendly error messages for failed operations |

### Non-Functional Requirements

| ID | Requirement | Criteria |
|----|-------------|----------|
| NFR-001 | Performance - Initial Load | First contentful paint < 2 seconds on 4G connection |
| NFR-002 | Performance - Interaction | Input response < 50ms, no jank during scrolling |
| NFR-003 | Browser Support | iOS Safari 15+, Android Chrome 90+, Desktop Chrome/Firefox/Safari latest |
| NFR-004 | Accessibility | WCAG 2.1 Level AA compliance for core chat functionality |
| NFR-005 | Mobile Touch | All interactive elements minimum 44x44px touch target |
| NFR-006 | Network Resilience | Graceful handling of network interruptions with user feedback |
| NFR-007 | Security | No sensitive data in localStorage; WebSocket over WSS in production |

---

## Technical Requirements

### Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Framework | Next.js 14+ (App Router) | Server components, streaming, React Server Components |
| Language | TypeScript | Type safety, better developer experience |
| Styling | Tailwind CSS | Mobile-first utility classes, small bundle size |
| State Management | React hooks (useState, useReducer, useContext) | Sufficient for MVP scope |
| WebSocket Client | Native WebSocket API or socket.io-client | Real-time communication |
| Code Highlighting | highlight.js or Prism.js | Syntax highlighting for code blocks |
| HTTP Client | Native fetch API | Built-in, no dependencies needed |

### Architecture

```
src/
├── app/
│   ├── layout.tsx          # Root layout with providers
│   ├── page.tsx            # Home/chat page
│   └── projects/
│       └── [id]/
│           └── page.tsx    # Individual project chat
├── components/
│   ├── chat/
│   │   ├── ChatContainer.tsx
│   │   ├── MessageList.tsx
│   │   ├── MessageBubble.tsx
│   │   ├── ChatInput.tsx
│   │   ├── CodeBlock.tsx
│   │   └── TypingIndicator.tsx
│   ├── projects/
│   │   ├── ProjectList.tsx
│   │   └── ProjectCard.tsx
│   └── shared/
│       ├── ConnectionStatus.tsx
│       └── LoadingSpinner.tsx
├── hooks/
│   ├── useWebSocket.ts
│   ├── useProjects.ts
│   └── useChat.ts
├── lib/
│   ├── api.ts              # REST API client
│   └── websocket.ts        # WebSocket connection manager
└── types/
    └── index.ts            # TypeScript interfaces
```

### Key Components

- **ChatContainer**: Main chat view orchestrating message list, input, and status
- **MessageBubble**: Renders individual messages with markdown/code parsing
- **CodeBlock**: Syntax-highlighted code with copy functionality
- **useWebSocket**: Hook managing WebSocket lifecycle, reconnection, and message handling
- **ProjectList**: Sidebar/drawer showing available projects

### WebSocket Message Protocol

```typescript
// Client to Server
interface ClientMessage {
  type: 'chat_message';
  projectId: string;
  content: string;
  timestamp: string;
}

// Server to Client
interface ServerMessage {
  type: 'message_start' | 'message_chunk' | 'message_complete' | 'error';
  projectId: string;
  messageId: string;
  content?: string;        // For chunks
  fullContent?: string;    // For complete message
  error?: string;
}
```

---

## Design Specifications

### Mobile Layout (320px - 768px)

- Full-screen chat with slide-out project drawer
- Sticky header with project name and menu button
- Floating input bar at bottom with safe area insets
- Messages take full width with appropriate padding

### Desktop Layout (769px+)

- Split view: Project list sidebar (280px) + Chat area
- Fixed sidebar, scrollable chat area
- Input bar at bottom of chat area

### Visual Design Tokens

| Token | Value | Usage |
|-------|-------|-------|
| User message bg | #3B82F6 (blue-500) | User message bubbles |
| AI message bg | #F3F4F6 (gray-100) | AI message bubbles |
| Code block bg | #1F2937 (gray-800) | Code block background |
| Primary action | #2563EB (blue-600) | Send button, primary actions |
| Text primary | #111827 (gray-900) | Main text |
| Text secondary | #6B7280 (gray-500) | Timestamps, metadata |

---

## API Dependencies

### REST Endpoints Required (from PRD-002)

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/projects` | GET | List all projects |
| `/api/projects` | POST | Create new project |
| `/api/projects/:id` | GET | Get project with messages |
| `/api/projects/:id` | DELETE | Delete project |

### WebSocket Endpoint

| Endpoint | Purpose |
|----------|---------|
| `ws://[host]/ws/chat` | Real-time chat streaming |

---

## Acceptance Criteria

### Phase 1 Exit Criteria (from Roadmap)

- [ ] User can type messages and see responses
- [ ] AI responds to messages with context
- [ ] Code is viewable in the interface
- [ ] Works on mobile browser (iOS Safari, Android Chrome)

### Detailed Acceptance Tests

1. **AC-001**: On mobile device, user can open app, type "Build me an inventory tracker", and see AI response stream in real-time
2. **AC-002**: AI response containing code blocks displays with syntax highlighting and visible copy button
3. **AC-003**: User can create 3 different projects and switch between them, each maintaining separate conversation history
4. **AC-004**: When network disconnects, user sees "Reconnecting..." indicator; when restored, chat resumes without lost messages
5. **AC-005**: On iPhone SE (375px width), all UI elements are usable without horizontal scrolling
6. **AC-006**: Send button and message input remain visible when mobile keyboard is open

---

## Dependencies

### Internal Dependencies

- PRD-002 Backend API must provide REST and WebSocket endpoints
- Backend WebSocket must support streaming responses

### External Dependencies

- Claude API (via backend) for AI responses
- WebSocket browser support (covered by target browsers)

---

## Out of Scope

The following are explicitly NOT included in this PRD:

1. User authentication/login
2. Message editing or deletion
3. File attachments or image uploads
4. Voice input
5. Offline mode / PWA functionality
6. Export or sharing of conversations
7. Dark mode (consider for future)
8. Keyboard shortcuts beyond Enter to send
9. Message search functionality
10. Typing indicators for multi-user (single user MVP)

---

## Risks and Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| WebSocket connection instability on mobile | High | Medium | Implement robust reconnection with exponential backoff; queue messages during disconnect |
| Large code blocks cause scroll issues | Medium | Medium | Implement horizontal scroll for code blocks; limit initial display height |
| Input hidden by mobile keyboard | High | Low | Use viewport units with keyboard detection; test on real devices |
| Slow initial load on 3G | Medium | Medium | Implement code splitting; lazy load syntax highlighter |

---

## Success Metrics

| Metric | Target | Measurement Method |
|--------|--------|-------------------|
| Time to first message sent | < 10 seconds | Analytics from app open to first send |
| Message streaming perceived latency | < 100ms | Time from send to first response chunk |
| Mobile usability score | 100% critical paths | Manual testing checklist |
| Connection recovery rate | > 95% | Automatic reconnection success rate |

---

## Timeline

| Milestone | Description | Target |
|-----------|-------------|--------|
| Design complete | Figma/sketches approved | Day 2 |
| Core chat functional | Send/receive messages working | Day 5 |
| WebSocket streaming | Real-time response display | Day 7 |
| Project switching | Multi-project support | Day 9 |
| Mobile polish | iOS/Android testing complete | Day 12 |
| Phase 1 complete | All acceptance criteria met | Day 14 |

---

## Approval

| Role | Name | Date | Status |
|------|------|------|--------|
| Product | | | [ ] |
| Engineering | | | [ ] |
| Design | | | [ ] |

---

**Document History**

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2025-12-24 | 1.0 | Initial PRD | Product Manager |
