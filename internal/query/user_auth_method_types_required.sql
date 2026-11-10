SELECT
    projections.users14.type
    , auth_methods_force_mfa.force_mfa
    , auth_methods_force_mfa.force_mfa_local_only
    , auth_methods_force_mfa.second_factors
    , user_auth_methods5.auth_method_types
FROM
    projections.users14
LEFT JOIN
    projections.login_policies5 AS auth_methods_force_mfa
ON
    auth_methods_force_mfa.instance_id = projections.users14.instance_id
    AND auth_methods_force_mfa.aggregate_id = ANY(ARRAY[projections.users14.instance_id, projections.users14.resource_owner])
LEFT JOIN LATERAL (
    SELECT
        ARRAY_AGG(projections.user_auth_methods5.method_type) AS auth_method_types
    FROM
        projections.user_auth_methods5
    WHERE
        projections.user_auth_methods5.user_id = projections.users14.id
        AND projections.user_auth_methods5.instance_id = projections.users14.instance_id
        AND projections.user_auth_methods5.state = 2
    ) AS user_auth_methods5 ON TRUE
WHERE
    projections.users14.id = $1
    AND projections.users14.instance_id = $2
ORDER BY
    auth_methods_force_mfa.is_default
LIMIT 1;