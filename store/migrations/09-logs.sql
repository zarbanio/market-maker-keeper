CREATE TABLE IF NOT EXISTS logs (
    id              SERIAL PRIMARY KEY,
    cycle_id        BIGINT NOT NULL,
    level    VARCHAR(20),
    message           TEXT,
    fields          JSONB,
    created_at      TIMESTAMP default current_timestamp
)