package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 70.sql
	idpTemplate6OIDCBranding string
)

type IDPTemplate6OIDCBranding struct {
	dbClient *database.DB
}

func (mig *IDPTemplate6OIDCBranding) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, idpTemplate6OIDCBranding)
	return err
}

func (mig *IDPTemplate6OIDCBranding) String() string {
	return "70_idp_templates6_oidc_add_icon_and_background_color"
}
