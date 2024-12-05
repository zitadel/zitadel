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
