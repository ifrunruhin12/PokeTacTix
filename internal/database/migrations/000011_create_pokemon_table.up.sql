CREATE TABLE IF NOT EXISTS pokemon (
    id                  INTEGER PRIMARY KEY,
    name                TEXT NOT NULL UNIQUE,
    species_id          INTEGER NOT NULL,
    evolution_chain_id  INTEGER NOT NULL,
    generation          SMALLINT NOT NULL,
    types               TEXT[] NOT NULL,
    base_stats          JSONB NOT NULL,
    abilities           JSONB NOT NULL,
    sprite_url          TEXT,
    raw_json            JSONB NOT NULL,
    fetched_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_pokemon_evolution_chain ON pokemon(evolution_chain_id);
CREATE INDEX IF NOT EXISTS idx_pokemon_generation      ON pokemon(generation);
