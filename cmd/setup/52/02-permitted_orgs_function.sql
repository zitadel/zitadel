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

    WITH tmp(member_type, aggregate_id, object_id ) AS (
      SELECT * FROM eventstore.get_system_permissions(
          system_user_perms,
          perm)
    )

    SELECT array_agg(DISTINCT o.org_id) INTO org_ids
      FROM eventstore.instance_orgs o, tmp
      WHERE
        CASE WHEN (SELECT TRUE WHERE tmp.member_type = 'System' LIMIT 1) THEN
          TRUE
        WHEN (SELECT TRUE WHERE tmp.member_type = 'IAM' LIMIT 1) THEN
          -- aggregate_id not present
          CASE WHEN (SELECT TRUE WHERE '' = ANY (
                              (
                                SELECT array_agg(tmp.aggregate_id)
                                  FROM tmp
                                  WHERE member_type = 'IAM'
                                  GROUP BY member_type
                                  LIMIT 1
                              )::TEXT[])) THEN
            TRUE
          ELSE 
            o.instance_id = ANY (
                            (
                              SELECT array_agg(tmp.aggregate_id)
                                FROM tmp
                                WHERE member_type = 'IAM'
                                GROUP BY member_type
                                LIMIT 1
                            )::TEXT[])
        END
        WHEN (SELECT TRUE WHERE tmp.member_type = 'Organization' LIMIT 1) THEN
          -- aggregate_id not present
          CASE WHEN (SELECT TRUE WHERE '' = ANY (
                              (
                                SELECT array_agg(tmp.aggregate_id)
                                  FROM tmp
                                  WHERE member_type = 'Organization'
                                  GROUP BY member_type
                                  LIMIT 1
                              )::TEXT[])) THEN
            TRUE
          ELSE 
            o.org_id = ANY (
                        (
                          SELECT array_agg(tmp.aggregate_id)
                            FROM tmp
                            WHERE member_type = 'Organization'
                            GROUP BY member_type
                            LIMIT 1
                        )::TEXT[])
          END
        END
        AND
        CASE WHEN filter_orgs != ''
        THEN o.org_id IN (filter_orgs) 
        ELSE TRUE END
      LIMIT 1;
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

~/go/src/zitadel/cmd/setup/52/02-permitted_orgs_function.sql  -- if system user
  IF system_user_perms IS NOT NULL THEN
    org_ids := eventstore.check_system_user_perms(system_user_perms, perm, filter_orgs);
  -- if human/machine user
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

