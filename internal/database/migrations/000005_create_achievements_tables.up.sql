-- Create achievements table
CREATE TABLE IF NOT EXISTS achievements (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    requirement_type VARCHAR(50) NOT NULL,
    requirement_value INTEGER NOT NULL CHECK (requirement_value >= 0)
);

-- Create user_achievements table
CREATE TABLE IF NOT EXISTS user_achievements (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    achievement_id INTEGER REFERENCES achievements(id) ON DELETE CASCADE,
    unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, achievement_id)
);

-- Create indexes
CREATE INDEX idx_user_achievements_user_id ON user_achievements(user_id);
CREATE INDEX idx_user_achievements_unlocked_at ON user_achievements(unlocked_at DESC);

-- Insert default achievements
INSERT INTO achievements (name, description, icon, requirement_type, requirement_value) VALUES
    ('First Victory', 'Win your first battle', '🏆', 'total_wins', 1),
    ('Veteran Trainer', 'Win 10 battles', '⭐', 'total_wins', 10),
    ('Elite Trainer', 'Win 50 battles', '🌟', 'total_wins', 50),
    ('Champion', 'Win 100 battles', '👑', 'total_wins', 100),
    ('Legendary Collector', 'Obtain a legendary Pokemon', '💎', 'legendary_owned', 1),
    ('Mythical Master', 'Obtain a mythical Pokemon', '✨', 'mythical_owned', 1),
    ('Max Level', 'Get a Pokemon to level 50', '🔥', 'max_level', 50),
    ('Coin Hoarder', 'Accumulate 5000 coins', '💰', 'coins_total', 5000);
