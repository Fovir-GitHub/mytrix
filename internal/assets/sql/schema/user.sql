CREATE TABLE IF NOT EXISTS user (
    id TEXT PRIMARY KEY,
    role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user'))
);
