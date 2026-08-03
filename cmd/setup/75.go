package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 75.sql
	addSecurityPolicyClientIDMetadataDocument string
)

type SecurityPolicies3AddClientIDMetadataDocument struct {
	dbClient *database.DB
}

func (mig *SecurityPolicies3AddClientIDMetadataDocument) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addSecurityPolicyClientIDMetadataDocument)
	return err
}

func (mig *SecurityPolicies3AddClientIDMetadataDocument) String() string {
	return "75_security_policies3_add_client_id_metadata_document"
}
