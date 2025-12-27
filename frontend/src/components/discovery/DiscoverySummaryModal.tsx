'use client';

import { useEffect, useRef } from 'react';
import { DiscoverySummary } from './DiscoverySummaryCard';

interface DiscoverySummaryModalProps {
  isOpen: boolean;
  onClose: () => void;
  summary: DiscoverySummary;
  /** If provided, shows action buttons for pre-confirmation state */
  onConfirm?: () => void;
  onEdit?: () => void;
  onStartOver?: () => void;
  isConfirming?: boolean;
}

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

function UserDot({ hasPermissions }: { hasPermissions: boolean }) {
  return (
    <span
      className={`inline-block w-2 h-2 rounded-full mr-2 flex-shrink-0 ${
        hasPermissions ? 'bg-teal-500' : 'bg-gray-400'
      }`}
      aria-hidden="true"
    />
  );
}

function SectionHeader({ children }: { children: React.ReactNode }) {
  return (
    <h4 className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">
      {children}
    </h4>
  );
}

/**
 * DiscoverySummaryModal - Modal to view discovery summary after completion
 *
 * Displays the project summary in a modal overlay, accessible from the header
 * after the user clicks "Start Building".
 */
function PencilIcon({ className }: { className?: string }) {
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
        d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
      />
    </svg>
  );
}

function ArrowRightIcon({ className }: { className?: string }) {
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
        d="M13 7l5 5m0 0l-5 5m5-5H6"
      />
    </svg>
  );
}

export function DiscoverySummaryModal({
  isOpen,
  onClose,
  summary,
  onConfirm,
  onEdit,
  onStartOver,
  isConfirming = false,
}: DiscoverySummaryModalProps) {
  // Determine if we're in action mode (pre-confirmation) or view-only mode
  const isActionMode = !!(onConfirm && onEdit && onStartOver);
  const modalRef = useRef<HTMLDivElement>(null);

  // Close on escape key
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape' && isOpen) {
        onClose();
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [isOpen, onClose]);

  // Focus trap
  useEffect(() => {
    if (isOpen && modalRef.current) {
      modalRef.current.focus();
    }
  }, [isOpen]);

  // Prevent body scroll when modal is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  if (!isOpen) return null;

  const { projectName, solvesStatement, users, mvpFeatures, futureFeatures } = summary;
  const hasFutureFeatures = futureFeatures.length > 0;
  const hasUsers = users.length > 0;
  const hasMvpFeatures = mvpFeatures.length > 0;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4"
      role="dialog"
      aria-modal="true"
      aria-labelledby="summary-modal-title"
    >
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/50 backdrop-blur-sm"
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Modal content */}
      <div
        ref={modalRef}
        className="relative bg-white rounded-xl shadow-xl max-w-lg w-full max-h-[85vh] overflow-y-auto"
        tabIndex={-1}
      >
        {/* Header */}
        <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between rounded-t-xl">
          <h2 id="summary-modal-title" className="text-lg font-semibold text-gray-900">
            Project Summary
          </h2>
          <button
            onClick={onClose}
            className="p-2 -mr-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors"
            aria-label="Close modal"
          >
            <CloseIcon className="w-5 h-5" />
          </button>
        </div>

        {/* Body */}
        <div className="px-6 py-4 space-y-5">
          {/* PROJECT Section */}
          <div>
            <SectionHeader>Project</SectionHeader>
            {projectName ? (
              <p className="text-sm text-gray-800 font-semibold">{projectName}</p>
            ) : (
              <p className="text-sm text-gray-400 italic">Project name will be generated</p>
            )}
          </div>

          {/* SOLVES Section */}
          <div>
            <SectionHeader>Solves</SectionHeader>
            {solvesStatement ? (
              <p className="text-sm text-gray-800">{solvesStatement}</p>
            ) : (
              <p className="text-sm text-gray-400 italic">Problem statement from discovery</p>
            )}
          </div>

          {/* USERS Section */}
          <div>
            <SectionHeader>Users</SectionHeader>
            {hasUsers ? (
              <ul className="space-y-1">
                {users.map((user) => (
                  <li key={user.id} className="flex items-start text-sm text-gray-800">
                    <UserDot hasPermissions={user.hasPermissions} />
                    <span>
                      {user.count === 1 ? 'You' : `${user.count} ${user.description}`}
                      {user.count === 1 && ` - ${user.description}`}
                      {user.permissionNotes && (
                        <span className="text-gray-500"> ({user.permissionNotes})</span>
                      )}
                    </span>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-sm text-gray-500 italic">No users defined</p>
            )}
          </div>

          {/* MVP FEATURES Section */}
          <div>
            <SectionHeader>MVP Features</SectionHeader>
            {hasMvpFeatures ? (
              <ul className="space-y-1">
                {mvpFeatures.map((feature) => (
                  <li key={feature.id} className="flex items-start text-sm text-gray-800">
                    <span className="mr-2 text-gray-400" aria-hidden="true">
                      &bull;
                    </span>
                    <span>{feature.name}</span>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-sm text-gray-500 italic">No MVP features defined</p>
            )}
          </div>

          {/* COMING LATER Section */}
          {hasFutureFeatures && (
            <div>
              <SectionHeader>Coming Later</SectionHeader>
              <ul className="space-y-1">
                {futureFeatures.map((feature) => (
                  <li key={feature.id} className="flex items-start text-sm text-gray-600">
                    <span className="mr-2 text-gray-400" aria-hidden="true">
                      &bull;
                    </span>
                    <span>
                      {feature.name}
                      <span className="ml-1 text-xs text-gray-400 uppercase">
                        ({feature.version})
                      </span>
                    </span>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="sticky bottom-0 bg-white border-t border-gray-200 px-6 py-4 rounded-b-xl">
          {isActionMode ? (
            <div className="flex flex-col gap-3">
              {/* Primary action */}
              <button
                onClick={onConfirm}
                disabled={isConfirming}
                className="w-full px-4 py-2.5 text-sm font-medium text-white bg-teal-500 rounded-lg hover:bg-teal-600 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:ring-offset-1 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
              >
                {isConfirming ? (
                  <>
                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    <span>Starting...</span>
                  </>
                ) : (
                  <>
                    <span>Start Building</span>
                    <ArrowRightIcon className="w-4 h-4" />
                  </>
                )}
              </button>

              {/* Secondary actions */}
              <div className="flex gap-2">
                <button
                  onClick={onEdit}
                  disabled={isConfirming}
                  className="flex-1 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-400 focus:ring-offset-1 transition-colors disabled:opacity-50 flex items-center justify-center gap-1.5"
                >
                  <PencilIcon className="w-4 h-4" />
                  <span>Edit in Chat</span>
                </button>
                <button
                  onClick={onStartOver}
                  disabled={isConfirming}
                  className="px-4 py-2 text-sm font-medium text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-400 focus:ring-offset-1 transition-colors disabled:opacity-50"
                >
                  Start Over
                </button>
              </div>
            </div>
          ) : (
            <button
              onClick={onClose}
              className="w-full px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:ring-offset-1 transition-colors"
            >
              Close
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
