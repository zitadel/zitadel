WITH recursive rel_tree AS
                   (SELECT instance_id                                                            AS instance_id,
                           id                                                                     AS id,
                           unnest(CASE WHEN "includes" <> '{}' THEN "includes" ELSE '{null}' END) AS include,
                           1                                                                      AS level,
                           array [id]                                                             AS path_info,
                           targets                                                                AS targets
                    FROM projections.executions
                    WHERE includes IS NULL
                      AND instance_id = $1
                    GROUP BY instance_id, id, includes, targets
                    UNION ALL
                    SELECT c.instance_id                   AS instance_id,
                           c.id                            AS id,
                           c.include                       AS include,
                           p.level + 1                     AS level,
                           p.path_info || c.id             AS path_info,
                           array_cat(p.targets, c.targets) AS targets
                    FROM (SELECT instance_id                                                        AS instance_id,
                                 id                                                                 AS id,
                                 unnest(CASE WHEN includes <> '{}' THEN includes ELSE '{null}' END) AS include,
                                 targets                                                            AS targets
                          FROM projections.executions) AS c
                             JOIN rel_tree p ON p.id = c.include AND p.instance_id = c.instance_id)
SELECT e.instance_id,
       e.id,
       target,
       t.target_type,
       t.url,
       t.timeout,
       t.interrupt_on_error
FROM ((SELECT instance_id, id, unnest(targets) AS target
       FROM rel_tree
       WHERE id = ANY (string_to_array($2, ','))
       ORDER BY id DESC
       LIMIT 1)
      UNION
      (SELECT instance_id, id, unnest(targets) AS target
       FROM rel_tree
       WHERE id = ANY (string_to_array($3, ','))
       ORDER BY id DESC
       LIMIT 1)) e
         LEFT JOIN projections.targets1 AS t ON t.instance_id = e.instance_id AND t.id = e.target;
