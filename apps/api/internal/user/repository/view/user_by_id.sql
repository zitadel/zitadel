WITH auth_methods AS (
  SELECT
    user_id
    , method_type
    , token_id
    , state
    , instance_id
    , name
  FROM
    projections.user_auth_methods5
  WHERE
    instance_id = $1
    AND user_id = $2
),
verified_auth_methods AS (
  SELECT
    method_type
  FROM
    auth_methods
  WHERE state = 2
)
SELECT
    u.id
    , u.creation_date
    , LEAST(u.change_date, au.change_date) AS change_date
    , u.resource_owner
    , u.state AS user_state
    , au.password_set
    , h.password_change_required
    , au.password_change
    , au.last_login
    , u.username AS user_name
    , (SELECT array_agg(ll.login_name) login_names FROM projections.login_names3 ll
                                                   WHERE u.instance_id = ll.instance_id AND u.id = ll.user_id
                                                   GROUP BY ll.user_id, ll.instance_id) AS login_names
    , l.login_name as preferred_login_name
    , h.first_name
    , h.last_name
    , h.nick_name
    , h.display_name
    , h.preferred_language
    , h.gender
    , h.email
    , h.is_email_verified
    , n.verified_email
    , h.phone
    , h.is_phone_verified
    , (SELECT COALESCE((SELECT state FROM auth_methods WHERE method_type = 1), 0)) AS otp_state
    , CASE
        WHEN EXISTS (SELECT true FROM verified_auth_methods WHERE method_type = 3) THEN 2
        WHEN EXISTS (SELECT true FROM verified_auth_methods WHERE method_type = 2) THEN 1
        ELSE 0
      END AS mfa_max_set_up
    , au.mfa_init_skipped
    , u.sequence
    , au.init_required
    , au.username_change_required
    , m.name AS machine_name
    , m.description AS machine_description
    , u.type AS user_type
    , (SELECT
          JSONB_AGG(json_build_object('webAuthNTokenId', token_id, 'webAuthNTokenName', name, 'state', state))
        FROM auth_methods
        WHERE method_type = 2
        ) AS u2f_tokens
    , (SELECT
        JSONB_AGG(json_build_object('webAuthNTokenId', token_id, 'webAuthNTokenName', name, 'state', state))
        FROM auth_methods
        WHERE method_type = 3
        ) AS passwordless_tokens
    , h.avatar_key
    , au.passwordless_init_required
    , au.password_init_required
    , u.instance_id
    , (SELECT EXISTS (SELECT true FROM verified_auth_methods WHERE method_type = 6)) AS otp_sms_added
    , (SELECT EXISTS (SELECT true FROM verified_auth_methods WHERE method_type = 7)) AS otp_email_added
FROM projections.users14 u
    LEFT JOIN projections.users14_humans h
        ON u.instance_id = h.instance_id
        AND u.id = h.user_id
    LEFT JOIN projections.users14_notifications n
        ON u.instance_id = n.instance_id
        AND u.id = n.user_id
    LEFT JOIN projections.login_names3 l
        ON u.instance_id = l.instance_id
        AND u.id = l.user_id
        AND l.is_primary = true
    LEFT JOIN projections.users14_machines m
        ON u.instance_id = m.instance_id
        AND u.id = m.user_id
    LEFT JOIN auth.users3 au
        ON u.instance_id = au.instance_id
        AND u.id = au.id
WHERE
  u.instance_id = $1
  AND u.id = $2
LIMIT 1;