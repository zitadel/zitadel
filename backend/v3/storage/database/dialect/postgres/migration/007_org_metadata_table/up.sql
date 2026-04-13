CREATE TABLE zitadel.organization_metadata (
    instance_id TEXT NOT NULL
    , organization_id TEXT NOT NULL
    , key TEXT NOT NULL CHECK (key <> '')
    , value BYTEA NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    
    , PRIMARY KEY (instance_id, organization_id, key)
    
    , CONSTRAINT fk_organization_metadata_org FOREIGN KEY (instance_id, organization_id) REFERENCES zitadel.organizations (instance_id, id) ON DELETE CASCADE
);

CREATE INDEX idx_organization_metadata_key ON zitadel.organization_metadata (key);
CREATE INDEX idx_organization_metadata_value ON zitadel.organization_metadata (sha256(value));

-- TODO(adlerhurst): these indexes can currently not be used by Postgres, because of the type conversion
-- the value can be a json but doesn't have to be.
-- CREATE INDEX idx_organization_metadata_value_number ON zitadel.organization_metadata ((value::NUMERIC)) WHERE jsonb_typeof(value) = 'number';
-- CREATE INDEX idx_organization_metadata_value_string ON zitadel.organization_metadata ((value#>>'{}')) WHERE jsonb_typeof(value) = 'string';
-- CREATE INDEX idx_organization_metadata_value_boolean ON zitadel.organization_metadata ((value::BOOLEAN)) WHERE jsonb_typeof(value) = 'boolean';

CREATE TRIGGER trg_set_updated_at_organization_metadata
  BEFORE INSERT OR UPDATE ON zitadel.organization_metadata
  FOR EACH ROW
  WHEN (NEW.updated_at IS NULL)
  EXECUTE FUNCTION zitadel.set_updated_at();