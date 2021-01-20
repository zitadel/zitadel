CREATE TABLE eventstore.unique_constraints (
    unique_type TEXT,
	unique_field TEXT,
	PRIMARY KEY (unique_type, unique_field)
);
