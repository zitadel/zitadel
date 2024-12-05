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

CREATE OR REPLACE FUNCTION eventstore.commands_to_events(commands eventstore.command[]) RETURNS SETOF eventstore.events2 VOLATILE AS $$
SELECT
    ("c").instance_id
    , ("c").aggregate_type
    , ("c").aggregate_id
    , ("c").command_type AS event_type
    , cs.sequence + ROW_NUMBER() OVER (PARTITION BY ("c").instance_id, ("c").aggregate_type, ("c").aggregate_id ORDER BY ("c").in_tx_order) AS sequence
    , ("c").revision
    , hlc_to_timestamp(cluster_logical_timestamp()) AS created_at
    , ("c").payload
    , ("c").creator
    , cs.owner
    , cluster_logical_timestamp() AS position
    , ("c").in_tx_order   
FROM (
    SELECT 
        ("c").instance_id
        , ("c").aggregate_type
        , ("c").aggregate_id
        , ("c").command_type
        , ("c").revision
        , ("c").payload
        , ("c").creator
        , ("c").owner
        , ROW_NUMBER() OVER () AS in_tx_order
    FROM 
        UNNEST(commands) AS "c"
) AS "c"
JOIN (
    SELECT
        cmds.instance_id
        , cmds.aggregate_type
        , cmds.aggregate_id
        , CASE WHEN (e.owner <> '') THEN e.owner ELSE command_owners.owner END AS owner
        , COALESCE(MAX(e.sequence), 0) AS sequence
    FROM (
        SELECT DISTINCT
            ("cmds").instance_id
            , ("cmds").aggregate_type
            , ("cmds").aggregate_id
            , ("cmds").owner
        FROM UNNEST(commands) AS "cmds"
    ) AS cmds
    LEFT JOIN eventstore.events2 AS e
        ON cmds.instance_id = e.instance_id
        AND cmds.aggregate_type = e.aggregate_type
        AND cmds.aggregate_id = e.aggregate_id
    JOIN (
        SELECT
            DISTINCT ON (
                ("c").instance_id
                , ("c").aggregate_type
                , ("c").aggregate_id
            )
            ("c").instance_id
            , ("c").aggregate_type
            , ("c").aggregate_id
            , ("c").owner
        FROM
            UNNEST(commands) AS "c"
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
    ON ("c").instance_id = cs.instance_id
    AND ("c").aggregate_type = cs.aggregate_type
    AND ("c").aggregate_id = cs.aggregate_id
ORDER BY
    in_tx_order
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION eventstore.push(commands eventstore.command[]) RETURNS SETOF eventstore.events2 AS $$
    INSERT INTO eventstore.events2
    SELECT * FROM eventstore.commands_to_events(commands)
    RETURNING *
$$ LANGUAGE SQL;