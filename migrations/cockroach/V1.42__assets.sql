CREATE TABLE eventstore.assets (
    id TEXT,
	asset TEXT,
	PRIMARY KEY (unique_type, unique_field)
);

GRANT DELETE ON TABLE eventstore.assets to eventstore;
