WITH current_state_floor AS (
    SELECT
        instance_id,
        MIN(position) AS floor_position
    FROM projections.current_states
    WHERE instance_id <> ''
      AND position IS NOT NULL
    GROUP BY instance_id
),
latest_events AS (
    SELECT DISTINCT ON (instance_id, aggregate_type, aggregate_id)
        instance_id,
        aggregate_type,
        aggregate_id,
        event_type AS last_event_type,
        created_at AS last_created_at,
        position AS last_position
    FROM eventstore.events2
    WHERE aggregate_type IN (
        'auth_request',
        'device_auth',
        'idpintent',
        'oidc_session',
        'saml_request'
    )
    ORDER BY instance_id, aggregate_type, aggregate_id, sequence DESC, in_tx_order DESC
),
latest_oidc_access_added AS (
    SELECT DISTINCT ON (instance_id, aggregate_id)
        instance_id,
        aggregate_id,
        created_at,
        COALESCE((payload->>'lifetime')::bigint, 0) AS lifetime_ns
    FROM eventstore.events2
    WHERE aggregate_type = 'oidc_session'
      AND event_type = 'oidc_session.access_token.added'
    ORDER BY instance_id, aggregate_id, sequence DESC, in_tx_order DESC
),
latest_oidc_access_revoked AS (
    SELECT DISTINCT ON (instance_id, aggregate_id)
        instance_id,
        aggregate_id,
        created_at
    FROM eventstore.events2
    WHERE aggregate_type = 'oidc_session'
      AND event_type = 'oidc_session.access_token.revoked'
    ORDER BY instance_id, aggregate_id, sequence DESC, in_tx_order DESC
),
latest_oidc_refresh_added AS (
    SELECT DISTINCT ON (instance_id, aggregate_id)
        instance_id,
        aggregate_id,
        created_at,
        COALESCE((payload->>'lifetime')::bigint, 0) AS lifetime_ns,
        COALESCE((payload->>'idleLifetime')::bigint, 0) AS idle_lifetime_ns
    FROM eventstore.events2
    WHERE aggregate_type = 'oidc_session'
      AND event_type = 'oidc_session.refresh_token.added'
    ORDER BY instance_id, aggregate_id, sequence DESC, in_tx_order DESC
),
latest_oidc_refresh_renewed AS (
    SELECT DISTINCT ON (instance_id, aggregate_id)
        instance_id,
        aggregate_id,
        created_at,
        COALESCE((payload->>'idleLifetime')::bigint, 0) AS idle_lifetime_ns
    FROM eventstore.events2
    WHERE aggregate_type = 'oidc_session'
      AND event_type = 'oidc_session.refresh_token.renewed'
    ORDER BY instance_id, aggregate_id, sequence DESC, in_tx_order DESC
),
latest_oidc_refresh_revoked AS (
    SELECT DISTINCT ON (instance_id, aggregate_id)
        instance_id,
        aggregate_id,
        created_at
    FROM eventstore.events2
    WHERE aggregate_type = 'oidc_session'
      AND event_type = 'oidc_session.refresh_token.revoked'
    ORDER BY instance_id, aggregate_id, sequence DESC, in_tx_order DESC
),
terminal_candidates AS (
    SELECT
        le.instance_id,
        le.aggregate_type,
        le.aggregate_id,
        le.last_created_at AS sort_time
    FROM latest_events le
    JOIN current_state_floor cs
        ON cs.instance_id = le.instance_id
    WHERE le.aggregate_type <> 'oidc_session'
      AND le.last_created_at < $1
      AND le.last_position <= cs.floor_position
      AND (
          (le.aggregate_type = 'auth_request' AND le.last_event_type IN (
              'auth_request.failed',
              'auth_request.succeeded'
          ))
          OR (le.aggregate_type = 'saml_request' AND le.last_event_type IN (
              'saml_request.failed',
              'saml_request.succeeded'
          ))
          OR (le.aggregate_type = 'device_auth' AND le.last_event_type IN (
              'device.authorization.canceled',
              'device.authorization.done'
          ))
          OR (le.aggregate_type = 'idpintent' AND le.last_event_type IN (
              'idpintent.consumed',
              'idpintent.failed'
          ))
      )
),
oidc_session_candidates AS (
    SELECT
        le.instance_id,
        le.aggregate_type,
        le.aggregate_id,
        oidc_state.invalid_after AS sort_time
    FROM latest_events le
    JOIN current_state_floor cs
        ON cs.instance_id = le.instance_id
    LEFT JOIN latest_oidc_access_added aaa
        ON aaa.instance_id = le.instance_id
       AND aaa.aggregate_id = le.aggregate_id
    LEFT JOIN latest_oidc_access_revoked aar
        ON aar.instance_id = le.instance_id
       AND aar.aggregate_id = le.aggregate_id
    LEFT JOIN latest_oidc_refresh_added rta
        ON rta.instance_id = le.instance_id
       AND rta.aggregate_id = le.aggregate_id
    LEFT JOIN latest_oidc_refresh_renewed rtr
        ON rtr.instance_id = le.instance_id
       AND rtr.aggregate_id = le.aggregate_id
    LEFT JOIN latest_oidc_refresh_revoked rrv
        ON rrv.instance_id = le.instance_id
       AND rrv.aggregate_id = le.aggregate_id
    CROSS JOIN LATERAL (
        SELECT
            CASE
                WHEN aaa.created_at IS NULL
                    AND aar.created_at IS NULL
                    AND rta.created_at IS NULL
                    AND rtr.created_at IS NULL
                    AND rrv.created_at IS NULL THEN le.last_created_at
                ELSE GREATEST(
                    COALESCE(
                        GREATEST(
                            COALESCE(aaa.created_at + ((aaa.lifetime_ns / 1000.0) * interval '1 microsecond'), '-infinity'::timestamptz),
                            COALESCE(aar.created_at, '-infinity'::timestamptz),
                            COALESCE(rrv.created_at, '-infinity'::timestamptz)
                        ),
                        '-infinity'::timestamptz
                    ),
                    COALESCE(
                        GREATEST(
                            LEAST(
                                COALESCE(rta.created_at + ((rta.lifetime_ns / 1000.0) * interval '1 microsecond'), 'infinity'::timestamptz),
                                GREATEST(
                                    COALESCE(rta.created_at + ((rta.idle_lifetime_ns / 1000.0) * interval '1 microsecond'), '-infinity'::timestamptz),
                                    COALESCE(rtr.created_at + ((rtr.idle_lifetime_ns / 1000.0) * interval '1 microsecond'), '-infinity'::timestamptz)
                                )
                            ),
                            COALESCE(rrv.created_at, '-infinity'::timestamptz)
                        ),
                        '-infinity'::timestamptz
                    )
                )
            END AS invalid_after
    ) AS oidc_state
    WHERE le.aggregate_type = 'oidc_session'
      AND le.last_position <= cs.floor_position
      AND oidc_state.invalid_after < $1
),
candidates AS (
    SELECT instance_id, aggregate_type, aggregate_id
    FROM (
        SELECT * FROM terminal_candidates
        UNION ALL
        SELECT * FROM oidc_session_candidates
    ) AS all_candidates
    ORDER BY sort_time
    LIMIT $2
),
deleted AS (
    DELETE FROM eventstore.events2 e
    USING candidates c
    WHERE c.instance_id = e.instance_id
      AND c.aggregate_type = e.aggregate_type
      AND c.aggregate_id = e.aggregate_id
    RETURNING e.event_type, e.instance_id, e.aggregate_type, e.aggregate_id
)
SELECT
    event_type,
    COUNT(*) AS row_count,
    COUNT(DISTINCT (instance_id, aggregate_type, aggregate_id)) AS aggregate_count
FROM deleted
GROUP BY event_type
ORDER BY event_type;
