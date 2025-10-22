CREATE TABLE zitadel.project_grants(
    instance_id TEXT NOT NULL
    , id TEXT NOT NULL

    , project_id TEXT NOT NULL
    , granting_organization_id TEXT NOT NULL
    , granted_organization_id TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

    , PRIMARY KEY (instance_id, id)

    -- we need both orgs because a project id is unique on org level
    , FOREIGN KEY (instance_id, granting_organization_id) REFERENCES zitadel.organizations(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, granted_organization_id) REFERENCES zitadel.organizations(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, granting_organization_id, project_id) REFERENCES zitadel.projects(instance_id, organization_id, id) ON DELETE CASCADE

    , UNIQUE (instance_id, project_id, granted_organization_id)
);

CREATE TRIGGER trg_set_updated_at_project_grants
  BEFORE UPDATE ON zitadel.project_grants
  FOR EACH ROW
  WHEN (NEW.updated_at IS NULL)
  EXECUTE FUNCTION zitadel.set_updated_at();

CREATE TABLE zitadel.project_grant_roles(
    instance_id TEXT NOT NULL
    , grant_id TEXT NOT NULL
    , project_id TEXT NOT NULL
    , role_key TEXT NOT NULL 
    
    , created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

    , PRIMARY KEY (instance_id, grant_id, role_key)
    , FOREIGN KEY (instance_id, grant_id) REFERENCES zitadel.project_grants(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, project_id, role_key) REFERENCES zitadel.project_roles(instance_id, project_id, key) ON DELETE CASCADE
);
