package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 29.sql
	addOAuth2CodeChallengeParams string
)

type IDPTemplate6OAuth2CodeChallengeParams struct {
	dbClient *database.DB
}

func (mig *IDPTemplate6OAuth2CodeChallengeParams) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addOAuth2CodeChallengeParams)
	return err
}

func (mig *IDPTemplate6OAuth2CodeChallengeParams) String() string {
	return "29_idp_templates6_add_oauth2_use_pkce"
}
