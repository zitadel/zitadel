with usr as (
	select u.id, u.creation_date, u.change_date, u.sequence, u.state, u.resource_owner, u.username, n.login_name as preferred_login_name
	from projections.users10 u
	left join projections.login_names3 n on u.id = n.user_id and u.instance_id = n.instance_id
	where u.id = $1
	and u.instance_id = $2
	and n.is_primary = true
),
human as (
	select $1 as user_id, row_to_json(r) as human from (
		select first_name, last_name, nick_name, display_name, avatar_key, preferred_language, gender, email, is_email_verified, phone, is_phone_verified
		from projections.users10_humans
		where user_id = $1
		and instance_id = $2
	) r
),
machine as (
	select $1 as user_id, row_to_json(r) as machine from (
		select name, description
		from projections.users10_machines
		where user_id = $1
		and instance_id = $2
	) r
),
-- find the user's metadata
metadata as (
	select json_agg(row_to_json(r)) as metadata from (
		select creation_date, change_date, sequence, resource_owner, key, encode(value, 'base64') as value
		from projections.user_metadata5
		where user_id = $1
		and instance_id = $2
	) r
),
-- get all user grants, needed for the orgs query
user_grants as (
	select id, grant_id, state, creation_date, change_date, sequence, user_id, roles, resource_owner, project_id
	from projections.user_grants3
	where user_id = $1
	and instance_id = $2
	and project_id = any($3)
),
-- filter all orgs we are interested in.
orgs as (
	select id, name, primary_domain
	from projections.orgs1
	where id in (
		select resource_owner from user_grants
		union
		select resource_owner from usr
	)
	and instance_id = $2
),
-- find the user's org
user_org as (
	select row_to_json(r) as organization from (
		select o.id, o.name, o.primary_domain
		from orgs o
		join usr u on o.id = u.resource_owner
	) r
),
-- join user grants to orgs, projects and user
grants as (
	select json_agg(row_to_json(r)) as grants from (
		select g.*,
			o.name as org_name, o.primary_domain as org_primary_domain,
			p.name as project_name, u.resource_owner as user_resource_owner
		from user_grants g
		left join orgs o on o.id = g.resource_owner
		left join projections.projects4 p on p.id = g.project_id
		left join usr u on u.id = g.user_id
		where p.instance_id = $2
	) r
)
-- build the final result JSON
select json_build_object(
	'user', (
		select row_to_json(r) as usr from (
			select u.*, h.human, m.machine
			from usr u
			left join human h on u.id = h.user_id
			left join machine m on u.id = m.user_id
		) r
	),
	'org', (select organization from user_org),
	'metadata', (select metadata from metadata),
	'user_grants', (select grants from grants)
);
