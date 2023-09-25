package setup

import (
	"context"
	"embed"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 16/cockroach/*.sql
	//go:embed 16/postgres/*.sql
	fillNewEventCols embed.FS
)

type FillPosition struct {
	step10   *CorrectCreationDate
	dbClient *database.DB
}

func (mig *FillPosition) Execute(ctx context.Context) error {
	// execute step 10 again because events created after the first execution of step 10
	// could still have the wrong ordering of sequences and creation date
	if err := mig.step10.Execute(ctx); err != nil {
		return err
	}

	migrations, err := fillNewEventCols.ReadDir("16/" + mig.dbClient.Type())
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		stmt, err := readStmt(fillNewEventCols, "16", mig.dbClient.Type(), migration.Name())
		if err != nil {
			return err
		}

		logging.WithFields("file", migration.Name(), "migration", mig.String()).Info("execute statement")

		_, err = mig.dbClient.ExecContext(ctx, stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mig *FillPosition) String() string {
	return "16_fill_position"
}
