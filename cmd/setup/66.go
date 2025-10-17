package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 66.sql
	recordServicePingResourceEvents string
)

type RecordServicePingResourceEvents struct {
	dbClient *database.DB
}

func (mig *RecordServicePingResourceEvents) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, recordServicePingResourceEvents)
	return err
}

func (mig *RecordServicePingResourceEvents) String() string {
	return "66_record_service_ping_resource_events"
}
