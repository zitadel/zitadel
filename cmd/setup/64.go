package setup

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 64.sql
	changePushPosition string
)

type ChangePushPosition struct {
	dbClient *database.DB
}

func (mig *ChangePushPosition) Execute(ctx context.Context, _ eventstore.Event) error {
	inTxOrderType, err := inTxOrderType(ctx, mig.dbClient)
	if err != nil {
		return err
	}
	stmt := fmt.Sprintf(changePushPosition, inTxOrderType)
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	for _, conn := range mig.dbClient.Pool.AcquireAllIdle(ctx) {
		logging.OnError(ctx, conn.Conn().Close(ctx)).Debug("failed to close idle connection")
	}
	return nil
}

func (mig *ChangePushPosition) String() string {
	return "64_change_push_position"
}

func inTxOrderType(ctx context.Context, client *database.DB) (typeName string, err error) {
	err = client.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&typeName)
	}, `SELECT data_type FROM information_schema.columns WHERE table_schema = 'eventstore' AND table_name = 'events2' AND column_name = 'in_tx_order'`)
	if err != nil {
		return "", fmt.Errorf("get in_tx_order_type: %w", err)
	}
	return typeName, nil
}
