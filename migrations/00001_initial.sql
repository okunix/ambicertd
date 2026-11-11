-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL CHECK (length(username) > 3),
    password_hash TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE requests (
    serial_number TEXT PRIMARY KEY, -- генерация в ambicertd
    status TEXT NOT NULL CHECK (status IN ('pending', 'issued')),
    pem TEXT NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
);

-- +goose Down
DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS users;
