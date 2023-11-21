with config as (
		select app_id, client_id, client_secret
		from projections.apps6_api_configs
		where instance_id = $1
			and client_id = $2
	union
		select app_id, client_id, client_secret
		from projections.apps6_oidc_configs
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
select config.client_id, config.client_secret, apps.project_id, keys.public_keys from config
join projections.apps6 apps on apps.id = config.app_id
left join keys on keys.client_id = config.client_id;
