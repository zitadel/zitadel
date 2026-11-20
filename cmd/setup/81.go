package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 81.sql
	addOIDCTokenEndpointAuthMethod string
)

type IDPTemplate6OIDCAddClientAuthMethod struct {
	dbClient *database.DB
}

func (mig *IDPTemplate6OIDCAddClientAuthMethod) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addOIDCTokenEndpointAuthMethod)
	return err
}

func (mig *IDPTemplate6OIDCAddClientAuthMethod) String() string {
	return "81_idp_templates6_oidc_add_token_endpoint_auth_method"
}
