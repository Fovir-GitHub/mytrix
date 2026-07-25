-- name: IsRoomJoined :one
SELECT
    state = 'accept'
FROM
    room
WHERE
    id = ?;

-- name: IsRoomLeft :one
SELECT
    state = 'leave'
FROM
    room
WHERE
    id = ?;

-- name: IsRoomExists :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            room
        WHERE
            id = ?);

-- name: CreateRoom :exec
INSERT
    OR IGNORE INTO room (id, state)
        VALUES (?, ?);

-- name: UpdateRoomState :exec
UPDATE
    room
SET
    state = ?
WHERE
    id = ?;

-- name: CreateEvent :one
INSERT
    OR IGNORE INTO event (event_id, room_id)
        VALUES (?, ?)
    RETURNING
        id;

-- name: GetEventByRoomIDAndEventID :one
SELECT
    *
FROM
    event
WHERE
    event_id = ?
    AND room_id = ?;

-- name: UpdateEventByID :exec
UPDATE
    event
SET
    room_id = ?,
    event_id = ?
WHERE
    id = ?;
