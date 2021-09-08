CREATE TABLE zitadel.projections.actions (
    id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    resource_owner TEXT,
    action_state SMALLINT,
    sequence BIGINT,

    name TEXT,

    PRIMARY KEY (id)
);

CREATE TABLE zitadel.projections.flows_actions (
    id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    resource_owner TEXT,
    sequence BIGINT,

    name TEXT,

    PRIMARY KEY (id)
);

CREATE TABLE zitadel.projections.flows_triggers (
    flow_type SMALLINT,
    trigger_type SMALLINT,
    resource_owner TEXT,
    action_id TEXT,

    PRIMARY KEY (flow_type, trigger_type, resource_owner, action_id),
    CONSTRAINT fk_action FOREIGN KEY (action_id) REFERENCES projections.flows_actions (id) ON DELETE CASCADE
);
--
-- CREATE VIEW zitadel.projections.flows AS (
--     SELECT o.id AS org_id,
--          o.name AS org_name,
--          o.creation_date,
--          u.owner_id,
--          u.language,
--          u.email,
--          u.first_name,
--          u.last_name,
--          u.gender
--     FROM projections.flows_orgs o
--            JOIN projections.flows_users u ON o.id = u.org_id
--       );

CREATE VIEW zitadel.projections.flows_actions_triggers AS (
    SELECT a.id AS action_id,
        a.name,
        a.creation_date,
        a.resource_owner,
        a.sequence,
        a.change_date,
        a.script,
        t.flow_type,
        t.trigger_type
    FROM zitadel.projections.flows_triggers t
           JOIN zitadel.projections.flows_actions a ON t.action_id = a.id
      );