package setup

import (
	"context"
	"embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 14/cockroach/*.sql
	//go:embed 14/postgres/*.sql
	currentProjectionState embed.FS
)

type CurrentProjectionState struct {
	dbClient *database.DB
}

func (mig *CurrentProjectionState) Execute(ctx context.Context) error {
	migrations, err := currentProjectionState.ReadDir("14/" + mig.dbClient.Type())
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		stmt, err := readStmt(currentProjectionState, "14", mig.dbClient.Type(), migration.Name())
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

func (mig *CurrentProjectionState) String() string {
	return "14_current_projection_state"
}
