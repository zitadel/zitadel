package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 50.sql
	addUsePKCE string
)

type IDPTemplate6UsePKCE struct {
	dbClient *database.DB
}

func (mig *IDPTemplate6UsePKCE) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addUsePKCE)
	return err
}

func (mig *IDPTemplate6UsePKCE) String() string {
	return "50_idp_templates6_add_use_pkce"
}
