-- Mark legacy evolution chains (multi-member chains with empty links) for
-- cold refresh by setting their fetched_at to epoch. The next time any Pokemon
-- in these families triggers an evolution check, loadEvolutionChain will detect
-- the stale data (via usable() returning false), attempt a PokéAPI fetch, and
-- store the proper links. This self-healing approach avoids hammering PokéAPI
-- at migration time for chains that may never be encountered in actual gameplay.

UPDATE evolution_chain
SET fetched_at = '1970-01-01 00:00:00+00'::timestamptz
WHERE links = '[]'::jsonb
  AND array_length(member_species_ids, 1) > 1;
