-- 005_prds.sql
-- PRD (Product Requirements Document) tables for feature tracking
-- Auto-generated from discovery features for agent context and build phase guidance

-- PRD status enum
CREATE TYPE prd_status AS ENUM (
    'pending',
    'generating',
    'draft',
    'ready',
    'in_progress',
    'complete',
    'failed'
);

-- Main PRD table
CREATE TABLE IF NOT EXISTS prds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discovery_id UUID NOT NULL REFERENCES project_discovery(id) ON DELETE CASCADE,
    feature_id UUID NOT NULL REFERENCES discovery_features(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

    -- Core Content
    title VARCHAR(255) NOT NULL,
    overview TEXT,
    version VARCHAR(10) NOT NULL DEFAULT 'v1',
    priority INTEGER NOT NULL DEFAULT 0,

    -- Structured Sections (JSONB)
    user_stories JSONB DEFAULT '[]'::JSONB,
    acceptance_criteria JSONB DEFAULT '[]'::JSONB,
    technical_notes JSONB DEFAULT '[]'::JSONB,

    -- Status Tracking
    status prd_status NOT NULL DEFAULT 'pending',
    generated_at TIMESTAMP WITH TIME ZONE,
    approved_at TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,

    -- Generation Metadata
    generation_attempts INTEGER DEFAULT 0,
    last_error TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT unique_feature_prd UNIQUE (feature_id),
    CONSTRAINT valid_priority CHECK (priority >= 0)
);

-- Add active PRD reference to projects
ALTER TABLE projects
ADD COLUMN IF NOT EXISTS active_prd_id UUID REFERENCES prds(id) ON DELETE SET NULL;

-- Indexes for prds
CREATE INDEX IF NOT EXISTS idx_prds_discovery ON prds(discovery_id);
CREATE INDEX IF NOT EXISTS idx_prds_project ON prds(project_id);
CREATE INDEX IF NOT EXISTS idx_prds_status ON prds(status);
CREATE INDEX IF NOT EXISTS idx_prds_version ON prds(version);
CREATE INDEX IF NOT EXISTS idx_prds_priority ON prds(priority);

-- Index for projects active PRD
CREATE INDEX IF NOT EXISTS idx_projects_active_prd ON projects(active_prd_id);

-- Updated_at trigger
CREATE OR REPLACE FUNCTION update_prds_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_prds_updated_at
    BEFORE UPDATE ON prds
    FOR EACH ROW
    EXECUTE FUNCTION update_prds_updated_at();

-- Comments
COMMENT ON TABLE prds IS 'Product Requirements Documents auto-generated from discovery features';
COMMENT ON COLUMN prds.discovery_id IS 'Source discovery that generated this PRD';
COMMENT ON COLUMN prds.feature_id IS 'Feature this PRD documents (one PRD per feature)';
COMMENT ON COLUMN prds.project_id IS 'Project this PRD belongs to';
COMMENT ON COLUMN prds.title IS 'PRD title (typically matches feature name)';
COMMENT ON COLUMN prds.overview IS 'Brief description of what this feature does and why it matters';
COMMENT ON COLUMN prds.version IS 'Version assignment: v1 for MVP, v2+ for future';
COMMENT ON COLUMN prds.priority IS 'Priority order within the version (lower = higher priority)';
COMMENT ON COLUMN prds.user_stories IS 'JSON array of user stories with As-a/I-want/So-that format';
COMMENT ON COLUMN prds.acceptance_criteria IS 'JSON array of Given/When/Then acceptance criteria';
COMMENT ON COLUMN prds.technical_notes IS 'JSON array of implementation notes by category';
COMMENT ON COLUMN prds.status IS 'Lifecycle status: pending->generating->draft->ready->in_progress->complete';
COMMENT ON COLUMN prds.generated_at IS 'Timestamp when Claude completed PRD generation';
COMMENT ON COLUMN prds.approved_at IS 'Timestamp when PRD was approved and marked ready';
COMMENT ON COLUMN prds.started_at IS 'Timestamp when implementation began';
COMMENT ON COLUMN prds.completed_at IS 'Timestamp when feature implementation was completed';
COMMENT ON COLUMN prds.generation_attempts IS 'Number of times Claude generation was attempted';
COMMENT ON COLUMN prds.last_error IS 'Error message from last failed generation attempt';
COMMENT ON COLUMN projects.active_prd_id IS 'Currently active PRD for focused agent context';
