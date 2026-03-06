-- Add game token fields to users table
ALTER TABLE users 
ADD COLUMN game_tokens INTEGER DEFAULT 10 CHECK (game_tokens >= 0),
ADD COLUMN last_token_reset DATE DEFAULT CURRENT_DATE,
ADD COLUMN tokens_purchased_today INTEGER DEFAULT 0 CHECK (tokens_purchased_today >= 0);

-- Create index for efficient token reset queries
CREATE INDEX idx_users_last_token_reset ON users(last_token_reset);
