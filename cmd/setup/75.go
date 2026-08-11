package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 75.sql
	addMinimalIntrospection string
)

type AppsConfigsAddMinimalIntrospection struct {
	dbClient *database.DB
}

func (mig *AppsConfigsAddMinimalIntrospection) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addMinimalIntrospection)
	return err
}

func (mig *AppsConfigsAddMinimalIntrospection) String() string {
	return "75_apps_configs_add_minimal_introspection"
}
