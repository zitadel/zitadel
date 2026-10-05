package setup

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
)

var (
	//go:embed 77.sql
	stampEventPositionAtInsert string
	//go:embed 77_current_states.sql
	backfillCurrentStatesInTxOrder string
)

// StampEventPositionAtInsert is historical setup 77. It is no longer executed.
// Existing Done events with this name remain the flag that filter_offset already
// stores events2.in_tx_order. New installs never record this name.
type StampEventPositionAtInsert struct {
	dbClient *database.DB
}

func (mig *StampEventPositionAtInsert) Execute(ctx context.Context, _ eventstore.Event) error {
	inTxOrderType, err := inTxOrderType(ctx, mig.dbClient)
	if err != nil {
		return err
	}

	tx, err := mig.dbClient.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt := fmt.Sprintf(stampEventPositionAtInsert, inTxOrderType)
	_, err = tx.ExecContext(ctx, stmt)
	if err == nil {
		_, err = tx.ExecContext(ctx, backfillCurrentStatesInTxOrder)
	}
	if err = database.CloseTransaction(tx, err); err != nil {
		return err
	}

	// close idle connections to prevent them from using the old prepared statement
	for _, conn := range mig.dbClient.Pool.AcquireAllIdle(ctx) {
		logging.OnError(ctx, conn.Conn().Close(ctx)).Debug("failed to close idle connection")
		conn.Release()
	}
	return nil
}

func (mig *StampEventPositionAtInsert) String() string {
	return migration.EventstorePositionClockTimestampStep
}
