select c.instance_id,
       c.app_id,
       a.state,
       c.entity_id,
       c.metadata,
       c.metadata_url,
       a.project_id,
       p.project_role_assertion,
       c.login_version,
       c.login_base_uri
from projections.apps7_saml_configs c
         join projections.apps7 a
              on a.id = c.app_id and a.instance_id = c.instance_id and a.state = 1
         join projections.projects4 p
              on p.id = a.project_id and p.instance_id = a.instance_id and p.state = 1
         join projections.orgs1 o
              on o.id = p.resource_owner and o.instance_id = c.instance_id and o.org_state = 1
where c.instance_id = $1
  and c.entity_id = $2
