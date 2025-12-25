'use client';

import { useEffect } from 'react';

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="flex h-full items-center justify-center">
      <div className="text-center p-8">
        <h2 className="text-xl font-semibold text-gray-900 mb-2">
          Something went wrong
        </h2>
        <p className="text-gray-500 mb-4">
          An unexpected error has occurred.
        </p>
        <button
          onClick={reset}
          className="px-6 py-3 bg-teal-400 text-white rounded-lg hover:bg-teal-500 transition-colors"
        >
          Try again
        </button>
      </div>
    </div>
  );
}
