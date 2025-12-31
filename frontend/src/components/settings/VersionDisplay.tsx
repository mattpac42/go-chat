'use client';

/**
 * VersionDisplay - Shows app version and git hash
 *
 * Format: v{version} ({hash})
 * Example: v0.1.0 (fb50661)
 */
export function VersionDisplay() {
  const version = process.env.NEXT_PUBLIC_APP_VERSION || '0.0.0';
  const gitHash = process.env.NEXT_PUBLIC_GIT_HASH || 'dev';

  return (
    <span
      data-testid="version-display"
      className="text-xs text-gray-400"
    >
      v{version} ({gitHash})
    </span>
  );
}
