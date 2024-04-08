select c.client_id, a.project_id, p.project_role_assertion
from projections.apps6_oidc_configs c
join projections.apps6 a on a.id = c.app_id and a.instance_id = $1
join projections.projects4 p on p.id = a.project_id and p.instance_id = $1
where c.instance_id = $1
    and c.client_id = $2;
