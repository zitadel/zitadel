CREATE INDEX CONCURRENTLY IF NOT EXISTS es_instance_position ON eventstore.events2 (instance_id, position);
