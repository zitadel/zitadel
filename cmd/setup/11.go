package setup

import (
	"context"
	"database/sql"
	"embed"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 11/11_add_column.sql
	addEventCreatedAt string
	//go:embed 11/11_fetch_events.sql
	fetchCreatedAt string
	//go:embed 11/11_fill_column.sql
	fillCreatedAt string
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
	logging.WithFields("step", "11").Info("ensure creation dates order")
	if err := mig.step10.Execute(ctx); err != nil {
		logging.WithFields("step", "11").WithError(err).Info("ensure creation dates order failed")
		return err
	}
	logging.WithFields("step", "11").Info("ensure creation dates order done")
	logging.WithFields("step", "11").Info("add created_at column")
	_, err := mig.dbClient.ExecContext(ctx, addEventCreatedAt)
	if err != nil {
		logging.WithFields("step", "11").WithError(err).Info("add created_at column failed")
		return err
	}
	logging.WithFields("step", "11").Info("created_at column added")

	createIndex, err := readStmt(createdAtIndexCreateStmt, "11", mig.dbClient.Type(), "create_index.sql")
	if err != nil {
		return err
	}
	logging.WithFields("step", "11").Info("create index")
	_, err = mig.dbClient.ExecContext(ctx, createIndex)
	logging.WithFields("step", "11").WithError(err).Info("create index failed")
	if err != nil {
		return err
	}
	logging.WithFields("step", "11").Info("index created")

	for i := 0; ; i++ {
		logging.WithFields("step", "11", "iteration", i).Info("begin set created_at iteration")
		var count int
		err = crdb.ExecuteTx(ctx, mig.dbClient.DB, nil, func(tx *sql.Tx) error {
			rows, err := tx.Query(fetchCreatedAt, mig.BulkAmount)
			if err != nil {
				return err
			}
			defer rows.Close()

			type date struct {
				instanceID    string
				eventSequence uint64
				creationDate  time.Time
			}
			dates := make([]*date, 0, 20)

			for rows.Next() {
				count++

				d := new(date)
				err = rows.Scan(&d.instanceID, &d.eventSequence, &d.creationDate)
				if err != nil {
					return err
				}
				dates = append(dates, d)

			}
			if err := rows.Err(); err != nil {
				return err
			}

			for _, d := range dates {
				_, err = tx.Exec(fillCreatedAt, d.creationDate, d.instanceID, d.eventSequence)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
		logging.WithFields("step", "11", "iteration", i, "count", count).Info("set created_at iteration done")
		if count < 20 {
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
