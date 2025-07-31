package setup

import (
	"context"
	"embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type InitPermittedOrgsFunction53 struct {
	dbClient *database.DB
}

//go:embed 53/*.sql
var permittedOrgsFunction53 embed.FS

func (mig *InitPermittedOrgsFunction53) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := readStatements(permittedOrgsFunction53, "53")
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

func (*InitPermittedOrgsFunction53) String() string {
	return "53_init_permitted_orgs_function_v2"
}
