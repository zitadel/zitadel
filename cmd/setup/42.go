package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 42.sql
	addOIDCAppLoginVersion string
)

type Apps7OIDCConfigsLoginVersion struct {
	dbClient *database.DB
}

func (mig *Apps7OIDCConfigsLoginVersion) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addOIDCAppLoginVersion)
	return err
}

func (mig *Apps7OIDCConfigsLoginVersion) String() string {
	return "40_apps7_oidc_configs_login_version"
}
