package setup

import (
	"context"
	"embed"
	_ "embed"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 45/cockroach/*.sql
	//go:embed 45/postgres/*.sql
	eventQueue embed.FS
)

type EventQueue struct {
	dbClient *database.DB
}

func (mig *EventQueue) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	statements, err := readStatements(eventQueue, "45", mig.dbClient.Type())
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		_, err = mig.dbClient.ExecContext(ctx, stmt.query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mig *EventQueue) String() string {
	return "45_event_queue"
}
