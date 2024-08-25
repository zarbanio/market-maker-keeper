CREATE TABLE IF NOT EXISTS arbitrage_opportunity (
    id SERIAL PRIMARY KEY,
    uuid BIGINT NOT NULL,
    opportunity JSON,
    created_at      TIMESTAMP default current_timestamp
);