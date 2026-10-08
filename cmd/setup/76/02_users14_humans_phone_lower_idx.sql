CREATE INDEX CONCURRENTLY IF NOT EXISTS users14_humans_phone_lower_idx ON projections.users14_humans (instance_id, LOWER(phone));
