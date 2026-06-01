package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 57.sql
	createResourceCounts string
)

type CreateResourceCounts struct {
	dbClient *database.DB
}

func (mig *CreateResourceCounts) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, createResourceCounts)
	return err
}

func (mig *CreateResourceCounts) String() string {
	return "57_create_resource_counts"
}
