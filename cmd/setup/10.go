package setup

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
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

func (mig *CorrectCreationDate) Execute(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, mig.FailAfter)
	defer cancel()

	for {
		affected := int64(0)
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

			_, err = tx.ExecContext(ctx, correctCreationDate10Truncate)
			if err != nil {
				return err
			}
			_, err = tx.ExecContext(ctx, correctCreationDate10FillTable)
			if err != nil {
				return err
			}

			res := tx.QueryRowContext(ctx, correctCreationDate10CountWrongEvents)
			if err := res.Scan(&affected); err != nil || affected == 0 {
				return err
			}

			_, err = tx.ExecContext(ctx, correctCreationDate10Update)
			if err != nil {
				return err
			}
			logging.WithFields("count", affected).Info("creation dates changed")
			return nil
		})
		if affected == 0 || err != nil {
			return err
		}
	}
}

func (mig *CorrectCreationDate) String() string {
	return "10_correct_creation_date"
}
