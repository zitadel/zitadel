SELECT
    s.user_agent_id
FROM auth.user_sessions s
WHERE
    s.id = $1
    AND s.instance_id = $2
LIMIT 1;