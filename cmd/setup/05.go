package setup

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 05.sql
	lastFailedStmts string
)

type LastFailed struct {
	dbClient *sql.DB
}

func (mig *LastFailed) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, lastFailedStmts)
	return err
}

func (mig *LastFailed) String() string {
	return "05_last_failed"
}
