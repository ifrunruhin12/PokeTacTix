-- Add draw columns to player_stats table
ALTER TABLE player_stats 
ADD COLUMN draws_1v1 INTEGER DEFAULT 0 CHECK (draws_1v1 >= 0),
ADD COLUMN draws_5v5 INTEGER DEFAULT 0 CHECK (draws_5v5 >= 0);

-- Update the check_stats_consistency trigger function to validate wins + losses + draws ≤ total_battles
CREATE OR REPLACE FUNCTION check_stats_consistency()
RETURNS TRIGGER AS $$
BEGIN
    -- Check 1v1 stats consistency
    IF (NEW.wins_1v1 + NEW.losses_1v1 + NEW.draws_1v1) > NEW.total_battles_1v1 THEN
        RAISE EXCEPTION 'Sum of wins, losses, and draws cannot exceed total battles for 1v1';
    END IF;
    
    -- Check 5v5 stats consistency
    IF (NEW.wins_5v5 + NEW.losses_5v5 + NEW.draws_5v5) > NEW.total_battles_5v5 THEN
        RAISE EXCEPTION 'Sum of wins, losses, and draws cannot exceed total battles for 5v5';
    END IF;
    
    -- Individual checks to maintain backward compatibility
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
