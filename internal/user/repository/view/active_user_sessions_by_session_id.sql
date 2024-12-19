SELECT
    s.user_agent_id,
    s.user_id,
    s.id
FROM auth.user_sessions s
    JOIN auth.user_sessions s2
        ON s.instance_id = s2.instance_id
        AND s.user_agent_id = s2.user_agent_id
WHERE
    s2.id = $1
    AND s.instance_id = $2
    AND s.state = 0;
