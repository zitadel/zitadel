CREATE TYPE zitadel.settings_type AS ENUM (
    'login',
    'label',
    'password_complexity', --4
    'password_expiry', --4
    'branding',
    'domain', -- 3
    'legal_and_support',
    'lockout', -- 3
    'general',
    'security', -- 2
    'organization' -- 4
);

CREATE TABLE zitadel.settings (
    instance_id TEXT NOT NULL
    , org_id TEXT
    , id TEXT NOT NULL CHECK (id <> '') DEFAULT gen_random_uuid()
    , type zitadel.settings_type NOT NULL
    , settings JSONB -- the storage does not really care about what is configured so we store it as json

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id)
    , FOREIGN KEY (instance_id, org_id) REFERENCES zitadel.organizations(instance_id, id)
);

-- CREATE UNIQUE INDEX idx_settings_unique_type ON zitadel.settings (instance_id, org_id, type) NULLS NOT DISTINCT WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_settings_unique_type ON zitadel.settings (instance_id, org_id, type) NULLS NOT DISTINCT;

CREATE INDEX idx_settings_type ON zitadel.settings(instance_id, type);


CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.settings
FOR EACH ROW
WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
EXECUTE FUNCTION zitadel.set_updated_at();
