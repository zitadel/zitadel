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


create index e_push_idx on eventstore.events2(instance_id, aggregate_type, aggregate_id, owner, sequence DESC);

DROP FUNCTION IF EXISTS eventstore.commands_to_events(eventstore.command[]);
CREATE OR REPLACE FUNCTION eventstore.commands_to_events(commands eventstore.command[]) RETURNS SETOF eventstore.events2 AS $$
WITH commands AS (
    SELECT 
        * 
    FROM 
        UNNEST(commands)
), aggregates AS (
    SELECT
        DISTINCT instance_id
        , aggregate_type
        , aggregate_id
        , "owner"
    FROM 
        commands
), current_states AS (
    SELECT
        a.instance_id,
        a.aggregate_type,
        a.aggregate_id,
        a.owner,
        COALESCE(MAX(e.sequence), 0) AS "sequence"
    FROM 
        aggregates a
    LEFT JOIN
        eventstore.events2 e
    ON
        a.instance_id = e.instance_id
        AND a.aggregate_type = e.aggregate_type
        AND a.aggregate_id = e.aggregate_id
        AND a.owner = e.owner
    GROUP BY
        a.instance_id,
        a.aggregate_type,
        a.aggregate_id,
        a.owner
)
SELECT
    c.instance_id
    , c.aggregate_type
    , c.aggregate_id
    , c.command_type AS event_type
    , cs.sequence + ROW_NUMBER() OVER (PARTITION BY c.instance_id, c.aggregate_type, c.aggregate_id) AS "sequence"
    , c.revision
    , c.created_at
    , c.payload
    , c.creator
    , c.owner
    , c.position
    , c.in_tx_order
FROM (
    SELECT
        c.*,
        NOW() AS created_at,
        EXTRACT(EPOCH FROM clock_timestamp()) AS position, 
        ROW_NUMBER() OVER () AS in_tx_order
    FROM commands c
) c
JOIN 
    current_states cs
ON
    c.instance_id = cs.instance_id
    AND c.aggregate_type = cs.aggregate_type
    AND c.aggregate_id = cs.aggregate_id
    AND c.owner = cs.owner
ORDER BY
    c.in_tx_order
$$ LANGUAGE SQL;

DROP FUNCTION IF EXISTS eventstore.push(eventstore.command[]);
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