CREATE USER queries WITH PASSWORD ${queriespassword};
GRANT SELECT ON DATABASE eventstore TO queries;
