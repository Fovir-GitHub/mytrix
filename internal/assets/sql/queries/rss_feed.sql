-- name: SelectFeedByID :one
SELECT
    *
FROM
    rss_feed
WHERE
    id = ?
LIMIT 1;

-- name: AllFeeds :many
SELECT
    *
FROM
    rss_feed
ORDER BY
    id;

-- name: CreateRSSFeed :exec
INSERT INTO rss_feed (url, title)
    VALUES (?, ?);

-- name: DeleteRSSFeed :exec
DELETE FROM rss_feed
WHERE id = ?;

-- name: GetFeedThreadRootEventByID :one
SELECT
    e.id,
    e.event_id,
    e.room_id
FROM
    event AS e
    INNER JOIN rss_feed AS f ON e.id = f.event_id
WHERE
    f.id = ?;

-- name: UpdateFeedThreadRootEventByID :exec
UPDATE
    rss_feed
SET
    event_id = ?
WHERE
    id = ?;
