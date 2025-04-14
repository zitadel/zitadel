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
    CREATE TYPE eventstore.command AS (
        instance_id TEXT
        , aggregate_type TEXT
        , aggregate_id TEXT
        , command_type TEXT
        , revision INT2
        , payload JSONB
        , creator TEXT
        , owner TEXT
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE OR REPLACE FUNCTION eventstore.commands_to_events(commands eventstore.command[]) RETURNS SETOF eventstore.events2 VOLATILE AS $$
SELECT
    c.instance_id
    , c.aggregate_type
    , c.aggregate_id
    , c.command_type AS event_type
    , cs.sequence + ROW_NUMBER() OVER (PARTITION BY c.instance_id, c.aggregate_type, c.aggregate_id ORDER BY c.in_tx_order) AS sequence
    , c.revision
    , NOW() AS created_at
    , c.payload
    , c.creator
    , cs.owner
    , EXTRACT(EPOCH FROM NOW()) AS position
    , c.in_tx_order
FROM (
    SELECT
        c.instance_id
        , c.aggregate_type
        , c.aggregate_id
        , c.command_type
        , c.revision
        , c.payload
        , c.creator
        , c.owner
        , ROW_NUMBER() OVER () AS in_tx_order
    FROM
        UNNEST(commands) AS c
) AS c
JOIN (
    SELECT
        cmds.instance_id
        , cmds.aggregate_type
        , cmds.aggregate_id
        , CASE WHEN (e.owner IS NOT NULL OR e.owner <> '') THEN e.owner ELSE command_owners.owner END AS owner
        , COALESCE(MAX(e.sequence), 0) AS sequence
    FROM (
        SELECT DISTINCT
            instance_id
            , aggregate_type
            , aggregate_id
            , owner
        FROM UNNEST(commands)
    ) AS cmds
    LEFT JOIN eventstore.events2 AS e
        ON cmds.instance_id = e.instance_id
        AND cmds.aggregate_type = e.aggregate_type
        AND cmds.aggregate_id = e.aggregate_id
    JOIN (
        SELECT
            DISTINCT ON (
                instance_id
                , aggregate_type
                , aggregate_id
            )
            instance_id
            , aggregate_type
            , aggregate_id
            , owner
        FROM
            UNNEST(commands)
    ) AS command_owners ON
        cmds.instance_id = command_owners.instance_id
        AND cmds.aggregate_type = command_owners.aggregate_type
        AND cmds.aggregate_id = command_owners.aggregate_id
    GROUP BY
        cmds.instance_id
        , cmds.aggregate_type
        , cmds.aggregate_id
        , 4 -- owner
) AS cs
    ON c.instance_id = cs.instance_id
    AND c.aggregate_type = cs.aggregate_type
    AND c.aggregate_id = cs.aggregate_id
ORDER BY
    in_tx_order;
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION eventstore.push(commands eventstore.command[]) RETURNS SETOF eventstore.events2 VOLATILE AS $$
INSERT INTO eventstore.events2
SELECT * FROM eventstore.commands_to_events(commands)
RETURNING *
$$ LANGUAGE SQL;
