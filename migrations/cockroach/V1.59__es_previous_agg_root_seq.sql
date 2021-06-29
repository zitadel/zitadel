BEGIN;

ALTER TABLE eventstore.events
    RENAME COLUMN previous_sequence TO previous_aggregate_sequence,
    ADD COLUMN previous_aggregate_type_sequence INT8, 
    ADD CONSTRAINT prev_agg_type_seq_unique UNIQUE(previous_aggregate_type_sequence);

COMMIT;

BEGIN;

-- --------------------------------------------------------
-- iam
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'iam' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- key_pair
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'key_pair' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- org
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'org' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- org.domain
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'org.domain' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- org.name
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'org' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- policy.password.complexity
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'policy.password.complexity' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- project
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'project' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- user
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'user' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- user.email
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'user.email' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- user.human.externalidp
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'user.human.externalidp' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- user.username
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'user.username' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- usergrant
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'usergrant' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;


-- --------------------------------------------------------
-- usergrant.unique
-- --------------------------------------------------------
WITH RECURSIVE set_previous (agg_type, agg_seq, agg_type_seq) AS (
    (SELECT aggregate_type, MIN(event_sequence), NULL::INT8 FROM eventstore.events WHERE aggregate_type = 'usergrant.unique' GROUP BY aggregate_type)
UNION ALL
    (SELECT 
        e.aggregate_type, e.event_sequence, sp.agg_seq
    FROM 
        eventstore.events e
    JOIN set_previous sp 
    ON
        e.aggregate_type = sp.agg_type
    AND 
        e.event_sequence > sp.agg_seq
    ORDER BY e.event_sequence ASC
    LIMIT 1)
)
UPDATE 
    eventstore.events 
SET 
    previous_aggregate_type_sequence = sp.agg_type_seq 
FROM 
    set_previous sp 
WHERE 
    event_sequence = sp.agg_seq;



COMMIT;
