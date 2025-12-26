-- Migration 006: Add agent_type to messages
-- Tracks which agent (product_manager, designer, developer) generated each assistant response

ALTER TABLE messages ADD COLUMN agent_type VARCHAR(50);

-- Add index for filtering by agent type
CREATE INDEX idx_messages_agent_type ON messages(agent_type) WHERE agent_type IS NOT NULL;

-- Add comment for documentation
COMMENT ON COLUMN messages.agent_type IS 'Agent that generated this message: product_manager, designer, developer. NULL for user messages.';
