-- Drop triggers and functions
DROP TRIGGER IF EXISTS enforce_stats_consistency ON player_stats;
DROP FUNCTION IF EXISTS check_stats_consistency();
DROP TRIGGER IF EXISTS update_player_stats_updated_at ON player_stats;

-- Drop player_stats table
DROP TABLE IF EXISTS player_stats;
