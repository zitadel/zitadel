-- replace the first %[1]q with the database in double quotes
-- replace the second \%[2]q with the user in double quotes$
-- For more information see technical advisory 10009 (https://zitadel.com/docs/support/advisory/a10009)
ALTER ROLE %[2]q IN DATABASE %[1]q SET enable_durable_locking_for_serializable = on;