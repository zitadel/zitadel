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

CREATE OR REPLACE FUNCTION eventstore.commands_to_events(commands eventstore.command2[]) 
    RETURNS SETOF eventstore.events2 
    LANGUAGE 'plpgsql'
    STABLE PARALLEL SAFE
    ROWS 10
AS $$
DECLARE
    latest_events eventstore.events2[];
BEGIN
    SELECT array_agg(e) INTO latest_events
    FROM (
        SELECT DISTINCT ON (e.instance_id, e.aggregate_type, e.aggregate_id) e
        FROM (
            SELECT
                c.instance_id
                , c.aggregate_type
                , c.aggregate_id
            FROM UNNEST(commands) AS c
            GROUP BY
                c.instance_id
                , c.aggregate_type
                , c.aggregate_id
        ) AS cmds
        LEFT JOIN eventstore.events2 AS e
            ON cmds.instance_id = e.instance_id
            AND cmds.aggregate_type = e.aggregate_type
            AND cmds.aggregate_id = e.aggregate_id
        ORDER BY
            e.instance_id
            , e.aggregate_type
            , e.aggregate_id
            , e.sequence DESC NULLS LAST
    );


    RETURN QUERY SELECT
        c.instance_id
        , c.aggregate_type
        , c.aggregate_id
        , c.command_type AS event_type
        , COALESCE(e.sequence, 0) + ROW_NUMBER() OVER (PARTITION BY c.instance_id, c.aggregate_type, c.aggregate_id ORDER BY c.in_tx_order) AS sequence
        , c.revision
        , NOW() AS created_at
        , c.payload
        , c.creator
        , CASE WHEN c.enforce_owner
            THEN c.owner
            ELSE COALESCE(e.owner, c.owner)
        END AS owner
        , EXTRACT(EPOCH FROM NOW()) AS position
        , c.in_tx_order
    FROM (
        SELECT c.*, CAST(ROW_NUMBER() OVER () AS INTEGER) AS in_tx_order
        FROM UNNEST(commands) AS c
    ) AS c
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
