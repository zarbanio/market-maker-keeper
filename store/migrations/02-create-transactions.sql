CREATE TABLE IF NOT EXISTS transactions(
    id                  SERIAL PRIMARY KEY,
    tx_hash             VARCHAR(255) NOT NULL UNIQUE,
    from_address        VARCHAR(255) NOT NULL,
    timestamp           TIMESTAMP,
    block_number        INTEGER,
    to_address          VARCHAR(255) NOT NULL,
    value               DECIMAL(80, 0) NOT NULL,
    gas_price           INTEGER NOT NULL,
    gas_usage           INTEGER,
    transaction_status  VARCHAR(20) NOT NULL CHECK (transaction_status IN ('pending', 'success', 'failed')),
    transaction_data    JSONB
);