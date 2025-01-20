WITH removed_domains AS (
    select
        e.instance_id
        , e.aggregate_type
        , e.aggregate_id
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
), delete_fields_dry as (
    delete from eventstore.fields as f
        using removed_domains as rd
        where f.instance_id = rd.instance_id
            and f.resource_owner = rd.aggregate_id
            and f.aggregate_type = rd.aggregate_type
            and f.object_type = 'org_domain'
            and f.object_id = rd.domain
            and f.field_name = 'verified'
           returning *
), delete_fields as(
    delete from eventstore.fields as f
    using delete_fields_dry as d
    where f.id = d.id
        and f.instance_id = d.instance_id
        and f.resource_owner = d.aggregate_id
        and f.aggregate_type = d.aggregate_type
        and f.object_type = d.object_type
        and f.object_id = d.object_id
        and f.field_name = d.field_name
)
SELECT
    rd.instance_id
    , rd.aggregate_type
    , rd.aggregate_id
    , STRING_AGG ( f.object_id, ',') deleted_fields
FROM removed_domains rd
    left join delete_fields_dry f
      on f.instance_id = rd.instance_id
          and f.resource_owner = rd.aggregate_id
          and f.aggregate_type = rd.aggregate_type
group by
    instance_id, aggregate_type, aggregate_id
;