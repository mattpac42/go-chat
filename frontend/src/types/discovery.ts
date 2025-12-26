/**
 * Discovery stage represents the current phase of the guided discovery process
 */
export type DiscoveryStage =
  | 'welcome'
  | 'problem'
  | 'personas'
  | 'mvp'
  | 'summary'
  | 'complete';

/**
 * Project discovery state returned from the API
 */
export interface ProjectDiscovery {
  id: string;
  projectId: string;
  stage: DiscoveryStage;
  stageStartedAt: string;
  businessContext?: string;
  problemStatement?: string;
  goals?: string[];
  projectName?: string;
  solvesStatement?: string;
  isReturningUser: boolean;
  confirmedAt?: string;
}

/**
 * User persona captured during discovery
 */
export interface DiscoveryUser {
  id: string;
  description: string;
  count: number;
  hasPermissions: boolean;
  permissionNotes?: string;
}

/**
 * Feature captured during MVP scoping
 */
export interface DiscoveryFeature {
  id: string;
  name: string;
  priority: number;
  version: string;
}

/**
 * Complete discovery summary for display
 */
export interface DiscoverySummary {
  projectName: string;
  solvesStatement: string;
  users: DiscoveryUser[];
  mvpFeatures: DiscoveryFeature[];
  futureFeatures: DiscoveryFeature[];
}
