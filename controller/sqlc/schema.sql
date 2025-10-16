CREATE TABLE users
(
    id              TEXT PRIMARY KEY,
    hashed_password TEXT    NOT NULL,
    created_at      INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    last_login      INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
