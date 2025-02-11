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

CREATE OR REPLACE FUNCTION eventstore.latest_aggregate_state(
    instance_id TEXT
    , aggregate_type TEXT
    , aggregate_id TEXT
    
    , sequence OUT BIGINT
    , owner OUT TEXT
)
    LANGUAGE 'plpgsql'
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

CREATE OR REPLACE FUNCTION eventstore.commands_to_events2(commands eventstore.command[])
    RETURNS eventstore.events2[]
    LANGUAGE 'plpgsql'
AS $$
DECLARE
    current_sequence BIGINT;
    current_owner TEXT;

    instance_id TEXT;
    aggregate_type TEXT;
    aggregate_id TEXT;

    _events eventstore.events2[];

    _aggregates CURSOR FOR 
        select 
            DISTINCT ("c").instance_id
            , ("c").aggregate_type
            , ("c").aggregate_id
        FROM
            UNNEST(commands) AS c;
BEGIN
    OPEN _aggregates;
    LOOP
        FETCH NEXT IN _aggregates INTO instance_id, aggregate_type, aggregate_id;
        -- crdb does not support EXIT WHEN NOT FOUND
        EXIT WHEN instance_id IS NULL;

        SELECT 
            * 
        INTO
            current_sequence
            , current_owner 
        FROM eventstore.latest_aggregate_state(
            instance_id
            , aggregate_type
            , aggregate_id
        );

        -- RETURN QUERY is not supported by crdb: https://github.com/cockroachdb/cockroach/issues/105240
        SELECT 
            ARRAY_CAT(_events, ARRAY_AGG(e))
        INTO
            _events
        FROM (
            SELECT
                ("c").instance_id
                , ("c").aggregate_type
                , ("c").aggregate_id
                , ("c").command_type -- AS event_type
                , COALESCE(current_sequence, 0) + ROW_NUMBER() OVER () -- AS sequence
                , ("c").revision
                , NOW() -- AS created_at
                , ("c").payload
                , ("c").creator
                , COALESCE(current_owner, ("c").owner) -- AS owner
                , cluster_logical_timestamp() -- AS position
                , ordinality::INT -- AS in_tx_order
            FROM
                UNNEST(commands) WITH ORDINALITY AS c
            WHERE
                ("c").instance_id = instance_id
                AND ("c").aggregate_type = aggregate_type
                AND ("c").aggregate_id = aggregate_id
        ) AS e;
    END LOOP;
    CLOSE _aggregates;
    RETURN _events;
END;
$$;

CREATE OR REPLACE FUNCTION eventstore.push(commands eventstore.command[]) RETURNS SETOF eventstore.events2 AS $$
    INSERT INTO eventstore.events2
    SELECT
        ("e").instance_id
        , ("e").aggregate_type
        , ("e").aggregate_id
        , ("e").event_type
        , ("e").sequence
        , ("e").revision
        , ("e").created_at
        , ("e").payload
        , ("e").creator
        , ("e").owner
        , ("e")."position"
        , ("e").in_tx_order
    FROM 
        UNNEST(eventstore.commands_to_events2(commands)) e 
    ORDER BY 
        in_tx_order
    RETURNING *
$$ LANGUAGE SQL;

/*
select (c).* from UNNEST(eventstore.commands_to_events2(
ARRAY[
    ROW('', 'system', 'SYSTEM', 'ct1', 1, '{"key": "value"}', 'c1', 'SYSTEM')
    , ROW('', 'system', 'SYSTEM', 'ct2', 1, '{"key": "value"}', 'c1', 'SYSTEM')
    , ROW('289525561255060732', 'org', '289575074711790844', 'ct3', 1, '{"key": "value"}', 'c1', '289575074711790844')
    , ROW('289525561255060732', 'user', '289575075164906748', 'ct3', 1, '{"key": "value"}', 'c1', '289575074711790844')
    , ROW('289525561255060732', 'oidc_session', 'V2_289575178579535100', 'ct3', 1, '{"key": "value"}', 'c1', '289575074711790844')
    , ROW('', 'system', 'SYSTEM', 'ct3', 1, '{"key": "value"}', 'c1', 'SYSTEM')
]::eventstore.command[]
) )c;
*/

