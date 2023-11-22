WITH casesensitive as (
    SELECT instance_id, unique_type, lower(unique_field)
    FROM eventstore.unique_constraints
    GROUP BY instance_id, unique_type, lower(unique_field)
    HAVING count(unique_field) < 2
)
UPDATE eventstore.unique_constraints c
    SET unique_field = casesensitive.lower
    FROM casesensitive
    WHERE c.instance_id = casesensitive.instance_id
        AND c.unique_type = casesensitive.unique_type
        AND lower(c.unique_field) = casesensitive.lower
        AND c.unique_field <> casesensitive.lower;