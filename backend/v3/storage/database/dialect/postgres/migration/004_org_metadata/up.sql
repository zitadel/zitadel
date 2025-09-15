CREATE TABLE zitadel.org_metadata (
    instance_id TEXT NOT NULL
    , org_id TEXT NOT NULL
    , key TEXT NOT NULL
    , value JSONB NOT NULL

    , created_at TIMESTAMPTZ NOT NULL
    , updated_at TIMESTAMPTZ NOT NULL
    
    , PRIMARY KEY (instance_id, org_id, key)
    
    , CONSTRAINT fk_org_metadata_org FOREIGN KEY (instance_id, org_id) REFERENCES zitadel.organizations (instance_id, id) ON DELETE CASCADE
);

CREATE INDEX idx_org_metadata_value_number ON zitadel.org_metadata ((value::NUMERIC)) WHERE jsonb_typeof(value) = 'number';
CREATE INDEX idx_org_metadata_value_string ON zitadel.org_metadata ((value#>>'{}')) WHERE jsonb_typeof(value) = 'string';
CREATE INDEX idx_org_metadata_value_boolean ON zitadel.org_metadata ((value::BOOLEAN)) WHERE jsonb_typeof(value) = 'boolean';

CREATE TRIGGER trg_set_updated_at_org_metadata
  BEFORE UPDATE ON zitadel.org_metadata
  FOR EACH ROW
  WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
  EXECUTE FUNCTION zitadel.set_updated_at();