# SQL initialisation

The sql-files in this folder initialize the ZITADEL database and user. These objects need to exist bevore ZITADEL is able to set and start up.

## files

- 01_user.sql: create the user zitadel uses to connect to the database
- 02_database.sql: create the database for zitadel
- 03_grant_user.sql: grants the user created bevore to have full access to it's database. The user needs full access to the database because zitadel makes ddl/dml on runtime
- 04_eventstore.sql: creates the schema needed for eventsourcing
- 05_projections.sql: creates the schema needed to read the data
- files 06_enable_hash_sharded_indexes.sql and 07_events_table.sql must run in the same session
  - 06_enable_hash_sharded_indexes.sql enables the (hash sharded index)[https://www.cockroachlabs.com/docs/stable/hash-sharded-indexes.html] feature for this session
  - 07_events_table.sql creates the table for eventsourcing
