-- +goose Up
-- 1. Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user' CHECK(role IN ('admin', 'user')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Таблица политик (Policies)
CREATE TABLE IF NOT EXISTS policies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    -- Разрешена ли автоподпись (0 - нет, требуется админ; 1 - да, подписывать сразу)
    auto_sign INTEGER NOT NULL DEFAULT 0 CHECK(auto_sign IN (0, 1)),
    -- Срок действия выпускаемого сертификата в днях
    validity_days INTEGER NOT NULL DEFAULT 365, -- ???
    -- Правила валидации CSR в формате JSON (разрешенные домены, O, OU, размер ключа и т.д.)
    rules_json TEXT NOT NULL, 
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 3. Связующая таблица Пользователи <-> Политики (Many-to-Many)
CREATE TABLE IF NOT EXISTS user_policies (
    user_id INTEGER NOT NULL,
    policy_id INTEGER NOT NULL,
    assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, policy_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE
);

-- 4. Таблица запросов на подпись (CSR)
CREATE TABLE IF NOT EXISTS csrs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    csr_pem TEXT NOT NULL, -- Сам CSR в формате PEM
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'approved', 'rejected')),
    policy_id INTEGER, -- По какой политике проверялся (опционально)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME,
    processed_by INTEGER, -- ID админа, который подписал/отклонил (NULL если автоподпись)
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE SET NULL,
    FOREIGN KEY (processed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 5. Таблица выпущенных сертификатов
CREATE TABLE IF NOT EXISTS certificates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    csr_id INTEGER UNIQUE NOT NULL, -- Ссылка на исходный CSR
    user_id INTEGER NOT NULL, -- Дублируем user_id для быстрых запросов GET /api/certs
    serial_number TEXT UNIQUE NOT NULL, -- Серийный номер серта (Hex), нужен для CRL
    common_name TEXT, -- CN из сертификата (для удобного отображения в списках)
    cert_pem TEXT NOT NULL, -- Сам сертификат в формате PEM
    not_before DATETIME NOT NULL, -- Действителен с
    not_after DATETIME NOT NULL, -- Действителен до
    -- Статусы: valid (действует), revoke_pending (запрос на отзыв), revoked (отозван)
    status TEXT NOT NULL DEFAULT 'valid' CHECK(status IN ('valid', 'revoke_pending', 'revoked')),
    revoked_at DATETIME, -- Дата отзыва (для генерации CRL)
    revocation_reason INTEGER, -- Код причины отзыва по RFC 5280 (например, 0 - unspecified, 1 - keyCompromise)
    FOREIGN KEY (csr_id) REFERENCES csrs(id) ON DELETE RESTRICT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Индексы для ускорения частых запросов
CREATE INDEX IF NOT EXISTS idx_csrs_status ON csrs(status);
CREATE INDEX IF NOT EXISTS idx_certs_user ON certificates(user_id);
CREATE INDEX IF NOT EXISTS idx_certs_status ON certificates(status);
CREATE INDEX IF NOT EXISTS idx_certs_serial ON certificates(serial_number);

-- +goose Down
DROP INDEX IF EXISTS idx_csrs_status;
DROP INDEX IF EXISTS idx_certs_user;
DROP INDEX IF EXISTS idx_certs_status;
DROP INDEX IF EXISTS idx_certs_serial;
