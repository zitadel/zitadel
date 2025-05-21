package setup

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/queue"
)

type RiverMigrateRepeatable struct {
	client *database.DB
}

func (mig *RiverMigrateRepeatable) Execute(ctx context.Context, _ eventstore.Event) error {
	return queue.NewMigrator(mig.client).Execute(ctx)
}

func (mig *RiverMigrateRepeatable) String() string {
	return "repeatable_migrate_river"
}

func (f *RiverMigrateRepeatable) Check(lastRun map[string]interface{}) bool {
	return true
}
