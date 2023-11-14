with usr as (
	select id, creation_date, change_date, sequence, state, resource_owner, username
	from projections.users9 u
	where id = $1
	and instance_id = $2
),
human as (
	select $1 as user_id, row_to_json(r) as human from (
		select first_name, last_name, nick_name, display_name, avatar_key, email, is_email_verified, phone, is_phone_verified
		from projections.users9_humans
		where user_id = $1
		and instance_id = $2
	) r
),
machine as (
	select $1 as user_id, row_to_json(r) as machine from (
		select name, description
		from projections.users9_machines
		where user_id = $1
		and instance_id = $2
	) r
),
metadata as (
	select json_agg(row_to_json(r)) as metadata from (
		select creation_date, change_date, sequence, resource_owner, key, encode(value, 'base64') as value
		from projections.user_metadata4
		where user_id = $1
		and instance_id = $2
	) r
),
org as (
	select row_to_json(r) as organization from (
		select name, primary_domain
		from projections.orgs1 o
		join usr u on o.id = u.resource_owner
		where instance_id = $2
	) r
)
select json_build_object(
	'user', (
		select row_to_json(r) as usr from (
			select u.*, h.human, m.machine
			from usr u
			left join human h on u.id = h.user_id
			left join machine m on u.id = m.user_id
		) r
	),
	'org', (select organization from org),
	'metadata', (select metadata from metadata)
);
