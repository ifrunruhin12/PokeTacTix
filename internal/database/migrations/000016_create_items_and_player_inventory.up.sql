-- Create items table: the purchasable game-item catalog (evolution stones now,
-- battle boosters later). item_id slugs intentionally match PokéAPI item names
-- (e.g. "thunder-stone") so evolution_chain links can reference them directly.
CREATE TABLE IF NOT EXISTS items (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('evolution', 'booster')),
    price INTEGER NOT NULL CHECK (price > 0),
    icon VARCHAR(255) NOT NULL DEFAULT '',
    effect JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_items_updated_at BEFORE UPDATE ON items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Player inventory: one row per (user, item) with a quantity. Purchases add,
-- evolutions consume; both run inside transactions that guard quantity > 0.
CREATE TABLE IF NOT EXISTS player_inventory (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id VARCHAR(50) NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, item_id)
);

CREATE TRIGGER update_player_inventory_updated_at BEFORE UPDATE ON player_inventory
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_player_inventory_user_id ON player_inventory(user_id);

-- Seed the initial evolution stones. Idempotent so re-running is safe.
INSERT INTO items (id, name, description, item_type, price, icon) VALUES
    ('thunder-stone', 'Thunder Stone',
     'A stone with a thunder-like pattern. Certain Pokemon, such as Pikachu, evolve when it is used on them.',
     'evolution', 800,
     'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/thunder-stone.png'),
    ('fire-stone', 'Fire Stone',
     'A stone with a fiery orange color. Certain Pokemon, such as Vulpix, evolve when it is used on them.',
     'evolution', 800,
     'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/fire-stone.png'),
    ('water-stone', 'Water Stone',
     'A stone with a blue, watery core. Certain Pokemon, such as Eevee, evolve when it is used on them.',
     'evolution', 800,
     'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/water-stone.png')
ON CONFLICT (id) DO NOTHING;
