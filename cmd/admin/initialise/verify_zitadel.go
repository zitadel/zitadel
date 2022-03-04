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
	searchTable          = "SELECT table_name FROM [SHOW TABLES] WHERE table_name = $1"
	searchSystemSequence = "SELECT sequence_name FROM [SHOW SEQUENCES] WHERE sequence_name = 'system_seq'"
	searchSchema         = "SELECT schema_name FROM [SHOW SCHEMAS] WHERE schema_name = $1"
	//go:embed sql/06_enable_hash_sharded_indexes.sql
	enableHashShardedIdx string
	//go:embed sql/07_events_table.sql
	createEventsStmt string
	//go:embed sql/05_projections.sql
	createProjectionsStmt string
	//go:embed sql/04_eventstore.sql
	createEventstoreStmt string
	//go:embed sql/08_system_sequence.sql
	createSystemSequenceStmt string
	//go:embed sql/09_unique_constraints_table.sql
	createUniqueConstraints string
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

	if err := verify(db, exists(searchTable, "events"), createEvents); err != nil {
		return err
	}

	if err := verify(db, exists(searchSystemSequence), exec(createSystemSequenceStmt)); err != nil {
		return err
	}

	if err := verify(db, exists(searchTable, "unique_constraints"), exec(createUniqueConstraints)); err != nil {
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
