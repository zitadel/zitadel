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
	//go:embed 40/cockroach/*.sql
	//go:embed 40/postgres/*.sql
	initPushFunc embed.FS
)

type InitPushFunc struct {
	dbClient *database.DB
}

func (mig *InitPushFunc) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	statements, err := readStatements(initPushFunc, "40", mig.dbClient.Type())
	if err != nil {
		return err
	}
	conn, err := mig.dbClient.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := conn.Close()
		logging.OnError(closeErr).Debug("failed to release connection")
		// Force the pool to reopen connections to apply the new types
		mig.dbClient.Pool.Reset()
	}()

	for _, stmt := range statements {
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		if _, err := conn.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}

	return nil
}

func (mig *InitPushFunc) String() string {
	return "40_init_push_func"
}
