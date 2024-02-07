with features as (
	select instance_id, json_object_agg(
		coalesce(i.key, s.key),
		coalesce(i.value, s.value)
	) features
	from (select $1 instance_id) x
	cross join projections.system_features s
	full outer join projections.instance_features i using (key, instance_id)
	group by instance_id
)
select
    i.id,
    i.default_org_id,
    i.iam_project_id,
    i.console_client_id,
    i.console_app_id,
    i.default_language,
    s.enabled,
    s.origins,
    l.audit_log_retention,
    l.block,
	f.features
from projections.instances i
left join projections.security_policies s on i.id = s.instance_id
left join projections.limits l on i.id = l.instance_id
left join features f on i.id = f.instance_id
where i.id = $1;
