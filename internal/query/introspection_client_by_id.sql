with config as (
		select instance_id, app_id, client_id, client_secret, 'api' as app_type
		from projections.apps7_api_configs
		where instance_id = $1
			and client_id = $2
	union
		select instance_id, app_id, client_id, client_secret, 'oidc' as app_type
		from projections.apps7_oidc_configs
		where instance_id = $1
			and client_id = $2
),
keys as (
	select identifier as client_id, json_object_agg(id, encode(public_key, 'base64')) as public_keys
	from projections.authn_keys2
	where $3 = true -- when argument is false, don't waste time on trying to query for keys.
		and instance_id = $1
		and identifier = $2
		and expiration > current_timestamp
	group by identifier
)
select config.app_id, config.client_id, config.client_secret, config.app_type, apps.project_id, apps.resource_owner, p.project_role_assertion, keys.public_keys
from config
join projections.apps7 apps on apps.id = config.app_id and apps.instance_id = config.instance_id
join projections.projects4 p on p.id = apps.project_id and p.instance_id = $1
left join keys on keys.client_id = config.client_id;
