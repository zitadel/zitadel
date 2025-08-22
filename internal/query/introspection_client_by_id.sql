with config as (
		select instance_id, app_id, client_id, client_secret, 'api' as app_type
		from projections.apps7_api_configs
		where instance_id = $1
			and client_id = $2
	union all
		select instance_id, app_id, client_id, client_secret, 'oidc' as app_type
		from projections.apps7_oidc_configs
		where instance_id = $1
			and client_id = $2
),
keys as (
	select $2::text as client_id, json_object_agg(id, encode(public_key, 'base64')) as public_keys
	from projections.authn_keys2
	where $3 = true -- when argument is false, don't waste time on trying to query for keys.
		and instance_id = $1
		and identifier = $2
		and expiration > current_timestamp
)
select c.app_id, c.client_id, c.client_secret, c.app_type, 
       a.project_id, a.resource_owner, p.project_role_assertion, 
       k.public_keys
from config c
join projections.apps7 a on a.id = c.app_id and a.instance_id = c.instance_id and a.state = 1
join projections.projects4 p on p.id = a.project_id and p.instance_id = c.instance_id and p.state = 1
join projections.orgs1 o on o.id = p.resource_owner and o.instance_id = c.instance_id and o.org_state = 1
left join keys k on k.client_id = c.client_id and $3 = true;
