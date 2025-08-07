CREATE INDEX CONCURRENTLY IF NOT EXISTS es_wm_temp
    ON eventstore.events2 (instance_id, aggregate_id, aggregate_type, event_type, "position");
