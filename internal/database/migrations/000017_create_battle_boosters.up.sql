-- Battle booster items: purchasable consumables that buff the player's deck
-- for the next N battles. The effect parameters live on the catalog row
-- (items.effect JSONB) so balancing is a data change, not a code change.
INSERT INTO items (id, name, description, item_type, price, icon, effect) VALUES
    ('attack-booster', 'Attack Booster',
     'All Pokemon in your deck gain +5 attack for the next 3 battles.',
     'booster', 300,
     'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/x-attack.png',
     '{"stat": "attack", "bonus": 5, "battles": 3}'),
    ('hp-booster', 'HP Booster',
     'All Pokemon in your deck gain +10 HP for the next 4 battles.',
     'booster', 300,
     'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/x-defense.png',
     '{"stat": "hp", "bonus": 10, "battles": 4}')
ON CONFLICT (id) DO NOTHING;

-- Active boosts are consumed when a booster is used; each row buffs the
-- player's whole deck and ticks down by one at the end of every completed
-- battle. Multiple active boosts stack (their bonuses add up).
CREATE TABLE IF NOT EXISTS active_boosts (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id VARCHAR(50) NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    stat VARCHAR(20) NOT NULL CHECK (stat IN ('attack', 'hp')),
    bonus INTEGER NOT NULL CHECK (bonus > 0),
    battles_remaining INTEGER NOT NULL CHECK (battles_remaining > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_active_boosts_updated_at BEFORE UPDATE ON active_boosts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_active_boosts_user_id ON active_boosts(user_id);
