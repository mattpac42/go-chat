-- 004_discovery.sql
-- Discovery flow tables for guided project setup
-- Stores discovery state and captured data for the 5-stage guided discovery process

-- Main discovery state table
CREATE TABLE IF NOT EXISTS project_discovery (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    stage VARCHAR(20) NOT NULL DEFAULT 'welcome'
        CHECK (stage IN ('welcome', 'problem', 'personas', 'mvp', 'summary', 'complete')),
    stage_started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Captured conversation data
    business_context TEXT,
    problem_statement TEXT,
    goals JSONB DEFAULT '[]'::JSONB,

    -- Summary data
    project_name VARCHAR(255),
    solves_statement TEXT,

    -- Metadata
    is_returning_user BOOLEAN DEFAULT FALSE,
    used_template_id UUID,
    confirmed_at TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT unique_project_discovery UNIQUE (project_id)
);

-- Discovery users (personas)
CREATE TABLE IF NOT EXISTS discovery_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discovery_id UUID NOT NULL REFERENCES project_discovery(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    user_count INTEGER DEFAULT 1,
    has_permissions BOOLEAN DEFAULT FALSE,
    permission_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Discovery features (MVP and future)
CREATE TABLE IF NOT EXISTS discovery_features (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discovery_id UUID NOT NULL REFERENCES project_discovery(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    version VARCHAR(10) NOT NULL DEFAULT 'v1',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Edit history for tracking changes
CREATE TABLE IF NOT EXISTS discovery_edit_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discovery_id UUID NOT NULL REFERENCES project_discovery(id) ON DELETE CASCADE,
    stage VARCHAR(20) NOT NULL,
    field_edited VARCHAR(100) NOT NULL,
    original_value TEXT,
    new_value TEXT,
    edited_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for project_discovery
CREATE INDEX IF NOT EXISTS idx_discovery_project ON project_discovery(project_id);
CREATE INDEX IF NOT EXISTS idx_discovery_stage ON project_discovery(stage);

-- Indexes for discovery_users
CREATE INDEX IF NOT EXISTS idx_discovery_users_discovery ON discovery_users(discovery_id);

-- Indexes for discovery_features
CREATE INDEX IF NOT EXISTS idx_discovery_features_discovery ON discovery_features(discovery_id);
CREATE INDEX IF NOT EXISTS idx_discovery_features_version ON discovery_features(version);

-- Indexes for discovery_edit_history
CREATE INDEX IF NOT EXISTS idx_discovery_edit_history_discovery ON discovery_edit_history(discovery_id);

-- Comments
COMMENT ON TABLE project_discovery IS 'Stores discovery flow state and captured data for each project';
COMMENT ON COLUMN project_discovery.stage IS 'Current stage in the discovery flow: welcome, problem, personas, mvp, summary, complete';
COMMENT ON COLUMN project_discovery.stage_started_at IS 'Timestamp when the current stage began';
COMMENT ON COLUMN project_discovery.business_context IS 'User''s business context captured during welcome stage';
COMMENT ON COLUMN project_discovery.problem_statement IS 'Pain points identified during problem discovery stage';
COMMENT ON COLUMN project_discovery.goals IS 'JSON array of user goals';
COMMENT ON COLUMN project_discovery.project_name IS 'Auto-generated or user-confirmed project name';
COMMENT ON COLUMN project_discovery.solves_statement IS 'Summary of what the project solves';
COMMENT ON COLUMN project_discovery.is_returning_user IS 'Whether user has completed discovery before';
COMMENT ON COLUMN project_discovery.used_template_id IS 'If fast-track, reference to template discovery used';
COMMENT ON COLUMN project_discovery.confirmed_at IS 'Timestamp when user confirmed summary and started building';

COMMENT ON TABLE discovery_users IS 'User personas identified during discovery';
COMMENT ON COLUMN discovery_users.description IS 'Description of the user type (e.g., "owner/baker", "order takers")';
COMMENT ON COLUMN discovery_users.user_count IS 'Number of users of this type';
COMMENT ON COLUMN discovery_users.has_permissions IS 'Whether this user type has elevated permissions';
COMMENT ON COLUMN discovery_users.permission_notes IS 'Notes about what permissions this user type has';

COMMENT ON TABLE discovery_features IS 'Features identified during MVP scoping with version assignments';
COMMENT ON COLUMN discovery_features.name IS 'Feature name (e.g., "Order list view")';
COMMENT ON COLUMN discovery_features.priority IS 'Priority order within the version (lower = higher priority)';
COMMENT ON COLUMN discovery_features.version IS 'Version assignment: v1 for MVP, v2+ for future';

COMMENT ON TABLE discovery_edit_history IS 'Tracks edits made to discovery data after initial capture';
COMMENT ON COLUMN discovery_edit_history.stage IS 'Stage the edit was made in';
COMMENT ON COLUMN discovery_edit_history.field_edited IS 'Field that was edited';
