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

CREATE OR REPLACE FUNCTION eventstore.commands_to_events(stmt_timestamp TIMESTAMPTZ, commands eventstore.command[])
    RETURNS SETOF eventstore.events2 
    LANGUAGE 'plpgsql'
    STABLE PARALLEL SAFE
    ROWS 10
AS $$
DECLARE
    "aggregate" RECORD;
    current_sequence BIGINT;
    current_owner TEXT;
BEGIN
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
            , stmt_timestamp -- AS created_at
            , c.payload
            , c.creator
            , COALESCE(current_owner, c.owner) -- AS owner
            , EXTRACT(EPOCH FROM stmt_timestamp) -- AS position
            , c.ordinality::integer -- AS in_tx_order
            , pg_current_xact_id()::text::numeric -- AS tx_id
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
DECLARE
    instance_ids TEXT[];
    i INTEGER;
BEGIN
    -- Extract distinct instance_ids from commands
    SELECT ARRAY_AGG(DISTINCT instance_id) 
    INTO instance_ids
    FROM UNNEST(commands);

    -- Handle empty commands array
    IF instance_ids IS NULL THEN
        RETURN;
    END IF;

    -- Acquire advisory locks for each instance_id
    FOR i IN 1..array_length(instance_ids, 1) LOOP
        PERFORM pg_advisory_xact_lock_shared('eventstore.events2'::REGCLASS::OID::INTEGER, hashtext(instance_ids[i]));
    END LOOP;

    -- Perform the insert and return results
    RETURN QUERY
    INSERT INTO eventstore.events2
    SELECT * FROM eventstore.commands_to_events(STATEMENT_TIMESTAMP(), commands)
    ORDER BY in_tx_order
    RETURNING *;

    RETURN;
END;
$$ LANGUAGE plpgsql;
