package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 28.sql
	addOIDCCodeChallengeParams string
)

type IDPTemplate6OIDCCodeChallengeParams struct {
	dbClient *database.DB
}

func (mig *IDPTemplate6OIDCCodeChallengeParams) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addOIDCCodeChallengeParams)
	return err
}

func (mig *IDPTemplate6OIDCCodeChallengeParams) String() string {
	return "28_idp_templates6_add_oidc_use_pkce"
}
