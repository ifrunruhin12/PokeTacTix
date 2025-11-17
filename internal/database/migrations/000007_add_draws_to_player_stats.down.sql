-- Restore original check_stats_consistency function (without draw validation)
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

-- Remove draw columns from player_stats table
ALTER TABLE player_stats 
DROP COLUMN IF EXISTS draws_1v1,
DROP COLUMN IF EXISTS draws_5v5;
