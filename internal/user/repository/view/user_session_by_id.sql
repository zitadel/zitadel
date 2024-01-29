SELECT
    auth.user_sessions.creation_date,
    auth.user_sessions.change_date,
    auth.user_sessions.resource_owner,
    auth.user_sessions.state,
    auth.user_sessions.user_agent_id,
    auth.user_sessions.user_id,
    u.username,
    l.login_name,
    h.display_name,
    h.avatar_key,
    auth.user_sessions.selected_idp_config_id,
    auth.user_sessions.password_verification,
    auth.user_sessions.passwordless_verification,
    auth.user_sessions.external_login_verification,
    auth.user_sessions.second_factor_verification,
    auth.user_sessions.second_factor_verification_type,
    auth.user_sessions.multi_factor_verification,
    auth.user_sessions.multi_factor_verification_type,
    auth.user_sessions.sequence,
    auth.user_sessions.instance_id
FROM auth.user_sessions
    LEFT JOIN projections.users10 u ON auth.user_sessions.user_id = u.id
    LEFT JOIN projections.users10_humans h ON auth.user_sessions.user_id = h.user_id
    LEFT JOIN projections.login_names3 l ON auth.user_sessions.user_id = l.user_id
WHERE (auth.user_sessions.user_agent_id = $1)
  AND (auth.user_sessions.user_id = $2)
  AND (auth.user_sessions.instance_id = $3)
    LIMIT 1