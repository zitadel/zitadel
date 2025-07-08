INSERT INTO projections.current_states AS cs ( instance_id
                                             , projection_name
                                             , last_updated
                                             , sequence
                                             , event_date
                                             , position
                                             , filter_offset)
SELECT instance_id
     , 'projections.execution_handler'
     , now()
     , $1
     , $2
     , $3
     , 0
FROM eventstore.events2 AS e
WHERE aggregate_type = 'instance'
  AND event_type = 'instance.added'
ON CONFLICT (instance_id, projection_name) DO UPDATE SET last_updated  = EXCLUDED.last_updated,
                                                         sequence      = EXCLUDED.sequence,
                                                         event_date    = EXCLUDED.event_date,
                                                         position      = EXCLUDED.position,
                                                         filter_offset = EXCLUDED.filter_offset;