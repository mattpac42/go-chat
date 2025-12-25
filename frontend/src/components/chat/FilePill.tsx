'use client';

interface FilePillProps {
  filename: string;
  onClick?: () => void;
}

function FileIcon({ className }: { className?: string }) {
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
        d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
      />
    </svg>
  );
}

export function FilePill({ filename, onClick }: FilePillProps) {
  return (
    <button
      onClick={onClick}
      className="inline-flex items-center gap-1.5 px-2 py-1 bg-gray-100 hover:bg-gray-200 rounded-md text-sm text-gray-700 font-medium transition-colors cursor-pointer border border-gray-200"
    >
      <FileIcon className="w-3.5 h-3.5 text-gray-500" />
      <span className="truncate max-w-[200px]">{filename}</span>
    </button>
  );
}
