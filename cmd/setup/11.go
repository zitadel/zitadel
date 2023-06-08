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

	for {
		var count int
		err = crdb.ExecuteTx(ctx, mig.dbClient.DB, nil, func(tx *sql.Tx) error {
			rows, err := tx.Query(fetchCreatedAt, mig.BulkAmount)
			if err != nil {
				return err
			}
			defer rows.Close()

			data := make(map[string]time.Time, 20)
			for rows.Next() {
				count++
				var (
					id           string
					creationDate time.Time
				)

				err = rows.Scan(&id, &creationDate)
				if err != nil {
					return err
				}

				data[id] = creationDate

			}
			if err := rows.Err(); err != nil {
				return err
			}

			for id, creationDate := range data {
				_, err = tx.Exec(fillCreatedAt, creationDate, id)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
		logging.WithFields("count", count).Info("creation dates set")
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
