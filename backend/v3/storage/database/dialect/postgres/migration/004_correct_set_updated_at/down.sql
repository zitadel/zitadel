CREATE OR REPLACE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.instances
FOR EACH ROW
WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
EXECUTE FUNCTION zitadel.set_updated_at();

CREATE OR REPLACE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.organizations
FOR EACH ROW
WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
EXECUTE FUNCTION zitadel.set_updated_at();

CREATE OR REPLACE TRIGGER trg_set_updated_at_instance_domains
  BEFORE UPDATE ON zitadel.instance_domains
  FOR EACH ROW
  WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
  EXECUTE FUNCTION zitadel.set_updated_at();

CREATE OR REPLACE TRIGGER trg_set_updated_at_org_domains
  BEFORE UPDATE ON zitadel.org_domains
  FOR EACH ROW
  WHEN (OLD.updated_at IS NOT DISTINCT FROM NEW.updated_at)
  EXECUTE FUNCTION zitadel.set_updated_at();