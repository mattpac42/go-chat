'use client';

import { useState, useEffect, useCallback } from 'react';
import { CompletenessReport } from '@/types';
import { CompletenessWarning } from '../projects/CompletenessWarning';

interface PreviewFile {
  path: string;
  content: string;
}

interface PreviewModalProps {
  files: PreviewFile[];
  isOpen: boolean;
  onClose: () => void;
  completenessReport?: CompletenessReport | null;
  onFixClick?: () => void;
}

type DeviceType = 'desktop' | 'tablet' | 'mobile';

const DEVICE_WIDTHS: Record<DeviceType, string> = {
  desktop: '100%',
  tablet: '768px',
  mobile: '375px',
};

function CloseIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M6 18L18 6M6 6l12 12"
      />
    </svg>
  );
}

function RefreshIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
      />
    </svg>
  );
}

function DesktopIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
      />
    </svg>
  );
}

function TabletIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M12 18h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"
      />
    </svg>
  );
}

function MobileIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"
      />
    </svg>
  );
}

/**
 * Combines HTML, CSS, and JS files into a single previewable document.
 */
function buildPreviewHtml(files: PreviewFile[]): string {
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

  const cssFiles = files.filter(f => f.path.endsWith('.css'));
  const jsFiles = files.filter(f => f.path.endsWith('.js'));

  let html = mainHtml.content;

  if (cssFiles.length > 0) {
    const cssContent = cssFiles
      .map(f => `<!-- ${f.path} -->\n<style>\n${f.content}\n</style>`)
      .join('\n');

    if (html.includes('</head>')) {
      html = html.replace('</head>', `${cssContent}\n</head>`);
    } else if (html.includes('<body')) {
      html = html.replace(/<body[^>]*>/, (match) => `${cssContent}\n${match}`);
    } else {
      html = `${cssContent}\n${html}`;
    }
  }

  if (jsFiles.length > 0) {
    const jsContent = jsFiles
      .map(f => `<!-- ${f.path} -->\n<script>\n${f.content}\n</script>`)
      .join('\n');

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
 * PreviewModal - Fullscreen preview modal with device frame selector
 *
 * Features:
 * - Centered overlay taking 90% of viewport (max-width: 1200px)
 * - Device frame selector: Desktop / Tablet / Mobile
 * - Refresh button to reload the preview
 * - Close on Escape key or backdrop click
 */
export function PreviewModal({ files, isOpen, onClose, completenessReport, onFixClick }: PreviewModalProps) {
  const [device, setDevice] = useState<DeviceType>('desktop');
  const [refreshKey, setRefreshKey] = useState(0);
  const [showPreviewAnyway, setShowPreviewAnyway] = useState(false);

  // Reset "show anyway" state when modal opens
  useEffect(() => {
    if (isOpen) {
      setShowPreviewAnyway(false);
    }
  }, [isOpen]);

  // Handle escape key
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };
    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
      document.body.style.overflow = 'hidden';
    }
    return () => {
      document.removeEventListener('keydown', handleEscape);
      document.body.style.overflow = '';
    };
  }, [isOpen, onClose]);

  const handleRefresh = useCallback(() => {
    setRefreshKey(prev => prev + 1);
  }, []);

  const handleBackdropClick = useCallback((e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  }, [onClose]);

  if (!isOpen) return null;

  const previewHtml = buildPreviewHtml(files);
  const deviceWidth = DEVICE_WIDTHS[device];
  const isFullWidth = device === 'desktop';

  // Determine if we should show the warning banner
  const hasIssues = completenessReport && completenessReport.status !== 'pass';
  const shouldBlockPreview = hasIssues && completenessReport.status === 'critical' && !showPreviewAnyway;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6"
      onClick={handleBackdropClick}
    >
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/60 transition-opacity" />

      {/* Modal */}
      <div
        className="relative bg-white rounded-xl shadow-2xl flex flex-col w-full max-h-[90vh]"
        style={{ maxWidth: '1200px', width: '90%', height: '90vh' }}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-gray-200 bg-gray-50 rounded-t-xl">
          <h2 className="text-lg font-semibold text-gray-900">Preview</h2>

          {/* Device selector */}
          <div className="flex items-center gap-1 bg-white rounded-lg border border-gray-200 p-1">
            <button
              onClick={() => setDevice('desktop')}
              className={`flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
                device === 'desktop'
                  ? 'bg-teal-100 text-teal-700'
                  : 'text-gray-500 hover:bg-gray-100 hover:text-gray-700'
              }`}
              title="Desktop view"
              aria-label="Desktop view"
            >
              <DesktopIcon className="w-4 h-4" />
              <span className="hidden sm:inline">Desktop</span>
            </button>
            <button
              onClick={() => setDevice('tablet')}
              className={`flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
                device === 'tablet'
                  ? 'bg-teal-100 text-teal-700'
                  : 'text-gray-500 hover:bg-gray-100 hover:text-gray-700'
              }`}
              title="Tablet view (768px)"
              aria-label="Tablet view"
            >
              <TabletIcon className="w-4 h-4" />
              <span className="hidden sm:inline">Tablet</span>
            </button>
            <button
              onClick={() => setDevice('mobile')}
              className={`flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
                device === 'mobile'
                  ? 'bg-teal-100 text-teal-700'
                  : 'text-gray-500 hover:bg-gray-100 hover:text-gray-700'
              }`}
              title="Mobile view (375px)"
              aria-label="Mobile view"
            >
              <MobileIcon className="w-4 h-4" />
              <span className="hidden sm:inline">Mobile</span>
            </button>
          </div>

          {/* Action buttons */}
          <div className="flex items-center gap-2">
            <button
              onClick={handleRefresh}
              className="p-2 rounded-lg hover:bg-gray-200 transition-colors"
              title="Refresh preview"
              aria-label="Refresh preview"
            >
              <RefreshIcon className="w-5 h-5 text-gray-500" />
            </button>
            <button
              onClick={onClose}
              className="p-2 rounded-lg hover:bg-gray-200 transition-colors"
              title="Close preview"
              aria-label="Close preview"
            >
              <CloseIcon className="w-5 h-5 text-gray-500" />
            </button>
          </div>
        </div>

        {/* Completeness warning */}
        {hasIssues && (
          <div className="px-4 pt-4">
            <CompletenessWarning
              report={completenessReport}
              onFixClick={onFixClick}
              onDismiss={() => setShowPreviewAnyway(true)}
            />
          </div>
        )}

        {/* Preview content */}
        <div className={`flex-1 overflow-hidden bg-gray-100 flex items-start justify-center p-4 rounded-b-xl ${shouldBlockPreview ? 'opacity-50 pointer-events-none' : ''}`}>
          {/* Device frame wrapper */}
          <div
            className={`bg-white h-full transition-all duration-300 ease-out ${
              !isFullWidth ? 'rounded-lg shadow-lg border border-gray-300' : ''
            }`}
            style={{
              width: deviceWidth,
              maxWidth: '100%',
            }}
          >
            {/* Device bezel for mobile/tablet */}
            {!isFullWidth && (
              <div className="bg-gray-800 rounded-t-lg px-4 py-2 flex items-center justify-center">
                <div className="w-16 h-1 bg-gray-600 rounded-full" />
              </div>
            )}

            <iframe
              key={refreshKey}
              srcDoc={previewHtml}
              sandbox="allow-scripts allow-same-origin"
              className={`w-full bg-white border-0 ${
                !isFullWidth
                  ? 'rounded-b-lg'
                  : 'rounded-lg'
              }`}
              style={{
                height: !isFullWidth ? 'calc(100% - 32px)' : '100%',
              }}
              title="App Preview"
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default PreviewModal;
