ALTER TABLE management.login_policies ADD COLUMN force_mfa BOOLEAN;
ALTER TABLE adminapi.login_policies ADD COLUMN force_mfa BOOLEAN;
ALTER TABLE auth.login_policies ADD COLUMN force_mfa BOOLEAN;

ALTER TABLE management.login_policies ADD COLUMN second_factors SMALLINT ARRAY;
ALTER TABLE adminapi.login_policies ADD COLUMN second_factors SMALLINT ARRAY;
ALTER TABLE auth.login_policies ADD COLUMN second_factors SMALLINT ARRAY;

ALTER TABLE management.login_policies ADD COLUMN multi_factors SMALLINT ARRAY;
ALTER TABLE adminapi.login_policies ADD COLUMN multi_factors SMALLINT ARRAY;
ALTER TABLE auth.login_policies ADD COLUMN multi_factors SMALLINT ARRAY;


ALTER TABLE auth.user_sessions RENAME COLUMN mfa_software_verification TO second_factor_verification;
ALTER TABLE auth.user_sessions RENAME COLUMN mfa_software_verification_type TO second_factor_verification_type;
ALTER TABLE auth.user_sessions RENAME COLUMN mfa_hardware_verification TO multi_factor_verification;
ALTER TABLE auth.user_sessions RENAME COLUMN mfa_hardware_verification_type TO multi_factor_verification_type;




