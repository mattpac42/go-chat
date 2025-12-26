'use client';

// Types for Discovery Summary
interface DiscoveryUser {
  id: string;
  description: string;
  count: number;
  hasPermissions: boolean;
  permissionNotes?: string;
}

interface DiscoveryFeature {
  id: string;
  name: string;
  priority: number;
  version: string; // "v1" or "v2"
}

interface DiscoverySummary {
  projectName: string;
  solvesStatement: string;
  users: DiscoveryUser[];
  mvpFeatures: DiscoveryFeature[];
  futureFeatures: DiscoveryFeature[];
}

interface DiscoverySummaryCardProps {
  summary: DiscoverySummary;
  onEdit: () => void;
  onConfirm: () => void;
  isConfirming?: boolean;
}

function LoadingSpinner({ className }: { className?: string }) {
  return (
    <svg
      className={`animate-spin ${className}`}
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        className="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        strokeWidth="4"
      />
      <path
        className="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
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
        d="M14 5l7 7m0 0l-7 7m7-7H3"
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
 * DiscoverySummaryCard - Displays a discovery summary inline in chat
 *
 * Shows project details, problem statement, users, and features in a
 * responsive 2-column (desktop) or stacked (mobile) layout.
 */
export function DiscoverySummaryCard({
  summary,
  onEdit,
  onConfirm,
  isConfirming = false,
}: DiscoverySummaryCardProps) {
  const { projectName, solvesStatement, users, mvpFeatures, futureFeatures } = summary;

  const hasFutureFeatures = futureFeatures.length > 0;
  const hasUsers = users.length > 0;
  const hasMvpFeatures = mvpFeatures.length > 0;

  return (
    <div className="bg-white rounded-lg border border-gray-200 shadow-sm p-4 md:p-6">
      {/* Desktop: 2-column grid / Mobile: stacked */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6">
        {/* Left column on desktop */}
        <div className="space-y-4">
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
        </div>

        {/* Right column on desktop */}
        <div className="space-y-4">
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
        </div>
      </div>

      {/* COMING LATER Section - Full width */}
      {hasFutureFeatures && (
        <div className="mt-4 pt-4 border-t border-gray-100">
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

      {/* Action Buttons */}
      <div className="mt-6 pt-4 border-t border-gray-100 flex flex-col md:flex-row md:justify-between gap-3">
        {/* Start Over - Ghost button (resets discovery) */}
        <button
          onClick={onEdit}
          disabled={isConfirming}
          className="px-4 py-2 text-sm font-medium text-gray-700 bg-transparent border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:ring-offset-1 disabled:opacity-50 disabled:cursor-not-allowed transition-colors order-2 md:order-1"
          aria-label="Start over with discovery"
        >
          Start Over
        </button>

        {/* Start Building - Primary button */}
        <button
          onClick={onConfirm}
          disabled={isConfirming}
          className="px-4 py-2 text-sm font-medium text-white bg-teal-600 rounded-lg hover:bg-teal-700 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:ring-offset-1 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2 order-1 md:order-2"
          aria-label="Confirm and start building project"
        >
          {isConfirming ? (
            <>
              <LoadingSpinner className="w-4 h-4" />
              <span>Starting...</span>
            </>
          ) : (
            <>
              <span>Start Building</span>
              <ArrowRightIcon className="w-4 h-4" />
            </>
          )}
        </button>
      </div>
    </div>
  );
}

// Export types for consumers
export type { DiscoverySummary, DiscoveryUser, DiscoveryFeature, DiscoverySummaryCardProps };
