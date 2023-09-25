package setup

import (
	"context"
	"embed"

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

	stmt, err := readStmt(fillNewEventCols, "16", mig.dbClient.Type(), "eventstore_columns.sql")
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *FillPosition) String() string {
	return "16_fill_position"
}
