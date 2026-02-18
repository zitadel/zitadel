with features as (
	select instance_id, json_object_agg(
		coalesce(i.key, s.key),
		coalesce(i.value, s.value)
	) features
	from (select $1::text instance_id) x
	cross join projections.system_features4 s
	full outer join projections.instance_features5 i using (key, instance_id)
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
), execution_targets as (
	select instance_id, json_agg(x.execution_targets) as execution_targets from (
		select e.instance_id, json_build_object(
			'execution_id', et.execution_id,
			'target_id', t.id,
			'target_type', t.target_type,
			'endpoint', t.endpoint,
			'timeout', t.timeout,
			'interrupt_on_error', t.interrupt_on_error,
			'signing_key', t.signing_key,
            'payload_type', t.payload_type,
            'encryption_key', encode(k.public_key, 'base64'),
            'encryption_key_id', k.id
		) as execution_targets
		from projections.executions1 e
		join projections.executions1_targets et
			on e.instance_id = et.instance_id
			and e.id = et.execution_id
		join projections.targets2 t
			on et.instance_id = t.instance_id
			and et.target_id = t.id
        left join projections.authn_keys2 k
            on k.instance_id = et.instance_id
            and k.object_id = t.id
            and k.enabled = true
            and (k.expiration IS NULL or k.expiration > now())
		where e.instance_id = $1
		order by et.position asc
	) as x
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
	td.domains as trusted_domains,
	et.execution_targets,
    r.allowed_languages
from projections.instances i
left join projections.security_policies2 s on i.id = s.instance_id
left join projections.limits l on i.id = l.instance_id
left join features f on i.id = f.instance_id
left join external_domains ed on i.id = ed.instance_id
left join trusted_domains td on i.id = td.instance_id
left join execution_targets et on i.id = et.instance_id
left join projections.restrictions2 r on i.id = r.instance_id
where i.id = $1;
