-- +goose Up

-- Все даты хранятся как TEXT в ISO-8601 UTC: 2026-11-20T12:34:56.789Z

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user' CHECK(role IN ('admin', 'user')),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;

CREATE TABLE IF NOT EXISTS policies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    auto_sign INTEGER NOT NULL DEFAULT 0 CHECK(auto_sign IN (0, 1)),
    subject_cn_regex TEXT, -- например '^[a-z0-9-]+\.corp\.example$'
    san_dns_regex TEXT, -- каждое DNS-имя в SAN должно подходить
    max_requested_days INTEGER CHECK (max_requested_days IS NULL OR max_requested_days > 0),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')) 
) STRICT;

CREATE TABLE IF NOT EXISTS user_policies (
    user_id INTEGER NOT NULL,
    policy_id INTEGER NOT NULL,
    assigned_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')), 
    PRIMARY KEY (user_id, policy_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE
) STRICT;

CREATE TABLE IF NOT EXISTS csrs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    csr_pem TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'approved', 'rejected')),
    requested_days INTEGER CHECK (requested_days IS NULL OR requested_days > 0),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    processed_at TEXT,  -- date iso8601
    processed_by INTEGER, -- ID админа, который подписал/отклонил (NULL если автоподпись)
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (processed_by) REFERENCES users(id) ON DELETE SET NULL
) STRICT;

CREATE TABLE IF NOT EXISTS certificates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    csr_id INTEGER UNIQUE NOT NULL,
    serial_number TEXT UNIQUE NOT NULL, -- нужен для CRL
    common_name TEXT, -- CN из сертификата (для удобного отображения в списках)
    cert_pem TEXT NOT NULL, -- Сам сертификат в формате PEM
    not_before TEXT NOT NULL, -- iso-8601 Действителен с
    not_after TEXT NOT NULL, -- iso-8601 Действителен до
    -- Статусы: valid (действует), revoke_pending (запрос на отзыв), revoked (отозван)
    status TEXT NOT NULL DEFAULT 'valid' CHECK(status IN ('valid', 'revoke_pending', 'revoked')),
    revoked_at TEXT, -- iso-8601 
    revocation_reason INTEGER, -- Код причины отзыва по RFC 5280 (например, 0 - unspecified, 1 - keyCompromise)
    FOREIGN KEY (csr_id) REFERENCES csrs(id) ON DELETE RESTRICT
) STRICT;

CREATE INDEX IF NOT EXISTS idx_csrs_status ON csrs(status);
CREATE INDEX IF NOT EXISTS idx_certs_status ON certificates(status);
CREATE INDEX IF NOT EXISTS idx_certs_serial ON certificates(serial_number);

-- +goose Down
DROP INDEX IF EXISTS idx_csrs_status;
DROP INDEX IF EXISTS idx_certs_user;
DROP INDEX IF EXISTS idx_certs_status;
DROP INDEX IF EXISTS idx_certs_serial;

DROP TABLE IF EXISTS certificates;
DROP TABLE IF EXISTS csrs;
DROP TABLE IF EXISTS user_policies;
DROP TABLE IF EXISTS policies;
DROP TABLE IF EXISTS users;
