CREATE TABLE IF NOT EXISTS zitadel.domains(
    id SERIAL PRIMARY KEY,
    instance_id TEXT NOT NULL,
    org_id TEXT,
    domain TEXT NOT NULL CHECK (LENGTH(domain) BETWEEN 1 AND 255),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    validation_type SMALLINT CHECK (validation_type >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT fk_domains_instance FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id) ON DELETE CASCADE,
    CONSTRAINT fk_domains_org FOREIGN KEY (instance_id, org_id) REFERENCES zitadel.organizations(instance_id, id) ON DELETE CASCADE,
    CONSTRAINT domain_unique UNIQUE NULLS NOT DISTINCT (instance_id, org_id, domain) WHERE deleted_at IS NULL
);

CREATE INDEX IF NOT EXISTS idx_domains_instance_id ON zitadel.domains(instance_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_domains_org_id ON zitadel.domains(instance_id, org_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_domains_domain ON zitadel.domains(domain) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_domains_is_primary ON zitadel.domains(is_primary) WHERE deleted_at IS NULL AND is_primary = true;
CREATE INDEX IF NOT EXISTS idx_domains_is_verified ON zitadel.domains(is_verified) WHERE deleted_at IS NULL;

CREATE TRIGGER trigger_set_updated_at_domains
BEFORE UPDATE ON zitadel.domains
FOR EACH ROW
WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
EXECUTE FUNCTION zitadel.set_updated_at();