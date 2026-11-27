-- Lets user grant lists of an instance read grant ids ordered by user_id from the index, without sorting the joined rows.
CREATE INDEX CONCURRENTLY IF NOT EXISTS user_grants5_instance_user_id_idx ON projections.user_grants5 (instance_id, user_id, id);
