DELETE FROM active_boosts;
DROP TABLE IF EXISTS active_boosts;
DELETE FROM items WHERE id IN ('attack-booster', 'hp-booster');
