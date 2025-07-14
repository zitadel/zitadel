package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 61.sql
	createDomainsTable string
)

type CreateDomainsTable struct {
	dbClient *database.DB
}

func (mig *CreateDomainsTable) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, createDomainsTable)
	return err
}

func (mig *CreateDomainsTable) String() string {
	return "61_create_domains_table"
}