package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 13.sql
	currentProjectionState string
)

type CurrentProjectionState struct {
	dbClient *database.DB
}

func (mig *CurrentProjectionState) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, currentProjectionState)
	return err
}

func (mig *CurrentProjectionState) String() string {
	return "13_current_projection_state"
}
