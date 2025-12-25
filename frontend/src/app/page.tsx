import { HomeClient } from '@/components/HomeClient';

// Force dynamic rendering (no static generation)
export const dynamic = 'force-dynamic';

export default function Home() {
  return <HomeClient />;
}
