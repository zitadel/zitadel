package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 11.sql
	currentProjectionState string
)

type CurrentProjectionState struct {
	dbClient *database.DB
}

func (mig *CurrentProjectionState) Execute(ctx context.Context) (err error) {
	tx, err := mig.dbClient.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			logging.OnError(tx.Rollback()).Debug("rollback failed")
			return
		}
		err = tx.Commit()
	}()
	for {
		res, err := tx.ExecContext(ctx, currentProjectionState)
		if err != nil {
			return err
		}
		affected, _ := res.RowsAffected()
		logging.WithFields("count", affected).Info("creation dates changed")
		if affected == 0 {
			return nil
		}
	}
}

func (mig *CurrentProjectionState) String() string {
	return "11_correct_creation_date"
}
