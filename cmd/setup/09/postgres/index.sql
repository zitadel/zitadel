-- replace agg_type_agg_id
BEGIN;
DROP INDEX IF EXISTS eventstore.agg_type_agg_id;
COMMIT;

BEGIN;
CREATE INDEX agg_type_agg_id ON eventstore.events (
    instance_id
    , aggregate_type
    , aggregate_id
);
COMMIT;

-- replace agg_type
BEGIN;
DROP INDEX IF EXISTS eventstore.agg_type;
COMMIT;

BEGIN;
CREATE INDEX agg_type ON eventstore.events (
    instance_id
    , aggregate_type
    , event_sequence
);
COMMIT;

-- drop unused index
BEGIN;
DROP INDEX IF EXISTS eventstore.agg_type_seq;
COMMIT;
