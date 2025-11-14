with application as (
    SELECT a.instance_id,
           a.resource_owner,
           a.project_id,
           a.id as app_id,
           p.project_role_check,
           p.has_project_check
    FROM projections.apps7 as a
         LEFT JOIN projections.apps7_oidc_configs as aoc
                   ON aoc.app_id = a.id
                   AND aoc.instance_id = a.instance_id
         INNER JOIN projections.projects4 as p
                    ON p.instance_id = a.instance_id
                    AND p.resource_owner = a.resource_owner
                    AND p.id = a.project_id
    WHERE a.instance_id = $1
      AND aoc.client_id = $2
      AND a.state = $3
      AND p.state = $4
), user_resourceowner as (
/* resourceowner of the active user */
     SELECT u.instance_id,
            u.resource_owner,
            u.id as user_id
     FROM projections.users14 as u
     WHERE u.instance_id = $1
       AND u.id = $5
       AND u.state = $6
), has_project_grant_check as (
/* all projectgrants active, then filtered with the project and user resourceowner */
     SELECT pg.instance_id,
            pg.resource_owner,
            pg.project_id,
            pg.granted_org_id
     FROM projections.project_grants4 as pg
     WHERE pg.instance_id = $1
       AND pg.state = $7
), project_role_check as (
/* all usergrants active and associated with the user, then filtered with the project */
     SELECT ug.instance_id,
            ug.resource_owner,
            ug.project_id
     FROM projections.user_grants5 as ug
     WHERE ug.instance_id = $1
       AND ug.user_id = $5
       AND ug.state = $8
)
SELECT
    /* project existence does not need to be checked, or resourceowner of user and project are equal, or resourceowner of user has project granted*/
       bool_and(COALESCE(
               (NOT a.has_project_check OR
                a.resource_owner = uro.resource_owner OR
                uro.resource_owner = hpgc.granted_org_id)
           , FALSE)
       ) as project_checked,
    /* authentication existence does not need to checked, or authentication for project is existing*/
       bool_and(COALESCE(
               (NOT a.project_role_check OR
                a.project_id = prc.project_id)
           , FALSE)
       ) as role_checked
FROM application as a
         LEFT JOIN user_resourceowner as uro
                   ON uro.instance_id = a.instance_id
         LEFT JOIN has_project_grant_check as hpgc
                   ON hpgc.instance_id = a.instance_id
                   AND hpgc.project_id = a.project_id
                   AND hpgc.granted_org_id = uro.resource_owner
         LEFT JOIN project_role_check as prc
                   ON prc.instance_id = a.instance_id
                   AND prc.project_id = a.project_id
GROUP BY a.instance_id, a.resource_owner, a.project_id, a.app_id, uro.resource_owner, hpgc.granted_org_id,
         prc.project_id
LIMIT 1;
