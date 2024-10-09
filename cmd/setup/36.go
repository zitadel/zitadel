package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 36.sql
	addEventstoreSnapshots string
)

type AddEventstoreSnapshots struct {
	dbClient *database.DB
}

func (mig *AddEventstoreSnapshots) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addEventstoreSnapshots)
	return err
}

func (mig *AddEventstoreSnapshots) String() string {
	return "36_add_eventstore_snapshots"
}
