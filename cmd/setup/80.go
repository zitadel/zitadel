package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

//go:embed 80_saml_metadata_url.sql
var samlMetadataURLColumn string

type SAMLMetadataURL struct {
	dbClient *database.DB
}

func (mig *SAMLMetadataURL) Execute(ctx context.Context, _ eventstore.Event) error {
	logging.Info(ctx, "add SAML metadata URL column", "migration", mig.String())
	_, err := mig.dbClient.ExecContext(ctx, samlMetadataURLColumn)
	return err
}

func (mig *SAMLMetadataURL) String() string {
	return "80_saml_metadata_url"
}
