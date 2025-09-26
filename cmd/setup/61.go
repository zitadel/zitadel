package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 61.sql
	addAllowedScopePrefixes string
)

type Apps7OIDCConfigsAddAllowedScopePrefixes struct {
	dbClient *database.DB
}

func (mig *Apps7OIDCConfigsAddAllowedScopePrefixes) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addAllowedScopePrefixes)
	return err
}

func (mig *Apps7OIDCConfigsAddAllowedScopePrefixes) String() string {
	return "61_apps7_oidc_configs_add_allowed_scope_prefixes"
}
