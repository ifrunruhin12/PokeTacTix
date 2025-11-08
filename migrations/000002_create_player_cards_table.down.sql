-- Drop triggers and functions
DROP TRIGGER IF EXISTS enforce_deck_size ON player_cards;
DROP FUNCTION IF EXISTS check_deck_size();
DROP TRIGGER IF EXISTS update_player_cards_updated_at ON player_cards;

-- Drop indexes
DROP INDEX IF EXISTS idx_player_cards_deck_position;
DROP INDEX IF EXISTS idx_player_cards_in_deck;
DROP INDEX IF EXISTS idx_player_cards_user_id;

-- Drop player_cards table
DROP TABLE IF EXISTS player_cards;
