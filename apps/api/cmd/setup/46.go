package setup

import (
	"context"
	"embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type InitPermissionFunctions struct {
	eventstoreClient *database.DB
}

var (
	//go:embed 46/*.sql
	permissionFunctions embed.FS
)

func (mig *InitPermissionFunctions) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := readStatements(permissionFunctions, "46")
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

func (*InitPermissionFunctions) String() string {
	return "46_init_permission_functions"
}
