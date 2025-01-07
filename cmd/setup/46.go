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
	//go:embed 46/cockroach/*.sql
	//go:embed 46/postgres/*.sql
	transactionalInstanceTable embed.FS
)

type TransactionalInstanceTable struct {
	dbClient *database.DB
}

func (mig *TransactionalInstanceTable) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	statements, err := readStatements(transactionalInstanceTable, "46", mig.dbClient.Type())
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

func (mig *TransactionalInstanceTable) String() string {
	return "46_transactional_instance_table"
}
