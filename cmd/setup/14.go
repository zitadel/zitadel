package setup

import (
	"context"
	"embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 14/cockroach/14.sql
	//go:embed 14/postgres/14.sql
	currentProjectionState embed.FS
)

type CurrentProjectionState struct {
	dbClient *database.DB
}

func (mig *CurrentProjectionState) Execute(ctx context.Context) error {
	stmt, err := readStmt(currentProjectionState, "14", mig.dbClient.Type(), "14.sql")
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *CurrentProjectionState) String() string {
	return "14_current_projection_state"
}
