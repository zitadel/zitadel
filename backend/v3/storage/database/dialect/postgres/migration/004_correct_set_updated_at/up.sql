CREATE OR REPLACE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.instances
FOR EACH ROW
WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();

CREATE OR REPLACE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.organizations
FOR EACH ROW
WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();

CREATE OR REPLACE TRIGGER trg_set_updated_at_instance_domains
  BEFORE UPDATE ON zitadel.instance_domains
  FOR EACH ROW
  WHEN (NEW.updated_at IS NULL)
  EXECUTE FUNCTION zitadel.set_updated_at();

CREATE OR REPLACE TRIGGER trg_set_updated_at_org_domains
  BEFORE UPDATE ON zitadel.org_domains
  FOR EACH ROW
  WHEN (NEW.updated_at IS NULL)
  EXECUTE FUNCTION zitadel.set_updated_at();