DELETE FROM zitadel.projections.apps 
WHERE project_id IN (
    SELECT aggregate_id
    FROM eventstore.events 
    WHERE aggregate_type = 'project' AND event_type = 'project.removed'
);
