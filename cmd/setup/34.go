package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 34/34_cache_schema.sql
	addCacheSchema string
)

type AddCacheSchema struct {
	dbClient *database.DB
}

func (mig *AddCacheSchema) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	_, err = mig.dbClient.ExecContext(ctx, addCacheSchema)
	return err
}

func (mig *AddCacheSchema) String() string {
	return "34_add_cache_schema"
}
