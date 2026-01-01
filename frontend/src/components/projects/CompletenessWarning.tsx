'use client';

import { CompletenessReport, CompletenessIssue } from '@/types';

interface CompletenessWarningProps {
  report: CompletenessReport;
  onFixClick?: () => void;
  onDismiss?: () => void;
}

function IssueIcon({ severity }: { severity: string }) {
  if (severity === 'critical') {
    return (
      <svg className="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
    );
  }
  return (
    <svg className="w-5 h-5 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
    </svg>
  );
}

function IssueItem({ issue }: { issue: CompletenessIssue }) {
  return (
    <li className="flex items-start gap-2 text-sm">
      <span className="text-gray-400 mt-0.5">
        {issue.referenceType === 'script' && '📜'}
        {issue.referenceType === 'stylesheet' && '🎨'}
        {issue.referenceType === 'import' && '📦'}
        {issue.referenceType === 'image' && '🖼️'}
        {!issue.referenceType && '📄'}
      </span>
      <div>
        <span className="font-medium text-gray-700">{issue.missingFile}</span>
        {issue.referencedBy && (
          <span className="text-gray-500"> (referenced by {issue.referencedBy})</span>
        )}
      </div>
    </li>
  );
}

export function CompletenessWarning({ report, onFixClick, onDismiss }: CompletenessWarningProps) {
  if (report.status === 'pass') {
    return null;
  }

  const isCritical = report.status === 'critical';
  const criticalIssues = report.issues.filter(i => i.severity === 'critical');
  const warningIssues = report.issues.filter(i => i.severity === 'warning');
  const missingFiles = Array.from(new Set(report.issues.filter(i => i.type === 'missing_file').map(i => i.missingFile)));

  return (
    <div className={`rounded-lg border p-4 mb-4 ${
      isCritical
        ? 'bg-red-50 border-red-200'
        : 'bg-yellow-50 border-yellow-200'
    }`}>
      <div className="flex items-start gap-3">
        <IssueIcon severity={isCritical ? 'critical' : 'warning'} />
        <div className="flex-1">
          <h3 className={`font-medium ${isCritical ? 'text-red-800' : 'text-yellow-800'}`}>
            {isCritical ? 'Some files are missing' : 'Minor issues detected'}
          </h3>
          <p className={`text-sm mt-1 ${isCritical ? 'text-red-600' : 'text-yellow-600'}`}>
            {isCritical
              ? 'Your app may not work correctly until these files are created.'
              : 'Your app should work, but some resources are missing.'}
          </p>

          {missingFiles.length > 0 && (
            <ul className="mt-3 space-y-1">
              {criticalIssues.slice(0, 5).map((issue) => (
                <IssueItem key={issue.id} issue={issue} />
              ))}
              {warningIssues.slice(0, 3).map((issue) => (
                <IssueItem key={issue.id} issue={issue} />
              ))}
              {report.issues.length > 8 && (
                <li className="text-sm text-gray-500">
                  ...and {report.issues.length - 8} more
                </li>
              )}
            </ul>
          )}

          <div className="mt-4 flex gap-2">
            {report.autoFixable > 0 && onFixClick && (
              <button
                onClick={onFixClick}
                className={`px-3 py-1.5 text-sm font-medium rounded-md ${
                  isCritical
                    ? 'bg-red-600 text-white hover:bg-red-700'
                    : 'bg-yellow-600 text-white hover:bg-yellow-700'
                }`}
              >
                Fix Now ({report.autoFixable} {report.autoFixable === 1 ? 'file' : 'files'})
              </button>
            )}
            {onDismiss && (
              <button
                onClick={onDismiss}
                className="px-3 py-1.5 text-sm font-medium text-gray-600 hover:text-gray-800"
              >
                {isCritical ? 'Show Preview Anyway' : 'Dismiss'}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default CompletenessWarning;
