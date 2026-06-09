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
