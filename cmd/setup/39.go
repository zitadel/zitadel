package setup

import (
	"context"
	"embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	es_v3 "github.com/zitadel/zitadel/internal/eventstore/v3"
)

var (
	//go:embed 39/cockroach/*.sql
	//go:embed 39/postgres/*.sql
	initPushFunc embed.FS
)

type InitPushFunc struct {
	dbClient *database.DB
}

func (mig *InitPushFunc) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	statements, err := readStatements(initPushFunc, "39", mig.dbClient.Type())
	if err != nil {
		return err
	}
	conn, err := mig.dbClient.Conn(ctx)
	if err != nil {
		return err
	}

	for _, stmt := range statements {
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		if _, err := conn.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}
	err = es_v3.CheckExecutionPlan(ctx, conn)
	logging.OnError(err).Debug("unable to register eventstore types")

	return nil
}

func (mig *InitPushFunc) String() string {
	return "39_init_push_func"
}
