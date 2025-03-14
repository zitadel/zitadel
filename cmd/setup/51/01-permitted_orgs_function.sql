DROP FUNCTION IF EXISTS eventstore.permitted_orgs;

CREATE OR REPLACE FUNCTION eventstore.permitted_orgs(
    instanceId TEXT
    , userId TEXT
    , perm TEXT
    , system_user_memeber_type INTEGER[]
    , system_user_instance_id  TEXT[]
    , system_user_aggregate_id  TEXT[]
    , system_user_permissions TEXT[][]
  , system_user_permissions_length INTEGER[]
    , filter_orgs TEXT

    , org_ids OUT TEXT[]
)
	LANGUAGE 'plpgsql'
	STABLE
AS $$
DECLARE
	matched_roles TEXT[]; -- roles containing permission
BEGIN

  IF system_user_memeber_type IS NOT NULL THEN
    DECLARE
      system_user_permission_found bool;
    BEGIN
      SELECT result.perm_found INTO system_user_permission_found
      FROM (SELECT eventstore.get_org_permission(perm, instanceId,filter_orgs, 
          system_user_memeber_type, system_user_instance_id, system_user_aggregate_id, 
          system_user_permissions, system_user_permissions_length) AS perm_found) AS result;

      IF system_user_permission_found THEN
        SELECT array_agg(o.org_id) INTO org_ids
        FROM eventstore.instance_orgs o
        WHERE o.instance_id = instanceId
        AND CASE WHEN filter_orgs != ''
          THEN o.org_id IN (filter_orgs) 
          ELSE TRUE END;
      END IF;
    END;
    RETURN;
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

DROP FUNCTION IF EXISTS eventstore.get_org_permission; 
CREATE OR REPLACE FUNCTION eventstore.get_org_permission(
  perm TEXT
  , instance_idd TEXT
  , org_id TEXT
  , system_user_memeber_type INTEGER[]
  , sustem_user_instance_id  TEXT[]
  , system_user_aggregate_id  TEXT[]
  , system_user_permissions TEXT[][]
  , system_user_permissions_length INTEGER[]
-- , outt OUT TEXT[]
  , outt OUT BOOL
)
	LANGUAGE 'plpgsql'
AS $$
DECLARE
  i INTEGER;
  length INTEGER;
  permission_length INTEGER;
BEGIN
  -- outt := FALSE;
  length := array_length(system_user_memeber_type, 1);
  -- length := 3;

  DROP TABLE IF EXISTS permissions; 
  CREATE TEMPORARY TABLE permissions (
    member_type INTEGER,
    instance_id TEXT,
    aggregate_id TEXT,
    permission TEXT
    );

  -- <<outer_loop>>
  FOR i IN 1..length LOOP
    -- only interested in organization level and higher
    IF system_user_memeber_type[i] > 3 THEN 
      CONTINUE;
    END IF;
    permission_length := system_user_permissions_length[i];

    FOR j IN 1..permission_length LOOP
    IF system_user_permissions[i][j] != perm THEN 
      CONTINUE;
    END IF;
      INSERT INTO permissions (member_type, instance_id, aggregate_id, permission) VALUES
        (system_user_memeber_type[i], sustem_user_instance_id[i], system_user_aggregate_id[i], system_user_permissions[i][j] );
-- outt := 555;
-- RETURN;
    END LOOP;
  END LOOP;

  -- outt := (SELECT permission FROM permissions LIMIT 1);
  SELECT result.res INTO outt
  FROM (SELECT TRUE AS res FROM permissions p
        WHERE 
          -- check instance id
          CASE WHEN p.member_type = 1 OR p.member_type = 2 THEN -- System or IAM
            p.aggregate_id = instance_idd 
            -- OR p.instance_id IS NULL
            OR p.instance_id = ''
          ELSE 
            p.instance_id = instance_idd 
            -- OR p.instance_id IS NULL
            OR p.instance_id = ''
          END
          AND
          -- check organization
          CASE WHEN p.member_type = 3 THEN -- organization
            p.aggregate_id = org_id
          ELSE
            TRUE
          END
         Limit 1
        ) AS result;
          
  DROP TABLE permissions;

END;
$$;


