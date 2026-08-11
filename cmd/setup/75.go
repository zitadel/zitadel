package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 75.sql
	addIDPAuthorizationParameters string
)

type IDPTemplate6AddAuthorizationParameters struct {
	dbClient *database.DB
}

func (mig *IDPTemplate6AddAuthorizationParameters) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addIDPAuthorizationParameters)
	return err
}

func (mig *IDPTemplate6AddAuthorizationParameters) String() string {
	return "75_idp_templates6_add_authorization_parameters"
}
