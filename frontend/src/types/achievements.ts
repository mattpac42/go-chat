/**
 * Achievement and Learning Journey types
 * Supports the gamified learning experience
 */

export type AchievementCategory = 'exploration' | 'understanding' | 'mastery' | 'graduation';

export interface Achievement {
  id: string;
  code: string;
  name: string;
  description: string;
  category: AchievementCategory;
  icon: string;
  points: number;
}

export interface UserAchievement {
  id: string;
  achievementId: string;
  unlockedAt: string;
  isSeen: boolean;
  achievement?: Achievement;
}

export type LearningLevel = 1 | 2 | 3 | 4;

export interface UserProgress {
  id: string;
  projectId: string;
  currentLevel: LearningLevel;
  totalPoints: number;
  filesViewedCount: number;
  codeViewsCount: number;
  treeExpansionsCount: number;
  levelChangesCount: number;
  firstCodeViewAt?: string;
  firstLevelUpAt?: string;
  lastActivityAt: string;
}

export interface Nudge {
  type: string;
  title: string;
  message: string;
  action: string;
  icon: string;
}

export interface LearningEvent {
  type: string;
  context?: Record<string, unknown>;
}
