'use client';

import { useState } from 'react';
import { DiscoveryProgress, DiscoverySummaryCard, DiscoveryStageDrawer } from '@/components/discovery';
import type { DiscoveryStage } from '@/components/discovery';

// Mock data for demo
const mockSummary = {
  projectName: 'Cake Order Manager',
  solvesStatement: 'Replaces paper and WhatsApp order tracking with a simple digital system that keeps everything in one place.',
  users: [
    { id: '1', description: 'owner/baker', count: 1, hasPermissions: true, permissionNotes: 'full access' },
    { id: '2', description: 'order takers', count: 2, hasPermissions: false, permissionNotes: 'orders only' },
  ],
  mvpFeatures: [
    { id: '1', name: 'Order list view', priority: 1, version: 'v1' },
    { id: '2', name: 'Order creation form', priority: 2, version: 'v1' },
    { id: '3', name: 'Due date tracking', priority: 3, version: 'v1' },
  ],
  futureFeatures: [
    { id: '4', name: 'Calendar view', priority: 1, version: 'v2' },
    { id: '5', name: 'Customer database', priority: 2, version: 'v2' },
  ],
};

function getPlaceholder(stage: DiscoveryStage): string {
  switch (stage) {
    case 'welcome':
      return 'Tell me about yourself...';
    case 'problem':
      return 'What challenges do you face?';
    case 'personas':
      return 'Who will use this?';
    case 'mvp':
      return 'What features are essential?';
    case 'summary':
      return 'Ready to start building?';
    case 'complete':
      return 'Describe what you want to build...';
    default:
      return 'Type a message...';
  }
}

export function DiscoveryDemoClient() {
  const [currentStage, setCurrentStage] = useState<DiscoveryStage>('personas');
  const [showDrawer, setShowDrawer] = useState(false);
  const [isConfirming, setIsConfirming] = useState(false);

  const stages: DiscoveryStage[] = ['welcome', 'problem', 'personas', 'mvp', 'summary', 'complete'];

  const handleConfirm = async () => {
    setIsConfirming(true);
    await new Promise(resolve => setTimeout(resolve, 1500));
    setIsConfirming(false);
    setCurrentStage('complete');
  };

  return (
    <div className="min-h-screen bg-gray-100 p-3 md:p-6">
      <div className="max-w-4xl mx-auto space-y-4">
        {/* Header with inline stage selector */}
        <div className="flex flex-wrap items-center justify-between gap-2 bg-white rounded-lg p-3 shadow-sm">
          <h1 className="text-lg font-bold text-gray-900">Discovery Demo</h1>
          <div className="flex flex-wrap gap-1.5">
            {stages.map((stage) => (
              <button
                key={stage}
                onClick={() => setCurrentStage(stage)}
                className={`px-2 py-1 rounded text-xs font-medium transition-colors ${
                  currentStage === stage
                    ? 'bg-teal-600 text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {stage}
              </button>
            ))}
          </div>
        </div>

        {/* Progress Indicators - Side by side */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
          {/* Desktop */}
          <div className="bg-white rounded-lg p-3 shadow-sm">
            <h2 className="text-xs font-semibold text-gray-400 uppercase mb-2">Desktop (≥768px)</h2>
            <div className="flex items-center justify-between border border-gray-200 rounded p-2">
              <span className="text-sm font-medium text-gray-900">New Project</span>
              <DiscoveryProgress
                currentStage={currentStage}
                onStageClick={() => setShowDrawer(true)}
                isMobile={false}
              />
            </div>
          </div>
          {/* Mobile */}
          <div className="bg-white rounded-lg p-3 shadow-sm">
            <h2 className="text-xs font-semibold text-gray-400 uppercase mb-2">Mobile (&lt;768px)</h2>
            <div className="flex items-center justify-between border border-gray-200 rounded p-2">
              <span className="text-sm font-medium text-gray-900">New Project</span>
              <DiscoveryProgress
                currentStage={currentStage}
                onStageClick={() => setShowDrawer(true)}
                isMobile={true}
              />
            </div>
          </div>
        </div>

        {/* Input Placeholder */}
        <div className="bg-white rounded-lg p-3 shadow-sm">
          <h2 className="text-xs font-semibold text-gray-400 uppercase mb-2">Stage-Aware Placeholder</h2>
          <input
            type="text"
            readOnly
            placeholder={getPlaceholder(currentStage)}
            className="w-full px-3 py-2 border border-gray-200 rounded text-sm text-gray-500 bg-gray-50"
          />
        </div>

        {/* Summary Card */}
        <div className="bg-white rounded-lg p-3 shadow-sm">
          <h2 className="text-xs font-semibold text-gray-400 uppercase mb-2">Summary Card</h2>
          <DiscoverySummaryCard
            summary={mockSummary}
            onEdit={() => alert('Edit clicked - would open edit mode')}
            onConfirm={handleConfirm}
            isConfirming={isConfirming}
          />
        </div>

        {/* Stage Drawer */}
        <DiscoveryStageDrawer
          isOpen={showDrawer}
          onClose={() => setShowDrawer(false)}
          currentStage={currentStage}
        />
      </div>
    </div>
  );
}
