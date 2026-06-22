package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 72.sql
	addOIDCConfigRegistrationToken string
)

type Apps7OIDCConfigsAddRegistrationToken struct {
	dbClient *database.DB
}

func (mig *Apps7OIDCConfigsAddRegistrationToken) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addOIDCConfigRegistrationToken)
	return err
}

func (mig *Apps7OIDCConfigsAddRegistrationToken) String() string {
	return "72_apps7_oidc_configs_add_registration_token"
}
