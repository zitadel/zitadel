DROP TABLE IF EXISTS zitadel.instance_domains;
DROP TABLE IF EXISTS zitadel.org_domains;
DROP TYPE IF EXISTS zitadel.domain_type;
DROP TYPE IF EXISTS zitadel.domain_validation_type;
DROP FUNCTION IF EXISTS zitadel.check_verified_org_domain();
DROP FUNCTION IF EXISTS zitadel.ensure_single_primary_instance_domain();
