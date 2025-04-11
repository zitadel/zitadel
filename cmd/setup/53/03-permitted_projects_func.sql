DROP FUNCTION IF EXISTS eventstore.permitted_projects;

CREATE OR REPLACE FUNCTION eventstore.permitted_projects(
    req_instance_id TEXT
    , auth_user_id TEXT
    , system_user_perms JSONB
    , perm TEXT
    , filter_org TEXT
    , filter_project TEXT

    , instance_permitted OUT BOOLEAN
    , org_ids OUT TEXT[]
    , project_ids OUT TEXT[]
)
	LANGUAGE 'plpgsql' STABLE
AS $$
BEGIN
    -- if system user
    IF system_user_perms IS NOT NULL THEN
        SELECT p.instance_permitted, p.org_ids INTO instance_permitted, org_ids, project_ids
        FROM eventstore.check_system_user_perms(system_user_perms, req_instance_id, perm) p;
        RETURN;
    END IF;

    -- if human/machine user
    SELECT * FROM eventstore.permitted_orgs(
        req_instance_id
        , auth_user_id
        , system_user_perms
        , perm
        , filter_org
    ) INTO instance_permitted, org_ids;
    IF instance_permitted THEN
        RETURN;
    END IF;
	DECLARE
    	matched_roles TEXT[] := eventstore.find_roles(req_instance_id, perm);
	BEGIN
	    -- Get the projects where permission were granted thru project-level roles
	    SELECT array_agg(sub.project_id) INTO project_ids
	    FROM (
	        SELECT DISTINCT pm.project_id
	        FROM eventstore.project_members pm
	        WHERE pm.role = ANY(matched_roles)
	        AND pm.instance_id = req_instance_id
	        AND pm.user_id = auth_user_id
	        AND (
	            pm.project_id = filter_project
	            OR filter_project IS NULL
	            OR filter_project = ''
	        )
	    ) AS sub;
	END;
END;
$$;
