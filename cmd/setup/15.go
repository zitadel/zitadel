package setup

import (
	"context"
	"embed"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 15/cockroach/*.sql
	//go:embed 15/postgres/*.sql
	currentProjectionState embed.FS
)

type CurrentProjectionState struct {
	dbClient *database.DB
}

func (mig *CurrentProjectionState) Execute(ctx context.Context, _ eventstore.Event) error {
	migrations, err := currentProjectionState.ReadDir("15/" + mig.dbClient.Type())
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		stmt, err := readStmt(currentProjectionState, "15", mig.dbClient.Type(), migration.Name())
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

func (mig *CurrentProjectionState) String() string {
	return "15_current_projection_state"
}
