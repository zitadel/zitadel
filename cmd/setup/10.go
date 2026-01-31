package setup

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 10/10_create_temp_table.sql
	correctCreationDate10CreateTable string
	//go:embed 10/10_fill_table.sql
	correctCreationDate10FillTable string
	//go:embed 10/10_update.sql
	correctCreationDate10Update string
	//go:embed 10/10_count_wrong_events.sql
	correctCreationDate10CountWrongEvents string
	//go:embed 10/10_empty_table.sql
	correctCreationDate10Truncate string
)

type CorrectCreationDate struct {
	dbClient  *database.DB
	FailAfter time.Duration
}

func (mig *CorrectCreationDate) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	ctx, cancel := context.WithTimeout(ctx, mig.FailAfter)
	defer cancel()

	for i := 0; ; i++ {
		logCtx := logging.With(ctx, "mig", mig.String(), "iteration", i)
		logging.Info(logCtx, "start iteration")
		var affected int64
		err = crdb.ExecuteTx(logCtx, mig.dbClient.DB, nil, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(logCtx, correctCreationDate10CreateTable)
			if err != nil {
				return err
			}
			logging.Debug(logCtx, "temp table created")

			_, err = tx.ExecContext(logCtx, correctCreationDate10Truncate)
			if err != nil {
				return err
			}
			_, err = tx.ExecContext(logCtx, correctCreationDate10FillTable)
			if err != nil {
				return err
			}
			logging.Debug(logCtx, "temp table filled")

			res := tx.QueryRowContext(logCtx, correctCreationDate10CountWrongEvents)
			if err := res.Scan(&affected); err != nil || affected == 0 {
				return err
			}

			_, err = tx.ExecContext(logCtx, correctCreationDate10Update)
			if err != nil {
				return err
			}
			logging.Debug(logCtx, "creation dates updated")
			return nil
		})
		logging.Debug(logCtx, "end iteration")
		if affected == 0 || err != nil {
			return err
		}
	}
}

func (mig *CorrectCreationDate) String() string {
	return "10_correct_creation_date"
}
