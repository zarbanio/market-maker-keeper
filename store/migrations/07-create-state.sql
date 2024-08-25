CREATE TABLE IF NOT EXISTS bot_state (
    id              SERIAL PRIMARY KEY,
    strategy_name      VARCHAR(255) NOT NULL,
    dx_id           INTEGER NULL,
    cx_id           INTEGER NULL
);