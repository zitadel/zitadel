WITH latest_groups AS MATERIALIZED (
    SELECT DISTINCT ON (aggregate_id) aggregate_id, owner, event_type, sequence
    FROM eventstore.events2
    WHERE instance_id = $1 AND aggregate_type = 'group' AND event_type IN ('group.added', 'group.removed')
    ORDER BY aggregate_id, sequence DESC
), live_groups AS MATERIALIZED (
    SELECT g.* FROM latest_groups g
    WHERE g.event_type = 'group.added'
        AND NOT EXISTS (SELECT 1 FROM eventstore.events2 e WHERE e.instance_id = $1 AND e.aggregate_type = 'org' AND e.aggregate_id = g.owner AND e.event_type = 'org.removed')
        AND NOT EXISTS (SELECT 1 FROM eventstore.events2 e WHERE e.instance_id = $1 AND e.aggregate_type = 'instance' AND e.event_type = 'instance.removed')
), names AS (
    SELECT DISTINCT ON (e.aggregate_id) e.aggregate_id, g.owner, e.payload->'name' AS value
    FROM eventstore.events2 e JOIN live_groups g ON e.aggregate_id = g.aggregate_id AND e.sequence >= g.sequence
    WHERE e.instance_id = $1 AND e.aggregate_type = 'group'
        AND e.event_type IN ('group.added', 'group.changed') AND e.payload->>'name' IS NOT NULL
    ORDER BY e.aggregate_id, e.sequence DESC
), latest_members AS (
    SELECT DISTINCT ON (e.aggregate_id, u.user_id) e.aggregate_id, g.owner, u.user_id, e.event_type IN ('group.user.added', 'group.users.added') AS added
    FROM eventstore.events2 e JOIN live_groups g ON e.aggregate_id = g.aggregate_id AND e.sequence >= g.sequence
    -- Releases before single-user events wrote batched group.users.* events.
    CROSS JOIN LATERAL jsonb_array_elements_text(CASE WHEN e.event_type IN ('group.users.added', 'group.users.removed')
        THEN e.payload->'userIds' ELSE jsonb_build_array(e.payload->'userId') END) AS u(user_id)
    WHERE e.instance_id = $1 AND e.aggregate_type = 'group'
        AND e.event_type IN ('group.user.added', 'group.user.removed', 'group.users.added', 'group.users.removed')
    ORDER BY e.aggregate_id, u.user_id, e.sequence DESC
), cleared AS (
    DELETE FROM eventstore.fields WHERE instance_id = $1 AND aggregate_type = 'group' AND object_type = 'group' RETURNING 1
), released AS (
    -- Older binaries removed groups without releasing their claims.
    DELETE FROM eventstore.unique_constraints c
    WHERE c.instance_id = $1 AND c.unique_type = 'group_external_id' AND c.claimant_type = 'group'
        AND NOT EXISTS (SELECT 1 FROM live_groups g WHERE g.aggregate_id = c.claimant_id)
    RETURNING 1
), current_fields AS (
    SELECT aggregate_id, owner, aggregate_id AS object_id, 'name' AS field_name, value, false AS should_index FROM names
    UNION ALL
    SELECT m.aggregate_id, m.owner, m.user_id, 'member', to_jsonb(m.user_id), true
    FROM latest_members m
    WHERE m.added
    UNION ALL
    SELECT g.aggregate_id, g.owner, c.unique_field, 'external_id_key', to_jsonb(c.unique_field), false
    FROM eventstore.unique_constraints c JOIN live_groups g ON c.claimant_id = g.aggregate_id
    WHERE c.instance_id = $1 AND c.unique_type = 'group_external_id' AND c.claimant_type = 'group'
)
INSERT INTO eventstore.fields (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, value_must_be_unique, should_index)
SELECT $1, f.owner, 'group', f.aggregate_id, 'group', f.object_id, 1, f.field_name, f.value, false, f.should_index
FROM current_fields f CROSS JOIN (SELECT count(*) FROM cleared) AS cleared_count CROSS JOIN (SELECT count(*) FROM released) AS released_count;
