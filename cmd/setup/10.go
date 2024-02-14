package setup

import (
	"context"
	"database/sql"
	"embed"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 10/10_create_temp_table.sql
	correctCreationDate10CreateTable string
	//go:embed 10/10_fill_table.sql
	correctCreationDate10FillTable string
	//go:embed 10/cockroach/10_update.sql
	//go:embed 10/postgres/10_update.sql
	correctCreationDate10Update embed.FS
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
		logging.WithFields("mig", mig.String(), "iteration", i).Debug("start iteration")
		var affected int64
		err = crdb.ExecuteTx(ctx, mig.dbClient.DB, nil, func(tx *sql.Tx) error {
			if mig.dbClient.Type() == "cockroach" {
				if _, err := tx.Exec("SET experimental_enable_temp_tables=on"); err != nil {
					return err
				}
			}
			_, err := tx.ExecContext(ctx, correctCreationDate10CreateTable)
			if err != nil {
				return err
			}
			logging.WithFields("mig", mig.String(), "iteration", i).Debug("temp table created")

			_, err = tx.ExecContext(ctx, correctCreationDate10Truncate)
			if err != nil {
				return err
			}
			_, err = tx.ExecContext(ctx, correctCreationDate10FillTable)
			if err != nil {
				return err
			}
			logging.WithFields("mig", mig.String(), "iteration", i).Debug("temp table filled")

			res := tx.QueryRowContext(ctx, correctCreationDate10CountWrongEvents)
			if err := res.Scan(&affected); err != nil || affected == 0 {
				return err
			}

			updateStmt, err := readStmt(correctCreationDate10Update, "10", mig.dbClient.Type(), "10_update.sql")
			if err != nil {
				return err
			}
			_, err = tx.ExecContext(ctx, updateStmt)
			if err != nil {
				return err
			}
			logging.WithFields("mig", mig.String(), "iteration", i, "count", affected).Debug("creation dates updated")
			return nil
		})
		logging.WithFields("mig", mig.String(), "iteration", i).Debug("end iteration")
		if affected == 0 || err != nil {
			return err
		}
	}
}

func (mig *CorrectCreationDate) String() string {
	return "10_correct_creation_date"
}
