CREATE USER queries WITH PASSWORD ${queriespassword};
GRANT SELECT ON TABLE eventstore.events TO queries;
