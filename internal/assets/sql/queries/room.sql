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
