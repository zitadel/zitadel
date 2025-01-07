package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 43.sql
	createFieldsDomainIndex string
)

type CreateFieldsDomainIndex struct {
	dbClient *database.DB
}

func (mig *CreateFieldsDomainIndex) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, createFieldsDomainIndex)
	return err
}

func (mig *CreateFieldsDomainIndex) String() string {
	return "43_create_fields_domain_index"
}
