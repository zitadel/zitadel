CREATE OR REPLACE FUNCTION eventstore.commands_to_events(commands eventstore.command[]) RETURNS SETOF eventstore.events2 AS $$
SELECT
    c.instance_id,
    c.aggregate_type,
    c.aggregate_id,
    c.command_type AS event_type,
    cs.sequence + ROW_NUMBER() OVER (PARTITION BY c.instance_id, c.aggregate_type, c.aggregate_id) AS sequence,
    c.revision,
    c.created_at,
    c.payload,
    c.creator,
    c.owner,
    c.position,
    c.in_tx_order
FROM (
    SELECT
        c.*,
        NOW() AS created_at,
        EXTRACT(EPOCH FROM NOW()) AS position,
        ROW_NUMBER() OVER () AS in_tx_order
    FROM UNNEST(commands) AS c
) AS c
JOIN (
    SELECT
        a.instance_id,
        a.aggregate_type,
        a.aggregate_id,
        a.owner,
        COALESCE(MAX(e.sequence), 0) AS sequence
    FROM (
        SELECT DISTINCT
            instance_id,
            aggregate_type,
            aggregate_id,
            owner
        FROM UNNEST(commands)
    ) AS a
    LEFT JOIN eventstore.events2 AS e
        ON a.instance_id = e.instance_id
        AND a.aggregate_type = e.aggregate_type
        AND a.aggregate_id = e.aggregate_id
        AND a.owner = e.owner
    GROUP BY
        a.instance_id,
        a.aggregate_type,
        a.aggregate_id,
        a.owner
) AS cs
    ON c.instance_id = cs.instance_id
    AND c.aggregate_type = cs.aggregate_type
    AND c.aggregate_id = cs.aggregate_id
    AND c.owner = cs.owner
ORDER BY
    c.in_tx_order;
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION eventstore.push(commands eventstore.command[]) RETURNS TABLE(
    instance_id TEXT
    , aggregate_type TEXT
    , aggregate_id TEXT
    , created_at TIMESTAMPTZ
    , "sequence" INT8
    , "position" DECIMAL
    , in_tx_order INT4
) AS $$
INSERT INTO eventstore.events2
SELECT * FROM eventstore.commands_to_events(commands)
RETURNING instance_id, aggregate_type, aggregate_id, created_at, "sequence", position, in_tx_order
$$ LANGUAGE SQL;
