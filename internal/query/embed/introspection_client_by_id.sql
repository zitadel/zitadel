with config as (
		select app_id, client_id, client_secret
		from projections.apps5_api_configs
		where instance_id = $1
			and client_id = $2
	union
		select app_id, client_id, client_secret
		from projections.apps5_oidc_configs
		where instance_id = $1
			and client_id = $2
),
keys as (
	select identifier as client_id, json_object_agg(id, public_key) as public_keys
	from projections.authn_keys2
	where $3 = true -- when argument is false, don't waste time on trying to query for keys.
		and instance_id = $1
		and identifier = $2
		and expiration > current_timestamp
	group by identifier
)
select apps.project_id, config.client_secret, keys.public_keys from config
join projections.apps5 apps on apps.id = config.app_id
left join keys on keys.client_id = config.client_id
where apps.owner_removed = false;
