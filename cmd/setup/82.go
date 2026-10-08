package setup

import (
	"context"
	_ "embed"
	"fmt"

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
	inTxOrderType, err := inTxOrderType(ctx, mig.dbClient)
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, fmt.Sprintf(currentStatesInTxOrder, inTxOrderType))
	return err
}

func (mig *CurrentStatesInTxOrder) String() string {
	return "82_current_states_in_tx_order"
}
