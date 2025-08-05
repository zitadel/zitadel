WITH corrupt_streams AS (
    select
        e.instance_id
        , e.aggregate_type
        , e.aggregate_id
        , min(e.sequence) as min_sequence
        , count(distinct e.owner) as owner_count
    from
        eventstore.events2 e
    where
        e.instance_id = $1
        and aggregate_type = 'project'
    group by 
        e.instance_id
        , e.aggregate_type
        , e.aggregate_id
    having
        count(distinct e.owner) > 1
), correct_owners AS (
    select
        e.instance_id
        , e.aggregate_type
        , e.aggregate_id
        , e.owner
    from
        eventstore.events2 e
    join
        corrupt_streams cs
    on
        e.instance_id = cs.instance_id
        and e.aggregate_type = cs.aggregate_type
        and e.aggregate_id = cs.aggregate_id
        and e.sequence = cs.min_sequence
), wrong_events AS (
    select
        e.instance_id
        , e.aggregate_type
        , e.aggregate_id
        , e.sequence
        , e.owner wrong_owner
        , co.owner correct_owner
    from
        eventstore.events2 e
    join
        correct_owners co
    on
        e.instance_id = co.instance_id
        and e.aggregate_type = co.aggregate_type
        and e.aggregate_id = co.aggregate_id
        and e.owner <> co.owner
), updated_events AS (
    UPDATE eventstore.events2 e
        SET owner = we.correct_owner
    FROM 
        wrong_events we
    WHERE
        e.instance_id = we.instance_id
        and e.aggregate_type = we.aggregate_type
        and e.aggregate_id = we.aggregate_id
        and e.sequence = we.sequence
    RETURNING 
        we.aggregate_id
        , we.correct_owner
        , we.sequence
        , we.wrong_owner
)
SELECT
    ue.aggregate_id
    , ue.correct_owner
    , jsonb_object_agg(
        ue.sequence::TEXT --formant to string because crdb is not able to handle int
        , ue.wrong_owner
    ) payload
FROM
    updated_events ue
GROUP BY
    ue.aggregate_id
    , ue.correct_owner
;
