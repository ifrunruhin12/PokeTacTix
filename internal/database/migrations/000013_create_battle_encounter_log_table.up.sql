CREATE TABLE IF NOT EXISTS battle_encounter_log (
    id                  BIGSERIAL PRIMARY KEY,
    player_id           INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pokemon_id          INTEGER NOT NULL REFERENCES pokemon(id),
    evolution_chain_id  INTEGER NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_encounter_player_time ON battle_encounter_log(player_id, created_at DESC);
