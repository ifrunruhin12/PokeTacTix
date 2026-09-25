CREATE TABLE item_purchase_requests (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    request_key UUID NOT NULL,
    item_id VARCHAR(50) NOT NULL REFERENCES items(id),
    quantity INTEGER NOT NULL,
    new_quantity INTEGER,
    remaining_coins INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, request_key)
);
