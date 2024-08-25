CREATE TABLE IF NOT EXISTS trades
(
    id             SERIAL PRIMARY KEY,
    pair_id        INTEGER REFERENCES pairs (id),
    order_id       INTEGER REFERENCES orders (id),
    transaction_id INTEGER REFERENCES transactions (id)
);