package setup

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

//go:embed 77.sql
var eventStatementPosition string

type EventStatementPosition struct{ dbClient *database.DB }

func (mig *EventStatementPosition) Execute(ctx context.Context, _ eventstore.Event) error {
	typ, err := inTxOrderType(ctx, mig.dbClient)
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, fmt.Sprintf(eventStatementPosition, typ))
	return err
}

func (*EventStatementPosition) String() string { return "77_event_statement_position" }
