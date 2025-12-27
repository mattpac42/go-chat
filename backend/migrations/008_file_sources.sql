-- Migration 008: Add file_sources table for tracking uploaded source files
-- This table tracks original uploaded files (images, PDFs, etc.) that were converted to markdown

CREATE TABLE IF NOT EXISTS file_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    original_filename VARCHAR(255) NOT NULL,
    original_mime_type VARCHAR(100) NOT NULL,
    original_size_bytes BIGINT NOT NULL,
    conversion_status VARCHAR(20) DEFAULT 'completed',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add index for looking up sources by file_id
CREATE INDEX IF NOT EXISTS idx_file_sources_file_id ON file_sources(file_id);

-- Add source_type column to files table to distinguish generated vs uploaded files
ALTER TABLE files ADD COLUMN IF NOT EXISTS source_type VARCHAR(20) DEFAULT 'generated';

-- Add index for filtering files by source_type
CREATE INDEX IF NOT EXISTS idx_files_source_type ON files(source_type);
