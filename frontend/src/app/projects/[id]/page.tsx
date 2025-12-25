import { ProjectPageClient } from '@/components/ProjectPageClient';

// Force dynamic rendering (no static generation)
export const dynamic = 'force-dynamic';

export default function ProjectPage() {
  return <ProjectPageClient />;
}
