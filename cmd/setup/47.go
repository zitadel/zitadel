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
	//go:embed 47/cockroach/*.sql
	//go:embed 47/postgres/*.sql
	transactionalInstanceDomainTable embed.FS
)

type TransactionalInstanceDomainTable struct {
	dbClient *database.DB
}

func (mig *TransactionalInstanceDomainTable) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	statements, err := readStatements(transactionalInstanceDomainTable, "47", mig.dbClient.Type())
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

func (mig *TransactionalInstanceDomainTable) String() string {
	return "47_transactional_instance_domain_table2"
}
