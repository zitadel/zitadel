package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 10.sql
	correctCreationDate10 string
)

type CorrectCreationDate struct {
	dbClient *database.DB
}

func (mig *CorrectCreationDate) Execute(ctx context.Context) (err error) {
	tx, err := mig.dbClient.Begin()
	if err != nil {
		return err
	}
	if mig.dbClient.Type() == "cockroach" {
		if _, err := tx.Exec("SET experimental_enable_temp_tables=on"); err != nil {
			return err
		}
	}
	defer func() {
		if err != nil {
			logging.OnError(tx.Rollback()).Debug("rollback failed")
			return
		}
		err = tx.Commit()
	}()
	for {
		res, err := tx.ExecContext(ctx, correctCreationDate10)
		if err != nil {
			return err
		}
		affected, _ := res.RowsAffected()
		logging.WithFields("rows", affected).Info("affected")
		if affected == 0 {
			return nil
		}
	}
}

func (mig *CorrectCreationDate) String() string {
	// TODO: reset to 10
	return "12_correct_creation_date"
}
