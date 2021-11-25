ALTER TABLE zitadel.projections.orgs RENAME COLUMN domain TO primary_domain;

DROP VIEW zitadel.projections.flows_actions_triggers;
ALTER TABLE zitadel.projections.flows_triggers DROP CONSTRAINT fk_action;
DROP TABLE zitadel.projections.flows_actions;
DELETE FROM zitadel.projections.current_sequences where projection_name in (
    'zitadel.projections.actions',
    'zitadel.projections.flows_actions',
    'zitadel.projections.flows_triggers');
