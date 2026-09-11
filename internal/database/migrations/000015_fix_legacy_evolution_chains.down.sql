-- Revert the fetched_at timestamp changes.
-- This rollback can't perfectly restore the original timestamps, so it sets
-- them to now() instead. The chains remain usable; they just won't trigger
-- an immediate refresh if the migration is rolled back.

UPDATE evolution_chain
SET fetched_at = now()
WHERE fetched_at = '1970-01-01 00:00:00+00'::timestamptz;
