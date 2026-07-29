CREATE TABLE IF NOT EXISTS evolution_chain (
    id                  INTEGER PRIMARY KEY,
    member_species_ids  INTEGER[] NOT NULL,
    fetched_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
