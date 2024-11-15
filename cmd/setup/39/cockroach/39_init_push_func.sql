-- represents an event to be created.
CREATE TYPE IF NOT EXISTS eventstore.command AS (
    instance_id TEXT
    , aggregate_type TEXT
    , aggregate_id TEXT
    , command_type TEXT
    , revision INT2
    , payload JSONB
    , creator TEXT
    , owner TEXT
);

/*
select * from eventstore.commands_to_events(
ARRAY[
    ROW('', 'system', 'SYSTEM', 'ct1', 1, '{"key": "value"}', 'c1', 'SYSTEM')
    , ROW('', 'system', 'SYSTEM', 'ct2', 1, '{"key": "value"}', 'c1', 'SYSTEM')
    , ROW('289525561255060732', 'org', '289575074711790844', 'ct3', 1, '{"key": "value"}', 'c1', '289575074711790844')
    , ROW('289525561255060732', 'user', '289575075164906748', 'ct3', 1, '{"key": "value"}', 'c1', '289575074711790844')
    , ROW('289525561255060732', 'oidc_session', 'V2_289575178579535100', 'ct3', 1, '{"key": "value"}', 'c1', '289575074711790844')
    , ROW('', 'system', 'SYSTEM', 'ct3', 1, '{"key": "value"}', 'c1', 'SYSTEM')
]::eventstore.command[]
);
*/

-- index is used for filtering for the current sequence of the aggregate
CREATE INDEX IF NOT EXISTS e_push_idx ON eventstore.events2(instance_id, aggregate_type, aggregate_id, "owner", "sequence" DESC);

CREATE OR REPLACE FUNCTION eventstore.commands_to_events(commands eventstore.command[]) RETURNS SETOF eventstore.events2 VOLATILE AS $$
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
        ("cmds").instance_id
        , ("cmds").aggregate_type
        , ("cmds").aggregate_id
        , ("cmds").command_type
        , ("cmds").revision
        , ("cmds").payload
        , ("cmds").creator
        , ("cmds").owner
        , NOW() AS created_at
        , EXTRACT(EPOCH FROM clock_timestamp())::DECIMAL AS position
        , ROW_NUMBER() OVER () AS in_tx_order
    FROM UNNEST(commands) AS "cmds"
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
            ("cmds").instance_id,
            ("cmds").aggregate_type,
            ("cmds").aggregate_id,
            ("cmds").owner
        FROM UNNEST(commands) AS "cmds"
    ) AS a
    LEFT JOIN eventstore.events2 AS e
        ON a.instance_id = e.instance_id
        AND a.aggregate_type = e.aggregate_type
        AND a.aggregate_id = e.aggregate_id
        AND (a.owner = '' OR a.owner = e.owner)
    GROUP BY
        a.instance_id,
        a.aggregate_type,
        a.aggregate_id,
        a.owner
) AS cs
    ON c.instance_id = cs.instance_id
    AND c.aggregate_type = cs.aggregate_type
    AND c.aggregate_id = cs.aggregate_id
    AND (cs.owner = '' OR cs.owner = c.owner)
ORDER BY
    c.in_tx_order
$$ LANGUAGE SQL;

CREATE TYPE IF NOT EXISTS eventstore.event AS (
    instance_id TEXT
    , aggregate_type TEXT
    , aggregate_id TEXT
    , created_at TIMESTAMPTZ
    , "sequence" INT8
    , "position" DECIMAL
    , in_tx_order INT4
);

CREATE OR REPLACE FUNCTION eventstore.push(commands eventstore.command[]) RETURNS SETOF eventstore.event AS $$
    INSERT INTO eventstore.events2
    SELECT * FROM eventstore.commands_to_events(commands)
    RETURNING instance_id, aggregate_type, aggregate_id, created_at, "sequence", position, in_tx_order
$$ LANGUAGE SQL;