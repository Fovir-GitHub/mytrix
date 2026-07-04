-- name: CreateRSSItem :one
INSERT
    OR IGNORE INTO rss_item (feed_id, guid, link, title, description)
        VALUES (?, ?, ?, ?, ?)
    RETURNING
        id;

-- name: DeleteRSSItemByFeedID :exec
DELETE FROM rss_item
WHERE feed_id = ?;

-- name: MarkRSSItemPushedByID :exec
UPDATE
    rss_item
SET
    pushed = 1
WHERE
    id = ?;
