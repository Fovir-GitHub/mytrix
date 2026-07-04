ALTER TABLE rss_item
    ADD COLUMN pushed INTEGER NOT NULL DEFAULT 0 CHECK (pushed IN (0, 1));

UPDATE
    rss_item
SET
    pushed = 1;
