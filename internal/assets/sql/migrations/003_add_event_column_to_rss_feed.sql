ALTER TABLE rss_feed
    ADD COLUMN event_id INTEGER REFERENCES event (id);
