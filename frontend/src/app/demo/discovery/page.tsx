import { DiscoveryDemoClient } from '@/components/demo/DiscoveryDemoClient';

// Force dynamic rendering (no static generation)
export const dynamic = 'force-dynamic';

export default function DiscoveryDemoPage() {
  return <DiscoveryDemoClient />;
}
