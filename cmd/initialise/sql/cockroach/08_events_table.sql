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
	, INDEX es_handler_idx (instance_id, aggregate_type, event_type)
	, INDEX es_agg_wm (instance_id, aggregate_type, aggregate_id, event_type, "position")
	, INDEX es_wm (instance_id, "owner", aggregate_type, aggregate_id, event_type)
	, INDEX es_active_instances (created_at DESC, instance_id) USING HASH
);