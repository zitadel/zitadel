CREATE TABLE IF NOT EXISTS eventstore.events2 (
    instance_id TEXT NOT NULL
    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    
    , event_type TEXT NOT NULL
    , "sequence" BIGINT NOT NULL
    , revision SMALLINT NOT NULL
    , created_at TIMESTAMPTZ NOT NULL
    , payload JSONB
    , creator TEXT NOT NULL
    , "owner" TEXT NOT NULL
    
    , "position" DECIMAL NOT NULL
    , in_tx_order INTEGER NOT NULL

    , PRIMARY KEY (instance_id, aggregate_type, aggregate_id, "sequence")
);

CREATE INDEX IF NOT EXISTS es_active_instances ON eventstore.events2 (created_at DESC, instance_id);
CREATE INDEX IF NOT EXISTS es_wm ON eventstore.events2 (aggregate_id, instance_id, aggregate_type, event_type);
CREATE INDEX IF NOT EXISTS es_projection ON eventstore.events2 (instance_id, aggregate_type, event_type, "position");