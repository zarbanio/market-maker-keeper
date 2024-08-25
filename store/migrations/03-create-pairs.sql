CREATE TABLE IF NOT EXISTS pairs
(
    id             SERIAL PRIMARY KEY,
    base_asset     VARCHAR(255) NOT NULL,
    quote_asset    VARCHAR(255) NOT NULL
);