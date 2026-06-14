CREATE TABLE IF NOT EXISTS room (
    id TEXT PRIMARY KEY,
    state TEXT NOT NULL CHECK (state IN ('left', 'joined'))
);
