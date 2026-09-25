CREATE TABLE battle_boost_reservations (
    battle_id VARCHAR(255) PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    finalized BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE battle_boost_reservation_items (
    battle_id VARCHAR(255) NOT NULL REFERENCES battle_boost_reservations(battle_id) ON DELETE CASCADE,
    boost_id BIGINT NOT NULL,
    item_id VARCHAR(50) NOT NULL REFERENCES items(id),
    stat VARCHAR(20) NOT NULL,
    bonus INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (battle_id, boost_id)
);
