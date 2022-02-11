package initialise

import (
	"database/sql"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/database"
)

const (
	eventstoreSchema  = "eventstore"
	projectionsSchema = "projections"
	eventsTable       = "events"
)

func prepareZitadel(config database.Config) error {
	db, err := database.Connect(config)
	if err != nil {
		return err
	}

	logging.Info("verify projections schema")
	if err := verifySchema(db, config, projectionsSchema); err != nil {
		return err
	}

	logging.Info("verify eventstore schema")
	if err := verifySchema(db, config, eventstoreSchema); err != nil {
		return err
	}

	logging.Info("verify events table")
	if err := verifyEvents(db, config); err != nil {
		return err
	}

	return db.Close()
}

func verifySchema(db *sql.DB, config database.Config, schema string) error {
	exists, err := existsSchema(db, config, schema)
	if exists || err != nil {
		return err
	}
	return createSchema(db, config, schema)
}

func existsSchema(db *sql.DB, config database.Config, schema string) (exists bool, err error) {
	row := db.QueryRow("SELECT EXISTS(SELECT schema_name FROM [SHOW SCHEMAS] WHERE schema_name = $1)", schema)
	err = row.Scan(&exists)
	return exists, err
}

func createSchema(db *sql.DB, config database.Config, schema string) error {
	_, err := db.Exec("CREATE SCHEMA " + schema)
	return err
}

func verifyEvents(db *sql.DB, config database.Config) error {
	exists, err := existsEvents(db, config)
	if exists || err != nil {
		return err
	}
	return createEvents(db, config)
}

func existsEvents(db *sql.DB, config database.Config) (exists bool, err error) {
	row := db.QueryRow("SELECT EXISTS(SELECT table_name FROM [SHOW TABLES] WHERE table_name = $1)", eventsTable)
	err = row.Scan(&exists)
	return exists, err
}

func createEvents(db *sql.DB, config database.Config) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec("SET experimental_enable_hash_sharded_indexes = on"); err != nil {
		tx.Rollback()
		return err
	}

	stmt := `CREATE TABLE eventstore.events (
	id UUID DEFAULT gen_random_uuid()
	, event_type TEXT NOT NULL
	, aggregate_type TEXT NOT NULL
	, aggregate_id TEXT NOT NULL
	, aggregate_version TEXT NOT NULL
	, event_sequence BIGINT NOT NULL
	, previous_aggregate_sequence BIGINT
	, previous_aggregate_type_sequence INT8
	, creation_date TIMESTAMPTZ NOT NULL DEFAULT now()
	, event_data JSONB
	, editor_user TEXT NOT NULL 
	, editor_service TEXT NOT NULL
	, resource_owner TEXT NOT NULL

	, PRIMARY KEY (event_sequence DESC) USING HASH WITH BUCKET_COUNT = 10
	, INDEX agg_type_agg_id (aggregate_type, aggregate_id)
	, INDEX agg_type (aggregate_type)
	, INDEX agg_type_seq (aggregate_type, event_sequence DESC) 
		STORING (id, event_type, aggregate_id, aggregate_version, previous_aggregate_sequence, creation_date, event_data, editor_user, editor_service, resource_owner, previous_aggregate_type_sequence)
	, INDEX changes_idx (aggregate_type, aggregate_id, creation_date) USING HASH WITH BUCKET_COUNT = 10
	, INDEX max_sequence (aggregate_type, aggregate_id, event_sequence DESC)
	, CONSTRAINT previous_sequence_unique UNIQUE (previous_aggregate_sequence DESC)
	, CONSTRAINT prev_agg_type_seq_unique UNIQUE(previous_aggregate_type_sequence)
)`
	if _, err = tx.Exec(stmt); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
