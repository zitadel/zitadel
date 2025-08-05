-- the query is so complex because we accidentally stored unique constraint case sensitive
-- the query checks first if there is a case sensitive match and afterwards if there is a case insensitive match
(instance_id = $%[1]d AND unique_type = $%[2]d AND unique_field = (
    SELECT unique_field from (
    SELECT instance_id, unique_type, unique_field
    FROM eventstore.unique_constraints
    WHERE instance_id = $%[1]d AND unique_type = $%[2]d AND unique_field = $%[3]d
    UNION ALL
    SELECT instance_id, unique_type, unique_field
    FROM eventstore.unique_constraints
    WHERE instance_id = $%[1]d AND unique_type = $%[2]d AND unique_field = LOWER($%[3]d)
    ) AS case_insensitive_constraints LIMIT 1)
)
