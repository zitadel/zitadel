CREATE TYPE zitadel.authorization_state AS ENUM (
    'active',
    'inactive'
    );

CREATE TABLE zitadel.authorizations
(
    instance_id TEXT                        NOT NULL,
    id          TEXT                        NOT NULL CHECK ( id <> '' ),
    state       zitadel.authorization_state NOT NULL DEFAULT 'active',
    project_id  TEXT                        NOT NULL,
    grant_id    TEXT,
    user_id     TEXT                        NOT NULL,
    created_at  TIMESTAMPTZ                 NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ                 NOT NULL DEFAULT NOW(),
    PRIMARY KEY (instance_id, id),
    FOREIGN KEY (instance_id, project_id) REFERENCES zitadel.projects (instance_id, id) ON DELETE CASCADE,
    FOREIGN KEY (instance_id, grant_id) REFERENCES zitadel.project_grants (instance_id, id) ON DELETE CASCADE
);

CREATE TABLE zitadel.authorization_roles
(
    instance_id      TEXT NOT NULL,
    authorization_id TEXT NOT NULL,
    role_key         TEXT NOT NULL CHECK ( role_key <> '' ),
    project_id       TEXT NOT NULL,
    grant_id         TEXT,
    PRIMARY KEY (instance_id, authorization_id, role_key),
    FOREIGN KEY (instance_id, authorization_id) REFERENCES zitadel.authorizations (instance_id, id) ON DELETE CASCADE,
    FOREIGN KEY (instance_id, project_id, role_key) REFERENCES zitadel.project_roles (instance_id, project_id, key) ON DELETE CASCADE,
    FOREIGN KEY (instance_id, grant_id, role_key) REFERENCES zitadel.project_grant_roles (instance_id, grant_id, key) ON DELETE CASCADE
);

CREATE TRIGGER trigger_set_updated_at
    BEFORE UPDATE
    ON zitadel.authorizations
    FOR EACH ROW
    WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();