CREATE INDEX CONCURRENTLY IF NOT EXISTS users14_instance_resource_owner_idx ON projections.users14 (instance_id, resource_owner, id) INCLUDE (type);
