CREATE TABLE zitadel.projections.actions (
    id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    resource_owner TEXT,
    action_state SMALLINT,
    sequence BIGINT,

    name TEXT,
    script TEXT,
    timeout BIGINT,
    allowed_to_fail BOOLEAN,

    PRIMARY KEY (id)
);

CREATE TABLE zitadel.projections.flows_actions (
    id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    resource_owner TEXT,
    sequence BIGINT,

    name TEXT,
    script TEXT,
    timeout BIGINT,
    allowed_to_fail BOOLEAN,

    PRIMARY KEY (id)
);

CREATE TABLE zitadel.projections.flows_triggers (
    flow_type SMALLINT,
    trigger_type SMALLINT,
    resource_owner TEXT,
    action_id TEXT,
    trigger_sequence SMALLINT,

    PRIMARY KEY (flow_type, trigger_type, resource_owner, action_id),
    CONSTRAINT fk_action FOREIGN KEY (action_id) REFERENCES zitadel.projections.flows_actions (id) ON DELETE CASCADE
);

CREATE VIEW zitadel.projections.flows_actions_triggers AS (
    SELECT a.id AS action_id,
        a.name,
        a.creation_date,
        a.resource_owner,
        a.sequence,
        a.change_date,
        a.script,
        a.timeout,
        a.allowed_to_fail,
        t.flow_type,
        t.trigger_type,
        t.trigger_sequence
    FROM zitadel.projections.flows_triggers t
           JOIN zitadel.projections.flows_actions a ON t.action_id = a.id
      );
