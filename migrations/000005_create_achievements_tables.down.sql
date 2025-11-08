-- Drop indexes
DROP INDEX IF EXISTS idx_user_achievements_unlocked_at;
DROP INDEX IF EXISTS idx_user_achievements_user_id;

-- Drop user_achievements table
DROP TABLE IF EXISTS user_achievements;

-- Drop achievements table
DROP TABLE IF EXISTS achievements;
