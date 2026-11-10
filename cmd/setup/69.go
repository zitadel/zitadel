package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 69.sql
	cacheTablesLogged string
)

type CacheTablesLogged struct {
	dbClient *database.DB
}

func (mig *CacheTablesLogged) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, cacheTablesLogged)
	return err
}

func (mig *CacheTablesLogged) String() string {
	return "69_cache_tables_logged"
}
