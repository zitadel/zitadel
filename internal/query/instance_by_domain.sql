with domain as (
	select instance_id from projections.instance_domains
	where domain = $1
), instance_features as (
	select i.*
	from domain d
	join projections.instance_features2 i on d.instance_id = i.instance_id
), features as (
	select instance_id, json_object_agg(
		coalesce(i.key, s.key),
		coalesce(i.value, s.value)
	) features
	from domain d
	cross join projections.system_features s
	full outer join instance_features i using (instance_id, key)
	group by instance_id
)
select
    i.id,
    i.default_org_id,
    i.iam_project_id,
    i.console_client_id,
    i.console_app_id,
    i.default_language,
    s.enable_iframe_embedding,
    s.origins,
	s.enable_impersonation,
    l.audit_log_retention,
    l.block,
	f.features
from domain d
join projections.instances i on i.id = d.instance_id
left join projections.security_policies2 s on i.id = s.instance_id
left join projections.limits l on i.id = l.instance_id
left join features f on i.id = f.instance_id;
