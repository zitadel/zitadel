package initialise

import (
	"database/sql"
	_ "embed"

	"github.com/caos/zitadel/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	eventstoreSchema  = "eventstore"
	projectionsSchema = "projections"
)

var (
	searchEventsTable = "SELECT table_name FROM [SHOW TABLES] WHERE table_name = 'events'"
	searchSchema      = "SELECT schema_name FROM [SHOW SCHEMAS] WHERE schema_name = $1"
	//go:embed sql/enable_hash_sharded_indexes.sql
	enableHashShardedIdx string
	//go:embed sql/events_table.sql
	createEventsStmt string
	//go:embed sql/projections.sql
	createProjectionsStmt string
	//go:embed sql/eventstore.sql
	createEventstoreStmt string
)

func newZitadel() *cobra.Command {
	return &cobra.Command{
		Use:   "zitadel",
		Short: "initialize ZITADEL internas",
		Long: `initialize ZITADEL internas.

Prereqesits:
- cockroachdb with user and database
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := new(Config)
			if err := viper.Unmarshal(config); err != nil {
				return err
			}
			return verifyZitadel(config.Database)
		},
	}
}

func verifyZitadel(config database.Config) error {
	db, err := database.Connect(config)
	if err != nil {
		return err
	}

	if err := verify(db, exists(searchSchema, projectionsSchema), exec(createProjectionsStmt)); err != nil {
		return err
	}

	if err := verify(db, exists(searchSchema, eventstoreSchema), exec(createEventstoreStmt)); err != nil {
		return err
	}

	if err := verify(db, exists(searchSchema, projectionsSchema), createEvents); err != nil {
		return err
	}

	return db.Close()
}

func createEvents(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec(enableHashShardedIdx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec(createEventsStmt); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
