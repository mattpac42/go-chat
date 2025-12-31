'use client';

import ReactMarkdown from 'react-markdown';
import { Message, AGENT_CONFIG } from '@/types';
import { CodeBlock } from './CodeBlock';
import { AgentHeader } from './AgentHeader';
import { CollapsibleContent } from './CollapsibleContent';

interface MessageBubbleProps {
  message: Message;
  showCodeBlocks?: boolean; // If false, replace code blocks with friendly summary
  showBadge?: boolean; // Show "NEW" badge for first agent introduction
}

/**
 * Simple preprocessing for user messages - just clean up extra newlines
 */
function preprocessMarkdown(content: string): string {
  return content.replace(/\n{3,}/g, '\n\n');
}

/**
 * Strip DISCOVERY_DATA metadata comments from content
 * These are backend-to-backend metadata that should never be shown to users
 */
function stripDiscoveryMetadata(content: string): string {
  // Remove the discovery data HTML comment completely
  // Pattern: <!--DISCOVERY_DATA:{...}--> where {...} is JSON (may contain nested braces)
  // Use a non-greedy match that ends at the --> delimiter
  return content.replace(/<!--DISCOVERY_DATA:.*?-->/g, '');
}

/**
 * Convert inline lists to proper markdown block lists
 * Handles patterns like:
 * - "1. Item 2. Item 3. Item" -> proper numbered list
 * - "* Item * Item" -> proper bullet list (only asterisks, not dashes)
 *
 * IMPORTANT: Only converts truly inline lists (no newlines between items).
 * If the content already has proper newlines, it's left alone.
 */
function convertInlineLists(content: string): string {
  let result = content;

  // Skip if content already has proper list formatting with newlines
  // This check prevents mangling properly formatted markdown
  if (/\n\s*[\*\-]\s+/.test(result) || /\n\s*\d+\.\s+/.test(result)) {
    return result;
  }

  // Convert inline numbered lists: "1. Item 2. Item 3. Item" (all on one line)
  // Only match if there are multiple numbered items without newlines between them
  if (/\d+\.\s+[^\n]+\s+\d+\.\s+/.test(result)) {
    // Split by numbered pattern and rebuild as block list
    const parts = result.split(/(\d+\.\s+)/);
    let inList = false;
    let listItems: string[] = [];
    let beforeList = '';
    let afterList = '';

    for (let i = 0; i < parts.length; i++) {
      const part = parts[i];
      const nextPart = parts[i + 1];

      // Check if this is a number marker (e.g., "1. ", "2. ")
      if (/^\d+\.\s+$/.test(part) && nextPart !== undefined) {
        if (!inList) {
          // Starting a list - everything before is preamble
          beforeList = listItems.join('');
          listItems = [];
          inList = true;
        }
        // Get the content after this number (next part)
        const itemContent = nextPart.trim();
        if (itemContent) {
          listItems.push(part.trim() + ' ' + itemContent);
        }
        i++; // Skip the content part we just consumed
      } else if (!inList) {
        // Before any list
        listItems.push(part);
      } else {
        // After list items - could be trailing content
        afterList += part;
      }
    }

    if (listItems.length > 1) {
      // Rebuild with newlines between items
      result = beforeList.trim();
      if (result) result += '\n\n';
      result += listItems.join('\n');
      if (afterList.trim()) result += '\n\n' + afterList.trim();
    }
  }

  // Convert inline bullet lists: "* Item * Item" (asterisks only, not dashes)
  // Dashes are commonly used in prose (e.g., "handle - like") and shouldn't be converted
  // Only match asterisks that appear inline without newlines
  if (/\*\s+[^\n]+\s+\*\s+/.test(result)) {
    // Split by asterisk bullet pattern and rebuild as block list
    const parts = result.split(/\s+(\*)\s+/);
    if (parts.length > 2) {
      let newContent = parts[0].trim(); // Text before first bullet
      const items: string[] = [];

      for (let i = 1; i < parts.length; i += 2) {
        const marker = parts[i];
        const itemText = parts[i + 1]?.trim();
        if (marker && itemText) {
          items.push(`${marker} ${itemText}`);
        }
      }

      if (items.length > 1) {
        if (newContent) newContent += '\n\n';
        result = newContent + items.join('\n');
      }
    }
  }

  return result;
}

/**
 * Extract file names from code blocks and create a summary message
 * Also removes code blocks from display since files are shown in the Files panel
 */
function processAssistantContent(content: string, showCode: boolean): string {
  // First, strip any discovery metadata - this should NEVER be shown to users
  let cleanContent = stripDiscoveryMetadata(content);

  // Extract filenames from code blocks with our metadata format
  // Pattern: ```lang:filename followed by --- and short_description
  const filenamePattern = /```(\w+):([^\n]+)\n---[\s\S]*?short_description:/g;
  const filenameMatches = Array.from(cleanContent.matchAll(filenamePattern));
  const filenames = filenameMatches.map(m => m[2].trim());

  // Remove duplicates
  const uniqueFilenames = Array.from(new Set(filenames));

  // If we found files with metadata, remove those code blocks from display
  // Strategy: Find the FIRST code block with metadata and remove everything from there to the end
  // This is simpler and more reliable than trying to find matching closing fences
  // (which is hard when the content contains nested code blocks)
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

  // Apply inline list conversion to fix Claude's tendency to output inline lists
  cleanContent = convertInlineLists(cleanContent);

  // If showCode is true, return whatever content we have
  if (showCode) {
    return cleanContent;
  }

  // If no files found, just return clean content
  if (uniqueFilenames.length === 0) {
    return cleanContent;
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

  // Track if message is currently streaming
  const isStreaming = message.isStreaming && !isUser;

  // Process content to hide code blocks for assistant messages
  // Preprocess user content too to fix markdown issues
  const displayContent = isUser
    ? preprocessMarkdown(message.content)
    : processAssistantContent(message.content, showCodeBlocks);

  // Get agent config for styling (for all assistant messages)
  // Default to 'product_manager' (Root) for assistant messages without explicit agentType
  const effectiveAgentType = !isUser ? (message.agentType || 'product_manager') : null;
  const agentConfig = effectiveAgentType ? AGENT_CONFIG[effectiveAgentType] : null;

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
        className={`max-w-[85%] md:max-w-[70%] rounded-2xl px-4 py-3 min-w-0 overflow-hidden ${
          isUser
            ? 'bg-teal-400 text-white rounded-br-md'
            : 'bg-gray-100 text-gray-900 rounded-bl-md'
        }`}
        style={bubbleStyle}
      >
        {/* Agent header for all assistant messages */}
        {!isUser && effectiveAgentType && (
          <AgentHeader agentType={effectiveAgentType} showBadge={showBadge} />
        )}

        {/* Content - render normally whether streaming or complete */}
        {/* Disable progressive disclosure during streaming to prevent flickering */}
        <CollapsibleContent
          paragraphThreshold={3}
          visibleParagraphs={2}
          isUserMessage={isUser}
          disabled={isStreaming}
        >
          <div className={`prose prose-sm max-w-none break-words [&_pre]:whitespace-pre-wrap [&_pre]:break-words [&_pre]:max-w-full [&_code]:break-words [&_code]:whitespace-pre-wrap ${
            isUser
              ? 'prose-invert text-white prose-p:text-white prose-headings:text-white prose-strong:text-white prose-code:text-white prose-li:text-white prose-ol:text-white prose-ul:text-white [&_ol>li]:marker:text-white [&_ul>li]:marker:text-white'
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
            </CollapsibleContent>

        {/* Streaming indicator - show after content */}
        {isStreaming && (
          <div className="flex items-center gap-2 mt-2 text-gray-500">
            <span className="inline-flex">
              <span className="w-1.5 h-1.5 bg-teal-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
              <span className="w-1.5 h-1.5 bg-teal-400 rounded-full animate-bounce ml-1" style={{ animationDelay: '150ms' }} />
              <span className="w-1.5 h-1.5 bg-teal-400 rounded-full animate-bounce ml-1" style={{ animationDelay: '300ms' }} />
            </span>
          </div>
        )}

        {/* Timestamp - only show when not streaming */}
        {!isStreaming && formattedTime && (
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
