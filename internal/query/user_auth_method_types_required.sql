SELECT 
    projections.users14.type
    , auth_methods_force_mfa.force_mfa
    , auth_methods_force_mfa.force_mfa_local_only 
FROM 
    projections.users14 
LEFT JOIN 
    projections.login_policies5 AS auth_methods_force_mfa
ON
    auth_methods_force_mfa.instance_id = projections.users14.instance_id
    AND auth_methods_force_mfa.aggregate_id = ANY(ARRAY[projections.users14.instance_id, projections.users14.resource_owner])
WHERE 
    projections.users14.id = $1
    AND projections.users14.instance_id = $2
ORDER BY 
    auth_methods_force_mfa.is_default 
LIMIT 1;