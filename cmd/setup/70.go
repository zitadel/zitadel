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
	conn, err := mig.dbClient.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := conn.Close()
		logging.OnError(ctx, closeErr).Debug("failed to release connection")
		mig.dbClient.Pool.Reset()
	}()

	inTxOrderType, err := (&ChangePushPosition{dbClient: mig.dbClient}).inTxOrderType(ctx)
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf(addEventStoreCommandEnforceOwnerColumn, inTxOrderType)
	_, err = conn.ExecContext(ctx, stmt)
	return err
}

func (mig *AddEventStoreCommandEnforceOwnerColumn) String() string {
	return "70_add_eventstore_command_type_enforce_owner_column"
}
