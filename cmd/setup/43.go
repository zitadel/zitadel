package setup

import (
	"context"
	"embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 43/*.sql
	createFieldsDomainIndex embed.FS
)

type CreateFieldsDomainIndex struct {
	dbClient *database.DB
}

func (mig *CreateFieldsDomainIndex) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := readStatements(createFieldsDomainIndex, "43")
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

func (mig *CreateFieldsDomainIndex) String() string {
	return "43_create_fields_domain_index"
}
