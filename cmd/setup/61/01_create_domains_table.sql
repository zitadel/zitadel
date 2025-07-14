CREATE TABLE IF NOT EXISTS zitadel.domains(
  instance_id TEXT NOT NULL
  , org_id TEXT
  , domain TEXT NOT NULL CHECK (LENGTH(domain) BETWEEN 1 AND 255)
  , is_verified BOOLEAN NOT NULL DEFAULT FALSE
  , is_primary BOOLEAN NOT NULL DEFAULT FALSE
  -- TODO make validation_type enum
  , validation_type SMALLINT CHECK (validation_type >= 0)

  , created_at TIMESTAMP DEFAULT NOW()
  , updated_at TIMESTAMP DEFAULT NOW()
  , deleted_at TIMESTAMP DEFAULT NULL

  , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id) ON DELETE CASCADE
  , FOREIGN KEY (instance_id, org_id) REFERENCES zitadel.organizations(instance_id, id) ON DELETE CASCADE

  , CONSTRAINT domain_unique UNIQUE NULLS NOT DISTINCT (instance_id, org_id, domain) WHERE deleted_at IS NULL
);

CREATE TRIGGER IF NOT EXISTS trigger_set_updated_at
BEFORE UPDATE ON zitadel.domains
FOR EACH ROW
WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
EXECUTE FUNCTION zitadel.set_updated_at();