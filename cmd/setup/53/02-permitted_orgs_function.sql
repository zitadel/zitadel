DROP FUNCTION IF EXISTS eventstore.permitted_orgs;
DROP FUNCTION IF EXISTS eventstore.find_roles;

-- find_roles finds all roles containing the permission
CREATE OR REPLACE FUNCTION eventstore.find_roles(
    req_instance_id TEXT
    , perm TEXT

    , roles OUT TEXT[]
)
LANGUAGE 'plpgsql' STABLE
AS $$
BEGIN
    SELECT array_agg(rp.role) INTO roles
    FROM eventstore.role_permissions rp
    WHERE rp.instance_id = req_instance_id
    AND rp.permission = perm;
END;
$$;

CREATE OR REPLACE FUNCTION eventstore.permitted_orgs(
    req_instance_id TEXT
    , auth_user_id TEXT
    , system_user_perms JSONB
    , perm TEXT
    , filter_org TEXT

    , instance_permitted OUT BOOLEAN
    , org_ids OUT TEXT[]
)
	LANGUAGE 'plpgsql' STABLE
AS $$
BEGIN
    -- if system user
    IF system_user_perms IS NOT NULL THEN
        SELECT p.instance_permitted, p.org_ids INTO instance_permitted, org_ids
        FROM eventstore.check_system_user_perms(system_user_perms, req_instance_id, perm) p;
        RETURN;
    END IF;
  
    -- if human/machine user
    DECLARE
    	matched_roles TEXT[] := eventstore.find_roles(req_instance_id, perm);
	BEGIN
        -- First try if the permission was granted thru an instance-level role
        SELECT true INTO instance_permitted
            FROM eventstore.instance_members im
            WHERE im.role = ANY(matched_roles)
            AND im.instance_id = req_instance_id
            AND im.user_id = auth_user_id
            LIMIT 1;
        
        org_ids := ARRAY[]::TEXT[];
        IF instance_permitted THEN
            RETURN;
        END IF;
        instance_permitted := FALSE;

        -- Return the organizations where permission were granted thru org-level roles
        SELECT array_agg(sub.org_id) INTO org_ids
        FROM (
            SELECT DISTINCT om.org_id
            FROM eventstore.org_members om
            WHERE om.role = ANY(matched_roles)
            AND om.instance_id = req_instance_id
            AND om.user_id = auth_user_id
            AND (filter_org IS NULL OR om.org_id = filter_org)
        ) AS sub;
    END;
END;
$$;
