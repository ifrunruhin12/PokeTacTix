-- Drop indexes
DROP INDEX IF EXISTS idx_battle_history_user_created;
DROP INDEX IF EXISTS idx_battle_history_created_at;
DROP INDEX IF EXISTS idx_battle_history_user_id;

-- Drop battle_history table
DROP TABLE IF EXISTS battle_history;
