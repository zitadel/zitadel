package setup

import (
	"context"
	_ "embed"
	"fmt"

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
	return err
}

func (mig *AddEventStoreCommandEnforceOwnerColumn) String() string {
	return "70_add_eventstore_command_type_enforce_owner_column"
}
