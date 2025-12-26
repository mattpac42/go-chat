'use client';

import { useEffect, useCallback, useState } from 'react';
import { createPortal } from 'react-dom';

type DiscoveryStage = 'welcome' | 'problem' | 'personas' | 'mvp' | 'summary' | 'complete';

interface StageInfo {
  stage: DiscoveryStage;
  label: string;
  description: string;
  status: 'completed' | 'current' | 'upcoming';
}

interface DiscoveryStageDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  currentStage: DiscoveryStage;
}

const STAGE_INFO: Omit<StageInfo, 'status'>[] = [
  { stage: 'welcome', label: 'Welcome', description: 'Set the stage' },
  { stage: 'problem', label: 'Problem Discovery', description: 'Identify pain points' },
  { stage: 'personas', label: 'User Personas', description: 'Define who uses this' },
  { stage: 'mvp', label: 'MVP Scope', description: 'Essential features' },
  { stage: 'summary', label: 'Summary', description: 'Confirm and begin' },
];

const STAGE_ORDER: DiscoveryStage[] = ['welcome', 'problem', 'personas', 'mvp', 'summary', 'complete'];

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

function getStageStatus(stage: DiscoveryStage, currentStage: DiscoveryStage): 'completed' | 'current' | 'upcoming' {
  const stageIndex = STAGE_ORDER.indexOf(stage);
  const currentIndex = STAGE_ORDER.indexOf(currentStage);

  if (stageIndex < currentIndex) {
    return 'completed';
  } else if (stageIndex === currentIndex) {
    return 'current';
  } else {
    return 'upcoming';
  }
}

function StageIndicator({ status }: { status: 'completed' | 'current' | 'upcoming' }) {
  if (status === 'completed') {
    return (
      <div className="w-4 h-4 rounded-full bg-teal-500 flex-shrink-0" />
    );
  }

  if (status === 'current') {
    return (
      <div className="w-4 h-4 rounded-full bg-teal-500 ring-2 ring-teal-300 ring-offset-2 flex-shrink-0" />
    );
  }

  return (
    <div className="w-4 h-4 rounded-full border-2 border-gray-300 flex-shrink-0" />
  );
}

function StatusBadge({ status }: { status: 'completed' | 'current' | 'upcoming' }) {
  if (status === 'completed') {
    return (
      <span className="text-xs font-medium uppercase text-teal-600">
        Completed
      </span>
    );
  }

  if (status === 'current') {
    return (
      <span className="text-xs font-medium uppercase text-teal-600 bg-teal-50 px-2 py-0.5 rounded">
        Current
      </span>
    );
  }

  return (
    <span className="text-xs font-medium uppercase text-gray-400">
      Upcoming
    </span>
  );
}

function StageRow({ stage, label, description, status }: StageInfo) {
  return (
    <div className="flex items-center gap-4 py-3 px-4 border-b border-gray-100 last:border-b-0">
      <StageIndicator status={status} />
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-gray-800">{label}</p>
        <p className="text-xs text-gray-500">{description}</p>
      </div>
      <StatusBadge status={status} />
    </div>
  );
}

export function DiscoveryStageDrawer({ isOpen, onClose, currentStage }: DiscoveryStageDrawerProps) {
  const [mounted, setMounted] = useState(false);
  const [isVisible, setIsVisible] = useState(false);

  // Handle mounting for portal
  useEffect(() => {
    setMounted(true);
  }, []);

  // Handle visibility animation
  useEffect(() => {
    if (isOpen) {
      // Small delay to trigger CSS transition
      const timer = setTimeout(() => setIsVisible(true), 10);
      return () => clearTimeout(timer);
    } else {
      setIsVisible(false);
    }
  }, [isOpen]);

  // Handle escape key
  const handleEscape = useCallback((e: KeyboardEvent) => {
    if (e.key === 'Escape') {
      onClose();
    }
  }, [onClose]);

  useEffect(() => {
    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
      document.body.style.overflow = 'hidden';
    }
    return () => {
      document.removeEventListener('keydown', handleEscape);
      document.body.style.overflow = '';
    };
  }, [isOpen, handleEscape]);

  // Build stage info with current status
  const stages: StageInfo[] = STAGE_INFO.map((info) => ({
    ...info,
    status: getStageStatus(info.stage, currentStage),
  }));

  if (!mounted || !isOpen) return null;

  const drawerContent = (
    <div className="fixed inset-0 z-50">
      {/* Backdrop */}
      <div
        className={`absolute inset-0 bg-black/50 transition-opacity duration-300 ${
          isVisible ? 'opacity-100' : 'opacity-0'
        }`}
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Bottom sheet drawer */}
      <div
        className={`absolute bottom-0 left-0 right-0 bg-white rounded-t-2xl shadow-xl transition-transform duration-300 ease-out ${
          isVisible ? 'translate-y-0' : 'translate-y-full'
        }`}
        style={{ maxHeight: '70vh' }}
        role="dialog"
        aria-modal="true"
        aria-labelledby="discovery-drawer-title"
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-200">
          <h2 id="discovery-drawer-title" className="text-lg font-semibold text-gray-900">
            Discovery Progress
          </h2>
          <button
            onClick={onClose}
            className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
            aria-label="Close drawer"
          >
            <CloseIcon className="w-5 h-5 text-gray-500" />
          </button>
        </div>

        {/* Stage list */}
        <div className="overflow-y-auto" style={{ maxHeight: 'calc(70vh - 65px)' }}>
          {stages.map((stage) => (
            <StageRow key={stage.stage} {...stage} />
          ))}
        </div>
      </div>
    </div>
  );

  return createPortal(drawerContent, document.body);
}
