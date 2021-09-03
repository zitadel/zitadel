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
