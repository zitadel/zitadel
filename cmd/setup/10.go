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
	//go:embed 10_create_temp_table.sql
	correctCreationDate10CreateTable string
	//go:embed 10_fill_table.sql
	correctCreationDate10FillTable string
	//go:embed 10_update.sql
	correctCreationDate10Update string
)

type CorrectCreationDate struct {
	dbClient  *database.DB
	FailAfter time.Duration
}

func (mig *CorrectCreationDate) Execute(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, mig.FailAfter)
	defer cancel()

	for {
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

			_, err = tx.ExecContext(ctx, correctCreationDate10FillTable)
			if err != nil {
				return err
			}

			res, err := tx.ExecContext(ctx, correctCreationDate10Update)
			if err != nil {
				return err
			}
			affected, _ = res.RowsAffected()
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
