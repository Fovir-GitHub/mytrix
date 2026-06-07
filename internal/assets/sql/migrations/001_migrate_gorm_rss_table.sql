CREATE TABLE IF NOT EXISTS rss_feeds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    url TEXT NOT NULL,
    title TEXT
);

CREATE TABLE IF NOT EXISTS rss_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    feed_id INTEGER,
    guid TEXT,
    link TEXT,
    title TEXT,
    description TEXT
);

INSERT
    OR IGNORE INTO rss_feed (id, url, title)
    SELECT
        id,
        url,
        title
    FROM
        rss_feeds
    WHERE
        url != ""
        AND deleted_at IS NULL;

INSERT
    OR IGNORE INTO rss_item (id, feed_id, guid, link, title, description)
    SELECT
        i.id,
        i.feed_id,
        i.guid,
        i.link,
        i.title,
        i.description
    FROM
        rss_items AS i
        JOIN rss_feeds AS f ON i.feed_id = f.id
    WHERE
        f.deleted_at IS NULL;

DROP TABLE IF EXISTS rss_feeds;

DROP TABLE IF EXISTS rss_items;
