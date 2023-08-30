
CREATE TABLE IF NOT EXISTS projections.failed_events2 (
    projection_name TEXT NOT NULL
    , instance_id TEXT NOT NULL

    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    , event_creation_date TIMESTAMPTZ NOT NULL
    , failed_sequence INT8 NOT NULL

    , failure_count INT2 NULL DEFAULT 0
    , error TEXT
    , last_failed TIMESTAMPTZ
    
    , PRIMARY KEY (projection_name, instance_id, aggregate_type, aggregate_id, failed_sequence)
);

CREATE INDEX IF NOT EXISTS fe2_instance_id_idx on projections.failed_events2 (instance_id);

INSERT INTO projections.failed_events2 (
    projection_name
    , instance_id
    , aggregate_type
    , aggregate_id
    , event_creation_date
    , failed_sequence
    , failure_count
    , error
    , last_failed
) SELECT 
    fe.projection_name
    , fe.instance_id
    , e.aggregate_type
    , e.aggregate_id
    , e.creation_date
    , e.event_sequence
    , fe.failure_count
    , fe.error
    , fe.last_failed
FROM 
    projections.failed_events fe
JOIN eventstore.events e ON
    e.instance_id = fe.instance_id 
    AND e.event_sequence = fe.failed_sequence 
;

INSERT INTO projections.failed_events2 (
    projection_name
    , instance_id
    , aggregate_type
    , aggregate_id
    , event_creation_date
    , failed_sequence
    , failure_count
    , error
    , last_failed
) SELECT 
    fe.view_name
    , fe.instance_id
    , e.aggregate_type
    , e.aggregate_id
    , e.creation_date
    , e.event_sequence
    , fe.failure_count
    , fe.err_msg
    , fe.last_failed
FROM 
    adminapi.failed_events fe
JOIN eventstore.events e ON
    e.instance_id = fe.instance_id 
    AND e.event_sequence = fe.failed_sequence 
;

INSERT INTO projections.failed_events2 (
    projection_name
    , instance_id
    , aggregate_type
    , aggregate_id
    , event_creation_date
    , failed_sequence
    , failure_count
    , error
    , last_failed
) SELECT 
    fe.view_name
    , fe.instance_id
    , e.aggregate_type
    , e.aggregate_id
    , e.creation_date
    , e.event_sequence
    , fe.failure_count
    , fe.err_msg
    , fe.last_failed
FROM 
    auth.failed_events fe
JOIN eventstore.events e ON
    e.instance_id = fe.instance_id 
    AND e.event_sequence = fe.failed_sequence 
;

CREATE TABLE IF NOT EXISTS projections.current_states (
    projection_name TEXT NOT NULL
    , instance_id TEXT NOT NULL

    , last_updated TIMESTAMPTZ

    , aggregate_id TEXT
    , aggregate_type TEXT
    , "sequence" INT8
    , event_date TIMESTAMPTZ
    , position DECIMAL

    , PRIMARY KEY (projection_name, instance_id)
);

INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , event_date
    , position
    , last_updated
) SELECT 
    cs.projection_name
    , cs.instance_id
    , e.creation_date
    , e.position
    , cs.timestamp
FROM 
    projections.current_sequences cs
JOIN eventstore.events e ON
    e.instance_id = cs.instance_id 
    AND e.aggregate_type = cs.aggregate_type
    AND e.event_sequence = cs.current_sequence 
    AND cs.current_sequence = (
        SELECT 
            MAX(cs2.current_sequence)
        FROM
            projections.current_sequences cs2
        WHERE
            cs.projection_name = cs2.projection_name
            AND cs.instance_id = cs2.instance_id
    )
;

INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , event_date
    , position
    , last_updated
) SELECT 
    cs.view_name
    , cs.instance_id
    , e.creation_date
    , e.position
    , cs.last_successful_spooler_run
FROM 
    adminapi.current_sequences cs
JOIN eventstore.events e ON
    e.instance_id = cs.instance_id 
    AND e.event_sequence = cs.current_sequence 
    AND cs.current_sequence = (
        SELECT 
            MAX(cs2.current_sequence)
        FROM
            adminapi.current_sequences cs2
        WHERE
            cs.view_name = cs2.view_name
            AND cs.instance_id = cs2.instance_id
    )
;

INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , event_date
    , position
    , last_updated
) SELECT 
    cs.view_name
    , cs.instance_id
    , e.creation_date
    , e.position
    , cs.last_successful_spooler_run
FROM 
    auth.current_sequences cs
JOIN eventstore.events e ON
    e.instance_id = cs.instance_id 
    AND e.event_sequence = cs.current_sequence 
    AND cs.current_sequence = (
        SELECT 
            MAX(cs2.current_sequence)
        FROM
            auth.current_sequences cs2
        WHERE
            cs.view_name = cs2.view_name
            AND cs.instance_id = cs2.instance_id
    )
;
