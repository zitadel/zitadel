CREATE temporary TABLE IF NOT EXISTS wrong_events (
    instance_id TEXT
    , event_sequence BIGINT
    , current_cd TIMESTAMPTZ
    , next_cd TIMESTAMPTZ
);

TRUNCATE wrong_events;

INSERT INTO wrong_events (
    SELECT * FROM (
        SELECT
            instance_id
            , event_sequence
            , creation_date AS current_cd
            , lead(creation_date) OVER (
                PARTITION BY instance_id
                ORDER BY event_sequence DESC
            ) AS next_cd
        FROM
            eventstore.events
    ) sub WHERE
        current_cd < next_cd
    ORDER BY
        event_sequence DESC
);

UPDATE eventstore.events e SET creation_date = we.next_cd FROM wrong_events we WHERE e.event_sequence = we.event_sequence and e.instance_id = we.instance_id;