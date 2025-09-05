CREATE TYPE zitadel.domain_validation_type AS ENUM (
    'dns'
    , 'http'
);

CREATE TYPE zitadel.domain_type AS ENUM (
    'custom'
    , 'trusted'
);

CREATE TABLE zitadel.instance_domains(
  instance_id TEXT NOT NULL
  , domain TEXT NOT NULL CHECK (LENGTH(domain) BETWEEN 1 AND 255)
  , is_primary BOOLEAN
  , is_generated BOOLEAN
  , type zitadel.domain_type NOT NULL

  , created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
  , updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL

  , PRIMARY KEY (domain)

  , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id) ON DELETE CASCADE

  , CONSTRAINT primary_cannot_be_trusted CHECK (is_primary IS NULL OR type != 'trusted')
  , CONSTRAINT generated_cannot_be_trusted CHECK (is_generated IS NULL OR type != 'trusted')
  , CONSTRAINT custom_values_set CHECK ((is_primary IS NOT NULL AND is_generated IS NOT NULL) OR type != 'custom')
);

CREATE INDEX idx_instance_domain_instance ON zitadel.instance_domains(instance_id);

CREATE TABLE zitadel.org_domains(
  instance_id TEXT NOT NULL
  , org_id TEXT NOT NULL
  , domain TEXT NOT NULL CHECK (LENGTH(domain) BETWEEN 1 AND 255)
  , is_verified BOOLEAN NOT NULL DEFAULT FALSE
  , is_primary BOOLEAN NOT NULL DEFAULT FALSE
  , validation_type zitadel.domain_validation_type

  , created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
  , updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL

  , PRIMARY KEY (instance_id, org_id, domain)

  , FOREIGN KEY (instance_id, org_id) REFERENCES zitadel.organizations(instance_id, id) ON DELETE CASCADE

  , UNIQUE (instance_id, org_id, domain)
);

CREATE INDEX idx_org_domain ON zitadel.org_domains(instance_id, domain);

-- Trigger to update the updated_at timestamp on instance_domains
CREATE TRIGGER trg_set_updated_at_instance_domains
  BEFORE UPDATE ON zitadel.instance_domains
  FOR EACH ROW
  WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
  EXECUTE FUNCTION zitadel.set_updated_at();

-- Trigger to update the updated_at timestamp on org_domains
CREATE TRIGGER trg_set_updated_at_org_domains
  BEFORE UPDATE ON zitadel.org_domains
  FOR EACH ROW
  WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
  EXECUTE FUNCTION zitadel.set_updated_at();

-- Function to check for already verified org domains
CREATE OR REPLACE FUNCTION zitadel.check_verified_org_domain()
RETURNS TRIGGER AS $$
BEGIN
  -- Check if there's already a verified domain within this instance (excluding the current record being updated)
  IF EXISTS (
    SELECT 1
    FROM zitadel.org_domains
    WHERE instance_id = NEW.instance_id
        AND domain = NEW.domain
        AND is_verified = TRUE
        AND (TG_OP = 'INSERT' OR (org_id != NEW.org_id))
    ) THEN
      RAISE EXCEPTION 'org domain is already taken';
  END IF;
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to enforce verified domain constraint on org_domains
CREATE TRIGGER trg_check_verified_org_domain
  BEFORE INSERT OR UPDATE ON zitadel.org_domains
  FOR EACH ROW
  WHEN (NEW.is_verified IS TRUE)
  EXECUTE FUNCTION zitadel.check_verified_org_domain();

-- Function to ensure only one primary domain per instance in instance_domains
CREATE OR REPLACE FUNCTION zitadel.ensure_single_primary_instance_domain()
RETURNS TRIGGER AS $$
BEGIN
  -- If setting this domain as primary, update all other domains in the same instance to non-primary
  UPDATE zitadel.instance_domains 
  SET is_primary = FALSE, updated_at = NOW()
  WHERE instance_id = NEW.instance_id 
    AND domain != NEW.domain 
    AND is_primary = TRUE
    AND type = 'custom';
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to enforce single primary domain constraint on instance_domains
CREATE TRIGGER trg_ensure_single_primary_instance_domain
  BEFORE INSERT OR UPDATE ON zitadel.instance_domains
  FOR EACH ROW
  WHEN (NEW.is_primary IS TRUE)
  EXECUTE FUNCTION zitadel.ensure_single_primary_instance_domain();

-- Function to ensure only one primary domain per organization in org_domains
CREATE OR REPLACE FUNCTION zitadel.ensure_single_primary_org_domain()
RETURNS TRIGGER AS $$
BEGIN
  -- If setting this domain as primary, update all other domains in the same organization to non-primary
  UPDATE zitadel.org_domains 
  SET is_primary = FALSE, updated_at = NOW()
  WHERE instance_id = NEW.instance_id 
    AND org_id = NEW.org_id
    AND domain != NEW.domain 
    AND is_primary = TRUE;
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to enforce single primary domain constraint on org_domains
CREATE TRIGGER trg_ensure_single_primary_org_domain
  BEFORE INSERT OR UPDATE ON zitadel.org_domains
  FOR EACH ROW
  WHEN (NEW.is_primary IS TRUE)
  EXECUTE FUNCTION zitadel.ensure_single_primary_org_domain();