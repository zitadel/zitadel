SET experimental_enable_hash_sharded_indexes = on;

CREATE TABLE IF NOT EXISTS eventstore.events (
	id UUID DEFAULT gen_random_uuid()
	, event_type TEXT NOT NULL
	, aggregate_type TEXT NOT NULL
	, aggregate_id TEXT NOT NULL
	, aggregate_version TEXT NOT NULL
	, creation_date TIMESTAMPTZ NOT NULL DEFAULT statement_timestamp()
	, event_data JSONB
	, editor_user TEXT NOT NULL 
	, editor_service TEXT NOT NULL
	, resource_owner TEXT NOT NULL
	, instance_id TEXT NOT NULL
	, previous_event_date TIMESTAMPTZ

	, PRIMARY KEY (creation_date DESC, instance_id) USING HASH WITH BUCKET_COUNT = 10
	, INDEX push_ro (aggregate_type, aggregate_id, instance_id) STORING (resource_owner)
	-- TODO: could be hashed for better sharding, focus on events returned
	-- , INDEX scheduled (creation_date, instance_id, aggregate_type) 
    -- 	STORING (id, event_type, aggregate_id, aggregate_version, event_data, editor_user, editor_service, resource_owner, previous_event_date)
	-- TODO: easier to shard by default, less focus on amount of event returned
	, INDEX scheduled4 (aggregate_type, instance_id, creation_date DESC) 
    	STORING (id, event_type, aggregate_id, aggregate_version, event_data, editor_user, editor_service, resource_owner, previous_event_date)



	, CONSTRAINT previous_event_unique UNIQUE (previous_event_date DESC, aggregate_type, aggregate_id, instance_id)
);