-- replace the first %[1]s with the database
-- replace the second \%[2]s with the user
DO $$
    DECLARE dblink_exists BOOLEAN;
BEGIN

GRANT ALL ON DATABASE "%[1]s" TO "%[2]s";

SELECT count(*) = 1 INTO dblink_exists FROM pg_extension WHERE extname = 'dblink';

IF NOT dblink_exists THEN
    CREATE EXTENSION dblink;
END IF;

PERFORM dblink_connect('%[3]s');

PERFORM dblink_exec('GRANT ALL ON SCHEMA public TO "%[2]s"');
PERFORM dblink_exec('GRANT ALL ON ALL TABLES IN SCHEMA public TO "%[2]s"');

PERFORM dblink_disconnect();

IF NOT dblink_exists THEN
    DROP EXTENSION dblink;
END IF;

END $$;
