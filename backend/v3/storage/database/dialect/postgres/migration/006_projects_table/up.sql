CREATE TYPE zitadel.project_state AS ENUM (
    'active',
    'inactive'
);

CREATE TABLE zitadel.projects(
    instance_id TEXT NOT NULL
    , organization_id TEXT NOT NULL
    , id TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

    , name TEXT NOT NULL CHECK(LENGTH(name) > 0)
    , state zitadel.project_state NOT NULL
    -- API: project_role_assertion
    , should_assert_role BOOLEAN NOT NULL DEFAULT FALSE
    -- API: authorization_required
    , is_authorization_required BOOLEAN NOT NULL DEFAULT FALSE
    -- API: project_access_required
    , is_project_access_required BOOLEAN NOT NULL DEFAULT FALSE
    --API: private_labeling_setting
    , used_labeling_setting_owner SMALLINT

    , PRIMARY KEY (instance_id, id)
    , UNIQUE (instance_id, organization_id, id)
    , FOREIGN KEY (instance_id, organization_id) REFERENCES zitadel.organizations(instance_id, id) ON DELETE CASCADE
);

CREATE TABLE zitadel.project_roles(
    instance_id TEXT NOT NULL
    , organization_id TEXT NOT NULL
    , project_id TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

    -- API: role_key
    , key TEXT NOT NULL CHECK(LENGTH(key) > 0)
    -- API: display_name
    , display_name TEXT NOT NULL CHECK(LENGTH(display_name) > 0)
    -- API: group
    -- group is a reserved keyword in PostgreSQL
    , role_group TEXT

    , PRIMARY KEY (instance_id, project_id, key)
    , UNIQUE (instance_id, organization_id, project_id, key)
    , FOREIGN KEY (instance_id, organization_id, project_id) REFERENCES zitadel.projects(instance_id, organization_id, id) ON DELETE CASCADE
);

CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.projects
FOR EACH ROW
WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();

CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.project_roles
FOR EACH ROW
WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();
