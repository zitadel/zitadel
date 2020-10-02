-- table manipulations need to be in seperate transactions 
-- (see: https://github.com/cockroachdb/cockroach/issues/13380#issuecomment-306560448)

BEGIN;
ALTER TABLE management.machine_keys DROP COLUMN IF EXISTS public_key;
COMMIT;

BEGIN;
ALTER TABLE management.machine_keys ADD COLUMN public_key BYTES;
COMMIT;

BEGIN;
ALTER TABLE auth.machine_keys DROP COLUMN IF EXISTS public_key;
COMMIT;

BEGIN;
ALTER TABLE auth.machine_keys ADD COLUMN public_key BYTES;
COMMIT;