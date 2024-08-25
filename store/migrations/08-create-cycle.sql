CREATE TABLE IF NOT EXISTS cycles (
    id serial NOT NULL,
    "start"     TIMESTAMP NULL,
    "end"       TIMESTAMP NULL,
    status      TEXT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE ("start", "end"),
    CHECK ("start" < "end")
)