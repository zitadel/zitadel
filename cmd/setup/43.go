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
	//go:embed 43/01_table_definition.sql
	createOutboxTable string

	//go:embed 43/cockroach/*.sql
	//go:embed 43/postgres/*.sql
	createOutboxTriggers embed.FS
)

type CreateOutbox struct {
	dbClient *database.DB
}

func (mig *CreateOutbox) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, createOutboxTable)
	if err != nil {
		return err
	}
	statements, err := readStatements(createOutboxTriggers, "43", mig.dbClient.Type())
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

func (mig *CreateOutbox) String() string {
	return "43_create_outbox"
}
