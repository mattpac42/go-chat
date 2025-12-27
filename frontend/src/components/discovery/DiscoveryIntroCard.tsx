'use client';

interface DiscoveryIntroCardProps {
  onSkip?: () => void;
  showSkipOption?: boolean;
}

export function DiscoveryIntroCard({
  onSkip,
  showSkipOption = false,
}: DiscoveryIntroCardProps) {
  return (
    <div className="bg-gradient-to-br from-teal-50 to-cyan-50 border border-teal-200 rounded-xl p-6 mx-4 my-4 shadow-sm">
      {/* Header - icon and text vertically centered */}
      <div className="flex items-center gap-4">
        <div className="w-12 h-12 rounded-full bg-teal-500 flex items-center justify-center flex-shrink-0">
          <CompassIcon className="w-6 h-6 text-white" />
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="text-lg font-semibold text-teal-800 leading-tight">
            Let&apos;s figure out what you need
          </h3>
          <div className="flex items-center gap-1.5 text-sm text-teal-600 mt-0.5">
            <ClockIcon className="w-4 h-4 flex-shrink-0" />
            <span>About 5 minutes</span>
          </div>
        </div>
      </div>

      {/* Content - aligned with header text using pl matching icon width + gap */}
      <div className="mt-5 pl-16">
        <p className="text-sm text-gray-600 mb-3">
          Root will help you clarify:
        </p>
        <ul className="space-y-2.5">
          <li className="flex items-center gap-2.5 text-sm text-gray-700">
            <CheckCircleIcon className="w-4 h-4 text-teal-500 flex-shrink-0" />
            <span>The problem you&apos;re solving</span>
          </li>
          <li className="flex items-center gap-2.5 text-sm text-gray-700">
            <CheckCircleIcon className="w-4 h-4 text-teal-500 flex-shrink-0" />
            <span>Who will use your product</span>
          </li>
          <li className="flex items-center gap-2.5 text-sm text-gray-700">
            <CheckCircleIcon className="w-4 h-4 text-teal-500 flex-shrink-0" />
            <span>Essential features to start with</span>
          </li>
        </ul>
      </div>

      {/* Skip option for returning users */}
      {showSkipOption && onSkip && (
        <div className="mt-5 pt-4 border-t border-teal-200 pl-16">
          <button
            onClick={onSkip}
            className="text-sm text-teal-600 hover:text-teal-800 font-medium flex items-center gap-1.5 transition-colors"
          >
            <span>Done this before? Skip to building</span>
            <ArrowRightIcon className="w-4 h-4" />
          </button>
        </div>
      )}
    </div>
  );
}

function CompassIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
      />
    </svg>
  );
}

function ClockIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  );
}

function CheckCircleIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="currentColor"
      viewBox="0 0 20 20"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fillRule="evenodd"
        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
        clipRule="evenodd"
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
      xmlns="http://www.w3.org/2000/svg"
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
