CREATE INDEX CONCURRENTLY IF NOT EXISTS fields_instance_domains_idx
ON eventstore.fields (object_id) INCLUDE (instance_id)
WHERE object_type = 'instance_domain' AND field_name = 'domain';