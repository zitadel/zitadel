INSERT INTO projections.current_states (
    instance_id
   , projection_name
   , last_updated
   , sequence
   , event_date
   , position
   , filter_offset
)
    SELECT instance_id
         , 'projections.notifications_back_channel_logout'
         , now()
         , $1
         , $2
         , $3
         , 0
    FROM eventstore.events2
        WHERE aggregate_type = 'instance'
            AND event_type = 'instance.added'
    ON CONFLICT DO NOTHING;