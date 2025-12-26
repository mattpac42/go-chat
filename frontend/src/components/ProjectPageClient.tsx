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
import { API_BASE_URL } from '@/lib/api';

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
      <aside className="hidden md:flex md:w-72 md:flex-shrink-0 border-r border-gray-200 bg-white flex-col">
        {isLoadingProjects ? (
          <div className="flex items-center justify-center w-full flex-1">
            <LoadingSpinner />
          </div>
        ) : (
          <>
            <div className="flex-shrink-0">
              <ProjectList
                projects={projects}
                activeProjectId={projectId}
                onNewProject={handleNewProject}
                onProjectSelect={handleProjectSelect}
                onProjectRename={handleProjectRename}
                onProjectDelete={handleProjectDelete}
              />
            </div>
            {/* Files section - visible on md screens where right panel is hidden */}
            {fileTree.length > 0 && (
              <div className="flex-1 flex flex-col border-t border-gray-200 lg:hidden overflow-hidden">
                <div className="px-4 py-3 border-b border-gray-200 flex-shrink-0">
                  <div className="flex items-center justify-between">
                    <h2 className="text-sm font-semibold text-gray-700">App Files</h2>
                    <button
                      onClick={() => {
                        window.open(`${API_BASE_URL}/api/projects/${projectId}/download`, '_blank');
                      }}
                      className="p-1.5 rounded-lg hover:bg-gray-100 transition-colors"
                      title="Download all files as ZIP"
                      aria-label="Download all files as ZIP"
                    >
                      <svg className="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                      </svg>
                    </button>
                  </div>
                </div>
                <div className="flex-1 overflow-hidden">
                  <FileExplorer
                    files={fileTree}
                    onLoadContent={getFile}
                    defaultViewMode="grouped"
                    showViewToggle={true}
                    isLoading={isLoadingFiles}
                  />
                </div>
              </div>
            )}
          </>
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
          <aside className="fixed inset-y-0 left-0 w-72 bg-white z-50 md:hidden shadow-xl flex flex-col">
            {isLoadingProjects ? (
              <div className="flex items-center justify-center h-full">
                <LoadingSpinner />
              </div>
            ) : (
              <>
                <div className="flex-shrink-0">
                  <ProjectList
                    projects={projects}
                    activeProjectId={projectId}
                    onNewProject={handleNewProject}
                    onProjectSelect={handleProjectSelect}
                    onProjectRename={handleProjectRename}
                    onProjectDelete={handleProjectDelete}
                  />
                </div>
                {/* Files section for mobile */}
                {fileTree.length > 0 && (
                  <div className="flex-1 flex flex-col border-t border-gray-200 overflow-hidden">
                    <div className="px-4 py-3 border-b border-gray-200 flex-shrink-0">
                      <div className="flex items-center justify-between">
                        <h2 className="text-sm font-semibold text-gray-700">App Files</h2>
                        <button
                          onClick={() => {
                            window.open(`${API_BASE_URL}/api/projects/${projectId}/download`, '_blank');
                          }}
                          className="p-1.5 rounded-lg hover:bg-gray-100 transition-colors"
                          title="Download all files as ZIP"
                          aria-label="Download all files as ZIP"
                        >
                          <svg className="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                          </svg>
                        </button>
                      </div>
                    </div>
                    <div className="flex-1 overflow-hidden">
                      <FileExplorer
                        files={fileTree}
                        onLoadContent={getFile}
                        defaultViewMode="grouped"
                        showViewToggle={true}
                        isLoading={isLoadingFiles}
                      />
                    </div>
                  </div>
                )}
              </>
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
          <div className="flex items-center justify-between">
            <h2 className="text-sm font-semibold text-gray-700">App Files</h2>
            {fileTree.length > 0 && (
              <button
                onClick={() => {
                  window.open(`${API_BASE_URL}/api/projects/${projectId}/download`, '_blank');
                }}
                className="p-1.5 rounded-lg hover:bg-gray-100 transition-colors"
                title="Download all files as ZIP"
                aria-label="Download all files as ZIP"
              >
                <svg className="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
              </button>
            )}
          </div>
          <p className="text-xs text-gray-400 mt-0.5">
            Click to see what each file does
          </p>
        </div>
        <div className="flex-1 overflow-hidden">
          <FileExplorer
            files={fileTree}
            onLoadContent={getFile}
            defaultViewMode="grouped"
            showViewToggle={true}
            isLoading={isLoadingFiles}
          />
        </div>
      </aside>
    </div>
  );
}
