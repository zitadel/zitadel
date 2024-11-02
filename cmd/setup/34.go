package setup

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 34/cockroach/34_cache_schema.sql
	addCacheSchemaCockroach string
	//go:embed 34/postgres/34_cache_schema.sql
	addCacheSchemaPostgres string
)

type AddCacheSchema struct {
	dbClient *database.DB
}

func (mig *AddCacheSchema) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	switch mig.dbClient.Type() {
	case "cockroach":
		_, err = mig.dbClient.ExecContext(ctx, addCacheSchemaCockroach)
	case "postgres":
		_, err = mig.dbClient.ExecContext(ctx, addCacheSchemaPostgres)
	default:
		err = fmt.Errorf("add cache schema: unsupported db type %q", mig.dbClient.Type())
	}
	return err
}

func (mig *AddCacheSchema) String() string {
	return "34_add_cache_schema"
}
