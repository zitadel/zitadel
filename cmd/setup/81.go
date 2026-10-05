package setup

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

// CommandsToEventsClockTimestamp replaces eventstore.commands_to_events so
// position is stamped with clock_timestamp() at INSERT. DBs that never ran
// setup 77 still need this; the 77 name is not reused so that Done event stays
// a historical cursor-mode flag. This step does not backfill current_states.
type CommandsToEventsClockTimestamp struct {
	dbClient *database.DB
}

func (mig *CommandsToEventsClockTimestamp) Execute(ctx context.Context, _ eventstore.Event) error {
	inTxOrderType, err := inTxOrderType(ctx, mig.dbClient)
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf(stampEventPositionAtInsert, inTxOrderType)
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	if mig.dbClient.Pool == nil {
		return nil
	}
	// close idle connections to prevent them from using the old prepared statement
	for _, conn := range mig.dbClient.Pool.AcquireAllIdle(ctx) {
		logging.OnError(ctx, conn.Conn().Close(ctx)).Debug("failed to close idle connection")
		conn.Release()
	}
	return nil
}

func (mig *CommandsToEventsClockTimestamp) String() string {
	return "81_eventstore_commands_to_events_clock_timestamp"
}
