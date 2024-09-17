WITH events AS (
  select 
    *
    , TRANSACTION_TIMESTAMP() AS created_at
    , EXTRACT(EPOCH FROM TRANSACTION_TIMESTAMP()) as position
    , ROW_NUMBER() OVER(PARTITION BY owner, aggregate_type, aggregate_id) AS increment
    , ROW_NUMBER() OVER() AS position_order
  FROM (VALUES
    %s
  ) AS events (
    instance_id
    , "owner"
    , aggregate_type
    , aggregate_id
    , revision

    , creator
    , event_type
    , payload
  )
), aggregates AS (
  SELECT
    instance_id
    , aggregate_type
    , aggregate_id
    , MAX("sequence") AS "sequence"
  FROM
    eventstore.events2
  WHERE
    (instance_id, aggregate_type, aggregate_id) IN (SELECT instance_id, aggregate_type, aggregate_id FROM events)
  GROUP BY
    instance_id
    , aggregate_type
    , aggregate_id
) INSERT INTO eventstore.events2 (
    instance_id
    , "owner"
    , aggregate_type
    , aggregate_id
    , revision

    , creator
    , event_type
    , payload
    , "sequence"
    , created_at

    , "position"
    , in_tx_order
) SELECT 
  e.instance_id
  , e."owner"
  , e.aggregate_type
  , e.aggregate_id
  , e.revision
  
  , e.creator
  , e.event_type
  , e.payload
  , COALESCE(a.sequence, 0) + e.increment AS "sequence"
  , e.created_at

  , e.position
  , e.position_order
FROM 
  aggregates a
RIGHT JOIN
  events e
ON
  a.instance_id = e.instance_id
  AND a.aggregate_type = e.aggregate_type
  AND a.aggregate_id = e.aggregate_id
RETURNING
  "sequence"
  , created_at
  , "position"
  , in_tx_order
;