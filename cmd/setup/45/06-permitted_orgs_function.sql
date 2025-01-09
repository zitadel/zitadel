CREATE OR REPLACE FUNCTION 
	eventstore.permitted_orgs(instanceId text, userId text, perm text)
RETURNS SETOF text AS $$
DECLARE
	matched_roles text[]; -- roles containing permission
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
			RETURN QUERY SELECT o.org_id
				FROM eventstore.instance_orgs o
				WHERE o.instance_id = instanceId;
			RETURN;
		END IF;
	END;
	
	-- Return the organizations where permission were granted thru org-level roles
	RETURN QUERY SELECT DISTINCT om.org_id
		FROM eventstore.org_members om
		WHERE om.role = ANY(matched_roles)
		AND om.instance_id = instanceID
		AND om.user_id = userId;
END;
$$
LANGUAGE plpgsql STABLE;
