package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 61.sql
	addSAMLSignatureAlgorithm string
)

type IDPTemplate6SAMLSignatureAlgorithm struct {
	dbClient *database.DB
}

func (mig *IDPTemplate6SAMLSignatureAlgorithm) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addSAMLSignatureAlgorithm)
	return err
}

func (mig *IDPTemplate6SAMLSignatureAlgorithm) String() string {
	return "61_idp_templates6_add_saml_signature_algorithm"
}
