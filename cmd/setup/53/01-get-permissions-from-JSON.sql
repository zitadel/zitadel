DROP FUNCTION IF EXISTS eventstore.check_system_user_perms;
DROP FUNCTION IF EXISTS eventstore.get_system_permissions;
DROP TYPE IF EXISTS eventstore.project_grant;


/*
	Function get_system_permissions unpacks an JSON array of system member permissions,
	into a table format. Each array entry maps to one row representing a membership which
	contained the req_permission.

	[
		{
		"member_type": "IAM",
		"aggregate_id": "310716990375453665",
		"object_id": "",
		"permissions": ["iam.read", "iam.write", "iam.policy.read"]
		},
		...
	]

	| member_type | aggregate_id         | object_id |
	| "IAM"       | "310716990375453665" | null      |
*/
CREATE OR REPLACE FUNCTION eventstore.get_system_permissions(
    permissions_json JSONB
    , permm TEXT
)
RETURNS TABLE (
    member_type TEXT,
    aggregate_id TEXT,
    object_id TEXT
)
  LANGUAGE 'plpgsql' IMMUTABLE
AS $$
BEGIN
    RETURN QUERY
    SELECT res.member_type, res.aggregate_id, res.object_id  FROM (
    SELECT 
        (perm)->>'member_type' AS member_type,
        (perm)->>'aggregate_id' AS aggregate_id,
        (perm)->>'object_id' AS object_id,
        permission
        FROM jsonb_array_elements(permissions_json) AS perm
        CROSS JOIN jsonb_array_elements_text(perm->'permissions') AS permission) AS res
        WHERE res.permission = permm;
END;
$$;

/*
	Type project_grant is composite identifier using its project and grant IDs.
*/
CREATE TYPE eventstore.project_grant AS (
    project_id TEXT -- mapped from a permission's aggregate_id
    , grant_id TEXT -- mapped from a permission's object_id
);

/*
	Function check_system_user_perms uses system member permissions to establish
	on which organization, project or project grant the user has the requested permission.
	The permission can also apply to the complete instance when a IAM membership matches
	the requested instance ID, or through system membership.

	See eventstore.get_system_permissions() on the supported JSON format.
*/
CREATE OR REPLACE FUNCTION eventstore.check_system_user_perms(
     system_user_perms JSONB
    , req_instance_id TEXT
    , perm TEXT

    , instance_permitted OUT BOOLEAN
    , org_ids OUT TEXT[]
    , project_ids OUT TEXT[]
    , project_grants OUT eventstore.project_grant[]
)
	LANGUAGE 'plpgsql' IMMUTABLE
AS $$
BEGIN
	-- make sure no nulls are returned
	instance_permitted := FALSE;
	org_ids := ARRAY[]::TEXT[];
	project_ids := ARRAY[]::TEXT[];
	project_grants := ARRAY[]::eventstore.project_grant[];
	DECLARE
	    p RECORD;
	BEGIN
        FOR p IN SELECT member_type, aggregate_id, object_id
	        FROM eventstore.get_system_permissions(system_user_perms, perm)
	    LOOP
	       CASE p.member_type
	            WHEN 'System' THEN
	                instance_permitted := TRUE;
	                RETURN;
	            WHEN 'IAM' THEN
	                IF p.aggregate_id = req_instance_id THEN
	                    instance_permitted := TRUE;
	                    RETURN;
	                END IF;
	            WHEN 'Organization' THEN
	                IF p.aggregate_id != '' THEN
	                    org_ids := array_append(org_ids, p.aggregate_id);
	                END IF;
	            WHEN 'Project' THEN
	                IF p.aggregate_id != '' THEN
	                    project_ids := array_append(project_ids, p.aggregate_id);
	                END IF;
	            WHEN 'ProjectGrant' THEN
	                IF p.aggregate_id != '' THEN
	                    project_grants := array_append(project_grants, ROW(p.aggregate_id, p.object_id)::eventstore.project_grant);
	                END IF;
	        END CASE;
	    END LOOP;
	END;
END;
$$;
