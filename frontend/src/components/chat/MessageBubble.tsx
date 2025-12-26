'use client';

import { useState } from 'react';
import ReactMarkdown from 'react-markdown';
import { Message, AGENT_CONFIG } from '@/types';
import { CodeBlock } from './CodeBlock';
import { AgentHeader } from './AgentHeader';

interface MessageBubbleProps {
  message: Message;
  showCodeBlocks?: boolean; // If false, replace code blocks with friendly summary
  showBadge?: boolean; // Show "NEW" badge for first agent introduction
}

function ChevronIcon({ isOpen, className }: { isOpen: boolean; className?: string }) {
  return (
    <svg
      className={`w-4 h-4 transition-transform duration-200 ${isOpen ? 'rotate-90' : ''} ${className || ''}`}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9 5l7 7-7 7"
      />
    </svg>
  );
}

/**
 * Simple preprocessing for user messages - just clean up extra newlines
 */
function preprocessMarkdown(content: string): string {
  return content.replace(/\n{3,}/g, '\n\n');
}

/**
 * Extract file names from code blocks and create a summary message
 * Also removes code blocks from display since files are shown in the Files panel
 */
function processAssistantContent(content: string, showCode: boolean): string {
  // Extract filenames from code blocks with our metadata format
  // Pattern: ```lang:filename followed by --- and short_description
  const filenamePattern = /```(\w+):([^\n]+)\n---[\s\S]*?short_description:/g;
  const filenameMatches = Array.from(content.matchAll(filenamePattern));
  const filenames = filenameMatches.map(m => m[2].trim());

  // Remove duplicates
  const uniqueFilenames = Array.from(new Set(filenames));

  // If we found files with metadata, remove those code blocks from display
  // Strategy: Find the FIRST code block with metadata and remove everything from there to the end
  // This is simpler and more reliable than trying to find matching closing fences
  // (which is hard when the content contains nested code blocks)
  let cleanContent = content;

  if (uniqueFilenames.length > 0) {
    // Find the first occurrence of a code block with our metadata format
    const metadataBlockStart = /```\w+:[^\n]+\n---[\s\S]*?short_description:/;
    const match = metadataBlockStart.exec(cleanContent);

    if (match) {
      // Keep only the content BEFORE the code block
      cleanContent = cleanContent.substring(0, match.index).trim();
    }

    // Clean up any stray backticks and extra newlines
    cleanContent = cleanContent
      .replace(/\n{3,}/g, '\n\n')
      .replace(/^\s*`\s*$/gm, '')
      .replace(/`\s*$/g, '')
      .trim();
  }

  // If showCode is true, return whatever content we have
  if (showCode) {
    return cleanContent || content;
  }

  // If no files found, just return clean content
  if (uniqueFilenames.length === 0) {
    return cleanContent || content;
  }

  // Determine what kind of files are being created
  const hasHtml = uniqueFilenames.some(f => f.endsWith('.html'));
  const hasCss = uniqueFilenames.some(f => f.endsWith('.css'));
  const hasJs = uniqueFilenames.some(f => f.endsWith('.js') || f.endsWith('.ts'));
  const hasMd = uniqueFilenames.some(f => f.endsWith('.md'));

  let fileType = 'files';
  if (hasHtml && hasCss && hasJs) {
    fileType = 'page structure';
  } else if (hasMd) {
    fileType = 'documentation';
  } else if (hasHtml) {
    fileType = 'HTML structure';
  } else if (hasCss) {
    fileType = 'styles';
  } else if (hasJs) {
    fileType = 'functionality';
  }

  // Create a friendly summary with filenames as inline code
  const fileList = uniqueFilenames.map(f => `\`${f}\``).join(', ');
  const summary = `\n\n📁 **Created ${uniqueFilenames.length} ${uniqueFilenames.length === 1 ? 'file' : 'files'}:** ${fileList}\n\n_Check the Files panel to see what each file does →_`;

  // If there's other content, keep it and append summary
  if (cleanContent) {
    return cleanContent + summary;
  }

  // If only code blocks, show a friendly message
  return `Creating your ${fileType}...` + summary;
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

export function MessageBubble({
  message,
  showCodeBlocks = false,
  showBadge = false,
}: MessageBubbleProps) {
  const isUser = message.role === 'user';
  const timestamp = getTimestamp(message);
  const formattedTime = formatTime(timestamp);
  const [showRawContent, setShowRawContent] = useState(false);

  // During streaming, show a generating message instead of raw content
  const isStreaming = message.isStreaming && !isUser;

  // Process content to hide code blocks for assistant messages
  // Preprocess user content too to fix markdown issues
  const displayContent = isUser
    ? preprocessMarkdown(message.content)
    : processAssistantContent(message.content, showCodeBlocks);

  // Get agent config for styling (only for assistant messages with agentType)
  const agentConfig = !isUser && message.agentType ? AGENT_CONFIG[message.agentType] : null;

  // Build style object for agent left border
  const bubbleStyle = agentConfig
    ? { borderLeft: `3px solid ${agentConfig.color}` }
    : undefined;

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
        style={bubbleStyle}
      >
        {/* Agent header for assistant messages with agentType */}
        {!isUser && message.agentType && !isStreaming && (
          <AgentHeader agentType={message.agentType} showBadge={showBadge} />
        )}

        {/* Streaming state: show generating message with expandable raw content */}
        {isStreaming ? (
          <div>
            {/* Generating indicator */}
            <div className="flex items-center gap-2 text-gray-600">
              <span className="inline-flex">
                <span className="w-2 h-2 bg-teal-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
                <span className="w-2 h-2 bg-teal-400 rounded-full animate-bounce ml-1" style={{ animationDelay: '150ms' }} />
                <span className="w-2 h-2 bg-teal-400 rounded-full animate-bounce ml-1" style={{ animationDelay: '300ms' }} />
              </span>
              <span className="text-sm font-medium">Generating response...</span>
            </div>

            {/* Expandable raw content toggle */}
            {message.content && (
              <button
                onClick={() => setShowRawContent(!showRawContent)}
                className="flex items-center gap-1 mt-2 text-xs text-gray-400 hover:text-gray-600 transition-colors"
              >
                <ChevronIcon isOpen={showRawContent} className="text-gray-400" />
                <span>{showRawContent ? 'Hide' : 'Show'} raw output</span>
              </button>
            )}

            {/* Raw content (collapsed by default) */}
            {showRawContent && message.content && (
              <div className="mt-2 p-2 bg-gray-200 rounded text-xs font-mono text-gray-600 max-h-40 overflow-auto whitespace-pre-wrap">
                {message.content}
              </div>
            )}
          </div>
        ) : (
          /* Normal rendering for completed messages */
          <>
            <div className={`prose prose-sm max-w-none break-words ${
              isUser
                ? 'prose-invert text-white prose-p:text-white prose-headings:text-white prose-strong:text-white prose-code:text-white prose-li:text-white prose-ol:text-white prose-ul:text-white'
                : 'prose-gray'
            }`}>
              <ReactMarkdown
                components={{
                  // Custom code block rendering
                  code: ({ className, children, ...props }) => {
                    const match = /language-(\w+)/.exec(className || '');
                    const isInline = !match && !className;

                    if (isInline) {
                      // Get text content, handling various React children types
                      let textContent = '';
                      if (typeof children === 'string') {
                        textContent = children;
                      } else if (Array.isArray(children)) {
                        textContent = children.map(c => String(c)).join('');
                      } else {
                        textContent = String(children || '');
                      }
                      // Strip any backticks from inline code content
                      const cleanContent = textContent.replace(/`/g, '');

                      // If content has newlines, render as a block-level code element
                      const hasNewlines = cleanContent.includes('\n');
                      if (hasNewlines) {
                        return (
                          <code
                            className={`block ${isUser ? 'bg-teal-500/50' : 'bg-gray-200'} px-3 py-2 rounded text-sm font-mono whitespace-pre my-2`}
                            {...props}
                          >
                            {cleanContent}
                          </code>
                        );
                      }

                      return (
                        <code
                          className={`${isUser ? 'bg-teal-500/50' : 'bg-gray-200'} px-1.5 py-0.5 rounded text-sm`}
                          {...props}
                        >
                          {cleanContent}
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
                  // Style lists - use list-outside with padding for proper alignment
                  ul: ({ children }) => (
                    <ul className="list-disc pl-5 mb-2 space-y-1">{children}</ul>
                  ),
                  ol: ({ children }) => (
                    <ol className="list-decimal pl-5 mb-2 space-y-1">{children}</ol>
                  ),
                  // Style list items to ensure inline display of number and content
                  li: ({ children }) => (
                    <li className="pl-1">{children}</li>
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
                {displayContent}
              </ReactMarkdown>
            </div>

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
          </>
        )}
      </div>
    </div>
  );
}
