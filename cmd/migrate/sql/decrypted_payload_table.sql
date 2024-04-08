CREATE TEMP TABLE eventstore.decrypted_payload (
    instance_id TEXT REFERENCES eventstore.events2 (instance_id),
    aggregate_type TEXT REFERENCES eventstore.events2 (aggregate_type),
    aggregate_id TEXT REFERENCES eventstore.events2 (aggregate_id),
    "sequence" BIGINT REFERENCES eventstore.events2 ("sequence"),

    decrypted_payload JSONB,
)