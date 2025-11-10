-- Create player_stats table
CREATE TABLE IF NOT EXISTS player_stats (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_battles_1v1 INTEGER DEFAULT 0 CHECK (total_battles_1v1 >= 0),
    wins_1v1 INTEGER DEFAULT 0 CHECK (wins_1v1 >= 0),
    losses_1v1 INTEGER DEFAULT 0 CHECK (losses_1v1 >= 0),
    total_battles_5v5 INTEGER DEFAULT 0 CHECK (total_battles_5v5 >= 0),
    wins_5v5 INTEGER DEFAULT 0 CHECK (wins_5v5 >= 0),
    losses_5v5 INTEGER DEFAULT 0 CHECK (losses_5v5 >= 0),
    total_coins_earned INTEGER DEFAULT 0 CHECK (total_coins_earned >= 0),
    highest_level INTEGER DEFAULT 1 CHECK (highest_level >= 1 AND highest_level <= 50),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger to update updated_at timestamp
CREATE TRIGGER update_player_stats_updated_at BEFORE UPDATE ON player_stats
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add constraint to ensure wins don't exceed total battles
CREATE OR REPLACE FUNCTION check_stats_consistency()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.wins_1v1 > NEW.total_battles_1v1 THEN
        RAISE EXCEPTION 'wins_1v1 cannot exceed total_battles_1v1';
    END IF;
    IF NEW.losses_1v1 > NEW.total_battles_1v1 THEN
        RAISE EXCEPTION 'losses_1v1 cannot exceed total_battles_1v1';
    END IF;
    IF NEW.wins_5v5 > NEW.total_battles_5v5 THEN
        RAISE EXCEPTION 'wins_5v5 cannot exceed total_battles_5v5';
    END IF;
    IF NEW.losses_5v5 > NEW.total_battles_5v5 THEN
        RAISE EXCEPTION 'losses_5v5 cannot exceed total_battles_5v5';
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER enforce_stats_consistency BEFORE INSERT OR UPDATE ON player_stats
    FOR EACH ROW EXECUTE FUNCTION check_stats_consistency();
