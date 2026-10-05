package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 82.sql
	currentStatesInTxOrder string
)

type CurrentStatesInTxOrder struct {
	dbClient *database.DB
}

func (mig *CurrentStatesInTxOrder) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, currentStatesInTxOrder)
	return err
}

func (mig *CurrentStatesInTxOrder) String() string {
	return "82_current_states_in_tx_order"
}
