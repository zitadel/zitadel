WITH removed_domains AS (
    select
        e.instance_id
        , e.payload::json->>'domain' as domain
    from eventstore.events2 e
    where e.instance_id = $1
        and e.aggregate_type = 'org'
        and e.event_type = 'org.domain.removed'
        and not exists (
            select 1
            from eventstore.events2 e2
            where e2.instance_id = e.instance_id
                and e2.aggregate_type = e.aggregate_type
                and e2.aggregate_id = e.aggregate_id
                and e2.payload = e.payload
                and e2.event_type = 'org.domain.verified'
                and e2.sequence > e.sequence
        )
), delete_unique_constraints_dry as (
    delete from eventstore.unique_constraints as uc
        using removed_domains as rd
    where uc.instance_id = rd.instance_id
           and uc.unique_type = 'org_domain'
           and uc.unique_field = rd.domain
           returning *
), delete_unique_constraints as (
    delete from eventstore.unique_constraints as uc
        using delete_unique_constraints_dry as d
    where uc.instance_id = d.instance_id
           and uc.unique_type = d.unique_type
           and uc.unique_field = d.unique_field
)
SELECT
    rd.instance_id
    , STRING_AGG ( uc.unique_field, ',') deleted_constraints
FROM removed_domains rd
    left join delete_unique_constraints_dry uc
        on uc.instance_id = rd.instance_id
group by
    instance_id
;


