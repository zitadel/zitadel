ALTER TABLE management.machine_keys DROP COLUMN IF EXISTS public_key;
ALTER TABLE management.machine_keys ADD COLUMN public_key BYTES;

ALTER TABLE auth.machine_keys DROP COLUMN IF EXISTS public_key;
ALTER TABLE auth.machine_keys ADD COLUMN public_key BYTES;
