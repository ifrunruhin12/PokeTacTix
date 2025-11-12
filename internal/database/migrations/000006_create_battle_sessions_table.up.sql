-- Create battle_sessions table to persist battle state
CREATE TABLE IF NOT EXISTS battle_sessions (
    session_id VARCHAR(255) PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    state_json JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on user_id for faster lookups
CREATE INDEX idx_battle_sessions_user_id ON battle_sessions(user_id);

-- Create index on updated_at for cleanup queries
CREATE INDEX idx_battle_sessions_updated_at ON battle_sessions(updated_at);
