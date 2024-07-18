CREATE TABLE IF NOT EXISTS users (
    'id' INTEGER PRIMARY KEY,
    'login' VARCHAR(50) NOT NULL,
    'pass_hash' VARCHAR(100) NOT NULL,
    'email' VARCHAR(50) NOT NULL
);
CREATE INDEX login_idx ON users ('login');