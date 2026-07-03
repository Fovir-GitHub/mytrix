-- name: CreateRSSItem :exec
INSERT
    OR IGNORE INTO rss_item (feed_id, guid, link, title, description)
        VALUES (?, ?, ?, ?, ?);

-- name: DeleteRSSItemByFeedID :exec
DELETE FROM rss_item
WHERE feed_id = ?;
