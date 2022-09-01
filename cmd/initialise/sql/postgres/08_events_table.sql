CREATE TABLE IF NOT EXISTS eventstore.events (
	id UUID DEFAULT gen_random_uuid()
	, event_type TEXT NOT NULL
	, aggregate_type TEXT NOT NULL
	, aggregate_id TEXT NOT NULL
	, aggregate_version TEXT NOT NULL
	, event_sequence BIGINT NOT NULL
	, previous_aggregate_sequence BIGINT
	, previous_aggregate_type_sequence INT8
	, creation_date TIMESTAMPTZ NOT NULL DEFAULT now()
	, event_data JSONB
	, editor_user TEXT NOT NULL 
	, editor_service TEXT NOT NULL
	, resource_owner TEXT NOT NULL
	, instance_id TEXT NOT NULL

	, PRIMARY KEY (event_sequence, instance_id)
	, CONSTRAINT previous_sequence_unique UNIQUE(previous_aggregate_sequence, instance_id)
	, CONSTRAINT prev_agg_type_seq_unique UNIQUE(previous_aggregate_type_sequence, instance_id)
);

CREATE INDEX IF NOT EXISTS agg_type_agg_id ON eventstore.events (aggregate_type, aggregate_id, instance_id);
CREATE INDEX IF NOT EXISTS agg_type ON eventstore.events (aggregate_type, instance_id);
CREATE INDEX IF NOT EXISTS agg_type_seq ON eventstore.events (aggregate_type, event_sequence DESC, instance_id);
CREATE INDEX IF NOT EXISTS max_sequence ON eventstore.events (aggregate_type, aggregate_id, event_sequence DESC, instance_id);
