CREATE TABLE IF NOT EXISTS eventstore.unique_constraints (
    instance_id TEXT,
    unique_type TEXT,
    unique_field TEXT,
    claimant_type TEXT,
    claimant_id TEXT,
    PRIMARY KEY (instance_id, unique_type, unique_field)
);
