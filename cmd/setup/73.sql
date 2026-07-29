-- Finds user grants whose stored roles include roles that are no longer valid
-- (the residue left by the buggy cascade removal, GHSA-v859-c572-qh5p), and
-- returns the corrected role set for each. $1 = instance_id, run once per instance.
--
-- Scope: only grant-based user grants (ug.grant_id set) can be affected by this
-- bug. It was introduced by ChangeProjectGrant (internal/command/project_grant.go),
-- which cascades a *multi*-role removal to the user grants tied to that project
-- grant. Direct user grants never go through a multi-role removal call, so any
-- role mismatch there is unrelated to this CVE (e.g. pre-existing role_key drift)
-- and must not be "fixed" here to avoid stripping unrelated, legitimate roles.
WITH computed AS (
    SELECT
        ug.id,
        ug.resource_owner,
        ug.roles AS current_roles,
        -- Recompute the roles this grant is actually allowed to keep by filtering
        -- its current roles down to the ones still granted by the project grant.
        COALESCE(
            ARRAY(
                SELECT r
                FROM unnest(ug.roles) AS r          -- expand the stored roles array
                WHERE r = ANY(pg.granted_role_keys)
            ),
            ARRAY[]::TEXT[]                          -- keep NULL out; use empty array
        ) AS valid_roles
    FROM projections.user_grants5 ug
    JOIN projections.project_grants4 pg
        ON pg.instance_id = ug.instance_id
       AND pg.grant_id = ug.grant_id
    WHERE ug.instance_id = $1
      AND ug.grant_id IS NOT NULL AND ug.grant_id <> ''  -- direct grants are out of scope, see above
      AND ug.roles IS NOT NULL
      AND cardinality(ug.roles) > 0                 -- nothing to correct on empty grants
)
-- valid_roles is always a subset of current_roles, so a size mismatch means at
-- least one stale role was dropped -> this grant needs a corrective event.
SELECT id, resource_owner, valid_roles
FROM computed
WHERE cardinality(valid_roles) <> cardinality(current_roles);
