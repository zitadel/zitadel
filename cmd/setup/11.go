package setup

import (
	"context"
	"database/sql"
	"embed"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 11/11_add_column.sql
	addEventCreatedAt string
	//go:embed 11/11_update_events.sql
	setCreatedAt string
	//go:embed 11/11_set_column.sql
	setCreatedAtDetails string
	//go:embed 11/postgres/create_index.sql
	//go:embed 11/cockroach/create_index.sql
	createdAtIndexCreateStmt embed.FS
	//go:embed 11/postgres/drop_index.sql
	//go:embed 11/cockroach/drop_index.sql
	createdAtIndexDropStmt embed.FS
)

type AddEventCreatedAt struct {
	BulkAmount int
	step10     *CorrectCreationDate
	dbClient   *database.DB
}

func (mig *AddEventCreatedAt) Execute(ctx context.Context) error {
	// execute step 10 again because events created after the first execution of step 10
	// could still have the wrong ordering of sequences and creation date
	if err := mig.step10.Execute(ctx); err != nil {
		return err
	}
	_, err := mig.dbClient.ExecContext(ctx, addEventCreatedAt)
	if err != nil {
		return err
	}

	createIndex, err := readStmt(createdAtIndexCreateStmt, "11", mig.dbClient.Type(), "create_index.sql")
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, createIndex)
	if err != nil {
		return err
	}

	for i := 0; ; i++ {
		var affected int64
		err = crdb.ExecuteTx(ctx, mig.dbClient.DB, nil, func(tx *sql.Tx) error {
			res, err := tx.Exec(setCreatedAt, mig.BulkAmount)
			if err != nil {
				return err
			}

			affected, _ = res.RowsAffected()
			return nil
		})
		if err != nil {
			return err
		}
		logging.WithFields("step", "11", "iteration", i, "affected", affected).Info("set created_at iteration done")
		if affected < int64(mig.BulkAmount) {
			break
		}
	}

	logging.Info("set details")
	_, err = mig.dbClient.ExecContext(ctx, setCreatedAtDetails)
	if err != nil {
		return err
	}

	dropIndex, err := readStmt(createdAtIndexDropStmt, "11", mig.dbClient.Type(), "drop_index.sql")
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, dropIndex)

	return err
}

func (mig *AddEventCreatedAt) String() string {
	return "11_event_created_at"
}
