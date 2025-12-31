'use client';

import { useMemo } from 'react';

interface PreviewFile {
  path: string;
  content: string;
}

interface ProjectPreviewProps {
  files: PreviewFile[];
  className?: string;
}

/**
 * Combines HTML, CSS, and JS files into a single previewable document.
 * Finds the main HTML file and injects CSS/JS content.
 */
function buildPreviewHtml(files: PreviewFile[]): string {
  // Find the main HTML file (prefer index.html, then any .html)
  const htmlFiles = files.filter(f => f.path.endsWith('.html'));
  const mainHtml = htmlFiles.find(f => f.path.endsWith('index.html')) || htmlFiles[0];

  if (!mainHtml) {
    return `
<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
      margin: 0;
      background: #f5f5f5;
      color: #666;
    }
    .empty-state {
      text-align: center;
      padding: 2rem;
    }
    .empty-state svg {
      width: 48px;
      height: 48px;
      margin-bottom: 1rem;
      opacity: 0.5;
    }
    h2 { margin: 0 0 0.5rem; color: #333; font-weight: 500; }
    p { margin: 0; font-size: 0.875rem; }
  </style>
</head>
<body>
  <div class="empty-state">
    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
    </svg>
    <h2>No HTML file to preview</h2>
    <p>Create an HTML file to see your app here</p>
  </div>
</body>
</html>`;
  }

  // Get CSS and JS files
  const cssFiles = files.filter(f => f.path.endsWith('.css'));
  const jsFiles = files.filter(f => f.path.endsWith('.js'));

  let html = mainHtml.content;

  // Inject CSS files into <head>
  if (cssFiles.length > 0) {
    const cssContent = cssFiles
      .map(f => `<!-- ${f.path} -->\n<style>\n${f.content}\n</style>`)
      .join('\n');

    // Try to inject before </head>, or before <body>, or at the start
    if (html.includes('</head>')) {
      html = html.replace('</head>', `${cssContent}\n</head>`);
    } else if (html.includes('<body')) {
      html = html.replace(/<body[^>]*>/, (match) => `${cssContent}\n${match}`);
    } else {
      html = `${cssContent}\n${html}`;
    }
  }

  // Inject JS files before </body>
  if (jsFiles.length > 0) {
    const jsContent = jsFiles
      .map(f => `<!-- ${f.path} -->\n<script>\n${f.content}\n</script>`)
      .join('\n');

    // Try to inject before </body>, or at the end
    if (html.includes('</body>')) {
      html = html.replace('</body>', `${jsContent}\n</body>`);
    } else if (html.includes('</html>')) {
      html = html.replace('</html>', `${jsContent}\n</html>`);
    } else {
      html = `${html}\n${jsContent}`;
    }
  }

  return html;
}

/**
 * ProjectPreview - Renders project files in a sandboxed iframe
 *
 * Security: Uses sandbox with limited permissions:
 * - allow-scripts: Execute JavaScript
 * - allow-same-origin: Access localStorage/sessionStorage (required for stateful apps)
 *
 * Still prevents:
 * - Form submissions to external servers
 * - Popups and new windows
 * - Top navigation
 * - Pointer lock
 * - Presentation API
 */
export function ProjectPreview({ files, className = '' }: ProjectPreviewProps) {
  const previewHtml = useMemo(() => buildPreviewHtml(files), [files]);

  return (
    <iframe
      srcDoc={previewHtml}
      sandbox="allow-scripts allow-same-origin"
      className={`w-full h-full border-0 bg-white ${className}`}
      title="App Preview"
    />
  );
}

export default ProjectPreview;
