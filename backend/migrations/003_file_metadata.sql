-- 003_file_metadata.sql
-- File metadata table for App Map feature
-- Stores descriptions and groupings for files in the project

-- File metadata table
CREATE TABLE IF NOT EXISTS file_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    short_description VARCHAR(255),
    long_description TEXT,
    functional_group VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_file_metadata_file FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

-- Unique constraint: one metadata record per file
CREATE UNIQUE INDEX IF NOT EXISTS idx_file_metadata_file_id ON file_metadata(file_id);

-- Index for querying by functional group
CREATE INDEX IF NOT EXISTS idx_file_metadata_functional_group ON file_metadata(functional_group);

-- Comments
COMMENT ON TABLE file_metadata IS 'Metadata for files in App Map feature';
COMMENT ON COLUMN file_metadata.file_id IS 'Reference to the file this metadata belongs to';
COMMENT ON COLUMN file_metadata.short_description IS 'Brief human-readable description of what the file does';
COMMENT ON COLUMN file_metadata.long_description IS 'Detailed explanation of the file purpose and functionality';
COMMENT ON COLUMN file_metadata.functional_group IS 'App Map group the file belongs to (e.g., Homepage, Backend Services)';
