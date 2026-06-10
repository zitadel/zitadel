-- represents an event to be created.
DO $$ BEGIN
    CREATE TYPE eventstore.command2 AS (
        instance_id TEXT
        , aggregate_type TEXT
        , aggregate_id TEXT
        , command_type TEXT
        , revision INT2
        , payload JSONB
        , creator TEXT
        , owner TEXT
        , enforce_owner BOOLEAN
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE eventstore.latest_command AS (
        instance_id TEXT
        , aggregate_type TEXT
        , aggregate_id TEXT
        , owner TEXT
        , sequence BIGINT
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE OR REPLACE FUNCTION eventstore.commands_to_events(commands eventstore.command2[]) 
    RETURNS SETOF eventstore.events2 
    LANGUAGE 'plpgsql'
    STABLE PARALLEL SAFE
    ROWS 10
AS $$
DECLARE
    _latest_events eventstore.latest_command[];
    command eventstore.command2;
    last_owner TEXT;
    i INTEGER;
BEGIN

    SELECT array_agg(res) INTO _latest_events
    FROM (
        SELECT
            ROW(
                c.instance_id,
                c.aggregate_type,
                c.aggregate_id,
                COALESCE(e.owner, c.owner),
                COALESCE(e.sequence, 0)
            )::eventstore.latest_command AS res
        FROM (
            SELECT DISTINCT ON (c.instance_id, c.aggregate_type, c.aggregate_id)
                c.instance_id,
                c.aggregate_type,
                c.aggregate_id,
                c.owner,
                c.ordinality
            FROM (
                SELECT *
                FROM UNNEST(commands) WITH ORDINALITY AS c
            ) AS c
            ORDER BY
                c.instance_id,
                c.aggregate_type,
                c.aggregate_id,
                c.ordinality
        ) AS c
        LEFT JOIN LATERAL (
            SELECT
                e.owner,
                e.sequence
            FROM eventstore.events2 AS e
            WHERE e.instance_id = c.instance_id
                AND e.aggregate_type = c.aggregate_type
                AND e.aggregate_id = c.aggregate_id
            ORDER BY e.sequence DESC
            LIMIT 1
        ) AS e ON TRUE
    ) AS latest;

    FOR command IN SELECT * FROM UNNEST(commands) LOOP
        IF NOT command.enforce_owner THEN
            SELECT owner INTO command.owner
            FROM UNNEST(_latest_events) AS e
            WHERE e.instance_id = command.instance_id
            AND e.aggregate_type = command.aggregate_type
            AND e.aggregate_id = command.aggregate_id;
        ELSE
            FOR i IN 1..array_length(_latest_events, 1) LOOP
                IF _latest_events[i].instance_id = command.instance_id
                    AND _latest_events[i].aggregate_type = command.aggregate_type
                    AND _latest_events[i].aggregate_id = command.aggregate_id
                THEN
                    _latest_events[i].owner := command.owner;
                    EXIT;
                END IF;
            END LOOP;
        END IF;

    END LOOP;


    RETURN QUERY SELECT
        c.instance_id
        , c.aggregate_type
        , c.aggregate_id
        , c.command_type AS event_type
        , COALESCE(e.sequence, 0) + ROW_NUMBER() OVER (PARTITION BY c.instance_id, c.aggregate_type, c.aggregate_id ORDER BY c.ordinality) AS sequence
        , c.revision
        , NOW() AS created_at
        , c.payload
        , c.creator
        , CASE WHEN c.enforce_owner
            THEN c.owner
            ELSE COALESCE(e.owner, c.owner)
        END AS owner
        , EXTRACT(EPOCH FROM NOW()) AS position
        , c.ordinality::%s -- AS in_tx_order
    FROM UNNEST(commands) WITH ORDINALITY AS c
    LEFT JOIN unnest(latest_events) AS e
        ON c.instance_id = e.instance_id
        AND c.aggregate_type = e.aggregate_type
        AND c.aggregate_id = e.aggregate_id
    ORDER BY
        in_tx_order;
END;
$$;

CREATE OR REPLACE FUNCTION eventstore.push(commands eventstore.command2[]) RETURNS SETOF eventstore.events2 VOLATILE AS $$
INSERT INTO eventstore.events2
SELECT * FROM eventstore.commands_to_events(commands)
ORDER BY in_tx_order
RETURNING *
$$ LANGUAGE SQL;