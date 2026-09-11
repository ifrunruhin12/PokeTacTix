# Evolution Chain Legacy Data Fix

## Problem Summary

Legacy evolution chains were stored in the database with:
- `member_species_ids` populated (e.g., `[4, 5, 6]` for Charmander family)
- `links = []` (empty array - no evolution edge information)

This caused evolution to silently fail for 228+ Pokemon families including:
- Bulbasaur → Ivysaur → Venusaur
- Charmander → Charmeleon → Charizard
- Squirtle → Wartortle → Blastoise
- And 225+ more families

## How Evolution Failed

1. Player's Pokemon gains XP and levels up ✅
2. Evolution check runs
3. `loadEvolutionChain()` retrieves chain with empty `links = []`
4. `PickEvolutionLink()` receives empty array, returns `nil`
5. No evolution happens (silently)
6. Player's Charmander stays Charmander forever at level 16+

## The Fix (Lazy Self-Healing)

**Migration Applied**: `000015_fix_legacy_evolution_chains.up.sql`

```sql
UPDATE evolution_chain
SET fetched_at = '1970-01-01 00:00:00+00'::timestamptz
WHERE links = '[]'::jsonb
  AND array_length(member_species_ids, 1) > 1;
```

**What This Does**:
- Marks 228 legacy chains with a very old `fetched_at` timestamp
- The existing self-healing code in `loadEvolutionChain()` already detects these as stale
- When evolution is checked for these families, the cold-refresh path triggers
- PokéAPI is queried and proper `links` are stored
- Future evolution checks work correctly

## Migration Results

- **Total chains updated**: 228
- **Chains already fixed**: 3 (Eevee, Seedot, Axew)
- **Chains still needing refresh**: 228

## What Happens Next

The next time any player:
1. Battles with a Pokemon from an affected family
2. That Pokemon levels up
3. Evolution check runs

The system will:
1. Detect the chain has empty links (via `usable()` check)
2. Fetch the proper evolution data from PokéAPI
3. Store it with correct `links` array
4. Apply the evolution if level threshold is met
5. All future evolution checks for that family work correctly

## Verification

Check chains that have been refreshed:
```sql
SELECT ec.id, ec.member_species_ids, p.name as pokemon_name, 
       jsonb_array_length(ec.links) as link_count, ec.fetched_at
FROM evolution_chain ec
LEFT JOIN pokemon p ON p.species_id = ec.member_species_ids[1]
WHERE ec.links != '[]'::jsonb 
  AND array_length(ec.member_species_ids, 1) > 1
ORDER BY ec.fetched_at DESC;
```

Check chains still awaiting refresh:
```sql
SELECT COUNT(*) as legacy_chains
FROM evolution_chain
WHERE links = '[]'::jsonb
  AND array_length(member_species_ids, 1) > 1;
```

## Trade-offs

**Pros**:
- No mass PokéAPI hammering at migration time
- Self-healing happens naturally during gameplay
- Only fetches data for families actually encountered
- Respects rate limits and backoff logic

**Cons**:
- Evolution won't work until first post-fix battle for each family
- Players who already experienced the bug won't get retroactive evolution
- Some rarely-encountered families may take longer to fix

## Alternative (Not Implemented)

An immediate fix script could fetch all 228 chains from PokéAPI right now:
```bash
# Not run - would require ~228 API calls
go run cmd/fix-evolution-chains/main.go
```

Decided against this to avoid:
- Rate limiting issues with PokéAPI
- Fetching data for Pokemon that may never be encountered
- Migration downtime

## Date Applied

- **Migration Created**: 2026-09-10
- **Migration Applied**: 2026-09-10
- **Chains Fixed**: 228 marked for lazy refresh

## Related Code

- `internal/pokemon/service.go` - `loadEvolutionChain()` contains self-healing logic
- `internal/battle/evolution.go` - `resolveEvolutionTargets()` triggers evolution checks
- `internal/battle/rewards.go` - `ApplyAllRewards()` applies evolution after battles
