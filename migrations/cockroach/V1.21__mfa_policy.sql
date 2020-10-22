ALTER TABLE management.login_policies ADD COLUMN force_mfa BOOLEAN;
ALTER TABLE adminapi.login_policies ADD COLUMN force_mfa BOOLEAN;
ALTER TABLE auth.login_policies ADD COLUMN force_mfa BOOLEAN;

ALTER TABLE management.login_policies ADD COLUMN software_mfas SMALLINT ARRAY;
ALTER TABLE adminapi.login_policies ADD COLUMN software_mfas SMALLINT ARRAY;
ALTER TABLE auth.login_policies ADD COLUMN software_mfas SMALLINT ARRAY;

ALTER TABLE management.login_policies ADD COLUMN hardware_mfas SMALLINT ARRAY;
ALTER TABLE adminapi.login_policies ADD COLUMN hardware_mfas SMALLINT ARRAY;
ALTER TABLE auth.login_policies ADD COLUMN hardware_mfas SMALLINT ARRAY;