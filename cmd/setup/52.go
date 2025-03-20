package setup

import (
	"context"
	"embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type InitPermittedOrgsFunction52 struct {
	eventstoreClient *database.DB
}

//go:embed 52/*.sql
var permittedOrgsFunction52 embed.FS

func (mig *InitPermittedOrgsFunction52) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := readStatements(permittedOrgsFunction52, "52", "")
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

func (*InitPermittedOrgsFunction52) String() string {
	return "52_init_permitted_orgs_function"
}
