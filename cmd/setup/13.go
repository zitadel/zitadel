package setup

import (
	"context"
	"embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 13/cockroach/*.sql
	//go:embed 13/postgres/*.sql
	changeEvents embed.FS
)

type ChangeEvents struct {
	dbClient *database.DB
}

func (mig *ChangeEvents) Execute(ctx context.Context) error {
	migrations, err := changeEvents.ReadDir("13/" + mig.dbClient.Type())
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		stmt, err := readStmt(changeEvents, "13", mig.dbClient.Type(), migration.Name())
		if err != nil {
			return err
		}
		_, err = mig.dbClient.ExecContext(ctx, stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mig *ChangeEvents) String() string {
	return "13_events_push"
}
