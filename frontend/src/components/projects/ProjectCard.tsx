'use client';

import { useState, useRef, useEffect, KeyboardEvent } from 'react';
import { Project } from '@/types';

interface ProjectCardProps {
  project: Project;
  isActive?: boolean;
  onClick?: (project: Project) => void;
  onRename?: (project: Project, newTitle: string) => Promise<void>;
  onDelete?: (project: Project) => Promise<void>;
}

export function ProjectCard({ project, isActive = false, onClick, onRename, onDelete }: ProjectCardProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [editTitle, setEditTitle] = useState(project.title);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [isEditing]);

  // Reset states when project changes
  useEffect(() => {
    setEditTitle(project.title);
    setIsEditing(false);
    setShowDeleteConfirm(false);
  }, [project.id, project.title]);

  const handleClick = () => {
    if (!isEditing && !showDeleteConfirm) {
      onClick?.(project);
    }
  };

  const handleEditClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    setIsEditing(true);
  };

  const handleTrashClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    setShowDeleteConfirm(true);
  };

  const handleDeleteCancel = (e: React.MouseEvent) => {
    e.stopPropagation();
    setShowDeleteConfirm(false);
  };

  const handleDeleteConfirm = async (e: React.MouseEvent) => {
    e.stopPropagation();
    if (!onDelete) return;
    setIsDeleting(true);
    try {
      await onDelete(project);
    } finally {
      setIsDeleting(false);
      setShowDeleteConfirm(false);
    }
  };

  const handleSave = async () => {
    const trimmedTitle = editTitle.trim();
    if (trimmedTitle && trimmedTitle !== project.title) {
      await onRename?.(project, trimmedTitle);
    }
    setIsEditing(false);
  };

  const handleCancel = () => {
    setEditTitle(project.title);
    setIsEditing(false);
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleSave();
    } else if (e.key === 'Escape') {
      e.preventDefault();
      handleCancel();
    }
  };

  const handleSaveClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    handleSave();
  };

  const handleCancelClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    handleCancel();
  };

  const formattedDate = new Date(project.updatedAt).toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });

  return (
    <div className="relative overflow-hidden rounded-lg">
      {/* Main card content */}
      <div
        onClick={handleClick}
        className={`relative w-full text-left p-4 transition-all duration-200 cursor-pointer group ${
          isActive
            ? 'bg-teal-50 border-2 border-teal-400 rounded-lg'
            : 'bg-white border border-gray-200 hover:bg-gray-50 hover:border-gray-300 rounded-lg'
        }`}
      >
        {/* Title row */}
        <div className="flex items-center justify-between gap-2">
          {isEditing ? (
            /* Edit mode - inline input with icon buttons */
            <>
              <input
                ref={inputRef}
                type="text"
                value={editTitle}
                onChange={(e) => setEditTitle(e.target.value)}
                onKeyDown={handleKeyDown}
                onClick={(e) => e.stopPropagation()}
                className="flex-1 min-w-0 px-2 py-0.5 text-sm font-medium text-gray-900 border border-teal-400 rounded focus:outline-none focus:ring-1 focus:ring-teal-400"
                maxLength={100}
              />
              <div className="flex items-center gap-1 flex-shrink-0">
                <button
                  onClick={handleSaveClick}
                  className="p-1.5 rounded-md hover:bg-teal-100 transition-colors"
                  aria-label="Save"
                  title="Save (Enter)"
                >
                  <CheckIcon className="w-4 h-4 text-teal-600" />
                </button>
                <button
                  onClick={handleCancelClick}
                  className="p-1.5 rounded-md hover:bg-gray-200 transition-colors"
                  aria-label="Cancel"
                  title="Cancel (Esc)"
                >
                  <XIcon className="w-4 h-4 text-gray-500" />
                </button>
                {onDelete && (
                  <button
                    onClick={handleTrashClick}
                    className="p-1.5 rounded-md hover:bg-red-100 transition-colors"
                    aria-label="Delete project"
                    title="Delete"
                  >
                    <TrashIcon className="w-4 h-4 text-gray-400 hover:text-red-500" />
                  </button>
                )}
              </div>
            </>
          ) : (
            /* Normal view mode */
            <>
              <h3 className="flex-1 min-w-0 font-medium text-gray-900 truncate">{project.title}</h3>
              {onRename && (
                <div className="flex items-center gap-1 flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button
                    onClick={handleEditClick}
                    className="p-1.5 rounded-md hover:bg-gray-200 transition-colors"
                    aria-label="Edit project"
                  >
                    <PencilIcon className="w-4 h-4 text-gray-500" />
                  </button>
                </div>
              )}
            </>
          )}
        </div>

        {/* Date row */}
        <p className="text-sm text-gray-500 mt-1">{formattedDate}</p>

        {/* Delete confirmation overlay - same size as card */}
        <div
          className={`absolute inset-0 flex items-center justify-center px-4 bg-red-50 border border-red-200 rounded-lg transition-all duration-200 ease-out ${
            showDeleteConfirm
              ? 'opacity-100 translate-x-0'
              : 'opacity-0 translate-x-full pointer-events-none'
          }`}
        >
          <div className="flex items-center gap-3">
            <span className="text-sm text-red-700 font-medium">Delete?</span>
            <button
              onClick={handleDeleteCancel}
              className="px-3 py-1 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleDeleteConfirm}
              disabled={isDeleting}
              className="px-3 py-1 text-sm font-medium text-white bg-red-500 rounded-md hover:bg-red-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {isDeleting ? '...' : 'Delete'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

function PencilIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
      />
    </svg>
  );
}

function CheckIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
    </svg>
  );
}

function XIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
    </svg>
  );
}

function TrashIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
      />
    </svg>
  );
}
