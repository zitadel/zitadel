CREATE INDEX CONCURRENTLY IF NOT EXISTS unique_constraints_owners_gin ON eventstore.unique_constraints USING GIN (owners);
