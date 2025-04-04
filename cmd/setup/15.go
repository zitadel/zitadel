package setup

import (
	"context"
	"embed"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 15/*.sql
	currentProjectionState embed.FS
)

type CurrentProjectionState struct {
	dbClient *database.DB
}

func (mig *CurrentProjectionState) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := readStatements(currentProjectionState, "15")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		_, err = mig.dbClient.ExecContext(ctx, stmt.query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mig *CurrentProjectionState) String() string {
	return "15_current_projection_state"
}
