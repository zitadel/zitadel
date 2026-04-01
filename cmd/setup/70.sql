DO $$ BEGIN
    ALTER TYPE eventstore.command ADD ATTRIBUTE enforce_owner BOOLEAN;
EXCEPTION
    WHEN duplicate_column THEN null;
END $$;

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
    command_owner TEXT;
    enforced_owner TEXT;
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

        SELECT
            c.owner
        INTO
            command_owner
        FROM
            UNNEST(commands) WITH ORDINALITY AS c
        WHERE
            c.instance_id = "aggregate".instance_id
            AND c.aggregate_type = "aggregate".aggregate_type
            AND c.aggregate_id = "aggregate".aggregate_id
        ORDER BY
            c.ordinality
        LIMIT 1;

        SELECT
            c.owner
        INTO
            enforced_owner
        FROM
            UNNEST(commands) WITH ORDINALITY AS c
        WHERE
            c.instance_id = "aggregate".instance_id
            AND c.aggregate_type = "aggregate".aggregate_type
            AND c.aggregate_id = "aggregate".aggregate_id
            AND c.enforce_owner
        ORDER BY
            c.ordinality
        LIMIT 1;

        RETURN QUERY
        SELECT
            c.instance_id
            , c.aggregate_type
            , c.aggregate_id
            , c.command_type
            , COALESCE(current_sequence, 0) + ROW_NUMBER() OVER (ORDER BY c.ordinality)
            , c.revision
            , created_at
            , c.payload
            , c.creator
            , COALESCE(enforced_owner, current_owner, command_owner)
            , EXTRACT(EPOCH FROM created_at)
            , c.ordinality::%s
        FROM
            UNNEST(commands) WITH ORDINALITY AS c
        WHERE
            c.instance_id = "aggregate".instance_id
            AND c.aggregate_type = "aggregate".aggregate_type
            AND c.aggregate_id = "aggregate".aggregate_id
        ORDER BY
            c.ordinality;
    END LOOP;
    RETURN;
END;
$$;
