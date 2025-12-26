-- 007_learning_journey.sql
-- Learning Journey & Progression System
-- Phase 3: Gamified learning milestones and user advancement tracking

-- ============================================================================
-- ACHIEVEMENTS DEFINITION TABLE
-- Stores all available achievements (seeded at startup, rarely changes)
-- ============================================================================
CREATE TABLE IF NOT EXISTS achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,          -- e.g., 'first_look', 'explorer'
    name VARCHAR(100) NOT NULL,                -- Display name: "First Look"
    description TEXT NOT NULL,                 -- What this achievement means
    category VARCHAR(30) NOT NULL              -- 'exploration', 'understanding', 'mastery'
        CHECK (category IN ('exploration', 'understanding', 'mastery', 'graduation')),
    icon VARCHAR(50) NOT NULL,                 -- Icon identifier for frontend
    points INTEGER NOT NULL DEFAULT 10,        -- Points value for gamification
    trigger_type VARCHAR(30) NOT NULL          -- How it's triggered
        CHECK (trigger_type IN ('event', 'count', 'milestone')),
    trigger_config JSONB DEFAULT '{}'::JSONB,  -- Configuration for trigger logic
    prerequisites TEXT[] DEFAULT '{}',         -- Array of achievement codes required first
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================================
-- USER PROGRESS TABLE
-- Tracks overall learning progression per project
-- ============================================================================
CREATE TABLE IF NOT EXISTS user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

    -- Current learning level (1=Functional, 2=Tree, 3=Technical, 4=Developer)
    current_level INTEGER NOT NULL DEFAULT 1
        CHECK (current_level BETWEEN 1 AND 4),

    -- Total points accumulated
    total_points INTEGER NOT NULL DEFAULT 0,

    -- Stats for trigger evaluation
    files_viewed_count INTEGER NOT NULL DEFAULT 0,
    code_views_count INTEGER NOT NULL DEFAULT 0,
    tree_expansions_count INTEGER NOT NULL DEFAULT 0,
    level_changes_count INTEGER NOT NULL DEFAULT 0,

    -- Timestamps
    first_code_view_at TIMESTAMP WITH TIME ZONE,
    first_level_up_at TIMESTAMP WITH TIME ZONE,
    last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT unique_project_progress UNIQUE (project_id)
);

-- ============================================================================
-- USER ACHIEVEMENTS TABLE
-- Junction table: which achievements each project has unlocked
-- ============================================================================
CREATE TABLE IF NOT EXISTS user_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,

    -- Context when achievement was earned
    unlocked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    trigger_context JSONB DEFAULT '{}'::JSONB,  -- What triggered it (file, action, etc.)

    -- Whether user has seen the notification
    is_seen BOOLEAN DEFAULT FALSE,
    seen_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_user_achievement UNIQUE (project_id, achievement_id)
);

-- ============================================================================
-- NUDGE HISTORY TABLE
-- Tracks which nudges have been shown to avoid repetition
-- ============================================================================
CREATE TABLE IF NOT EXISTS nudge_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    nudge_type VARCHAR(50) NOT NULL,           -- 'explore_code', 'try_tree_view', etc.

    -- When nudge was shown and user response
    shown_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    dismissed_at TIMESTAMP WITH TIME ZONE,
    clicked_at TIMESTAMP WITH TIME ZONE,       -- If user followed the nudge

    -- Context
    context JSONB DEFAULT '{}'::JSONB
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Achievements lookup
CREATE INDEX IF NOT EXISTS idx_achievements_code ON achievements(code);
CREATE INDEX IF NOT EXISTS idx_achievements_category ON achievements(category);
CREATE INDEX IF NOT EXISTS idx_achievements_active ON achievements(is_active) WHERE is_active = TRUE;

-- User progress
CREATE INDEX IF NOT EXISTS idx_user_progress_project ON user_progress(project_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_level ON user_progress(current_level);

-- User achievements
CREATE INDEX IF NOT EXISTS idx_user_achievements_project ON user_achievements(project_id);
CREATE INDEX IF NOT EXISTS idx_user_achievements_unseen ON user_achievements(project_id, is_seen)
    WHERE is_seen = FALSE;
CREATE INDEX IF NOT EXISTS idx_user_achievements_recent ON user_achievements(unlocked_at DESC);

-- Nudge history
CREATE INDEX IF NOT EXISTS idx_nudge_history_project ON nudge_history(project_id);
CREATE INDEX IF NOT EXISTS idx_nudge_history_type ON nudge_history(nudge_type);

-- ============================================================================
-- SEED DATA: Initial Achievements
-- ============================================================================
INSERT INTO achievements (code, name, description, category, icon, points, trigger_type, trigger_config) VALUES
-- Exploration achievements
('first_look', 'First Look', 'Viewed code for the first time', 'exploration', 'eye', 10,
 'event', '{"event": "code_viewed", "count": 1}'::JSONB),

('curious_mind', 'Curious Mind', 'Viewed code for 5 different files', 'exploration', 'search', 20,
 'count', '{"metric": "files_viewed_count", "threshold": 5}'::JSONB),

('deep_diver', 'Deep Diver', 'Viewed code for 10 different files', 'exploration', 'layers', 30,
 'count', '{"metric": "files_viewed_count", "threshold": 10}'::JSONB),

-- Understanding achievements
('connection_maker', 'Connection Maker', 'Expanded a relationship in the tree view', 'understanding', 'git-branch', 15,
 'event', '{"event": "tree_expanded", "count": 1}'::JSONB),

('big_picture', 'Big Picture', 'Viewed the full application tree', 'understanding', 'network', 25,
 'event', '{"event": "full_tree_viewed"}'::JSONB),

-- Mastery achievements
('level_up', 'Level Up', 'Advanced from Level 1 to Level 2 view', 'mastery', 'trending-up', 25,
 'event', '{"event": "level_changed", "from": 1, "to": 2}'::JSONB),

('explorer', 'Explorer', 'Viewed the full technical tree (Level 3)', 'mastery', 'compass', 35,
 'event', '{"event": "level_changed", "to": 3}'::JSONB),

('technologist', 'Technologist', 'Spent time in developer view (Level 4)', 'mastery', 'code', 50,
 'event', '{"event": "level_changed", "to": 4}'::JSONB),

-- Graduation achievements
('graduate', 'Graduate', 'Exported project to VS Code or local development', 'graduation', 'graduation-cap', 100,
 'event', '{"event": "project_exported"}'::JSONB),

('self_sufficient', 'Self-Sufficient', 'Made a direct code edit after learning', 'graduation', 'edit', 75,
 'event', '{"event": "code_edited"}'::JSONB)

ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- COMMENTS
-- ============================================================================
COMMENT ON TABLE achievements IS 'Defines all available achievements users can unlock';
COMMENT ON TABLE user_progress IS 'Tracks learning progression and stats per project';
COMMENT ON TABLE user_achievements IS 'Records which achievements each project has unlocked';
COMMENT ON TABLE nudge_history IS 'Tracks nudge interactions to avoid repetition';

COMMENT ON COLUMN achievements.trigger_type IS 'event=single action, count=accumulation threshold, milestone=complex condition';
COMMENT ON COLUMN achievements.trigger_config IS 'JSON configuration for trigger logic evaluation';
COMMENT ON COLUMN user_progress.current_level IS '1=Functional, 2=Tree, 3=Technical, 4=Developer';
