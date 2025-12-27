# Architect Session: Chat Input Enhancements

**Date**: 2025-12-27T19:45:00Z
**Agent**: architect
**Task**: Assess technical feasibility and provide implementation approach for clipboard paste screenshots and up-arrow message history

## Work Completed

Analyzed current architecture:
- `/workspace/frontend/src/components/chat/ChatInput.tsx` - Simple textarea component with send handler
- `/workspace/frontend/src/components/chat/ChatContainer.tsx` - Parent orchestrator with discovery integration
- `/workspace/frontend/src/hooks/useChat.ts` - WebSocket-based chat state management
- `/workspace/frontend/src/types/index.ts` - Type definitions including Message interface

## Decisions Made

### Feature 1: Clipboard Paste for Screenshots

**Feasibility**: Fully feasible with standard Web APIs

**Architecture Recommendation**:

1. **Clipboard Event Handling**
   - Use `onPaste` event on the textarea
   - Access `clipboardData.items` to detect image content
   - Images come as `Blob` objects via `getAsFile()`

2. **Data Flow**
   ```
   Paste Event -> Extract Blob -> Create Object URL for preview -> Convert to Base64 -> Upload
   ```

3. **Component Changes**:
   - Extend `ChatInput` to accept `onImagePaste?: (file: File) => void` callback
   - Add preview state for pending image
   - Modify `onSend` signature to accept `{ content: string, attachment?: File }`

4. **Type Extension** (in `/workspace/frontend/src/types/index.ts`):
   ```typescript
   export interface ChatAttachment {
     type: 'image';
     file: File;
     previewUrl: string;  // Object URL for display
   }

   export interface SendMessagePayload {
     content: string;
     attachments?: ChatAttachment[];
   }
   ```

5. **Backend Integration**:
   - Option A: Base64 encode and send inline with WebSocket message
   - Option B: Upload to storage first, send URL reference
   - Recommend Option B for larger files (>100KB)

### Feature 2: Up Arrow Message History

**Feasibility**: Fully feasible, low complexity

**Architecture Recommendation**:

1. **Storage Strategy** (Decision: React state + sessionStorage)
   - Use React state for current session history
   - Persist to sessionStorage for browser-session persistence
   - Do NOT use localStorage (would persist indefinitely)
   - Limit to 50 messages (sufficient for session, minimal memory)

2. **Custom Hook**: Create `useMessageHistory(projectId)`
   ```typescript
   interface UseMessageHistoryReturn {
     history: string[];
     historyIndex: number;
     navigateHistory: (direction: 'up' | 'down') => string | null;
     addToHistory: (message: string);
     resetNavigation: () => void;
   }
   ```

3. **Behavior Specification**:
   - Up arrow: Only activate when cursor at position 0 (start of input)
   - Down arrow: Only when navigating through history
   - Store original input when user starts navigating
   - Reset index on send or when user types new content

4. **Edge Cases**:
   - Cursor not at start: Normal arrow key behavior
   - Partial text exists: Store as "draft" before navigating
   - Empty history: No-op
   - At history boundaries: No-op (don't wrap)

5. **Session Scoping**:
   - Key by projectId: `chat-history-${projectId}`
   - Each project has independent history

## Files Modified

None - architecture assessment only

## Recommendations

### Implementation Priority
1. **Up Arrow History** (Low effort, high value) - 2-4 hours
2. **Clipboard Paste** (Medium effort, medium value) - 4-8 hours

### Implementation Approach - Up Arrow History

```typescript
// /workspace/frontend/src/hooks/useMessageHistory.ts
import { useState, useCallback, useEffect } from 'react';

const MAX_HISTORY = 50;
const STORAGE_KEY_PREFIX = 'chat-history-';

export function useMessageHistory(projectId: string) {
  const [history, setHistory] = useState<string[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const [draftMessage, setDraftMessage] = useState('');

  // Load from sessionStorage on mount
  useEffect(() => {
    const stored = sessionStorage.getItem(`${STORAGE_KEY_PREFIX}${projectId}`);
    if (stored) {
      setHistory(JSON.parse(stored));
    }
  }, [projectId]);

  // Save to sessionStorage on change
  useEffect(() => {
    if (history.length > 0) {
      sessionStorage.setItem(
        `${STORAGE_KEY_PREFIX}${projectId}`,
        JSON.stringify(history)
      );
    }
  }, [history, projectId]);

  const addToHistory = useCallback((message: string) => {
    setHistory(prev => {
      // Don't add duplicates of the most recent message
      if (prev[0] === message) return prev;
      const updated = [message, ...prev].slice(0, MAX_HISTORY);
      return updated;
    });
    setHistoryIndex(-1);
    setDraftMessage('');
  }, []);

  const navigateHistory = useCallback((
    direction: 'up' | 'down',
    currentMessage: string
  ): string | null => {
    if (history.length === 0) return null;

    if (direction === 'up') {
      // Save current input as draft when first navigating
      if (historyIndex === -1) {
        setDraftMessage(currentMessage);
      }
      const newIndex = Math.min(historyIndex + 1, history.length - 1);
      if (newIndex !== historyIndex) {
        setHistoryIndex(newIndex);
        return history[newIndex];
      }
    } else {
      if (historyIndex > 0) {
        const newIndex = historyIndex - 1;
        setHistoryIndex(newIndex);
        return history[newIndex];
      } else if (historyIndex === 0) {
        setHistoryIndex(-1);
        return draftMessage;
      }
    }
    return null;
  }, [history, historyIndex, draftMessage]);

  const resetNavigation = useCallback(() => {
    setHistoryIndex(-1);
    setDraftMessage('');
  }, []);

  return {
    history,
    historyIndex,
    addToHistory,
    navigateHistory,
    resetNavigation,
    isNavigating: historyIndex >= 0,
  };
}
```

### Implementation Approach - Clipboard Paste

```typescript
// Addition to ChatInput.tsx

const handlePaste = useCallback((e: React.ClipboardEvent) => {
  const items = e.clipboardData.items;

  for (let i = 0; i < items.length; i++) {
    if (items[i].type.indexOf('image') !== -1) {
      e.preventDefault();
      const file = items[i].getAsFile();
      if (file && onImagePaste) {
        onImagePaste(file);
      }
      return;
    }
  }
  // Let text paste through normally
}, [onImagePaste]);

// In render:
<textarea
  onPaste={handlePaste}
  // ... other props
/>
```

**Preview Component Pattern**:
```typescript
interface ImagePreviewProps {
  file: File;
  onRemove: () => void;
}

function ImagePreview({ file, onRemove }: ImagePreviewProps) {
  const [previewUrl, setPreviewUrl] = useState<string>('');

  useEffect(() => {
    const url = URL.createObjectURL(file);
    setPreviewUrl(url);
    return () => URL.revokeObjectURL(url);
  }, [file]);

  return (
    <div className="relative inline-block">
      <img src={previewUrl} className="max-h-32 rounded" alt="Preview" />
      <button onClick={onRemove} className="absolute -top-2 -right-2">
        <XIcon />
      </button>
    </div>
  );
}
```

### Backend Considerations

For clipboard images, the backend needs:
1. Endpoint to receive image uploads (multipart/form-data or base64)
2. Storage solution (local filesystem, S3, etc.)
3. Return URL/ID that can be referenced in chat messages
4. Message type extension to support attachments

### ADR Summary

| Decision | Choice | Rationale |
|----------|--------|-----------|
| History storage | sessionStorage | Browser-session persistence, clears on close |
| History limit | 50 messages | Balance of utility vs memory |
| History scope | Per project | Independent context per project |
| Image transfer | Separate upload | Better for large files, progress tracking |
| Preview storage | Object URL | Efficient, auto-cleanup on revoke |

## Next Steps for Developer

1. Create `useMessageHistory` hook
2. Integrate into `ChatInput` with keyboard event handling
3. Create image paste handler and preview component
4. Extend `onSend` signature for attachments
5. Coordinate with backend on upload endpoint
