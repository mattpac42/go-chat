'use client';

import { Project } from '@/types';
import { ProjectCard } from './ProjectCard';

interface ProjectListProps {
  projects: Project[];
  activeProjectId?: string;
  onNewProject?: () => void;
  onProjectSelect?: (project: Project) => void;
  onProjectRename?: (project: Project, newTitle: string) => Promise<void>;
  onProjectDelete?: (project: Project) => Promise<void>;
}

export function ProjectList({
  projects,
  activeProjectId,
  onNewProject,
  onProjectSelect,
  onProjectRename,
  onProjectDelete,
}: ProjectListProps) {
  return (
    <div className="flex flex-col h-full w-full">
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b border-gray-200">
        <h2 className="text-lg font-semibold text-gray-900">Projects</h2>
        <button
          onClick={onNewProject}
          className="flex items-center justify-center w-10 h-10 rounded-full bg-teal-400 text-white hover:bg-teal-500 focus:outline-none focus:ring-2 focus:ring-teal-400 focus:ring-offset-2 transition-colors"
          aria-label="New project"
        >
          <PlusIcon className="w-5 h-5" />
        </button>
      </div>

      {/* Project list */}
      <div className="flex-1 overflow-y-auto p-4 space-y-3">
        {projects.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <p className="mb-2">No projects yet</p>
            <p className="text-sm">Create your first project to get started</p>
          </div>
        ) : (
          projects.map((project) => (
            <ProjectCard
              key={project.id}
              project={project}
              isActive={project.id === activeProjectId}
              onClick={onProjectSelect}
              onRename={onProjectRename}
              onDelete={onProjectDelete}
            />
          ))
        )}
      </div>
    </div>
  );
}

function PlusIcon({ className }: { className?: string }) {
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
        d="M12 4v16m8-8H4"
      />
    </svg>
  );
}
