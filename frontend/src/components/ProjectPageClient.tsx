'use client';

import { useState, useCallback, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { ProjectList } from '@/components/projects/ProjectList';
import { ChatContainer } from '@/components/chat/ChatContainer';
import { FileExplorer } from '@/components/projects/FileExplorer';
import { ProjectPreview } from '@/components/preview/ProjectPreview';
import { PreviewModal } from '@/components/preview/PreviewModal';
import { LoadingSpinner } from '@/components/shared/LoadingSpinner';
import { useProjects, useProject } from '@/hooks/useProjects';
import { useFiles } from '@/hooks/useFiles';
import { usePreviewFiles } from '@/hooks/usePreviewFiles';
import { Project } from '@/types';
import { API_BASE_URL } from '@/lib/api';

type RightPanelView = 'files' | 'preview';

export function ProjectPageClient() {
  const params = useParams();
  const router = useRouter();
  const projectId = params.id as string;

  const {
    projects,
    isLoading: isLoadingProjects,
    fetchProjects,
    createProject,
    renameProject,
    deleteProject,
  } = useProjects();

  const {
    project: currentProject,
    messages: initialMessages,
    isLoading: isLoadingProject,
    error: projectError,
    fetchProject,
  } = useProject(projectId, projects);

  const {
    files,
    fileTree,
    isLoading: isLoadingFiles,
    fetchFiles,
    getFile,
  } = useFiles(projectId);

  console.log('[ProjectPageClient] fileTree:', fileTree, 'isLoadingFiles:', isLoadingFiles);

  const { previewFiles, isLoading: isLoadingPreview, loadAllFiles } = usePreviewFiles();

  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [rightPanelView, setRightPanelView] = useState<RightPanelView>('files');
  const [isPreviewModalOpen, setIsPreviewModalOpen] = useState(false);

  // Load preview files when files change and preview tab is active
  useEffect(() => {
    if (rightPanelView === 'preview' && files.length > 0) {
      loadAllFiles(files);
    }
  }, [rightPanelView, files, loadAllFiles]);

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

  const handleTitleUpdate = useCallback(async (newTitle: string) => {
    await renameProject(projectId, newTitle);
    // Refresh project list to update sidebar
    await fetchProjects();
  }, [renameProject, projectId, fetchProjects]);

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
      <main className="flex-1 flex flex-col min-w-0 overflow-visible">
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
            onFilesUpdated={fetchFiles}
            onDiscoveryConfirmed={fetchProjects}
            onRefetchMessages={fetchProject}
            onTitleUpdate={handleTitleUpdate}
          />
        )}
      </main>

      {/* Right sidebar - Files/Preview with tabs */}
      <aside className="hidden lg:flex lg:w-80 lg:flex-shrink-0 border-l border-gray-200 bg-white flex-col">
        {/* Tab header */}
        <div className="px-4 py-2 border-b border-gray-200">
          <div className="flex items-center gap-1">
            <button
              onClick={() => setRightPanelView('files')}
              className={`px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
                rightPanelView === 'files'
                  ? 'bg-teal-100 text-teal-700'
                  : 'text-gray-500 hover:bg-gray-100 hover:text-gray-700'
              }`}
            >
              Files
            </button>
            <button
              onClick={() => setRightPanelView('preview')}
              className={`px-3 py-1.5 text-sm font-medium rounded-md transition-colors flex items-center gap-1.5 ${
                rightPanelView === 'preview'
                  ? 'bg-teal-100 text-teal-700'
                  : 'text-gray-500 hover:bg-gray-100 hover:text-gray-700'
              }`}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              Preview
            </button>
            {rightPanelView === 'files' && fileTree.length > 0 && (
              <button
                onClick={() => {
                  window.open(`${API_BASE_URL}/api/projects/${projectId}/download`, '_blank');
                }}
                className="ml-auto p-1.5 rounded-lg hover:bg-gray-100 transition-colors"
                title="Download all files as ZIP"
                aria-label="Download all files as ZIP"
              >
                <svg className="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
              </button>
            )}
            {rightPanelView === 'preview' && (
              <button
                onClick={() => setIsPreviewModalOpen(true)}
                className="ml-auto p-1.5 rounded-lg hover:bg-gray-100 transition-colors"
                title="Open fullscreen preview"
                aria-label="Open fullscreen preview"
              >
                <svg className="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
                </svg>
              </button>
            )}
          </div>
        </div>

        {/* Files view */}
        {rightPanelView === 'files' && (
          <>
            <div className="px-4 py-2 border-b border-gray-100">
              <p className="text-xs text-gray-400">
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
          </>
        )}

        {/* Preview view */}
        {rightPanelView === 'preview' && (
          <div className="flex-1 overflow-hidden">
            {isLoadingPreview ? (
              <div className="flex items-center justify-center h-full">
                <LoadingSpinner size="lg" />
              </div>
            ) : (
              <ProjectPreview files={previewFiles} />
            )}
          </div>
        )}
      </aside>

      {/* Fullscreen Preview Modal */}
      <PreviewModal
        files={previewFiles}
        isOpen={isPreviewModalOpen}
        onClose={() => setIsPreviewModalOpen(false)}
      />
    </div>
  );
}
