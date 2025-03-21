DROP FUNCTION IF EXISTS eventstore.check_system_user_perms;

CREATE OR REPLACE FUNCTION eventstore.check_system_user_perms(
     system_user_perms JSONB
    , perm TEXT
    , filter_orgs TEXT

    , org_ids OUT TEXT[]
)
	LANGUAGE 'plpgsql'
AS $$
DECLARE
      has_instance_or_iam_permission bool;
      member_type_found bool;
      aggregate_ids TEXT;
      check_aggregates bool;
BEGIN

  DROP TABLE IF EXISTS matching_system_user_perms;
  CREATE TEMPORARY TABLE matching_system_user_perms (
          member_type TEXT,
          aggregate_id TEXT,
          object_id TEXT
      ) ON COMMIT DROP;

  INSERT INTO matching_system_user_perms 
  (SELECT * FROM eventstore.get_system_permissions(system_user_perms, perm));

  -- System member type
  SELECT TRUE INTO member_type_found
    FROM matching_system_user_perms msup
    WHERE msup.member_type = 'System'
    LIMIT 1;

  IF member_type_found THEN
    -- Return all organizations or only those in filter_orgs
    SELECT array_agg(o.org_id) INTO org_ids
      FROM eventstore.instance_orgs o
      WHERE
        CASE WHEN filter_orgs != ''
        THEN o.org_id IN (filter_orgs) 
        ELSE TRUE END;
    RETURN;
  END IF;

  -- IAM member type
  SELECT TRUE, array_agg(msup.aggregate_id) INTO member_type_found, aggregate_ids
    FROM matching_system_user_perms msup
    WHERE msup.member_type = 'IAM'
    GROUP BY msup.member_type
    LIMIT 1;

  IF member_type_found THEN
    IF (SELECT FALSE WHERE '' = ANY (aggregate_ids::TEXT[])) = FALSE THEN 
      check_aggregates := FALSE;
    ELSE
      check_aggregates := TRUE;
    END IF;

    -- Return all organizations or only those in filter_orgs
    SELECT array_agg(o.org_id) INTO org_ids
      FROM eventstore.instance_orgs o
      WHERE CASE  
        WHEN check_aggregates THEN 
          o.instance_id = ANY (aggregate_ids::TEXT[])
        ELSE
          TRUE
        END
      AND CASE WHEN filter_orgs != ''
        THEN o.org_id IN (filter_orgs) 
        ELSE TRUE END;
    RETURN;
  END IF;

  -- Organization member type
  SELECT TRUE, array_agg(msup.aggregate_id) INTO member_type_found, aggregate_ids
    FROM matching_system_user_perms msup
    WHERE msup.member_type = 'Organization'
    GROUP BY msup.member_type
    LIMIT 1;

  IF member_type_found THEN
    member_type_found := FALSE;
    -- if any of the aggregate_ids = '', then we search on all organization
    IF (SELECT FALSE WHERE '' = ANY (aggregate_ids::TEXT[])) = FALSE THEN 
      check_aggregates := FALSE;
    ELSE
      check_aggregates := TRUE;
    END IF;

    -- Return all organizations or only those in filter_orgs
    SELECT array_agg(o.org_id) INTO org_ids
      FROM eventstore.instance_orgs o
      WHERE 
        CASE  
          WHEN check_aggregates THEN 
            o.org_id = ANY (aggregate_ids::TEXT[])
          ELSE
            TRUE
          END
      AND CASE WHEN filter_orgs != ''
        THEN o.org_id IN (filter_orgs) 
        ELSE TRUE END;
    RETURN;
  END IF;
END;
$$;

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
AS $$
BEGIN

  -- if system user
  IF system_user_perms IS NOT NULL THEN
    org_ids := eventstore.check_system_user_perms(system_user_perms, perm, filter_orgs)
    RETURN;
  ELSE
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
    END;
  END IF;
END;
$$;

