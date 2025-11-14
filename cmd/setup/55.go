package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 55.sql
	executionHandlerCurrentState string
)

type ExecutionHandlerStart struct {
	dbClient *database.DB
}

func (mig *ExecutionHandlerStart) Execute(ctx context.Context, e eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, executionHandlerCurrentState, e.Sequence(), e.CreatedAt(), e.Position())
	return err
}

func (mig *ExecutionHandlerStart) String() string {
	return "55_execution_handler_start"
}
