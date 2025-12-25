'use client';

import { useState, useCallback, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { ProjectList } from '@/components/projects/ProjectList';
import { ChatContainer } from '@/components/chat/ChatContainer';
import { LoadingSpinner } from '@/components/shared/LoadingSpinner';
import { useProjects, useProject } from '@/hooks/useProjects';
import { Project } from '@/types';

export function HomeClient() {
  const router = useRouter();
  const { projects, isLoading, error, createProject, renameProject, deleteProject } = useProjects();
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  // Load project with messages when selected
  const {
    messages: projectMessages,
    isLoading: isLoadingProject,
  } = useProject(selectedProject?.id || '');

  // Select first project when projects load
  useEffect(() => {
    if (projects.length > 0 && !selectedProject) {
      setSelectedProject(projects[0]);
    }
  }, [projects, selectedProject]);

  // Update selectedProject reference when projects change
  useEffect(() => {
    if (selectedProject) {
      const updated = projects.find(p => p.id === selectedProject.id);
      if (updated && updated !== selectedProject) {
        setSelectedProject(updated);
      }
    }
  }, [projects, selectedProject]);

  const handleNewProject = useCallback(async () => {
    const newProject = await createProject();
    if (newProject) {
      router.push(`/projects/${newProject.id}`);
      setIsSidebarOpen(false);
    }
  }, [createProject, router]);

  const handleProjectSelect = useCallback((project: Project) => {
    router.push(`/projects/${project.id}`);
    setIsSidebarOpen(false);
  }, [router]);

  const handleMenuClick = useCallback(() => {
    setIsSidebarOpen(true);
  }, []);

  const handleCloseSidebar = useCallback(() => {
    setIsSidebarOpen(false);
  }, []);

  const handleProjectRename = useCallback(async (project: Project, newTitle: string) => {
    await renameProject(project.id, newTitle);
  }, [renameProject]);

  const handleProjectDelete = useCallback(async (project: Project) => {
    const success = await deleteProject(project.id);
    if (success && selectedProject?.id === project.id) {
      // If deleted project was selected, clear selection or select another
      const remaining = projects.filter(p => p.id !== project.id);
      setSelectedProject(remaining.length > 0 ? remaining[0] : null);
    }
  }, [deleteProject, selectedProject, projects]);

  // Show loading state
  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  // Show error state
  if (error) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center p-8">
          <h2 className="text-xl font-semibold text-gray-900 mb-2">
            Failed to load projects
          </h2>
          <p className="text-gray-500 mb-4">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="px-6 py-3 bg-teal-400 text-white rounded-lg hover:bg-teal-500 transition-colors"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full">
      {/* Sidebar - Desktop */}
      <aside className="hidden md:flex md:w-72 md:flex-shrink-0 border-r border-gray-200 bg-white">
        <ProjectList
          projects={projects}
          activeProjectId={selectedProject?.id}
          onNewProject={handleNewProject}
          onProjectSelect={handleProjectSelect}
          onProjectRename={handleProjectRename}
          onProjectDelete={handleProjectDelete}
        />
      </aside>

      {/* Sidebar - Mobile overlay */}
      {isSidebarOpen && (
        <>
          {/* Backdrop */}
          <div
            className="fixed inset-0 bg-black/50 z-40 md:hidden"
            onClick={handleCloseSidebar}
          />
          {/* Drawer */}
          <aside className="fixed inset-y-0 left-0 w-72 bg-white z-50 md:hidden shadow-xl">
            <ProjectList
              projects={projects}
              activeProjectId={selectedProject?.id}
              onNewProject={handleNewProject}
              onProjectSelect={handleProjectSelect}
              onProjectRename={handleProjectRename}
              onProjectDelete={handleProjectDelete}
            />
          </aside>
        </>
      )}

      {/* Main chat area */}
      <main className="flex-1 flex flex-col min-w-0">
        {selectedProject ? (
          isLoadingProject ? (
            <div className="flex-1 flex items-center justify-center">
              <LoadingSpinner size="lg" />
            </div>
          ) : (
            <ChatContainer
              projectId={selectedProject.id}
              projectTitle={selectedProject.title}
              initialMessages={projectMessages}
              onMenuClick={handleMenuClick}
            />
          )
        ) : (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center p-8">
              <h2 className="text-xl font-semibold text-gray-900 mb-2">
                Welcome to Go Chat
              </h2>
              <p className="text-gray-500 mb-4">
                Select a project or create a new one to get started
              </p>
              <button
                onClick={handleNewProject}
                className="px-6 py-3 bg-teal-400 text-white rounded-lg hover:bg-teal-500 transition-colors"
              >
                Create New Project
              </button>
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
