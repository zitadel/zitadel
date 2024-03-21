with recursive rel_tree as (select instance_id,
                                   id,
                                   UNNEST(CASE WHEN "includes" <> '{}' THEN "includes" ELSE '{null}' END) as include,
                                   1                                                                      as level,
                                   array [id]                                                             as path_info,
                                   targets
                            from projections.executions
                            where includes IS NULL
                              and instance_id = $1
                            group by instance_id, id, includes, targets
                            union all
                            select c.instance_id,
                                   c.id,
                                   c.include,
                                   p.level + 1,
                                   p.path_info || c.id,
                                   array_cat(p.targets, c.targets)
                            from (select instance_id,
                                         id,
                                         UNNEST(CASE WHEN includes <> '{}' THEN includes ELSE '{null}' END) as include,
                                         targets
                                  from projections.executions) as c
                                     join rel_tree p on p.id = c.include and p.instance_id = c.instance_id)
select id, array_agg(target) as targets
from (SELECT id,
             unnest(targets) AS target
      FROM rel_tree) x
where id = $2
group by id;
