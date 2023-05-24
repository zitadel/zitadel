TRUNCATE wrong_events;

INSERT INTO wrong_events (
    SELECT * FROM (
        SELECT
            instance_id
            , event_sequence
            , created_at AS current_cd
            , lead(created_at) OVER (
                PARTITION BY instance_id
                ORDER BY event_sequence DESC
            ) AS next_cd
        FROM
            eventstore.events
        WHERE 
            creation_date > '2023-05-23 13:00'
    ) sub AS OF SYSTEM TIME follower_read_timestamp() WHERE
        current_cd < next_cd
    ORDER BY
        event_sequence DESC
);
