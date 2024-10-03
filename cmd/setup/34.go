package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 34.sql
	addCacheSchema string
)

type AddCacheSchema struct {
	dbClient *database.DB
}

func (mig *AddCacheSchema) Execute(ctx context.Context, _ eventstore.Event) error {
	if mig.dbClient.Type() == "postgres" {
		_, err := mig.dbClient.ExecContext(ctx, addCacheSchema)
		return err
	}
	logging.WithFields("name", mig.String(), "dialect", mig.dbClient.Type()).Info("unlogged tables not supported")
	return nil
}

func (mig *AddCacheSchema) String() string {
	return "34_add_cache_schema"
}
