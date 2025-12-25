'use client';

import ReactMarkdown from 'react-markdown';
import { Message } from '@/types';
import { CodeBlock } from './CodeBlock';

interface MessageBubbleProps {
  message: Message;
}

// Get timestamp string from message, handling both 'timestamp' and 'createdAt' fields
function getTimestamp(message: Message): string {
  // Check for timestamp first (from WebSocket messages)
  if (message.timestamp) {
    return message.timestamp;
  }
  // Fall back to createdAt (from API messages)
  if ('createdAt' in message && (message as { createdAt?: string }).createdAt) {
    return (message as { createdAt: string }).createdAt;
  }
  // If neither exists, return current time
  return new Date().toISOString();
}

// Format timestamp for display
function formatTime(timestamp: string): string {
  const date = new Date(timestamp);
  if (isNaN(date.getTime())) {
    return '';
  }
  return date.toLocaleTimeString([], {
    hour: '2-digit',
    minute: '2-digit',
  });
}

export function MessageBubble({ message }: MessageBubbleProps) {
  const isUser = message.role === 'user';
  const timestamp = getTimestamp(message);
  const formattedTime = formatTime(timestamp);

  return (
    <div
      className={`flex ${isUser ? 'justify-end' : 'justify-start'} mb-4`}
      data-testid="message-bubble"
      data-role={message.role}
    >
      <div
        className={`max-w-[85%] md:max-w-[70%] rounded-2xl px-4 py-3 ${
          isUser
            ? 'bg-teal-400 text-white rounded-br-md'
            : 'bg-gray-100 text-gray-900 rounded-bl-md'
        }`}
      >
        <div className={`prose prose-sm max-w-none break-words ${
          isUser
            ? 'prose-invert prose-p:text-white prose-headings:text-white prose-strong:text-white prose-code:text-white'
            : 'prose-gray'
        }`}>
          <ReactMarkdown
            components={{
              // Custom code block rendering
              code: ({ className, children, ...props }) => {
                const match = /language-(\w+)/.exec(className || '');
                const isInline = !match && !className;

                if (isInline) {
                  return (
                    <code
                      className={`${isUser ? 'bg-teal-500/50' : 'bg-gray-200'} px-1.5 py-0.5 rounded text-sm`}
                      {...props}
                    >
                      {children}
                    </code>
                  );
                }

                const codeContent = String(children).replace(/\n$/, '');
                return (
                  <CodeBlock code={codeContent} language={match?.[1] || 'text'} />
                );
              },
              // Handle pre elements to avoid double wrapping
              pre: ({ children }) => <>{children}</>,
              // Style paragraphs
              p: ({ children }) => (
                <p className="mb-2 last:mb-0 whitespace-pre-wrap">{children}</p>
              ),
              // Style headers
              h1: ({ children }) => (
                <h1 className={`text-xl font-bold mb-2 ${isUser ? 'text-white' : 'text-gray-900'}`}>{children}</h1>
              ),
              h2: ({ children }) => (
                <h2 className={`text-lg font-bold mb-2 ${isUser ? 'text-white' : 'text-gray-900'}`}>{children}</h2>
              ),
              h3: ({ children }) => (
                <h3 className={`text-base font-bold mb-2 ${isUser ? 'text-white' : 'text-gray-900'}`}>{children}</h3>
              ),
              // Style lists
              ul: ({ children }) => (
                <ul className="list-disc list-inside mb-2 space-y-1">{children}</ul>
              ),
              ol: ({ children }) => (
                <ol className="list-decimal list-inside mb-2 space-y-1">{children}</ol>
              ),
              // Style strong/bold
              strong: ({ children }) => (
                <strong className={`font-bold ${isUser ? 'text-white' : 'text-gray-900'}`}>{children}</strong>
              ),
              // Style emphasis/italic
              em: ({ children }) => (
                <em className="italic">{children}</em>
              ),
              // Style links
              a: ({ href, children }) => (
                <a
                  href={href}
                  className={`underline ${isUser ? 'text-white hover:text-teal-100' : 'text-teal-600 hover:text-teal-800'}`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {children}
                </a>
              ),
            }}
          >
            {message.content}
          </ReactMarkdown>
        </div>

        {/* Streaming indicator */}
        {message.isStreaming && (
          <span className="inline-block ml-1 animate-pulse">|</span>
        )}

        {/* Timestamp */}
        {formattedTime && (
          <div
            className={`text-xs mt-2 ${
              isUser ? 'text-teal-100' : 'text-gray-500'
            }`}
          >
            {formattedTime}
          </div>
        )}
      </div>
    </div>
  );
}
