package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 59.sql
	addHostedLoginTranslations string
)

type HostedLoginTranslation struct {
	dbClient *database.DB
}

func (mig *HostedLoginTranslation) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addHostedLoginTranslations)
	return err
}

func (mig *HostedLoginTranslation) String() string {
	return "59_hosted_login_translations"
}
