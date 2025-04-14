CREATE OR REPLACE FUNCTION eventstore.permitted_orgs(
    instanceId TEXT
    , userId TEXT
    , perm TEXT

    , org_ids OUT TEXT[]
)
	LANGUAGE 'plpgsql'
	STABLE
AS $$
DECLARE
	matched_roles TEXT[]; -- roles containing permission
BEGIN
	SELECT array_agg(rp.role) INTO matched_roles
	FROM eventstore.role_permissions rp
	WHERE rp.instance_id = instanceId
	AND rp.permission = perm;
	
	-- First try if the permission was granted thru an instance-level role
	DECLARE
		has_instance_permission bool;
	BEGIN
		SELECT true INTO has_instance_permission
			FROM eventstore.instance_members im
			WHERE im.role = ANY(matched_roles)
			AND im.instance_id = instanceId
			AND im.user_id = userId
			LIMIT 1;
		
		IF has_instance_permission THEN
			-- Return all organizations
			SELECT array_agg(o.org_id) INTO org_ids
				FROM eventstore.instance_orgs o
				WHERE o.instance_id = instanceId;
			RETURN;
		END IF;
	END;
	
	-- Return the organizations where permission were granted thru org-level roles
	SELECT array_agg(sub.org_id) INTO org_ids
	FROM (
		SELECT DISTINCT om.org_id
		FROM eventstore.org_members om
		WHERE om.role = ANY(matched_roles)
		AND om.instance_id = instanceID
		AND om.user_id = userId
	) AS sub;
    RETURN;
END;
$$;
