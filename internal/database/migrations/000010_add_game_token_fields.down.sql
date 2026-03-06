-- Remove game token fields from users table
DROP INDEX IF EXISTS idx_users_last_token_reset;

ALTER TABLE users 
DROP COLUMN IF EXISTS tokens_purchased_today,
DROP COLUMN IF EXISTS last_token_reset,
DROP COLUMN IF EXISTS game_tokens;
