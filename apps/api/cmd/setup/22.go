package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 22.sql
	activeInstanceEvents string
)

type ActiveInstanceEvents struct {
	dbClient *database.DB
}

func (mig *ActiveInstanceEvents) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, activeInstanceEvents)
	return err
}

func (mig *ActiveInstanceEvents) String() string {
	return "22_active_instance_events_index"
}
