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
	//go:embed 76.sql
	stampEventPositionAtInsert string
)

type StampEventPositionAtInsert struct {
	dbClient *database.DB
}

func (mig *StampEventPositionAtInsert) Execute(ctx context.Context, _ eventstore.Event) error {
	inTxOrderType, err := inTxOrderType(ctx, mig.dbClient)
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf(stampEventPositionAtInsert, inTxOrderType)
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	if err != nil {
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
	return "76_eventstore_position_clock_timestamp"
}
