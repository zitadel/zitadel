package setup

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 70.sql
	addEventStoreCommandEnforceOwnerColumn string
)

type AddEventStoreCommandEnforceOwnerColumn struct {
	dbClient *database.DB
}

func (mig *AddEventStoreCommandEnforceOwnerColumn) Execute(ctx context.Context, _ eventstore.Event) error {
	inTxOrderType, err := inTxOrderType(ctx, mig.dbClient)
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf(addEventStoreCommandEnforceOwnerColumn, inTxOrderType)
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	// close idle connections to prevent them from using the old prepared statement with the wrong type
	// and having wrong plan of `eventstore.command2`-type
	for _, conn := range mig.dbClient.Pool.AcquireAllIdle(ctx) {
		logging.OnError(ctx, conn.Conn().Close(ctx)).Debug("failed to close idle connection")
		conn.Release()
	}
	return nil
}

func (mig *AddEventStoreCommandEnforceOwnerColumn) String() string {
	return "70_add_eventstore_command_type_enforce_owner_column"
}
