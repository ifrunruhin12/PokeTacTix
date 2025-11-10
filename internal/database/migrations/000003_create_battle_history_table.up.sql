-- Create battle_history table
CREATE TABLE IF NOT EXISTS battle_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mode VARCHAR(10) NOT NULL CHECK (mode IN ('1v1', '5v5')),
    result VARCHAR(10) NOT NULL CHECK (result IN ('win', 'loss', 'draw')),
    coins_earned INTEGER DEFAULT 0 CHECK (coins_earned >= 0),
    duration INTEGER CHECK (duration >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for frequently queried columns
CREATE INDEX idx_battle_history_user_id ON battle_history(user_id);
CREATE INDEX idx_battle_history_created_at ON battle_history(created_at DESC);
CREATE INDEX idx_battle_history_user_created ON battle_history(user_id, created_at DESC);
