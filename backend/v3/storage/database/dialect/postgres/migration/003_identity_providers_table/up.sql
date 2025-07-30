CREATE TYPE zitadel.idp_state AS ENUM (
    'active',
    'inactive'
);

CREATE TYPE zitadel.idp_type AS ENUM (
    'oidc',
    'jwt',
    'oauth',
    'saml',
    'ldap',
    'github',
    'google',
    'microsoft',
    'apple'
);

CREATE TABLE zitadel.identity_providers (
    instance_id TEXT NOT NULL
    , org_id TEXT
    , id TEXT NOT NULL CHECK (id <> '')
    , state zitadel.idp_state NOT NULL DEFAULT 'active'
    , name TEXT NOT NULL CHECK (name <> '')
    , type zitadel.idp_type -- NOT NULL
    , auto_register BOOLEAN NOT NULL DEFAULT TRUE
    , allow_creation BOOLEAN NOT NULL DEFAULT TRUE
    , allow_auto_creation BOOLEAN NOT NULL DEFAULT TRUE
    , allow_auto_update BOOLEAN NOT NULL DEFAULT TRUE
    , allow_linking BOOLEAN NOT NULL DEFAULT TRUE
    , allow_auto_linking BOOLEAN NOT NULL DEFAULT TRUE
    , styling_type SMALLINT
    , payload JSONB
    
    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    
    , PRIMARY KEY (instance_id, id)
    , CONSTRAINT identity_providers_id_unique UNIQUE NULLS NOT DISTINCT (instance_id, org_id, id)
    , CONSTRAINT identity_providers_name_unique UNIQUE NULLS NOT DISTINCT (instance_id, org_id, name)
    , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id)
    , FOREIGN KEY (instance_id, org_id) REFERENCES zitadel.organizations(instance_id, id)
);

-- CREATE INDEX idx_identity_providers_org_id ON identity_providers(instance_id, org_id) WHERE org_id IS NOT NULL;
CREATE INDEX idx_identity_providers_state ON zitadel.identity_providers(instance_id, state);
CREATE INDEX idx_identity_providers_type ON zitadel.identity_providers(instance_id, type);
-- CREATE INDEX idx_identity_providers_created_at ON identity_providers(created_at);
-- CREATE INDEX idx_identity_providers_deleted_at ON identity_providers(deleted_at) WHERE deleted_at IS NOT NULL;


CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.identity_providers
FOR EACH ROW
WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
EXECUTE FUNCTION zitadel.set_updated_at();
