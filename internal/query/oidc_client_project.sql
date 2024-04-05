select a.project_id
from projections.apps6_oidc_configs c
join projections.apps6 a on a.id = c.app_id and a.instance_id = c.instance_id
where c.instance_id = $1
    and c.client_id = $2;
