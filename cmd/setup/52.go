package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 52.sql
	renameTableIfNotExisting string
)

type IDPTemplate6LDAP2 struct {
	dbClient *database.DB
}

func (mig *IDPTemplate6LDAP2) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, renameTableIfNotExisting)
	return err
}

func (mig *IDPTemplate6LDAP2) String() string {
	return "52_idp_templates6_ldap2"
}
