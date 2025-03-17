DROP FUNCTION IF EXISTS eventstore.permitted_orgs;

CREATE OR REPLACE FUNCTION eventstore.permitted_orgs(
    instanceId TEXT
    , userId TEXT
    , system_user_perms JSONB
    , perm TEXT
    , filter_orgs TEXT

    , org_ids OUT TEXT[]
)
	LANGUAGE 'plpgsql'
	-- STABLE
AS $$
DECLARE
	matched_roles TEXT[]; -- roles containing permission
BEGIN

  -- if system user
  IF jsonb_array_length(system_user_perms) = 0 THEN
    DECLARE
      has_instance_or_iam_permission bool;
    BEGIN

    DROP TABLE IF EXISTS matching_system_user_perms;
    CREATE TEMPORARY TABLE matching_system_user_perms (
            member_type TEXT,
            -- instance_id TEXT,
            aggregate_id TEXT,
            object_id TEXT,
            permission TEXT
        ) ON COMMIT DROP;


    INSERT INTO matching_system_user_perms 
    (SELECT * FROM eventstore.get_system_permissions(system_user_perms, perm));

      -- check instance or iam level
      SELECT true INTO has_instance_or_iam_permission
        FROM matching_system_user_perms msup
        WHERE (msup.member_type = 'System' AND msup.aggregate_id = '')
        OR (msup.member_type = 'System' AND msup.aggregate_id = instanceId)
        OR (msup.member_type = 'IAM' AND msup.aggregate_id = instanceId)
        LIMIT 1;

      IF has_instance_or_iam_permission THEN
        -- Return all organizations or only those in filter_orgs
        SELECT array_agg(o.org_id) INTO org_ids
          FROM eventstore.instance_orgs o
          WHERE o.instance_id = instanceId
          AND CASE WHEN filter_orgs != ''
            THEN o.org_id IN (filter_orgs) 
            ELSE TRUE END;
        RETURN;
      END IF;

    -- SELECT array_agg(msup.aggregate_id) INTO org_ids
    --   FROM matching_system_user_perms msup
    --   WHERE msup.instance_id = instanceId
    --   AND msup.member_type = 'Organization'
    --   AND CASE WHEN filter_orgs != ''
    --     THEN msup.aggregate_id IN (filter_orgs) 
    --     ELSE TRUE END;
    --   RETURN;

    END;
  END IF;

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
			-- Return all organizations or only those in filter_orgs
			SELECT array_agg(o.org_id) INTO org_ids
				FROM eventstore.instance_orgs o
				WHERE o.instance_id = instanceId
				AND CASE WHEN filter_orgs != ''
					THEN o.org_id IN (filter_orgs) 
					ELSE TRUE END;
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




DROP FUNCTION IF EXISTS eventstore.get_system_permissions;
CREATE OR REPLACE FUNCTION eventstore.get_system_permissions(
    permissions_json JSONB
    , permm TEXT
    -- , res OUT eventstore.system_perms
)
RETURNS TABLE (
    member_type TEXT,
    -- instance_id TEXT,
    aggregate_id TEXT,
    object_id TEXT,
    permission TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT *  FROM (
    SELECT 
        (perm)->>'member_type' AS member_type,
        (perm)->>'aggregate_id' AS agregate_id,
        (perm)->>'object_id' AS objectId,
        permis AS permission
        FROM jsonb_array_elements(permissions_json) AS perm
        CROSS JOIN jsonb_array_elements_text(perm->'permissions') AS permis) as p
        WHERE p.permission = permm;
END;
$$ LANGUAGE plpgsql;
