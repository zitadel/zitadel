CREATE OR REPLACE FUNCTION eventstore.latest_aggregate_state(
    instance_id TEXT
    , aggregate_type TEXT
    , aggregate_id TEXT
    
    , sequence OUT BIGINT
    , owner OUT TEXT
)
    LANGUAGE 'plpgsql'
    STABLE PARALLEL SAFE
AS $$
    BEGIN
        SELECT
            COALESCE(e.sequence, 0) AS sequence
            , e.owner
        INTO
            sequence
            , owner
        FROM
            eventstore.events2 e
        WHERE
            e.instance_id = $1
            AND e.aggregate_type = $2
            AND e.aggregate_id = $3
        ORDER BY 
            e.sequence DESC
        LIMIT 1;

        RETURN;
    END;
$$;

CREATE OR REPLACE FUNCTION eventstore.commands_to_events(commands eventstore.command[])
    RETURNS SETOF eventstore.events2 
    LANGUAGE 'plpgsql'
    STABLE PARALLEL SAFE
    ROWS 10
AS $$
DECLARE
    "aggregate" RECORD;
    current_sequence BIGINT;
    current_owner TEXT;
    created_at TIMESTAMPTZ;
BEGIN
    created_at := statement_timestamp();
    FOR "aggregate" IN 
        SELECT DISTINCT
            instance_id
            , aggregate_type
            , aggregate_id
        FROM UNNEST(commands)
    LOOP
        SELECT 
            * 
        INTO
            current_sequence
            , current_owner 
        FROM eventstore.latest_aggregate_state(
            "aggregate".instance_id
            , "aggregate".aggregate_type
            , "aggregate".aggregate_id
        );

        RETURN QUERY
        SELECT
            c.instance_id
            , c.aggregate_type
            , c.aggregate_id
            , c.command_type -- AS event_type
            , COALESCE(current_sequence, 0) + ROW_NUMBER() OVER () -- AS sequence
            , c.revision
            , created_at
            , c.payload
            , c.creator
            , COALESCE(current_owner, c.owner) -- AS owner
            , EXTRACT(EPOCH FROM created_at) -- AS position
            , c.ordinality::{{ .InTxOrderType }} -- AS in_tx_order
        FROM
            UNNEST(commands) WITH ORDINALITY AS c
        WHERE
            c.instance_id = aggregate.instance_id
            AND c.aggregate_type = aggregate.aggregate_type
            AND c.aggregate_id = aggregate.aggregate_id;
    END LOOP;
    RETURN;
END;
$$;

CREATE OR REPLACE FUNCTION eventstore.push(commands eventstore.command[]) RETURNS SETOF eventstore.events2 VOLATILE AS $$
INSERT INTO eventstore.events2
SELECT * FROM eventstore.commands_to_events(commands)
ORDER BY in_tx_order
RETURNING *
$$ LANGUAGE SQL;
