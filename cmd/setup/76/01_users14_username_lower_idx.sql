CREATE INDEX CONCURRENTLY IF NOT EXISTS users14_username_lower_idx ON projections.users14 (instance_id, LOWER(username));
