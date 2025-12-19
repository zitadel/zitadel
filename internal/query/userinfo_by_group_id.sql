with grp as (
	select 
	g.id, 
	g.creation_date, 
	g.change_date,
	g.sequence,
	g.state, 
	g.resource_owner, 
	g.name,
	g.description
	from projections.group g
	where g.id = $1 and g.state = 1 -- only allow active groups
	and g.instance_id = $2
),
-- get all group grants, now for all group_ids from the user record
group_grants as (
    select 
        id, 
        grant_id, 
        state, 
        creation_date, 
        change_date, 
        sequence, 
        group_id, 
        roles, 
        resource_owner, 
        project_id
    from projections.group_grants
    where group_id = $1
      and instance_id = $2
      and project_id = any($3)
      and state = 1
    {{ if . -}}
      and resource_owner = any($4)
    {{- end }}
),
-- get all group metadata, now for all group_ids from the user record
group_metadata as (
    select 
        creation_date, 
        change_date, 
        sequence, 
        group_id, 
        resource_owner, 
        key,
        value
    from projections.group_metadata
    where group_id = $1
      and instance_id = $2
    {{ if . -}}
      and resource_owner = any($4)
    {{- end }}
),
-- filter all orgs we are interested in.
orgs as (
	select id, name, primary_domain, instance_id
	from projections.orgs1
	where id in (
		select resource_owner from group_grants
		union
		select resource_owner from grp
	)
	and instance_id = $2
),
-- join group grants to orgs, projects
groupgrants as (
    select json_agg(row_to_json(r)) as grants 
    from (
        select 
            g.*,
            o.name as org_name, 
            o.primary_domain as org_primary_domain,
            p.name as project_name
        from group_grants g
        left join orgs o on o.id = g.resource_owner
        left join projections.projects4 p on p.id = g.project_id
        where p.instance_id = $2
    ) r
),
-- aggregate group metadata to orgs.
groupmetadata as (
    select json_agg(row_to_json(r)) as metadata
    from (
        select 
            gm.group_id,
            gm.resource_owner,
            gm.key,
            gm.value
        from group_metadata gm
        left join orgs o on o.id = gm.resource_owner
    ) r
)
-- build the final result JSON
select json_build_object(
	'group', (select row_to_json(r) as grp from (select * from grp) r),
	'group_grants', (select grants from groupgrants),
	'group_metadata', (select metadata from groupmetadata)
);