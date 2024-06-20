SELECT s.creation_date,
       s.change_date,
       s.resource_owner,
       s.state,
       s.user_agent_id,
       s.user_id,
       u.username,
       l.login_name,
       h.display_name,
       h.avatar_key,
       s.selected_idp_config_id,
       s.password_verification,
       s.passwordless_verification,
       s.external_login_verification,
       s.second_factor_verification,
       s.second_factor_verification_type,
       s.multi_factor_verification,
       s.multi_factor_verification_type,
       s.sequence,
       s.instance_id
FROM auth.user_sessions s
         LEFT JOIN projections.users13 u ON s.user_id = u.id AND s.instance_id = u.instance_id
         LEFT JOIN projections.users13_humans h ON s.user_id = h.user_id AND s.instance_id = h.instance_id
         LEFT JOIN projections.login_names3 l ON s.user_id = l.user_id AND s.instance_id = l.instance_id AND l.is_primary = true
WHERE (s.user_agent_id = $1)
  AND (s.user_id = $2)
  AND (s.instance_id = $3)
LIMIT 1
;