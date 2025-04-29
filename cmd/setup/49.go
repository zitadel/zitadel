package setup

import (
	"context"
	"embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type InitPermittedOrgsFunction struct {
	eventstoreClient *database.DB
}

var (
	//go:embed 49/*.sql
	permittedOrgsFunction embed.FS
)

func (mig *InitPermittedOrgsFunction) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := readStatements(permittedOrgsFunction, "49")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		if _, err := mig.eventstoreClient.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}
	return nil
}

func (*InitPermittedOrgsFunction) String() string {
	return "49_init_permitted_orgs_function"
}
