package setup

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 61/*.sql
	createDomainsTable embed.FS
)

type CreateDomainsTable struct {
	dbClient *database.DB
}

func (mig *CreateDomainsTable) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := readStatements(createDomainsTable, "61")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		if _, err := mig.dbClient.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}
	return nil
}

func (mig *CreateDomainsTable) String() string {
	return "61_create_domains_table"
}