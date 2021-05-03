CREATE TABLE eventstore.assets (
    id TEXT,
	asset TEXT,
	PRIMARY KEY (id)
);

GRANT DELETE ON TABLE eventstore.assets to eventstore;
