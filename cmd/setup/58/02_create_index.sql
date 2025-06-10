CREATE INDEX CONCURRENTLY IF NOT EXISTS login_names3_policies_is_default_owner_idx ON projections.login_names3_policies (instance_id, is_default, resource_owner) INCLUDE (must_be_domain)
