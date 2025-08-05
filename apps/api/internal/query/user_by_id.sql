SELECT 
  u.id
  , u.creation_date
  , u.change_date
  , u.resource_owner
  , u.sequence
  , u.state
  , u.type
  , u.username
  , login_names.login_names AS login_names
  , login_names.preferred_login_name AS preferred_login_name
  , h.user_id
  , h.first_name
  , h.last_name
  , h.nick_name
  , h.display_name
  , h.preferred_language
  , h.gender
  , h.avatar_key
  , h.email
  , h.is_email_verified
  , h.phone
  , h.is_phone_verified
  , h.password_change_required
  , h.password_changed
  , h.mfa_init_skipped
  , m.user_id
  , m.name
  , m.description
  , m.secret
  , m.access_token_type
  , count(*) OVER ()
FROM projections.users14 u
LEFT JOIN
  projections.users14_humans h
  ON
    u.id = h.user_id
    AND u.instance_id = h.instance_id
LEFT JOIN
  projections.users14_machines m
  ON
    u.id = m.user_id
    AND u.instance_id = m.instance_id
LEFT JOIN LATERAL (
    SELECT
        ARRAY_AGG(ln.login_name ORDER BY ln.login_name) AS login_names,
        MAX(CASE WHEN ln.is_primary THEN ln.login_name ELSE NULL END) AS preferred_login_name
    FROM
        projections.login_names3 AS ln
    WHERE
        ln.user_id = u.id
        AND ln.instance_id = u.instance_id
) AS login_names ON TRUE
WHERE 
  u.id = $1
  AND (u.resource_owner = $2 OR $2 = '')
  AND u.instance_id = $3
LIMIT 1
;