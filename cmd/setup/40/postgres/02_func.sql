CREATE OR REPLACE FUNCTION eventstore.commands_to_events(commands eventstore.command[])
    RETURNS SETOF eventstore.events2
    STABLE
    PARALLEL SAFE
AS $$
WITH cmds AS (
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
), aggregates AS (
	SELECT DISTINCT
		c.instance_id
        , c.aggregate_type
        , c.aggregate_id
	FROM cmds c
), max_seq AS (
    SELECT
        e.instance_id
        , e.aggregate_type
        , e.aggregate_id
		, e.owner
        , MAX(e.sequence) sequence
    FROM
        eventstore.events2 e
    JOIN
        aggregates a
    ON
        e.instance_id = a.instance_id
        AND e.aggregate_type = a.aggregate_type
        AND e.aggregate_id = a.aggregate_id
    GROUP BY
        e.instance_id
        , e.aggregate_type
        , e.aggregate_id
		, e.owner
)
SELECT
    cmds.instance_id
    , cmds.aggregate_type
    , cmds.aggregate_id
    , cmds.command_type AS event_type
    , COALESCE(max_seq.sequence, 0) + ROW_NUMBER() OVER (PARTITION BY cmds.instance_id, cmds.aggregate_type, cmds.aggregate_id ORDER BY cmds.in_tx_order) AS sequence
    , cmds.revision
    , NOW() AS created_at
    , cmds.payload
    , cmds.creator
	, COALESCE(max_seq.owner, cmds.owner)
    , EXTRACT(EPOCH FROM NOW()) AS position
    , cmds.in_tx_order
FROM cmds
LEFT JOIN max_seq 
    ON cmds.instance_id = max_seq.instance_id
    AND cmds.aggregate_type = max_seq.aggregate_type
    AND cmds.aggregate_id = max_seq.aggregate_id
ORDER BY
    in_tx_order
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION eventstore.push(commands eventstore.command[]) RETURNS SETOF eventstore.events2 VOLATILE AS $$
INSERT INTO eventstore.events2
SELECT * FROM eventstore.commands_to_events(commands)
RETURNING *
$$ LANGUAGE SQL;
