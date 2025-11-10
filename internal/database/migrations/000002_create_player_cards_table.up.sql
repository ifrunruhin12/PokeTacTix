-- Create player_cards table
CREATE TABLE IF NOT EXISTS player_cards (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pokemon_name VARCHAR(100) NOT NULL,
    level INTEGER DEFAULT 1 CHECK (level >= 1 AND level <= 50),
    xp INTEGER DEFAULT 0 CHECK (xp >= 0),
    base_hp INTEGER NOT NULL CHECK (base_hp > 0),
    base_attack INTEGER NOT NULL CHECK (base_attack > 0),
    base_defense INTEGER NOT NULL CHECK (base_defense > 0),
    base_speed INTEGER NOT NULL CHECK (base_speed > 0),
    types JSONB NOT NULL,
    moves JSONB NOT NULL,
    sprite VARCHAR(255),
    is_legendary BOOLEAN DEFAULT FALSE,
    is_mythical BOOLEAN DEFAULT FALSE,
    in_deck BOOLEAN DEFAULT FALSE,
    deck_position INTEGER CHECK (deck_position >= 1 AND deck_position <= 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for frequently queried columns
CREATE INDEX idx_player_cards_user_id ON player_cards(user_id);
CREATE INDEX idx_player_cards_in_deck ON player_cards(user_id, in_deck);
CREATE INDEX idx_player_cards_deck_position ON player_cards(user_id, deck_position) WHERE in_deck = TRUE;

-- Create trigger to update updated_at timestamp
CREATE TRIGGER update_player_cards_updated_at BEFORE UPDATE ON player_cards
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add constraint to ensure exactly 5 cards in deck per user
CREATE OR REPLACE FUNCTION check_deck_size()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.in_deck = TRUE THEN
        IF (SELECT COUNT(*) FROM player_cards WHERE user_id = NEW.user_id AND in_deck = TRUE AND id != NEW.id) >= 5 THEN
            RAISE EXCEPTION 'User can only have 5 cards in deck';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER enforce_deck_size BEFORE INSERT OR UPDATE ON player_cards
    FOR EACH ROW EXECUTE FUNCTION check_deck_size();
