CREATE TYPE zitadel.settings_type AS ENUM (
    'login',
    'label',
    'password_complexity',
    'password_expiry',
    'domain',
    'lockout',
    'security',
    'organization'
);

CREATE TYPE zitadel.label_state AS ENUM (
    'preview',
    'activated'
);

CREATE TABLE zitadel.settings (
    instance_id TEXT NOT NULL
    , org_id TEXT
    , id TEXT NOT NULL CHECK (id <> '') DEFAULT gen_random_uuid()
    , type zitadel.settings_type NOT NULL
    , label_state zitadel.label_state DEFAULT NULL
    , settings JSONB -- the storage does not really care about what is configured so we store it as json

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id)
    , FOREIGN KEY (instance_id, org_id) REFERENCES zitadel.organizations(instance_id, id)
);

-- CREATE UNIQUE INDEX idx_settings_unique_type ON zitadel.settings (instance_id, org_id, type) NULLS NOT DISTINCT WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_settings_unique_type ON zitadel.settings (instance_id, org_id, type) NULLS NOT DISTINCT WHERE type != 'label';
CREATE UNIQUE INDEX idx_settings_label_unique_type ON zitadel.settings (instance_id, org_id, type, label_state) NULLS NOT DISTINCT WHERE type = 'label';

CREATE INDEX idx_settings_type ON zitadel.settings(instance_id, type, label_state) NULLS NOT DISTINCT;


CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.settings
FOR EACH ROW
WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();
