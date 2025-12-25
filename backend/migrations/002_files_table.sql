-- 002_files_table.sql
-- Files table for storing extracted code from AI responses

-- Files table
CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    path VARCHAR(500) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    language VARCHAR(50),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_files_project_id ON files(project_id);
CREATE INDEX IF NOT EXISTS idx_files_path ON files(path);

-- Unique constraint on project_id + path to prevent duplicates
CREATE UNIQUE INDEX IF NOT EXISTS idx_files_project_path ON files(project_id, path);

-- Comments
COMMENT ON TABLE files IS 'Extracted code files from AI responses';
COMMENT ON COLUMN files.path IS 'Full file path (e.g., src/components/Button.tsx)';
COMMENT ON COLUMN files.filename IS 'Just the filename (e.g., Button.tsx)';
COMMENT ON COLUMN files.language IS 'Programming language detected from code block';
COMMENT ON COLUMN files.content IS 'File content/code';
