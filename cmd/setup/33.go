package setup

import (
	"context"
	"embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 33/cockroach/33.sql
	//go:embed 33/postgres/33.sql
	eventstoreIndexes embed.FS
)

type EventstoreIndexes struct {
	dbClient *database.DB
}

func (mig *EventstoreIndexes) Execute(ctx context.Context, _ eventstore.Event) error {
	stmt, err := readStmt(eventstoreIndexes, "33", mig.dbClient.Type(), "33.sql")
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *EventstoreIndexes) String() string {
	return "33_eventstore_indexes"
}
