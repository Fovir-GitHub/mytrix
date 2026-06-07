-- name: CreateRSSItem :exec
INSERT INTO rss_item (feed_id, guid, link, title, description)
    VALUES (?, ?, ?, ?, ?);

-- name: DeleteRSSItemByFeedID :exec
DELETE FROM rss_item
WHERE feed_id = ?;
