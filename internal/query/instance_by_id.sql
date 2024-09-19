with features as (
	select instance_id, json_object_agg(
		coalesce(i.key, s.key),
		coalesce(i.value, s.value)
	) features
	from (select $1::text instance_id) x
	cross join projections.system_features s
	full outer join projections.instance_features2 i using (key, instance_id)
	group by instance_id
), external_domains as (
	select instance_id, array_agg(domain) as domains
	from projections.instance_domains
    where instance_id = $1
	group by instance_id
), trusted_domains as (
	select instance_id, array_agg(domain) as domains
	from projections.instance_trusted_domains
    where instance_id = $1
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
	f.features,
    ed.domains as external_domains,
	td.domains as trusted_domains
from projections.instances i
left join projections.security_policies2 s on i.id = s.instance_id
left join projections.limits l on i.id = l.instance_id
left join features f on i.id = f.instance_id
left join external_domains ed on i.id = ed.instance_id
left join trusted_domains td on i.id = td.instance_id
where i.id = $1;
