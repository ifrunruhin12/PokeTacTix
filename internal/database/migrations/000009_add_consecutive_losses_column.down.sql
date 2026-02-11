-- Remove consecutive_losses column from player_stats table
ALTER TABLE player_stats 
DROP COLUMN IF EXISTS consecutive_losses;
