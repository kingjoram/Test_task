CREATE TABLE IF NOT EXISTS url (
    long TEXT NOT NULL UNIQUE,
    short TEXT NOT NULL UNIQUE
        CONSTRAINT length_constraint CHECK(LENGTH(short) = 10),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expiry_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '7 day'
)