SELECT u.id
FROM projections.login_names3_users u
LEFT JOIN LATERAL (
    SELECT p.must_be_domain
    FROM projections.login_names3_policies AS p
    WHERE
        (
            p.instance_id = ?
            AND NOT p.is_default
            AND p.resource_owner = u.resource_owner
        ) OR (
            p.instance_id = ?
            AND p.is_default
        )
    ORDER BY p.is_default
    LIMIT 1
) AS p ON TRUE
LEFT JOIN projections.login_names3_domains d
    ON p.must_be_domain
    AND u.resource_owner = d.resource_owner
    AND u.instance_id = d.instance_id
    AND d.name_lower = ?
WHERE
    u.instance_id = ?
    AND u.user_name_lower IN (?, ?)
    AND (
        (p.must_be_domain AND u.user_name_lower = ? AND d.name_lower = ?)
        OR (NOT COALESCE(p.must_be_domain, FALSE) AND u.user_name_lower = ?)
    )
