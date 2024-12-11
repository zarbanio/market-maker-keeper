CREATE TABLE if not exists orders (
    id SERIAL PRIMARY KEY,
    order_id INTEGER,
    execution VARCHAR(20) NOT NULL CHECK (
        execution IN ('market', 'limit', 'stop_market', 'stop_limit')
    ),
    side VARCHAR(20) NOT NULL CHECK (side IN ('buy', 'sell')),
    srcCurrency VARCHAR(20) NOT NULL,
    dstCurrency VARCHAR(20) NOT NULL,
    price VARCHAR(50) NOT NULL,
    amount VARCHAR(50) NOT NULL,
    totalPrice VARCHAR(50) NOT NULL,
    totalOrderPrice VARCHAR(50) NOT NULL,
    stopPrice VARCHAR(50) NOT NULL,
    matchedAmount NUMERIC(28, 8) NOT NULL,
    unmatchedAmount NUMERIC(28, 8) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (
        status IN (
            'open',
            'filled',
            'partially_filled',
            'canceled',
            'draft'
        )
    ),
    partial BOOLEAN,
    fee VARCHAR(50) NOT NULL,
    feeCurrency VARCHAR(50) NOT NULL,
    account VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP NOT NULL,
    CONSTRAINT orders_currency_check CHECK (srcCurrency <> dstCurrency)
);
