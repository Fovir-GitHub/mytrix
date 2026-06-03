-- name: SelectFeedByID :one
SELECT
    *
FROM
    rss_feed
WHERE
    id = ?
LIMIT 1;

-- name: AllFeds :many
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
