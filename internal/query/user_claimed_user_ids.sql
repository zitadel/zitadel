SELECT u.id
FROM projections.login_names3_users u
    LEFT JOIN projections.login_names3_policies p_custom
        ON  u.instance_id = p_custom.instance_id
            AND p_custom.instance_id = $1
            AND p_custom.resource_owner = u.resource_owner
    JOIN projections.login_names3_policies p_default
        ON  u.instance_id = p_default.instance_id
            AND p_default.instance_id = $1 AND p_default.is_default IS TRUE
WHERE u.instance_id = $1
    AND COALESCE(p_custom.must_be_domain, p_default.must_be_domain) = false
    AND u.user_name_lower like $2
    AND u.resource_owner <> $3;