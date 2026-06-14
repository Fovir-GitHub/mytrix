-- name: SelectUserByID :one
SELECT
    *
FROM
    user
WHERE
    id = ?
LIMIT 1;

-- name: IsUserAdmin :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            user
        WHERE
            id = ?
            AND role = 'admin');

-- name: CreateUser :exec
INSERT
    OR IGNORE INTO user (id, role)
        VALUES (?, ?);

-- name: UpdateUserRoleByID :exec
UPDATE
    user
SET
    role = ?
WHERE
    id = ?;
