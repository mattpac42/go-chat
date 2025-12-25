import React from 'react';

interface ReactMarkdownProps {
  children: string;
  components?: Record<string, React.ComponentType<{ children?: React.ReactNode }>>;
}

function MockReactMarkdown({ children, components }: ReactMarkdownProps) {
  // Simple mock that renders content with basic parsing
  const content = children || '';

  // Parse code blocks
  const codeBlockRegex = /```(\w+)?\n?([\s\S]*?)```/g;
  let lastIndex = 0;
  const parts: React.ReactNode[] = [];
  let match;

  while ((match = codeBlockRegex.exec(content)) !== null) {
    // Add text before code block
    if (match.index > lastIndex) {
      const textContent = content.slice(lastIndex, match.index);
      parts.push(<span key={`text-${lastIndex}`}>{textContent}</span>);
    }

    // Add code block
    const language = match[1] || 'text';
    const code = match[2].trim();
    if (components?.code) {
      const CodeComponent = components.code;
      parts.push(
        <CodeComponent key={`code-${match.index}`} className={`language-${language}`}>
          {code}
        </CodeComponent>
      );
    } else {
      parts.push(
        <code key={`code-${match.index}`} className={`language-${language}`}>
          {code}
        </code>
      );
    }

    lastIndex = match.index + match[0].length;
  }

  // Add remaining text
  if (lastIndex < content.length) {
    const remaining = content.slice(lastIndex);
    parts.push(<span key={`text-${lastIndex}`}>{remaining}</span>);
  }

  if (parts.length === 0) {
    return <>{content}</>;
  }

  return <>{parts}</>;
}

export default MockReactMarkdown;
