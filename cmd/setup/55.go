package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 55.sql
	addSAMLFederatedLogout string
)

type IDPTemplate6SAMLFederatedLogout struct {
	dbClient *database.DB
}

func (mig *IDPTemplate6SAMLFederatedLogout) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addSAMLFederatedLogout)
	return err
}

func (mig *IDPTemplate6SAMLFederatedLogout) String() string {
	return "55_idp_templates6_add_saml_federated_logout"
}
