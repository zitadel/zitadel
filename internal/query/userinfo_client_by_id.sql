select a.project_id, p.project_role_assertion
from projections.apps7_oidc_configs c
join projections.apps7 a on a.id = c.app_id and a.instance_id = c.instance_id
join projections.projects4 p on p.id = a.project_id and p.instance_id = a.instance_id
where c.instance_id = $1
    and c.client_id = $2;
