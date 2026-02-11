-- Add consecutive_losses column to player_stats table
ALTER TABLE player_stats 
ADD COLUMN consecutive_losses INTEGER DEFAULT 0 CHECK (consecutive_losses >= 0);
