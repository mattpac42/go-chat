'use client';

import { useState, useCallback } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { ProjectList } from '@/components/projects/ProjectList';
import { ChatContainer } from '@/components/chat/ChatContainer';
import { FileExplorer } from '@/components/projects/FileExplorer';
import { LoadingSpinner } from '@/components/shared/LoadingSpinner';
import { useProjects, useProject } from '@/hooks/useProjects';
import { useFiles } from '@/hooks/useFiles';
import { Project } from '@/types';

export function ProjectPageClient() {
  const params = useParams();
  const router = useRouter();
  const projectId = params.id as string;

  const {
    projects,
    isLoading: isLoadingProjects,
    createProject,
    renameProject,
    deleteProject,
  } = useProjects();

  const {
    project: currentProject,
    messages: initialMessages,
    isLoading: isLoadingProject,
    error: projectError,
  } = useProject(projectId, projects);

  const {
    fileTree,
    isLoading: isLoadingFiles,
    fetchFiles,
    getFile,
  } = useFiles(projectId);

  console.log('[ProjectPageClient] fileTree:', fileTree, 'isLoadingFiles:', isLoadingFiles);

  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

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
    if (success) {
      // If deleted project was the current one, navigate to another or home
      if (project.id === projectId) {
        const remaining = projects.filter(p => p.id !== project.id);
        if (remaining.length > 0) {
          router.push(`/projects/${remaining[0].id}`);
        } else {
          router.push('/');
        }
      }
    }
  }, [deleteProject, projectId, projects, router]);

  // Show project not found error
  if (!isLoadingProject && projectError) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center p-8">
          <h2 className="text-xl font-semibold text-gray-900 mb-2">
            Project not found
          </h2>
          <p className="text-gray-500 mb-4">
            The project you are looking for does not exist.
          </p>
          <Link
            href="/"
            className="px-6 py-3 bg-teal-400 text-white rounded-lg hover:bg-teal-500 transition-colors inline-block"
          >
            Go to Home
          </Link>
        </div>
      </div>
    );
  }

  const projectTitle = currentProject?.title || 'New Project';

  return (
    <div className="flex h-full">
      {/* Sidebar - Desktop */}
      <aside className="hidden md:flex md:w-72 md:flex-shrink-0 border-r border-gray-200 bg-white">
        {isLoadingProjects ? (
          <div className="flex items-center justify-center w-full">
            <LoadingSpinner />
          </div>
        ) : (
          <ProjectList
            projects={projects}
            activeProjectId={projectId}
            onNewProject={handleNewProject}
            onProjectSelect={handleProjectSelect}
            onProjectRename={handleProjectRename}
            onProjectDelete={handleProjectDelete}
          />
        )}
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
            {isLoadingProjects ? (
              <div className="flex items-center justify-center h-full">
                <LoadingSpinner />
              </div>
            ) : (
              <ProjectList
                projects={projects}
                activeProjectId={projectId}
                onNewProject={handleNewProject}
                onProjectSelect={handleProjectSelect}
                onProjectRename={handleProjectRename}
                onProjectDelete={handleProjectDelete}
              />
            )}
          </aside>
        </>
      )}

      {/* Main chat area */}
      <main className="flex-1 flex flex-col min-w-0">
        {isLoadingProject ? (
          <div className="flex-1 flex items-center justify-center">
            <LoadingSpinner size="lg" />
          </div>
        ) : (
          <ChatContainer
            projectId={projectId}
            projectTitle={projectTitle}
            initialMessages={initialMessages}
            onMenuClick={handleMenuClick}
            onStreamingComplete={fetchFiles}
          />
        )}
      </main>

      {/* Right sidebar - Files with 2-tier reveal system */}
      <aside className="hidden lg:flex lg:w-80 lg:flex-shrink-0 border-l border-gray-200 bg-white flex-col">
        <div className="px-4 py-3 border-b border-gray-200">
          <h2 className="text-sm font-semibold text-gray-700">App Files</h2>
          <p className="text-xs text-gray-400 mt-0.5">
            Click to see what each file does
          </p>
        </div>
        <div className="flex-1 overflow-hidden">
          <FileExplorer
            files={fileTree}
            onLoadContent={getFile}
            defaultViewMode="reveal"
            showViewToggle={true}
            isLoading={isLoadingFiles}
          />
        </div>
      </aside>
    </div>
  );
}
