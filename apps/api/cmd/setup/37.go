package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 37.sql
	addBackChannelLogoutURI string
)

type Apps7OIDConfigsBackChannelLogoutURI struct {
	dbClient *database.DB
}

func (mig *Apps7OIDConfigsBackChannelLogoutURI) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addBackChannelLogoutURI)
	return err
}

func (mig *Apps7OIDConfigsBackChannelLogoutURI) String() string {
	return "37_apps7_oidc_configs_add_back_channel_logout_uri"
}
