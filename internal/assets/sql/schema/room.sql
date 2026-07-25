CREATE TABLE IF NOT EXISTS room (
    id TEXT PRIMARY KEY,
    state TEXT NOT NULL CHECK (state IN ('left', 'joined'))
);

CREATE TABLE IF NOT EXISTS event (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id TEXT NOT NULL,
    room_id TEXT NOT NULL REFERENCES room (id),
    UNIQUE (event_id, room_id)
);
