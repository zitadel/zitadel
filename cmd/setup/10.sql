create temporary table IF NOT EXISTS wrong_events (
    instance_id string
    , event_sequence int8
    , current_cd TIMESTAMPTZ
    , next_cd TIMESTAMPTZ
);

TRUNCATE wrong_events;

insert into wrong_events (
    select * from (
        select
            instance_id
            , event_sequence
            , creation_date as current_cd
            , lead(creation_date) over (
                partition by instance_id
                order by event_sequence desc
            ) as next_cd
        from
            eventstore.events
    ) where
        current_cd < next_cd
    order by
        event_sequence desc
);

UPDATE eventstore.events e SET creation_date = we.next_cd FROM wrong_events we WHERE e.event_sequence = we.event_sequence and e.instance_id = we.instance_id;