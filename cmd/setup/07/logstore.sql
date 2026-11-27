CREATE SCHEMA IF NOT EXISTS logstore;

-- These statements run as the service user itself in every flow
-- ('zitadel init', 'zitadel init schema' and 'setup'), which makes this grant
-- a self-grant. Some managed PostgreSQL services (e.g. Fly.io Managed
-- Postgres) reject GRANT statements entirely, so only execute it when the
-- current role actually differs from the target role.
DO $$
BEGIN
    IF current_user <> '%[1]s' THEN
        EXECUTE format('GRANT ALL ON ALL TABLES IN SCHEMA logstore TO %I', '%[1]s');
    END IF;
END
$$;
