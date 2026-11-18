package setup

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 82_add_column.sql
	currentStatesInTxOrderAddColumn string
	//go:embed 82_triggers.sql
	currentStatesInTxOrderTriggers string
	//go:embed 82_backfill.sql
	currentStatesInTxOrderBackfill string
)

type CurrentStatesInTxOrder struct {
	dbClient *database.DB
}

func (mig *CurrentStatesInTxOrder) Execute(ctx context.Context, _ eventstore.Event) error {
	inTxOrderType, err := inTxOrderType(ctx, mig.dbClient)
	if err != nil {
		return err
	}
	// Separate Execs so ADD COLUMN's ACCESS EXCLUSIVE lock is not held for the backfill.
	if _, err = mig.dbClient.ExecContext(ctx, fmt.Sprintf(currentStatesInTxOrderAddColumn, inTxOrderType)); err != nil {
		return err
	}
	if _, err = mig.dbClient.ExecContext(ctx, currentStatesInTxOrderTriggers); err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, currentStatesInTxOrderBackfill)
	return err
}

func (mig *CurrentStatesInTxOrder) String() string {
	return "82_current_states_in_tx_order"
}
